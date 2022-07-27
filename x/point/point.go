package point

import (
	"github.com/downflux/go-geometry/nd/vector"
)

type P interface {
	P() vector.V
}
