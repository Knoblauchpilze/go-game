package dtos

import (
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/google/uuid"
)

type UserDto struct {
	Mail     string
	Name     string
	Password string
}

func (dto UserDto) Convert() users.User {
	return users.User{
		Mail:     dto.Mail,
		Name:     dto.Name,
		Password: dto.Password,
	}
}

type SignUpResponse struct {
	Id uuid.UUID
}
