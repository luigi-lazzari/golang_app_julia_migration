package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"julia-conversation-api/internal/appcontext"
)

func TestHeaderPropagationRoundTripper(t *testing.T) {
	// Mock server to receive the request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the expected headers
		if r.Header.Get("X-Request-Id") != "test-request-id" {
			t.Errorf("Expected X-Request-Id header to be 'test-request-id', got '%s'", r.Header.Get("X-Request-Id"))
		}
		if r.Header.Get("X-Correlation-Id") != "test-correlation-id" {
			t.Errorf("Expected X-Correlation-Id header to be 'test-correlation-id', got '%s'", r.Header.Get("X-Correlation-Id"))
		}
		if r.Header.Get("X-App-Platform") != "ios" {
			t.Errorf("Expected X-App-Platform header to be 'ios', got '%s'", r.Header.Get("X-App-Platform"))
		}
		if r.Header.Get("Authorization") != "Bearer test-jwt-token" {
			t.Errorf("Expected Authorization header to be 'Bearer test-jwt-token', got '%s'", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a context with headers
	ctx := context.Background()
	ctx = appcontext.WithHeader(ctx, appcontext.RequestIDKey, "test-request-id")
	ctx = appcontext.WithHeader(ctx, appcontext.CorrelationIDKey, "test-correlation-id")
	ctx = appcontext.WithHeader(ctx, appcontext.AppPlatformKey, "ios")
	ctx = appcontext.WithAuthToken(ctx, "test-jwt-token")

	// Create a client with the propagation transport
	client := NewHeaderPropagationClient(&http.Client{})

	// Create a request with the context
	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
