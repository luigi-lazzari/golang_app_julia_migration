package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server      ServerConfig
	Services    ServicesConfig
	Auth        AuthConfig
	Environment string
	LogLevel    string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
}

// ServicesConfig holds downstream service URLs
type ServicesConfig struct {
	SuperAgent          string
	ConversationService string
	ProfileService      string
}

// AuthConfig holds authentication settings
type AuthConfig struct {
	JWTSecret         string
	JWTIssuer         string
	JWTAudience       string
	ValidationEnabled bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Services: ServicesConfig{
			SuperAgent:          getEnv("SUPERAGENT_SERVICE_BASE_URL", "http://localhost:8888"),
			ConversationService: getEnv("CONVERSATION_SERVICE_BASE_URL", "http://localhost:8888"),
			ProfileService:      getEnv("PROFILE_SERVICE_BASE_URL", "http://localhost:8090"),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("AUTH_JWT_SECRET", ""),
			JWTIssuer:         getEnv("AUTH_JWT_ISSUER", "https://ssopre.comune.roma.it:443/ssoservice/oauth2/realms/root/realms/public"),
			JWTAudience:       getEnv("AUTH_JWT_AUDIENCE", "julia"),
			ValidationEnabled: getEnvBool("AUTH_VALIDATION_ENABLED", true),
		},
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolVal
	}
	return defaultValue
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Services.SuperAgent == "" {
		return fmt.Errorf("SUPERAGENT_SERVICE_BASE_URL is required")
	}
	if c.Services.ConversationService == "" {
		return fmt.Errorf("CONVERSATION_SERVICE_BASE_URL is required")
	}
	if c.Auth.ValidationEnabled && c.Auth.JWTSecret == "" && c.Environment == "production" {
		return fmt.Errorf("AUTH_JWT_SECRET is required in production when validation is enabled")
	}
	return nil
}
