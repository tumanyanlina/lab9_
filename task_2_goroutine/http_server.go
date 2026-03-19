package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ProcessRequest struct {
	Data string `json:"data"`
}

type ProcessResponse struct {
	Status string `json:"status"`
}

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

	// Process in background
	go func(data string) {
		time.Sleep(2 * time.Second)
		fmt.Printf("Processed: %s\n", data)
	}(req.Data)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/process", loggingMiddleware(http.HandlerFunc(processHandler)))
	log.Println("Server listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
