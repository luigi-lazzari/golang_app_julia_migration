package appcontext

import (
	"context"
)

type contextKey string

const (
	RequestIDKey     contextKey = "X-Request-Id"
	CorrelationIDKey contextKey = "X-Correlation-Id"
	AppPlatformKey   contextKey = "X-App-Platform"
	AppVersionKey    contextKey = "X-App-Version"
	AuthTokenKey     contextKey = "Authorization"
)

// WithHeader returns a new context with the given header value.
func WithHeader(ctx context.Context, key contextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetHeader returns the header value from the context.
func GetHeader(ctx context.Context, key contextKey) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return ""
}

// WithAuthToken returns a new context with the authorization token.
func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AuthTokenKey, token)
}

// GetAuthToken returns the authorization token from the context.
func GetAuthToken(ctx context.Context) string {
	return GetHeader(ctx, AuthTokenKey)
}
