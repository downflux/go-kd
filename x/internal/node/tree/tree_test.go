package tree

import (
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/node"
	"github.com/downflux/go-kd/x/internal/node/util"
	"github.com/downflux/go-kd/x/point"
	"github.com/downflux/go-kd/x/point/mock"
)

var _ node.N[mock.P] = &N[mock.P]{}

func equal[T point.P](n *N[T], m *N[T]) bool {
	if n == nil || m == nil {
		return n == m
	}
	return n.low == m.low && n.high == m.high && n.pivot == m.pivot && equal(n.left, m.left) && equal(n.right, m.right)
}

func TestNew(t *testing.T) {
	type config struct {
		name string
		opts O[mock.P]

		want *N[mock.P]
	}

	configs := []config{
		{
			name: "NullNode",
			opts: O[mock.P]{
				Data: nil,
				K:    2,
				N:    1,
				Axis: 0,
				Low:  0,
				High: 0,
			},
			want: nil,
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
				Axis: 0,
				Low:  0,
				High: 1,
			},
			want: &N[mock.P]{
				low:   0,
				high:  1,
				pivot: -1,
				axis:  0,
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
				Axis: 0,
				Low:  0,
				High: 2,
			},
			want: &N[mock.P]{
				low:   0,
				high:  2,
				pivot: 1,
				axis:  0,
				left: &N[mock.P]{
					low:   0,
					high:  1,
					pivot: -1,
					axis:  0,
				},
			},
		},
		{
			name: "TripleElement/Unbalanced/BigLeaf",
			opts: O[mock.P]{
				Data: []mock.P{
					{
						X:    mock.U(-100),
						Data: "foo",
					},
					{
						X:    mock.U(1),
						Data: "bar",
					},
					{
						X:    mock.U(0),
						Data: "baz",
					},
				},
				K:    1,
				N:    2,
				Axis: 0,
				Low:  0,
				High: 3,
			},
			want: &N[mock.P]{
				low:   0,
				high:  3,
				pivot: 0,
				axis:  0,
				right: &N[mock.P]{
					low:   1,
					high:  3,
					pivot: -1,
					axis:  0,
				},
			},
		},
		{
			name: "TripleElement/Unbalanced/BigLeaf/BigK",
			opts: O[mock.P]{
				Data: []mock.P{
					{
						X:    mock.V(*vector.New(-100, 1)),
						Data: "foo",
					},
					{
						X:    mock.V(*vector.New(1, 50)),
						Data: "bar",
					},
					{
						X:    mock.V(*vector.New(0, 75)),
						Data: "baz",
					},
				},
				K:    2,
				N:    2,
				Axis: 0,
				Low:  0,
				High: 3,
			},
			want: &N[mock.P]{
				low:   0,
				high:  3,
				pivot: 0,
				axis:  0,
				right: &N[mock.P]{
					low:   1,
					high:  3,
					pivot: -1,
					axis:  1,
				},
			},
		},
		{
			name: "TripleElement/Unbalanced",
			opts: O[mock.P]{
				Data: []mock.P{
					{
						X:    mock.U(-100),
						Data: "foo",
					},
					{
						X:    mock.U(1),
						Data: "bar",
					},
					{
						X:    mock.U(0),
						Data: "baz",
					},
				},
				K:    1,
				N:    1,
				Axis: 0,
				Low:  0,
				High: 3,
			},
			want: &N[mock.P]{
				low:   0,
				high:  3,
				pivot: 0,
				axis:  0,
				right: &N[mock.P]{
					low:   1,
					high:  3,
					pivot: 2,
					axis:  1,
					left: &N[mock.P]{
						low:   1,
						high:  2,
						pivot: -1,
						axis:  0,
					},
				},
			},
		},
		{
			name: "TripleElement/Unbalanced/BigK",
			opts: O[mock.P]{
				Data: []mock.P{
					{
						X:    mock.V(*vector.New(-100, 1)),
						Data: "foo",
					},
					{
						X:    mock.V(*vector.New(1, 50)),
						Data: "bar",
					},
					{
						X:    mock.V(*vector.New(0, 75)),
						Data: "baz",
					},
				},
				K:    2,
				N:    1,
				Axis: 0,
				Low:  0,
				High: 3,
			},
			want: &N[mock.P]{
				low:   0,
				high:  3,
				pivot: 0,
				axis:  0,
				right: &N[mock.P]{
					low:   1,
					high:  3,
					pivot: 1,
					axis:  1,
					right: &N[mock.P]{
						low:   2,
						high:  3,
						pivot: -1,
						axis:  0,
					},
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

			if got == nil {
				return
			}

			if util.Validate[mock.P](got) != true {
				t.Errorf("Validate() = %v, want = %v", false, true)
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
		less  func(a point.V, b point.V) bool

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
			less:  util.Cmp(point.AXIS_X).Less,
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
			less:  util.Cmp(point.AXIS_X).Less,
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
			less:  util.Cmp(point.AXIS_X).Less,
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
			less:  util.Cmp(point.AXIS_X).Less,
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
			less:  util.Cmp(point.AXIS_X).Less,
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
