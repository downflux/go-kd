package pq

import (
	"container/heap"
	"math"

	"github.com/downflux/go-kd/point"
)

type node[T point.P] struct {
	p      T
	index  int
	weight float64
}

// max implements the heap.Interface.
//
// See https://pkg.go.dev/container/heap for the heap implementation from which
// this was shamelessly stolen.
type max[T point.P] []*node[T]

func (heap *max[T]) Len() int           { return len(*heap) }
func (heap *max[T]) Less(i, j int) bool { return (*heap)[i].weight > (*heap)[j].weight }
func (heap *max[T]) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
	(*heap)[i].index = i
	(*heap)[j].index = j
}
func (heap *max[T]) Push(x any) {
	m := x.(*node[T])
	m.index = heap.Len()
	*heap = append(*heap, m)
}
func (heap *max[T]) Pop() any {
	h := *heap
	n := len(h)
	m := h[n-1]

	h[n-1] = nil
	m.index = -1

	*heap = h[0 : n-1]
	return m
}

// PQ is an exported priority queue of k-D tree nodes with a maximum queue size.
// PQ tracks a specified number points with the smallest given priorities --
// attempting to add a point with a larger priority will result in an effective
// no-op.
type PQ[T point.P] struct {
	heap *max[T]
	n    int
}

func New[T point.P](size int) *PQ[T] {
	h := max[T](make([]*node[T], 0, size))
	pq := &PQ[T]{
		heap: &h,
		n:    size,
	}
	heap.Init(pq.heap)
	return pq
}

func (pq *PQ[T]) Len() int    { return pq.heap.Len() }
func (pq *PQ[T]) Empty() bool { return pq.Len() == 0 }
func (pq *PQ[T]) Full() bool  { return pq.Len() >= pq.n }

// Priority calculates the current highest priority of the queue.
func (pq *PQ[T]) Priority() float64 {
	if pq.Empty() {
		return math.Inf(0)
	}

	// See https://groups.google.com/g/golang-nuts/c/sy1p8SfyPoY.
	return (*pq.heap)[0].weight
}

// Push adds a new point into the queue.
//
// The queue will enforce the struct size constraint by removing elements frmo
// itself until the constraint is satisfied.
func (pq *PQ[T]) Push(p T, priority float64) {
	heap.Push(pq.heap, &node[T]{
		p:      p,
		weight: priority,
	})
	for !pq.Empty() && pq.Len() > pq.n {
		pq.Pop()
	}
}

// Pop removes the node with the highest priority from the queue.
func (pq *PQ[T]) Pop() T { return heap.Pop(pq.heap).(*node[T]).p }
