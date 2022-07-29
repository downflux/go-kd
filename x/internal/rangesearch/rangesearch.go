package rangesearch

import (
	"math"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
)

func RangeSearch[T point.P](n node.N[T], q hyperrectangle.R) []T {
	if n.Nil() {
		return nil
	}

	min := make([]float64, n.K())
	max := make([]float64, n.K())

	for i := vector.D(0); i < n.K(); i++ {
		min[i] = math.Inf(-1)
		max[i] = math.Inf(0)
	}

	return search(n, q, *hyperrectangle.New(vector.V(min), vector.V(max)))
}

func search[T point.P](n node.N[T], q hyperrectangle.R, bound hyperrectangle.R) []T {
	if _, ok := q.Intersect(bound); n.Nil() || !ok {
		return nil
	}

	var data []T
	for _, p := range n.Data() {
		if q.In(p.P()) {
			data = append(data, p)
		}
	}

	if n.Leaf() {
		return data
	}

	l := make(chan []T)
	r := make(chan []T)

	go func(ch chan<- []T) {
		max := make([]float64, n.K())
		copy(max, bound.Max())
		max[n.Axis()] = n.Pivot().X(n.Axis())

		bound := *hyperrectangle.New(bound.Min(), max)
		ch <- search(n.L(), q, bound)
		close(ch)
	}(l)
	go func(ch chan<- []T) {
		min := make([]float64, n.K())
		copy(min, bound.Min())
		min[n.Axis()] = n.Pivot().X(n.Axis())

		bound := *hyperrectangle.New(min, bound.Max())
		ch <- search(n.R(), q, bound)
		close(ch)
	}(r)

	data = append(data, <-l...)
	data = append(data, <-r...)

	return data

}
