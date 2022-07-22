package util

import (
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
)

type Cmp[p point.P] point.D

func (c Cmp[p]) Less(a p, b p) bool {
	return a.P().X(point.D(c)) < b.P().X(point.D(c))
}

func Validate[p point.P](t *node.N, data []p) bool {
	open := []*node.N{t}
	for len(open) > 0 {
		var n *node.N
		n, open = open[0], open[1:]

		if n.Left != nil {
			open = append(open, n.Left)
		}
		if n.Right != nil {
			open = append(open, n.Right)
		}

		for i := n.Low; n.Pivot >= 0 && i < n.Pivot; i++ {
			if Cmp[p](n.Axis).Less(data[n.Pivot], data[i]) {
				return false
			}
		}
		for i := n.Pivot; n.Pivot >= 0 && i < n.High; i++ {
			if Cmp[p](n.Axis).Less(data[i], data[n.Pivot]) {
				return false
			}
		}
	}
	return true
}
