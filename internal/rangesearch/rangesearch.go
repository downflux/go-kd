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

	return search(n, r, *hyperrectangle.New(*vector.New(min...), *vector.New(max...)))
}

// search recursively travels through the tree node to look for nodes within the
// input rectangle. The bounding rectangle is a recursion artifact which keeps
// track of the bounding box of the current node.
func search(n *node.N, r hyperrectangle.R, bound hyperrectangle.R) []*node.N {
	if _, ok := r.Intersect(bound); n == nil || !ok {
		return nil
	}

	var ns []*node.N

	if r.In(n.P()) {
		ns = append(ns, n)
	}

	// Calculate the new bounding boxes of the child nodes.
	lbMin := make([]float64, n.P().Dimension())
	lbMax := make([]float64, n.P().Dimension())

	for i := vector.D(0); i < n.P().Dimension(); i++ {
		lbMin[i] = bound.Min().X(i)
		lbMax[i] = bound.Max().X(i)

		if i == n.Axis() {
			lbMax[i] = n.P().X(i)
		}
	}

	// Calculate the new bounding boxes of the child nodes.
	rbMin := make([]float64, n.P().Dimension())
	rbMax := make([]float64, n.P().Dimension())

	for i := vector.D(0); i < n.P().Dimension(); i++ {
		rbMin[i] = bound.Min().X(i)
		rbMax[i] = bound.Max().X(i)

		if i == n.Axis() {
			rbMin[i] = n.P().X(i)
		}
	}

	if c := n.L(); c != nil {
		ns = append(ns, search(c, r, *hyperrectangle.New(
			*vector.New(lbMin...),
			*vector.New(lbMax...)))...)
	}
	if c := n.R(); c != nil {
		ns = append(ns, search(c, r, *hyperrectangle.New(
			*vector.New(rbMin...),
			*vector.New(rbMax...)))...)
	}

	return ns
}
