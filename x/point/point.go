package point

import (
	"github.com/downflux/go-kd/x/vector"
)

type P interface {
	P() vector.V
}
