package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	requestIDContextKey = &struct{}{}
)

// RequestID adds a HTTP request ID to the context & response headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make([]byte, 4)
		rand.Read(data)
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 32)
		id := fmt.Sprintf("%s-%02x", timestamp, data)

		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDContextKey, id)))
	})
}

// GetRequestID returns the Request ID stored in the context. It will return an
// empty string if the context does not contain a request ID.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDContextKey).(string); ok {
		return id
	}
	return ""
}
