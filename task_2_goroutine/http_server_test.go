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
	// Create the handler with loggingMiddleware (but can mute logs in test)
	handler := loggingMiddleware(http.HandlerFunc(processHandler))

	srv := httptest.NewServer(handler)
	defer srv.Close()

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

	// Ответ должен прийти сразу (не ждать 2 секунды)
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