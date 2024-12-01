package users

import (
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
	api.HandleFunc(http.MethodGet, "/users/current", func(w http.ResponseWriter, r *http.Request, _ app.Route) {
		c := controller.Query[*entity.User](w, r)
		u, err := m.User().Find(r.Context(), model.DummyUserID)
		if err != nil {
			if model.IsErrorNotFound(err) {
				c.Response().NotFound(err)
				return
			}
			c.Response().InternalServerError(err)
			return
		}

		c.Render(u)
	})
}
