package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/comune-roma/bff-julia-profile-api/internal/config"
	"github.com/comune-roma/bff-julia-profile-api/internal/model"
	"github.com/comune-roma/bff-julia-profile-api/internal/repository"
	"github.com/comune-roma/bff-julia-profile-api/pkg/cache"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserProfileService handles business logic for user profiles
type UserProfileService struct {
	repo  *repository.UserProfileRepository
	cache cache.Cache
	cfg   config.RedisConfig
	log   *zap.Logger
}

// NewUserProfileService creates a new UserProfileService
func NewUserProfileService(repo *repository.UserProfileRepository, cache cache.Cache, cfg config.RedisConfig, log *zap.Logger) *UserProfileService {
	return &UserProfileService{
		repo:  repo,
		cache: cache,
		cfg:   cfg,
		log:   log,
	}
}

// GetUserProfile retrieves a user's profile
func (s *UserProfileService) GetUserProfile(ctx context.Context, userID string) (*model.UserProfileResponse, error) {
	cacheKey := fmt.Sprintf("profile:%s", userID)

	// Try to get from cache if enabled
	if s.cfg.Enabled {
		val, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var profile model.UserProfileResponse
			if err := json.Unmarshal([]byte(val), &profile); err == nil {
				s.log.Debug("Profile found in cache", zap.String("userID", userID))
				return &profile, nil
			}
		}
	}

	s.log.Info("Fetching user profile from database", zap.String("userID", userID))

	// TODO: Implement actual Cosmos DB query
	// For now, return a mock response
	profile := &model.UserProfileResponse{
		ID:        uuid.New().String(),
		UserID:    userID,
		FirstName: "Mario",
		LastName:  "Rossi",
		Email:     "mario.rossi@example.com",
		Phone:     "+39 123 456 7890",
		Address: &model.Address{
			Street:     "Via Roma 1",
			City:       "Roma",
			PostalCode: "00100",
			Country:    "Italia",
		},
		CreatedAt: time.Now().Add(-365 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Store in cache if enabled
	if s.cfg.Enabled {
		data, err := json.Marshal(profile)
		if err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, time.Duration(s.cfg.TTL)*time.Second); err != nil {
				s.log.Warn("Failed to set profile in cache", zap.Error(err))
			}
		}
	}

	return profile, nil
}

// UpdateUserProfile updates a user's profile
func (s *UserProfileService) UpdateUserProfile(ctx context.Context, userID string, req *model.UpdateUserProfileRequest) (*model.UserProfileResponse, error) {
	s.log.Info("Updating user profile", zap.String("userID", userID))

	// TODO: Implement actual Cosmos DB update
	// For now, return a mock response with updated data
	profile := &model.UserProfileResponse{
		ID:        uuid.New().String(),
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: time.Now().Add(-365 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Invalidate cache if enabled
	if s.cfg.Enabled {
		cacheKey := fmt.Sprintf("profile:%s", userID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.log.Warn("Failed to invalidate profile cache", zap.Error(err))
		}
	}

	return profile, nil
}
