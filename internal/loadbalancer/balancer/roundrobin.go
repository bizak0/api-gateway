package balancer

import (
	"sync"

	"github.com/bizak0/api-gateway/internal/loadbalancer/registry"
)

type RoundRobin struct {
	mu      sync.Mutex
	current int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		current: 0,
	}
}

func (rr *RoundRobin) Next(services []*registry.Service) *registry.Service {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(services) == 0 {
		return nil
	}

	service := services[rr.current%len(services)]
	rr.current++
	return service
}
