package main

import (
    "fmt"
    "time"
)

func main() {

    c1 := make(chan string)
    c2 := make(chan string)
	c3 := make(chan string)
    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
	}()
	
	go func() {
        c3 <- "two"
    }()

    for {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
			fmt.Println("received", msg2)
		case <-c3:
			break
        }
    }
}