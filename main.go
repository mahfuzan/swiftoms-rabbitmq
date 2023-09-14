package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mahfuzan/swiftoms-rabbitmq/client/omsclient"
	"github.com/mahfuzan/swiftoms-rabbitmq/config"
	"github.com/mahfuzan/swiftoms-rabbitmq/consumer"
)

func main() {
	configs := config.NewConfig()

	// Create a channel to listen for OS signals
	sigCh := make(chan os.Signal, 1)

	// Notify the sigCh channel for the specified signals
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, conf := range configs.ServiceConfigs {
			// Setup oms integration service
			omsClient := setupOmsClient(conf.SwiftomsHost, conf.SwiftomsToken)
			omsIntegrationService := consumer.NewOmsService(omsClient)

			// Create new consumer
			consumer, err := consumer.NewRabbitMQConsumer(conf, omsIntegrationService)
			if err != nil {
				log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
			}

			// Consume messages
			go consumer.ConsumeMessages()
		}
	}()

	// Wait for a signal to exit
	<-sigCh

	log.Printf("Received %v signal. Exiting...", sigCh)
}

// setupOmsClient setup a new oms client.
func setupOmsClient(baseUrl, token string) omsclient.Client {
	return omsclient.NewClient(baseUrl, token)
}
