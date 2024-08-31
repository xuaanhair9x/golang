package main

import (
	"fmt"
	"sync"
	"time"
)
/**

Start gotoutine 1
Start gotoutine 3
Start gotoutine 2
Process gotoutine 1
0  is even
Process gotoutine 3
Process gotoutine 2
1  is positive

**/
func isEven(n int) bool {
	return n%2 == 0
}

func main() {
	n := 0
	var m sync.RWMutex

	// now, both goroutines call m.Lock() before accessing `n`
	// and call m.Unlock once they are done
	go func() {
		fmt.Println("Start gotoutine 1")
		m.RLock()
		defer m.RUnlock()
		nIsEven := isEven(n)
		time.Sleep(5 * time.Millisecond)
		fmt.Println("Process gotoutine 1")
		if nIsEven {
			fmt.Println(n, " is even")
			return
		}
		fmt.Println(n, "is odd")
	}()

	go func() {
		fmt.Println("Start gotoutine 2")
		m.RLock()
		defer m.RUnlock()
		nIsPositive := n > 0
		fmt.Println("Process gotoutine 2")
		time.Sleep(5 * time.Millisecond)
		if nIsPositive {
			fmt.Println(n, " is positive")
			return
		}
		fmt.Println(n, "is not positive")
	}()

	go func() {
		fmt.Println("Start gotoutine 3")
		m.Lock()
		fmt.Println("Process gotoutine 3")
		n++
		m.Unlock()
	}()

	time.Sleep(time.Second)
}
