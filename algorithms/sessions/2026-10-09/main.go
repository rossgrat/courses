package main

import "fmt"

func main() {
	// a := []int{1, 2, 3, 1, 2, 3}
	// out := countLongestSubarray(a)
	// fmt.Println(a, out)

	assertEqual(binarySearch([]int{}, 5), -1)
	assertEqual(binarySearch([]int{5}, 5), 0)
	assertEqual(binarySearch([]int{5, 7}, 5), 0)
	assertEqual(binarySearch([]int{4, 5}, 5), 1)
	assertEqual(binarySearch([]int{7}, 5), -1)
}

func assertEqual(a any, b any) {
	if a != b {
		fmt.Println("assertion failed")
	}
}

func binarySearch(a []int, target int) int {
	lo := 0
	hi := len(a) - 1

	for lo <= hi {
		mid := lo + (hi-lo)/2

		m := a[mid]

		if m == target {
			return mid
		}
		if m > target {
			hi = mid - 1
		}
		if m < target {
			lo = mid + 1
		}
	}
	return -1
}

func countLongestSubarray(a []int) int {
	current := 1
	count := 0
	longest := 0

	check := func() {
		switch {
		case current > longest:
			longest = current
			count = 1
		case current == longest:
			count++
		case current < longest:
			// noop
		}
		current = 1
	}

	for i, val := range a {
		if i+1 >= len(a) {
			check()
			break
		}
		switch {
		case a[i+1] > val:
			current++

		case a[i+1] < val:
			check()

		case a[i+1] == val:
			check()
		}
	}

	return count
}
