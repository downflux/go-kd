package kd

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock"
)

const (
	f = .15
)

var (
	kRange = []int{2, 10, 100}
	nRange = []int{1e5, 2e5, 3e5}
)

func TestConsistentK(t *testing.T) {
	if _, err := New([]point.P{
		*mock.New(*vector.New(1), "A"),
		*mock.New(*vector.New(1, 2), "A"),
	}); err == nil {
		t.Errorf("New() = _, %v, want a non-nil error", err)
	}

	kd, _ := New([]point.P{
		*mock.New(*vector.New(1, 2, 3), "A"),
		*mock.New(*vector.New(1, 2, 3), "B"),
	})

	if _, err := KNN(kd, *vector.New(1), 2); err == nil {
		t.Errorf("KNN() = _, %v, want a non-nil error", err)
	}

	if err := kd.Insert(*mock.New(*vector.New(1), "A")); err == nil {
		t.Errorf("Insert() = _, %v, want a non-nil error", err)
	}

	if _, err := kd.Remove(*vector.New(1), func(point.P) bool { return true }); err == nil {
		t.Errorf("Remove() = _, %v, want a non-nil error", err)
	}

	if _, err := RadialFilter(kd, *hypersphere.New(*vector.New(1), 5), func(p point.P) bool { return true }); err == nil {
		t.Errorf("RadialFilter() = _, %v, want a non-nil error", err)
	}
}

func rn() float64 { return rand.Float64()*200 - 100 }

func rv(d int) vector.V {
	xs := make([]float64, d)
	for i := 0; i < d; i++ {
		xs[i] = rn()
	}

	return *vector.New(xs...)
}

func rp(n int, d int) []point.P {
	var ps []point.P
	for i := 0; i < n; i++ {
		ps = append(ps, *mock.New(rv(d), fmt.Sprintf("Random-%v", i)))
	}
	return ps
}

func BenchmarkNew(b *testing.B) {
	type config struct {
		name string

		// k is the number of dimensions of the input data, i.e. the "K"
		// in K-D tree.
		k int

		// n is the number of points to generate.
		n int
	}

	var testConfigs []config

	for _, k := range kRange {
		for _, n := range nRange {
			testConfigs = append(testConfigs, config{
				name: fmt.Sprintf("K=%v/N=%v", k, n),
				k:    k,
				n:    n,
			})
		}
	}

	for _, c := range testConfigs {
		b.Run(c.name, func(b *testing.B) {
			ps := rp(c.n, c.k)
			b.ResetTimer()
			New(ps)
		})
	}
}

func BenchmarkKNN(b *testing.B) {
	type config struct {
		name string

		// k is the number of dimensions of the input data, i.e. the "K"
		// in K-D tree.
		k int

		// knn is the number of points to look for in the KNN search.
		knn int

		kd *T
	}

	var testConfigs []config

	for _, k := range kRange {
		for _, n := range nRange {
			fmt.Printf("BenchmarkKNN/K=%v/N=%v\n", k, n)
			knn := int(float64(n) * f)
			kd, _ := New(rp(n, k))
			testConfigs = append(testConfigs, config{
				name: fmt.Sprintf("K=%v/N=%v", k, n),
				kd:   kd,
				k:    k,
				knn:  knn,
			})
		}
	}

	for _, c := range testConfigs {
		b.Run(c.name, func(b *testing.B) {
			if _, err := KNN(c.kd, rv(c.k), c.knn); err != nil {
				b.Errorf("KNN() = _, %v, want = _, %v", err, nil)
			}
		})
	}
}

func BenchmarkRadialFilter(b *testing.B) {
	type config struct {
		name string

		// k is the number of dimensions of the input data, i.e. the "K"
		// in K-D tree.
		k int

		// n is the number of points to generate.
		n int

		// r is the ball radius in the RadialFilter search.
		r float64

		kd *T
	}

	var testConfigs []config

	for _, k := range kRange {
		for _, n := range nRange {
			// We generate points in the interval
			//
			//   [-100, 100]
			//
			// along each axis in K-dimensional ambient
			// space.
			r := 100.0 * f

			kd, _ := New(rp(n, k))
			testConfigs = append(testConfigs, config{
				name: fmt.Sprintf("K=%v/N=%v", k, n),
				k:    k,
				kd:   kd,
				r:    r,
			})
		}
	}

	for _, c := range testConfigs {
		b.Run(c.name, func(b *testing.B) {
			if _, err := RadialFilter(
				c.kd,
				*hypersphere.New(rv(c.k), c.r),
				func(point.P) bool { return true },
			); err != nil {
				b.Errorf("RadialSearch() = _, %v, want = _, %v", err, nil)
			}
		})
	}
}
