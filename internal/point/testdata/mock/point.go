// Package point implements a mock point.P with some arbitrary encapsulated
// data.
package point

import (
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"
)

type P struct {
	p    vector.V
	hash string
}

func (p P) P() vector.V  { return p.p }
func (p P) Hash() string { return p.hash }

func New(p vector.V, hash string) *P     { return &P{p: p, hash: hash} }

// CheckHash may be passed into the K-D tree API as the filter function for this
// data struct.
func CheckHash(p point.P, h string) bool { return p.(P).hash == h }
