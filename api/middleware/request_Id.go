package middleware

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"tapera/util/constant"
)

// AddRequestID func
func AddRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewV4().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, constant.RequestIDKey, id)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
