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

func withTestServer(proc *MockProcessor, f func(url string)) {
	srv := httptest.NewServer(loggingMiddleware(processHandler(proc)))
	defer srv.Close()
	f(srv.URL)
}

func doRequest(t *testing.T, method, url string, body []byte) *http.Response {
	if method == http.MethodPost {
		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST error: %v", err)
		}
		return resp
	}
	req, _ := http.NewRequest(method, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s request error: %v", method, err)
	}
	return resp
}

func TestProcessHandlerVariants(t *testing.T) {
	mockProc := &MockProcessor{}
	cases := []struct {
		name          string
		method        string
		body          []byte
		wantStatus    int
		wantJSON      *ProcessResponse
		wantProcess   string
	}{
		{"POST valid JSON", http.MethodPost, mustMarshal(ProcessRequest{Data: "test"}), http.StatusOK, &ProcessResponse{Status: "accepted"}, "test"},
		{"GET returns 405", http.MethodGet, nil, http.StatusMethodNotAllowed, nil, ""},
		{"POST invalid JSON", http.MethodPost, []byte(`{invalid json}`), http.StatusBadRequest, nil, ""},
		{"POST empty body", http.MethodPost, nil, http.StatusBadRequest, nil, ""},
	}

	withTestServer(mockProc, func(url string) {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				resp := doRequest(t, tc.method, url, tc.body)
				defer resp.Body.Close()
				if resp.StatusCode != tc.wantStatus {
					t.Errorf("want status %d, got %d", tc.wantStatus, resp.StatusCode)
				}
				if tc.wantJSON != nil {
					var got ProcessResponse
					if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
						t.Errorf("decode error: %v", err)
					}
					if got != *tc.wantJSON {
						t.Errorf("want %+v, got %+v", *tc.wantJSON, got)
					}
				}
				if tc.wantProcess != "" {
					time.Sleep(30 * time.Millisecond)
					mockProc.mu.Lock()
					defer mockProc.mu.Unlock()
					if len(mockProc.CalledWith) == 0 || mockProc.CalledWith[len(mockProc.CalledWith)-1] != tc.wantProcess {
						t.Errorf("expected Process with %q, got %v", tc.wantProcess, mockProc.CalledWith)
					}
				}
			})
		}
	})
}

func TestBackgroundProcessing(t *testing.T) {
	mockProc := &MockProcessor{}
	withTestServer(mockProc, func(url string) {
		body := mustMarshal(ProcessRequest{Data: "background"})
		resp := doRequest(t, http.MethodPost, url, body)
		resp.Body.Close()
		time.Sleep(30 * time.Millisecond)

		mockProc.mu.Lock()
		defer mockProc.mu.Unlock()
		if len(mockProc.CalledWith) != 1 {
			t.Fatalf("want 1 process call, got %d", len(mockProc.CalledWith))
		}
		if mockProc.CalledWith[0] != "background" {
			t.Errorf("want 'background', got %q", mockProc.CalledWith[0])
		}
	})
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}