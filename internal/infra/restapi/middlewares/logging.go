package middlewares

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Now, you'll want to add your request-specific details
		subLog := log.With().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Logger()

		// Save the logger into the context
		ctx = subLog.WithContext(ctx)

		// And pass execution along
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
