package handler

import (
	"net/http"
	"strconv"

	"julia-conversation-api/internal/api/models"
	"julia-conversation-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	service *service.ConversationService
}

func NewConversationHandler(service *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service: service,
	}
}

func (h *ConversationHandler) GetConversation(c *gin.Context) {
	id := c.Param("conversationId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	jwt := c.GetString("jwt")

	result, err := h.service.GetConversation(id, page, size, jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ConversationHandler) GetUserConversations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	jwt := c.GetString("jwt")

	result, err := h.service.GetUserConversations(page, size, jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ConversationHandler) AssociateUserConversation(c *gin.Context) {
	var req models.ConversationAssociationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt := c.GetString("jwt")
	if err := h.service.AssociateUserConversation(req.ConversationID, jwt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ConversationHandler) GetSuggestions(c *gin.Context) {
	id := c.Param("conversationId")
	jwt := c.GetString("jwt")

	result, err := h.service.GetSuggestions(id, jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	id := c.Param("conversationId")
	jwt := c.GetString("jwt")

	if err := h.service.DeleteConversation(id, jwt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ConversationHandler) ConversationInteract(c *gin.Context) {
	var req models.ConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt := c.GetString("jwt")
	result, err := h.service.ConversationInteract(req, jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
