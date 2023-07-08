package users

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserValidate_NoEmail(t *testing.T) {
	assert := assert.New(t)

	user := User{
		Id:       uuid.New(),
		Name:     "someName",
		Password: "password",
	}

	err := user.validate()

	errCode, ok := err.(errors.ErrorWithCode)
	assert.True(ok)
	assert.Equal(errors.ErrInvalidUserMail, errCode.Code())
}

func TestUserValidate_NoName(t *testing.T) {
	assert := assert.New(t)

	user := User{
		Id:       uuid.New(),
		Mail:     "some@mail",
		Password: "password",
	}

	err := user.validate()

	errCode, ok := err.(errors.ErrorWithCode)
	assert.True(ok)
	assert.Equal(errors.ErrInvalidUserName, errCode.Code())
}

func TestUserValidate_NoPassword(t *testing.T) {
	assert := assert.New(t)

	user := User{
		Id:   uuid.New(),
		Mail: "some@mail",
		Name: "someName",
	}

	err := user.validate()

	errCode, ok := err.(errors.ErrorWithCode)
	assert.True(ok)
	assert.Equal(errors.ErrInvalidPassword, errCode.Code())
}

func TestUserValidate(t *testing.T) {
	assert := assert.New(t)

	user := User{
		Id:       uuid.New(),
		Mail:     "some@mail",
		Name:     "someName",
		Password: "password",
	}

	err := user.validate()
	assert.Nil(err)
}
