package registry

import (
	"sync"
)

type Service struct {
	ID      string
	Address string
	Healthy bool
}

type Registry struct {
	mu       sync.Mutex
	services map[string]*Service
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*Service),
	}
}

func (r *Registry) Register(id string, address string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.services[id] = &Service{
		ID:      id,
		Address: address,
		Healthy: true,
	}
}

func (r *Registry) Deregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.services, id)
}

func (r *Registry) GetHealthy() []*Service {
	r.mu.Lock()
	defer r.mu.Unlock()

	var healthy []*Service
	for _, s := range r.services {
		if s.Healthy {
			healthy = append(healthy, s)
		}
	}
	return healthy
}

func (r *Registry) SetHealth(id string, healthy bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s, exists := r.services[id]; exists {
		s.Healthy = healthy
	}
}
