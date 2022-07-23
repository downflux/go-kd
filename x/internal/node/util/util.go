package util

import (
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
)

func Map[T point.P](n node.N[T], f func(n node.N[T])) {
	open := []node.N[T]{n}
	for len(open) > 0 {
		var n node.N[T]
		n, open = open[0], open[1:]

		if n.Nil() {
			continue
		}

		if !n.L().Nil() {
			open = append(open, n.L())
		}
		if !n.R().Nil() {
			open = append(open, n.R())
		}

		f(n)
	}
}

type Cmp point.D

func (c Cmp) Less(a point.V, b point.V) bool {
	return a.X(point.D(c)) < b.X(point.D(c))
}

func Validate[T point.P](t node.N[T]) bool {
	equal := true
	f := func(n node.N[T]) {
		if n.Nil() {
			return
		}

		if !n.L().Nil() {
			for _, p := range n.L().Data() {
				equal = equal && Cmp(n.Axis()).Less(p.P(), n.Pivot())
			}
		}
		if !n.R().Nil() {
			for _, p := range n.R().Data() {
				equal = equal && Cmp(n.Axis()).Less(n.Pivot(), p.P())
			}
		}
	}
	Map[T](t, f)
	return equal
}
