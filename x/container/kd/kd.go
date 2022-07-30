package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/filter"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point"
)

type KD[T point.P] kd.KD[T]

func (t *KD[T]) KNN(p vector.V, k int, f filter.F[T]) []T { return kd.KNN((*kd.KD[T])(t), p, k, f) }
func (t *KD[T]) RangeSearch(q hyperrectangle.R, f filter.F[T]) []T {
	return kd.RangeSearch((*kd.KD[T])(t), q, f)
}
func (t *KD[T]) Data() []T { return kd.Data((*kd.KD[T])(t)) }
func (t *KD[T]) Balance()  { (*kd.KD[T])(t).Balance() }
