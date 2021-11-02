# go-kd

Golang k-d tree implementation with duplicate coordinate support

See https://en.wikipedia.org/wiki/K-d_tree for more information.

## Example

```golang
package main

import (
	"fmt"

	"github.com/downflux/go-geometry/circle"
	"github.com/downflux/go-geometry/vector"

	"github.com/downflux/go-kd/kd"
	"github.com/downflux/go-kd/point"
)

// P implements the point.P interface, which needs to provide a coordinate
// vector function V().
type P struct {
	v   vector.V
	tag string
}

func (p P) V() vector.V { return p.v }

func main() {
	t := kd.New([]point.P{
		P{v: *vector.New(1, 2), tag: "A"},
		P{v: *vector.New(1, 2), tag: "B"},
	}, 1e-10)

	for _, p := range kd.KNN(t, *vector.New(0, 0), 2) {
		fmt.Println(p)
	}

	// Remove deletes the first data point at the given input coordinate and
	// matches the input check function.
	t.Remove(*vector.New(1, 2), func(p point.P) bool {
		return p.(P).tag == "B"
	})

	// RadialFilter returns all points within the circle range and match the
	// input filter function.
	for _, p := range kd.RadialFilter(
		t,
		*circle.New(*vector.New(0, 0), 5),
		func(p point.P) bool { return true }) {
		fmt.Println(p)
	}
}
```
