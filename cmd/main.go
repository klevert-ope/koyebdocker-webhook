package main

import (
	"context"
	"errors"
	"koyebdocker-webhook/config"
	"koyebdocker-webhook/controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	if err := config.LoadServices(); err != nil {
		log.Fatalf("Error loading services: %v", err)
	}

	// Create a new server
	addr := ":8088"
	srv := &http.Server{
		Addr:    addr,
		Handler: controller.SetupRouter(),
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on %s: %v\n", addr, err)
		}
	}()

	// Ensure the server is ready before accepting requests
	waitForServerReady()

	// Block until a signal is received
	<-stop

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// waitForServerReady waits until the server is ready to handle requests
func waitForServerReady() {
	for {
		resp, err := http.Get("http://localhost:8088/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("Server is ready to handle requests")
			return
		}
		log.Println("Waiting for server to be ready...")
		time.Sleep(1 * time.Second)
	}
}
