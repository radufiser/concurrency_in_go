package main

import "fmt"

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

	done := make(chan interface{})
	defer close(done)
	myChan := make(chan interface{})
	go func() {
		myChan <- "cici"
		close(myChan)
	}()
	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
}
