package wrapper

import (
	"github.com/downflux/go-kd/x/kd/mock"

	pmock "github.com/downflux/go-kd/x/point/mock"
)

var _ mock.I[pmock.P] = &T[pmock.P]{}