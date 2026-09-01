package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Item is the unit of work exchanged between a producer and a consumer.
type Item struct {
	producer string
	id       int
}

// produce runs forever, pushing items onto queue until ctx is cancelled.
// A send on a full buffered channel blocks — that's the backpressure.
func produce(name string, queue chan<- Item, done <-chan struct{}) {
	id := 0
	for {
		select {
		case <-done:
			return
		case queue <- Item{producer: name, id: id}:
			id++
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		}
	}
}

// consume runs forever, pulling items off queue until ctx is cancelled.
// A receive on an empty channel blocks until a producer sends.
func consume(name string, queue <-chan Item, done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case item := <-queue:
			fmt.Printf("Consumer [%s] consume item [%d] produced by [%s]\n", name, item.id, item.producer)
		}
	}
}
