package main

import (
	"fmt"
	"sync"
	"time"
)
/**

Start gotoutine 1
Start gotoutine 2
Process gotoutine 1
0  is even
Process gotoutine 2

**/
func isEven(n int) bool {
	return n%2 == 0
}

func main() {
	n := 2
	var m sync.Mutex

	go func() {
		fmt.Println("Start gotoutine 1")
		m.Lock()
		defer m.Unlock()
		nIsEven := isEven(n)
		time.Sleep(5 * time.Millisecond)
		if nIsEven {
			fmt.Println("Process gotoutine 1")
			fmt.Println(n, " is even")
		} else {
			fmt.Println("Process gotoutine 1")
			fmt.Println(n, "is odd")
		}
	}()

	go func() {
		m.Lock()
		fmt.Println("Start gotoutine 2")
		fmt.Println("Process gotoutine 2")
		n++
		m.Unlock()
	}()

	// just waiting for the goroutines to finish before exiting
	time.Sleep(time.Second)
}
