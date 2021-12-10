package rangesearch

import (
	"fmt"
	"sort"
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/internal/testdata/generator"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock/simple"
)

type lex []point.P

func (l lex) Len() int { return len(l) }
func (l lex) Less(i, j int) bool {
	p := l[i]
	q := l[j]

	for d := vector.D(0); d < p.P().Dimension(); d++ {
		if p.P().X(d) < q.P().X(d) {
			return true
		}
		if p.P().X(d) > q.P().X(d) {
			return false
		}
	}
	return p.(mock.P).Hash() < q.(mock.P).Hash()
}
func (l lex) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func TestSearchPoints(t *testing.T) {
	const n = 0

	type config struct {
		name string
		ps   []point.P
		r    hyperrectangle.R
	}

	testConfigs := []config{
		{
			name: "Trivial",
			ps: []point.P{
				*mock.New(*vector.New(1, 2), ""),
			},
			r: *hyperrectangle.New(
				*vector.New(0, 0),
				*vector.New(2, 3),
			),
		},
		{
			name: "Trivial/Disjoint",
			ps: []point.P{
				*mock.New(*vector.New(1, 2), ""),
			},
			r: *hyperrectangle.New(
				*vector.New(-1, -1),
				*vector.New(0, 0),
			),
		},
		{
			name: "Multipoint",
			ps: []point.P{
				*mock.New(*vector.New(-63.415016709218314, -14.328583638638435), ""),
				*mock.New(*vector.New(-55.54211659864245, 36.21566247851419), ""),
				*mock.New(*vector.New(51.69698229056947, -37.69551113789503), ""),
				*mock.New(*vector.New(60.21100853053224, 46.04629545896165), ""),
				*mock.New(*vector.New(86.5692857036868, 48.369791998364605), ""),
			},
			r: *hyperrectangle.New(
				*vector.New(79.39839151237453, 36.53069760264876),
				*vector.New(95.78587111533753, 84.44245178434537),
			),
		},
	}

	for i := 0; i < n; i++ {
		testConfigs = append(
			testConfigs,
			config{
				name: fmt.Sprintf("Random/%v", i),
				ps:   generator.P(n, 2),
				r: *hyperrectangle.New(
					generator.V(2),
					generator.V(2),
				),
			},
		)
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			n := node.New(c.ps, 0)

			var got []point.P
			for _, n := range Search(n, c.r) {
				got = append(got, n.Data()...)
			}

			var want []point.P
			for _, p := range c.ps {
				if c.r.In(p.P()) {
					want = append(want, p)
				}
			}

			sort.Sort(lex(got))
			sort.Sort(lex(want))

			if diff := cmp.Diff(
				want,
				got,
				cmp.AllowUnexported(
					mock.P{},
				),
			); diff != "" {
				t.Errorf("Search() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		r    hyperrectangle.R
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Trivial",
			n:    nil,
			r:    *hyperrectangle.New(*vector.New(0, 0), *vector.New(1, 2)),
			want: nil,
		},
	}

	testConfigs = append(
		testConfigs,
		func() []config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				0,
			)

			return []config{
				{
					name: "Trivial/Embedded",
					n:    n,
					r: *hyperrectangle.New(
						*vector.New(0, 0),
						*vector.New(2, 3),
					),
					want: []*node.N{n},
				},
				{
					name: "Trivial/Disjoint",
					n:    n,
					r: *hyperrectangle.New(
						*vector.New(2, 3),
						*vector.New(3, 4),
					),
					want: nil,
				},
			}
		}()...,
	)

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := Search(c.n, c.r)
			if diff := cmp.Diff(
				got,
				c.want,
				cmp.AllowUnexported(
					node.N{},
					mock.P{},
				),
			); diff != "" {
				t.Errorf("Search() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
