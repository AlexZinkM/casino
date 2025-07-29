package logging

import (
	"context"
)

type Logger interface {
	Error(ctx context.Context, errs ...error)
	Info(ctx context.Context, messages ...string)
}

