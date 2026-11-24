package main

import "fmt"

func mergeSort(arr []int) []int {
	if len(arr) == 0 {
		return []int{}
	}
	if len(arr) == 1 {
		return arr
	}

	// Divide array into top and bottom half
	mid := (len(arr) / 2)

	left := arr[:mid]
	right := arr[mid:]

	// Sort left and right
	leftSorted := mergeSort(left)
	rightSorted := mergeSort(right)

	// Combine sorted left and right
	combinedSorted := []int{}
	leftIter := 0
	rightIter := 0

	appendRest := func(rest []int, start int) {
		for i := start; i < len(rest); i++ {
			combinedSorted = append(combinedSorted, rest[i])
		}
	}

	for {
		// Left iter over range
		if leftIter >= len(leftSorted) {
			appendRest(rightSorted, rightIter)
			break
		}

		// Right iter over range
		if rightIter >= len(rightSorted) {
			appendRest(leftSorted, leftIter)
			break
		}

		// Iters both still in range
		leftVal := leftSorted[leftIter]
		rightVal := rightSorted[rightIter]
		if leftVal <= rightVal {
			combinedSorted = append(combinedSorted, leftVal)
			leftIter++
		} else {
			combinedSorted = append(combinedSorted, rightVal)
			rightIter++
		}

	}

	fmt.Println("combined", combinedSorted)
	return combinedSorted
}
