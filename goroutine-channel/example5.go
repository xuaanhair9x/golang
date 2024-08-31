package main

import "fmt"

var c = make(chan int)
var a string

func f() {
    counter := 10;
    for {
        counter ++;
        if counter == 10000000000 {
            break;
        }
    }
    a = "hello, world"
    x := <- c
    fmt.Println(x)
}

func main() {
    go f()
    c <- 0
    print(a)
}