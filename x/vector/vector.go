package vector

import (
	"github.com/downflux/go-geometry/nd/vector"
)

type Comparator vector.D

func (c Comparator) Less(v vector.V, u vector.V) bool { return v.X(vector.D(c)) < u.X(vector.D(c)) }
