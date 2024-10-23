package main

import (
	"fmt"
	"sync"
	"time"
)


func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func (delay time.Duration)  {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Remove from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}

		fmt.Println("Adding to queue")

		queue = append(queue, struct{}{})

		go removeFromQueue(1*time.Second)

		c.L.Unlock()
	}


	type Button struct {
		Clicked *sync.Cond
	}

	button :=  Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func (c *sync.Cond, fn func())  {
		var goroutineRouning sync.WaitGroup
		goroutineRouning.Add(1)

		go func ()  {
			goroutineRouning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()	
		}()
		goroutineRouning.Wait()
	}

	var clickedRegistered sync.WaitGroup 
	clickedRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximising Window")
		clickedRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Display annoying dialog box")
		clickedRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked")
		clickedRegistered.Done()
	})

	button.Clicked.Broadcast()


	clickedRegistered.Wait()

}