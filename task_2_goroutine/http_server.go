package main

import (
	"context"
	"encoding/json"
	"fmt"
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

var (
	bgWg sync.WaitGroup
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func processHandler(w http.ResponseWriter, r *http.Request) {
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
		time.Sleep(2 * time.Second)
		fmt.Printf("Processed: %s\n", data)
	}(req.Data)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/process", loggingMiddleware(http.HandlerFunc(processHandler)))
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
