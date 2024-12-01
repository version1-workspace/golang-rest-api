package posts

import (
	"context"
	"net/http"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/controller"
	"github.com/version-1/golang-rest-api/internal/model"
	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type registerHandler interface {
	HandleFunc(string, string, app.HandlerFunc)
}

func Register(api registerHandler, m *model.Model) {
	api.HandleFunc(http.MethodGet, "/posts", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.Query[[]*entity.Post](w, r)
		posts, err := m.Post().FindAll(r.Context(), 10)
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		c.Render(posts)
	})

	api.HandleFunc(http.MethodGet, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.Query[*entity.Post](w, r)
		var post *entity.Post
		var ok bool
		err := c.Validate(func(input *controller.VoidInput, req *controller.Request[controller.VoidInput]) error {
			id, err := retrieveID(ro)
			if err != nil {
				return err
			}

			post, ok = ensureExist(r.Context(), id, m.Post(), c.Response())
			return nil
		})

		if !ok || err != nil {
			return
		}

		c.Render(post)
	})
	api.HandleFunc(http.MethodPost, "/posts", create(m))
	api.HandleFunc(http.MethodPatch, "/posts/{id}", patch(m))
	api.HandleFunc(http.MethodDelete, "/posts/{id}", destroy(m))
}

type errorRenderer interface {
	NotFound(error)
	InternalServerError(error)
}

func retrieveID(ro app.Route) (int, error) {
	p, err := ro.Params()
	if err != nil {
		return 0, err
	}

	return p.Int("id")
}

func ensureExist(ctx context.Context, id int, pm *model.PostModel, er errorRenderer) (*entity.Post, bool) {
	post, err := pm.Find(ctx, id)
	if err != nil {
		if model.IsErrorNotFound(err) {
			er.NotFound(err)
			return post, false
		}
		er.InternalServerError(err)
		return post, false
	}

	return post, true
}
