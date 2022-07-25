package kd

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"

	"github.com/downflux/go-kd/x/internal/node/util"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/downflux/go-kd/x/vector"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

var (
	kRange    = []vector.D{2, 10, 100, 500}
	nRange    = []int{1e3, 1e4, 1e5, 1e6}
	sizeRange = []int{1, 4, 16, 64}
)

func rv(k vector.D, min float64, max float64) mock.V {
	var xs []float64
	for i := 0; i < int(k); i++ {
		xs = append(xs, rand.Float64()*(max-min)+min)
	}
	return mock.V(vnd.V(xs))
}

func generate(n int, k vector.D) []*mock.P {
	// Generating large number of points in tests will mess with data
	// collection figures. We should ignore these allocs.
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	var ps []*mock.P
	for i := 0; i < n; i++ {
		ps = append(ps, &mock.P{
			X: rv(k, -100, 100),
		})
	}

	return ps
}

func TestNew(t *testing.T) {
	type config struct {
		name string
		k    vector.D
		n    int

		size int
	}

	var configs []config
	for _, k := range kRange {
		for _, n := range nRange {
			for _, size := range sizeRange {
				configs = append(configs, config{
					name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v", k, n, size),
					k:    k,
					n:    n,
					size: size,
				})
			}
		}
	}

	for _, c := range configs {
		ps := generate(c.n, c.k)
		t.Run(c.name, func(t *testing.T) {
			tree := New[*mock.P](O[*mock.P]{
				Data: ps,
				K:    c.k,
				N:    c.size,
			})
			if !util.Validate(tree.root) {
				t.Errorf("validate() = %v, want = %v", false, true)
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	type config struct {
		name string
		k    vector.D
		n    int

		size int
	}

	var configs []config
	for _, k := range kRange {
		for _, n := range nRange {
			for _, size := range sizeRange {
				configs = append(configs, config{
					name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v", k, n, size),
					k:    k,
					n:    n,
					size: size,
				})
			}
		}
	}

	for _, c := range configs {
		ps := generate(c.n, c.k)
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New[*mock.P](O[*mock.P]{
					Data: ps,
					K:    c.k,
					N:    c.size,
				})
			}
		})
	}
}
