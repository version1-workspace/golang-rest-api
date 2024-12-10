package controller

import (
	"encoding/json"
	"net/http"
)

type VoidInput struct{}

type Request[I any] struct {
	r    *http.Request
	body *I
}

func (r *Request[I]) Body() (*I, error) {
	if r.body != nil {
		return r.body, nil
	}

	if err := json.NewDecoder(r.r.Body).Decode(&r.body); err != nil {
		return r.body, err
	}

	defer r.r.Body.Close()

	return r.body, nil
}

func (r *Request[I]) Method() string {
	return r.r.Method
}
