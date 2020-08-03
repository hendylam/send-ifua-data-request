package appcontext

import (
	"context"

	"github.com/rs/zerolog"

	"tapera/util/constant"
	"tapera/util/identity"
)

var nologger *zerolog.Logger

func init() {
	lgr := zerolog.Nop()
	nologger = &lgr
}

// ReqID returns a request ID from the given context if one is present.
// Returns the empty uuid if a request ID cannot be found.
func ReqID(ctx context.Context) string {
	if ctx != nil {
		if reqID, ok := ctx.Value(constant.RequestIDKey).(string); ok {
			return reqID
		}
	}
	return ""
}

// Logger func, get logger by key from context
func Logger(ctx context.Context) *zerolog.Logger {
	if ctx != nil {
		if logger, ok := ctx.Value(constant.ContextLoggerKey).(*zerolog.Logger); ok {
			return logger
		}
	}
	return nologger
}

// LoggerWithInf func
func LoggerWithInf(ctx context.Context) (*zerolog.Logger, bool) {
	logger := Logger(ctx)
	return logger, logger != nologger
}

// User func
func User(ctx context.Context) identity.UserInfo {
	if ctx != nil {
		if user, ok := ctx.Value(constant.UserKey).(identity.UserInfo); ok {
			return user
		}
	}
	return nil
}
