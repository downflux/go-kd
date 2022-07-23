package node

import (
	"github.com/downflux/go-kd/x/point"
)

type N[T point.P] interface {
	L() N[T]
	R() N[T]
	Data() []T
	Pivot() point.V
	Axis() point.D
	Leaf() bool
	Nil() bool
}
