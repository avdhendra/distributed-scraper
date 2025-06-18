package main

import (
	"context"
	"distributed-web-scrapper/services/scraper/internal/auth"
	"distributed-web-scrapper/services/scraper/internal/config"
	"distributed-web-scrapper/services/scraper/internal/kafka"
	"distributed-web-scrapper/services/scraper/internal/logger"
	"distributed-web-scrapper/services/scraper/internal/scraper"
	"distributed-web-scrapper/services/scraper/internal/tracing"
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

	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatal("failed to initialize kafka producer", zap.Error(err))
	}
	defer producer.Close()

	oauthClient, err := auth.NewOAuthClient(cfg.OAuthConfig)
	if err != nil {
		log.Fatal("failed to initialize oauth client", zap.Error(err))
	}

	scraperFactory := scraper.NewFactory(producer, oauthClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scrapers := []string{"linkedin", "instagram", "youtube"}
	for _, platform := range scrapers {
		s, err := scraperFactory.CreateScraper(platform)
		if err != nil {
			log.Error("failed to create scraper", zap.String("platform", platform), zap.Error(err))
			continue
		}
		go s.Start(ctx)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("received shutdown signal, initiating graceful shutdown")
	cancel()
	time.Sleep(5 * time.Second)
	log.Info("scraper service shutdown complete")
}