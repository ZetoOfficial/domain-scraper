package rabbitmq

import (
	"fmt"

	"github.com/ZetoOfficial/domain-scraper/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp091.Connection
	Channel    *amqp091.Channel
	QueueName  string
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)
	conn, err := amqp091.Dial(connStr)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.RabbitMQQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &Consumer{
		connection: conn,
		Channel:    ch,
		QueueName:  cfg.RabbitMQQueue,
	}, nil
}

func (c *Consumer) Consume() (<-chan amqp091.Delivery, error) {
	msgs, err := c.Channel.Consume(
		c.QueueName,
		"",
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (c *Consumer) Acknowledge(delivery amqp091.Delivery) error {
	return delivery.Ack(false)
}

func (c *Consumer) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}
