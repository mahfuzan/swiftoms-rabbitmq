package consumer

import (
	"fmt"
	"log"

	"github.com/mahfuzan/swiftoms-rabbitmq/config"
	"github.com/streadway/amqp"
)

// RabbitMQConsumer is a struct that represents the consumer with its dependencies.
type RabbitMQConsumer struct {
	conn                  *amqp.Connection
	channel               *amqp.Channel
	conf                  config.ServiceConfig
	omsIntegrationService *OmsIntegrationService
}

// NewRabbitMQConsumer creates a new RabbitMQConsumer with the provided connection and queue name.
func NewRabbitMQConsumer(
	conf config.ServiceConfig,
	omsIntegrationService *OmsIntegrationService) (*RabbitMQConsumer, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", conf.AmqpUsername, conf.AmqpPassword, conf.AmqpHost, conf.AmqpPort, conf.AmqpVhost)
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:                  conn,
		channel:               channel,
		conf:                  conf,
		omsIntegrationService: omsIntegrationService,
	}, nil
}

// ConsumeMessages begins consuming messages.
func (c *RabbitMQConsumer) ConsumeMessages() {
	for _, queueConfig := range c.conf.QueueConfigs {
		go func(queueConf config.QueueConfig) {
			// consume messages
			messages, err := c.channel.Consume(
				queueConf.QueueName,
				"",
				false, // manual acknowledgment
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Fatalf("Failed to register a consumer: %v", err)
			}

			// Create a worker pool with the specified number of workers
			for i := 0; i < queueConf.NumberOfWorkers; i++ {
				go func(workerID int) {
					fmt.Printf("Service '%s': Worker %d for queue %s is waiting for messages. \n", c.conf.ServiceName, workerID, queueConf.QueueName)
					for message := range messages {
						// log.Printf("Worker %d for queue %s of service %s received a message: %s", workerID, queueConf.QueueName, c.conf.ServiceName, message.Body)
						log.Printf("Service '%s': Worker %d for queue %s received a message", c.conf.ServiceName, workerID, queueConf.QueueName)

						_, err := c.handleMessage(c.conf.ServiceName, queueConf.QueueName, message.Body, workerID)
						if err != nil {
							log.Printf("Error handling message for service '%s', worker %d, queue %s: %v", c.conf.ServiceName, workerID, queueConf.QueueName, err)
							return
						}

						if err := message.Ack(false); err != nil {
							log.Printf("Service '%s': Worker %d for queue %s error acknowledging message: %v", c.conf.ServiceName, workerID, queueConf.QueueName, err)
						}
					}
				}(i)
			}
		}(queueConfig)
	}
}

// handleMessage handles each message and process them.
func (c *RabbitMQConsumer) handleMessage(serviceName, queueName string, messageBody []byte, workerId int) (bool, error) {
	// process message
	success, err := c.omsIntegrationService.OmsIntegration(queueName, messageBody)
	if err != nil {
		log.Printf("OmsIntegration fails: %v", err)
		return false, err
	}
	log.Printf("Success: %v, Service '%s': workerId: %d, queue: %s", success, serviceName, workerId, queueName)
	return success, nil
}

// Close closes the RabbitMQ connection and channel.
func (c *RabbitMQConsumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
