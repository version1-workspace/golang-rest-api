package app

import (
	"fmt"
	"net/http"
)

type app struct {
	basePath string

	m           *http.ServeMux
	routes      []Route
	middlewares []Middleware
}

type Middleware func(w http.ResponseWriter, r *http.Request) bool

func New(basePath string) *app {
	return &app{
		basePath: basePath,
		m:        http.NewServeMux(),
	}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range a.middlewares {
		if !m(w, r) {
			return
		}
	}
	a.m.ServeHTTP(w, r)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request, route Route)

func (a *app) HandleFunc(method, path string, handler HandlerFunc) {
	fullPath := a.basePath + path
	matcher := fmt.Sprintf("%s %s", method, fullPath)
	route := Route{
		Method:  method,
		Matcher: fullPath,
	}
	a.m.HandleFunc(matcher, func(w http.ResponseWriter, r *http.Request) {
		route.Path = r.URL.Path
		handler(w, r, route)
	})
	a.routes = append(a.routes, route)
}

func (a *app) Walk(cb func(r Route)) {
	for _, route := range a.routes {
		cb(route)
	}
}

func (a *app) Use(middlewares ...Middleware) {
	a.middlewares = append(a.middlewares, middlewares...)
}
