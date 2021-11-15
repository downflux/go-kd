// Package knn implements k-nearest neighbors search algorithm on a K-D tree node.
package knn

import (
	"math"

	"github.com/downflux/go-geometry/nd/vector"
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
func path(n *node.N, v vector.V) []*node.N {
	if n == nil {
		return nil
	}

	if n.Leaf() {
		return []*node.N{n}
	}

	// Note that we are bypassing the v == n.P() stop condition check -- we
	// are always continuing to the leaf node.

	x := v.X(n.Axis())
	nx := n.P().X(n.Axis())

	if x < nx {
		return append(path(n.L(), v), n)
	}
	return append(path(n.R(), v), n)
}

// KNN returns the k-nearest neighbors of a given node. Nodes are returned in
// sorted order, with nodes closest to the input vector at the head of the
// slice.
func KNN(n *node.N, p vector.V, k int) []*node.N {
	q := pq.New(k)
	knn(n, p, q)

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
func knn(n *node.N, p vector.V, q *pq.Q) {
	if n == nil {
		return
	}

	for _, n := range path(n, p) {
		if d := vector.Magnitude(vector.Sub(p, n.P())); !q.Full() || d < q.Priority() {
			q.Push(n, d)
		}

		x := p.X(n.Axis())
		nx := n.P().X(n.Axis())

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

			knn(c, p, q)
		}
	}
}
