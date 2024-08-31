package main

import (
	"fmt"
	"sync"
	"time"
)

type intLock struct {
	val int
	sync.Mutex
}

func (n *intLock) isEven() bool {
	return n.val%2 == 0
}

func checkNumber(n *intLock)  {
	fmt.Println("Start gotoutine 1")
	n.Lock()
	defer n.Unlock()
	fmt.Println("Process gotoutine 1")

	go func() {
		for i := 0; i < 450; i++ {
			fmt.Printf(" %d",i)
		}
	}()


	nIsEven := n.isEven()
	if nIsEven {
		fmt.Println(n.val, " is even")
		return
	}
	fmt.Println(n.val, "is odd")
}

func main() {
	n := &intLock{val: 0}

	go checkNumber(n)

	go func() {
		fmt.Println("Start gotoutine 2")
		n.Lock()
		fmt.Println("Process gotoutine 2")
		n.val++
		n.Unlock()
	}()

	time.Sleep(time.Second)
}
