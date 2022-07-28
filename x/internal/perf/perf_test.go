package perf

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/perf/util"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/kd/mock"
	"github.com/downflux/go-kd/x/kd/mock/bruteforce"
	"github.com/downflux/go-kd/x/kd/mock/wrapper"

	pmock "github.com/downflux/go-kd/x/point/mock"
)

func BenchmarkKNN(b *testing.B) {
	type config struct {
		name string
		t    mock.I[*pmock.P]
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
				t:    bruteforce.New[*pmock.P](ps),
				p:    vector.V(make([]float64, k)),
				knn:  n,
			})

			for _, f := range []float64{0.05, 0.1, 0.25} {
				knn := int(float64(n) * f)

				for _, size := range util.SizeRange {

					configs = append(configs, config{
						name: fmt.Sprintf("Real/K=%v/N=%v/LeafSize=%v/KNN=%v", k, n, size, f),
						t: (*wrapper.T[*pmock.P])(unsafe.Pointer(
							kd.New[*pmock.P](kd.O[*pmock.P]{
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
