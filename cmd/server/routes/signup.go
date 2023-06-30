package routes

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/middlewares"
	"github.com/KnoblauchPilze/go-game/pkg/types"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
)

var SignUpURLRoute = "/signup"

func SignUpRouter(udb users.UserManager) http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.RequestCtx)
		r.Post("/", generateSignUpHandler(udb))
	})

	return r
}

func generateSignUpHandler(udb users.UserManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, ok := middlewares.GetRequestDataFromContextOrFail(w, r)
		if !ok {
			return
		}

		var err error
		var ud types.UserData

		if ud, err = getUserDataFromRequest(r); err != nil {
			reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
			return
		}

		out, err := signUpUser(ud, udb)
		if err != nil {
			errCode := interpretSignUpFailure(err)
			reqData.FailWithErrorAndCode(err, errCode, w)
			return
		}

		reqData.WriteDetails(out, w)
	}
}

func signUpUser(ud types.UserData, udb users.UserManager) (types.SignUpResponse, error) {
	var err error
	var out types.SignUpResponse

	out.Id, err = udb.AddUser(ud.Name, ud.Password)
	return out, err
}

func interpretSignUpFailure(err error) int {
	errCode := http.StatusBadRequest
	if errors.IsErrorWithCode(err, errors.ErrUserCreationFailure) {
		errCode = http.StatusInternalServerError
	}

	return errCode
}
