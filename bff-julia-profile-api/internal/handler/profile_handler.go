package handler

import (
	"net/http"

	"github.com/comune-roma/bff-julia-profile-api/internal/model"
	"github.com/comune-roma/bff-julia-profile-api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserProfileHandler handles user profile requests
type UserProfileHandler struct {
	service *service.UserProfileService
	log     *zap.Logger
}

// NewUserProfileHandler creates a new UserProfileHandler
func NewUserProfileHandler(service *service.UserProfileService, log *zap.Logger) *UserProfileHandler {
	return &UserProfileHandler{
		service: service,
		log:     log,
	}
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get user profile information
// @Tags profile
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (iOS/Android)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-ID header string false "Request ID"
// @Param X-Correlation-ID header string false "Correlation ID"
// @Success 200 {object} model.UserProfileResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /user/profile [get]
func (h *UserProfileHandler) GetUserProfile(c *gin.Context) {
	platform := c.GetHeader("X-App-Platform")
	version := c.GetHeader("X-App-Version")
	requestID := c.GetHeader("X-Request-ID")
	correlationID := c.GetHeader("X-Correlation-ID")

	h.log.Info("Getting user profile",
		zap.String("platform", platform),
		zap.String("version", version),
		zap.String("requestID", requestID),
		zap.String("correlationID", correlationID),
	)

	// TODO: Extract user ID from authentication context
	userID := "user-001"

	profile, err := h.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve user profile",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update user profile information
// @Tags profile
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (iOS/Android)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-ID header string false "Request ID"
// @Param X-Correlation-ID header string false "Correlation ID"
// @Param request body model.UpdateUserProfileRequest true "Update profile request"
// @Success 200 {object} model.UserProfileResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /user/profile [put]
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	platform := c.GetHeader("X-App-Platform")
	version := c.GetHeader("X-App-Version")
	requestID := c.GetHeader("X-Request-ID")
	correlationID := c.GetHeader("X-Correlation-ID")

	var req model.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("Updating user profile",
		zap.String("platform", platform),
		zap.String("version", version),
		zap.String("requestID", requestID),
		zap.String("correlationID", correlationID),
	)

	// TODO: Extract user ID from authentication context
	userID := "user-001"

	profile, err := h.service.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		h.log.Error("Failed to update user profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}
