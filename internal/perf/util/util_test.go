package util

import (
	"flag"
)

var (
	s            = SizeSmall
	_ flag.Value = &s
)
