// Package knn implements k-nearest neighbors search algorithm on a K-D tree node.
package knn

import (
	"math"
	"fmt"

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
	stats.recursion = 0
	stats.total = 0
	if n == nil {
		stats.size = 0
	} else {
		stats.size = n.Size()
	}

	q := pq.New(k)
	knn(n, p, q, map[*node.N]bool{})

	ns := make([]*node.N, 0, k)
	for !q.Empty() {
		ns = append(ns, q.Pop())
	}

	// The priority queue q sorts data with highest priority (read:
	// distance) first. We want to return data closest the the query point
	// at the head of the slice, so we need to reverse this ordering.
	for i, j := 0, len(ns)-1; i < j; i, j = i+1, j-1 {
		ns[i], ns[j] = ns[j], ns[i]
	}

	fmt.Println("DEBUG: statsingleton ==", stats)
	return ns
}

// DEBUG
type statSingleton struct {
	size int
	recursion int
	total int
}

var stats *statSingleton = &statSingleton{}

// knn recursively searches for the k-nearest neighbors of a node. The priority
// queue input q in effect tracks the searched space.
func knn(n *node.N, p vector.V, q *pq.Q, closed map[*node.N]bool) {
	stats.total += 1

	if n == nil || closed[n] {
		return
	}

	stats.recursion += 1

	closed[n] = true

        // TODO THIS IS NOT CORRECT -- WE ONLY GENERATE PATH ONCEE, PASS INTO knn, AND THAT REPEATEDLY POPS HEAD, TAIL, EXPANDING INTO HEAD IF NECESSARY
	for _, n := range path(n, p) {
		if d := vector.Magnitude(vector.Sub(p, n.P())); !q.Full() || d < q.Priority() {
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

			knn(c, p, q, closed)
		}
	}
}
