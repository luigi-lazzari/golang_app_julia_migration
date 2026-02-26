package client

import (
	"net/http"

	"julia-conversation-api/internal/appcontext"
)

// HeaderPropagationRoundTripper is an http.RoundTripper that propagates headers from the context.
type HeaderPropagationRoundTripper struct {
	Base http.RoundTripper
}

func (l *HeaderPropagationRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Propagate tracking headers
	if rid := appcontext.GetHeader(ctx, appcontext.RequestIDKey); rid != "" {
		req.Header.Set(string(appcontext.RequestIDKey), rid)
	}
	if cid := appcontext.GetHeader(ctx, appcontext.CorrelationIDKey); cid != "" {
		req.Header.Set(string(appcontext.CorrelationIDKey), cid)
	}
	if plat := appcontext.GetHeader(ctx, appcontext.AppPlatformKey); plat != "" {
		req.Header.Set(string(appcontext.AppPlatformKey), plat)
	}
	if ver := appcontext.GetHeader(ctx, appcontext.AppVersionKey); ver != "" {
		req.Header.Set(string(appcontext.AppVersionKey), ver)
	}

	// Propagate Authorization header
	if auth := appcontext.GetAuthToken(ctx); auth != "" {
		req.Header.Set("Authorization", "Bearer "+auth)
	}

	return l.Base.RoundTrip(req)
}

// NewHeaderPropagationClient returns an http.Client configured with the HeaderPropagationRoundTripper.
func NewHeaderPropagationClient(baseClient *http.Client) *http.Client {
	if baseClient == nil {
		baseClient = http.DefaultClient
	}

	transport := baseClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &http.Client{
		Transport: &HeaderPropagationRoundTripper{
			Base: transport,
		},
		Timeout: baseClient.Timeout,
	}
}
