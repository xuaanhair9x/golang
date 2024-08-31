package main

import (
	"fmt"
	//"sort"
)
// https://leetcode.com/discuss/interview-question/1362915/amazon-hackerrank-question-priority-assignment
func main() {
	bb := []int{1,3,7,3}
	reassignPriority(bb)
}

func reassignPriority(priorities []int) {

	var tempArr [100]int
	for _, value := range priorities {
		tempArr[value] = 1
	}

	rank := 0;
	for i := 0; i < 100; i ++ {
		if (tempArr[i] > 0) {
			rank ++
			tempArr[i] = rank
		} else {
			tempArr[i] = rank
		}
	}

	for _, i := range priorities {
		fmt.Println(tempArr[i])
	}
}