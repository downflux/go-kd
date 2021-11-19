package generator

import (
	"fmt"
	"math/rand"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock"
)

const Interval = 200

func F() float64 { return rand.Float64()*Interval - (Interval / 2) }

func V(d vector.D) vector.V {
	xs := make([]float64, d)
	for i := 0; i < int(d); i++ {
		xs[i] = F()
	}

	return vector.V(xs)
}

func P(n int, d vector.D) []point.P {
	ps := make([]point.P, n)
	for i := 0; i < n; i++ {
		ps[i] = *mock.New(V(d), fmt.Sprintf("Random-%v", i))
	}
	return ps
}
