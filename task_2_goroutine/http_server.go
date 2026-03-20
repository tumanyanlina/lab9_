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

type Processor interface {
	Process(string) error
}

type DefaultProcessor struct{}

func (p *DefaultProcessor) Process(data string) error {
	log.Printf("Start processing: %s", data)
	time.Sleep(2 * time.Second)
	log.Printf("Done processing: %s", data)
	return nil
}

var bgWg sync.WaitGroup

func respondJSON(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func decodeProcessRequest(r *http.Request) (ProcessRequest, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ProcessRequest{}, err
	}
	var req ProcessRequest
	err = json.Unmarshal(body, &req)
	return req, err
}

func processHandler(processor Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		req, err := decodeProcessRequest(r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		respondJSON(w, http.StatusOK, ProcessResponse{Status: "accepted"})
		bgWg.Add(1)
		go func(data string) {
			defer bgWg.Done()
			if err := processor.Process(data); err != nil {
				log.Printf("Background processing failed: %v", err)
			}
		}(req.Data)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func runServer(addr string, handler http.Handler) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		bgWg.Wait()
		close(idleConnsClosed)
	}()
	log.Println("Server listening on", addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	log.Println("Server stopped")
}

func main() {
	processor := &DefaultProcessor{}
	mux := http.NewServeMux()
	mux.Handle("/process", loggingMiddleware(processHandler(processor)))
	runServer(":8081", mux)
}