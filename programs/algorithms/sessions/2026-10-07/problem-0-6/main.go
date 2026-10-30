package main

import "fmt"

func main() {
	arrays := [][]int{
		{1, 3, 4, 2, 7, 5, 6, 9, 8},
		{4, 1, 2, 3, 1, 2, 3},
		{1, 1, 1},
		{0},
		{3, 2, 1},
		{1, 5, 7, 1, 4, 7, 1, 6, 8, 9, 9, 3, 2},
	}

	for _, a := range arrays {
		answer := countLongSubarrays(a)
		fmt.Println(a, answer)
	}
}

func countLongSubarrays(a []int) int {
	count := 0
	longest := 0
	current := 1

	check := func() {
		if current > longest {
			count = 1
			longest = current
		} else if current == longest {
			count++
		}
		current = 1
	}

	for i, val := range a {
		if i+1 == len(a) {
			check()
			break
		}
		if a[i+1] > val {
			current++
		} else {
			check()
		}
	}

	return count
}
