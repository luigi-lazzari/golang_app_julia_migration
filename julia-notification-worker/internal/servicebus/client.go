package servicebus

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type MessageProcessor interface {
	ProcessMessage(messageID string, contentType string, body []byte) error
}

type Client struct {
	client   *azservicebus.Client
	receiver *azservicebus.Receiver
	handler  MessageProcessor
}

func NewClient(connectionString, topicName, subscriptionName string, maxConcurrentCalls int, handler MessageProcessor) (*Client, error) {
	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create service bus client: %w", err)
	}

	receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create receiver: %w", err)
	}

	return &Client{
		client:   client,
		receiver: receiver,
		handler:  handler,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	for {
		messages, err := c.receiver.ReceiveMessages(ctx, 10, nil)
		if err != nil {
			if ctx.Err() != nil {
				return nil // Context cancelled
			}
			log.Printf("Error receiving messages: %v", err)
			continue
		}

		for _, message := range messages {
			c.processSingleMessage(ctx, message)
		}
	}
}

func (c *Client) processSingleMessage(ctx context.Context, message *azservicebus.ReceivedMessage) {
	messageID := message.MessageID
	var contentType string
	if message.ContentType != nil {
		contentType = *message.ContentType
	}

	err := c.handler.ProcessMessage(messageID, contentType, message.Body)
	if err != nil {
		log.Printf("Error processing message %s: %v. Abandoning message for retry.", messageID, err)
		// Explicitly abandon the message so it becomes available for another worker immediately
		abandonErr := c.receiver.AbandonMessage(ctx, message, nil)
		if abandonErr != nil {
			log.Printf("Error abandoning message %s: %v", messageID, abandonErr)
		}
		return
	}

	err = c.receiver.CompleteMessage(ctx, message, nil)
	if err != nil {
		log.Printf("Error completing message %s: %v", messageID, err)
	}
}

func (c *Client) Close() error {
	return c.client.Close(context.Background())
}
