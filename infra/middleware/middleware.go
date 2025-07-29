package infra

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"casino/boundary/logging"
	"casino/utils"
)

func LoggingMiddleware(handler http.HandlerFunc, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info(r.Context(), fmt.Sprintf("Request started: %s %s",
			r.Method,
			r.URL.Path),
		)

		ctx := r.Context()
		if ctx.Value(utils.CtxKeyRequestID) == nil {
			ctx = context.WithValue(ctx, utils.CtxKeyRequestID, utils.GenerateUUID())
		}

		handler.ServeHTTP(w, r.WithContext(ctx))

		duration := time.Since(start)
		logger.Info(r.Context(), fmt.Sprintf("Request completed: %s %s - Duration: %s",
			r.Method,
			r.URL.Path,
			duration.String()),
		)
	}
}
