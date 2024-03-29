package util

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"
	"github.com/downflux/go-kd/point/mock"
	"github.com/google/go-cmp/cmp"
)

type PerfTestSize int

const (
	SizeUnknown PerfTestSize = iota
	SizeUnit
	SizeSmall
	SizeLarge
)

func (s *PerfTestSize) String() string {
	return map[PerfTestSize]string{
		SizeUnit:  "unit",
		SizeSmall: "small",
		SizeLarge: "large",
	}[*s]
}

func (s *PerfTestSize) Set(v string) error {
	size, ok := map[string]PerfTestSize{
		"unit":  SizeUnit,
		"small": SizeSmall,
		"large": SizeLarge,
	}[v]
	if !ok {
		return fmt.Errorf("invalid test size value: %v", v)
	}
	*s = size
	return nil
}

func (s PerfTestSize) F() []float64 {
	return map[PerfTestSize][]float64{
		SizeUnit:  []float64{0.05},
		SizeSmall: []float64{0.05},
		SizeLarge: []float64{0.05, 0.1},
	}[s]
}

func (s PerfTestSize) LeafSize() []int {
	return map[PerfTestSize][]int{
		SizeUnit:  []int{1, 16},
		SizeSmall: []int{1, 32, 512},
		SizeLarge: []int{1, 16, 256},
	}[s]
}

func (s PerfTestSize) N() []int {
	return map[PerfTestSize][]int{
		SizeUnit:  []int{1e3},
		SizeSmall: []int{1e3, 1e4},
		SizeLarge: []int{1e3, 1e4, 1e6},
	}[s]
}

func (s PerfTestSize) K() []vector.D {
	return map[PerfTestSize][]vector.D{
		SizeUnit:  []vector.D{2},
		SizeSmall: []vector.D{2, 16},

		// Large tests phyically cannot store enough point data in
		// memory with high-dimensional data.
		SizeLarge: []vector.D{16},
	}[s]
}

func TrivialFilter(p *mock.P) bool { return true }

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
