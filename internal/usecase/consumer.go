package usecase

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"

	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/htmlparser"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/httpclient"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/rabbitmq"
	"github.com/ZetoOfficial/domain-scraper/pkg/linkutils"
)

type ConsumerUseCase struct {
	httpClient *httpclient.HTTPClient
	htmlParser *htmlparser.HTMLParser
	rabbitMQ   *rabbitmq.Consumer
	linkUtils  *linkutils.LinkUtils
}

func NewConsumerUseCase(client *httpclient.HTTPClient, parser *htmlparser.HTMLParser, consumer *rabbitmq.Consumer, utils *linkutils.LinkUtils) *ConsumerUseCase {
	return &ConsumerUseCase{
		httpClient: client,
		htmlParser: parser,
		rabbitMQ:   consumer,
		linkUtils:  utils,
	}
}

func (c *ConsumerUseCase) Start(ctx context.Context) error {
	msgs, err := c.rabbitMQ.Consume()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Timeout reached or context canceled. Exiting consumer.")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				log.Println("No messages received. Exiting consumer.")
				return nil
			}
			link := string(msg.Body)
			log.Printf("Processing link: %s", link)
			body, title, err := c.httpClient.Get(link)
			if err != nil {
				log.Printf("Failed to get URL '%s': %v", link, err)
				c.rabbitMQ.Acknowledge(msg)
				continue
			}
			log.Printf("Page Title: %s", title)

			links, err := c.htmlParser.ParseLinks(body, link)
			if err != nil {
				log.Printf("Failed to parse HTML for '%s': %v", link, err)
				c.rabbitMQ.Acknowledge(msg)
				continue
			}

			for _, l := range links {
				log.Printf("Found link: %s", l.Href)
				err = c.rabbitMQ.Channel.Publish(
					"",                   // exchange
					c.rabbitMQ.QueueName, // routing key
					false,
					false,
					amqp091.Publishing{
						ContentType: "text/plain",
						Body:        []byte(l.Href),
					},
				)
				if err != nil {
					log.Printf("Failed to publish link '%s': %v", l.Href, err)
				}
			}

			c.rabbitMQ.Acknowledge(msg)
		}
	}
}
