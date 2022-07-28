package correctness

import (
	"fmt"
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/validation/util"
	"github.com/downflux/go-kd/x/kd"
	"github.com/downflux/go-kd/x/kd/mock/bruteforce"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"
)

func TestKNN(t *testing.T) {
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
		t.Run(c.name, func(t *testing.T) {
			knn := int(float64(c.n) * 0.1)
			p := vector.V(make([]float64, c.k))

			got := kd.KNN(
				kd.New[*mock.P](kd.O[*mock.P]{
					Data: ps,
					K:    c.k,
					N:    c.size,
				}),
				p,
				knn,
			)
			want := bruteforce.New[*mock.P](ps).KNN(p, knn)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("KNN mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
