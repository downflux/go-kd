package point

import (
	"github.com/downflux/go-geometry/nd/vector"
)

// TODO(minkezhang): Refactor using generics.
type P interface {
	// P is the coordinate on the XY-plane at which the data is embedded.
	P() vector.V
}
