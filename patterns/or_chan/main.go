package main

import (
	"fmt"
	"time"
)

func main() {

	var or func(channels ...<-chan interface{}) <-chan interface{} 

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		go func ()  {
			defer close(orDone)	
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
			
		}()
		return orDone
	}

	sig := func (after time.Duration) <-chan interface{}  {
		c := make(chan interface{})
		go func ()  {
			time.Sleep(after)
			close(c)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2 * time.Second),
		sig(4 * time.Second),
		sig(8 * time.Second),
		sig(9 * time.Second),
		sig(3 * time.Second),
	)
	fmt.Printf("done after %v", time.Since(start))
	
}