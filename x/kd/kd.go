package kd

import (
	"github.com/downflux/go-kd/x/internal/knn"
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/internal/node/tree"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/vector"
)

type O[U point.P] struct {
	Data []U
	K    vector.D
	N    int
}

type T[U point.P] struct {
	k    vector.D
	n    int
	data []U

	root node.N[U]
}

func New[U point.P](o O[U]) *T[U] {
	data := make([]U, len(o.Data))
	if l := copy(data, o.Data); l != len(o.Data) {
		panic("could not copy data into k-D tree")
	}
	if o.K < 1 {
		panic("k-D tree must contain points with non-zero length vectors")
	}
	if o.N < 1 {
		panic("k-D tree minimum leaf node size must be positive")
	}

	t := &T[U]{
		k:    o.K,
		n:    o.N,
		data: data,
		root: tree.New[U](tree.O[U]{
			Data: data,
			Axis: 0,
			K:    o.K,
			N:    o.N,
			Low:  0,
			High: len(data),
		}),
	}

	return t
}

func KNN[U point.P](t *T[U], p vector.V, k int) []U { return knn.KNN(t.root, p, k) }
