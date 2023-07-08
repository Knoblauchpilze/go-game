package users

import (
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Mail      string
	Name      string
	Password  string
	CreatedAt time.Time
}

func (u User) validate() error {
	if len(u.Mail) == 0 {
		return errors.NewCode(errors.ErrInvalidUserMail)
	}
	if len(u.Name) == 0 {
		return errors.NewCode(errors.ErrInvalidUserName)
	}
	if len(u.Password) == 0 {
		return errors.NewCode(errors.ErrInvalidPassword)
	}

	return nil
}
