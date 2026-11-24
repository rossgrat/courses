package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Main(t *testing.T) {
	cases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "1,2,3",
			input:    []int{3, 2, 1},
			expected: []int{1, 2, 3},
		},
		{
			name:     "1",
			input:    []int{1},
			expected: []int{1},
		},
		{
			name:     "0",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "1,2,3,4,5,6,7",
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			expected: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:     "9,5,1,4,7,0",
			input:    []int{9, 5, 1, 4, 7, 0},
			expected: []int{0, 1, 4, 5, 7, 9},
		},
		{
			name:     "0,0,0,0,5,4,3,2,1",
			input:    []int{0, 0, 0, 0, 5, 4, 3, 2, 1},
			expected: []int{0, 0, 0, 0, 1, 2, 3, 4, 5},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mergeSort(tt.input))
		})
	}
}
