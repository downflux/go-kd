// Package knn implements k-nearest neighbors search algorithm on a K-D tree node.
package knn

import (
	"math"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/pq"
)

// path generates a list of nodes to the root, starting from a leaf.
func path(n *node.N, v vector.V) []*node.N {
	if n == nil {
		return nil
	}

	if n.Leaf() {
		return []*node.N{n}
	}

	// Note that we are bypassing the v == n.P() stop condition check -- we
	// are always continuing to the leaf node. This is necessary for finding
	// multiple closest neighbors, as we care about points in the tree which
	// do not have to coincide with the point coordinates.

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
	knn(n, p, q, make([]float64, p.Dimension()))

	// q.Pop() returns furthest distance first, so we need to reverse the
	// queue for the KNN output.
	ns := make([]*node.N, q.Len())
	for i := q.Len() - 1; i >= 0; i-- {
		ns[i] = q.Pop()
	}

	return ns
}

// sub uses the scratch space to calculate the difference between two vectors.
//
// Care should be taken that the scratch space is not used concurrently, as
// subsequent calls will modify the slice.
func sub(v vector.V, u vector.V, scratch []float64) vector.V {
	for i := vector.D(0); i < v.Dimension(); i++ {
		scratch[i] = v.X(i) - u.X(i)
	}
	return vector.V(scratch)
}

// knn recursively searches for the k-nearest neighbors of a node. The priority
// queue input q in effect tracks the searched space.
//
// The scratch input is used to reduce the amount of vector.Sub operations, as
// the returned vector is copy-constructed.
//
// This does not seem have a significant impact relative to the overall
// execution time, but we do see the overall allocs cut in half.
//
// TODO(minkezhang): Determine if this optimization is actually useful.
func knn(n *node.N, p vector.V, q *pq.Q, scratch []float64) {
	if n == nil {
		return
	}

	for _, n := range path(n, p) {
		if d := vector.Magnitude(sub(p, n.P(), scratch)); !q.Full() || d < q.Priority() {
			q.Push(n, d)
		}

		x := p.X(n.Axis())
		nx := n.P().X(n.Axis())

		// The minimal distance so far exceeds the current node split
		// plane, so we need to expand into the complement (i.e.
		// unexplored) child nodes.
		if q.Priority() > math.Abs(nx-x) {

			// We normally will expand the left child node if
			//
			//   x < nx
			//
			// however, this was already expanded while
			// generating the queue; therefore, we want to expand to
			// the unexplored, complement child instead.
			var c *node.N
			if x < nx {
				c = n.R()
			} else {
				c = n.L()
			}

			knn(c, p, q, scratch)
		}
	}
}
