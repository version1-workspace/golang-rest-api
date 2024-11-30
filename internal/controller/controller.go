package controller

import (
	"net/http"
)

type Controller[I, O any] struct {
	req *Request[I]
	res *Response[I, O]
}

type VoidInput *any

func New[I, O any](w http.ResponseWriter, r *http.Request) *Controller[I, O] {
	req := NewRequest[I](r)
	res := NewResponse[I, O](w, r)
	c := &Controller[I, O]{req: req, res: res}

	return c
}

func Query[O any](w http.ResponseWriter, r *http.Request) *Controller[VoidInput, O] {
	return New[VoidInput, O](w, r)
}

func (c *Controller[I, O]) Response() *Response[I, O] {
	return c.res
}

func (c *Controller[I, O]) Request() *Request[I] {
	return c.req
}

func (c *Controller[I, O]) Body() (*I, error) {
	return c.req.Body()
}

func (c *Controller[I, O]) Render(o O) *Response[I, O] {
	return c.res.Render(c.res.w, o)
}

func (c *Controller[I, O]) Validate(cb func(input *I, req *Request[I]) error) error {
	if isMutationMethod(c.req.Method()) {
		v, err := c.req.Body()
		if err != nil {
			c.res.BadRequest(err)
			return err
		}

		if err := cb(v, c.req); err != nil {
			c.res.BadRequest(err)
			return err
		}
		return nil
	}

	if err := cb(nil, c.req); err != nil {
		c.res.BadRequest(err)
		return err
	}

	return nil
}

func isMutationMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}
