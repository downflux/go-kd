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
	// KNN returns the k-nearest neighbors of the given search coordinates.
	//
	// N.B.: KNN will return at max k neighbors; in the degenerate case that
	// multiple data points reside at the same spacial coordinate, this
	// function will arbitrarily return a subset of these to fulfill the
	// k-neighbors constraint.
	KNN(p vector.V, k int, f filter.F[T]) []T

	// Data returns all data stored in the K-D tree.
	Data() []T

	// RangeSearch returns a set of data points in the given bounding box.
	// Data points are added to the returned set if they fall inside the
	// bounding box and passes the given filter function.
	RangeSearch(q hyperrectangle.R, f filter.F[T]) []T

	// Balance ensures the tree has minimal height by reconstructing the
	// tree. Note that in general, mutations to the structure of the tree
	// are expensive and complicated, so much so that it's easier to just
	// redo the tree from scratch than to worry about shifting nodes around.
	Balance()
}
