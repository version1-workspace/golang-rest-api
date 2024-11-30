package main

import (
	"log"
	"net/http"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/controller"
	"github.com/version-1/golang-rest-api/internal/swagger"
)

func main() {
	api := app.New()
	api.Use(loggingMiddleware)

	api.HandleFunc(http.MethodGet, "/spec/swagger.yml", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(swagger.SwaggerYAML()))
	})

	api.HandleFunc(http.MethodGet, "/spec", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(swagger.IndexHTML(r.URL.Host)))
	})

	api.HandleFunc(http.MethodGet, "/posts", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.Query[map[string]string](w, r)
		c.Render(map[string]string{"message": "Hello, World!"})
	})

	api.HandleFunc(http.MethodPost, "/posts", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		err := c.Validate(func(input *controller.VoidInput, req *controller.Request[controller.VoidInput]) error {
			_, err := c.Request().Body()
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return
		}

		c.Render(map[string]string{"message": "ok"})
	})

	api.HandleFunc(http.MethodGet, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.Query[map[string]string](w, r)
		err := c.Validate(func(input *controller.VoidInput, req *controller.Request[controller.VoidInput]) error {
			p, err := ro.Params()
			if err != nil {
				return err
			}

			_, err = p.Int("id")
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return
		}

		c.Render(map[string]string{"message": "ok"})
	})

	api.HandleFunc(http.MethodPatch, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.Query[map[string]string](w, r)
		err := c.Validate(func(input *controller.VoidInput, req *controller.Request[controller.VoidInput]) error {
			p, err := ro.Params()
			if err != nil {
				return err
			}

			_, err = p.Int("id")
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return
		}

		c.Render(map[string]string{"message": "ok"})
	})

	api.HandleFunc(http.MethodDelete, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[controller.VoidInput, map[string]string](w, r)
		err := c.Validate(func(input *controller.VoidInput, req *controller.Request[controller.VoidInput]) error {
			p, err := ro.Params()
			if err != nil {
				return err
			}

			_, err = p.Int("id")
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return
		}

		c.Render(map[string]string{"message": "ok"})
	})

	api.HandleFunc(http.MethodGet, "/users/current", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.Query[map[string]string](w, r)
		c.Render(map[string]string{"message": "ok"})
	})

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
