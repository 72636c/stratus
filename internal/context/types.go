package context

import (
	"context"
)

type (
	Context = context.Context
)

var (
	Background  = context.Background
	WithTimeout = context.WithTimeout
	WithValue   = context.WithValue
)
