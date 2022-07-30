package mock

import (
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point"
)

var _ point.P = P{}

type P struct {
	X    vector.V
	Data string
}

func (p P) P() vector.V { return p.X }

func Equal(a P, b P) bool {
	return a.Data == b.Data && vector.Within(a.P(), b.P())
}

func U(x float64) vector.V   { return vector.V([]float64{x}) }
func V(v []float64) vector.V { return vector.V(v) }
