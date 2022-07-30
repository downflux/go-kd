// Package container exports the expected storage API used for querying a set of
// objects in a system. This may be used to more freely move between different
// implementations as the conditions of the system change, e.g. when the number
// or density of agents reach some threshold.
package container

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/filter"
	"github.com/downflux/go-kd/x/point"
)

type C[T point.P] interface {
	KNN(p vector.V, k int, f filter.F[T]) []T
	RangeSearch(q hyperrectangle.R, f filter.F[T]) []T
	Data() []T
	Balance()
}
