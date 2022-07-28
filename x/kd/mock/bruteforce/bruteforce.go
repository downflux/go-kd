package bruteforce

import (
	"math"
	"sort"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point"
)

type L[T point.P] []T

func New[T point.P](d []T) *L[T] {
	data := make([]T, len(d))
	if l := copy(data, d); l != len(d) {
		panic("could not copy data into brute force list")
	}
	m := L[T](data)
	return &m
}

func (m *L[T]) KNN(p vector.V, k int) []T {
	sort.Sort(l[T]{
		data: *m,
		p:    p,
	})

	return []T(*m)[0:int(math.Min(float64(k), float64(len(*m))-1))]
}

func (m *L[T]) Data() []T { return *m }

type l[T point.P] struct {
	data []T
	p    vector.V
}

func (l l[T]) Len() int      { return len(l.data) }
func (l l[T]) Swap(i, j int) { l.data[i], l.data[j] = l.data[j], l.data[i] }

func (l l[T]) Less(i, j int) bool {
	return vector.SquaredMagnitude(
		vector.Sub(l.data[i].P(), l.p),
	) < vector.SquaredMagnitude(
		vector.Sub(l.data[j].P(), l.p),
	)
}
