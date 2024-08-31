// Go program to illustrate send
// and receive operation
package main
  
import "fmt"
import "time"
  
func myfunc(ch chan int) {
    for v := range ch {
        fmt.Println(234 + v)
    }
}
func main() {
    fmt.Println("start Main method")
    // Creating a channel
    ch := make(chan int)
  
    go myfunc(ch)
    go func() {
        ch <- 23
        ch <- 13
        close(ch)
    }()
    time.Sleep(5 * time.Second)
    fmt.Println("End Main method")
}

/**
Output

start Main method
257
End Main method


**/