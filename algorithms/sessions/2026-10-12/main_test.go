package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		q := Queue{}
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		q.Enqueue(4)
		v, o := q.Dequeue()
		assert.Equal(t, v, 1, "dequeue 1")
		q.Enqueue(10)
		v, o = q.Dequeue()
		assert.Equal(t, v, 2, "dequeue 2")
		v, o = q.Dequeue()
		assert.Equal(t, v, 3, "dequeue 3")
		v, o = q.Dequeue()
		assert.Equal(t, v, 4, "dequeue 4")
		v, o = q.Dequeue()
		assert.Equal(t, v, 10, "dequeue 10")
		v, o = q.Dequeue()
		assert.Equal(t, o, false, "dequeue empty queue")
		q.Enqueue(0)
		q.Enqueue(-1)
		v, o = q.Dequeue()
		assert.Equal(t, v, 0, "dequeue 0")
		v, o = q.Dequeue()
		assert.Equal(t, v, -1, "dequeue -1")
		v, o = q.Dequeue()
		assert.Equal(t, o, false, "dequeue empty queue")
	})
}
