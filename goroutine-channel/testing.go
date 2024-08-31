package main

import (
	//"fmt"
)

// Client request to replicate data
type ClientRequest struct {
    value string
    reply chan int
}
type Propose struct {
	block string
	client chan ClientRequest
}

func ()  {
	
}

func main()  {
	pro := &Propose {}
	replyChannel := make(chan int)
	request := ClientRequest{"123", replyChannel}
	pro.client <- request
}