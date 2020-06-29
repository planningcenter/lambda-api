package middleware

import (
	"log"
	"net/http"

	api "github.com/planningcenter/lambda-api"
)

// Recovery handles any panic that is thrown later in the middleware stack. By
// default, the only thing done is to assign a HTTP status code of 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("API-Recovery: %v", rec)
				api.Abort()
			}
		}()

		next.ServeHTTP(w, r)
	})
}
