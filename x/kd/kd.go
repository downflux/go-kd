package kd

import (
	"github.com/downflux/go-kd/x/point"
)

type O[p point.P] struct {
	Data []p
	K    point.D
	N    int
}

type N struct {
	Low  int
	High int
	Mid  int

	Left  *N
	Right *N
}

type T[p point.P] struct {
	k    point.D
	n    int
	data []p
}

func New[p point.P](o O[p]) *T[p] {
	data := make([]p, len(o.Data))
	if l := copy(data, o.Data); l != len(o.Data) {
		panic("could not copy data into k-D tree")
	}
	return nil
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
