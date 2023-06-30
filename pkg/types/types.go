package types

import (
	"github.com/KnoblauchPilze/go-game/pkg/auth"
	"github.com/google/uuid"
)

type UserData struct {
	Name     string
	Password string
}

type SignUpResponse struct {
	Id uuid.UUID
}

type LoginResponse struct {
	Token auth.Token
}
