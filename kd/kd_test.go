package kd

import (
	"fmt"
	"math"
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/testdata/generator"
	"github.com/downflux/go-kd/point"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock"
)

const (
	// f defines the rough percentage of points benchmark tests should seek
	// for. We deem this an arbitrary but reasonable enough heuristic for
	// normal data access patterns.
	f = 0.05
)

var (
	// kRange is used for several benchmark tests to specify the dimension
	// of the ambient space. Vectors will have the specified number of
	// elements.
	//
	// N.B.: The size of a tree is dominated by the total amount of data
	// stored at each point, and for large K, is dominated by the size of
	// the vector. A float64 is 8B; a K=100 vector is therefore 800B, and
	// benchmarking with ~1M elements is around 800MB. Keep this lower limit
	// in mind when trying to test for more stressful conditions.
	kRange = []int{2, 10}

	// nRange is used for several benchmark tests to specify the number of
	// elements that should be added to the K-D tree.
	nRange = []int{1e4, 1e5, 2e5, 3e5, 1e6}
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

// BenchmarkNew measures the construction time of the tree.
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
		ps := generator.P(c.n, vector.D(c.k))
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New(ps)
			}
		})
	}
}

// BenchmarkKNN measures the average KNN search time.
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
			knn := int(float64(n) * f)
			kd, _ := New(generator.P(n, vector.D(k)))
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
			for i := 0; i < b.N; i++ {
				if _, err := KNN(c.kd, generator.V(vector.D(c.k)), c.knn); err != nil {
					b.Errorf("KNN() = _, %v, want = _, %v", err, nil)
				}
			}
		})
	}
}

// BenchmarkFilter measures the average radial search time.
func BenchmarkFilter(b *testing.B) {
	type config struct {
		name string

		// k is the number of dimensions of the input data, i.e. the "K"
		// in K-D tree.
		k int

		// n is the number of points to generate.
		n int

		// l is the length of the hyperrectangle.
		l float64

		kd *T
	}

	var testConfigs []config

	for _, k := range kRange {
		for _, n := range nRange {
			// We want to search approximately f = 0.05 of the
			// total space, so we will define a ball with this
			// constraint in mind.
			volume := math.Pow(generator.Interval, float64(k)) * f
			l := math.Pow(volume, 1/float64(k))

			kd, _ := New(generator.P(n, vector.D(k)))
			testConfigs = append(testConfigs, config{
				name: fmt.Sprintf("K=%v/N=%v", k, n),
				k:    k,
				kd:   kd,
				l:    l,
			})
		}
	}

	for _, c := range testConfigs {
		p := generator.V(vector.D(c.k))

		min := make([]float64, c.k)
		max := make([]float64, c.k)

		for i := vector.D(0); i < vector.D(c.k); i++ {
			min[i] = p.X(i) - (c.l / 2)
			max[i] = p.X(i) + (c.l / 2)
		}

		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := Filter(
					c.kd,
					*hyperrectangle.New(min, max),
					func(point.P) bool { return true },
				); err != nil {
					b.Errorf("Filter() = _, %v, want = _, %v", err, nil)
				}
			}
		})
	}
}
