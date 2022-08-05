package perf

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/container"
	"github.com/downflux/go-kd/container/bruteforce"
	"github.com/downflux/go-kd/container/kyroy"
	"github.com/downflux/go-kd/internal/perf/util"
	"github.com/downflux/go-kd/kd"
	"github.com/downflux/go-kd/point/mock"

	ckd "github.com/downflux/go-kd/container/kd"
)

func BenchmarkNew(b *testing.B) {
	type config struct {
		name string
		k    vector.D
		n    int

		// kyroy implementation does not take a leaf-size parameter.
		kyroy bool

		size int
	}

	var configs []config
	for _, k := range util.BenchmarkKRange {
		for _, n := range util.BenchmarkNRange {
			configs = append(configs, config{
				name:  fmt.Sprintf("kyroy/K=%v/N=%v", k, n),
				k:     k,
				n:     n,
				kyroy: true,
			})
			for _, size := range util.BenchmarkSizeRange {
				configs = append(configs, config{
					name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v", k, n, size),
					k:    k,
					n:    n,
					size: size,
				})
			}
		}
	}

	for _, c := range configs {
		ps := util.Generate(c.n, c.k)

		if c.kyroy {
			b.Run(c.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					kyroy.New[*mock.P](ps)
				}
			})
		} else {
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
}

func BenchmarkKNN(b *testing.B) {
	type config struct {
		name string
		t    container.C[*mock.P]
		p    vector.V
		knn  int
	}

	var configs []config
	for _, k := range util.BenchmarkKRange {
		for _, n := range util.BenchmarkNRange {
			ps := util.Generate(n, k)

			// Brute force approach sorts all data, meaning that the
			// KNN factor does not matter.
			configs = append(configs, config{
				name: fmt.Sprintf("BruteForce/K=%v/N=%v", k, n),
				t:    bruteforce.New[*mock.P](ps),
				p:    vector.V(make([]float64, k)),
				knn:  n,
			})

			for _, f := range util.BenchmarkFRange {
				knn := int(float64(n) * f)

				// kyroy implementation does not take a
				// leaf-size parameter.
				configs = append(configs, config{
					name: fmt.Sprintf("kyroy/K=%v/N=%v/KNN=%v", k, n, f),
					t:    kyroy.New[*mock.P](ps),
					p:    vector.V(make([]float64, k)),
					knn:  knn,
				})

				for _, size := range util.BenchmarkSizeRange {
					configs = append(configs, config{
						name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v/KNN=%v", k, n, size, f),
						t: (*ckd.KD[*mock.P])(unsafe.Pointer(
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
				c.t.KNN(c.p, c.knn, util.TrivialFilter)
			}
		})
	}
}

func BenchmarkRangeSearch(b *testing.B) {
	type config struct {
		name string
		t    container.C[*mock.P]
		q    hyperrectangle.R
	}

	var configs []config
	for _, k := range util.BenchmarkKRange {
		for _, n := range util.BenchmarkNRange {
			ps := util.Generate(n, k)

			// Brute force approach sorts all data, meaning that the
			// query range factor does not matter.
			configs = append(configs, config{
				name: fmt.Sprintf("BruteForce/K=%v/N=%v", k, n),
				t:    bruteforce.New[*mock.P](ps),
				q:    util.RH(k, 1),
			})

			for _, f := range util.BenchmarkFRange {
				q := util.RH(k, f)

				// kyroy implementation does not take a
				// leaf-size parameter.
				configs = append(configs, config{
					name: fmt.Sprintf("kyroy/K=%v/N=%v/Coverage=%v", k, n, f),
					t:    kyroy.New[*mock.P](ps),
					q:    q,
				})

				for _, size := range util.BenchmarkSizeRange {
					configs = append(configs, config{
						name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v/Coverage=%v", k, n, size, f),
						t: (*ckd.KD[*mock.P])(unsafe.Pointer(
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
				c.t.RangeSearch(c.q, util.TrivialFilter)
			}
		})
	}
}
