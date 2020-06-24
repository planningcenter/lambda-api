package middleware

import (
	"log"
	"net/http"

	api "github.com/planningcenter/lambda-api"
)

// Recovery handles any panic that is thrown later in the middleware
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				api.Abort()
				log.Printf("API-Recovery: %v", rec)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
