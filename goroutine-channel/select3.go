package main

import (
    "fmt"
	"time"
	"math/rand"
)

func goroutineone(endpoint <-chan int)  {
	for {
		select {
		case rep := <-endpoint:
			fmt.Println("--- goroutineone ---: ", rep)

		case <-time.After(3*time.Second):
			fmt.Println("We already waited for 3 hours !")
			return
		}
	}
}

func goroutinetwo(endpoint <-chan int)  {
	for {
		select {
		case rep := <-endpoint:
			fmt.Println("--- goroutinetwo ---: ", rep)

		case <-time.After(4 * time.Second):
			fmt.Println("We already waited for 4 hours !")
			return
		}
	}
}

func putvalueendpoint() <-chan int{
	value := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			fmt.Println("Count: ", i+1)
			value <- rand.Intn(100000)
		}
		//close(value)
	}()
	
	return value
}

func main() {
	value := putvalueendpoint()
	//go goroutineone(value)
	temp(value)
	time.Sleep(6 * time.Second)
	fmt.Println("done")
}
func temp(value <-chan int)  {
	for i := 0 ; i < 5; i ++ {
		select {
		case rep := <-value:
			fmt.Println("--- goroutineone ---: ", rep)

		case <-time.After(3*time.Second):
			fmt.Println("We already waited for 3 hours !")
			return
		}
	}
	go goroutinetwo(value)
}