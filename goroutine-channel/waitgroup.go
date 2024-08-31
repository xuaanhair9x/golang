package main

import (
	"fmt"
	//"math/rand"
	"sync"
	//"time"
)

func main() {
	//r := rand.New(rand.NewSource(time.Now().Unix()))
	wc := new(sync.WaitGroup)
	wc.Add(2)

	go func() {
		for i := 0; i < 200; i++ {
			fmt.Printf("11 ")
		}
		wc.Done()
	}()

	go func() {
		for i := 0; i < 200; i++ {
			fmt.Printf("22 ")
		}
		wc.Done()
	}()

	wc.Wait()
	fmt.Println("All Goroutines done")
}