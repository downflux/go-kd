package knn

import (
	"math"

	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/point"
	"github.com/downflux/go-kd/vector"
	"github.com/downflux/go-pq/pq"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

func path[T point.P](n node.N[T], p vnd.V) []node.N[T] {
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

func KNN[T point.P](n node.N[T], p vnd.V, k int, f func(p T) bool) []T {
	q := pq.New[T](k, pq.PMax)
	knn(n, p, q, vnd.M(make([]float64, p.Dimension())), f)

	ps := make([]T, q.Len())
	for i := q.Len() - 1; i >= 0; i-- {
		ps[i], _ = q.Pop()
	}
	return ps
}

func knn[T point.P](n node.N[T], p vnd.V, q *pq.PQ[T], buf vnd.M, f func(p T) bool) {
	for _, n := range path[T](n, p) {
		for _, datum := range n.Data() {
			buf.Copy(p)
			buf.Sub(datum.P())

			if d := vnd.SquaredMagnitude(buf.V()); (!q.Full() || d < q.Priority()) && f(datum) {
				q.Push(datum, d)
			}
		}

		if !n.Leaf() {
			buf.Copy(p)
			buf.Sub(n.Pivot())

			if q.Priority() > math.Pow(buf.X(n.Axis()), 2) {
				if vector.Comparator(n.Axis()).Less(p, n.Pivot()) {
					knn(n.R(), p, q, buf, f)
				} else {
					knn(n.L(), p, q, buf, f)
				}
			}
		}
	}
}
