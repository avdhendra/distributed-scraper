package main

import (
	"context"
	"distributed-web-scrapper/services/metrics/internal/config"
	"distributed-web-scrapper/services/metrics/internal/logger"
	"distributed-web-scrapper/services/metrics/internal/metrics"
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

	_, err = config.LoadFromConsul()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	metrics.Init()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("received shutdown signal, initiating graceful shutdown")
	cancel()
	time.Sleep(5 * time.Second)
	log.Info("metrics service shutdown complete")
}