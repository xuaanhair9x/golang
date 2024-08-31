package main


import "fmt"

// Lưu trữ dữ liệu theo thứ tự FIFO (First-In-First-Out)

func main() {
	bufferedChan := make(chan int, 5)

	for i := 1; i <= 5; i++ {
		bufferedChan <- i
	}

	for i := 1; i <= 5; i++ {
		fmt.Println(<-bufferedChan)
	}
}