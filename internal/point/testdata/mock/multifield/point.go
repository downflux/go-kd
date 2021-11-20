package point

import (
	"time"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"
)

const (
	// etag is base58-encoded "flames on the side of my face".
	etag = "5ege8ja1nQTnQ8rvMtbyJ9FKgTtibTPJc4CUKJZa"
)

var _ point.P = &P{}

type P struct {
	id        int
	comment   string
	etag      string
	modified  bool
	metadata  map[string]string
	timestamp time.Time

	p vector.V
}

func New(p vector.V, name string) *P {
	return &P{
		p:       p,
		comment: name,
		timestamp: time.Now(),
		metadata: map[string]string{
			"Mrs. Peacock":    "Eileen Brennan",
			"Mrs. White":      "Madeline Kahn",
			"Professor Plum":  "Christopher Lloyd",
			"Mr. Green":       "Michael McKean",
			"Colonel Mustard": "Martin Mull",
			"Miss Scarlet":    "Lesley Ann Warren",

			"Mr. Boddy": "Lee Ving",
			"Yvette":    "Colleen Camp",
			"Wadsworth": "SPAAAAAAAACE",
		},

		etag: etag,
	}
}

func (p *P) P() vector.V { return p.p }
