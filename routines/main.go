package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup 

	sayHellow :=  func ()  {
		defer wg.Done()
		fmt.Println("hi")
	}

	wg.Add(1)
	go sayHellow()
	wg.Wait()
}