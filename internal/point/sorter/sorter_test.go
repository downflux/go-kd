package sorter

import (
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock"
)

const (
	tolerance = 1e-10
)

var _ point.P = mock.P{}

func TestSorterLen(t *testing.T) {
	testConfigs := []struct {
		name string
		data []point.P
		axis vector.D
		want int
	}{
		{
			name: "Null",
			data: nil,
			axis: vector.AXIS_X,
			want: 0,
		},
		{
			name: "Simple",
			data: []point.P{
				*mock.New(*vector.New(1, 1), ""),
			},
			axis: vector.AXIS_X,
			want: 1,
		},
		{
			name: "Duplicate",
			data: []point.P{
				*mock.New(*vector.New(1, 1), ""),
				*mock.New(*vector.New(1, 1), ""),
			},
			axis: vector.AXIS_X,
			want: 2,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			s := s{
				data: c.data,
				axis: c.axis,
			}
			if got := s.Len(); got != c.want {
				t.Errorf("Len() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestSorterLess(t *testing.T) {
	testConfigs := []struct {
		name string
		s    *s
		i    int
		j    int
		want bool
	}{
		{
			name: "OneElement/X",
			s: &s{
				data: []point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				axis: vector.AXIS_X,
			},
			i:    0,
			j:    0,
			want: false,
		},
		{
			name: "OneElement/Y",
			s: &s{
				data: []point.P{
					*mock.New(*vector.New(1, 2), ""),
				},
				axis: vector.AXIS_Y,
			},
			i:    0,
			j:    0,
			want: false,
		},
		{
			name: "Simple/X",
			s: &s{
				data: []point.P{
					*mock.New(*vector.New(1, 2), ""),
					*mock.New(*vector.New(2, 1), ""),
				},
				axis: vector.AXIS_X,
			},
			i:    0,
			j:    1,
			want: true,
		},
		{
			name: "Simple/Y",
			s: &s{
				data: []point.P{
					*mock.New(*vector.New(1, 2), ""),
					*mock.New(*vector.New(2, 1), ""),
				},
				axis: vector.AXIS_Y,
			},
			i:    0,
			j:    1,
			want: false,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.s.Less(c.i, c.j); got != c.want {
				t.Errorf("Less = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestSorterSwap(t *testing.T) {
	testConfigs := []struct {
		name string
		data []point.P
		i    int
		j    int
		want []vector.V
	}{
		{
			name: "OneElement",
			data: []point.P{
				*mock.New(*vector.New(1, 2), ""),
			},
			i: 0,
			j: 0,
			want: []vector.V{
				*vector.New(1, 2),
				*vector.New(1, 2),
			},
		},
		{
			name: "TwoElements",
			data: []point.P{
				*mock.New(*vector.New(1, 2), ""),
				*mock.New(*vector.New(2, 1), ""),
			},
			i: 0,
			j: 1,
			want: []vector.V{
				*vector.New(2, 1),
				*vector.New(1, 2),
			},
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			s := &s{data: c.data}
			s.Swap(c.i, c.j)

			got := []vector.V{
				c.data[c.i].P(),
				c.data[c.j].P(),
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(mock.P{})); diff != "" {
				t.Errorf("Swap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSort(t *testing.T) {
	testConfigs := []struct {
		name string
		data []point.P
		axis vector.D
		want []point.P
	}{
		{
			name: "Trivial/NoData/X",
			data: []point.P{*mock.New(*vector.New(1, 2), "")},
			axis: vector.AXIS_X,
			want: []point.P{*mock.New(*vector.New(1, 2), "")},
		},
		{
			name: "Trivial/NoData/Y",
			data: []point.P{*mock.New(*vector.New(1, 2), "")},
			axis: vector.AXIS_Y,
			want: []point.P{*mock.New(*vector.New(1, 2), "")},
		},

		{
			name: "Trivial/WithData/X",
			data: []point.P{*mock.New(*vector.New(1, 2), "foo")},
			axis: vector.AXIS_X,
			want: []point.P{*mock.New(*vector.New(1, 2), "foo")},
		},
		{
			name: "Trivial/NoData/Y",
			data: []point.P{*mock.New(*vector.New(1, 2), "foo")},
			axis: vector.AXIS_Y,
			want: []point.P{*mock.New(*vector.New(1, 2), "foo")},
		},

		{
			name: "Simple/NoData/X",
			data: []point.P{
				*mock.New(*vector.New(3, 1), ""),
				*mock.New(*vector.New(2, 2), ""),
				*mock.New(*vector.New(1, 3), ""),
			},
			axis: vector.AXIS_X,
			want: []point.P{
				*mock.New(*vector.New(1, 3), ""),
				*mock.New(*vector.New(2, 2), ""),
				*mock.New(*vector.New(3, 1), ""),
			},
		},
		{
			name: "Simple/NoData/Y",
			data: []point.P{
				*mock.New(*vector.New(1, 3), ""),
				*mock.New(*vector.New(2, 2), ""),
				*mock.New(*vector.New(3, 1), ""),
			},
			axis: vector.AXIS_Y,
			want: []point.P{
				*mock.New(*vector.New(3, 1), ""),
				*mock.New(*vector.New(2, 2), ""),
				*mock.New(*vector.New(1, 3), ""),
			},
		},

		{
			name: "Simple/WithData/X",
			data: []point.P{
				*mock.New(*vector.New(3, 1), "foo3"),
				*mock.New(*vector.New(2, 2), "foo2"),
				*mock.New(*vector.New(1, 3), "foo1"),
			},
			axis: vector.AXIS_X,
			want: []point.P{
				*mock.New(*vector.New(1, 3), "foo1"),
				*mock.New(*vector.New(2, 2), "foo2"),
				*mock.New(*vector.New(3, 1), "foo3"),
			},
		},
		{
			name: "Simple/WithData/Y",
			data: []point.P{
				*mock.New(*vector.New(1, 3), "foo1"),
				*mock.New(*vector.New(2, 2), "foo2"),
				*mock.New(*vector.New(3, 1), "foo3"),
			},
			axis: vector.AXIS_Y,
			want: []point.P{
				*mock.New(*vector.New(3, 1), "foo3"),
				*mock.New(*vector.New(2, 2), "foo2"),
				*mock.New(*vector.New(1, 3), "foo1"),
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			Sort(c.data, c.axis)

			if diff := cmp.Diff(
				c.want,
				c.data,
				cmp.AllowUnexported(mock.P{})); diff != "" {
				t.Errorf("Swap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
