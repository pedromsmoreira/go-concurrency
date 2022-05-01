package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	batchRoutinesCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_routine_created_total",
		Help: "The total number of batch routine created",
	})
	batchDispatched = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_dispatched_total",
		Help: "The total number of batches dispatched",
	})
)

func main() {

	var wg sync.WaitGroup
	wg.Add(2)
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.Handler())
	srv := http.Server{
		Addr:    ":8000",
		Handler: m,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			return
		}
	}()

	batchSize := 50
	tb := New(500*time.Millisecond, batchSize)
	batches := make([][]interface{}, 0)
	numberOfMessagesToPublish := 100
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

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
		batchRoutinesCreated.Inc()
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				if len(batch) != 0 {
					dispatcher <- batch
					batchDispatched.Inc()
				}
				close(dispatcher)
				return
			case msg := <-receiver:
				batch = append(batch, msg)
				fmt.Println(fmt.Sprintf("batch has %d elements", len(batch)))
				if len(batch) == batchSize {
					dispatcher <- batch
					batchDispatched.Inc()
					fmt.Println("resetting batch...")
					batch = make([]interface{}, 0, batchSize)
					break
				}
			case <-ticker.C:
				if len(batch) != 0 {
					dispatcher <- batch
					batchDispatched.Inc()
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
