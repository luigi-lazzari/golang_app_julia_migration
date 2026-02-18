package resilience

import (
	"time"

	"github.com/sony/gobreaker"
)

// NewCircuitBreaker creates a configured circuit breaker
func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}
