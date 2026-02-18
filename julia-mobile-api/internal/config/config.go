package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server      ServerConfig
	CosmosDB    CosmosDBConfig
	AppConfig   AppConfigConfig
	Environment string
	LogLevel    string
	Defaults    DefaultConfig
}

// DefaultConfig holds default values for app configuration
type DefaultConfig struct {
	Features map[string]bool
	Config   map[string]interface{}
	Locale   map[string]string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
}

// CosmosDBConfig holds Cosmos DB configuration
type CosmosDBConfig struct {
	Endpoint  string
	Key       string
	Database  string
	Container string
	Emulator  bool
}

// AppConfigConfig holds Azure App Configuration settings
type AppConfigConfig struct {
	Endpoint      string
	ConnectionStr string
	LabelFilter   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		CosmosDB: CosmosDBConfig{
			Endpoint:  getEnv("COSMOS_DB_ENDPOINT", "https://localhost:8182"),
			Key:       getEnv("COSMOS_DB_KEY", "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="),
			Database:  getEnv("COSMOS_DB_DATABASE", "bff_julia_db"),
			Container: getEnv("COSMOS_DB_CONTAINER", "app_config"),
			Emulator:  getEnvBool("COSMOS_EMULATOR_ENABLED", true),
		},
		AppConfig: AppConfigConfig{
			Endpoint:      getEnv("AZURE_APPCONFIG_ENDPOINT", "http://localhost:8484"),
			ConnectionStr: getEnv("AZURE_APPCONFIG_CONNECTION_STRING", "Endpoint=http://localhost:8484;Id=local;Secret=c2VjcmV0"),
			LabelFilter:   getEnv("AZURE_APPCONFIG_LABEL", "local"),
		},
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Defaults: DefaultConfig{
			Features: map[string]bool{
				"newUI":         true,
				"darkMode":      true,
				"notifications": true,
			},
			Config: map[string]interface{}{},
			Locale: map[string]string{
				"default": "it-IT",
			},
		},
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
	if c.CosmosDB.Endpoint == "" {
		return fmt.Errorf("cosmos DB endpoint is required")
	}
	if c.CosmosDB.Database == "" {
		return fmt.Errorf("cosmos DB database is required")
	}
	return nil
}
