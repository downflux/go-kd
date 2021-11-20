package sorter

import (
	"sort"
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/testdata/generator"
	"github.com/downflux/go-kd/point"
	"github.com/google/go-cmp/cmp"
	"github.com/kyroy/kdtree"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock/simple"
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

type reference struct {
	dimension int
	points    []kdtree.Point
}

func (b *reference) Len() int { return len(b.points) }
func (b *reference) Less(i, j int) bool {
	return b.points[i].Dimension(b.dimension) < b.points[j].Dimension(b.dimension)
}
func (b *reference) Swap(i, j int) { b.points[i], b.points[j] = b.points[j], b.points[i] }

// BenchmarkSorter provides a sanity check that our implementation is not too
// far worse than the reference.
//
// The most travelled path in K-D tree construction is sorting the slice of
// coordinates by a particular axis. We want to be sure that our sorting
// operation on our specific data structure is at least baseline comparable to
// that of a more widely used alternative.
func BenchmarkSorter(b *testing.B) {
	const n = 1e6

	// N.B.: generator.R creates a list of kyroy/tree/points.Point
	// instances. Our generator creates a list of N-dimensional points, vs.
	// the specific Points2D struct used in the reference example. This is
	// because our implementation only supports N-dimensional vectors, and
	// the reference handles the 2D vector case separately.
	ps := generator.P(n, 2)
	rs := generator.R(ps)

	b.Run("Reference", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sort.Sort(&reference{dimension: 0, points: rs})
		}
	})
	b.Run("Sorter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Sort(ps, 0)
		}
	})
}
