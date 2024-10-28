package main

import (
	"fmt"
	"math/rand"
)

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}

	}()
	return valueStream
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		defer fmt.Println("repeatFn exited")
		for {
			select {
			case <-done:
				fmt.Println("Done repeatFn")
				return
			case valueStream <- fn():
			}

		}
	}()
	return valueStream
}

func toString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {

	stringStream := make(chan string)

	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()

	return stringStream
}

func main() {

	done := make(chan interface{})
	defer close(done)

	randomInt := func() interface{} {
		return rand.Int()
	}
	for num := range take(done, repeatFn(done, randomInt), 10) {
		fmt.Println(num)
	}

	done1 := make(chan interface{})
	defer close(done1)
	var message string
	for token := range toString(done1, take(done1, repeat(done1, "I", "am."), 50)) {
		message += token
	}

	fmt.Printf("message: %s...", message)

}
