package main
  
import "fmt"
/**
Length of the Channel: In channel, you can find the length of the 
channel using len() function. Here, the length indicates the number 
of value queued in the channel buffer.
**/


// Main function
func main() {
  
    // Creating a channel
    // Using make() function
    mychnl := make(chan string, 4)
    mychnl <- "GFG"
    mychnl <- "gfg"
    mychnl <- "Geeks"
    mychnl <- "GeeksforGeeks"
  
    // Finding the length of the channel
    // Using len() function
    fmt.Println("Length of the channel is: ", len(mychnl))
}

// Length of the channel is:  4