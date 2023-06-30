package auth

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	User       uuid.UUID
	Value      string
	Expiration time.Time
}

type Authenticater interface {
	GenerateToken(user uuid.UUID, password string) (Token, error)
	GetToken(user uuid.UUID) (Token, error)
}
