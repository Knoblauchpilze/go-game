package users

import "github.com/google/uuid"

// https://threedots.tech/post/repository-pattern-in-go/
type Repository interface {
	Create(user User) (uuid.UUID, error)
	Get(id uuid.UUID) (User, error)
	Patch(id uuid.UUID, patch User) (User, error)
	Delete(id uuid.UUID) error

	GetAll() ([]uuid.UUID, error)
}
