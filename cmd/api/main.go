package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/model"
	"github.com/version-1/golang-rest-api/internal/module/posts"
	"github.com/version-1/golang-rest-api/internal/module/users"
	"github.com/version-1/golang-rest-api/internal/swagger"
)

func main() {
	api := app.New()
	api.Use(loggingMiddleware)
	// dsn := "psql://gorest:password@127.0.0.1:5432/gorest_development?sslmode=disable"
	connStr := "host=127.0.0.1 port=5432 user=gorest password=password dbname=gorest_development sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	m := model.New(conn)
	api.HandleFunc(http.MethodGet, "/spec/swagger.yml", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(swagger.SwaggerYAML()))
	})

	api.HandleFunc(http.MethodGet, "/spec", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(swagger.IndexHTML(r.URL.Host)))
	})

	posts.Register(api, m)
	users.Register(api, m)
	api.Walk(func(r app.Route) {
		log.Printf("%s %s", r.Method, r.Path)
	})

	if err := http.ListenAndServe(":8000", api); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request Received: %s %s", r.Method, r.URL.Path)
}
