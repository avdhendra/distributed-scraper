package main

import (
	"context"
	"distributed-web-scrapper/services/consumer/internal/config"
	"distributed-web-scrapper/services/consumer/internal/kafka"
	"distributed-web-scrapper/services/consumer/internal/logger"
	"distributed-web-scrapper/services/consumer/internal/storage"
	"distributed-web-scrapper/services/consumer/internal/tracing"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	tracer, err := tracing.InitTracer()
	if err != nil {
		log.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer tracer.Close()

	cfg, err := config.LoadFromConsul()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	store, err := storage.NewPostgresStorage(cfg.PostgresURL)
	if err != nil {
		log.Fatal("failed to initialize storage", zap.Error(err))
	}
	defer store.Close()

	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, "scraper-group")
	if err != nil {
		log.Fatal("failed to initialize kafka consumer", zap.Error(err))
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Consume(ctx, []string{"linkedin_data", "youtube_data", "instagram_data"}, store)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("received shutdown signal, initiating graceful shutdown")
	cancel()
	time.Sleep(5 * time.Second)
	log.Info("consumer service shutdown complete")
}