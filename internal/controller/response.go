package controller

import (
	"encoding/json"
	"net/http"
)

type Response[I, O any] struct {
	status      int
	contentType string
	w           http.ResponseWriter
	r           *http.Request
}

type NilResponse map[string]any

var nilResponse = NilResponse{"data": nil}

func NewResponse[I, O any](w http.ResponseWriter, r *http.Request) *Response[I, O] {
	contentType := "application/json"
	w.Header().Add("Content-Type", contentType)

	return &Response[I, O]{
		w:           w,
		r:           r,
		contentType: contentType,
		status:      http.StatusOK,
	}
}

// Ensure if Status is called after determining the content body
func (c *Response[I, O]) Status(status int) *Response[I, O] {
	c.status = status
	c.w.WriteHeader(status)
	return c
}

func (c *Response[I, O]) ContentType(t string) *Response[I, O] {
	c.contentType = t
	c.w.Header().Add("Content-Type", t)
	return c
}

func (c *Response[I, O]) Render(w http.ResponseWriter, payload O) *Response[I, O] {
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		c.renderError(w, err).Status(http.StatusInternalServerError)
		return c
	}

	return c
}

func (c *Response[I, O]) renderEmpty(w http.ResponseWriter) *Response[I, O] {
	if err := json.NewEncoder(w).Encode(nilResponse); err != nil {
		panic(err)
	}

	return c
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (c *Response[I, O]) renderError(w http.ResponseWriter, err error) *Response[I, O] {
	if err := json.NewEncoder(w).Encode(ErrorResponse{err.Error()}); err != nil {
		panic(err)
	}

	return c
}

func (c *Response[I, O]) Created() *Response[I, O] {
	return c.Status(http.StatusCreated).renderEmpty(c.w)
}

func (c *Response[I, O]) NotFound(payload error) *Response[I, O] {
	return c.Status(http.StatusNotFound).renderError(c.w, payload)
}

func (c *Response[I, O]) InternalServerError(payload error) *Response[I, O] {
	return c.Status(http.StatusInternalServerError).renderError(c.w, payload)
}

func (c *Response[I, O]) Unauthorized(payload error) *Response[I, O] {
	return c.Status(http.StatusUnauthorized).renderError(c.w, payload)
}

func (c *Response[I, O]) UnprocessableEntity(payload error) *Response[I, O] {
	return c.Status(http.StatusUnprocessableEntity).renderError(c.w, payload)
}

func (c *Response[I, O]) BadRequest(payload error) *Response[I, O] {
	return c.Status(http.StatusBadRequest).renderError(c.w, payload)
}
