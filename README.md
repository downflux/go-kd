# go-kd

Golang K-D tree implementation with duplicate coordinate support

See [Wikipedia](https://en.wikipedia.org/wiki/K-d_tree) for more information.

## Testing

```bash
go test github.com/downflux/go-kd/...
go test github.com/downflux/go-kd/internal/perf \
  -bench . \
  -benchmem \
  -timeout=60m \
  -args -performance_test_size=large
```

## Example

```golang
package main

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"

	"github.com/downflux/go-kd/kd"
)

// P implements the point.P interface, which needs to provide a coordinate
// vector function P().
var _ point.P = &P{}

type P struct {
	p   vector.V
	tag string
}

func (p *P) P() vector.V     { return p.p }
func (p *P) Equal(q *P) bool { return vector.Within(p.P(), q.P()) && p.tag == q.tag }

func main() {
	data := []*P{
		&P{p: vector.V{1, 2}, tag: "A"},
		&P{p: vector.V{2, 100}, tag: "B"},
	}

	// Data is copy-constructed and may be read from outside the k-D tree.
	t := kd.New[*P](kd.O[*P]{
		Data: data,
		K:    2,
		N:    1,
	})

	fmt.Println("KNN search")
	for _, p := range kd.KNN(
		t,
		/* v = */ vector.V{0, 0},
		/* k = */ 2,
		func(p *P) bool { return true }) {
		fmt.Println(p)
	}

	// Remove deletes the first data point at the given input coordinate and
	// matches the input check function.
	p, ok := t.Remove(data[0].P(), data[0].Equal)
	fmt.Printf("removed %v (found = %v)\n", p, ok)

	// RangeSearch returns all points within the k-D bounds and matches the
	// input filter function.
	fmt.Println("range search")
	for _, p := range kd.RangeSearch(
		t,
		*hyperrectangle.New(
			/* min = */ vector.V{0, 0},
			/* max = */ vector.V{100, 100},
		),
		func(p *P) bool { return true },
	) {
		fmt.Println(p)
	}
}
```

## Performance (@v1.0.0)

This k-D tree implementation was compared against a brute force method, as well
as with the leading Golang k-D tree implementation
(http://github.com/kyroy/kdtree). Overall, we have found that

* tree construction is about 10x faster for large N.

  ```
  BenchmarkNew/kyroy/K=16/N=1000-8               758980 ns/op  146777 B/op
  BenchmarkNew/Real/K=16/N=1000/LeafSize=16-8    200749 ns/op   32637 B/op

  BenchmarkNew/kyroy/K=16/N=1000000-8                7407144200 ns/op  184813784 B/op
  BenchmarkNew/Real/K=16/N=1000000/LeafSize=256-8     588456300 ns/op   12462912 B/op
  ```

* KNN is significantly faster; for small N, we have found our implementation is
  ~10x faster than the reference implementation and ~20x faster than brute
  force. For large N, we have found up to ~15x faster than brute force, and a
  staggering _~1500x_ speedup when compared to the reference implementation.

  ```
  BenchmarkKNN/BruteForce/K=16/N=1000-8                   1563019 ns/op  2220712 B/op
  BenchmarkKNN/kyroy/K=16/N=1000/KNN=0.05-8                791415 ns/op    21960 B/op
  BenchmarkKNN/Real/K=16/N=1000/LeafSize=16/KNN=0.05-8      69537 ns/op    12024 B/op

  BenchmarkKNN/BruteForce/K=16/N=1000000-8                       5030811400 ns/op  5347687464 B/op
  BenchmarkKNN/kyroy/K=16/N=1000000/KNN=0.05-8                 529703585200 ns/op    23755688 B/op
  BenchmarkKNN/Real/K=16/N=1000000/LeafSize=256/KNN=0.05-8        335845533 ns/op     6044016 B/op
  ```

* RangeSearch is slower for small N -- we are approximately at parity for brute
  force, and ~10x slower than the reference implementation. However, at large N,
  we are ~300x faster than brute force, and ~100x faster than the reference
  implementation.

  ```
  BenchmarkRangeSearch/BruteForce/K=16/N=1000-8                        154712 ns/op   25208 B/op
  BenchmarkRangeSearch/kyroy/K=16/N=1000/Coverage=0.05-8                13373 ns/op     496 B/op
  BenchmarkRangeSearch/Real/K=16/N=1000/LeafSize=16/Coverage=0.05-8    193276 ns/op  101603 B/op

  BenchmarkRangeSearch/BruteForce/K=16/N=1000000-8                         173427000 ns/op  41678072 B/op
  BenchmarkRangeSearch/kyroy/K=16/N=1000000/Coverage=0.05-8                 56820240 ns/op       496 B/op
  BenchmarkRangeSearch/Real/K=16/N=1000000/LeafSize=256/Coverage=0.05-8       530937 ns/op    212134 B/op
  ```

Raw data on these results may be found [here](/internal/perf/results/v0.5.5.txt).
