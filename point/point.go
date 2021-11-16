// Package point defines the data interface for the K-D tree.
package point

import (
	"github.com/downflux/go-geometry/nd/vector"
)

type P interface {
	// P is the coordinate in the K-dimensional ambient space at which the
	// data is embedded.
	P() vector.V
}
