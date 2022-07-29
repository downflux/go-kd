package container

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point"
)

type I[T point.P] interface {
	KNN(p vector.V, k int) []T
	RangeSearch(q hyperrectangle.R) []T
	Data() []T
}
