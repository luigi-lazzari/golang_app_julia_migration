package service

import (
	"context"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/comune-roma/bff-julia-mobile-api/internal/config"
	"github.com/comune-roma/bff-julia-mobile-api/internal/model"
	"github.com/comune-roma/bff-julia-mobile-api/internal/repository"
	"go.uber.org/zap"
)

// AppConfigService handles business logic for app configuration
type AppConfigService struct {
	appConfigClient interface{} // Azure App Config client
	repo            *repository.CosmosRepository
	cfg             *config.Config
	log             *zap.Logger
}

// NewAppConfigService creates a new AppConfigService
func NewAppConfigService(appConfigClient interface{}, repo *repository.CosmosRepository, cfg *config.Config, log *zap.Logger) *AppConfigService {
	return &AppConfigService{
		appConfigClient: appConfigClient,
		repo:            repo,
		cfg:             cfg,
		log:             log,
	}
}

// GetAppConfig retrieves app configuration based on platform and version
func (s *AppConfigService) GetAppConfig(ctx context.Context, platformStr, versionStr, requestID, correlationID string) (*model.AppConfigResponse, error) {
	s.log.Info("Fetching app config",
		zap.String("platform", platformStr),
		zap.String("version", versionStr),
		zap.String("requestID", requestID),
		zap.String("correlationID", correlationID),
	)

	platform := model.AppPlatform(platformStr)

	// TODO: Fetch configuration from Azure App Configuration
	// For now, using placeholders as in Java before gateway implementation

	// Mocking configuration values that would come from Azure App Config
	minVersionStr := "1.0.0"
	latestVersionStr := "1.2.0"
	storeURL := "https://apps.apple.com/app/id123456789"
	if platform == model.PlatformAndroid {
		storeURL = "https://play.google.com/store/apps/details?id=com.example.app"
	}

	maintenanceEnabled := false
	retryAfter := 3600

	response := &model.AppConfigResponse{
		ServerTime: time.Now(),
		Maintenance: model.MaintenanceStatus{
			Enabled:           maintenanceEnabled,
			RetryAfterSeconds: &retryAfter,
		},
		Update:   s.buildUpdatePolicy(platform, versionStr, minVersionStr, latestVersionStr, storeURL),
		Config:   s.cfg.Defaults.Config,
		Locale:   s.cfg.Defaults.Locale,
		Features: s.cfg.Defaults.Features,
	}

	return response, nil
}

func (s *AppConfigService) buildUpdatePolicy(platform model.AppPlatform, version, minVersion, latestVersion, storeURL string) model.UpdatePolicy {
	action := s.computeUpdateAction(version, minVersion, latestVersion)

	return model.UpdatePolicy{
		StoreURL: storeURL,
		Action:   action,
	}
}

func (s *AppConfigService) computeUpdateAction(versionStr, minVersionStr, latestVersionStr string) model.UpdateAction {
	v, err := semver.NewVersion(versionStr)
	if err != nil {
		s.log.Error("Error parsing app version", zap.String("version", versionStr), zap.Error(err))
		return model.ActionNone
	}

	if minVersionStr != "" {
		minV, err := semver.NewVersion(minVersionStr)
		if err == nil && v.LessThan(minV) {
			return model.ActionRequire
		}
	}

	if latestVersionStr != "" {
		latestV, err := semver.NewVersion(latestVersionStr)
		if err == nil && v.LessThan(latestV) {
			return model.ActionRecommend
		}
	}

	return model.ActionNone
}

// ValidateVersion validates if the app version meets minimum requirements
func (s *AppConfigService) ValidateVersion(appVersion, minVersion string) (bool, error) {
	v, err := semver.NewVersion(appVersion)
	if err != nil {
		return false, err
	}
	minV, err := semver.NewVersion(minVersion)
	if err != nil {
		return false, err
	}

	if v.LessThan(minV) {
		return false, nil
	}
	return true, nil
}
