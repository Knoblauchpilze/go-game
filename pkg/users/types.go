package users

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Name     string
	Password string
}

type UserManager interface {
	AddUser(name string, password string) (uuid.UUID, error)
	GetUser(id uuid.UUID) (User, error)
	GetUserFromName(name string) (User, error)
	GetUsers() []uuid.UUID
}
