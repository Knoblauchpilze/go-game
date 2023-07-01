package routes

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var userIdDataKey = "user"

func getUserDtoFromRequest(r *http.Request) (dtos.UserDto, error) {
	var ud dtos.UserDto
	err := rest.GetBodyFromHttpRequestAs(r, &ud)
	return ud, err
}

func getUserIdFromHttpRequest(r *http.Request) (uuid.UUID, error) {
	var err error
	var id uuid.UUID

	qp := chi.URLParam(r, userIdDataKey)
	if len(qp) == 0 {
		return id, errors.New("no user Id provided")
	}

	id, err = uuid.Parse(qp)
	if err != nil {
		return id, errors.Wrap(err, "invalid user id provided")
	}

	return id, nil
}
