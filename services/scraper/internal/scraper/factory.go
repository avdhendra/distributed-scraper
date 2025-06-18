package scraper

import (
	"distributed-web-scrapper/services/scraper/internal/auth"
	"distributed-web-scrapper/services/scraper/internal/kafka"
	"fmt"
	"time"
)

type Factory struct {
	producer    *kafka.Producer
	oauthClient *auth.OAuthClient
	proxyRotator *ProxyRotator
	rateLimiter  *RateLimiter
}

func NewFactory(producer *kafka.Producer, oauthClient *auth.OAuthClient) *Factory {
	return &Factory{
		producer:    producer,
		oauthClient: oauthClient,
		proxyRotator: NewProxyRotator(),
		rateLimiter:  NewRateLimiter(10, time.Second),
	}
}

func (f *Factory) CreateScraper(platform string) (Scraper, error) {
	cb := NewCircuitBreaker(platform)
	switch platform {
	case "linkedin":
		return NewLinkedInScraper(f.producer, f.oauthClient, f.proxyRotator, f.rateLimiter, cb), nil
	case "instagram":
		return NewInstagramScraper(f.producer, f.oauthClient, f.proxyRotator, f.rateLimiter, cb), nil
	case "youtube":
		return NewYouTubeScraper(f.producer, f.oauthClient, f.proxyRotator, f.rateLimiter, cb), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}



