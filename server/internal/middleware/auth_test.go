package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth_AllowsOptionsWithoutCredentials(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	handler := BasicAuth("user", "password", next)
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/owners", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected next handler to be called for OPTIONS request")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestBasicAuth_AllowsPublicRouteWithoutCredentials(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := BasicAuth("user", "password", next)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/actuator/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected next handler to be called for public path")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestBasicAuth_RejectsInvalidCredentials(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called with invalid credentials")
	})

	handler := BasicAuth("user", "password", next)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/owners", nil)
	req.SetBasicAuth("user", "wrong")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if got := rr.Header().Get("WWW-Authenticate"); got != `Basic realm="VetHub"` {
		t.Fatalf("expected WWW-Authenticate header to be set, got %q", got)
	}
}

func TestBasicAuth_AllowsValidCredentials(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := BasicAuth("user", "password", next)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/owners", nil)
	req.SetBasicAuth("user", "password")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected next handler to be called for valid credentials")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
