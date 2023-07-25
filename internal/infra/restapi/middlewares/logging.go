package middlewares

import (
	"github.com/rs/zerolog"
	"net/http"
)

func loggingMiddleware(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Now, you'll want to add your request-specific details
			subLog := log.With().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Logger()

			// Save the logger into the context
			ctx := subLog.WithContext(r.Context())

			// And pass execution along
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
