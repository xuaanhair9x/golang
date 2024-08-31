package main

import (
	"fmt"
)
//https://www.geeksforgeeks.org/find-substrings-contain-vowels/

/*
 * Complete the 'vowelsubstring' function below.
 *
 * The function is expected to return a LONG_INTEGER.
 * The function accepts STRING s as parameter.
 */
func isVowel(x rune) bool {
	if x == 'a' || x == 'e' || x == 'i' || x == 'o' || x == 'u' {
		return true
	}
	fmt.Printf("%c \n", x)

	return false
}
func vowelsubstring(s string) int64 {
	var result int64 = 0
	for i, _ := range s {
		fmt.Println("dasndjahsdjas")
		data := make(map[rune]int)
		ChildLoop:
		for j := i; j < len(s); j++ {
			if !isVowel([]rune(s)[j]) {
				break ChildLoop
			}
			data[[]rune(s)[j]] = 1
			if len(data) == 5 {
				result = result + 1
			}
		}

	}
	return result
}

func main() {

	result := vowelsubstring("aaeiouxa")

	fmt.Println( "result:", result)

}
