// Package vector defines a vector interface for the K-D tree.
package vector

type V interface {
	D() int
	X(i int) float64
}
