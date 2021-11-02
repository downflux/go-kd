package point

import (
	"github.com/downflux/go-geometry/vector"
)

// TODO(minkezhang): Refactor using generics.
type P interface {
	V() vector.V
}
