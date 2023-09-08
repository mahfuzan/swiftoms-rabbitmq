package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AmqpHost     string `envconfig:"amqp_host" default:""`
	AmqpPort     string `envconfig:"amqp_port" default:""`
	AmqpUser     string `envconfig:"amqp_user" default:""`
	AmqpPassword string `envconfig:"amqp_user" default:""`
	AmqpVhost    string `envconfig:"amqp_user" default:"\\"`

	OmsNewOrderQueue string `envconfig:"oms_new_order_queue" default:"swiftoms.order-queue.new"`

	OmsHost  string `envconfig:"oms_token" default:""`
	OmsToken string `envconfig:"oms_token" default:""`
}

func NewConfig() Config {
	var c Config
	err := envconfig.Process("OMS_RABBITMQ", &c)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return c
}
