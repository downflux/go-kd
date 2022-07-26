package vector

import (
	"github.com/downflux/go-geometry/epsilon"
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

type Buffer []float64

func (v Buffer) X(d D) float64 { return v[d] }
func (v Buffer) D() D          { return D(len(v)) }

type Comparator D

func (c Comparator) Less(v V, u V) bool { return v.X(D(c)) < u.X(D(c)) }

func SquaredMagnitude(v V) float64 { return vector.SquaredMagnitude(convert(v)) }

func Sub(v V, u V) V { return Buffer(vector.Sub(convert(v), convert(u))) }
func SubBuf(v V, u V, buf Buffer) {
	for i := D(0); i < v.D(); i++ {
		buf[i] = v.X(i) - u.X(i)
	}
}

func Within(v V, u V) bool { return vector.Within(convert(v), convert(u)) }
func WithinEpsilon(v V, u V, e epsilon.E) bool {
	return vector.WithinEpsilon(convert(v), convert(u), e)
}
