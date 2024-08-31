package main

import (
	"fmt"
)

func main() {
	bufferedChan := make(chan int, 5)
	bufferedChan <- 1
	bufferedChan <- 2
	bufferedChan <- 3
	bufferedChan <- 4
	bufferedChan <- 5

	fmt.Printf("BufferChan has len = %d, cap = %d\n", len(bufferedChan), cap(bufferedChan))
	// <-bufferedChan fix deadlock
	bufferedChan <- 6 // deadlock here

	

}