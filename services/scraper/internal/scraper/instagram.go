package scraper

import (
	"context"
	"distributed-web-scrapper/services/scraper/internal/auth"
	"distributed-web-scrapper/services/scraper/internal/kafka"
	"time"

	"github.com/chromedp/chromedp"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type InstagramScraper struct {
	producer     *kafka.Producer
	logger       *zap.Logger
	oauthClient  *auth.OAuthClient
	proxyRotator *ProxyRotator
	rateLimiter  *RateLimiter
	cb           *CircuitBreaker
}

func NewInstagramScraper(producer *kafka.Producer, oauthClient *auth.OAuthClient, proxyRotator *ProxyRotator, rateLimiter *RateLimiter, cb *CircuitBreaker) Scraper {
	logger, _ := zap.NewProduction()
	return &InstagramScraper{
		producer:     producer,
		logger:       logger,
		oauthClient:  oauthClient,
		proxyRotator: proxyRotator,
		rateLimiter:  rateLimiter,
		cb:           cb,
	}
}

func (s *InstagramScraper) Start(ctx context.Context) {
	tracer := otel.Tracer("instagram-scraper")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("shutting down instagram scraper")
			return
		case <-ticker.C:
			_, span := tracer.Start(ctx, "scrape-instagram")
			if err := s.Scrape(ctx); err != nil {
				s.logger.Error("scrape failed", zap.Error(err))
			}
			span.End()
		}
	}
}

func (s *InstagramScraper) Scrape(ctx context.Context) error {
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return err
	}

	proxy := s.proxyRotator.GetProxy()
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(proxy),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var postCaption string
	_, err := s.cb.Execute(func() (interface{}, error) {
		return nil, chromedp.Run(browserCtx,
			chromedp.Navigate("https://www.instagram.com/p/some-post/"),
			chromedp.Text("div._a9zs", &postCaption),
		)
	})
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"platform": "instagram",
		"caption":  postCaption,
		"timestamp": time.Now(),
	}

	return s.producer.Publish("instagram_data", data)
}