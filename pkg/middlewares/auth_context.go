package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/auth"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/google/uuid"
)

var authenticationScheme = "bearer"
var authenticationUserKey = "user"
var authenticationTokenKey = "token"

var genErrMsg = "invalid authentication header"

func GenerateAuthenticationContext(tokens auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// https://stackoverflow.com/questions/33265812/best-http-authorization-header-type-for-jwt
			// https://reqbin.com/req/5k564bhv/get-request-bearer-token-authorization-header-example
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization
			reqData, ok := GetRequestDataFromContextOrFail(w, r)
			if !ok {
				return
			}

			authData, err := rest.GetSingleHeaderFromHttpRequest(r, "Authorization")
			if err != nil {
				reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
				return
			}

			token, err := parseAuthenticationHeader(authData)
			if err != nil {
				reqData.FailWithErrorAndCode(err, http.StatusBadRequest, w)
				return
			}

			check, err := tokens.GetToken(token.User)
			if err != nil {
				reqData.FailWithErrorAndCode(errors.NewCode(errors.ErrAuthenticationFailure), http.StatusUnauthorized, w)
				return
			}
			if token.Value != check.Value {
				reqData.FailWithErrorAndCode(errors.NewCode(errors.ErrAuthenticationFailure), http.StatusUnauthorized, w)
				return
			}
			if time.Now().After(check.Expiration) {
				reqData.FailWithErrorAndCode(errors.NewCode(errors.ErrAuthenticationExpired), http.StatusUnauthorized, w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseAuthenticationHeader(authData string) (auth.Token, error) {
	var out auth.Token

	tokens := strings.Split(authData, " ")
	if len(tokens) != 3 {
		return out, errors.New(genErrMsg)
	}

	if strings.ToLower(tokens[0]) != authenticationScheme {
		return out, errors.New(genErrMsg)
	}

	props := make(map[string]string)
	tokens = tokens[1:]
	for _, prop := range tokens {
		keyValue := strings.Split(prop, "=")
		if len(keyValue) != 2 {
			err := errors.Newf("Ill-formed prop in Authorization header: \"%s\"", prop)
			return out, errors.Wrap(err, genErrMsg)
		}
		if len(keyValue[0]) == 0 || len(keyValue[1]) == 0 {
			err := errors.Newf("Ill-formed prop in Authorization header: \"%s\"", prop)
			return out, errors.Wrap(err, genErrMsg)
		}

		props[keyValue[0]] = keyValue[1]
	}

	userId, ok := props[authenticationUserKey]
	if !ok {
		err := errors.Newf("no \"%s\" key in Authorization header", authenticationUserKey)
		return out, errors.Wrap(err, genErrMsg)
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		err := errors.Wrapf(err, "failed to parse user id \"%s\"", userId)
		return out, errors.Wrap(err, genErrMsg)
	}

	token, ok := props[authenticationTokenKey]
	if !ok {
		err := errors.Newf("no \"%s\" key in Authorization header", authenticationTokenKey)
		return out, errors.Wrap(err, genErrMsg)
	}

	out.User = id
	out.Value = token

	return out, nil
}
