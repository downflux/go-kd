package axis

import (
	"github.com/downflux/go-geometry/vector"
)

type Type int

const (
	Axis_X Type = iota
	Axis_Y
)

func A(depth int) Type             { return Type(depth % 2) }
func X(v vector.V, a Type) float64 { return []float64{v.X(), v.Y()}[a] }
