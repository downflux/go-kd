package kd

import (
	"fmt"
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/internal/testdata/generator"

	multifield "github.com/downflux/go-kd/internal/point/testdata/mock/multifield"
	simple "github.com/downflux/go-kd/internal/point/testdata/mock/simple"
)

var (
	kRange = []int{2, 10}
	nRange = []int{1e4, 1e5, 2e5, 3e5, 1e6}
)

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
		ps := make([]simple.P, 0, c.n)
		for _, p := range generator.P(c.n, vector.D(c.k)) {
			ps = append(ps, p.(simple.P))
		}

		cs := make([]*multifield.P, 0, c.n)
		for _, p := range generator.C(c.n, vector.D(c.k)) {
			cs = append(cs, p.(*multifield.P))
		}

		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New[simple.P](ps)
			}
		})

		// Check tree construction times with more complex point.P
		// implementations.
		b.Run(fmt.Sprintf("%v/MultiField", c.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New[*multifield.P](cs)
			}
		})
	}
}
