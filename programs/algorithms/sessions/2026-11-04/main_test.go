package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cases := []struct {
		name     string
		trips    [][]int
		capacity int
		expected bool
	}{
		{
			"example 1",
			[][]int{
				{2, 1, 5},
				{3, 3, 7},
			},
			4,
			false,
		},
		{
			"example 2",
			[][]int{
				{2, 1, 5},
				{3, 3, 7},
			},
			5,
			true,
		},
		{
			"dropoff before pickup",
			[][]int{
				{1, 1, 2},
				{1, 2, 3},
			},
			1,
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, carPooling(tt.trips, tt.capacity))
		})
	}
}
