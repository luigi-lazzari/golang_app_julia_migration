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
	Auth        AuthConfig
	Redis       RedisConfig
	Defaults    DefaultPreferences
}

// DefaultPreferences holds default values for user preferences
type DefaultPreferences struct {
	Chat          []PreferenceDefinition
	Notifications []string
}

// PreferenceDefinition defines a single preference
type PreferenceDefinition struct {
	ID       string
	Category string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
	TTL      int // in seconds
}

type AuthConfig struct {
	JWTSecret         string
	JWTIssuer         string
	JWTAudience       string
	ValidationEnabled bool
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
}

// CosmosDBConfig holds Cosmos DB configuration
type CosmosDBConfig struct {
	Endpoint                 string
	Key                      string
	Database                 string
	UserProfileContainer     string
	UserPreferencesContainer string
	Emulator                 bool
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
			Port: getEnv("SERVER_PORT", "8090"),
		},
		CosmosDB: CosmosDBConfig{
			Endpoint:                 getEnv("COSMOS_DB_ENDPOINT", ""),
			Key:                      getEnv("COSMOS_DB_KEY", ""),
			Database:                 getEnv("COSMOS_DB_DATABASE", ""),
			UserProfileContainer:     getEnv("COSMOS_DB_PROFILE_CONTAINER", "user_profiles"),
			UserPreferencesContainer: getEnv("COSMOS_DB_PREFERENCES_CONTAINER", "user_preferences"),
			Emulator:                 getEnvBool("COSMOS_EMULATOR_ENABLED", false),
		},
		AppConfig: AppConfigConfig{
			Endpoint:      getEnv("AZURE_APPCONFIG_ENDPOINT", ""),
			ConnectionStr: getEnv("AZURE_APPCONFIG_CONNECTION_STRING", ""),
			LabelFilter:   getEnv("AZURE_APPCONFIG_LABEL", "dev"),
		},
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Auth: AuthConfig{
			JWTSecret:         getEnv("AUTH_JWT_SECRET", ""),
			JWTIssuer:         getEnv("AUTH_JWT_ISSUER", "bff-julia"),
			JWTAudience:       getEnv("AUTH_JWT_AUDIENCE", "julia-app"),
			ValidationEnabled: getEnvBool("AUTH_VALIDATION_ENABLED", true),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			Enabled:  getEnvBool("REDIS_ENABLED", false),
			TTL:      getEnvInt("REDIS_TTL", 3600),
		},
		Defaults: DefaultPreferences{
			Chat: []PreferenceDefinition{
				{ID: "documents", Category: "SERVICES"},
				{ID: "school", Category: "SERVICES"},
				{ID: "taxes", Category: "SERVICES"},
				{ID: "housing", Category: "SERVICES"},
				{ID: "work", Category: "SERVICES"},
				{ID: "environment", Category: "SERVICES"},
				{ID: "business", Category: "SERVICES"},
				{ID: "art_culture", Category: "LEISURE"},
				{ID: "music_shows", Category: "LEISURE"},
				{ID: "food_dining", Category: "LEISURE"},
				{ID: "sports_nature", Category: "LEISURE"},
				{ID: "cinema", Category: "LEISURE"},
				{ID: "family", Category: "LEISURE"},
				{ID: "libraries", Category: "LEISURE"},
				{ID: "transit", Category: "TRANSPORT"},
				{ID: "car", Category: "TRANSPORT"},
				{ID: "bike", Category: "TRANSPORT"},
				{ID: "walking", Category: "TRANSPORT"},
				{ID: "motorcycle", Category: "TRANSPORT"},
			},
			Notifications: []string{
				"push_municipality",
				"push_mobility",
				"push_news",
				"push_appio",
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

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return intVal
	}
	return defaultValue
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.CosmosDB.Endpoint == "" {
		return fmt.Errorf("COSMOS_DB_ENDPOINT is required")
	}
	if c.CosmosDB.Key == "" && !c.CosmosDB.Emulator {
		return fmt.Errorf("COSMOS_DB_KEY is required when emulator is disabled")
	}
	if c.CosmosDB.Database == "" {
		return fmt.Errorf("COSMOS_DB_DATABASE is required")
	}
	if c.AppConfig.ConnectionStr == "" && c.Environment == "production" {
		return fmt.Errorf("AZURE_APPCONFIG_CONNECTION_STRING is required in production")
	}
	if c.Auth.ValidationEnabled && c.Auth.JWTSecret == "" && c.Environment == "production" {
		return fmt.Errorf("AUTH_JWT_SECRET is required in production when validation is enabled")
	}
	return nil
}
