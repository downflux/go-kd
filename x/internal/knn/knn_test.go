package knn

import (
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/internal/node/tree"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"
)

func TestKNN(t *testing.T) {
	type config[T point.P] struct {
		name string
		n    node.N[T]
		p    vector.V
		k    int
		want []T
	}

	configs := []config[mock.P]{
		{
			name: "Trivial",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: nil,
				K:    1,
				N:    10,
			}),
			p:    mock.V(*vector.New(100, 200)),
			k:    100,
			want: []mock.P{},
		},
		{
			name: "Simple",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.U(10)},
				},
				Low:  0,
				High: 1,
				K:    1,
				N:    10,
			}),
			p: mock.U(-1000),
			k: 100,
			want: []mock.P{
				mock.P{X: mock.U(10)},
			},
		},
		{
			name: "Simple/2D",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.V(*vector.New(100, 1))},
				},
				Low:  0,
				High: 1,
				K:    2,
				N:    1,
			}),
			p: mock.V(*vector.New(0, -100)),
			k: 100,
			want: []mock.P{
				mock.P{X: mock.V(*vector.New(100, 1))},
			},
		},
		{
			name: "Simple/MultiK",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.U(101)},
					mock.P{X: mock.U(102)},
					mock.P{X: mock.U(103)},
					mock.P{X: mock.U(99)},
				},
				Low:  0,
				High: 4,
				K:    1,
				N:    1,
			}),
			p: mock.U(100),
			k: 2,
			want: []mock.P{
				mock.P{X: mock.U(101)},
				mock.P{X: mock.U(99)},
			},
		},
		{
			name: "Simple/MultiK/Degenerate",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.U(99), Data: "A"},
					mock.P{X: mock.U(99), Data: "B"},
					mock.P{X: mock.U(99), Data: "C"},
					mock.P{X: mock.U(99), Data: "D"},
				},
				Low:  0,
				High: 4,
				K:    1,
				N:    1,
			}),
			p: mock.U(100),
			k: 2,
			want: []mock.P{
				mock.P{X: mock.U(99), Data: "C"},
				mock.P{X: mock.U(99), Data: "D"},
			},
		},
		{
			name: "Simple/MultiK/2D/Degenerate",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.V(*vector.New(99, 100)), Data: "A"},
					mock.P{X: mock.V(*vector.New(99, 100)), Data: "B"},
					mock.P{X: mock.V(*vector.New(99, 100)), Data: "C"},
					mock.P{X: mock.V(*vector.New(99, 100)), Data: "D"},
				},
				Low:  0,
				High: 4,
				K:    2,
				N:    1,
			}),
			p: mock.V(*vector.New(0, 0)),
			k: 2,
			want: []mock.P{
				mock.P{X: mock.V(*vector.New(99, 100)), Data: "B"},
				mock.P{X: mock.V(*vector.New(99, 100)), Data: "A"},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := KNN(c.n, c.p, c.k)
			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("KNN mismatch (-want +got):\n%v", diff)
			}

			for i, p := range got {
				if !mock.Equal(p, c.want[i]) {
					t.Errorf("Equal() = %v, want = %v", false, true)
				}
			}
		})
	}
}
