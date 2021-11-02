package axis

import (
	"testing"

	"github.com/downflux/go-geometry/vector"
)

func TestA(t *testing.T) {
	testConfigs := []struct {
		name  string
		depth int
		want  Type
	}{
		{
			name:  "Null",
			depth: 0,
			want:  Axis_X,
		},
		{
			name:  "Simple/1",
			depth: 1,
			want:  Axis_Y,
		},
		{
			name:  "Simple/2",
			depth: 2,
			want:  Axis_X,
		},
		{
			name:  "Simple/3",
			depth: 3,
			want:  Type(1), // Axis_Y,
		},
		{
			name:  "Simple/4",
			depth: 4,
			want:  Axis_X,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := A(c.depth); got != c.want {
				t.Errorf("A() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestX(t *testing.T) {
	testConfigs := []struct {
		name string
		axis Type
		v    vector.V
		want float64
	}{
		{
			name: "X",
			axis: Axis_X,
			v:    *vector.New(1, 2),
			want: 1,
		},
		{
			name: "Y",
			axis: Axis_Y,
			v:    *vector.New(1, 2),
			want: 2,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := X(c.v, c.axis); got != c.want {
				t.Errorf("X() = %v, want = %v", got, c.want)
			}
		})
	}
}
