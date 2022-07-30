// Package kd implements a k-D tree with arbitrary data packing and duplicate
// data coordinate support.
//
// k-D trees are generally a cacheing layer representation of the local state --
// we do not expect to be making frequent mutations to this tree once
// constructed.
//
// Read operations on this k-D tree may be done in parallel. Mutations on the
// k-D tree must be done serially.
//
// N.B.: Mutating the data point positions must be accompanied by mutating the
// k-D tree. For large numbers of points, and for a large number of queries, the
// time taken to build the tree will be offset by the speedup of subsequent
// reads.
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

	// N is the nominal leaf size of the k-D tree. Leaf nodes are checked
	// via bruteforce methods.
	//
	// Note that individual nodes (including non-leaf nodes) may contain
	// elements that exceed this size constraint after inserts and removes.
	//
	// Leaf size will significantly impact performance -- users should
	// tailor this value to their specific use-case. We recommend setting
	// this value to 16 and up as the size of the data set increases.
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

// Balance reconstructs the k-D tree.
//
// This k-D tree implementation does not support concurrent mutations.
func (t *KD[T]) Balance() {
	t.root = tree.New[T](tree.O[T]{
		Data: Data(t),
		Axis: vector.AXIS_X,
		K:    t.k,
		N:    t.n,
	})
}

// Insert adds a new point into the k-D tree.
//
// Insert is not a balanced operation -- after many mutations, the tree should
// be explicitly reconstructed.
//
// This k-D tree implementation does not support concurrent mutations.
func (t *KD[T]) Insert(p T) { t.root.Insert(p) }

// Remove pops a point from the k-D tree which lies at the input vector v and
// matches the filter. Note that if multiple points match both the location
// vector and the filter, an arbitrary one will be removed. This function will
// pop at most one element from the k-D tree.
//
// Remove is not a balanced operation -- after many mutations, the tree should
// be explicitly reconstructed.
//
// This k-D tree implementation does not support concurrent mutations.
//
// If there is no matching point, the returned bool will be false.
func (t *KD[T]) Remove(v vector.V, f filter.F[T]) (T, bool) { return t.root.Remove(v, f) }

// KNN returns the k nearest neighbors to the input vector p and matches the
// filter function.
//
// This k-D tree implementation supports concurrent read operations.
func KNN[T point.P](t *KD[T], p vector.V, k int, f filter.F[T]) []T {
	return knn.KNN(t.root, p, k, f)
}

// RangeSearch returns all points which are found in the given bounds and
// matches the filter function.
//
// This k-D tree implementation supports concurrent read operations.
func RangeSearch[T point.P](t *KD[T], q hyperrectangle.R, f filter.F[T]) []T {
	return rangesearch.RangeSearch(t.root, q, f)
}

// Data returns all points in the k-D tree.
//
// This k-D tree implementation supports concurrent read operations.
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
