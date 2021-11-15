// Package kd implements a 2D K-D tree with arbitrary data packing and duplicate
// data coordinate support.
//
// K-D trees are generally a cacheing layer representation of the local state --
// we do not expect to be making frequent mutations to this tree once
// constructed.
//
// The tree will have to be rebalanced as data points move in the simulation.
// For large numbers of points, and for a large number of queries, the time
// taken to build the tree will be offset by the speedup of subsequent reads.
package kd

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/knn"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/rangesearch"
	"github.com/downflux/go-kd/point"
)

// T is a K-D tree implementation.
type T struct {
	root *node.N
}

func New(data []point.P) (*T, error) {
	if len(data) > 0 {
		k := data[0].P().Dimension()
		for _, c := range data {
			if c.P().Dimension() != k {
				return nil, fmt.Errorf("cannot create a K-D tree with data of inconsistent dimensions")
			}
		}
	}
	return &T{
		root: node.New(data, 0),
	}, nil
}

// Balance rebalances the tree; note that in general, tree node mutations are
// expensive and messy, so much so that it's easier to just redo the tree from
// scratch.
func (t *T) Balance() { t.root = node.New(node.Points(t.root), 0) }

// Insert adds a new data point into the tree.
//
// N.B.: This is not guaranteed to be a balanced insertion -- adding lots of
// points may unbalance the tree, causing queries to become much slower than in
// theory. If this function must be called a number of times, it's generally
// good practice to call Balance() in order to ensure optimal lookup times.
func (t *T) Insert(datum point.P) { t.root.Insert(datum) }

// Remove deletes an existing data point from the tree. This function will
// delete the first matching point with the given coordinates.
func (t *T) Remove(position vector.V, f func(datum point.P) bool) bool {
	return t.root.Remove(position, f)
}

// Filter returns a set of data points in the given bounding box. Data points
// are added to the returned set if they fall inside the bounding box and passes
// the given filter function.
func Filter(t *T, r hyperrectangle.R, f func(datum point.P) bool) []point.P {
	var data []point.P
	for _, n := range rangesearch.Search(t.root, r) {
		for _, p := range n.Data() {
			if f(p) {
				data = append(data, p)
			}
		}
	}
	return data
}

// RadialFilter returns a set of data points in the given bounding circle. Data
// points are added to the returned set if they fall inside the bounding circle
// and passes the given filter function.
func RadialFilter(t *T, c hypersphere.C, f func(datum point.P) bool) []point.P {
	r := *hyperrectangle.New(
		vector.Sub(c.P(), *vector.New(c.R(), c.R())),
		vector.Add(c.P(), *vector.New(c.R(), c.R())),
	)
	return Filter(t, r, func(datum point.P) bool {
		return vector.SquaredMagnitude(vector.Sub(datum.P(), c.P())) <= c.R()*c.R() && f(datum)
	})
}

// KNN returns the k nearest neighbors of the given search coordinates.
//
// N.B.: KNN will return at max k neighbors; in the degenerate case that
// multiple data points reside at the same spacial coordinate, this function
// will arbitrarily return a subset of these to fulfill the k neighbors
// criteria.
func KNN(t *T, position vector.V, k int) []point.P {
	var data []point.P
	for _, n := range knn.KNN(t.root, position, k) {
		if len(data) < k {
			data = append(data, n.Data()...)
		} else {
			break
		}
	}
	return data
}
