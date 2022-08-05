package kyroy

import (
	"github.com/downflux/go-kd/container"
	"github.com/downflux/go-kd/point/mock"
)

var _ container.C[mock.P] = &KD[mock.P]{}
