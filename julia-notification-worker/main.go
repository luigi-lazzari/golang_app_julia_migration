package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"julia-notification-worker/internal/config"
	"julia-notification-worker/internal/service"
	"julia-notification-worker/internal/servicebus"
	"julia-notification-worker/internal/worker"
)

func main() {
	log.Println("Initializing Julia Notification Worker (Go version)...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize services
	dedupeService := service.NewDeduplicationService()
	hubService := service.NewNotificationHubService(cfg.NotificationHub, dedupeService)

	// Initialize message processor
	processor := worker.NewServiceBusNotificationProcessor(hubService)

	// Initialize Service Bus client
	sbClient, err := servicebus.NewClient(
		cfg.ServiceBus.ConnectionString,
		cfg.ServiceBus.TopicName,
		cfg.ServiceBus.SubscriptionName,
		cfg.ServiceBus.MaxConcurrentCalls,
		processor,
	)
	if err != nil {
		log.Fatalf("Error initializing Service Bus client: %v", err)
	}
	defer sbClient.Close()

	// Create a context that is cancelled on OS interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received...")
		cancel()
	}()

	log.Println("Worker is running. Listening for messages...")
	if err := sbClient.Start(ctx); err != nil {
		log.Printf("Worker stopped with error: %v", err)
	}

	log.Println("Shutdown complete.")
}
