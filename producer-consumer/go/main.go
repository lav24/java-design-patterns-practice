package main

import (
	"fmt"
	"time"
)

func main() {
	queue := make(chan Item, 5) // bounded buffer, same capacity as the Java LinkedBlockingQueue
	done := make(chan struct{}) // closing this signals every goroutine to stop

	for i := 0; i < 2; i++ {
		go produce(fmt.Sprintf("Producer_%d", i), queue, done)
	}
	for i := 0; i < 3; i++ {
		go consume(fmt.Sprintf("Consumer_%d", i), queue, done)
	}

	time.Sleep(10 * time.Second)
	close(done)
	time.Sleep(100 * time.Millisecond) // give goroutines a moment to notice and exit
}
