package kd

import (
	"github.com/downflux/go-kd/x/container"
	"github.com/downflux/go-kd/x/point/mock"
)

var _ container.I[mock.P] = &T[mock.P]{}
