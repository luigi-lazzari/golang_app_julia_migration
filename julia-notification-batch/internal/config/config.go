package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Batch JobConfig  `mapstructure:"julia_batch_jobs_notification"`
	Rest  RestConfig `mapstructure:"julia_rest_services"`
}

type JobConfig struct {
	Cron       string `mapstructure:"cron"`
	MaxRetries int    `mapstructure:"maxRetries"`
}

type RestConfig struct {
	Notification RestService `mapstructure:"notification"`
	External     RestService `mapstructure:"external"`
}

type RestService struct {
	BaseURL        string        `mapstructure:"baseUrl"`
	RequestTimeout time.Duration `mapstructure:"requestTimeout"`
	ConnectTimeout time.Duration `mapstructure:"connectTimeout"`
}

func LoadConfig() *Config {
	viper.SetConfigName("application") // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")           // look for config in the working directory
	viper.AutomaticEnv()               // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config file found. Using environment variables.")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config
}
