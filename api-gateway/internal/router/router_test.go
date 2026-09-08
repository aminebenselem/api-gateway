package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aminebenselem/api-gateway/internal/config"
)

func TestRouterProxiesRequest(t *testing.T) {
	backend := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello from backend")
		}),
	)
	defer backend.Close()

	routes := &config.RoutesConfig{
		Routes: []config.Route{
			{
				Path: "/orders",
				Targets: []string{
					backend.URL,
				},
			},
		},
	}

	r := New(routes)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/123",
		nil,
	)

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if rec.Body.String() != "hello from backend" {
		t.Fatalf(
			"unexpected response: %s",
			rec.Body.String(),
		)
	}
}
func TestRouterSkipsFailedBackend(t *testing.T) {
	healthyBackend := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}

			fmt.Fprint(w, "healthy backend")
		}),
	)
	defer healthyBackend.Close()

	failingBackend := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}

			fmt.Fprint(w, "failing backend")
		}),
	)

	routes := &config.RoutesConfig{
		Routes: []config.Route{
			{
				Path: "/orders",
				Targets: []string{
	healthyBackend.URL,
	failingBackend.URL,
},
			},
		},
	}

	r := New(routes)

	// First request goes to failingBackend.
	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/123",
		nil,
	)

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// Simulate the backend dying.
	failingBackend.Close()

	// Second request attempts the failing backend.
	rec = httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected 502, got %d",
			rec.Code,
		)
	}

	// Third request should use the healthy backend.
	rec = httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}

	if rec.Body.String() != "healthy backend" {
		t.Fatalf(
			"unexpected response: %s",
			rec.Body.String(),
		)
	}
}
func TestRateLimit(t *testing.T) {
	backend := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "backend")
		}),
	)
	defer backend.Close()

	routes := &config.RoutesConfig{
		Routes: []config.Route{
			{
				Path: "/orders",
				Targets: []string{
					backend.URL,
				},
			},
		},
	}

	r := New(routes)

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(
			http.MethodGet,
			"/orders/123",
			nil,
		)

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"request %d: expected 200, got %d",
				i+1,
				rec.Code,
			)
		}
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/123",
		nil,
	)

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected 429, got %d",
			rec.Code,
		)
	}
}