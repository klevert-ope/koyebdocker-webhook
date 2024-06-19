package controller

import (
	"encoding/json"
	"koyebdocker-webhook/model"
	"koyebdocker-webhook/service"
	"log"
	"net/http"
)

// SetupRouter sets up the HTTP routes
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", webhookHandler)
	mux.HandleFunc("/health", healthHandler)
	return mux
}

// webhookHandler handles incoming webhook requests
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Printf("Invalid request method: %s", r.Method)
		return
	}

	// Parse the JSON payload
	var payload model.DockerHubPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	// Handle the webhook
	if err := service.HandleWebhook(payload); err != nil {
		http.Error(w, "Failed to handle webhook: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to handle webhook: %v", err)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Webhook received successfully")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}
