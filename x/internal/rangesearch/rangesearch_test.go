package rangesearch

import (
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/node/tree"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"
)

func TestRangeSearch(t *testing.T) {
	type config struct {
		name string
		data []*mock.P
		k    vector.D
		n    int
		q    hyperrectangle.R
		want []*mock.P
	}

	configs := []config{
		{
			name: "Nil",
			data: nil,
			k:    1,
			n:    1,
			q:    *hyperrectangle.New(mock.U(1), mock.U(2)),
			want: nil,
		},
		{
			name: "Simple",
			data: []*mock.P{
				&mock.P{X: mock.U(1.5)},
			},
			k: 1,
			n: 1,
			q: *hyperrectangle.New(mock.U(1), mock.U(2)),
			want: []*mock.P{
				&mock.P{X: mock.U(1.5)},
			},
		},
		{
			name: "LR",
			data: []*mock.P{
				&mock.P{X: mock.U(1.5)},
				&mock.P{X: mock.U(1)},
				&mock.P{X: mock.U(2)},
			},
			k: 1,
			n: 1,
			q: *hyperrectangle.New(mock.U(1), mock.U(2)),
			want: []*mock.P{
				&mock.P{X: mock.U(1.5)},
				&mock.P{X: mock.U(1)},
				&mock.P{X: mock.U(2)},
			},
		},
		{
			name: "Partial",
			data: []*mock.P{
				&mock.P{X: mock.U(1.5)},
				&mock.P{X: mock.U(1)},
				&mock.P{X: mock.U(2)},
			},
			k: 1,
			n: 1,
			q: *hyperrectangle.New(mock.U(1.9), mock.U(2.1)),
			want: []*mock.P{
				&mock.P{X: mock.U(2)},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := RangeSearch[*mock.P](
				tree.New[*mock.P](tree.O[*mock.P]{
					Data: c.data,
					K:    c.k,
					N:    c.n,
					Axis: vector.AXIS_X,
				}),
				c.q,
				func(*mock.P) bool { return true },
			)
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("RangeSearch() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
