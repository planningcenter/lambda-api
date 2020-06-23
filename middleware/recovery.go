package middleware

import (
	"log"
	"net/http"

	"github.com/maddiesch/api"
)

// Recovery handles any panic that is thrown later in the middleware
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("recovery")

		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				api.Abort()
			}
		}()

		next.ServeHTTP(w, r)
	})
}
