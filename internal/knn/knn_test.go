package knn

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/downflux/go-geometry/vector"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"

	test "github.com/downflux/go-kd/internal/node/test"
	mock "github.com/downflux/go-kd/internal/point/mock"
)

const (
	tolerance = 1e-10
)

func TestPath(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		v    vector.V
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Null",
			n:    nil,
			v:    *vector.New(1, 2),
			want: nil,
		},

		// Ensure a path to a leaf node is returned if the query
		// point is the same as the leaf.
		func() config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				0,
				tolerance,
			)
			return config{
				name: "Leaf/Match",
				n:    n,
				v:    n.V(),
				want: []*node.N{n},
			}
		}(),

		// Ensure a path to a leaf node is returned even if the
		// query point is the different from the leaf.
		func() config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				0,
				tolerance,
			)
			return config{
				name: "Leaf/NoMatch",
				n:    n,
				v:    *vector.New(0, 2),
				want: []*node.N{n},
			}
		}(),

		// Ensure the generated path starts from the leaf node.
		func() config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
					*mock.New(*vector.New(2, 2), ""),
				},
				0,
				tolerance,
			)
			return config{
				name: "AssertLeafFirst",
				n:    n,
				v:    *vector.New(1, 2),
				want: []*node.N{
					n.L(),
					n,
				},
			}
		}(),

		// Ensure that even if a non-leaf node matches the query, the
		// path continues to generate towards a leaf node.
		func() config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
					*mock.New(*vector.New(2, 2), ""),
					*mock.New(*vector.New(3, 2), ""),
				},
				0,
				tolerance,
			)
			return config{
				name: "AssertAlwaysLeaf",
				n:    n,
				v:    *vector.New(2, 2),
				want: []*node.N{
					n.R(),
					n,
				},
			}
		}(),
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := path(c.n, c.v, tolerance)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(node.N{}, vector.V{}, mock.P{})); diff != "" {
				t.Errorf("queue() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

type nl struct {
	ns []*node.N
	v  vector.V
}

func (s *nl) Dist(i int) float64 { return vector.Magnitude(vector.Sub(s.ns[i].V(), s.v)) }
func (s *nl) Len() int           { return len(s.ns) }
func (s *nl) Less(i, j int) bool { return s.Dist(i) < s.Dist(j) }
func (s *nl) Swap(i, j int)      { s.ns[i], s.ns[j] = s.ns[j], s.ns[i] }

func sortNodes(n *node.N, v vector.V) []*node.N {
	s := &nl{
		ns: test.Flatten(n),
		v:  v,
	}

	sort.Sort(s)

	return s.ns
}

func rn() float64  { return rand.Float64()*200 - 100 }
func rv() vector.V { return *vector.New(rn(), rn()) }
func rp() point.P  { return *mock.New(rv(), "") }

func TestKNN(t *testing.T) {
	const k = 1000

	type config struct {
		name string
		n    *node.N
		v    vector.V
		k    int
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Null",
			n:    nil,
			v:    *vector.New(1, 2),
			k:    1,
			want: nil,
		},
		{
			name: "Null/k=0",
			n: node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				0,
				tolerance,
			),
			v:    *vector.New(1, 2),
			k:    0,
			want: nil,
		},
	}

	testConfigs = append(
		testConfigs,

		// Ensure for a single element, the input query point does not
		// matter.
		func() []config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				0,
				tolerance,
			)
			return []config{
				config{
					name: "Trivial/Near",
					n:    n,
					v:    n.V(),
					k:    1,
					want: []*node.N{n},
				},
				config{
					name: "Trivial/Far",
					n:    n,
					v:    *vector.New(1000, 1000),
					k:    1,
					want: []*node.N{n},
				},
			}
		}()...,
	)

	testConfigs = append(
		testConfigs,

		func() []config {
			ps := []point.P{
				*mock.New(*vector.New(1, 60), "A"),
				*mock.New(*vector.New(2, 42), "B"),
				*mock.New(*vector.New(3, 40), "C"),
				*mock.New(*vector.New(4, 39), "D"),
				*mock.New(*vector.New(5, 20), "E"),
			}

			//     C
			//    / \
			//   A   D
			//  /   /
			// B   E
			n := node.New(ps, 0, tolerance)

			cs := []config{
				config{
					name: "Multiple/k=1/Near",
					n:    n,
					v:    *vector.New(4, 39),
					k:    1,
					want: sortNodes(n, *vector.New(4, 39))[:1],
				},
			}

			for i := 0; i < k; i++ {
				k := rand.Intn(len(ps)) + 1
				v := rv()
				cs = append(
					cs,
					config{
						name: fmt.Sprintf("Multiple/k=%v/%v", k, i),
						n:    n,
						v:    v,
						k:    k,
						want: sortNodes(n, v)[:k],
					},
				)
			}

			return cs
		}()...,
	)

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := KNN(c.n, c.v, c.k, tolerance)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(node.N{}, vector.V{}, mock.P{})); diff != "" {
				t.Errorf("KNN() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
