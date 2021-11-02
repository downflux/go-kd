package point

import (
	"github.com/downflux/go-geometry/vector"
	"github.com/downflux/go-kd/point"
)

type P struct {
	v    vector.V
	hash string
}

func (p P) V() vector.V  { return p.v }
func (p P) Hash() string { return p.hash }

func New(v vector.V, hash string) *P     { return &P{v: v, hash: hash} }
func CheckHash(p point.P, h string) bool { return p.(P).hash == h }
