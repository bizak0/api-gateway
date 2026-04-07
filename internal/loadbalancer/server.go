package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/bizak0/api-gateway/internal/loadbalancer/balancer"
	"github.com/bizak0/api-gateway/internal/loadbalancer/health"
	"github.com/bizak0/api-gateway/internal/loadbalancer/registry"
	"github.com/bizak0/api-gateway/internal/loadbalancer/resilience"
	"github.com/bizak0/api-gateway/internal/loadbalancer/router"
)

type LoadBalancer struct {
	registry *registry.Registry
	balancer *balancer.RoundRobin
	router   *router.Router
	breaker  *resilience.CircuitBreaker
	retry    *resilience.RetryConfig
}

func NewLoadBalancer() *LoadBalancer {
	reg := registry.NewRegistry()
	checker := health.NewHealthChecker(reg, 10*time.Second)
	checker.Start()

	return &LoadBalancer{
		registry: reg,
		balancer: balancer.NewRoundRobin(),
		router:   router.NewRouter(),
		breaker:  resilience.NewCircuitBreaker(3, 30*time.Second),
		retry:    resilience.NewRetryConfig(3, 1*time.Second),
	}
}

func (lb *LoadBalancer) Register(id string, address string) {
	lb.registry.Register(id, address)
	log.Printf("Registered service %s at %s", id, address)
}

func (lb *LoadBalancer) AddRoute(prefix string, service string) {
	lb.router.AddRoute(prefix, service)
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !lb.breaker.Allow() {
		http.Error(w, "Service unavailable - circuit breaker open", http.StatusServiceUnavailable)
		return
	}

	services := lb.registry.GetHealthy()
	if len(services) == 0 {
		http.Error(w, "No healthy services available", http.StatusServiceUnavailable)
		return
	}

	service := lb.balancer.Next(services)
	if service == nil {
		http.Error(w, "No service available", http.StatusServiceUnavailable)
		return
	}

	target, err := url.Parse(service.Address)
	if err != nil {
		lb.breaker.Failure()
		http.Error(w, "Invalid service address", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Routing request %s to %s\n", r.URL.Path, service.Address)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		lb.breaker.Failure()
		http.Error(w, "Service error", http.StatusBadGateway)
	}

	lb.breaker.Success()
	proxy.ServeHTTP(w, r)
}
