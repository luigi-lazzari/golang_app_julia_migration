package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// ServiceBusNotificationDto represents the raw JSON from Service Bus
type ServiceBusNotificationDto struct {
	Title   string                 `json:"title"`
	Body    string                 `json:"body"`
	Message string                 `json:"message"` // Fallback field
	Tags    string                 `json:"tags"`    // Comma separated tags/categories
	UserID  string                 `json:"userId"`
	Data    map[string]interface{} `json:"data"`
}

type NotificationHubService interface {
	SendNotification(msg NotificationMessage, messageId string) error
}

type ServiceBusNotificationProcessor struct {
	notificationHubService NotificationHubService
}

func NewServiceBusNotificationProcessor(hubService NotificationHubService) *ServiceBusNotificationProcessor {
	return &ServiceBusNotificationProcessor{
		notificationHubService: hubService,
	}
}

func (p *ServiceBusNotificationProcessor) ProcessMessage(messageID string, contentType string, body []byte) error {
	log.Printf("Received message from Service Bus: MessageId=%s", messageID)

	var dto ServiceBusNotificationDto
	if err := json.Unmarshal(body, &dto); err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}

	// Preprocess DTO (fallback logic)
	p.preprocessDto(&dto)

	// Map DTO to internal model
	notification := NotificationMessage{
		Title:  dto.Title,
		Body:   dto.Body,
		UserID: dto.UserID,
		Data:   dto.Data,
	}

	if dto.Tags != "" {
		notification.Categories = strings.Split(dto.Tags, ",")
		for i := range notification.Categories {
			notification.Categories[i] = strings.TrimSpace(notification.Categories[i])
		}
	}

	// Validate
	if err := p.validateNotificationMessage(&notification); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	log.Printf("Sending notification to Hub: MessageId=%s, Title=%s, UserID=%s, Categories=%v", messageID, notification.Title, notification.UserID, notification.Categories)

	if err := p.notificationHubService.SendNotification(notification, messageID); err != nil {
		// Detect duplicate message error (defined in worker/model.go now)
		if _, ok := err.(*DuplicateMessageError); ok {
			log.Printf("Duplicate message skipped, will complete: MessageId=%s, reason=%v", messageID, err)
			return nil // Return nil correctly to trigger CompleteMessage in listener
		}

		log.Printf("Failed to send notification to Hub: %v", err)
		return fmt.Errorf("failed to send notification to hub: %w", err)
	}

	log.Printf("Notification sent successfully: MessageId=%s", messageID)
	return nil
}

func (p *ServiceBusNotificationProcessor) preprocessDto(dto *ServiceBusNotificationDto) {
	// Fallback: if 'body' is empty but 'message' is present, use 'message' as body
	if strings.TrimSpace(dto.Body) == "" && strings.TrimSpace(dto.Message) != "" {
		dto.Body = dto.Message
	}
}

func (p *ServiceBusNotificationProcessor) validateNotificationMessage(msg *NotificationMessage) error {
	if strings.TrimSpace(msg.Title) == "" {
		return fmt.Errorf("notification title is required")
	}
	if strings.TrimSpace(msg.Body) == "" {
		return fmt.Errorf("notification body is required")
	}
	return nil
}
