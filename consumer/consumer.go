package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

// RabbitMQConsumer is a struct that represents the consumer with its dependencies.
type RabbitMQConsumer struct {
	conn                  *amqp.Connection
	channel               *amqp.Channel
	queueName             string
	workers               int
	omsIntegrationService *OmsIntegrationService
}

// NewRabbitMQConsumer creates a new RabbitMQConsumer with the provided connection and queue name.
func NewRabbitMQConsumer(conn *amqp.Connection,
	queueName string,
	workers int,
	omsIntegrationService *OmsIntegrationService) (*RabbitMQConsumer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:                  conn,
		channel:               channel,
		queueName:             queueName,
		workers:               workers,
		omsIntegrationService: omsIntegrationService,
	}, nil
}

// ConsumeMessages starts consuming messages from the queue and handles them.
func (c *RabbitMQConsumer) ConsumeMessages() {
	// consume messages
	msgs, err := c.channel.Consume(
		c.queueName,
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

	// create signal for closing rabbitmq connection if signal received
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a worker pool with the specified number of workers
	for i := 0; i < c.workers; i++ {
		go func(workerID int) {
			fmt.Printf("Worker %d is waiting for messages. To exit, press Ctrl+C\n", workerID)
			for msg := range msgs {
				log.Printf("Worker %d received a message: %s", workerID, msg.Body)

				_, err := c.handleMessage(msg.Body, workerID)
				if err != nil {
					log.Printf("Error handling message: %v", err)
					return
				}

				if err := msg.Ack(false); err != nil {
					log.Printf("Worker %d error acknowledging message: %v", workerID, err)
				}
			}
		}(i)
	}

	sig := <-sigs
	log.Printf("Received %v signal. Closing connection...", sig)
	c.conn.Close()
}

// handleMessage handles each message and process them.
func (c *RabbitMQConsumer) handleMessage(messageBody []byte, workerId int) (bool, error) {
	// mapping message
	var orderQueueMessage NewOrderQueueMessage
	err := json.Unmarshal(messageBody, &orderQueueMessage)
	if err != nil {
		log.Printf("Fail to unmarshal %v %s", err.Error(), messageBody)
		return false, nil
	}

	// process message
	success, err := c.omsIntegrationService.OmsIntegration(orderQueueMessage)
	log.Printf("success: %v, workerId: %v", success, workerId)
	if err != nil {
		log.Printf("OmsIntegration fails %v", err)
		return false, err
	}
	return success, nil
}
