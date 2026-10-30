package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_canPlaceFlowers(t *testing.T) {

	cases := []struct {
		name       string
		flowerbed  []int
		newFlowers int
		expected   bool
	}{
		{
			"1",
			[]int{0, 0, 0},
			1,
			true,
		},
		{
			"2",
			[]int{0, 0, 0},
			2,
			true,
		},
		{
			"3",
			[]int{1, 0, 0, 1},
			1,
			false,
		},
		{
			"4",
			[]int{1, 0, 0, 0, 1},
			1,
			true,
		},
		{
			"5",
			[]int{0},
			1,
			true,
		},
		{
			"6",
			[]int{0, 0},
			1,
			true,
		},
		{
			"7",
			[]int{0, 1},
			1,
			false,
		},
		{
			"8",
			[]int{1},
			0,
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out := canPlaceFlowers(tt.flowerbed, tt.newFlowers)
			assert.Equal(t, tt.expected, out)
		})
	}

}
