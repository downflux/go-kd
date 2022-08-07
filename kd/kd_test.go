package kd

import (
	"fmt"
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-kd/container/bruteforce"
	"github.com/downflux/go-kd/internal/node/util"
	"github.com/downflux/go-kd/point/mock"
	"github.com/downflux/go-kd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	vnd "github.com/downflux/go-geometry/nd/vector"
	putil "github.com/downflux/go-kd/internal/perf/util"
)

func TestNew(t *testing.T) {
	type config struct {
		name string
		k    vnd.D
		n    int

		size int
	}

	var configs []config
	for _, k := range putil.KRange {
		for _, n := range putil.NRange {
			for _, size := range putil.SizeRange {
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
		ps := putil.Generate(c.n, c.k)
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

func TestData(t *testing.T) {
	type config struct {
		name string
		data []*mock.P
		k    vnd.D
		want []*mock.P
	}

	configs := []config{
		{
			name: "Nil",
			data: nil,
			want: nil,
			k:    1,
		},
		{
			name: "Simple",
			data: []*mock.P{
				&mock.P{X: mock.U(1)},
			},
			want: []*mock.P{
				&mock.P{X: mock.U(1)},
			},
			k: 1,
		},
		{
			name: "LR",
			data: []*mock.P{
				&mock.P{X: mock.U(1)},
				&mock.P{X: mock.U(0)},
				&mock.P{X: mock.U(2)},
			},
			want: []*mock.P{
				&mock.P{X: mock.U(1)},
				&mock.P{X: mock.U(0)},
				&mock.P{X: mock.U(2)},
			},
			k: 1,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			kd := New(O[*mock.P]{
				Data: c.data,
				K:    c.k,
				N:    1,
			})
			got := Data(kd)
			if diff := cmp.Diff(c.want, got, cmpopts.SortSlices(
				func(p, q *mock.P) bool {
					return vector.Less(p.P(), q.P())
				})); diff != "" {
				t.Errorf("KNN mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
func TestKNN(t *testing.T) {
	type config struct {
		name string
		k    vnd.D
		n    int
		size int

		knn int
	}

	var configs []config
	for _, k := range putil.KRange {
		for _, n := range putil.NRange {
			for _, size := range putil.SizeRange {
				for _, f := range putil.FRange {
					configs = append(configs, config{
						name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v/KNN=%v", k, n, size, f),
						k:    k,
						n:    n,
						knn:  int(float64(n) * f),
						size: size,
					})
				}
			}
		}
	}

	for _, c := range configs {
		ps := putil.Generate(c.n, c.k)
		t.Run(c.name, func(t *testing.T) {
			p := vnd.V(make([]float64, c.k))

			got := KNN(
				New[*mock.P](O[*mock.P]{
					Data: ps,
					K:    c.k,
					N:    c.size,
				}),
				p,
				c.knn,
				putil.TrivialFilter,
			)
			want := bruteforce.New[*mock.P](ps).KNN(p, c.knn, putil.TrivialFilter)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("KNN mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestRangeSearch(t *testing.T) {
	type config struct {
		name string
		k    vnd.D
		n    int
		size int
		q    hyperrectangle.R
	}

	var configs []config
	for _, k := range putil.KRange {
		for _, n := range putil.NRange {
			for _, size := range putil.SizeRange {
				for _, f := range putil.FRange {
					configs = append(configs, config{
						name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v/Coverage=%v", k, n, size, f),
						k:    k,
						n:    n,
						size: size,
						q:    putil.RH(k, f),
					})
				}
			}
		}
	}

	for _, c := range configs {
		ps := putil.Generate(c.n, c.k)
		t.Run(c.name, func(t *testing.T) {
			got := RangeSearch(
				New[*mock.P](O[*mock.P]{
					Data: ps,
					K:    c.k,
					N:    c.size,
				}),
				c.q,
				putil.TrivialFilter,
			)
			want := bruteforce.New[*mock.P](ps).RangeSearch(c.q, putil.TrivialFilter)

			if diff := cmp.Diff(want, got, putil.Transformer(vnd.V(make([]float64, c.k)))); diff != "" {
				t.Errorf("RangeSearch mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
