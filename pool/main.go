package main

import (
	"fmt"
	"sync"
)

func main() {
	count := 0
	myPool := &sync.Pool{
		New: func() any {
			count++
			fmt.Printf("Creating new instance %d \n", count)
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()

	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}
	// Seed the pool with 4KB
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	
	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
			// Assume something interesting, but quick is being done with
			// this memory.
		}()
	}
	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)

}
