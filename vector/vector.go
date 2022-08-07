package vector

import (
	"github.com/downflux/go-geometry/nd/vector"
)

type Comparator vector.D

func (c Comparator) Less(v vector.V, u vector.V) bool { return v.X(vector.D(c)) < u.X(vector.D(c)) }

// Less returns the lexicographical ordering between two vectors.
func Less(v vector.V, u vector.V) bool {
	if v.Dimension() != u.Dimension() {
		panic("mismatching vector dimensions")
	}
	for i := vector.D(0); i < v.Dimension(); i++ {
		if !Comparator(i).Less(v, u) {
			return false
		}
	}
	return true
}
