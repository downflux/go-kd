package tree

import (
	"fmt"

	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/vector"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

type O[T point.P] struct {
	Data []T
	K    vnd.D

	// N is the nominal leaf size of a node.
	N int

	Axis vnd.D
}

type N[T point.P] struct {
	data []T

	k     vnd.D
	pivot vnd.V
	axis  vnd.D
	left  *N[T]
	right *N[T]
}

func validate[T point.P](o O[T]) error {
	if o.Axis >= o.K {
		return fmt.Errorf("given node dimension greater than vnd dimension: %v > %v", o.Axis, o.K)
	}
	if o.N < 1 {
		return fmt.Errorf("given leaf node size must be a positive integer")
	}
	return nil
}

// New recursively constructs a node object given the input data.
func New[T point.P](o O[T]) *N[T] {
	if err := validate(o); err != nil {
		panic(fmt.Sprintf("could not construct node: %v", err))
	}

	if len(o.Data) <= 0 {
		return nil
	}
	if len(o.Data) <= o.N {
		return &N[T]{
			data: o.Data,
			axis: o.Axis,
			k:    o.K,
		}
	}
	pivot := hoare(o.Data, 0, 0, len(o.Data), func(a vnd.V, b vnd.V) bool { return a.X(o.Axis) < b.X(o.Axis) })

	node := &N[T]{
		data:  []T{o.Data[pivot]},
		pivot: o.Data[pivot].P(),
		axis:  o.Axis,
		k:     o.K,
	}

	// Node construction can be concurrent since we guarantee child nodes
	// will never access data across the high / low boundary. Note that this
	// does increase the number of allocs ~3x, and is less performant for
	// low data sizes (i.e. for less than ~10k points). We can optionally
	// branch tree construction on the length of incoming data if we find
	// the performance prohibitively slow.
	l := make(chan *N[T])
	r := make(chan *N[T])

	go func(ch chan<- *N[T]) {
		ch <- New[T](O[T]{
			Data: o.Data[0:pivot],
			Axis: (o.Axis + 1) % o.K,
			K:    o.K,
			N:    o.N,
		})
		close(ch)
	}(l)
	go func(ch chan<- *N[T]) {
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

func (n *N[T]) Nil() bool    { return n == nil }
func (n *N[T]) L() node.N[T] { return n.left }
func (n *N[T]) R() node.N[T] { return n.right }
func (n *N[T]) Leaf() bool   { return n.pivot == nil }
func (n *N[T]) Axis() vnd.D  { return n.axis }
func (n *N[T]) Pivot() vnd.V { return n.pivot }
func (n *N[T]) Data() []T    { return n.data }
func (n *N[T]) K() vnd.D     { return n.k }

func (n *N[T]) Insert(p T) {
	if n.Leaf() || !n.Leaf() && vnd.Within(n.Pivot(), p.P()) {
		n.data = append(n.data, p)
		return
	}

	if vector.Comparator(n.Axis()).Less(p.P(), n.Pivot()) {
		n.L().Insert(p)
	}
	n.R().Insert(p)
}

func (n *N[T]) Remove(v vnd.V, f func(p T) bool) (bool, T) {
	var blank T
	if n.Leaf() || !n.Leaf() && vnd.Within(n.Pivot(), v) {
		for i, p := range n.Data() {
			if f(p) {
				n.data[i], n.data[len(n.data)-1] = n.data[len(n.data)-1], blank
				n.data = n.data[:len(n.data)-1]

				return true, p
			}
		}
		return false, blank
	}

	if vector.Comparator(n.Axis()).Less(v, n.Pivot()) {
		return n.L().Remove(v, f)
	}
	return n.R().Remove(v, f)
}

// hoare partitions the input data by the pivot.
//
// N.B.: The high index is exclusive -- that is, when partitioning an entire
// array, high should be set to len(data).
func hoare[T point.P](data []T, pivot int, low int, high int, less func(a vnd.V, b vnd.V) bool) int {
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
