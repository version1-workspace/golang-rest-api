package middleware

import (
	"log"
	"net/http"

	"github.com/version-1/golang-rest-api/internal/app"
)

func Logging() app.Middleware {
	return func(w http.ResponseWriter, r *http.Request) bool {
		log.Printf("Request Recieved: %s %s\n", r.Method, r.URL.Path)
		return true
	}
}
