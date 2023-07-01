package dtos

import (
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/google/uuid"
)

type UserDto struct {
	Id       uuid.UUID
	Mail     string
	Name     string
	Password string
}

func (dto UserDto) Convert() users.User {
	return users.User{
		Id:       dto.Id,
		Mail:     dto.Mail,
		Name:     dto.Name,
		Password: dto.Password,
	}
}

type PostResponse string
type GetResponse UserDto
type GetAllResponse []string
type PatchResponse UserDto
type DeleteResponse string
