package kd

import (
	"testing"

	"github.com/downflux/go-geometry/nd/hypersphere"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/point"

	mock "github.com/downflux/go-kd/internal/point/testdata/mock"
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
