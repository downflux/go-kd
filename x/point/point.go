package point

type D int

const (
	AXIS_X D = iota
	AXIS_Y
	AXIS_Z
	AXIS_W
)

type V interface {
	X(d D) float64
	D() D
}

type P interface {
	P() V
}
