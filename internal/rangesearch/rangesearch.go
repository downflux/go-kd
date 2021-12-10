// Package rangesearch implements a range search algorithm for a K-D tree.
package rangesearch

import (
	"math"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/node"
)

// Search traverses the K-D tree node and returns all child nodes which are
// contained in the bounding rectangle.
func Search(n *node.N, r hyperrectangle.R) []*node.N {
	if n == nil {
		return nil
	}

	min := make([]float64, n.P().Dimension())
	max := make([]float64, n.P().Dimension())

	for i := vector.D(0); i < n.P().Dimension(); i++ {
		min[i] = math.Inf(-1)
		max[i] = math.Inf(0)
	}

	return search(n, &r, *hyperrectangle.New(*vector.New(min...), *vector.New(max...)))
}

// search recursively travels through the tree node to look for nodes within the
// input K-dimensional rectangle. The K-dimensional rectangle input bound is a
// recursion artifact which keeps track of the bounding box of the current node.
func search(n *node.N, r *hyperrectangle.R, bound hyperrectangle.R) []*node.N {
	if _, ok := (*r).Intersect(bound); n == nil || !ok {
		return nil
	}

	var ns []*node.N

	if r.In(n.P()) {
		ns = append(ns, n)
	}

	if c := n.L(); c != nil {
		b := *hyperrectangle.New(
			bound.Min(),
			bound.Max(),
		)
		// WLOG the bounding box of the left child node will not
		// contain data from points in the right node. We know the
		// upper bound of the data in the left node due by definition
		// of a K-D tree node.
		b.Max()[n.Axis()] = n.P().X(n.Axis())
		ns = append(ns, search(c, r, b)...)
	}
	if c := n.R(); c != nil {
		b := *hyperrectangle.New(
			bound.Min(),
			bound.Max(),
		)
		b.Min()[n.Axis()] = n.P().X(n.Axis())
		ns = append(ns, search(c, r, b)...)
	}

	return ns
}
