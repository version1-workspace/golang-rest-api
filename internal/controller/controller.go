package controller

import "net/http"

type Controller[I, O any] struct {
	w   http.ResponseWriter
	req *Request[I]
	res *Response[I, O]
}

func New[I, O any](w http.ResponseWriter, req *http.Request) *Controller[I, O] {
	r := &Request[I]{
		r: req,
	}
	return &Controller[I, O]{
		req: r,
		res: NewResponse[I, O](w, r),
	}
}

func (c *Controller[I, O]) Request() *Request[I] {
	return c.req
}

func (c *Controller[I, O]) Response() *Response[I, O] {
	return c.res
}

func (c *Controller[I, O]) Render(payload O) {
	c.res.Render(payload)
}
