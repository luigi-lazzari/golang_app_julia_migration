package models

import (
	"time"
)

// AppPlatform defines the app platform
type AppPlatform string

const (
	AppPlatformIOS     AppPlatform = "IOS"
	AppPlatformAndroid AppPlatform = "ANDROID"
)

// MessageRole defines roles in conversation items
type MessageRole string

const (
	MessageRoleAgent MessageRole = "AGENT"
	MessageRoleUser  MessageRole = "USER"
)

// Action defines conversation actions
type Action string

const (
	ActionAppTaxi  Action = "APP_TAXY"
	ActionLocation Action = "LOCATION"
	ActionAuth     Action = "AUTH"
)

// MapItemType defines the map provider
type MapItemType string

const (
	MapItemTypeGoogle MapItemType = "GOOGLE"
	MapItemTypeApple  MapItemType = "APPLE"
	MapItemTypeBing   MapItemType = "BING"
)

// ContactItemType defines the contact type
type ContactItemType string

const (
	ContactItemTypePhone  ContactItemType = "PHONE"
	ContactItemTypeLink   ContactItemType = "LINK"
	ContactItemTypeEmail  ContactItemType = "EMAIL"
	ContactItemTypeSocial ContactItemType = "SOCIAL"
)

// GeoCoordinates represents geographic coordinates
type GeoCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// MessageRawData represents plain text message data
type MessageRawData struct {
	Content string `json:"content"`
}

// MessageStructuredData represents structured message data with actions and components
type MessageStructuredData struct {
	Content    string      `json:"content"`
	Summary    string      `json:"summary,omitempty"`
	Actions    []Action    `json:"actions,omitempty"`
	Components []Component `json:"components,omitempty"`
}

// Component represents a UI component in the message
type Component struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"` // ComponentPoiData or ComponentGenericData
}

// ComponentPoiData represents data for a POI component
type ComponentPoiData struct {
	Name        string        `json:"name,omitempty"`
	Category    string        `json:"category"`
	Description string        `json:"description,omitempty"`
	Contacts    []ContactItem `json:"contacts,omitempty"`
	Maps        []MapItem     `json:"maps,omitempty"`
	Addresses   []AddressItem `json:"addresses,omitempty"`
}

// ComponentGenericData represents data for a generic component
type ComponentGenericData struct {
	Name        string        `json:"name,omitempty"`
	Category    string        `json:"category,omitempty"`
	Description string        `json:"description,omitempty"`
	Contacts    []ContactItem `json:"contacts,omitempty"`
	Maps        []MapItem     `json:"maps,omitempty"`
}

// MapItem represents a link to a map
type MapItem struct {
	Type  MapItemType `json:"type"`
	Value string      `json:"value"`
}

// ContactItem represents contact information
type ContactItem struct {
	Type  ContactItemType `json:"type"`
	Value string          `json:"value"`
}

// AddressItem represents an address with position
type AddressItem struct {
	AddressLine string         `json:"addressLine"`
	Position    GeoCoordinates `json:"position"`
}

// Message represents a message in a conversation
type Message struct {
	Type string      `json:"type"` // "message" or "structured"
	Role MessageRole `json:"role"`
	Data interface{} `json:"data"` // MessageRawData or MessageStructuredData
}

// ConversationRequest represents a request to interact with a conversation
type ConversationRequest struct {
	ConversationID string          `json:"conversationId,omitempty"`
	Message        MessageRawData  `json:"message"`
	Location       *GeoCoordinates `json:"location,omitempty"`
}

// ConversationResponse represents a response from a conversation interaction
type ConversationResponse struct {
	ConversationID string                `json:"conversationId"`
	Message        MessageStructuredData `json:"message"`
}

// ConversationSummary represents a summary of a conversation
type ConversationSummary struct {
	ConversationID string    `json:"conversationId"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ConversationSummaryPage represents a paginated list of conversation summaries
type ConversationSummaryPage struct {
	TotalElements int                   `json:"totalElements"`
	TotalPages    int                   `json:"totalPages"`
	PageSize      int                   `json:"pageSize"`
	Data          []ConversationSummary `json:"data"`
}

// ConversationPage represents a paginated list of messages in a conversation
type ConversationPage struct {
	TotalElements int       `json:"totalElements"`
	TotalPages    int       `json:"totalPages"`
	PageSize      int       `json:"pageSize"`
	Data          []Message `json:"data"`
}

// Suggestion represents a conversation suggestion
type Suggestion struct {
	Label  string `json:"label"`
	Action string `json:"action"`
}

// ConversationAssociationRequest represents a request to associate a conversation with a user
type ConversationAssociationRequest struct {
	ConversationID string `json:"conversationId"`
}

// ProblemDetails represents RFC 7807 problem details
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}
