package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProcessEndpoint(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// 1. Основной happy-path тест
	payload := ProcessRequest{Data: "test"}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	start := time.Now()
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to do POST: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	if elapsed > 300*time.Millisecond {
		t.Errorf("response took too long: %v", elapsed)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}

	var got ProcessResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	expected := ProcessResponse{Status: "accepted"}
	if got != expected {
		t.Errorf("expected response %+v, got %+v", expected, got)
	}
}

func TestProcessMethodNotAllowed(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestProcessInvalidJSON(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	invalidJSON := []byte(`{invalid json}`)
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(invalidJSON))
	if err != nil {
		t.Fatalf("failed to POST invalid JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestProcessEmptyBody(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("failed to POST with empty body: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}