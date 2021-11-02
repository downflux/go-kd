// Package knn implements k-nearest neighbors search algorithm on a K-D tree node.
package knn

import (
	"math"

	"github.com/downflux/go-geometry/vector"
	"github.com/downflux/go-kd/internal/axis"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/pq"
)

// path generates a list of nodes to the root, starting from a leaf node,
// with a node guaranteed to contain the input coordinates.
//
// N.B.: We do not stop the recursion if we reach a node with matching
// coordinates; this is necessary for finding multiple closest neighbors, as we
// care about points in the tree which do not have to coincide with the point
// coordinates.
func path(n *node.N, v vector.V, tolerance float64) []*node.N {
	if n == nil {
		return nil
	}

	if n.Leaf() {
		return []*node.N{n}
	}

	// Note that we are bypassing the v == n.V() stop condition check -- we
	// are always continuing to the leaf node.

	x := axis.X(v, n.Axis())
	nx := axis.X(n.V(), n.Axis())

	if x < nx {
		return append(path(n.L(), v, tolerance), n)
	}
	return append(path(n.R(), v, tolerance), n)
}

// KNN returns the k-nearest neighbors of a given node. Nodes are returned in
// sorted order, with nodes closest to the input vector at the head of the
// slice.
func KNN(n *node.N, v vector.V, k int, tolerance float64) []*node.N {
	q := pq.New(k)
	knn(n, v, q, tolerance)

	var ns []*node.N
	for !q.Empty() {
		ns = append(ns, q.Pop())
	}

	// Return nearest nodes first by reversing the flattened queue order.
	for i, j := 0, len(ns)-1; i < j; i, j = i+1, j-1 {
		ns[i], ns[j] = ns[j], ns[i]
	}
	return ns
}

// knn recursively searches for the k-nearest neighbors of a node. The priority
// queue input effectively tracks the searched space.
func knn(n *node.N, v vector.V, q *pq.Q, tolerance float64) {
	if n == nil {
		return
	}

	p := path(n, v, tolerance)

	for _, n := range p {
		if d := vector.Magnitude(vector.Sub(v, n.V())); !q.Full() || d < q.Priority() {
			q.Push(n, d)
		}

		x := axis.X(v, n.Axis())
		nx := axis.X(n.V(), n.Axis())

		// The minimal distance so far exceeds the current node split
		// plane -- we need to expand into the child nodes.
		if q.Priority() > math.Abs(nx-x) {
			// We normally will expand the left child node if
			// x < nx; however, this was already expanded while
			// generating the queue; therefore, we want to expand to
			// the unexplored, complement child instead.
			var c *node.N
			if x < nx {
				c = n.R()
			} else {
				c = n.L()
			}

			knn(c, v, q, tolerance)
		}
	}
}
