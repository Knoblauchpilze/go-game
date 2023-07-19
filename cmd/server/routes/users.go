package routes

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
)

func UsersRouter(repo users.Repository) http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
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
		users, err := repo.GetAll()
		if err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusInternalServerError, w)
			return
		}

		rest.WriteDetails(r.Context(), users, w)
	}
}

func createUser(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var dto dtos.UserDto
		if dto, err = getUserDtoFromRequest(r); err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusBadRequest, w)
			return
		}

		id, err := repo.Create(dto.Convert())
		if err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusBadRequest, w)
			return
		}

		rest.WriteDetails(r.Context(), id, w)
	}
}

func getUser(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUserIdFromHttpRequest(r)
		if err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusBadRequest, w)
			return
		}

		user, err := repo.Get(id)
		if err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusBadRequest, w)
			return
		}

		rest.WriteDetails(r.Context(), user, w)
	}
}

func deleteUser(repo users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUserIdFromHttpRequest(r)
		if err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusBadRequest, w)
			return
		}

		if err := repo.Delete(id); err != nil {
			rest.FailWithErrorAndCode(r.Context(), err, http.StatusBadRequest, w)
			return
		}

		rest.WriteDetails(r.Context(), id, w)
	}
}
