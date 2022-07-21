package node

import (
	"testing"

	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/point/mock"
)

type cmp point.D

func (f cmp) Less(a mock.P, b mock.P) bool {
	return a.P().X(point.D(f)) < b.P().X(point.D(f))
}

func equal(n *N, m *N) bool {
	if n == nil || m == nil {
		return n == m
	}
	return n.Low == m.Low && n.High == m.High && n.Pivot == m.Pivot && equal(n.Left, m.Left) && equal(n.Right, m.Right)
}

func TestNew(t *testing.T) {
	type config struct {
		name string
		opts O[mock.P]

		want *N
	}

	configs := []config{
		{
			name: "NullNode",
			opts: O[mock.P]{
				Data: nil,
				K:    2,
				N:    1,
				Dim:  0,
				Low:  0,
				High: 0,
			},
			want: &N{
				Low:   0,
				High:  0,
				Pivot: -1,
				Dim:   0,
			},
		},
		{
			name: "SingleElement",
			opts: O[mock.P]{
				Data: []mock.P{
					{
						X:    mock.U(1),
						Data: "foo",
					},
				},
				K:    1,
				N:    1,
				Dim:  0,
				Low:  0,
				High: 1,
			},
			want: &N{
				Low:   0,
				High:  1,
				Pivot: -1,
				Dim:   0,
			},
		},
		{
			name: "DoubleElement",
			opts: O[mock.P]{
				Data: []mock.P{
					{
						X:    mock.U(1),
						Data: "bar",
					},
					{
						X:    mock.U(-100),
						Data: "foo",
					},
				},
				K:    1,
				N:    1,
				Dim:  0,
				Low:  0,
				High: 2,
			},
			want: &N{
				Low:   1,
				High:  2,
				Pivot: 1,
				Dim:   0,
				Left: &N{
					Low:   0,
					High:  1,
					Pivot: 0,
					Dim:   1,
				},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := New[mock.P](c.opts)
			if !equal(got, c.want) {
				t.Errorf("New() = %v, want = %v", got, c.want)
			}

			open := []*N{got}
			for len(open) > 0 {
				var n *N
				n, open = open[0], open[1:]

				if n.Left != nil {
					open = append(open, n.Left)
				}
				if n.Right != nil {
					open = append(open, n.Right)
				}

				for i := n.Low; n.Pivot >= 0 && i < n.Pivot; i++ {
					if cmp(n.Dim).Less(c.opts.Data[n.Pivot], c.opts.Data[i]) {
						t.Errorf("Less(%v, %v) = false, want = true", c.opts.Data[n.Pivot], c.opts.Data[i])
					}
				}
				for i := n.Pivot; n.Pivot >= 0 && i < n.High; i++ {
					if cmp(n.Dim).Less(c.opts.Data[i], c.opts.Data[n.Pivot]) {
						t.Errorf("Less(%v, %v) = false, want = true", c.opts.Data[i], c.opts.Data[n.Pivot])
					}
				}
			}
		})
	}
}

func TestHoare(t *testing.T) {
	type result struct {
		data  []mock.P
		pivot int
	}

	type config struct {
		name string

		data  []mock.P
		pivot int
		low   int
		high  int
		less  func(a mock.P, b mock.P) bool

		want result
	}

	configs := []config{
		{
			name: "Trivial",
			data: []mock.P{
				mock.P{
					X:    mock.V(*vector.New(100, 80)),
					Data: "foo",
				},
			},
			pivot: 0,
			low:   0,
			high:  1,
			less:  cmp(point.AXIS_X).Less,
			want: result{
				data: []mock.P{
					mock.P{
						X:    mock.V(*vector.New(100, 80)),
						Data: "foo",
					},
				},
				pivot: 0,
			},
		},
		{
			name: "Simple/NoSwap",
			data: []mock.P{
				mock.P{
					X:    mock.U(0),
					Data: "foo",
				},
				mock.P{
					X:    mock.U(1),
					Data: "bar",
				},
			},
			pivot: 0,
			low:   0,
			high:  2,
			less:  cmp(point.AXIS_X).Less,
			want: result{
				data: []mock.P{
					mock.P{
						X:    mock.U(0),
						Data: "foo",
					},
					mock.P{
						X:    mock.U(1),
						Data: "bar",
					},
				},
				pivot: 0,
			},
		},
		{
			name: "Simple/Swap",
			data: []mock.P{
				mock.P{
					X:    mock.U(1),
					Data: "bar",
				},
				mock.P{
					X:    mock.U(0),
					Data: "foo",
				},
			},
			pivot: 0,
			low:   0,
			high:  2,
			less:  cmp(point.AXIS_X).Less,
			want: result{
				data: []mock.P{
					mock.P{
						X:    mock.U(0),
						Data: "foo",
					},
					mock.P{
						X:    mock.U(1),
						Data: "bar",
					},
				},
				pivot: 1,
			},
		},
		{
			name: "Pivot",
			data: []mock.P{
				mock.P{
					X:    mock.U(100),
					Data: "2",
				},
				mock.P{
					X:    mock.U(0),
					Data: "0",
				},
				mock.P{
					X:    mock.U(50),
					Data: "1",
				},
			},
			pivot: 1,
			low:   0,
			high:  3,
			less:  cmp(point.AXIS_X).Less,
			want: result{
				data: []mock.P{
					mock.P{
						X:    mock.U(0),
						Data: "0",
					},
					mock.P{
						X:    mock.U(100),
						Data: "2",
					},
					mock.P{
						X:    mock.U(50),
						Data: "1",
					},
				},
				pivot: 0,
			},
		},
		{
			name: "Pivot/Partial",
			data: []mock.P{
				mock.P{
					X:    mock.U(100),
					Data: "2",
				},
				mock.P{
					X:    mock.U(0),
					Data: "0",
				},
				mock.P{
					X:    mock.U(50),
					Data: "1",
				},
			},
			pivot: 2,
			low:   1,
			high:  3,
			less:  cmp(point.AXIS_X).Less,
			want: result{
				data: []mock.P{
					mock.P{
						X:    mock.U(100),
						Data: "2",
					},
					mock.P{
						X:    mock.U(0),
						Data: "0",
					},
					mock.P{
						X:    mock.U(50),
						Data: "1",
					},
				},
				pivot: 2,
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := hoare(c.data, c.pivot, c.low, c.high, c.less); got != c.want.pivot {
				t.Errorf("hoare() = %v, want = %v", got, c.want.pivot)
			}

			for i, got := range c.data {
				if !mock.Equal(got, c.want.data[i]) {
					t.Errorf("data[i] = %v, want = %v", got, c.want.data[i])
				}
			}
		})
	}
}
