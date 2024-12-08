package controller

import (
	"encoding/json"
	"net/http"
)

type Response[I, O any] struct {
	status      int
	contentType string
	w           http.ResponseWriter
	req         *Request[I]
}

func NewResponse[I, O any](w http.ResponseWriter, req *Request[I]) *Response[I, O] {
	contentType := "application/json"
	w.Header().Set("Content-Type", contentType)

	return &Response[I, O]{
		w:   w,
		req: req,
	}
}

func (r *Response[I, O]) Status(status int) *Response[I, O] {
	r.status = status
	r.w.WriteHeader(status)
	return r
}

func (r *Response[I, O]) Render(payload any) {
	if err := json.NewEncoder(r.w).Encode(payload); err != nil {
		r.Status(http.StatusInternalServerError).Render(map[string]any{"error": err.Error()})
	}
}

var nilResponse = map[string]any{}

func (r *Response[I, O]) renderEmpty(w http.ResponseWriter) {
	w.WriteHeader(r.status)
	if err := json.NewEncoder(w).Encode(nilResponse); err != nil {
		panic(err)
	}
}

func (r *Response[I, O]) renderError(err error) *Response[I, O] {
	r.Status(r.status).Render(map[string]any{"error": err.Error()})
	return r
}

func (c *Response[I, O]) Created() {
	c.Status(http.StatusCreated).renderEmpty(c.w)
}

func (c *Response[I, O]) NotFound(payload error) {
	c.Status(http.StatusNotFound).renderError(payload)
}

func (c *Response[I, O]) InternalServerError(payload error) {
	c.Status(http.StatusInternalServerError).renderError(payload)
}

func (c *Response[I, O]) Unauthorized(payload error) {
	c.Status(http.StatusUnauthorized).renderError(payload)
}

func (c *Response[I, O]) UnprocessableEntity(payload error) {
	c.Status(http.StatusUnprocessableEntity).renderError(payload)
}

func (c *Response[I, O]) BadRequest(payload error) *Response[I, O] {
	return c.Status(http.StatusBadRequest).renderError(payload)
}
