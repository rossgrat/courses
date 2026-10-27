package main

import (
	"fmt"
)

func main() {
	q := Queue{}

	q.Enqueue(1)
	q.Enqueue(2)
	v, t := q.Dequeue()
	fmt.Println(v, t)
	v, t = q.Dequeue()
	fmt.Println(v, t)
	v, t = q.Dequeue()
	fmt.Println(v, t)
}

type Queue struct {
	s []int
}

func (q *Queue) Enqueue(v int) {
	q.s = append(q.s, v)
}

func (q *Queue) Dequeue() (int, bool) {
	if len(q.s) == 0 {
		return 0, false
	}
	v := q.s[0]
	q.s = q.s[1:]

	return v, true
}
