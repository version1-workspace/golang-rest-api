package posts

import (
	"net/http"

	"github.com/version-1/golang-rest-api/internal/app"
	"github.com/version-1/golang-rest-api/internal/controller"
	"github.com/version-1/golang-rest-api/internal/model"
	"github.com/version-1/golang-rest-api/internal/model/entity"
)

func create(m *model.Model) app.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[entity.PostBody, *entity.Post](w, r)
		var body *entity.PostBody
		err := c.Validate(func(input *entity.PostBody, req *controller.Request[entity.PostBody]) error {
			b, err := c.Request().Body()
			if err != nil {
				return err
			}

			body = b

			return nil
		})

		if err != nil {
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
			}

			return nil
		})
		if err != nil {
			c.Response().InternalServerError(err)
			return
		}

		c.Render(post)
	}
}

func patch(m *model.Model) app.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[entity.PostBody, *entity.Post](w, r)
		var body *entity.PostBody
		var post *entity.Post
		var ok bool
		err := c.Validate(func(input *entity.PostBody, req *controller.Request[entity.PostBody]) error {
			b, err := c.Request().Body()
			if err != nil {
				return err
			}

			body = b

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

		ctx := r.Context()
		err = m.Transaction(ctx, func(tx model.Executor) error {
			_, err = m.Post(tx).Update(ctx, model.DummyUserID, body.Title, body.Content)
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

		post, ok = ensureExist(r.Context(), post.ID, m.Post(), c.Response())
		if !ok {
			return
		}

		c.Render(post)
	}
}

func destroy(m *model.Model) app.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ro app.Route) {
		c := controller.New[controller.VoidInput, *entity.Post](w, r)
		var post *entity.Post
		err := c.Validate(func(input *controller.VoidInput, req *controller.Request[controller.VoidInput]) error {
			id, err := retrieveID(ro)
			if err != nil {
				return err
			}

			post, _ = ensureExist(r.Context(), id, m.Post(), c.Response())
			return nil
		})

		if err != nil {
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
	}
}
