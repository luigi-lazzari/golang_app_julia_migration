package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"julia-conversation-api/internal/client"
	"julia-conversation-api/internal/config"
	"julia-conversation-api/internal/handler"
	"julia-conversation-api/internal/middleware"
	"julia-conversation-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize Clients
	timeout := 30 * time.Second
	convClient := client.NewConversationClient(cfg.Services.ConversationService, timeout)
	saClient := client.NewSuperAgentClient(cfg.Services.SuperAgent, timeout)
	profileClient := client.NewProfileClient(cfg.Services.ProfileService, timeout)

	// Initialize Services
	convService := service.NewConversationService(convClient, saClient, profileClient)

	// Initialize Handlers
	convHandler := handler.NewConversationHandler(convService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Middlewares
	router.Use(middleware.RequestHeaders())

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API Routes
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Auth(cfg.Auth))
	{
		// Conversations
		v1.POST("/conversations", convHandler.ConversationInteract)
		v1.GET("/conversations/:conversationId", convHandler.GetConversation)
		v1.DELETE("/conversations/:conversationId", convHandler.DeleteConversation)
		v1.GET("/conversations/:conversationId/suggestions", convHandler.GetSuggestions)

		// Users
		v1.GET("/users/me/conversations", convHandler.GetUserConversations)
		v1.POST("/users/me/conversations/associate", convHandler.AssociateUserConversation)
	}

	// Start Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
