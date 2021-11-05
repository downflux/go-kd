package vector

import (
	"github.com/downflux/go-geometry/vector"
)

type V struct {
	v vector.V
}

func New(x, y float64) *V { return &V{v: *vector.New(x, y)} }

func (v V) D() int { return 2 }
func (v V) X(i int) float64 {
	if i == 0 {
		return v.v.X()
	}
	return v.v.Y()
}
