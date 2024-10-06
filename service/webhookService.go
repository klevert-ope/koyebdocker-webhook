package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"koyebdocker-webhook/config"
	"koyebdocker-webhook/model"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	koyebAPIBaseURL  = "https://app.koyeb.com/v1/services/"
	koyebAPITokenEnv = "KOYEB_API_TOKEN"
)

// HandleWebhook processes the incoming webhook payload
func HandleWebhook(payload model.DockerHubPayload) error {
	imageName := fmt.Sprintf("%s/%s:%s", payload.Repository.Namespace, payload.Repository.Name, payload.PushData.Tag)
	log.Printf("Processing webhook for image: %s", imageName)

	serviceID, exists := config.Services[imageName]
	if !exists {
		log.Printf("No service configured for image: %s", imageName)
		return fmt.Errorf("no service configured for this image")
	}

	log.Printf("Received webhook for image: %s, triggering redeployment for service ID: %s", imageName, serviceID)
	return triggerRedeployment(serviceID, imageName)
}

// triggerRedeployment triggers a redeployment for the given service ID and image name
func triggerRedeployment(serviceID, imageName string) error {
	koyebAPIToken := os.Getenv(koyebAPITokenEnv)
	if koyebAPIToken == "" {
		log.Printf("%s environment variable is not set", koyebAPITokenEnv)
		return fmt.Errorf("%s environment variable is not set", koyebAPITokenEnv)
	}

	koyebAPIURL := fmt.Sprintf("%s%s/redeploy", koyebAPIBaseURL, serviceID)

	payload := map[string]string{"image": imageName}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload for service %s: %v", serviceID, err)
		return fmt.Errorf("failed to marshal payload for service %s: %w", serviceID, err)
	}

	req, err := http.NewRequest("POST", koyebAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Failed to create request for service %s: %v", serviceID, err)
		return fmt.Errorf("failed to create request for service %s: %w", serviceID, err)
	}
	req.Header.Set("Authorization", "Bearer "+koyebAPIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request for service %s: %v", serviceID, err)
		return fmt.Errorf("failed to send request for service %s: %w", serviceID, err)
	}
	defer closeResponseBody(resp.Body, serviceID)

	if resp.StatusCode != http.StatusOK {
		return handleNonOKResponse(resp, serviceID)
	}

	log.Printf("Successfully triggered redeployment for service ID: %s with image: %s", serviceID, imageName)
	return nil
}

// closeResponseBody closes the response body and logs any error
func closeResponseBody(body io.Closer, serviceID string) {
	if err := body.Close(); err != nil {
		log.Printf("Error closing response body for service %s: %v", serviceID, err)
	}
}

// handleNonOKResponse handles non-OK responses from the Koyeb API
func handleNonOKResponse(resp *http.Response, serviceID string) error {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Koyeb API returned status %d for service %s, and failed to read response body: %v", resp.StatusCode, serviceID, err)
		return fmt.Errorf("koyeb API returned status %d for service %s, and failed to read response body: %w", resp.StatusCode, serviceID, err)
	}
	log.Printf("Koyeb API returned status %d for service %s: %s", resp.StatusCode, serviceID, string(respBody))
	return fmt.Errorf("koyeb API returned status %d for service %s: %s", resp.StatusCode, serviceID, string(respBody))
}
