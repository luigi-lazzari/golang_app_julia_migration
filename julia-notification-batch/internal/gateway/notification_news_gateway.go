package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"julia-notification-batch/internal/config"
	"julia-notification-batch/internal/service"
	"net/http"
)

type NotificationNewsGateway struct {
	client  *http.Client
	baseURL string
}

func NewNotificationNewsGateway(cfg config.RestService) *NotificationNewsGateway {
	return &NotificationNewsGateway{
		client: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		baseURL: cfg.BaseURL,
	}
}

func (g *NotificationNewsGateway) UpdateNotificationNewsPreferences(news []service.NewsItem) error {
	url := fmt.Sprintf("%s/api/v1/news", g.baseURL)

	body, err := json.Marshal(news)
	if err != nil {
		return fmt.Errorf("error marshaling news request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling notification news API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification news API returned status: %s", resp.Status)
	}

	return nil
}
