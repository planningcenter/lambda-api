package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/planningcenter/lambda-api"
	"github.com/planningcenter/lambda-api/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	a := api.New()

	a.Draw(func(r api.Router) {
		r.Add(middleware.Recovery, middleware.RequestID, middleware.Logger)

		v1 := r.Group("/v1")

		v1.Handle("GET", "/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		v1.Handle("GET", "/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("testing")
		})
	})

	t.Run("http no-content", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/v1/test", nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("panic recovery", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/v1/panic", nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
	})

	t.Run("404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/v1/foo-bar", nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "application/json", w.HeaderMap.Get("Content-Type"))
		assert.Equal(t, "85", w.HeaderMap.Get("Content-Length"))
	})
}

func ExampleAPI() {
	app := api.New()

	app.Draw(func(router api.Router) {
		router.Add(middleware.Recovery, middleware.RequestID, middleware.Logger)

		router.Handle("GET", "/", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	server.ListenAndServe()
}
