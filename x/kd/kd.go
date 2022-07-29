package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/knn"
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/internal/node/tree"
	"github.com/downflux/go-kd/x/internal/rangesearch"
	"github.com/downflux/go-kd/x/point"
)

type O[U point.P] struct {
	Data []U
	K    vector.D
	N    int
}

type T[U point.P] struct {
	k vector.D
	n int

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
		k: o.K,
		n: o.N,
		root: tree.New[U](tree.O[U]{
			Data: data,
			Axis: 0,
			K:    o.K,
			N:    o.N,
		}),
	}

	return t
}

func KNN[U point.P](t *T[U], p vector.V, k int, f func(p U) bool) []U {
	return knn.KNN(t.root, p, k, f)
}
func RangeSearch[U point.P](t *T[U], q hyperrectangle.R, f func(p U) bool) []U {
	return rangesearch.RangeSearch(t.root, q, f)
}

func Data[U point.P](t *T[U]) []U {
	if t.root.Nil() {
		return nil
	}
	var data []U

	var n node.N[U]
	open := []node.N[U]{t.root}
	for len(open) > 0 {
		n, open = open[0], open[1:]

		data = append(data, n.Data()...)
		if !n.L().Nil() {
			open = append(open, n.L())
		}
		if !n.R().Nil() {
			open = append(open, n.R())
		}
	}

	return data
}
