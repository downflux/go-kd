package kd

import (
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/point"
)

type O[p point.P] struct {
	Data []p
	K    point.D
	N    int
}

type T[p point.P] struct {
	k    point.D
	n    int
	data []p

	root *node.N
}

func New[p point.P](o O[p]) *T[p] {
	data := make([]p, len(o.Data))
	l := copy(data, o.Data)
	if l != len(o.Data) {
		panic("could not copy data into k-D tree")
	}
	if o.K < 1 {
		panic("k-D tree must contain points with non-zero length vectors")
	}
	if o.N < 1 {
		panic("k-D tree minimum leaf node size must be positive")
	}

	t := &T[p]{
		k:    o.K,
		n:    o.N,
		data: data,
		root: node.New[p](node.O[p]{
			Data: data,
			Axis: 0,
			K:    o.K,
			N:    o.N,
			Low:  0,
			High: len(data),
		}),
	}

	return t
}
