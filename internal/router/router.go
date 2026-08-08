package router

import (
	"net/http"
	"strings"

	"github.com/aminebenselem/api-gateway/internal/config"
)

type Router struct {
	routes *config.RoutesConfig
}

func New(routes *config.RoutesConfig) *Router {
	return &Router{
		routes: routes,
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
	route, ok := r.Match(req.URL.Path)

	if !ok {
		http.NotFound(w, req)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		"matched route: " + route.Path +
			"\ntargets: " + strings.Join(route.Targets, ", "),
	))
}