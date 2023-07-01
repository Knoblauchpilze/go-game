package routes

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/middlewares"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
)

func UsersRouter(repo users.Repository) http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.RequestCtx)
		r.Get("/", getUsers(repo))

		r.Route("/{user}", func(r chi.Router) {
			r.Post("/", createUser(repo))
			r.Get("/", getUser(repo))
			r.Patch("/", patchUser(repo))
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

		logger.Infof("Hello from getUsers")
		reqData.FailWithErrorAndCode(errors.NewCode(errors.ErrNotImplemented), http.StatusInternalServerError, w)
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

func patchUser(repo users.Repository) http.HandlerFunc {
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

		id, err := getUserIdFromHttpRequest(r)
		if err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		patch := dto.Convert()
		patch.Id = id

		user, err := repo.Patch(id, patch)
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
