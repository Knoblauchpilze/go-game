package auth

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	auth := NewAuthenticater()
	id, _ := uuid.NewUUID()

	_, err := auth.GenerateToken(id, "")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidPassword))
}

func TestGenerateToken(t *testing.T) {
	assert := assert.New(t)

	auth := NewAuthenticater()
	user, _ := uuid.NewUUID()

	token, err := auth.GenerateToken(user, "foo")

	assert.Nil(err)
	assert.Equal(user, token.User)
	assert.Less(1, len(token.Value))
	assert.True(time.Now().Before(token.Expiration))

	_, err = auth.GenerateToken(user, "foo")
	assert.True(errors.IsErrorWithCode(err, errors.ErrTokenAlreadyExists))
}

func TestGetToken(t *testing.T) {
	assert := assert.New(t)

	auth := NewAuthenticater()
	user, _ := uuid.NewUUID()

	check, _ := auth.GenerateToken(user, "foo")
	token, err := auth.GetToken(user)

	assert.Nil(err)
	assert.Equal(check.User, token.User)
	assert.Equal(check.Value, token.Value)
	assert.Equal(check.Expiration, token.Expiration)
}

func TestGetToken_InvalidId(t *testing.T) {
	assert := assert.New(t)

	auth := NewAuthenticater()
	id, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()

	auth.GenerateToken(id, "foo")
	_, err := auth.GetToken(id2)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchToken))
}
