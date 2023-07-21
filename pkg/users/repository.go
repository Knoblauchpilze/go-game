package users

import (
	"context"

	"github.com/google/uuid"
)

// https://threedots.tech/post/repository-pattern-in-go/
type Repository interface {
	Create(ctx context.Context, user User) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAll(ctx context.Context) ([]uuid.UUID, error)
}
