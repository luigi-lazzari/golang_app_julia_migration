package model

// UserPreferenceCategory represents the category of a preference
type UserPreferenceCategory string

const (
	CategoryServices  UserPreferenceCategory = "SERVICES"
	CategoryLeisure   UserPreferenceCategory = "LEISURE"
	CategoryTransport UserPreferenceCategory = "TRANSPORT"
)

// UserPreference represents a single user preference
type UserPreference struct {
	ID          string                 `json:"id"`
	Category    UserPreferenceCategory `json:"category,omitempty"`
	Enabled     bool                   `json:"enabled"`
	Description string                 `json:"description,omitempty"`
}

// CustomPreference represents a user's custom preference
type CustomPreference struct {
	Description string `json:"description"`
}

// ChatPreferences represents the response with user preferences
type ChatPreferences struct {
	Preferences       []UserPreference   `json:"preferences"`
	CustomPreferences []CustomPreference `json:"customPreferences,omitempty"`
}

// LanguagePreference represents the user's preferred language
type LanguagePreference struct {
	Language string `json:"language,omitempty"`
}

// NotificationPreferenceItem represents a single notification setting
type NotificationPreferenceItem struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// NotificationPreferences represents the user's notification preferences
type NotificationPreferences struct {
	Notifications []NotificationPreferenceItem `json:"notifications"`
	Language      string                       `json:"language,omitempty"`
}

// UserPreferenceUpdate represents a preference update in a request (keeping it for internal use if needed, but ChatPreferences is used in API)
type UserPreferenceUpdate struct {
	ID          string                 `json:"id"`
	Category    UserPreferenceCategory `json:"category,omitempty"`
	Enabled     bool                   `json:"enabled"`
	Description string                 `json:"description,omitempty"`
}

// UpdateUserPreferencesRequest represents the request to update user preferences
type UpdateUserPreferencesRequest struct {
	Preferences       []UserPreferenceUpdate `json:"preferences"`
	CustomPreferences []CustomPreference     `json:"customPreferences,omitempty"`
}

// ProblemDetails represents RFC 7807 problem details
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
