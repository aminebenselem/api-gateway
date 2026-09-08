package router

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"time"

	"github.com/aminebenselem/api-gateway/internal/config"
	"github.com/aminebenselem/api-gateway/internal/loadbalancer"
	"github.com/aminebenselem/api-gateway/internal/ratelimiter"
)

type Router struct {
	routes        *config.RoutesConfig
	loadBalancers map[string]*loadbalancer.LoadBalancer
	rateLimiter   *ratelimiter.RateLimiter
}

func New(routes *config.RoutesConfig) *Router {
	loadBalancers := make(map[string]*loadbalancer.LoadBalancer)

	for _, route := range routes.Routes {
		lb := loadbalancer.New(route.Targets)

		lb.StartHealthChecks(5 * time.Second)

		loadBalancers[route.Path] = lb
	}

	return &Router{
		routes:        routes,
		loadBalancers: loadBalancers,
		rateLimiter:   ratelimiter.New(5, 10),
	}
}
func (r *Router) Match(path string) (*config.Route, bool) {
	var matched *config.Route

	for i := range r.routes.Routes {
		route := &r.routes.Routes[i]

		if strings.HasPrefix(path, route.Path) {
			if matched == nil || len(route.Path) > len(matched.Path) {
				matched = route
			}
		}
	}

	if matched == nil {
		return nil, false
	}

	return matched, true
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.rateLimiter.Allow() {
		http.Error(
			w,
			"rate limit exceeded",
			http.StatusTooManyRequests,
		)
		return
	}

	route, ok := r.Match(req.URL.Path)

	if !ok {
		http.NotFound(w, req)
		return
	}

	lb := r.loadBalancers[route.Path]

	backend, ok := lb.Next()
	if !ok {
		http.Error(
			w,
			"no healthy backends available",
			http.StatusServiceUnavailable,
		)
		return
	}

	log.Printf(
		"routing %s %s → %s",
		req.Method,
		req.URL.Path,
		backend.URL,
	)

	target, err := url.Parse(backend.URL)
	if err != nil {
		http.Error(w, "invalid target", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf(
			"backend failure %s %s → %s: %v",
			req.Method,
			req.URL.Path,
			backend.URL,
			err,
		)

		lb.MarkUnhealthy(backend)

		http.Error(
			w,
			"backend unavailable",
			http.StatusBadGateway,
		)
	}

	proxy.ServeHTTP(w, req)
}