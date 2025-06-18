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

type LinkedInScraper struct {
	producer     *kafka.Producer
	logger       *zap.Logger
	oauthClient  *auth.OAuthClient
	proxyRotator *ProxyRotator
	rateLimiter  *RateLimiter
	cb           *CircuitBreaker
}

func NewLinkedInScraper(producer *kafka.Producer, oauthClient *auth.OAuthClient, proxyRotator *ProxyRotator, rateLimiter *RateLimiter, cb *CircuitBreaker) Scraper {
	logger, _ := zap.NewProduction()
	return &LinkedInScraper{
		producer:     producer,
		logger:       logger,
		oauthClient:  oauthClient,
		proxyRotator: proxyRotator,
		rateLimiter:  rateLimiter,
		cb:           cb,
	}
}

func (s *LinkedInScraper) Start(ctx context.Context) {
	tracer := otel.Tracer("linkedin-scraper")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("shutting down linkedin scraper")
			return
		case <-ticker.C:
			_, span := tracer.Start(ctx, "scrape-linkedin")
			if err := s.Scrape(ctx); err != nil {
				s.logger.Error("scrape failed", zap.Error(err))
			}
			span.End()
		}
	}
}

func (s *LinkedInScraper) Scrape(ctx context.Context) error {
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

	var profileData string
	_, err := s.cb.Execute(func() (interface{}, error) {
		return nil, chromedp.Run(browserCtx,
			chromedp.Navigate("https://www.linkedin.com/in/some-profile"),
			chromedp.Text(".pv-top-card", &profileData),
		)
	})
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"platform": "linkedin",
		"profile":  profileData,
		"timestamp": time.Now(),
	}

	return s.producer.Publish("linkedin_data", data)
}