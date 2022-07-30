package filter

import (
	"github.com/downflux/go-kd/point"
)

type F[T point.P] func(p T) bool
