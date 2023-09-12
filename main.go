package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mahfuzan/swiftoms-rabbitmq/client/omsclient"
	"github.com/mahfuzan/swiftoms-rabbitmq/config"
	"github.com/mahfuzan/swiftoms-rabbitmq/consumer"
	"github.com/streadway/amqp"
)

func main() {
	conf := config.NewConfig()

	// get number of workers from cli argument
	var numWorkers int
	flag.IntVar(&numWorkers, "w", 1, "number of workers")
	flag.Parse()

	// create connection to rabbitmq
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.AmqpUser, conf.AmqpPassword, conf.AmqpHost, conf.AmqpPort)
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// setup oms integration service
	omsClient := setupOmsClient(conf.OmsHost, conf.OmsToken)
	omsIntegrationService := consumer.NewOmsService(omsClient)

	// create new consumer
	consumer, err := consumer.NewRabbitMQConsumer(conn, conf.OmsNewOrderQueue, numWorkers, omsIntegrationService)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}

	// consume message
	consumer.ConsumeMessages()
}

// setupOmsClient setup a new oms client.
func setupOmsClient(baseUrl, token string) omsclient.Client {
	return omsclient.NewClient(baseUrl, token)
}
