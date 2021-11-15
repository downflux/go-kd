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
func CheckHash(p point.P, h string) bool { return p.(P).hash == h }
