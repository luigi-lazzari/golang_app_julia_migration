package model

import (
	"time"
)

// AppConfigResponse represents the app configuration response
type AppConfigResponse struct {
	ServerTime  time.Time              `json:"serverTime"`
	Maintenance MaintenanceStatus      `json:"maintenance"`
	Update      UpdatePolicy           `json:"update"`
	Config      map[string]interface{} `json:"config"`
	Locale      map[string]string      `json:"locale"`
	Features    map[string]bool        `json:"features"`
}

// MaintenanceStatus represents the maintenance status
type MaintenanceStatus struct {
	Enabled           bool `json:"enabled"`
	RetryAfterSeconds *int `json:"retryAfterSeconds,omitempty"`
}

// UpdatePolicy represents the update policy
type UpdatePolicy struct {
	StoreURL string       `json:"storeUrl"`
	Action   UpdateAction `json:"action"`
}

// UpdateAction represents the update action
type UpdateAction string

const (
	ActionRequire   UpdateAction = "REQUIRE"
	ActionRecommend UpdateAction = "RECOMMEND"
	ActionNone      UpdateAction = "NONE"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// AppPlatform represents the app platform
type AppPlatform string

const (
	PlatformIOS     AppPlatform = "IOS"
	PlatformAndroid AppPlatform = "ANDROID"
)
