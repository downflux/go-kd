// Package sorter provides a way to sort a list of points pivoting on a
// specified axis.
package sorter

import (
	"sort"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"
)

type s struct {
	axis vector.D
	data []point.P
}

func (s *s) Len() int           { return len(s.data) }
func (s *s) Less(i, j int) bool { return s.data[i].P().X(s.axis) < s.data[j].P().X(s.axis) }
func (s *s) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }

// Sort sorts a list of points in-place by the given axis.
func Sort(ps []point.P, a vector.D) { sort.Sort(&s{axis: a, data: ps}) }
