// Package node defines test-only functions on nodes.
package node

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
