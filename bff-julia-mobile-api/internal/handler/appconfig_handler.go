package handler

import (
	"net/http"

	"github.com/comune-roma/bff-julia-mobile-api/internal/model"
	"github.com/comune-roma/bff-julia-mobile-api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AppConfigHandler handles app configuration requests
type AppConfigHandler struct {
	service *service.AppConfigService
	log     *zap.Logger
}

// NewAppConfigHandler creates a new AppConfigHandler
func NewAppConfigHandler(service *service.AppConfigService, log *zap.Logger) *AppConfigHandler {
	return &AppConfigHandler{
		service: service,
		log:     log,
	}
}

// GetAppConfig godoc
// @Summary Get app configuration
// @Description Get application configuration based on platform and version
// @Tags appconfig
// @Accept json
// @Produce json
// @Param X-App-Platform header string true "App Platform (iOS/Android)"
// @Param X-App-Version header string true "App Version (semver)"
// @Success 200 {object} model.AppConfigResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /app-config [get]
func (h *AppConfigHandler) GetAppConfig(c *gin.Context) {
	platform := c.GetHeader("X-App-Platform")
	version := c.GetHeader("X-App-Version")
	requestID := c.GetHeader("X-Request-Id")
	correlationID := c.GetHeader("X-Correlation-Id")

	if platform == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Bad Request",
			Message: "X-App-Platform header is required",
		})
		return
	}

	if version == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Bad Request",
			Message: "X-App-Version header is required",
		})
		return
	}

	h.log.Info("Getting app config",
		zap.String("platform", platform),
		zap.String("version", version),
		zap.String("requestID", requestID),
		zap.String("correlationID", correlationID),
	)

	config, err := h.service.GetAppConfig(c.Request.Context(), platform, version, requestID, correlationID)
	if err != nil {
		h.log.Error("Failed to get app config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve app configuration",
		})
		return
	}

	c.JSON(http.StatusOK, config)
}
