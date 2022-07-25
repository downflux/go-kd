package vector

import (
	"github.com/downflux/go-geometry/nd/vector"
)

type D int

const (
	AXIS_X D = iota
	AXIS_Y
	AXIS_Z
	AXIS_W
)

type V interface {
	X(d D) float64
	D() D
}

func convert(v V) vector.V {
	u := make([]float64, v.D())
	for i := D(0); i < v.D(); i++ {
		u[i] = v.X(i)
	}
	return u
}

type internal vector.V

func (v internal) X(d D) float64 { return vector.V(v).X(vector.D(d)) }
func (v internal) D() D          { return D(vector.V(v).Dimension()) }

type Comparator D

func (c Comparator) Less(v V, u V) bool { return v.X(D(c)) < u.X(D(c)) }

func SquaredMagnitude(v V) float64 { return vector.SquaredMagnitude(convert(v)) }
func Sub(v V, u V) V               { return internal(vector.Sub(convert(v), convert(u))) }
