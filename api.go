// Package api provides a simple muxing & middleware layer for a http.Handler
package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// API contains a single API's muxer, middleware, and handler.
type API struct {
	mux      *mux.Router                           // Routing mux
	mu       sync.Mutex                            // Execution mutex
	context  func(context.Context) context.Context // Provides a new context for every request
	hasDrawn bool                                  // Stores whether the Draw function has been called
}

var (
	// DefaultNotFoundHandler is the default handler for a 404 Not Found
	DefaultNotFoundHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		WriteJSON(w, map[string]interface{}{
			"Error": map[string]string{
				"Title":  http.StatusText(http.StatusNotFound),
				"Detail": "The requested resource could not be found.",
			},
		})
	}))

	// DefaultMethodNotAllowedHandler is the default handler for a 405 Method Not
	// Allowed
	DefaultMethodNotAllowedHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)

		WriteJSON(w, map[string]interface{}{
			"Error": map[string]string{
				"Title":  http.StatusText(http.StatusMethodNotAllowed),
				"Detail": "The HTTP method is not allowed.",
			},
		})
	}))
)

// New returns a new API instance with the default configuration for the
// NotFoundHandler and MethodNotAllowedHandler
func New() *API {
	router := mux.NewRouter()
	router.NotFoundHandler = DefaultNotFoundHandler
	router.MethodNotAllowedHandler = DefaultMethodNotAllowedHandler

	return &API{
		mux: router,
		context: func(ctx context.Context) context.Context {
			return ctx
		},
	}
}

// Draw allows routes to be defined on the API using the passed Router
func (a *API) Draw(fn func(Router), middleware ...MiddlewareFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.hasDrawn {
		panic("multiple calls of Draw are not supported")
	}

	router := &subRouter{mux: a.mux, mid: make([]MiddlewareFunc, 0)}
	fn(router)

	a.hasDrawn = true
}

func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var waiter sync.WaitGroup
	waiter.Add(1)

	go func() {
		a.mu.Lock()
		defer a.mu.Unlock()
		defer waiter.Done()

		a.mux.ServeHTTP(w, req.WithContext(Context(req.Context(), a)))
	}()

	waiter.Wait()
}
