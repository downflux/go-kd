package wrapper

import (
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point"
)

type T[U point.P] kd.T[U]

func (t *T[U]) KNN(p vector.V, k int) []U { return kd.KNN((*kd.T[U])(t), p, k) }
