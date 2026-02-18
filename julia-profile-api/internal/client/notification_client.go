package client

import (
	"context"

	"github.com/comune-roma/bff-julia-profile-api/pkg/resilience"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// NotificationClient handles communication with the Notification Service
type NotificationClient struct {
	baseURL string
	cb      *gobreaker.CircuitBreaker
	log     *zap.Logger
}

// NewNotificationClient creates a new NotificationClient
func NewNotificationClient(baseURL string, log *zap.Logger) *NotificationClient {
	return &NotificationClient{
		baseURL: baseURL,
		cb:      resilience.NewCircuitBreaker("notification-service"),
		log:     log,
	}
}

// SyncUserPreferences synchronizes user preferences with the Notification Service
func (c *NotificationClient) SyncUserPreferences(ctx context.Context, language string, topics []string) error {
	_, err := c.cb.Execute(func() (interface{}, error) {
		c.log.Info("Syncing user preferences to Notification Service",
			zap.String("language", language),
			zap.Int("topicsCount", len(topics)),
		)

		// TODO: Implement actual HTTP call to Notification Service
		// In a real implementation:
		// resp, err := http.Post(c.baseURL + "/sync", ...)
		// if err != nil || resp.StatusCode >= 500 { return nil, err }

		return nil, nil
	})

	if err != nil {
		c.log.Warn("Notification service sync failed or circuit open", zap.Error(err))
	}

	return err
}
