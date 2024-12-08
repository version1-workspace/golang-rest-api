package controller

import "net/http"

type VoidInput struct{}

type Request[I any] struct {
	r    *http.Request
	body *I
}
