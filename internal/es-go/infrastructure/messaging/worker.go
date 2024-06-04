package messaging

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"os"
)

func NewRabbitMQConnection() (*amqp091.Connection, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_DEFAULT_USER"),
		os.Getenv("RABBITMQ_DEFAULT_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to message queue: %s", err)
	}

	return conn, nil
}
