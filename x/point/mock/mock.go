package mock

import (
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-kd/x/point"
)

var _ point.P = P{}

var _ point.V = V{}
var _ point.V = U(0)

type V vector.V

func (v V) D() point.D { return point.D(2) }
func (v V) X(d point.D) float64 {
	return map[point.D]float64{
		0: vector.V(v).X(),
		1: vector.V(v).Y(),
	}[d]
}

type U float64

func (u U) D() point.D          { return point.D(1) }
func (u U) X(d point.D) float64 { return float64(u) }

type P struct {
	X    point.V
	Data string
}

func (p P) P() point.V { return p.X }

func Equal(a P, b P) bool {
	if a.Data != b.Data {
		return false
	}
	if a.P().D() != b.P().D() {
		return false
	}

	for i := point.D(0); i < a.P().D(); i++ {
		if !epsilon.Within(a.P().X(i), b.P().X(i)) {
			return false
		}
	}
	return true
}
