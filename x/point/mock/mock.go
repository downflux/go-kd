package mock

import (
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/vector"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

var _ point.P = P{}

var _ vector.V = V{}
var _ vector.V = U(0)

type V vnd.V

func (v V) D() vector.D          { return vector.D(vnd.V(v).Dimension()) }
func (v V) X(d vector.D) float64 { return vnd.V(v).X(vnd.D(d)) }

type U float64

func (u U) D() vector.D          { return vector.D(1) }
func (u U) X(d vector.D) float64 { return float64(u) }

type P struct {
	X    vector.V
	Data string
}

func (p P) P() vector.V { return p.X }

func Equal(a P, b P) bool {
	if a.Data != b.Data {
		return false
	}
	if a.P().D() != b.P().D() {
		return false
	}

	for i := vector.D(0); i < a.P().D(); i++ {
		if !epsilon.Within(a.P().X(i), b.P().X(i)) {
			return false
		}
	}
	return true
}
