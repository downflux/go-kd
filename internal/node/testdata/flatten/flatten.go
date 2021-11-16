// Package flatten unnests a node and returns an ordered list of node instances.
package flatten

import (
	"github.com/downflux/go-kd/internal/node"
)

func Flatten(n *node.N) []*node.N {
	if n == nil {
		return nil
	}

	ns := []*node.N{n}
	if l := n.L(); l != nil {
		ns = append(Flatten(l), ns...)
	}
	if r := n.R(); r != nil {
		ns = append(Flatten(r), ns...)
	}

	return ns
}
