package pq

import (
	"container/heap"
	"testing"

	"github.com/downflux/go-kd/x/point/mock"
	"github.com/google/go-cmp/cmp"
)

var _ heap.Interface = &max[mock.P]{}

func TestPQ(t *testing.T) {
	type item struct {
		p        mock.P
		priority float64
	}

	type config struct {
		name string
		data []item
		size int
		want []mock.P
	}

	configs := []config{
		{
			name: "Null",
			data: nil,
			want: nil,
		},

		{
			name: "Trivial",
			data: []item{
				{
					p: mock.P{
						X: mock.U(1),
					},
					priority: 0,
				},
			},
			size: 1,
			want: []mock.P{
				mock.P{
					X: mock.U(1),
				},
			},
		},
		{
			name: "Trivial/NoSize",
			data: []item{
				{
					p: mock.P{
						X: mock.U(1),
					},
					priority: 0,
				},
			},
			size: 0,
			want: nil,
		},

		{
			name: "Sorted",
			data: []item{
				{
					p:        mock.P{X: mock.U(0)},
					priority: 5,
				},
				{
					p:        mock.P{X: mock.U(-1)},
					priority: 4,
				},
				{
					p:        mock.P{X: mock.U(-2)},
					priority: 3,
				},
				{
					p:        mock.P{X: mock.U(-3)},
					priority: 2,
				},
				{
					p:        mock.P{X: mock.U(-4)},
					priority: 1,
				},
			},
			size: 5,
			want: []mock.P{
				mock.P{X: mock.U(0)},
				mock.P{X: mock.U(-1)},
				mock.P{X: mock.U(-2)},
				mock.P{X: mock.U(-3)},
				mock.P{X: mock.U(-4)},
			},
		},
		{
			name: "Sorted/Reverse",
			data: []item{
				{
					p:        mock.P{X: mock.U(0)},
					priority: 1,
				},
				{
					p:        mock.P{X: mock.U(-1)},
					priority: 2,
				},
				{
					p:        mock.P{X: mock.U(-2)},
					priority: 3,
				},
				{
					p:        mock.P{X: mock.U(-3)},
					priority: 4,
				},
				{
					p:        mock.P{X: mock.U(-4)},
					priority: 5,
				},
			},
			size: 5,
			want: []mock.P{
				mock.P{X: mock.U(-4)},
				mock.P{X: mock.U(-3)},
				mock.P{X: mock.U(-2)},
				mock.P{X: mock.U(-1)},
				mock.P{X: mock.U(0)},
			},
		},
		{
			name: "Sorted/Shuffled",
			data: []item{
				{
					p:        mock.P{X: mock.U(0)},
					priority: 4,
				},
				{
					p:        mock.P{X: mock.U(-1)},
					priority: 2,
				},
				{
					p:        mock.P{X: mock.U(-2)},
					priority: 5,
				},
				{
					p:        mock.P{X: mock.U(-3)},
					priority: 1,
				},
				{
					p:        mock.P{X: mock.U(-4)},
					priority: 3,
				},
			},
			size: 5,
			want: []mock.P{
				mock.P{X: mock.U(-2)},
				mock.P{X: mock.U(0)},
				mock.P{X: mock.U(-4)},
				mock.P{X: mock.U(-1)},
				mock.P{X: mock.U(-3)},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			pq := New[mock.P](c.size)
			for _, d := range c.data {
				pq.Push(d.p, d.priority)
			}

			var got []mock.P
			for !pq.Empty() {
				got = append(got, pq.Pop())
			}

			if diff := cmp.Diff(
				c.want, got,
			); diff != "" {
				t.Errorf("Pop() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
