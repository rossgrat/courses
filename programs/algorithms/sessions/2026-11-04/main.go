package main

func carPooling(trips [][]int, capacity int) bool {
	pickup := make(map[int]int)
	dropoff := make(map[int]int)
	// stops?

	for _, val := range trips {
		passengers := val[0]
		start := val[1]
		end := val[2]

		if val, ok := pickup[start]; ok {
			pickup[start] = val + passengers
		} else {
			pickup[start] = passengers
		}

		if val, ok := dropoff[end]; ok {
			dropoff[end] = val + passengers
		} else {
			dropoff[end] = passengers
		}
	}

	passengers := 0
	for i := 0; i <= 1000; i++ {
		if val, ok := dropoff[i]; ok {
			passengers = passengers - val
		}

		if val, ok := pickup[i]; ok {
			if passengers+val > capacity {
				return false
			} else {
				passengers = passengers + val
			}
		}
	}
	return true
}

// Sort trips by start, break ties with end?
//
// There are 1000 stops, so we could stop at every stop?
// Two maps, start to value, end to value?
// Keep a running index of actual stops that need to be made
//
