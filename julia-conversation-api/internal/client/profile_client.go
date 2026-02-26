package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProfileClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProfileClient(baseURL string, timeout time.Duration) *ProfileClient {
	return &ProfileClient{
		baseURL: baseURL,
		httpClient: NewHeaderPropagationClient(&http.Client{
			Timeout: timeout,
		}),
	}
}

type UserPreferencesResponse struct {
	Chat          map[string][]string `json:"chat"`
	Notifications []string            `json:"notifications"`
}

func (c *ProfileClient) GetUserPreferences(ctx context.Context) (*UserPreferencesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/me/preferences", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var preferences UserPreferencesResponse
	if err := json.NewDecoder(resp.Body).Decode(&preferences); err != nil {
		return nil, err
	}

	return &preferences, nil
}
