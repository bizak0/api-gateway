package health

import (
	"log"
	"net/http"
	"time"

	"github.com/bizak0/api-gateway/internal/loadbalancer/registry"
)

type HealthChecker struct {
	registry *registry.Registry
	interval time.Duration
}

func NewHealthChecker(reg *registry.Registry, interval time.Duration) *HealthChecker {
	return &HealthChecker{
		registry: reg,
		interval: interval,
	}
}

func (hc *HealthChecker) Start() {
	go func() {
		for {
			hc.checkAll()
			time.Sleep(hc.interval)
		}
	}()
}

func (hc *HealthChecker) checkAll() {
	services := hc.registry.GetHealthy()
	for _, service := range services {
		healthy := hc.check(service.Address)
		hc.registry.SetHealth(service.ID, healthy)
		if !healthy {
			log.Printf("Service %s is DOWN at %s", service.ID, service.Address)
		} else {
			log.Printf("Service %s is UP at %s", service.ID, service.Address)
		}
	}
}

func (hc *HealthChecker) check(address string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(address + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
