package context

import (
	"github.com/72636c/stratus/internal/log"
)

func Logger(ctx Context) log.Logger {
	logger, ok := ctx.Value(loggerKey).(log.Logger)
	if ok {
		return logger
	}

	return log.StandardLogger
}

func WithLogger(ctx Context, logger log.Logger) Context {
	return WithValue(ctx, loggerKey, logger)
}
