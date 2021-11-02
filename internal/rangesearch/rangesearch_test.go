package rangesearch

import (
	"testing"

	"github.com/downflux/go-geometry/rectangle"
	"github.com/downflux/go-geometry/vector"
	"github.com/downflux/go-kd/internal/node"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"

	mock "github.com/downflux/go-kd/internal/point/mock"
)

const (
	tolerance = 1e-10
)

func TestSearch(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		r    rectangle.R
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Trivial",
			n:    nil,
			r:    *rectangle.New(*vector.New(0, 0), *vector.New(1, 2)),
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
				tolerance,
			)

			return []config{
				{
					name: "Trivial/Embedded",
					n:    n,
					r: *rectangle.New(
						*vector.New(0, 0),
						*vector.New(2, 3),
					),
					want: []*node.N{n},
				},
				{
					name: "Trivial/Disjoint",
					n:    n,
					r: *rectangle.New(
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
					vector.V{},
					mock.P{},
				),
			); diff != "" {
				t.Errorf("Search() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
