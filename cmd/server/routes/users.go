package routes

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/middlewares"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
)

func UsersRouter(repo users.Repository) http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.RequestCtx)
		r.Use(middlewares.TimingCtx)
		r.Get("/", getUsers(repo))
		r.Post("/", createUser(repo))

		r.Route("/{user}", func(r chi.Router) {
			r.Get("/", getUser(repo))
			r.Delete("/", deleteUser(repo))
		})
	})

	return r
}

func getUsers(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, ok := middlewares.GetRequestDataFromContextOrFail(w, r)
		if !ok {
			return
		}

		users, err := repo.GetAll()
		if err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusInternalServerError, w)
			return
		}

		reqData.WriteDetails(users, w)
	}
}

func createUser(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, ok := middlewares.GetRequestDataFromContextOrFail(w, r)
		if !ok {
			return
		}

		var err error
		var dto dtos.UserDto
		if dto, err = getUserDtoFromRequest(r); err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		id, err := repo.Create(dto.Convert())
		if err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		reqData.WriteDetails(id, w)
	}
}

func getUser(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, ok := middlewares.GetRequestDataFromContextOrFail(w, r)
		if !ok {
			return
		}

		id, err := getUserIdFromHttpRequest(r)
		if err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		user, err := repo.Get(id)
		if err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		reqData.WriteDetails(user, w)
	}
}

func deleteUser(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, ok := middlewares.GetRequestDataFromContextOrFail(w, r)
		if !ok {
			return
		}

		id, err := getUserIdFromHttpRequest(r)
		if err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		if err := repo.Delete(id); err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		reqData.WriteDetails(id, w)
	}
}
