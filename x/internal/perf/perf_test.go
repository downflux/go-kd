package perf

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/downflux/go-kd/x/internal/perf/bruteforce"
	"github.com/downflux/go-kd/x/internal/perf/util"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/downflux/go-kd/x/vector"
)

var _ I[*mock.P] = &T[*mock.P]{}
var _ I[*mock.P] = &bruteforce.L[*mock.P]{}

func BenchmarkKNN(b *testing.B) {
	type config struct {
		name string
		t    I[*mock.P]
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
				p:    mock.V(make([]float64, k)),
				knn:  n,
			})

			for _, f := range []float64{0.05, 0.1, 0.25} {
				knn := int(float64(n) * f)

				for _, size := range util.SizeRange {

					configs = append(configs, config{
						name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v/KNN=%v", k, n, size, f),
						t: (*T[*mock.P])(unsafe.Pointer(
							kd.New[*mock.P](kd.O[*mock.P]{
								Data: ps,
								K:    k,
								N:    size,
							}),
						)),
						p:   mock.V(make([]float64, k)),
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
