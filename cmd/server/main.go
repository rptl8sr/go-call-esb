package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UpdatePriceRequest struct {
	ApiKey    string `json:"api_key"`
	Timestamp int64  `json:"timestamp"`
}

func updatePriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/api/updatePriceNew/error500" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/api/updatePriceNew/noresponse" {
		// Simulate network issue - no response
		return
	}

	if r.URL.Path == "/api/updatePriceNew" {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		var req UpdatePriceRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if req.ApiKey == "" || req.Timestamp == 0 {
			http.Error(w, "Missing fields in JSON body", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Prices successfully updated"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/api/updatePriceNew", updatePriceHandler)
	http.HandleFunc("/api/updatePriceNew/error500", updatePriceHandler)
	http.HandleFunc("/api/updatePriceNew/noresponse", updatePriceHandler)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on port %d...\n", port)

	server := &http.Server{
		Addr:    addr,
		Handler: nil, // Use the default DefaultServeMux
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", addr, err)
	}

	fmt.Println("Server stopped.")
}