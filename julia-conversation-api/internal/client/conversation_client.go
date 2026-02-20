package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"julia-conversation-api/internal/api/models"
)

type ConversationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewConversationClient(baseURL string, timeout time.Duration) *ConversationClient {
	return &ConversationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *ConversationClient) GetConversation(id string, limit, offset int, jwt string) (*models.ConversationPage, error) {
	url := fmt.Sprintf("%s/v1/conversations/%s?limit=%d&offset=%d", c.baseURL, id, limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var page models.ConversationPage
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}

func (c *ConversationClient) ListUserConversations(limit, offset int, jwt string) (*models.ConversationSummaryPage, error) {
	url := fmt.Sprintf("%s/v1/conversations/user?limit=%d&offset=%d", c.baseURL, limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var page models.ConversationSummaryPage
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}

func (c *ConversationClient) AssociateConversation(id string, jwt string) error {
	url := fmt.Sprintf("%s/v1/conversations/%s/associate", c.baseURL, id)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *ConversationClient) GenerateSuggestions(id string, jwt string, preferences map[string][]string) ([]models.Suggestion, error) {
	url := fmt.Sprintf("%s/v1/conversations/%s/suggestions", c.baseURL, id)

	payload := struct {
		ConversationID string              `json:"conversationId"`
		Preferences    map[string][]string `json:"preferences,omitempty"`
	}{
		ConversationID: id,
		Preferences:    preferences,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var suggestions []models.Suggestion
	if err := json.NewDecoder(resp.Body).Decode(&suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}
