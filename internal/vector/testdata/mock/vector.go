package vector

import (
	"github.com/downflux/go-kd/vector"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

type V vnd.V

func New(xs... float64) *V {
	v := V(*vnd.New(xs...))
	return &v
}

func (v V) Dimension() vector.D { return vector.D(vnd.V(v).Dimension()) }
func (v V) X(i vector.D) float64 { return vnd.V(v).X(vnd.D(i)) }
