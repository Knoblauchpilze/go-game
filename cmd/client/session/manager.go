package session

import (
	"fmt"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/auth"
	"github.com/KnoblauchPilze/go-game/pkg/connection"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/KnoblauchPilze/go-game/pkg/types"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/google/uuid"
)

type Manager interface {
	SignUp(in types.UserData) error
	Login(in types.UserData) error

	Authenticate(token auth.Token) error

	ListUsers() ([]uuid.UUID, error)
	ListUser(id uuid.UUID) (users.User, error)
}

type managerImpl struct {
	userId uuid.UUID
	token  auth.Token
	url    string
}

func NewManager(url string) Manager {
	return &managerImpl{url: url}
}

func (mi *managerImpl) SignUp(in types.UserData) error {
	var out types.SignUpResponse

	signUpUrl := fmt.Sprintf("%s/signup", mi.url)

	rb := connection.NewHttpPostRequestBuilder()
	rb.SetUrl(signUpUrl)
	rb.SetBody("application/json", in)

	req, err := rb.Build()
	if err != nil {
		return err
	}
	resp, err := req.Perform()
	if err != nil {
		return err
	}

	err = rest.GetBodyFromHttpResponseAs(resp, &out)
	if err != nil {
		return err
	}

	mi.userId = out.Id
	logger.Infof("Signed up with id %v", mi.userId)

	return nil
}

func (mi *managerImpl) Login(in types.UserData) error {
	var out types.LoginResponse

	loginUrl := fmt.Sprintf("%s/login", mi.url)

	rb := connection.NewHttpPostRequestBuilder()
	rb.SetUrl(loginUrl)
	rb.SetBody("application/json", in)

	req, err := rb.Build()
	if err != nil {
		return err
	}
	resp, err := req.Perform()
	if err != nil {
		return err
	}

	err = rest.GetBodyFromHttpResponseAs(resp, &out)
	if err != nil {
		return err
	}

	mi.token = out.Token
	logger.Infof("Logged in, active token is %+v", mi.token)

	return nil
}

func (mi *managerImpl) Authenticate(token auth.Token) error {
	if len(token.Value) == 0 {
		return errors.NewCode(errors.ErrNotLoggedIn)
	}
	if time.Now().After(token.Expiration) {
		return errors.NewCode(errors.ErrAuthenticationExpired)
	}

	mi.token = token

	return nil
}

func (mi *managerImpl) ListUsers() ([]uuid.UUID, error) {
	var out []uuid.UUID

	listUsersURL := fmt.Sprintf("%s/users", mi.url)

	auth, err := mi.generateAuthenticationHeader()
	if err != nil {
		return out, err
	}

	rb := connection.NewHttpGetRequestBuilder()
	rb.SetUrl(listUsersURL)
	rb.SetHeaders(map[string][]string{
		"Authorization": {auth},
	})

	req, err := rb.Build()
	if err != nil {
		return out, err
	}
	resp, err := req.Perform()
	if err != nil {
		return out, err
	}

	err = rest.GetBodyFromHttpResponseAs(resp, &out)
	if err != nil {
		return out, errors.WrapCode(err, errors.ErrGetRequestFailed)
	}

	return out, nil
}

func (mi *managerImpl) ListUser(id uuid.UUID) (users.User, error) {
	var out users.User

	listUserURL := fmt.Sprintf("%s/users/%s", mi.url, id)

	auth, err := mi.generateAuthenticationHeader()
	if err != nil {
		return out, err
	}

	rb := connection.NewHttpGetRequestBuilder()
	rb.SetUrl(listUserURL)
	rb.SetHeaders(map[string][]string{
		"Authorization": {auth},
	})

	req, err := rb.Build()
	if err != nil {
		return out, err
	}
	resp, err := req.Perform()
	if err != nil {
		return out, err
	}

	err = rest.GetBodyFromHttpResponseAs(resp, &out)
	if err != nil {
		return out, errors.WrapCode(err, errors.ErrGetRequestFailed)
	}

	return out, nil
}

func (mi *managerImpl) generateAuthenticationHeader() (string, error) {
	if len(mi.token.Value) == 0 {
		return "", errors.NewCode(errors.ErrNotLoggedIn)
	}
	if time.Now().After(mi.token.Expiration) {
		return "", errors.NewCode(errors.ErrAuthenticationExpired)
	}

	auth := fmt.Sprintf("bearer user=%v token=%v", mi.token.User, mi.token.Value)

	return auth, nil
}
