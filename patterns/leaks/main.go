package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() { // will remain in memory for the lifetime of this process
			defer fmt.Println("do work exited")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return completed
	}
	doWork(nil)
	fmt.Println("Done")

	doWork2 := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork2 exited")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)

				case <-done:
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	terminated := doWork2(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("canceling doWork routine ...")
		close(done)

	}()

	<-terminated
	fmt.Println("Done.")

	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream close exited")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}

		}()
		return randStream
	}

	done2 := make(chan interface{})
	randStream := newRandStream(done2)

	fmt.Println("Print 3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	close(done2)
	time.Sleep(2 * time.Second)
}
