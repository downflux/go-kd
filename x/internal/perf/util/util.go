package util

import (
	"math/rand"
	"runtime"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point/mock"
)

var (
	KRange    = []vector.D{2, 32, 256}
	NRange    = []int{1e3, 1e4, 1e5, 1e6}
	SizeRange = []int{1, 16, 128, 512}
)

func RV(k vector.D, min float64, max float64) vector.V {
	var xs []float64
	for i := 0; i < int(k); i++ {
		xs = append(xs, rand.Float64()*(max-min)+min)
	}
	return vector.V(xs)
}

func Generate(n int, k vector.D) []*mock.P {
	// Generating large number of points in tests will mess with data
	// collection figures. We should ignore these allocs.
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	var ps []*mock.P
	for i := 0; i < n; i++ {
		ps = append(ps, &mock.P{
			X: RV(k, -100, 100),
		})
	}

	return ps
}
