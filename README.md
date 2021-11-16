# go-kd

Golang k-d tree implementation with duplicate coordinate support

See https://en.wikipedia.org/wiki/K-d_tree for more information.

## Testing

```bash
go test -v github.com/downflux/go-kd/... -bench .
```

## Example

```golang
package main

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/kd"
	"github.com/downflux/go-kd/point"
)

// P implements the point.P interface, which needs to provide a coordinate
// vector function P().
type P struct {
	p   vector.V
	tag string
}

func (p P) P() vector.V { return p.p }

func main() {
	// N.B.: KD operations will return non-nil errors if the input vectors
	// are not a consistent length.
	v := *vector.New(1, 2, 3)
	origin := *vector.New(0, 0, 0)

	t, _ := kd.New([]point.P{
		P{p: v, tag: "A"},
		P{p: v, tag: "B"},
	})

	fmt.Println("KNN search")
	ns, _ := kd.KNN(t, origin, 2)
	for _, p := range ns {
		fmt.Println(p)
	}

	// Remove deletes the first data point at the given input coordinate and
	// matches the input check function.
	t.Remove(v, func(p point.P) bool {
		return p.(P).tag == "B"
	})

	// RadialFilter returns all points within the circle range and match the
	// input filter function.
	fmt.Println("radial search")
	ns, _ = kd.RadialFilter(
		t,
		*hypersphere.New(origin, 5),
		func(p point.P) bool { return true })
	for _, p := range ns {
		fmt.Println(p)
	}
}
```
