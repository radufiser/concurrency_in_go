package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()
	fmt.Println("Blocking on read")

	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}

	c1 := make(chan chan interface{})
	close(c1)
	c2 := make(chan chan interface{})
	close(c2)

	var c1Count, c2Count int

	for i := 10000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++

		case <-c2:
			c2Count++
		}

	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)

	var ch <-chan int
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out")

	}

	done := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0

loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}
