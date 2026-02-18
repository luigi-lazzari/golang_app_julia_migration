package gateway

import (
	"encoding/json"
	"fmt"
	"julia-notification-batch/internal/config"
	"julia-notification-batch/internal/service"
	"net/http"
)

type ExternalNewsGateway struct {
	client  *http.Client
	baseURL string
}

func NewExternalNewsGateway(cfg config.RestService) *ExternalNewsGateway {
	return &ExternalNewsGateway{
		client: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		baseURL: cfg.BaseURL,
	}
}

func (g *ExternalNewsGateway) GetNotificationPreferences() ([]service.NewsExternalItem, error) {
	url := fmt.Sprintf("%s/api/v1/external/news", g.baseURL)
	resp, err := g.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling external news API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external news API returned status: %s", resp.Status)
	}

	var news []service.NewsExternalItem
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, fmt.Errorf("error decoding external news response: %w", err)
	}

	return news, nil
}
