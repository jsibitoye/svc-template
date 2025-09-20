package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/echo?msg=hi", nil)
	rr := httptest.NewRecorder()
	Echo(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var out echoResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if out.Message != "hi" {
		t.Fatalf("unexpected message: %q", out.Message)
	}
}

func TestEchoMissingParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/echo", nil)
	rr := httptest.NewRecorder()
	Echo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
