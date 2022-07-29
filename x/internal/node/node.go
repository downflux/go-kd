package node

import (
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point"
)

type N[T point.P] interface {
	// L consists of points strictly less than the current pivot for the
	// current axis.
	L() N[T]

	// R consists of points greater than or equal to the current pivot for
	// the current axis.
	R() N[T]

	// Data returns the points stored in the current node -- note that this
	// does not include data from child nodes.
	Data() []T

	Pivot() vector.V
	K() vector.D
	Axis() vector.D

	Leaf() bool
	Nil() bool
}
