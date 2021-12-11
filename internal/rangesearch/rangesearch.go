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

	return search(
		n,
		&r,
		hyperrectangle.New(vector.V(min), vector.V(max)),
		hyperrectangle.New(
			vector.V(make([]float64, n.P().Dimension())),
			vector.V(make([]float64, n.P().Dimension())),
		),
	)
}

// search recursively travels through the tree node to look for nodes within the
// input K-dimensional rectangle. The K-dimensional rectangle input bound is a
// recursion artifact which keeps track of the bounding box of the current node.
func search(n *node.N, r *hyperrectangle.R, bound *hyperrectangle.R, buf *hyperrectangle.R) []*node.N {
	if ok := (*r).IntersectBuf(*bound, buf); n == nil || !ok {
		return nil
	}

	var ns []*node.N

	if r.In(n.P()) {
		ns = append(ns, n)
	}

	if c := n.L(); c != nil {
		// bound is a pointer to a hyperrectangle and is re-used if the
		// right branch is also not empty -- we need to make sure this
		// branch does not have any non-hermetic effects on the sibling
		// branch.
		max := bound.Max()[n.Axis()]

		// WLOG the bounding box of the left child node will not
		// contain data from points in the right node. We know the
		// upper bound of the data in the left node due by definition
		// of a K-D tree node.
		bound.Max()[n.Axis()] = n.P().X(n.Axis())
		ns = append(ns, search(c, r, bound, buf)...)

		// Restore the bounding box dimension for the right child.
		bound.Max()[n.Axis()] = max
	}
	if c := n.R(); c != nil {
		min := bound.Min()[n.Axis()]

		bound.Min()[n.Axis()] = n.P().X(n.Axis())
		ns = append(ns, search(c, r, bound, buf)...)

		bound.Min()[n.Axis()] = min
	}

	return ns
}
