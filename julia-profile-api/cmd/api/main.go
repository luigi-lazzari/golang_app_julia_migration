package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/comune-roma/bff-julia-profile-api/docs"
	"github.com/comune-roma/bff-julia-profile-api/internal/client"
	"github.com/comune-roma/bff-julia-profile-api/internal/config"
	"github.com/comune-roma/bff-julia-profile-api/internal/handler"
	"github.com/comune-roma/bff-julia-profile-api/internal/middleware"
	"github.com/comune-roma/bff-julia-profile-api/internal/repository"
	"github.com/comune-roma/bff-julia-profile-api/internal/service"
	"github.com/comune-roma/bff-julia-profile-api/pkg/azure"
	"github.com/comune-roma/bff-julia-profile-api/pkg/cache"
	"github.com/comune-roma/bff-julia-profile-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title BFF Julia Profile API
// @version 1.0.0
// @description Backend For Frontend API for Julia User Profile Management
// @contact.name Comune di Roma
// @BasePath /api/v1
// @schemes http https
func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration", zap.Error(err))
	}

	// Initialize Azure clients
	cosmosClient, err := azure.NewCosmosClient(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Cosmos DB client", zap.Error(err))
	}

	appConfigClient, err := azure.NewAppConfigClient(cfg)
	if err != nil {
		log.Fatal("Failed to initialize App Configuration client", zap.Error(err))
	}

	// Initialize repositories
	userProfileRepo := repository.NewUserProfileRepository(cosmosClient, cfg.CosmosDB.Database)
	userPreferencesRepo := repository.NewUserPreferencesRepository(cosmosClient, cfg.CosmosDB.Database)

	// Initialize notification client
	notificationClient := client.NewNotificationClient("http://notification-service", log)

	// Initialize Redis cache
	redisCache := cache.NewRedisCache(cfg.Redis)

	// Initialize services
	userProfileService := service.NewUserProfileService(userProfileRepo, redisCache, cfg.Redis, log)
	userPreferencesService := service.NewUserPreferencesService(appConfigClient, userPreferencesRepo, notificationClient, cfg, log)

	// Initialize handlers
	profileHandler := handler.NewUserProfileHandler(userProfileService, log)
	preferencesHandler := handler.NewUserPreferencesHandler(userPreferencesService, log)
	installationHandler := handler.NewInstallationHandler(userPreferencesService, log)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Middlewares
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := http.StatusOK
		response := gin.H{
			"status":  "UP",
			"service": "bff-julia-profile-api",
			"time":    time.Now().Format(time.RFC3339),
		}

		// Check database connection
		if err := userProfileRepo.Ping(c.Request.Context()); err != nil {
			status = http.StatusServiceUnavailable
			response["status"] = "DOWN"
			response["error"] = "Database connection error"
			log.Error("Health check failed: database", zap.Error(err))
		}

		// Check Redis connection if enabled
		if cfg.Redis.Enabled {
			if err := redisCache.Ping(c.Request.Context()); err != nil {
				response["cache_status"] = "DOWN"
				log.Warn("Health check warning: Redis", zap.Error(err))
			} else {
				response["cache_status"] = "UP"
			}
		}

		c.JSON(status, response)
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Auth(cfg.Auth))
	{
		// Profile
		v1.GET("/users/me", profileHandler.GetUserProfile)
		v1.PUT("/users/me", profileHandler.UpdateUserProfile)

		// Preferences
		v1.GET("/users/me/preferences/chat", preferencesHandler.GetUserPreferences)
		v1.PUT("/users/me/preferences/chat", preferencesHandler.UpdateUserPreferences)
		v1.GET("/users/me/preferences/language", preferencesHandler.GetPreferredLanguage)
		v1.PUT("/users/me/preferences/language", preferencesHandler.SetPreferredLanguage)

		// Notifications
		v1.GET("/users/me/notifications/preferences", preferencesHandler.GetNotificationPreferences)
		v1.PUT("/users/me/notifications/preferences", preferencesHandler.UpdateNotificationPreferences)
		v1.PUT("/users/me/notifications/installations/:installationId", installationHandler.UpsertInstallation)
		v1.DELETE("/users/me/notifications/installations/:installationId", installationHandler.DeleteInstallation)
	}

	// Swagger documentation (enabled only if not in production or explicitly allowed)
	if cfg.Environment != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
