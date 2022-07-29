package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/filter"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point"
)

type T[U point.P] kd.T[U]

func (t *T[U]) KNN(p vector.V, k int, f filter.F[U]) []U { return kd.KNN((*kd.T[U])(t), p, k, f) }
func (t *T[U]) RangeSearch(q hyperrectangle.R, f filter.F[U]) []U {
	return kd.RangeSearch((*kd.T[U])(t), q, f)
}
func (t *T[U]) Data() []U { return kd.Data((*kd.T[U])(t)) }
