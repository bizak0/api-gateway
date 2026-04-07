package router

import (
	"strings"
)

type Route struct {
	Prefix  string
	Service string
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{
		routes: []Route{},
	}
}

func (r *Router) AddRoute(prefix string, service string) {
	r.routes = append(r.routes, Route{
		Prefix:  prefix,
		Service: service,
	})
}

func (r *Router) Match(path string) string {
	for _, route := range r.routes {
		if strings.HasPrefix(path, route.Prefix) {
			return route.Service
		}
	}
	return ""
}
