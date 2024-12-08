package main

import (
	"log"
	"net/http"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/controller"
	"github.com/version-1/golang-rest-api/internal/middleware"
	"github.com/version-1/golang-rest-api/internal/swagger"
)

func main() {
	api := app.New("/api/v1")
	api.Use(middleware.Logging())
	api.HandleFunc(http.MethodGet, "/spec", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(swagger.IndexHTML(r.URL.Host)))
	})
	api.HandleFunc(http.MethodGet, "/spec/swagger.yml", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(swagger.SwaggerYAML()))
	})

	api.HandleFunc(http.MethodGet, "/posts", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		c.Render(map[string]string{"message": "Hello, World!"})
	})

	api.HandleFunc(http.MethodPost, "/posts", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		c.Render(map[string]string{})
	})

	api.HandleFunc(http.MethodPatch, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		c.Render(map[string]string{"message": "Hello, World!"})
	})

	api.HandleFunc(http.MethodGet, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		c.Render(map[string]string{"message": "Hello, World!"})
	})

	api.HandleFunc(http.MethodDelete, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		c.Render(map[string]string{"message": "Hello, World!"})
	})

	api.HandleFunc(http.MethodGet, "/user/current", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		c.Render(map[string]string{"message": "Hello, World!"})
	})
	api.Walk(func(r app.Route) {
		log.Printf("%s %s\n", r.Method, r.Matcher)
	})

	if err := http.ListenAndServe(":8080", api); err != nil {
		panic(err)
	}
}
