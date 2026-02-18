package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceBus      ServiceBusConfig      `mapstructure:"azure_servicebus"`
	NotificationHub NotificationHubConfig `mapstructure:"azure_notificationhub"`
}

type ServiceBusConfig struct {
	ConnectionString   string `mapstructure:"connectionString"`
	TopicName          string `mapstructure:"topicName"`
	SubscriptionName   string `mapstructure:"subscriptionName"`
	MaxConcurrentCalls int    `mapstructure:"maxConcurrentCalls"`
}

type NotificationHubConfig struct {
	ConnectionString   string `mapstructure:"connectionString"`
	HubName            string `mapstructure:"hubName"`
	Enabled            bool   `mapstructure:"enabled"`
	SendTimeoutSeconds int    `mapstructure:"sendTimeoutSeconds"`
}

func LoadConfig() *Config {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config file found. Using environment variables.")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if config.ServiceBus.MaxConcurrentCalls == 0 {
		config.ServiceBus.MaxConcurrentCalls = 1
	}

	if config.NotificationHub.SendTimeoutSeconds == 0 {
		config.NotificationHub.SendTimeoutSeconds = 60
	}

	return &config
}
