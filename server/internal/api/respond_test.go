package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON_WritesStatusHeaderAndBody(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	JSON(rr, http.StatusCreated, map[string]string{"ok": "yes"})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json, got %q", got)
	}
	if body := rr.Body.String(); !strings.Contains(body, `"ok":"yes"`) {
		t.Fatalf("expected json body to contain payload, got %q", body)
	}
}

func TestErrorHelpers_WriteExpectedStatusCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		call       func(http.ResponseWriter)
		wantStatus int
	}{
		{
			name:       "Error helper",
			call:       func(w http.ResponseWriter) { Error(w, http.StatusConflict, "ERR-X", "boom") },
			wantStatus: http.StatusConflict,
		},
		{
			name:       "NotFound helper",
			call:       func(w http.ResponseWriter) { NotFound(w, "ERR-404", "missing") },
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "BadRequest helper",
			call:       func(w http.ResponseWriter) { BadRequest(w, "ERR-400", "invalid") },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rr := httptest.NewRecorder()
			tt.call(rr)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}
			if body := rr.Body.String(); !strings.Contains(body, `"code"`) || !strings.Contains(body, `"message"`) {
				t.Fatalf("expected error json body, got %q", body)
			}
		})
	}
}

func TestNoContent_Writes204(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	NoContent(rr)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestDecodeJSON_DecodesPayload(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test", strings.NewReader(`{"name":"leo"}`))
	var payload struct {
		Name string `json:"name"`
	}

	if err := DecodeJSON(req, &payload); err != nil {
		t.Fatalf("expected decode success, got error: %v", err)
	}
	if payload.Name != "leo" {
		t.Fatalf("expected decoded name leo, got %q", payload.Name)
	}
}

func TestPathID_ReturnsPathValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pets/12", nil)
	req.SetPathValue("id", "12")

	if got := PathID(req, "id"); got != "12" {
		t.Fatalf("expected path value 12, got %q", got)
	}
}

func TestDecodeJSON_ReturnsErrorOnInvalidJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test", io.NopCloser(strings.NewReader(`{"name"`)))
	var payload map[string]any
	if err := DecodeJSON(req, &payload); err == nil {
		t.Fatalf("expected decode error for invalid JSON")
	}
}
