package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/filter"
	"github.com/downflux/go-kd/kd"
	"github.com/downflux/go-kd/point"
)

type KD[T point.P] kd.KD[T]

func (t *KD[T]) KNN(p vector.V, k int, f filter.F[T]) []T { return kd.KNN((*kd.KD[T])(t), p, k, f) }
func (t *KD[T]) RangeSearch(q hyperrectangle.R, f filter.F[T]) []T {
	return kd.RangeSearch((*kd.KD[T])(t), q, f)
}
func (t *KD[T]) Data() []T                                  { return kd.Data((*kd.KD[T])(t)) }
func (t *KD[T]) Balance()                                   { (*kd.KD[T])(t).Balance() }
func (t *KD[T]) Insert(p T)                                 { (*kd.KD[T])(t).Insert(p) }
func (t *KD[T]) Remove(v vector.V, f filter.F[T]) (bool, T) { return (*kd.KD[T])(t).Remove(v, f) }
