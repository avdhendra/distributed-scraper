package scraper

import (
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(name string) *CircuitBreaker {
	return &CircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        name,
			MaxRequests: 5,
			Interval:    60 * time.Second,
			Timeout:     30 * time.Second,
		}),
	}
}

func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return cb.cb.Execute(fn)
}