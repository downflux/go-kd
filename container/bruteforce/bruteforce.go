package bruteforce

import (
	"sort"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/filter"
	"github.com/downflux/go-kd/internal/perf/util"
	"github.com/downflux/go-kd/point"
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

func (m *L[T]) KNN(p vector.V, k int, f filter.F[T]) []T {
	sort.Sort(util.L[T]{
		Data: *m,
		P:    p,
	})

	var data []T
	for _, p := range *m {
		if f(p) {
			data = append(data, p)
		}
		if len(data) == k {
			return data
		}
	}
	return data
}

func (m *L[T]) RangeSearch(q hyperrectangle.R, f filter.F[T]) []T {
	var data []T
	for _, p := range m.Data() {
		if q.In(p.P()) && f(p) {
			data = append(data, p)
		}
	}
	return data
}

func (m *L[T]) Insert(p T) { *m = append(*m, p) }
func (m *L[T]) Remove(p vector.V, f filter.F[T]) (T, bool) {
	var blank T
	for i, q := range *m {
		if vector.Within(p, q.P()) && f(q) {
			(*m)[i], (*m)[len(*m)-1] = (*m)[len(*m)-1], blank
			*m = (*m)[:len(*m)-1]
			return q, true
		}
	}
	return blank, false
}

func (m *L[T]) Data() []T { return *m }
func (m *L[T]) Balance()  {}
