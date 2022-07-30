package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/filter"
	"github.com/downflux/go-kd/internal/knn"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/node/tree"
	"github.com/downflux/go-kd/internal/rangesearch"
	"github.com/downflux/go-kd/point"
)

type O[T point.P] struct {
	Data []T
	K    vector.D

	// N is the leaf size of the k-D tree. Leaf nodes are checked via
	// bruteforce methods.
	N int
}

type KD[T point.P] struct {
	k vector.D
	n int

	root node.N[T]
}

func New[T point.P](o O[T]) *KD[T] {
	data := make([]T, len(o.Data))
	if l := copy(data, o.Data); l != len(o.Data) {
		panic("could not copy data into k-D tree")
	}
	if o.K < 1 {
		panic("k-D tree must contain points with non-zero length vectors")
	}
	if o.N < 1 {
		panic("k-D tree minimum leaf node size must be positive")
	}

	t := &KD[T]{
		k: o.K,
		n: o.N,
		root: tree.New[T](tree.O[T]{
			Data: data,
			Axis: vector.AXIS_X,
			K:    o.K,
			N:    o.N,
		}),
	}

	return t
}

func (t *KD[T]) Balance() {
	t.root = tree.New[T](tree.O[T]{
		Data: Data(t),
		Axis: vector.AXIS_X,
		K:    t.k,
		N:    t.n,
	})
}

func (t *KD[T]) Insert(p T)                                 { t.root.Insert(p) }
func (t *KD[T]) Remove(v vector.V, f filter.F[T]) (T, bool) { return t.root.Remove(v, f) }

func KNN[T point.P](t *KD[T], p vector.V, k int, f filter.F[T]) []T {
	return knn.KNN(t.root, p, k, f)
}
func RangeSearch[T point.P](t *KD[T], q hyperrectangle.R, f filter.F[T]) []T {
	return rangesearch.RangeSearch(t.root, q, f)
}

func Data[T point.P](t *KD[T]) []T {
	if t.root.Nil() {
		return nil
	}
	var data []T

	var n node.N[T]
	open := []node.N[T]{t.root}
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
