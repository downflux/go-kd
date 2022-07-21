package node

import (
	"fmt"

	"github.com/downflux/go-kd/x/point"
)

type O[p point.P] struct {
	Data []p
	K    point.D
	N    int

	Dim  point.D
	Low  int
	High int
}

type N struct {
	Low   int
	High  int
	Pivot int

	Dim   point.D
	Left  *N
	Right *N
}

func New[p point.P](o O[p]) *N {
	if o.Dim > o.K {
		panic(fmt.Sprintf("given node dimension greater than vector dimension: %v > %v", o.Dim, o.K))
	}
	if o.N < 1 {
		panic("given leaf node size must be a positive integer")
	}
	if o.High-o.Low <= o.N {
		return &N{
			Low:   o.Low,
			High:  o.High,
			Pivot: -1,
			Dim:   o.Dim,
		}
	}
	pivot := hoare(o.Data, o.Low, o.Low, o.High, func(a p, b p) bool { return a.P().X(o.Dim) < b.P().X(o.Dim) })
	return &N{
		Low:   o.Low,
		High:  o.High,
		Pivot: pivot,
		Dim:   o.Dim,

		Left: New[p](O[p]{
			Data: o.Data,
			Dim:  (o.Dim + 1) % o.K,
			K:    o.K,
			N:    o.N,
			Low:  o.Low,
			High: pivot,
		}),
		Right: New[p](O[p]{
			Data: o.Data,
			Dim:  (o.Dim + 1) % o.K,
			K:    o.K,
			N:    o.N,
			Low:  pivot + 1,
			High: o.High,
		}),
	}
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
