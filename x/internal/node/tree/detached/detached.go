package detached

import (
	"fmt"

	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/vector"
)

type O[T point.P] struct {
	Data []T
	K    vector.D
	N    int

	Axis vector.D
}

type N[T point.P] struct {
	data []T

	pivot vector.V
	axis  vector.D
	left  *N[T]
	right *N[T]
}

// New recursively constructs a node object given the input data.
func New[T point.P](o O[T]) *N[T] {
	if o.Axis >= o.K {
		panic(fmt.Sprintf("given node dimension greater than vector dimension: %v > %v", o.Axis, o.K))
	}
	if o.N < 1 {
		panic("given leaf node size must be a positive integer")
	}
	if len(o.Data) <= 0 {
		return nil
	}
	if len(o.Data) <= o.N {
		return &N[T]{
			data: o.Data,
			axis: o.Axis,
		}
	}
	pivot := hoare(o.Data, 0, 0, len(o.Data), func(a vector.V, b vector.V) bool { return a.X(o.Axis) < b.X(o.Axis) })

	node := &N[T]{
		data:  []T{o.Data[pivot]},
		pivot: o.Data[pivot].P(),
		axis:  o.Axis,
	}

	// Node construction can be concurrent since we guarantee child nodes
	// will never access data across the high / low boundary. Note that this
	// does increase the number of allocs ~3x, and is less performant for
	// low data sizes (i.e. for less than ~10k points). We can optionally
	// branch tree construction on the length of incoming data if we find
	// the performance prohibitively slow.
	l := make(chan *N[T])
	r := make(chan *N[T])

	go func(ch chan *N[T]) {
		ch <- New[T](O[T]{
			Data: o.Data[0:pivot],
			Axis: (o.Axis + 1) % o.K,
			K:    o.K,
			N:    o.N,
		})
		close(ch)
	}(l)
	go func(ch chan *N[T]) {
		ch <- New[T](O[T]{
			Data: o.Data[pivot+1 : len(o.Data)],
			Axis: (o.Axis + 1) % o.K,
			K:    o.K,
			N:    o.N,
		})
		close(ch)
	}(r)

	node.left = <-l
	node.right = <-r

	return node
}

func (n *N[T]) Nil() bool       { return n == nil }
func (n *N[T]) L() node.N[T]    { return n.left }
func (n *N[T]) R() node.N[T]    { return n.right }
func (n *N[T]) Leaf() bool      { return n.pivot == nil }
func (n *N[T]) Axis() vector.D  { return n.axis }
func (n *N[T]) Pivot() vector.V { return n.pivot }
func (n *N[T]) Data() []T       { return n.data }

// hoare partitions the input data by the pivot.
//
// N.B.: The high index is exclusive -- that is, when partitioning an entire
// array, high should be set to len(data).
func hoare[T point.P](data []T, pivot int, low int, high int, less func(a vector.V, b vector.V) bool) int {
	if pivot < 0 || low < 0 || high < 0 || pivot >= len(data) || low >= len(data) || high > len(data) {
		return -1
	}

	// hoare partitioning requires the pivot at the beginning of the array.
	data[pivot], data[low] = data[low], data[pivot]

	// i and j are the left and right tracker indices, respectively. i is
	// strictly increasing, while j is strictly decreasing
	i := low + 1
	j := high - 1

	for i <= j {
		// Skip array elements which are already sorted. Note that we
		// are partitioning such that
		//
		// l < p <= r
		for ; less(data[i].P(), data[low].P()) && i < j; i++ {
		}
		for ; !less(data[j].P(), data[low].P()) && j > 0; j-- {
		}

		if i > j {
			break
		}

		data[i], data[j] = data[j], data[i]

		i++
		j--

	}

	// Since the pivot is stored at the beginning of the array, we need to
	// do a final swap to ensure the pivot is at the right position.
	data[low], data[i-1] = data[i-1], data[low]
	return i - 1
}
