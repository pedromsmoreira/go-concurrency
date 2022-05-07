package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/pedromsmoreira/go-concurrency/metrics"
	"github.com/pedromsmoreira/go-concurrency/server"
	"os"
	"sync"
)

type User struct {
	Id string `json:"id"`
}

func main() {

	srv := server.New().WithMetrics()
	srv.Start()
	defer srv.Stop()

	var wg sync.WaitGroup
	natsURL := os.Getenv("NATS_URL")
	nc, _ := nats.Connect(natsURL)
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "test",
		Subjects: []string{"test.*"},
	})

	if err != nil {
		fmt.Println("error creating stream")
		os.Exit(1)
	}

	// move to another go routine
	for i := 0; i < 500; i++ {
		j := i
		go func() {
			wg.Add(1)
			defer wg.Done()
			ack, err := js.Publish("test.sub1", []byte(fmt.Sprintf("msg: %d", j)))
			if err != nil {
				metrics.PublishErrorCount.WithLabelValues("test.error").Inc()
			}
			if ack == nil {
				metrics.NackTotalMessages.WithLabelValues("test.sub1").Inc()
			}
		}()
	}

	fmt.Println("waiting to publish all messages")
	wg.Wait()

	_, err = js.Subscribe("test.sub1", func(msg *nats.Msg) {
		data := string(msg.Data)
		fmt.Println(data)
	}, nats.Durable("test_consumer"))

	if err != nil {
		fmt.Println(fmt.Sprintf("error: %v", err.Error()))
		os.Exit(1)
	}
}
