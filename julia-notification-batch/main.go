package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"julia-notification-batch/internal/config"
	"julia-notification-batch/internal/gateway"
	"julia-notification-batch/internal/jobs"
	"julia-notification-batch/internal/scheduler"
	"julia-notification-batch/internal/service"
)

func main() {
	log.Println("Initializing Julia Notification Batch (Go version)...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Gateways
	extGateway := gateway.NewExternalNewsGateway(cfg.Rest.External)
	intGateway := gateway.NewNotificationNewsGateway(cfg.Rest.Notification)

	// Initialize Service
	orchestrator := service.NewOrchestratorService(extGateway, intGateway)

	// Initialize job
	maxRetries := cfg.Batch.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3 // Default
	}
	job := jobs.NewNotificationJob(orchestrator, maxRetries)

	// Initialize scheduler
	s := scheduler.NewScheduler()

	// Add job to scheduler
	cronExpr := cfg.Batch.Cron
	if cronExpr == "" {
		log.Println("No cron expression found in config, using default (every minute)")
		cronExpr = "* * * * *"
	}

	_, err := s.AddJob(cronExpr, job)
	if err != nil {
		log.Fatalf("Error adding job to scheduler: %v", err)
	}

	// Start scheduler
	s.Start()

	// Keep application running until signal received
	log.Println("Application is running. Press CTRL+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down application...")
	s.Stop()
	log.Println("Shutdown complete.")
}
