package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				fmt.Println("take done")
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
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

func isPrime(v int) bool {
	if v <= 1 {
		return false // 0 and 1 are not prime numbers
	}
	if v <= 3 {
		return true // 2 and 3 are prime numbers
	}
	if v%2 == 0 || v%3 == 0 {
		return false // Eliminates most non-primes early
	}
	for i := 5; i*i <= v; i += 6 { // Only check odd divisors up to sqrt(v)
		if v%i == 0 || v%(i+2) == 0 {
			return false
		}
	}
	return true
}

func isPrimeNaive(v int) bool {
	for i := 2; i < v; i++ {
		if v%i == 0 {
			return false
		}
	}
	return true
}

func primeFinder(done chan interface{}, randIntStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})

	go func() {
		defer close(primeStream)

		for i := range randIntStream {
			select {
			case <-done:
				fmt.Println("primeFinder done")
				return
			default:
				if isPrimeNaive(i) {
					primeStream <- i
				}
			}
		}
	}()
	return primeStream
}

func toInt(done chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)

	go func() {
		defer close(intStream)

		for v := range valueStream {
			select {
			case <-done:
				fmt.Println("toInt done")
				return
			case intStream <- v.(int):
			}
		}

	}()
	return intStream
}

func fanIn(done chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	fanInStream := make(chan interface{})
	var wg sync.WaitGroup

	multiplex := func(channel <-chan interface{}) {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("fanIn done")
				return

			case fanInStream <- <-channel:
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go multiplex(ch)
	}

	go func() {
		wg.Wait()
		close(fanInStream)
	}()

	return fanInStream
}

func main() {
	randomInt := func() interface{} {
		return rand.Intn(500000000)
	}

	done := make(chan interface{})
	defer close(done)

	startTime := time.Now()

	randomIntStream := toInt(done, repeatFn(done, randomInt))

	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randomIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(startTime))

	//fan_in
	fmt.Println("Primes FanIn:")
	start := time.Now()
	numFinders := runtime.NumCPU()
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randomIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))

}
