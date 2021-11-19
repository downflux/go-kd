package knn

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/node/testdata/flatten"
	"github.com/downflux/go-kd/internal/testdata/generator"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock"
)

func TestPath(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		p    vector.V
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Null",
			n:    nil,
			p:    *vector.New(1, 2),
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
			)
			return config{
				name: "Leaf/Match",
				n:    n,
				p:    n.P(),
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
			)
			return config{
				name: "Leaf/NoMatch",
				n:    n,
				p:    *vector.New(0, 2),
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
			)
			return config{
				name: "AssertLeafFirst",
				n:    n,
				p:    *vector.New(1, 2),
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
			)
			return config{
				name: "AssertAlwaysLeaf",
				n:    n,
				p:    *vector.New(2, 2),
				want: []*node.N{
					n.R(),
					n,
				},
			}
		}(),
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := path(c.n, c.p)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(node.N{}, mock.P{})); diff != "" {
				t.Errorf("queue() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

type nl struct {
	ns []*node.N
	p  vector.V
}

func (s *nl) Dist(i int) float64 { return vector.Magnitude(vector.Sub(s.ns[i].P(), s.p)) }
func (s *nl) Len() int           { return len(s.ns) }
func (s *nl) Less(i, j int) bool { return s.Dist(i) < s.Dist(j) }
func (s *nl) Swap(i, j int)      { s.ns[i], s.ns[j] = s.ns[j], s.ns[i] }

func sortNodes(n *node.N, p vector.V) []*node.N {
	s := &nl{
		ns: flatten.Flatten(n),
		p:  p,
	}

	sort.Sort(s)

	return s.ns
}

func TestKNN(t *testing.T) {
	const k = 1000

	type config struct {
		name string
		n    *node.N
		p    vector.V
		k    int
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Null",
			n:    nil,
			p:    *vector.New(1, 2),
			k:    1,
			want: []*node.N{},
		},
		{
			name: "Null/k=0",
			n: node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				0,
			),
			p:    *vector.New(1, 2),
			k:    0,
			want: []*node.N{},
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
			)
			return []config{
				config{
					name: "Trivial/Near",
					n:    n,
					p:    n.P(),
					k:    1,
					want: []*node.N{n},
				},
				config{
					name: "Trivial/Far",
					n:    n,
					p:    *vector.New(1000, 1000),
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
			n := node.New(ps, 0)

			cs := []config{
				config{
					name: "Multiple/k=1/Near",
					n:    n,
					p:    *vector.New(4, 39),
					k:    1,
					want: sortNodes(n, *vector.New(4, 39))[:1],
				},
			}

			for i := 0; i < k; i++ {
				k := rand.Intn(len(ps)) + 1
				p := generator.V(2)
				cs = append(
					cs,
					config{
						name: fmt.Sprintf("Multiple/k=%v/%v", k, i),
						n:    n,
						p:    p,
						k:    k,
						want: sortNodes(n, p)[:k],
					},
				)
			}

			return cs
		}()...,
	)

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := KNN(c.n, c.p, c.k)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(node.N{}, mock.P{})); diff != "" {
				t.Errorf("KNN() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
