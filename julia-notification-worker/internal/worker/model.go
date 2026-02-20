package worker

import (
	"fmt"
)

// NotificationMessage represents a notification message to be sent to Azure Notification Hub.
type NotificationMessage struct {
	Title      string                 `json:"title" validate:"required"`
	Body       string                 `json:"body" validate:"required"`
	UserID     string                 `json:"userId,omitempty"`
	Categories []string               `json:"categories,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

type DuplicateMessageError struct {
	MessageID string
}

func (e *DuplicateMessageError) Error() string {
	return fmt.Sprintf("duplicate message detected: %s", e.MessageID)
}
