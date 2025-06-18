package scraper

import "context"

type Scraper interface {
	Start(ctx context.Context)
	Scrape(ctx context.Context) error
}