package kd

import (
	"fmt"
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/x/internal/node/util"
	"github.com/downflux/go-kd/x/point/mock"

	putil "github.com/downflux/go-kd/x/internal/validation/util"
)

func TestNew(t *testing.T) {
	type config struct {
		name string
		k    vector.D
		n    int

		size int
	}

	var configs []config
	for _, k := range putil.KRange {
		for _, n := range putil.NRange {
			for _, size := range putil.SizeRange {
				configs = append(configs, config{
					name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v", k, n, size),
					k:    k,
					n:    n,
					size: size,
				})
			}
		}
	}

	for _, c := range configs {
		ps := putil.Generate(c.n, c.k)
		t.Run(c.name, func(t *testing.T) {
			tree := New[*mock.P](O[*mock.P]{
				Data: ps,
				K:    c.k,
				N:    c.size,
			})
			if !util.Validate(tree.root) {
				t.Errorf("validate() = %v, want = %v", false, true)
			}
		})
	}
}
