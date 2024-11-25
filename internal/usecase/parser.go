package usecase

import (
	"context"
	"log"
	"strings"

	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/htmlparser"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/httpclient"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/rabbitmq"
	"github.com/ZetoOfficial/domain-scraper/pkg/linkutils"
)

type ParserUseCase struct {
	httpClient *httpclient.HTTPClient
	htmlParser *htmlparser.HTMLParser
	rabbitMQ   *rabbitmq.Producer
	linkUtils  *linkutils.LinkUtils
}

func NewParserUseCase(client *httpclient.HTTPClient, parser *htmlparser.HTMLParser, producer *rabbitmq.Producer, utils *linkutils.LinkUtils) *ParserUseCase {
	return &ParserUseCase{
		httpClient: client,
		htmlParser: parser,
		rabbitMQ:   producer,
		linkUtils:  utils,
	}
}

func (p *ParserUseCase) ParseAndSend(ctx context.Context, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	body, title, err := p.httpClient.Get(url)
	if err != nil {
		return err
	}
	log.Printf("Page Title: %s", strings.TrimSpace(strings.ReplaceAll(title, "\n", "")))

	links, err := p.htmlParser.ParseLinks(body, url)
	if err != nil {
		return err
	}

	for _, l := range links {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Printf("Found link: %s", l.Href)
			err = p.rabbitMQ.Publish(l.Href)
			if err != nil {
				log.Printf("Failed to publish link '%s': %v", l.Href, err)
			}
		}
	}

	return nil
}
