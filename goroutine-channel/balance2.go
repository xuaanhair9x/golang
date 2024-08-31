package main

import (
	"fmt"
	//"time"
	"sync"
)

func publisher() <-chan int {
	c := make(chan int)
	go func() {
		for i := 1; i <= 500; i++ {
			fmt.Printf("Consumer add task %d\n", i)
			c <- i
		}

		close(c)
	}()
	return c
}

func main() {
	var wg sync.WaitGroup
	myChan := publisher()
	maxConsumer := 8

	wg.Add(maxConsumer)
	doIncrement := func(c <-chan int, name string) {
		counter := 0
		for value := range c {
			fmt.Printf("Consumer %s is doing task %d\n", name, value)
			counter++
			time.Sleep(time.Millisecond)
		}
		fmt.Printf("Consumer %s has finished %d task(s)\n", name, counter)
		wg.Done()
	}

	for i := 1; i <= maxConsumer; i++ {
		go doIncrement(myChan, fmt.Sprintf("%d", i))
	}
	wg.Wait()
}