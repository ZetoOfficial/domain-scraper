package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	timeoutArg := flag.String("timeout", "30", "Timeout in seconds for the consumer to process messages")
	flag.Parse()

	timeoutSeconds, err := strconv.Atoi(*timeoutArg)
	if err != nil || timeoutSeconds <= 0 {
		log.Fatalf("Invalid timeout value: %s. Please provide a positive integer.", *timeoutArg)
	}
	timeout := time.Duration(timeoutSeconds) * time.Second
	log.Printf("Using timeout: %s", timeout)

	cfg, err := config.LoadConfig("configs/config.env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	httpClient := httpclient.NewHTTPClient()
	htmlParser := htmlparser.NewHTMLParser()
	rabbitConsumer, err := rabbitmq.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer rabbitConsumer.Close()

	consumerUseCase := usecase.NewConsumerUseCase(httpClient, htmlParser, rabbitConsumer, linkutils.NewLinkUtils())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := consumerUseCase.Start(ctx); err != nil {
			log.Printf("Error during consuming: %v", err)
			cancel()
		}
	}()

	<-sigChan
	log.Println("Shutdown signal received")
	cancel()
	log.Println("Gracefully shutting down")
}
