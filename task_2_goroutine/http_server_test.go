package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// --- Мок-процессор для тестирования фоновой обработки ---
type MockProcessor struct {
	CalledWith []string
	mu         sync.Mutex
}

func (m *MockProcessor) Process(data string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CalledWith = append(m.CalledWith, data)
	return nil
}

func TestProcessEndpoint(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
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
	// Ответ должен прийти сразу (не ждать 2 секунды).
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

	// Проверяем, что Process вызван с нужными данными
	time.Sleep(30 * time.Millisecond)
	mockProc.mu.Lock()
	defer mockProc.mu.Unlock()
	if len(mockProc.CalledWith) != 1 {
		t.Fatalf("Expected Process to be called once, got %d", len(mockProc.CalledWith))
	}
	if mockProc.CalledWith[0] != "test" {
		t.Errorf("Expected Process called with 'test', got %q", mockProc.CalledWith[0])
	}
}

func TestProcessMethodNotAllowed(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
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
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
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
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
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

func TestBackgroundProcessing(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	payload := ProcessRequest{Data: "background"}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to POST: %v", err)
	}
	resp.Body.Close()

	// Даем время на запуск фоновой горутины (Process вызывается асинхронно)
	time.Sleep(30 * time.Millisecond)

	mockProc.mu.Lock()
	defer mockProc.mu.Unlock()
	if len(mockProc.CalledWith) != 1 {
		t.Fatalf("process should be called once, got %d calls", len(mockProc.CalledWith))
	}
	if mockProc.CalledWith[0] != "background" {
		t.Errorf("expected process called with %q, got %q", "background", mockProc.CalledWith[0])
	}
}
