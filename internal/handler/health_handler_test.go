package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_Health_NoDB(t *testing.T) {
	h := NewHealthHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	h.Health(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("status want 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type want application/json, got %s", rr.Header().Get("Content-Type"))
	}
	body := rr.Body.String()
	if body == "" || body == "null" {
		t.Error("expected JSON body with status")
	}
}
