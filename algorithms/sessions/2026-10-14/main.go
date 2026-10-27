package main

import (
	"fmt"
	"slices"
)

func main() {
	g := [][]int{
		{1, 2},
		{3},
		{3},
		{1},
		{},
	}

	a := reachable(g, 0)
	if !slices.Equal(a, []int{0, 1, 2, 3}) {
		fmt.Println("test failed")
	} else {
		fmt.Println("test succeeded")
	}
}

// Given a graph, return the list of nodes that are reachable from the start
func reachable(graph [][]int, start int) []int {
	q := []int{}
	q = append(q, start)
	visited := make(map[int]struct{})
	visited[start] = struct{}{}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		nodes := graph[curr]

		for _, node := range nodes {
			if _, ok := visited[node]; ok {
				continue
			}
			visited[node] = struct{}{}
			q = append(q, node)
		}
	}

	visitedArr := []int{}
	for k := range visited {
		visitedArr = append(visitedArr, k)
	}

	slices.Sort(visitedArr)

	return visitedArr
}
