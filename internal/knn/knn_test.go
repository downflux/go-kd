package knn

import (
	"testing"

	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/node/tree"
	"github.com/downflux/go-kd/point"
	"github.com/downflux/go-kd/point/mock"
	"github.com/downflux/go-kd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

func TestKNN(t *testing.T) {
	type config[T point.P] struct {
		name string
		n    node.N[T]
		p    vnd.V
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
			p:    mock.V(*vnd.New(100, 200)),
			k:    100,
			want: []mock.P{},
		},
		{
			name: "SmallD",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.U(0.1)},
					mock.P{X: mock.U(0.01)},
				},
				K: 1,
				N: 10,
			}),
			p: mock.U(0),
			k: 1,
			want: []mock.P{
				mock.P{X: mock.U(0.01)},
			},
		},
		{
			name: "Simple",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.U(10)},
				},
				K: 1,
				N: 10,
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
					mock.P{X: mock.V(*vnd.New(100, 1))},
				},
				K: 2,
				N: 1,
			}),
			p: mock.V(*vnd.New(0, -100)),
			k: 100,
			want: []mock.P{
				mock.P{X: mock.V(*vnd.New(100, 1))},
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
				K: 1,
				N: 1,
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
				K: 1,
				N: 1,
			}),
			p: mock.U(100),
			k: 2,
			want: []mock.P{
				// We don't care what data we get here, as long
				// as it's two of the input set. The ordering
				// was matched manually.
				mock.P{X: mock.U(99), Data: "C"},
				mock.P{X: mock.U(99), Data: "B"},
			},
		},
		{
			name: "Simple/MultiK/2D/Degenerate",
			n: tree.New[mock.P](tree.O[mock.P]{
				Data: []mock.P{
					mock.P{X: mock.V(*vnd.New(99, 100)), Data: "A"},
					mock.P{X: mock.V(*vnd.New(99, 100)), Data: "B"},
					mock.P{X: mock.V(*vnd.New(99, 100)), Data: "C"},
					mock.P{X: mock.V(*vnd.New(99, 100)), Data: "D"},
				},
				K: 2,
				N: 1,
			}),
			p: mock.V(*vnd.New(0, 0)),
			k: 2,
			want: []mock.P{
				mock.P{X: mock.V(*vnd.New(99, 100)), Data: "C"},
				mock.P{X: mock.V(*vnd.New(99, 100)), Data: "B"},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := KNN(c.n, c.p, c.k, func(mock.P) bool { return true })
			if diff := cmp.Diff(c.want, got, cmpopts.SortSlices(
				func(p, q mock.P) bool {
					return vector.Less(p.P(), q.P())
				})); diff != "" {
				t.Errorf("KNN mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
