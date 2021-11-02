// Package rangesearch implements a range search algorithm for a K-D tree.

package rangesearch

import (
	"math"

	"github.com/downflux/go-geometry/rectangle"
	"github.com/downflux/go-geometry/vector"
	"github.com/downflux/go-kd/axis"
	"github.com/downflux/go-kd/node"
)

// Search traverses the K-D tree node and returns all child nodes which are
// contained in the bounding rectangle.
func Search(n *node.N, r rectangle.R) []*node.N {
	return search(n, r, *rectangle.New(
		*vector.New(math.Inf(-1), math.Inf(-1)),
		*vector.New(math.Inf(0), math.Inf(0)),
	))
}

// search recursively travels through the tree node to look for nodes within the
// input rectangle. The bounding rectangle is a recursion artifact which keeps
// track of the bounding box of the current node.
func search(n *node.N, r rectangle.R, bound rectangle.R) []*node.N {
	if _, ok := r.Intersect(bound); n == nil || !ok {
		return nil
	}

	var ns []*node.N

	if r.In(n.V()) {
		ns = append(ns, n)
	}

	// Calculate the new bounding boxes of the child nodes.
	lb := map[axis.Type]rectangle.R{
		axis.Axis_X: *rectangle.New(
			bound.Min(),
			*vector.New(n.V().X(), bound.Max().Y()),
		),
		axis.Axis_Y: *rectangle.New(
			bound.Min(),
			*vector.New(bound.Max().X(), n.V().Y()),
		),
	}[n.Axis()]
	rb := map[axis.Type]rectangle.R{
		axis.Axis_X: *rectangle.New(
			*vector.New(n.V().X(), bound.Min().Y()),
			bound.Max(),
		),
		axis.Axis_Y: *rectangle.New(
			*vector.New(bound.Min().X(), n.V().Y()),
			bound.Max(),
		),
	}[n.Axis()]

	if c := n.L(); c != nil {
		ns = append(ns, search(c, r, lb)...)
	}
	if c := n.R(); c != nil {
		ns = append(ns, search(c, r, rb)...)
	}

	return ns
}
