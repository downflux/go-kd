package node

import (
	"fmt"

	"github.com/downflux/go-kd/x/point"
)

type O[p point.P] struct {
	Data []p
	K    point.D
	N    int

	Axis point.D
	Low  int
	High int
}

// N is a tree node structure which tracks slices of an input data array that
// are managed in the node and children. Leaf nodes set the pivot index to -1.
// Data to the left of current node have smaller values of the given axis.
type N struct {
	Low   int
	High  int
	Pivot int

	Axis  point.D
	Left  *N
	Right *N
}

// New recursively constructs a node object given the input data.
func New[p point.P](o O[p]) *N {
	if o.Axis > o.K {
		panic(fmt.Sprintf("given node dimension greater than vector dimension: %v > %v", o.Axis, o.K))
	}
	if o.N < 1 {
		panic("given leaf node size must be a positive integer")
	}
	n := o.High - o.Low
	if n <= 0 {
		return nil
	}
	if n <= o.N {
		return &N{
			Low:   o.Low,
			High:  o.High,
			Pivot: -1,
			Axis:  o.Axis,
		}
	}
	pivot := hoare(o.Data, o.Low, o.Low, o.High, func(a p, b p) bool { return a.P().X(o.Axis) < b.P().X(o.Axis) })

	node := &N{
		Low:   o.Low,
		High:  o.High,
		Pivot: pivot,
		Axis:  o.Axis,
	}

	// Node construction can be concurrent since we guarantee child nodes
	// will never access data across the high / low boundary. Note that this
	// does increase the number of allocs ~3x, and is less performant for
	// low data sizes (i.e. for less than ~10k points). We can optionally
	// branch tree construction on the length of incoming data if we find
	// the performance prohibitively slow.
	l := make(chan *N)
	r := make(chan *N)

	go func(ch chan *N) {
		ch <- New[p](O[p]{
			Data: o.Data,
			Axis: (o.Axis + 1) % o.K,
			K:    o.K,
			N:    o.N,
			Low:  o.Low,
			High: pivot,
		})
		close(ch)
	}(l)
	go func(ch chan *N) {
		ch <- New[p](O[p]{
			Data: o.Data,
			Axis: (o.Axis + 1) % o.K,
			K:    o.K,
			N:    o.N,
			Low:  pivot + 1,
			High: o.High,
		})
		close(ch)
	}(r)

	node.Left = <-l
	node.Right = <-r

	return node
}

// hoare partitions the input data by the pivot.
//
// N.B.: The high index is exclusive -- that is, when partitioning an entire
// array, high should be set to len(data).
func hoare[p point.P](data []p, pivot int, low int, high int, less func(a p, b p) bool) int {
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
		// Skip array elements which are already sorted.
		for ; less(data[i], data[low]) && i < j; i++ {
		}
		for ; less(data[low], data[j]); j-- {
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
