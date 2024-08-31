package main

import (
	"fmt"
)

func main() {
	bufferedChan := make(chan int, 5)
	fmt.Printf("BufferChan has len = %d, cap = %d\n", len(bufferedChan), cap(bufferedChan))
	<-bufferedChan // deadlock here
}