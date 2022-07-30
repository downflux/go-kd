package bruteforce

import (
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/container"
	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"
)

var _ container.C[mock.P] = &L[mock.P]{}

func TestDelete(t *testing.T) {
	type config struct {
		name string
		data []*mock.P
		vs   []vector.V

		want []*mock.P
	}

	configs := []config{
		{
			name: "Nil",
			data: nil,
			vs: []vector.V{
				mock.U(100),
			},
			want: []*mock.P{},
		},
		{
			name: "Simple",
			data: []*mock.P{
				&mock.P{X: mock.U(50)},
				&mock.P{X: mock.U(100)},
			},
			vs: []vector.V{
				mock.U(100),
			},
			want: []*mock.P{
				&mock.P{X: mock.U(50)},
			},
		},
		{
			name: "Degenerate",
			data: []*mock.P{
				&mock.P{X: mock.U(100), Data: "A"},
				&mock.P{X: mock.U(100), Data: "B"},
			},
			vs: []vector.V{
				mock.U(100),
			},
			want: []*mock.P{
				&mock.P{X: mock.U(100), Data: "B"},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			l := New(c.data)
			for _, v := range c.vs {
				l.Remove(v, func(p *mock.P) bool { return vector.Within(v, p.P()) })
			}

			got := l.Data()
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("Data() mismatch (-want +got):\n%v", diff)
			}

		})
	}
}

func TestInsert(t *testing.T) {
	type config struct {
		name string
		data []*mock.P
		ps   []*mock.P

		want []*mock.P
	}

	configs := []config{
		{
			name: "Trivial",
			data: nil,
			ps: []*mock.P{
				&mock.P{X: mock.U(100)},
			},
			want: []*mock.P{
				&mock.P{X: mock.U(100)},
			},
		},
		{
			name: "MultipleInsert",
			data: nil,
			ps: []*mock.P{
				&mock.P{X: mock.U(101)},
				&mock.P{X: mock.U(100)},
				&mock.P{X: mock.U(202)},
			},
			want: []*mock.P{
				&mock.P{X: mock.U(101)},
				&mock.P{X: mock.U(100)},
				&mock.P{X: mock.U(202)},
			},
		},
		{
			name: "MultipleInsert/NonNil",
			data: []*mock.P{
				&mock.P{X: mock.U(4)},
				&mock.P{X: mock.U(5)},
			},
			ps: []*mock.P{
				&mock.P{X: mock.U(101)},
				&mock.P{X: mock.U(100)},
				&mock.P{X: mock.U(202)},
			},
			want: []*mock.P{
				&mock.P{X: mock.U(4)},
				&mock.P{X: mock.U(5)},
				&mock.P{X: mock.U(101)},
				&mock.P{X: mock.U(100)},
				&mock.P{X: mock.U(202)},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			l := New(c.data)
			for _, p := range c.ps {
				l.Insert(p)
			}

			got := l.Data()
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("Data() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
