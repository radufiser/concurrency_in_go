package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	dataRace()
	dataRace_sync()
}

func dataRace() {
	var data int
	go func() { data++ }() // critical section
	time.Sleep(1 * time.Second)
	if data == 0 { // critical section
		fmt.Printf("the value is %v.\n", data) // critical section
	}
}

// It doesn't solve the logical correctness and data race. 
// we haven't solved the race condition.
// we still don't know deterministically what will run first 
// performance ramifications: Lock() make our programs slow
func dataRace_sync() {
	var memoryAccess sync.Mutex
	var data int
	go func() {
		memoryAccess.Lock()
		data++
		memoryAccess.Unlock()
	}()
	memoryAccess.Lock()
	if data == 0 {
		fmt.Printf("the value is %v.\n", data) // critical section
	}
	memoryAccess.Unlock()
}
