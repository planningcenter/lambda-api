package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	loggerContextValueKey = &struct{}{}
)

// Logger handles HTTP requests by logging the HTTP request. Additionally, it
// assigns a log.Logger instance into the context before passing it to the next
// middleware.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		id := GetRequestID(r.Context())

		logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", id), log.LstdFlags|log.Lmsgprefix)

		defer func(s time.Time) {
			logger.Println(r.URL.Path, time.Since(start).String())
		}(start)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), loggerContextValueKey, logger)))
	})
}

// GetLogger returns the logger stored in the context. It will panic if the
// logger instance can't be found in given the context.
func GetLogger(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerContextValueKey).(*log.Logger); ok {
		return logger
	}
	panic("attempting to find non-existing logger")
}
