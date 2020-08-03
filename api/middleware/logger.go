package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"

	"tapera/util/appcontext"
	"tapera/util/constant"
)

// Logger func
func Logger(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := appcontext.ReqID(r.Context())
			ctxLogger := logger.With().Str("reqId", reqID).Logger()
			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.ContextLoggerKey, &ctxLogger)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
