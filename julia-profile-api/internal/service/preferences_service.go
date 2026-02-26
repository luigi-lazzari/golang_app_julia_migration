package service

import (
	"context"

	"github.com/comune-roma/bff-julia-profile-api/internal/client"
	"github.com/comune-roma/bff-julia-profile-api/internal/config"
	"github.com/comune-roma/bff-julia-profile-api/internal/model"
	"github.com/comune-roma/bff-julia-profile-api/internal/repository"
	"go.uber.org/zap"
)

// UserPreferencesService handles business logic for user preferences
type UserPreferencesService struct {
	appConfigClient    interface{} // Azure App Config client
	repo               *repository.UserPreferencesRepository
	notificationClient *client.NotificationClient
	cfg                *config.Config
	log                *zap.Logger
}

// NewUserPreferencesService creates a new UserPreferencesService
func NewUserPreferencesService(appConfigClient interface{}, repo *repository.UserPreferencesRepository, notificationClient *client.NotificationClient, cfg *config.Config, log *zap.Logger) *UserPreferencesService {
	return &UserPreferencesService{
		appConfigClient:    appConfigClient,
		repo:               repo,
		notificationClient: notificationClient,
		cfg:                cfg,
		log:                log,
	}
}

// GetChatPreferences retrieves user chat preferences
func (s *UserPreferencesService) GetChatPreferences(ctx context.Context, userID string) (*model.ChatPreferences, error) {
	s.log.Info("Fetching chat preferences", zap.String("userID", userID))

	// Get saved preferences from repo (mocked for now)
	userPrefs := []string{} // Placeholder for IDs of enabled preferences from DB
	customPrefDesc := ""    // Placeholder for custom preference description

	// Default preferences from configuration
	defaultPrefs := make([]model.UserPreference, len(s.cfg.Defaults.Chat))
	for i, d := range s.cfg.Defaults.Chat {
		defaultPrefs[i] = model.UserPreference{
			ID:       d.ID,
			Category: model.UserPreferenceCategory(d.Category),
		}
	}

	merged := s.mergeChatPreferences(userPrefs, defaultPrefs)

	var customPrefs []model.CustomPreference
	if customPrefDesc != "" {
		customPrefs = append(customPrefs, model.CustomPreference{Description: customPrefDesc})
	}

	return &model.ChatPreferences{
		Preferences:       merged,
		CustomPreferences: customPrefs,
	}, nil
}

func (s *UserPreferencesService) mergeChatPreferences(userPrefIDs []string, defaultPrefs []model.UserPreference) []model.UserPreference {
	merged := make([]model.UserPreference, len(defaultPrefs))

	userPrefMap := make(map[string]bool)
	for _, id := range userPrefIDs {
		userPrefMap[id] = true
	}

	for i, dp := range defaultPrefs {
		dp.Enabled = userPrefMap[dp.ID]
		merged[i] = dp
	}

	return merged
}

// UpdateChatPreferences updates user chat preferences
func (s *UserPreferencesService) UpdateChatPreferences(ctx context.Context, userID string, req *model.ChatPreferences) (*model.ChatPreferences, error) {
	s.log.Info("Updating chat preferences", zap.String("userID", userID))

	// TODO: Save to Cosmos DB via repo
	// s.repo.UpdateChatPreferences(ctx, userID, req)

	s.log.Info("Chat preferences updated successfully")

	// Sync to Notification Service (fire and forget)
	enabledIDs := []string{}
	for _, p := range req.Preferences {
		if p.Enabled {
			enabledIDs = append(enabledIDs, p.ID)
		}
	}
	go s.notificationClient.SyncUserPreferences(context.Background(), "it", enabledIDs)

	return s.GetChatPreferences(ctx, userID)
}

// GetPreferredLanguage retrieves user's preferred language
func (s *UserPreferencesService) GetPreferredLanguage(ctx context.Context, userID string) (*model.LanguagePreference, error) {
	s.log.Info("Fetching preferred language", zap.String("userID", userID))
	// TODO: Fetch from DB
	return &model.LanguagePreference{Language: "it-IT"}, nil
}

// UpdatePreferredLanguage updates user's preferred language
func (s *UserPreferencesService) UpdatePreferredLanguage(ctx context.Context, userID string, req *model.LanguagePreference) (*model.LanguagePreference, error) {
	s.log.Info("Updating preferred language", zap.String("userID", userID), zap.String("language", req.Language))
	// TODO: Save to DB
	return req, nil
}

// GetNotificationPreferences retrieves user notification preferences
func (s *UserPreferencesService) GetNotificationPreferences(ctx context.Context, userID string) (*model.NotificationPreferences, error) {
	s.log.Info("Fetching notification preferences", zap.String("userID", userID))
	// TODO: Fetch from Notification Service or DB
	// Default notifications from configuration
	notifications := make([]model.NotificationPreferenceItem, len(s.cfg.Defaults.Notifications))
	for i, id := range s.cfg.Defaults.Notifications {
		notifications[i] = model.NotificationPreferenceItem{
			ID:      id,
			Enabled: true, // Defaulting to true as in Java properties for now
		}
	}

	return &model.NotificationPreferences{
		Notifications: notifications,
		Language:      "it-IT",
	}, nil
}

// UpdateNotificationPreferences updates user notification preferences
func (s *UserPreferencesService) UpdateNotificationPreferences(ctx context.Context, userID string, req *model.NotificationPreferences) (*model.NotificationPreferences, error) {
	s.log.Info("Updating notification preferences", zap.String("userID", userID))
	// TODO: Sync with Notification Service
	return req, nil
}

// UpsertInstallation registers or updates a device installation
func (s *UserPreferencesService) UpsertInstallation(ctx context.Context, userID, installationID string, req *model.DeviceInstallationRequest) error {
	s.log.Info("Upserting installation", zap.String("userID", userID), zap.String("installationID", installationID))
	// TODO: Call Notification Service to register installation
	return nil
}

// DeleteInstallation removes a device installation
func (s *UserPreferencesService) DeleteInstallation(ctx context.Context, userID, installationID string) error {
	s.log.Info("Deleting installation", zap.String("userID", userID), zap.String("installationID", installationID))
	// TODO: Call Notification Service to delete installation
	return nil
}
