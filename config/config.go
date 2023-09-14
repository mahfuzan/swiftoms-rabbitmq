package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config represents the configuration for services.
type Config struct {
	ServiceConfigs []ServiceConfig `mapstructure:"service_configs"`
}

// ServiceConfig represents the configuration for a single service connection.
type ServiceConfig struct {
	ServiceName   string        `mapstructure:"service_name" default:"oms-a"`
	AmqpHost      string        `mapstructure:"amqp_host" default:"localhost"`
	AmqpPort      string        `mapstructure:"amqp_port" default:"5672"`
	AmqpUsername  string        `mapstructure:"amqp_username" default:"guest"`
	AmqpPassword  string        `mapstructure:"amqp_password" default:"guest"`
	AmqpVhost     string        `mapstructure:"amqp_vhost" default:"/"`
	QueueConfigs  []QueueConfig `mapstructure:"queue_configs"`
	SwiftomsHost  string        `mapstructure:"swiftoms_host" default:""`
	SwiftomsToken string        `mapstructure:"swiftoms_token" default:""`
}

type QueueConfig struct {
	QueueName       string `mapstructure:"queue_name" default:""`
	NumberOfWorkers int    `mapstructure:"number_of_workers" default:"1"`
}

func NewConfig() Config {
	// Load configuration from a config file or environment variables using Viper.
	viper.SetConfigName("config") // Name of your config file (config.yaml, config.json, etc.)
	viper.AddConfigPath(".")      // Path to look for the config file
	viper.SetEnvPrefix("MYAPP")   // Prefix for environment variables
	viper.AutomaticEnv()          // Automatically bind environment variables

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return config
}
