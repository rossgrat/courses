package main

import (
	"testing"
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
	}

	for _, tt :=  range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mergeSort(tt.input))
		}
	}
}
