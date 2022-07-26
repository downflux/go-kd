package knn

import (
	"math"

	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/point/pq"
	"github.com/downflux/go-kd/x/vector"
)

func path[T point.P](n node.N[T], p vector.V) []node.N[T] {
	if n.Nil() {
		return nil
	}
	if n.Leaf() {
		return []node.N[T]{n}
	}

	// Note that we are bypassing the v == n.Pivot() stop condition check --
	// we are always continuing to the leaf ndoe. This is necessary for
	// finding multiple closest neighbors, as we care about points in the
	// tree which do not have to coincide with the point coordinates.
	if vector.Comparator(n.Axis()).Less(p, n.Pivot()) {
		return append(path(n.L(), p), n)
	}
	return append(path(n.R(), p), n)
}

func KNN[T point.P](n node.N[T], p vector.V, k int) []T {
	q := pq.New[T](k)
	knn(n, p, q, vector.Buffer(make([]float64, p.D())))

	ps := make([]T, q.Len())
	for i := q.Len() - 1; i >= 0; i-- {
		ps[i] = q.Pop()
	}
	return ps
}

func knn[T point.P](n node.N[T], p vector.V, q *pq.PQ[T], buf vector.Buffer) {
	for _, n := range path[T](n, p) {
		for _, datum := range n.Data() {
			vector.SubBuf(p, datum.P(), buf)
			if d := vector.SquaredMagnitude(buf); !q.Full() || d < q.Priority() {
				q.Push(datum, d)
			}
		}

		if !n.Leaf() {
			if q.Priority() > math.Abs(p.X(n.Axis())-n.Pivot().X(n.Axis())) {
				if vector.Comparator(n.Axis()).Less(p, n.Pivot()) {
					knn(n.R(), p, q, buf)
				} else {
					knn(n.L(), p, q, buf)
				}
			}
		}
	}
}
