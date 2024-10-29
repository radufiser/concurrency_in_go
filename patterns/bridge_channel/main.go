package main

import (
	"fmt"
)

func main() {

	orDone := func(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {

		resultStream := make(chan interface{})

		go func() {
			defer close(resultStream)
			for {
				select {
				case <-done:
					fmt.Println("orDone done")
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case resultStream <- v:
					case <-done:
					}
				}
			}

		}()
		return resultStream
	}

	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {

		valStream := make(chan interface{})

		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case <-done:
					return
				case maybeStream, ok := <-chanStream:
					if !ok {
						return
					}
					stream = maybeStream
				}

				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}

			}
		}()

		return valStream
	}

	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 3)
				stream <- i
				stream <- i * i
				stream <- i * i * i
				chanStream <- stream
				close(stream)

			}
		}()

		return chanStream

	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
