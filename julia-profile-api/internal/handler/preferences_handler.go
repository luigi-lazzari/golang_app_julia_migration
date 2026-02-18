package handler

import (
	"net/http"

	"github.com/comune-roma/bff-julia-profile-api/internal/model"
	"github.com/comune-roma/bff-julia-profile-api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserPreferencesHandler handles user preferences requests
type UserPreferencesHandler struct {
	service *service.UserPreferencesService
	log     *zap.Logger
}

// NewUserPreferencesHandler creates a new UserPreferencesHandler
func NewUserPreferencesHandler(service *service.UserPreferencesService, log *zap.Logger) *UserPreferencesHandler {
	return &UserPreferencesHandler{
		service: service,
		log:     log,
	}
}

// GetUserPreferences godoc
// @Summary Get user preferences
// @Description Get user application preferences
// @Tags preferences
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Success 200 {object} model.ChatPreferences
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/preferences/chat [get]
func (h *UserPreferencesHandler) GetUserPreferences(c *gin.Context) {
	platform := c.GetHeader("X-App-Platform")
	version := c.GetHeader("X-App-Version")
	requestID := c.GetHeader("X-Request-Id")
	correlationID := c.GetHeader("X-Correlation-Id")

	h.log.Info("Getting user preferences",
		zap.String("platform", platform),
		zap.String("version", version),
		zap.String("requestID", requestID),
		zap.String("correlationID", correlationID),
	)

	// TODO: Extract user ID from authentication context (JWT)
	userID := "user-001"

	preferences, err := h.service.GetChatPreferences(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve user preferences",
		})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

// UpdateUserPreferences godoc
// @Summary Update user preferences
// @Description Update user application preferences
// @Tags preferences
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Param request body model.ChatPreferences true "Update preferences request"
// @Success 200 {object} model.ChatPreferences
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/preferences/chat [put]
func (h *UserPreferencesHandler) UpdateUserPreferences(c *gin.Context) {
	platform := c.GetHeader("X-App-Platform")
	version := c.GetHeader("X-App-Version")
	requestID := c.GetHeader("X-Request-Id")
	correlationID := c.GetHeader("X-Correlation-Id")

	h.log.Info("Updating user preferences",
		zap.String("platform", platform),
		zap.String("version", version),
		zap.String("requestID", requestID),
		zap.String("correlationID", correlationID),
	)

	var req model.ChatPreferences
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	// TODO: Extract user ID from authentication context
	userID := "user-001"

	preferences, err := h.service.UpdateChatPreferences(c.Request.Context(), userID, &req)
	if err != nil {
		h.log.Error("Failed to update user preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to update user preferences",
		})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

// GetPreferredLanguage godoc
// @Summary Get preferred language
// @Description Get user preferred language
// @Tags preferences
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Success 200 {object} model.LanguagePreference
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/preferences/language [get]
func (h *UserPreferencesHandler) GetPreferredLanguage(c *gin.Context) {
	userID := "user-001" // TODO: Auth
	pref, err := h.service.GetPreferredLanguage(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal Error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, pref)
}

// SetPreferredLanguage godoc
// @Summary Set preferred language
// @Description Set user preferred language
// @Tags preferences
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Param request body model.LanguagePreference true "Language preference"
// @Success 200 {object} model.LanguagePreference
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/preferences/language [put]
func (h *UserPreferencesHandler) SetPreferredLanguage(c *gin.Context) {
	var req model.LanguagePreference
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		return
	}
	userID := "user-001" // TODO: Auth
	pref, err := h.service.UpdatePreferredLanguage(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal Error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, pref)
}

// GetNotificationPreferences godoc
// @Summary Get notification preferences
// @Description Get user notification preferences
// @Tags Notifications
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Success 200 {object} model.NotificationPreferences
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/notifications/preferences [get]
func (h *UserPreferencesHandler) GetNotificationPreferences(c *gin.Context) {
	userID := "user-001" // TODO: Auth
	prefs, err := h.service.GetNotificationPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal Error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, prefs)
}

// UpdateNotificationPreferences godoc
// @Summary Update notification preferences
// @Description Update user notification preferences
// @Tags Notifications
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Param request body model.NotificationPreferences true "Notification preferences"
// @Success 200 {object} model.NotificationPreferences
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/notifications/preferences [put]
func (h *UserPreferencesHandler) UpdateNotificationPreferences(c *gin.Context) {
	var req model.NotificationPreferences
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		return
	}
	userID := "user-001" // TODO: Auth
	prefs, err := h.service.UpdateNotificationPreferences(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal Error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, prefs)
}
