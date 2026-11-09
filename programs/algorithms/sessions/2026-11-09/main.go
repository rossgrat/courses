package main

func mergeSort(arr []int) []int {
	if len(arr) == 0 {
		return []int{}
	}
	if len(arr) == 1 {
		return arr
	}

	// Divide array in top and bottom half
	mid := (len(arr) / 2) + 1

	left := arr[:mid]
	right := arr[mid:]

	// Sort left and right
	leftSorted := mergeSort(left)
	rightSorted := mergeSort(right)

	// Combine sorted left and right
	combinedSorted := []int{}
	leftIter := 0
	rightIter := 0

	appendRest := func(rest []int) {
		for _, val := range rest {
			combinedSorted = append(combinedSorted, val)
		}
	}

	for {
		// Left iter over range
		if leftIter >= len(leftSorted) {
			appendRest(rightSorted)
			break
		}

		// Right iter over range
		if rightIter >= len(rightSorted) {
			appendRest(leftSorted)
			break
		}

		// Iters both still in range
		if len(leftSorted) < leftIter && len(rightSorted) < rightIter {
			leftVal := leftSorted[leftIter]
			rightVal := rightSorted[rightIter]
			if leftVal <= rightVal {
				combinedSorted = append(combinedSorted, leftVal)
				leftIter++
			} else {
				combinedSorted = append(combinedSorted, rightVal)
				rightVal++
			}
		}

	}

	return combinedSorted

}
