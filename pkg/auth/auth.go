package auth

import (
	"sync"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type authImpl struct {
	lock   sync.Mutex
	tokens map[uuid.UUID]Token
}

var TokenDefaultExpirationTime = 10 * time.Minute

func NewAuthenticator() Authenticator {
	return &authImpl{
		tokens: make(map[uuid.UUID]Token),
	}
}

func (auth *authImpl) GenerateToken(user uuid.UUID, password string) (Token, error) {
	if len(password) == 0 {
		return Token{}, errors.NewCode(errors.ErrInvalidPassword)
	}

	auth.lock.Lock()
	defer auth.lock.Unlock()

	if _, ok := auth.tokens[user]; ok {
		return Token{}, errors.NewCode(errors.ErrTokenAlreadyExists)
	}

	token := Token{
		User:       user,
		Expiration: time.Now().Add(TokenDefaultExpirationTime),
		Value:      "dummy-token",
	}

	auth.tokens[user] = token

	return token, nil
}

func (auth *authImpl) GetToken(user uuid.UUID) (Token, error) {
	auth.lock.Lock()
	defer auth.lock.Unlock()

	token, ok := auth.tokens[user]
	if !ok {
		return Token{}, errors.NewCode(errors.ErrNoSuchToken)
	}

	return token, nil
}
