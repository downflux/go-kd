package util

import (
	"math"
	"math/rand"
	"runtime"
	"sort"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"
)

var (
	KRange    = []vector.D{2, 32, 256}
	NRange    = []int{1e3, 1e4, 1e5, 1e6}
	SizeRange = []int{1, 16, 128, 512}
)

// Transformer sorts a list of points.
func Transformer(p vector.V) cmp.Option {
	return cmp.Transformer("Sort", func(in []*mock.P) []*mock.P {
		out := append([]*mock.P(nil), in...)
		sort.Sort(L[*mock.P]{
			Data: out,
			P:    p,
		})
		return out
	})
}

type L[T point.P] struct {
	Data []T
	P    vector.V
}

func (l L[T]) Len() int      { return len(l.Data) }
func (l L[T]) Swap(i, j int) { l.Data[i], l.Data[j] = l.Data[j], l.Data[i] }

func (l L[T]) Less(i, j int) bool {
	return vector.SquaredMagnitude(
		vector.Sub(l.Data[i].P(), l.P),
	) < vector.SquaredMagnitude(
		vector.Sub(l.Data[j].P(), l.P),
	)
}

func RH(k vector.D, f float64) hyperrectangle.R {
	min := make([]float64, k)
	max := make([]float64, k)
	for i := vector.D(0); i < k; i++ {
		min[i] = -100 * math.Sqrt(f)
		max[i] = 100 * math.Sqrt(f)
	}
	return *hyperrectangle.New(vector.V(min), vector.V(max))
}

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
