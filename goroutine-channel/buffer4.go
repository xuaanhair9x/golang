package main

func main() {
	bufferedChan := make(chan int, 1)
	unbufferedChan := make(chan int)

	bufferedChan <- 1   // OK
	unbufferedChan <- 1 // deadlock
}