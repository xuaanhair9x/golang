package main

import (
	"fmt"
	"sync"
)

func streamNumbers(numbers ...int) <-chan int {
	c := make(chan int)

	go func() {
		for n := range numbers {
			fmt.Println("streamNumbers:", n);
			c <- n
		}

		close(c)
	}()

	return c
}

func sumAllStreams(streams ...<-chan int) <-chan int {
	sumChan := make(chan int)
	counter := 0
	wc := new(sync.WaitGroup)

	wc.Add(len(streams))

	for i := 0; i < len(streams); i++ {
		go func(s <-chan int) {
			for n := range s {
				fmt.Println("sumAllStreams:", n);

				counter += n
			}
			wc.Done()
		}(streams[i])
	}

	go func() {
		wc.Wait()
		sumChan <- counter
		fmt.Println("counter -------------- ", counter)

	}()
	fmt.Println("finish ")
	return sumChan
}

func main() {
	fmt.Println("Start 1");
	s := sumAllStreams(
		streamNumbers(1, 2, 3, 4, 5),
		streamNumbers(8, 8, 3, 3, 10, 12, 14),
		streamNumbers(1, 1, 2, 2, 4, 4, 6),
	)
	fmt.Println("Start 2");

	fmt.Println(<-s)
	fmt.Println("Start 3");

}