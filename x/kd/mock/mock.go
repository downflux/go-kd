package mock

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/point"
)

// I is a mock interface used for faciliating benchmark tests. This
// should not be consumed by external users.
type I[U point.P] interface {
	KNN(p vector.V, k int) []U
	RangeSearch(q hyperrectangle.R) []U
	Data() []U
}
