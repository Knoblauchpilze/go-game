package users

import (
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type userDbRepo struct {
}

func NewDbRepository() Repository {
	return &userDbRepo{}
}

func (repo *userDbRepo) Create(user User) (uuid.UUID, error) {
	return uuid.UUID{}, errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) Get(id uuid.UUID) (User, error) {
	return User{}, errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) Patch(id uuid.UUID, patch User) (User, error) {
	return User{}, errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) Delete(id uuid.UUID) error {
	return errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) GetAll() ([]uuid.UUID, error) {
	return []uuid.UUID{}, errors.NewCode(errors.ErrNotImplemented)
}
