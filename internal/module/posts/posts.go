package posts

import (
	"context"
	"net/http"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/controller"
	"github.com/version-1/golang-rest-api/internal/model"
	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type handlerRegister interface {
	HandleFunc(string, string, app.HandlerFunc)
}

func Register(api handlerRegister, m *model.Model) {
	api.HandleFunc(http.MethodGet, "/posts", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[controller.VoidInput, []*entity.Post](w, r)
		posts, err := m.Post().FindAll(r.Context(), 10)
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		c.Render(posts)
	})

	api.HandleFunc(http.MethodPost, "/posts", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.New[entity.PostBody, *entity.Post](w, r)
		body, err := c.Request().Body()
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		var post *entity.Post
		ctx := r.Context()
		err = m.Transaction(ctx, func(tx model.Executor) error {
			post, err = m.Post(tx).Create(ctx, model.DummyUserID, body.Title, body.Content)
			if err != nil {
				return err
			}

			for _, tag := range body.Tags {
				_, err = m.Tag(tx).Attach(ctx, post.ID, tag.Slug, tag.Name)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		post, ok := ensureExist(ctx, post.ID, m.Post(), c.Response())
		if !ok {
			return
		}

		c.Render(post)
	})

	api.HandleFunc(http.MethodPatch, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[entity.PostBody, *entity.Post](w, r)
		body, err := c.Request().Body()
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		id, err := retrieveID(ro)
		if err != nil {
			c.Response().BadRequest(err)
			return
		}

		if _, ok := ensureExist(r.Context(), id, m.Post(), c.Response()); !ok {
			return
		}

		ctx := r.Context()
		var post *entity.Post
		err = m.Transaction(ctx, func(tx model.Executor) error {
			post, err = m.Post(tx).Update(ctx, model.DummyUserID, body.Title, body.Content)
			if err != nil {
				return err
			}

			err = m.Tag(tx).DetachAll(ctx, post.ID)
			if err != nil {
				return err
			}

			for _, tag := range body.Tags {
				_, err = m.Tag(tx).Attach(ctx, post.ID, tag.Slug, tag.Name)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		c.Render(post)
	})

	api.HandleFunc(http.MethodGet, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[controller.VoidInput, *entity.Post](w, r)
		var post *entity.Post
		var ok bool
		id, err := retrieveID(ro)
		if err != nil {
			c.Response().BadRequest(err)
			return
		}

		post, ok = ensureExist(r.Context(), id, m.Post(), c.Response())
		if !ok {
			return
		}

		c.Render(post)
	})

	api.HandleFunc(http.MethodDelete, "/posts/{id}", func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[controller.VoidInput, *entity.Post](w, r)
		var post *entity.Post
		id, err := retrieveID(ro)
		if err != nil {
			c.Response().BadRequest(err)
			return
		}

		post, ok := ensureExist(r.Context(), id, m.Post(), c.Response())
		if !ok {
			return
		}

		ctx := r.Context()
		err = m.Transaction(ctx, func(tx model.Executor) error {
			if _, err = m.Post(tx).Delete(ctx, post.ID); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		c.Render(post)
	})

}

func retrieveID(ro app.Route) (int, error) {
	return ro.Params().Int("id")
}

type errorRenderer interface {
	InternalServerError(error)
	NotFound(error)
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
