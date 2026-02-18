package handler

import (
	"net/http"

	"github.com/comune-roma/bff-julia-profile-api/internal/model"
	"github.com/comune-roma/bff-julia-profile-api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InstallationHandler handles device installation requests
type InstallationHandler struct {
	service *service.UserPreferencesService // Reusing UserPreferencesService for simplicity or create a new one
	log     *zap.Logger
}

// NewInstallationHandler creates a new InstallationHandler
func NewInstallationHandler(service *service.UserPreferencesService, log *zap.Logger) *InstallationHandler {
	return &InstallationHandler{
		service: service,
		log:     log,
	}
}

// UpsertInstallation godoc
// @Summary Register or update the current device installation
// @Description Register or update the device installation of the authenticated user for push notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Param installationId path string true "Installation ID"
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Param request body model.DeviceInstallationRequest true "Installation request"
// @Success 201 "Created"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/notifications/installations/{installationId} [put]
func (h *InstallationHandler) UpsertInstallation(c *gin.Context) {
	installationID := c.Param("installationId")
	var req model.DeviceInstallationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		return
	}

	userID := "user-001" // TODO: Auth
	err := h.service.UpsertInstallation(c.Request.Context(), userID, installationID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal Error", Message: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// DeleteInstallation godoc
// @Summary Delete the current device installation
// @Description Remove the device installation of the authenticated user
// @Tags Notifications
// @Accept json
// @Produce json
// @Param installationId path string true "Installation ID"
// @Param X-App-Platform header string true "App Platform (IOS/ANDROID)"
// @Param X-App-Version header string true "App Version (semver)"
// @Param X-Request-Id header string true "Request ID"
// @Param X-Correlation-Id header string true "Correlation ID"
// @Success 204 "No Content"
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/notifications/installations/{installationId} [delete]
func (h *InstallationHandler) DeleteInstallation(c *gin.Context) {
	installationID := c.Param("installationId")
	userID := "user-001" // TODO: Auth
	err := h.service.DeleteInstallation(c.Request.Context(), userID, installationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal Error", Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
