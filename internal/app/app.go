package app

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type Middleware func(http.ResponseWriter, *http.Request)

type app struct {
	basePath    string
	m           *http.ServeMux
	middlewares []Middleware
	routes      []Route
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range a.middlewares {
		middleware(w, r)
	}
	a.m.ServeHTTP(w, r)
}

func (a *app) Use(middleware Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}

func (a *app) Walk(cb func(r Route)) {
	for _, route := range a.routes {
		cb(route)
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request, Route)

func (a *app) HandleFunc(method, path string, handler HandlerFunc) {
	fullPath := filepath.Clean(fmt.Sprintf("%s%s", a.basePath, path))
	matcher := fmt.Sprintf("%s %s", method, fullPath)
	a.m.HandleFunc(matcher, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, Route{Method: method, Matcher: fullPath, Path: r.URL.Path})
	})
	a.routes = append(a.routes, Route{Method: method, Path: fullPath})
}

func New() *app {
	return &app{
		m:        http.NewServeMux(),
		basePath: "/api/v1",
	}
}
