package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"julia-conversation-api/internal/api/models"
)

type SuperAgentClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewSuperAgentClient(baseURL string, timeout time.Duration) *SuperAgentClient {
	return &SuperAgentClient{
		baseURL: baseURL,
		httpClient: NewHeaderPropagationClient(&http.Client{
			Timeout: timeout,
		}),
	}
}

type ChatRequest struct {
	SessionID           string                 `json:"sessionId,omitempty"`
	ChannelName         string                 `json:"channelName"`
	Structured          bool                   `json:"structured"`
	ChannelCapabilities []string               `json:"channelCapabilities"`
	Message             RequestMessage         `json:"message"`
	Preferences         map[string]interface{} `json:"preferences,omitempty"`
}

type RequestMessage struct {
	Parts []RequestMessagePart `json:"parts"`
}

type RequestMessagePart struct {
	Type     string                 `json:"type"`
	Text     string                 `json:"text,omitempty"`
	Location *models.GeoCoordinates `json:"location,omitempty"`
}

type ChatResponse struct {
	SessionID string                       `json:"sessionId"`
	Message   models.MessageStructuredData `json:"message"`
}

func (c *SuperAgentClient) ConversationInteract(ctx context.Context, request models.ConversationRequest, preferences map[string]interface{}) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/chat", c.baseURL)

	parts := []RequestMessagePart{
		{
			Type: "text",
			Text: request.Message.Content,
		},
	}

	if request.Location != nil {
		parts = append(parts, RequestMessagePart{
			Type:     "location",
			Location: request.Location,
		})
	}

	chatReq := ChatRequest{
		SessionID:           request.ConversationID,
		ChannelName:         "julia-app",
		Structured:          true,
		ChannelCapabilities: []string{"Location"},
		Message: RequestMessage{
			Parts: parts,
		},
		Preferences: preferences,
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// ... rest of the file ...
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}
