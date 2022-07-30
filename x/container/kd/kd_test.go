package kd

import (
	"github.com/downflux/go-kd/x/container"
	"github.com/downflux/go-kd/x/point/mock"
)

var _ container.C[mock.P] = &KD[mock.P]{}
