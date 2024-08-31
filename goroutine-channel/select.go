package main

import (
    "fmt"
    "time"
)
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			time.Sleep(time.Second)
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return

		// case <-time.After(3 * time.Second):
		// 	fmt.Println("There's no more time to this. Exiting!")
		// 	return
		}

		fmt.Printf("sss ")
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			//time.Sleep(4*time.Second)
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}