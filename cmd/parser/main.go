package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZetoOfficial/domain-scraper/internal/config"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/htmlparser"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/httpclient"
	"github.com/ZetoOfficial/domain-scraper/internal/infrastructure/rabbitmq"
	"github.com/ZetoOfficial/domain-scraper/internal/usecase"
	"github.com/ZetoOfficial/domain-scraper/pkg/linkutils"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <URL>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	url := flag.Arg(0)

	cfg, err := config.LoadConfig("configs/config.env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	httpClient := httpclient.NewHTTPClient()
	htmlParser := htmlparser.NewHTMLParser()
	rabbitProducer, err := rabbitmq.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer rabbitProducer.Close()

	parserUseCase := usecase.NewParserUseCase(httpClient, htmlParser, rabbitProducer, linkutils.NewLinkUtils())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func(url string) {
		if err := parserUseCase.ParseAndSend(ctx, url); err != nil {
			log.Printf("Error during parsing: %v", err)
			cancel()
		}
	}(url)

	<-sigChan
	log.Println("Shutdown signal received")
	cancel()
	log.Println("Gracefully shutting down")
}
