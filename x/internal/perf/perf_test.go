package perf

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/container"
	"github.com/downflux/go-kd/x/container/bruteforce"
	"github.com/downflux/go-kd/x/internal/perf/util"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point/mock"

	ckd "github.com/downflux/go-kd/x/container/kd"
)

func BenchmarkNew(b *testing.B) {
	type config struct {
		name string
		k    vector.D
		n    int

		size int
	}

	var configs []config
	for _, k := range util.KRange {
		for _, n := range util.NRange {
			for _, size := range util.SizeRange {
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
		ps := util.Generate(c.n, c.k)
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				kd.New[*mock.P](kd.O[*mock.P]{
					Data: ps,
					K:    c.k,
					N:    c.size,
				})
			}
		})
	}
}

func BenchmarkKNN(b *testing.B) {
	type config struct {
		name string
		t    container.I[*mock.P]
		p    vector.V
		knn  int
	}

	var configs []config
	for _, k := range util.KRange {
		for _, n := range util.NRange {
			ps := util.Generate(n, k)

			// Brute force approach sorts all data, meaning that the
			// KNN factor does not matter.
			configs = append(configs, config{
				name: fmt.Sprintf("BruteForce/K=%v/N=%v", k, n),
				t:    bruteforce.New[*mock.P](ps),
				p:    vector.V(make([]float64, k)),
				knn:  n,
			})

			for _, f := range []float64{0.05, 0.1, 0.25} {
				knn := int(float64(n) * f)

				for _, size := range util.SizeRange {
					configs = append(configs, config{
						name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v/KNN=%v", k, n, size, f),
						t: (*ckd.T[*mock.P])(unsafe.Pointer(
							kd.New[*mock.P](kd.O[*mock.P]{
								Data: ps,
								K:    k,
								N:    size,
							}),
						)),
						p:   vector.V(make([]float64, k)),
						knn: knn,
					})
				}

			}
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.t.KNN(c.p, c.knn)
			}
		})
	}
}

func BenchmarkRangeSearch(b *testing.B) {
	type config struct {
		name string
		t    container.I[*mock.P]
		q    hyperrectangle.R
	}

	var configs []config
	for _, k := range util.KRange {
		for _, n := range util.NRange {
			ps := util.Generate(n, k)

			// Brute force approach sorts all data, meaning that the
			// query range factor does not matter.
			configs = append(configs, config{
				name: fmt.Sprintf("BruteForce/K=%v/N=%v", k, n),
				t:    bruteforce.New[*mock.P](ps),
				q:    util.RH(k, 1),
			})

			for _, f := range []float64{0.05, 0.1, 0.25} {
				q := util.RH(k, f)

				for _, size := range util.SizeRange {
					configs = append(configs, config{
						name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v/Coverage=%v", k, n, size, f),
						t: (*ckd.T[*mock.P])(unsafe.Pointer(
							kd.New[*mock.P](kd.O[*mock.P]{
								Data: ps,
								K:    k,
								N:    size,
							}),
						)),
						q: q,
					})
				}
			}
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.t.RangeSearch(c.q)
			}
		})
	}
}
