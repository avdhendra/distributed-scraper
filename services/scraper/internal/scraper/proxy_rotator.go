package scraper

import (
	"distributed-web-scrapper/services/scraper/internal/config"
	"sync"
)

type ProxyRotator struct {
	proxies []string
	index   int
	mu      sync.Mutex
}

func NewProxyRotator() *ProxyRotator {
	cfg, _ := config.LoadFromConsul()
	return &ProxyRotator{
		proxies: cfg.ProxyList,
	}
}

func (pr *ProxyRotator) GetProxy() string {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	proxy := pr.proxies[pr.index]
	pr.index = (pr.index + 1) % len(pr.proxies)
	return proxy
}