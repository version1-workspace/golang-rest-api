package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/middleware"
	"github.com/version-1/golang-rest-api/internal/model"
	"github.com/version-1/golang-rest-api/internal/module/posts"
	"github.com/version-1/golang-rest-api/internal/module/users"
	"github.com/version-1/golang-rest-api/internal/swagger"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=127.0.0.1 port=5432 user=gorest password=password dbname=gorest_development sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	m := model.New(conn, log.New(os.Stdout, "model: ", log.LstdFlags))

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

	posts.Register(api, m)
	users.Register(api, m)

	api.Walk(func(r app.Route) {
		log.Printf("%s %s\n", r.Method, r.Matcher)
	})

	if err := http.ListenAndServe(":8080", api); err != nil {
		panic(err)
	}
}
