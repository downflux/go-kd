package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point"
)

type T[U point.P] kd.T[U]

func (t *T[U]) KNN(p vector.V, k int) []U          { return kd.KNN((*kd.T[U])(t), p, k) }
func (t *T[U]) RangeSearch(q hyperrectangle.R) []U { return kd.RangeSearch((*kd.T[U])(t), q) }
func (t *T[U]) Data() []U                          { return kd.Data((*kd.T[U])(t)) }
