package loadbalancer

import (
	"net/http"
	"sync"
	"time"
)

type Backend struct {
	URL     string
	Healthy bool
}

type LoadBalancer struct {
	backends []Backend
	current  int
	mu       sync.Mutex
}

func New(targets []string) *LoadBalancer {
	backends := make([]Backend, len(targets))

	for i, target := range targets {
		backends[i] = Backend{
			URL:     target,
			Healthy: true,
		}
	}

	return &LoadBalancer{
		backends: backends,
	}
}

func (lb *LoadBalancer) Next() (*Backend, bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.backends) == 0 {
		return nil, false
	}

	for i := 0; i < len(lb.backends); i++ {
		index := (lb.current + i) % len(lb.backends)

		if lb.backends[index].Healthy {
			lb.current = (index + 1) % len(lb.backends)

			return &lb.backends[index], true
		}
	}

	return nil, false
}

func (lb *LoadBalancer) MarkUnhealthy(backend *Backend) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := range lb.backends {
		if &lb.backends[i] == backend {
			lb.backends[i].Healthy = false
			return
		}
	}
}

func (lb *LoadBalancer) CheckHealth() {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	for i := range lb.backends {
		healthy := true

		resp, err := client.Get(lb.backends[i].URL + "/health")

		if err != nil {
			healthy = false
		} else {
			resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				healthy = false
			}
		}

		lb.mu.Lock()
		lb.backends[i].Healthy = healthy
		lb.mu.Unlock()
	}
}

func (lb *LoadBalancer) StartHealthChecks(interval time.Duration) {
	go func() {
		for {
			lb.CheckHealth()
			time.Sleep(interval)
		}
	}()
}