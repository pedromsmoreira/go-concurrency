package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	batchSize := 50
	tb := New(500*time.Millisecond, batchSize)
	var wg sync.WaitGroup
	batches := make([][]interface{}, 0)
	numberOfMessagesToPublish := 100
	wg.Add(1)
	go func() {
		defer wg.Done()
		for b := range tb.dispatcher {
			fmt.Println(fmt.Sprintf("received batch with size %d", len(b)))
			batches = append(batches, b)
		}
	}()

	for i := 0; i < numberOfMessagesToPublish; i++ {
		tb.Publish(i)
		time.Sleep(50 * time.Millisecond)
	}
	tb.Close()

	wg.Wait()
	fmt.Println("done")
}

type timedBatchManager struct {
	ticker     *time.Ticker
	receiver   chan interface{}
	dispatcher chan []interface{}
	wg         *sync.WaitGroup
	close      sync.Once
	isClosed   bool
	cancel     context.CancelFunc
}

func New(interval time.Duration, batchSize int) *timedBatchManager {
	ticker := time.NewTicker(interval)
	receiver := make(chan interface{})
	dispatcher := make(chan []interface{})
	batch := make([]interface{}, 0, batchSize)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	tb := &timedBatchManager{
		ticker:     ticker,
		receiver:   receiver,
		dispatcher: dispatcher,
		wg:         &wg,
		cancel:     cancel,
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				if len(batch) != 0 {
					dispatcher <- batch
				}
				close(dispatcher)
				return
			case msg := <-receiver:
				batch = append(batch, msg)
				fmt.Println(fmt.Sprintf("batch has %d elements", len(batch)))
				if len(batch) == batchSize {
					dispatcher <- batch
					fmt.Println("resetting batch...")
					batch = make([]interface{}, 0, batchSize)
					break
				}
			case <-ticker.C:
				if len(batch) != 0 {
					dispatcher <- batch
					fmt.Println(fmt.Sprintf("Sent batch with %d", len(batch)))
					batch = make([]interface{}, 0, batchSize)
				}
			}
		}
	}()

	return tb
}

func (tb *timedBatchManager) Publish(message interface{}) {
	if tb.isClosed {
		return
	}
	fmt.Println(fmt.Sprintf("sending %v", message))

	tb.receiver <- message
}

func (tb *timedBatchManager) Close() {
	tb.close.Do(func() {
		tb.isClosed = true
		tb.cancel()
		tb.wg.Wait()
	})
}
