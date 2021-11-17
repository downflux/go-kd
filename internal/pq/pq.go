// Package pq implements a priority queue.
package pq

import (
	"container/heap"
	"math"

	"github.com/downflux/go-kd/internal/node"
)

type d struct {
	n *node.N

	index    int
	priority float64
}

// max implements a max heap by fulilling heap.Interface.
//
// See https://pkg.go.dev/container/heap for the implementation from which this
// was shamelessly stolen.
type max struct {
	data []*d
}

func (h *max) Len() int           { return len(h.data) }
func (h *max) Less(i, j int) bool { return h.data[i].priority > h.data[j].priority }
func (h *max) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.data[i].index = i
	h.data[j].index = j
}

func (h *max) Push(x interface{}) {
	i := x.(*d)
	i.index = len(h.data)
	h.data = append(h.data, i)
}

func (h *max) Pop() interface{} {
	i := h.data[len(h.data)-1]

	h.data[len(h.data)-1] = nil
	h.data = h.data[:len(h.data)-1]

	return i
}

// Q represents a priority queue of K-D tree nodes with a set size. Q tracks
// some specified number of nodes with the smallest given priorites --
// attempting to add a node with a larger priority will result in an effective
// no-op.
type Q struct {
	h    *max
	size int
}

func New(size int) *Q {
	q := &Q{
		h: &max{
			data: make([]*d, 0, size),
		},
		size: size,
	}
	heap.Init(q.h)
	return q

}

// Empty checks if the queue contains any elements.
func (q *Q) Empty() bool { return q.h.Len() == 0 }

// Full checks if the queue is at capacity, and if it may reject future points.
func (q *Q) Full() bool { return q.h.Len() >= q.size }

// Priority calculates the current highest priority of queue.
func (q *Q) Priority() float64 {
	if q.Empty() {
		return math.Inf(0)
	}
	// See https://groups.google.com/g/golang-nuts/c/sy1p8SfyPoY.
	return q.h.data[0].priority
}

// Push adds a new node into the queue with the given priority.
//
// The queue will enforce the struct size constraint by removing elements from
// itself until the constraint is satisfied.
func (q *Q) Push(n *node.N, priority float64) {
	heap.Push(q.h, &d{
		n:        n,
		priority: priority,
	})
	for !q.Empty() && q.h.Len() > q.size {
		q.Pop()
	}
}

// Pop removes the node with highest priority from the queue.
func (q *Q) Pop() *node.N { return heap.Pop(q.h).(*d).n }
