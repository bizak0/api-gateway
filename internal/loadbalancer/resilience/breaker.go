package resilience

import (
	"errors"
	"log"
	"sync"
	"time"
)

type State int

const (
	StateClosed   State = iota
	StateOpen     State = iota
	StateHalfOpen State = iota
)

type CircuitBreaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	maxFailures  int
	lastFailure  time.Time
	resetTimeout time.Duration
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        StateClosed,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			log.Println("Circuit Breaker: Half-Open")
			cb.state = StateHalfOpen
			return true
		}
		log.Println("Circuit Breaker: Open - request blocked")
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.state = StateClosed
	log.Println("Circuit Breaker: Closed")
}

func (cb *CircuitBreaker) Failure() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = StateOpen
		log.Println("Circuit Breaker: Open - too many failures")
	}

	return errors.New("circuit breaker: request failed")
}
