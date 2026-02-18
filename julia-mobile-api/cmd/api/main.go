package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/comune-roma/bff-julia-mobile-api/internal/config"
	"github.com/comune-roma/bff-julia-mobile-api/internal/handler"
	"github.com/comune-roma/bff-julia-mobile-api/internal/middleware"
	"github.com/comune-roma/bff-julia-mobile-api/internal/repository"
	"github.com/comune-roma/bff-julia-mobile-api/internal/service"
	"github.com/comune-roma/bff-julia-mobile-api/pkg/azure"
	"github.com/comune-roma/bff-julia-mobile-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// @title BFF Julia Mobile API
// @version 1.0.0
// @description Backend For Frontend API for Julia Mobile Application
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

	// Initialize Azure clients
	cosmosClient, err := azure.NewCosmosClient(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Cosmos DB client", zap.Error(err))
	}

	appConfigClient, err := azure.NewAppConfigClient(cfg)
	if err != nil {
		log.Fatal("Failed to initialize App Configuration client", zap.Error(err))
	}

	// Initialize repository
	repo := repository.NewCosmosRepository(cosmosClient, cfg.CosmosDB.Database)

	// Initialize service
	appConfigService := service.NewAppConfigService(appConfigClient, repo, cfg, log)

	// Initialize handler
	appConfigHandler := handler.NewAppConfigHandler(appConfigService, log)

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
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "bff-julia-mobile-api",
		})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/app-config", appConfigHandler.GetAppConfig)
	}

	// Swagger documentation
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
