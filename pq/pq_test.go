package pq

import (
	"testing"

	"github.com/downflux/go-geometry/vector"
	"github.com/downflux/go-kd/node"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"

	mock "github.com/downflux/go-kd/point/mock"
)

const (
	tolerance = 1e-10
)

type item struct {
	n *node.N
	p float64
}

func TestHeap(t *testing.T) {
	type config struct {
		name string
		data []item
		size int
		want []*node.N
	}

	testConfigs := []config{
		{
			name: "Null",
			data: nil,
			want: nil,
		},
	}

	testConfigs = append(
		testConfigs,
		func() []config {
			n := node.New(
				[]point.P{
					*mock.New(*vector.New(1, 2), "A"),
				},
				0,
				tolerance,
			)
			return []config{
				{
					name: "Trivial",
					data: []item{
						{n: n, p: 0},
					},
					size: 1,
					want: []*node.N{n},
				},
				{
					name: "Trivial/NoSize",
					data: []item{
						{n: n, p: 0},
					},
					size: 0,
					want: nil,
				},
			}
		}()...,
	)

	testConfigs = append(
		testConfigs,

		func() []config {
			ns := []*node.N{
				node.New(
					[]point.P{
						*mock.New(*vector.New(1, 5), "A"),
					},
					0,
					tolerance,
				),
				node.New(
					[]point.P{
						*mock.New(*vector.New(2, 4), "B"),
					},
					0,
					tolerance,
				),
				node.New(
					[]point.P{
						*mock.New(*vector.New(3, 3), "C"),
					},
					0,
					tolerance,
				),
				node.New(
					[]point.P{
						*mock.New(*vector.New(4, 2), "D"),
					},
					0,
					tolerance,
				),
				node.New(
					[]point.P{
						*mock.New(*vector.New(5, 1), "E"),
					},
					0,
					tolerance,
				),
			}
			return []config{
				{
					name: "Order/Reverse",
					data: []item{
						{n: ns[0], p: 0},
						{n: ns[1], p: 1},
						{n: ns[2], p: 2},
						{n: ns[3], p: 3},
						{n: ns[4], p: 4},
					},
					size: 5,
					want: []*node.N{
						ns[4],
						ns[3],
						ns[2],
						ns[1],
						ns[0],
					},
				},
				{
					name: "Order/InOrder",
					data: []item{
						{n: ns[0], p: 4},
						{n: ns[1], p: 3},
						{n: ns[2], p: 2},
						{n: ns[3], p: 1},
						{n: ns[4], p: 0},
					},
					size: 5,
					want: []*node.N{
						ns[0],
						ns[1],
						ns[2],
						ns[3],
						ns[4],
					},
				},
				{
					name: "Order/Shuffled",
					data: []item{
						{n: ns[0], p: 3},
						{n: ns[1], p: 1},
						{n: ns[2], p: 4},
						{n: ns[3], p: 0},
						{n: ns[4], p: 2},
					},
					size: 5,
					want: []*node.N{
						ns[2],
						ns[0],
						ns[4],
						ns[1],
						ns[3],
					},
				},
			}
		}()...,
	)

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			q := New(c.size)
			for _, d := range c.data {
				q.Push(d.n, d.p)
			}

			var got []*node.N
			for !q.Empty() {
				got = append(got, q.Pop())
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(node.N{}, vector.V{}, mock.P{}),
			); diff != "" {
				t.Errorf("Pop() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
