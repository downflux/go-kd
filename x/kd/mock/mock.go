package mock

import (
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point"
)

// I is a mock interface used for faciliating benchmark tests. This
// should not be consumed by external users.
type I[U point.P] interface {
	KNN(p vector.V, k int) []U
}

type T[U point.P] kd.T[U]

func (t *T[U]) KNN(p vector.V, k int) []U { return kd.KNN((*kd.T[U])(t), p, k) }
func (t *T[U]) Data() []U                 { return kd.Data((*kd.T[U])(t)) }
