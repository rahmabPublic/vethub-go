package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchOrigin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		origin  string
		pattern string
		want    bool
	}{
		{name: "exact match", origin: "http://localhost:5173", pattern: "http://localhost:5173", want: true},
		{name: "exact mismatch", origin: "http://localhost:4173", pattern: "http://localhost:5173", want: false},
		{name: "wildcard match", origin: "http://localhost:4173", pattern: "http://localhost:*", want: true},
		{name: "wildcard mismatch", origin: "http://127.0.0.1:4173", pattern: "http://localhost:*", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := matchOrigin(tt.origin, tt.pattern); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestCORS_SetsHeadersForAllowedOrigin(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CORS("http://localhost:*", next)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/owners", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("expected Access-Control-Allow-Origin to be set, got %q", got)
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestCORS_DoesNotSetHeadersForDisallowedOrigin(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CORS("http://localhost:*", next)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/owners", nil)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Access-Control-Allow-Origin header, got %q", got)
	}
}

func TestCORS_OptionsReturnsNoContent(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	handler := CORS("http://localhost:*", next)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/owners", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Fatalf("expected next handler not to be called for OPTIONS")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}
