// Package kyroy is a wrapper around the @kyroy k-D tree implementation. This is
// used for performance testing.
package kyroy

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/filter"
	"github.com/downflux/go-kd/point"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

type KD[T point.P] kdtree.KDTree

func New[T point.P](data []T) *KD[T] {
	var ps []kdtree.Point
	for _, p := range data {
		ps = append(ps, points.NewPoint([]float64(p.P()), p))
	}
	return (*KD[T])(kdtree.New(ps))
}

func (t *KD[T]) KNN(p vector.V, k int, f filter.F[T]) []T {
	var data []T
	for _, p := range (*kdtree.KDTree)(t).KNN(
		points.NewPoint([]float64(p), nil),
		k,
	) {
		if f(p.(*points.Point).Data.(T)) {
			data = append(data, p.(*points.Point).Data.(T))
		}
	}
	return data
}

func (t *KD[T]) Data() []T {
	var data []T
	for _, p := range (*kdtree.KDTree)(t).Points() {
		data = append(data, p.(T))
	}
	return data
}

func (t *KD[T]) RangeSearch(q hyperrectangle.R, f filter.F[T]) []T {
	var r [][2]float64
	for i := vector.D(0); i < q.Min().Dimension(); i++ {
		r = append(r, [2]float64{q.Min().X(i), q.Max().X(i)})
	}

	var data []T
	for _, p := range (*kdtree.KDTree)(t).RangeSearch(r) {
		if f(p.(*points.Point).Data.(T)) {
			data = append(data, p.(*points.Point).Data.(T))
		}
	}
	return data
}

func (t *KD[T]) Balance()   { (*kdtree.KDTree)(t).Balance() }
func (t *KD[T]) Insert(p T) { (*kdtree.KDTree)(t).Insert(points.NewPoint([]float64(p.P()), p)) }

func (t *KD[T]) Remove(p vector.V, f filter.F[T]) (T, bool) {
	v := (*kdtree.KDTree)(t).Remove(points.NewPoint([]float64(p), nil)).(*points.Point).Data.(T)
	if !f(v) {
		var blank T
		(*kdtree.KDTree)(t).Insert(points.NewPoint([]float64(p), v))
		return blank, false
	}
	return v, true
}
