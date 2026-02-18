package main

import (
	"context"
	"fmt"
	"log"

	"julia-notification-worker/internal/config"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
	fmt.Println("Service Bus Test Publisher")

	// Load configuration
	cfg := config.LoadConfig()

	connectionString := cfg.ServiceBus.ConnectionString
	topicName := cfg.ServiceBus.TopicName

	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close(context.Background())

	sender, err := client.NewSender(topicName, nil)
	if err != nil {
		log.Fatalf("Failed to create sender: %v", err)
	}
	defer sender.Close(context.Background())

	messageBody := `{"id": "test-123", "content": "Hello from Go Publisher!", "type": "notification"}`
	contentType := "application/json"

	message := &azservicebus.Message{
		Body:        []byte(messageBody),
		ContentType: &contentType,
	}

	err = sender.SendMessage(context.Background(), message, nil)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent successfully to topic: %s\n", topicName)
}
