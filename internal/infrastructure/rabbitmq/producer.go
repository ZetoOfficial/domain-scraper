package rabbitmq

import (
	"fmt"

	"github.com/ZetoOfficial/domain-scraper/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	queueName  string
}

func NewProducer(cfg *config.Config) (*Producer, error) {
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

	return &Producer{
		connection: conn,
		channel:    ch,
		queueName:  cfg.RabbitMQQueue,
	}, nil
}

func (p *Producer) Publish(message string) error {
	err := p.channel.Publish(
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.connection != nil {
		p.connection.Close()
	}
}
