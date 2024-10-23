package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	salutation := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "ahoj"
	}()
	wg.Wait()

	fmt.Println(salutation) // goroutines execute in the same address space they were created in

	
	for _, salutation := range []string{"1", "2", "3", "4", "5", "6", "7"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()


	for _, salutation := range []string{"hello1", "greetings2", "good day2"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()

}
