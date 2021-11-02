// Package sorter provides a way to sort a list of points pivoting on a
// specified axis.
package sorter

import (
	"sort"

	"github.com/downflux/go-kd/axis"
	"github.com/downflux/go-kd/point"
)

type s struct {
	axis axis.Type
	data []point.P
}

func (s *s) Len() int { return len(s.data) }
func (s *s) Less(i, j int) bool {
	return axis.X(s.data[i].V(), s.axis) < axis.X(s.data[j].V(), s.axis)
}
func (s *s) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }

// Sort sorts a list of points in-place by the given axis.
func Sort(ps []point.P, a axis.Type) { sort.Sort(&s{axis: a, data: ps}) }
