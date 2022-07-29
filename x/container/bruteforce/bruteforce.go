package bruteforce

import (
	"sort"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/filter"
	"github.com/downflux/go-kd/x/internal/perf/util"
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

func (m *L[T]) Data() []T { return *m }
