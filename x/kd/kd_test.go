package kd

import (
	"fmt"
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/node/util"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"

	putil "github.com/downflux/go-kd/x/internal/validation/util"
)

func TestNew(t *testing.T) {
	type config struct {
		name string
		k    vector.D
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
		k    vector.D
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
			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("KNN mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
