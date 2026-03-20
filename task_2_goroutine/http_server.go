package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ProcessRequest struct {
	Data string `json:"data"`
}

type ProcessResponse struct {
	Status string `json:"status"`
}

// Processor interface
type Processor interface {
	Process(data string) error
}

// DefaultProcessor implementation
type DefaultProcessor struct{}

func (p *DefaultProcessor) Process(data string) error {
	log.Printf("Start processing: %s", data)
	time.Sleep(2 * time.Second)
	log.Printf("Done processing: %s", data)
	return nil
}

var (
	bgWg sync.WaitGroup
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Handler factory to use injected Processor
func processHandler(processor Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req ProcessRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Respond immediately
		resp := ProcessResponse{Status: "accepted"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

		// Process in background and track completion
		bgWg.Add(1)
		go func(data string) {
			defer bgWg.Done()
			if err := processor.Process(data); err != nil {
				log.Printf("Background processing failed: %v", err)
			}
		}(req.Data)
	}
}

func main() {
	mux := http.NewServeMux()
	processor := &DefaultProcessor{}
	mux.Handle("/process", loggingMiddleware(processHandler(processor)))
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// Signal handling for graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Stop accepting new connections, allow up to 5s for existing requests
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		// Wait for background goroutines to finish
		bgWg.Wait()

		close(idleConnsClosed)
	}()

	log.Println("Server listening on :8081")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped")
}
