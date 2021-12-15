// Package kd exports K-D trees as a generic collection.
package kd

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/kd"
	"github.com/downflux/go-kd/point"
)

type T[T point.P] kd.T

func New[P point.P](gd []P) (*T[P], error) {
	data := make([]point.P, 0, len(gd))
	for _, d := range gd {
		data = append(data, point.P(d))
	}

	t, err := kd.New(data)
	if err != nil {
		return nil, err
	}
	gt := T[P](*t)
	return &gt, nil
}

func (gt *T[P]) K() vector.D {
	t := kd.T(*gt)
	return (&t).K()
}

func (gt *T[P]) Balance() {
	t := kd.T(*gt)
	(&t).Balance()
}

func (gt *T[P]) Insert(gd P) error {
	t := kd.T(*gt)
	return (&t).Insert(gd)
}

func (gt *T[P]) Remove(position vector.V, gf func(gd P) bool) (bool, error) {
	t := kd.T(*gt)
	f := func(d point.P) bool { return gf(d.(P)) }
	return (&t).Remove(position, f)
}

func Filter[P point.P](gt *T[P], r hyperrectangle.R, gf func(gd P) bool) ([]P, error) {
	t := kd.T(*gt)
	f := func(d point.P) bool { return gf(d.(P)) }

	d, err := kd.Filter(&t, r, f)
	if err != nil {
		return nil, err
	}

	gd := make([]P, 0, len(d))
	for _, d := range d {
		gd = append(gd, d.(P))
	}
	return gd, nil
}

func RadialFilter[P point.P](gt *T[P], c hypersphere.C, gf func(gd P) bool) ([]P, error) {
	t := kd.T(*gt)
	f := func(d point.P) bool { return gf(d.(P)) }

	d, err := kd.RadialFilter(&t, c, f)
	if err != nil {
		return nil, err
	}

	gd := make([]P, 0, len(d))
	for _, d := range d {
		gd = append(gd, d.(P))
	}
	return gd, nil
}

func KNN[P point.P](gt *T[P], position vector.V, k int) ([]P, error) {
	t := kd.T(*gt)

	d, err := kd.KNN(&t, position, k)
	if err != nil {
		return nil, err
	}

	gd := make([]P, 0, len(d))
	for _, d := range d {
		gd = append(gd, d.(P))
	}
	return gd, nil
}

func Data[P point.P](gt *T[P]) []P {
	t := kd.T(*gt)

	d := kd.Data(&t)

	gd := make([]P, 0, len(d))
	for _, d := range d {
		gd = append(gd, d.(P))
	}
	return gd
}
