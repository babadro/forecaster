package middlewares

import (
	"net/http"

	"github.com/rs/zerolog"
)

func Logging(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subLog := log.With().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Logger()

			ctx := subLog.WithContext(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
