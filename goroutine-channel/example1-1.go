package main

import (
	"fmt"
)

func main() {
	// The code is blocked until something gets pushed into the returned channel
	// As opposed to the previous method, we block in the main function, instead
	// of the function itself
	i := <-getNumberChan()
	fmt.Println(i)
}

// return an integer channel instead of an integer
func getNumberChan() <-chan int {
	// create the channel
	c := make(chan int)
	go func() {
		// push the result into the channel
		c <- 5
	}()
	// immediately return the channel
	return c
}