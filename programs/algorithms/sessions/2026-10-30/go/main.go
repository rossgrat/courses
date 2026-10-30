package main

const (
	EMPTY  = 0
	FLOWER = 1
)

func canPlaceFlowers(flowerbed []int, n int) bool {
	if n == 0 {
		return true
	}
	for i := range flowerbed {
		if flowerbed[i] == FLOWER {
			continue
		}
		if (i > 0) && (flowerbed[i-1] == FLOWER) {
			continue
		}
		if (i < len(flowerbed)-1) && (flowerbed[i+1] == FLOWER) {
			continue
		}
		n--
		flowerbed[i] = FLOWER
		if n == 0 {
			return true
		}
	}

	return false
}
