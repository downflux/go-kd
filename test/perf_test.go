package perf

import (
	"fmt"
	"sort"
	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/kd"
	"github.com/downflux/go-kd/point"
	"math/rand"
	"testing"
)

var _ point.P = P{}

type P struct {
	v vector.V
	i int
}

func (p P) P() vector.V { return p.v }

func rn() float64 { return rand.Float64()*200 - 100 }

func rv(d int) vector.V {
	xs := make([]float64, d)
	for i := 0; i < d; i++ {
		xs[i] = rn()
	}

	return vector.V(xs)
}

func rp(n int, d int) []point.P {
	ps := make([]point.P, n)
	for i := 0; i < n; i++ {
		ps[i] = P{v: rv(d), i: i}
	}
	return ps
}

func KNN(ps []point.P, v vector.V, k int) ([]point.P, error) {
	result := make([]point.P, 0, k)
	rs := map[float64][]point.P{}
	for _, p := range ps {
		d := vector.SquaredMagnitude(vector.Sub(v, p.P()))
		rs[d] = append(rs[d], p)
	}

	ks := make([]float64, len(rs))
	for k := range rs {
		ks = append(ks, k)
	}

	sort.Float64s(ks)

	for _, k := range ks {
		result = append(result, rs[k]...)
	}
	return result, nil
}

func RadialFilter(ps []point.P, c hypersphere.C, f func(datum point.P) bool) ([]point.P, error) {
	var rs []point.P
	for _, p := range ps {
		if c.In(p.P()) && f(p) {
			rs = append(rs, p)
		}
	}
	return ps, nil
}

func BenchmarkKNN(b *testing.B) {
	const k = 3
	const n = 100000

	ps := rp(n, k)

	t, _ := kd.New(ps)

	for f := 1; f <= 10; f++ {
		b.Run(fmt.Sprintf("Naive/R=%v", float64(f)/10), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				KNN(ps, *vector.New(0, 0, 0), int(n * (float64(f) / 10)))
			}
		})

		b.Run(fmt.Sprintf("KD/R=%v", float64(f)/10), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				kd.KNN(t, *vector.New(0, 0, 0), int(n * (float64(f) / 10)))
			}
		})
	}
}

/*
func BenchmarkRadialFilter(b *testing.B) {
	const k = 3
	const n = 100000

	ps := rp(n, k)
	f := func(point.P) bool { return true }

	t, _ := kd.New(ps)

	for i := 1; i <= 10; i++ {
		b.Run(fmt.Sprintf("Naive/R=%v", float64(i)/10), func(b *testing.B) {
			c := *hypersphere.New(*vector.New(0, 0, 0), 200*(float64(i)/10))
			for i := 0; i < b.N; i++ {
				RadialFilter(ps, c, f)
			}
		})

		b.Run(fmt.Sprintf("KD/R=%v", float64(i)/10), func(b *testing.B) {
			c := *hypersphere.New(*vector.New(0, 0, 0), 200*(float64(i)/10))
			for i := 0; i < b.N; i++ {
				kd.RadialFilter(t, c, f)
			}
		})
	}
}
 */

