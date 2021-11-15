// Package node implements a K-D tree node.
package node

import (
	"math"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/point/sorter"
	"github.com/downflux/go-kd/point"
)

// N represents a K-D tree node. Child nodes are sorted along the same axis,
// with coordinates in the left node preceeding right node coordinates (along
// the given axis).
type N struct {
	l *N
	r *N

	// depth represents the current level of the tree; the root node has a
	// depth of 0.
	depth int

	// p represents the coordinate of all data stored in this node.
	p vector.V

	// data is a list of data points stored in this node. All data here are
	// located at the same spacial coordinate.
	data []point.P

	// sizeCache tracks the current subtree size of the current node. A size
	// of 1 indicates this is a leaf node.
	//
	// A node contributes to the tree size if it contains some data.
	//
	// As calculating the tree size may only occur when during an insert or
	// delete operation, and such operations are rarely called in a K-D
	// tree, it is sensible to provide a size cache.
	sizeCache      int
	sizeCacheValid bool
}

func (n *N) Leaf() bool { return n.size() <= 1 }

// size returns the number of meaningful nodes in the current subtree. A size of
// 0 or 1 indicates n is a leaf node.
func (n *N) size() int {
	if n.sizeCacheValid {
		return n.sizeCache
	}

	var s int

	if l := n.L(); l != nil {
		s += l.size()
	}
	if r := n.R(); r != nil {
		s += r.size()
	}

	if len(n.data) > 0 {
		s += 1
	}

	n.sizeCacheValid = true
	n.sizeCache = s

	return n.sizeCache
}

// Axis is the discriminant dimension for this tree node -- if the node is
// split on the X-axis, then all points left of this node in the XY-plane are in
// the left child, and all points on or right of this node are in the right
// child.
func (n *N) Axis() vector.D { return vector.D(n.depth % int(n.p.Dimension())) }

// P is the point on the XY-plane to which this node is embedded. All data in
// this node are located at the same spacial coordinate, within a small margin
// of error.
func (n *N) P() vector.V { return n.p }

// Data returns a list of data stored in this node.
func (n *N) Data() []point.P { return append([]point.P{}, n.data...) }

// Child is the appropriately expanded child node given the input coordinates --
// that is, this function wraps the normal branching pattern for e.g. search
// operations.
func (n *N) Child(v vector.V) *N {
	if vector.Within(n.P(), v) {
		return nil
	}

	x := v.X(n.Axis())
	nx := n.P().X(n.Axis())

	if x < nx {
		return n.L()
	}
	return n.R()
}

// L is the meaningful left child node of the current node.
//
// This function returns nil if the left node does not contain any data.
func (n *N) L() *N {
	if n.l == nil || n.l.size() == 0 {
		return nil
	}
	return n.l
}

// R is the meaningful right child node of the current node.
func (n *N) R() *N {
	if n.r == nil || n.r.size() == 0 {
		return nil
	}
	return n.r
}

// Insert inserts a data point into the node. The point may be stored inside a
// child node.
func (n *N) Insert(p point.P) {
	// The number of meaningful child nodes may increase after this
	// operation, so we need to ensure this cache is updated.
	defer func() { n.sizeCacheValid = false }()

	if vector.Within(p.P(), n.P()) {
		n.data = append(n.data, p)
		return
	}

	x := p.P().X(n.Axis())
	nx := n.P().X(n.Axis())

	if x < nx {
		if n.l == nil {
			n.l = &N{
				depth: n.depth + 1,
				p:     p.P(),
			}
		}
		n.l.Insert(p)
		return
	}
	if n.r == nil {
		n.r = &N{
			depth: n.depth + 1,
			p:     p.P(),
		}
	}
	n.r.Insert(p)
}

// Remove deletes a data point from the node or child nodes. A returned value of
// false indicates the given point was not found.
//
// N.B.: Remove does not actually delete the underlying k-d tree node. Manually
// removing and shifting the nodes is a non-trivial task. We generally expect
// k-d trees to be relatively stable once created, and that insert and remove
// operations are kept at a minimum.
func (n *N) Remove(v vector.V, f func(p point.P) bool) bool {
	// The number of meaningful child nodes may decrease after this
	// operation, so we need to ensure this cache is updated.
	defer func() { n.sizeCacheValid = false }()

	if vector.Within(v, n.P()) {
		for i := range n.data {
			if f(n.data[i]) {
				// Remove the i-th element and set the data to
				// be a shortened slice.
				//
				// N.B.: we can mutate the data slice in this
				// manner only because we are guaranteed to
				// return from the function immediately after,
				// skipping any subsequent iterations.
				n.data[i], n.data[len(n.data)-1] = n.data[len(n.data)-1], nil
				n.data = n.data[:len(n.data)-1]
				return true
			}
		}
	}

	c := n.Child(v)

	return c != nil && c.Remove(v, f)
}

// New returns a new K-D tree node instance.
func New(data []point.P, depth int) *N {
	if len(data) == 0 {
		return nil
	}

	// Sort is not stable -- the order may be shuffled, meaning that while
	// the axis coordinates are in order, the complement dimension is not.
	//
	// That is, give we are sorting on the X-axis,
	//
	//   [(1, 3), (1, 1)]
	//
	// is a valid ordering.
	sorter.Sort(data, vector.D(vector.D(depth)%data[0].P().Dimension()))

	m := len(data) / 2
	p := data[m].P()

	// Find adjacent elements in the sorted list that have the same
	// coordinates as the median, as they all should be in the same node.
	var l int
	var r int
	for i := m; i >= 0 && vector.Within(p, data[i].P()); i-- {
		l = i
	}
	for i := m; i < len(data) && vector.Within(p, data[i].P()); i++ {
		r = i
	}

	l = int(math.Max(0, float64(l)))
	r = int(math.Min(float64(len(data)-1), float64(r)))

	n := &N{
		l:     New(data[:l], depth+1),
		r:     New(data[r+1:], depth+1),
		depth: depth,

		p:    p,
		data: data[l : r+1],
	}

	return n
}

// Points returns all data stored in the node and subnodes.
func Points(n *N) []point.P {
	ps := n.Data()
	if n.L() != nil {
		ps = append(ps, Points(n.L())...)
	}
	if n.R() != nil {
		ps = append(ps, Points(n.R())...)
	}
	return ps
}
