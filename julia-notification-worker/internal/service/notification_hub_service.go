package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"julia-notification-worker/internal/config"
	"julia-notification-worker/internal/worker"
)

type NotificationHubService struct {
	cfg                  config.NotificationHubConfig
	httpClient           *http.Client
	deduplicationService *DeduplicationService
}

func NewNotificationHubService(cfg config.NotificationHubConfig, dedupeService *DeduplicationService) *NotificationHubService {
	return &NotificationHubService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.SendTimeoutSeconds) * time.Second,
		},
		deduplicationService: dedupeService,
	}
}

// SendNotification sends a template notification to Azure Notification Hub.
func (s *NotificationHubService) SendNotification(msg worker.NotificationMessage, messageId string) error {
	if !s.cfg.Enabled {
		log.Printf("Notification Hub is disabled, skipping notification: %s", messageId)
		return nil
	}

	// Step 1: Deduplication check
	if s.deduplicationService.IsDuplicate(messageId) {
		return &worker.DuplicateMessageError{MessageID: messageId}
	}

	endpoint, _, _, err := s.parseConnectionString()
	if err != nil {
		return err
	}

	// Construct REST URL: https://{namespace}.servicebus.windows.net/{hubname}/messages/?api-version=2015-01
	// For template notifications, we use the /messages/ endpoint.
	baseUrl := strings.Replace(endpoint, "sb://", "https://", 1)
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	u := fmt.Sprintf("%s%s/messages/?api-version=2015-01", baseUrl, s.cfg.HubName)

	sasToken, err := s.generateSasToken(u)
	if err != nil {
		return fmt.Errorf("failed to generate SAS token: %w", err)
	}

	// Prepare payload (template properties)
	properties := make(map[string]string)
	properties["title"] = msg.Title
	properties["message"] = msg.Body
	properties["messageId"] = messageId

	for k, v := range msg.Data {
		properties[k] = fmt.Sprintf("%v", v)
	}

	body, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal notification properties: %w", err)
	}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", sasToken)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("ServiceBusNotification-Format", "template")

	// Add TagExpression header if present
	if msg.TagExpression != "" {
		req.Header.Set("ServiceBusNotification-Tags", msg.TagExpression)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification hub returned error status: %s", resp.Status)
	}

	// Step 5: Mark as processed
	s.deduplicationService.MarkAsProcessed(messageId)

	return nil
}

func (s *NotificationHubService) generateSasToken(uri string) (string, error) {
	_, keyName, key, err := s.parseConnectionString()
	if err != nil {
		return "", err
	}

	// Target URI: convert to lowercase and URL-encode
	targetUri := strings.ToLower(url.QueryEscape(uri))

	// Expiration: 1 hour from now
	expires := time.Now().Add(time.Hour).Unix()
	toSign := fmt.Sprintf("%s\n%d", targetUri, expires)

	// HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(toSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Construct token
	token := fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%d&skn=%s",
		targetUri,
		url.QueryEscape(signature),
		expires,
		keyName)

	return token, nil
}

func (s *NotificationHubService) parseConnectionString() (endpoint, keyName, key string, err error) {
	parts := strings.Split(s.cfg.ConnectionString, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "Endpoint=") {
			endpoint = strings.TrimPrefix(part, "Endpoint=")
		} else if strings.HasPrefix(part, "SharedAccessKeyName=") {
			keyName = strings.TrimPrefix(part, "SharedAccessKeyName=")
		} else if strings.HasPrefix(part, "SharedAccessKey=") {
			key = strings.TrimPrefix(part, "SharedAccessKey=")
		}
	}

	if endpoint == "" || keyName == "" || key == "" {
		return "", "", "", fmt.Errorf("invalid connection string")
	}
	return endpoint, keyName, key, nil
}
