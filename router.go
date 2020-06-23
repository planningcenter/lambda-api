package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MiddlewareFunc is the HTTP middleware handler
type MiddlewareFunc func(http.Handler) http.Handler

// Router represents the interface for a router
type Router interface {
	Handle(string, string, http.HandlerFunc) // Adds a HTTP handler

	Group(string) Router // Creates a new sub-router hash

	Add(...MiddlewareFunc) // Adds a new middleware
}

type subRouter struct {
	mux *mux.Router

	mid []MiddlewareFunc
}

func (s *subRouter) Add(m ...MiddlewareFunc) {
	s.mid = append(s.mid, m...)
}

func (s *subRouter) Handle(method, path string, fn http.HandlerFunc) {
	handler := http.Handler(http.HandlerFunc(fn))

	// Build middleware stack handler in reverse order
	for i := len(s.mid) - 1; i >= 0; i-- {
		handler = s.mid[i](handler)
	}

	muxHandler := func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}

	s.mux.HandleFunc(path, muxHandler).Methods(method)
}

func (s *subRouter) Group(path string) Router {
	return &subRouter{mux: s.mux.PathPrefix(path).Subrouter(), mid: s.mid}
}
