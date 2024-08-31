// Go program to illustrate how to
// find the capacity of the channel
  
package main
  
import "fmt"
/**
Capacity of the Channel: In channel, you can find the capacity 
of the channel using cap() function. 
Here, the capacity indicates the size of the buffer.

**/


// Main function
func main() {
  
    // Creating a channel
    // Using make() function
    mychnl := make(chan string, 5)
    mychnl <- "GFG"
    mychnl <- "gfg"
    mychnl <- "Geeks"
    mychnl <- "GeeksforGeeks"
  
    // Finding the capacity of the channel
    // Using cap() function
    fmt.Println("Capacity of the channel is: ", cap(mychnl))
}

// Capacity of the channel is:  5