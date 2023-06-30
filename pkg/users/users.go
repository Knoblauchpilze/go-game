package users

import (
	"sync"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type userDbImpl struct {
	lock  sync.Mutex
	users map[string]string
	ids   map[uuid.UUID]string
	names map[string]uuid.UUID
}

func NewUserManager() UserManager {
	return &userDbImpl{
		users: make(map[string]string),
		ids:   make(map[uuid.UUID]string),
		names: make(map[string]uuid.UUID),
	}
}

func (udb *userDbImpl) AddUser(name string, password string) (uuid.UUID, error) {
	if len(name) == 0 {
		return uuid.UUID{}, errors.NewCode(errors.ErrInvalidUserName)
	}
	if len(password) == 0 {
		return uuid.UUID{}, errors.NewCode(errors.ErrInvalidPassword)
	}

	udb.lock.Lock()
	defer udb.lock.Unlock()

	if _, ok := udb.users[name]; ok {
		return uuid.UUID{}, errors.NewCode(errors.ErrUserAlreadyExists)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return uuid.UUID{}, errors.NewCode(errors.ErrUserCreationFailure)
	}

	udb.users[name] = password
	udb.ids[id] = name
	udb.names[name] = id

	return id, nil
}

func (udb *userDbImpl) GetUser(id uuid.UUID) (User, error) {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	user := User{
		Id: id,
	}

	name, ok := udb.ids[id]
	if !ok {
		return user, errors.NewCode(errors.ErrNoSuchUser)
	}

	password, ok := udb.users[name]
	if !ok {
		return user, errors.NewCode(errors.ErrNoSuchUser)
	}

	user.Name = name
	user.Password = password

	return user, nil
}

func (udb *userDbImpl) GetUserFromName(name string) (User, error) {
	user := User{
		Name: name,
	}

	if len(name) == 0 {
		return user, errors.NewCode(errors.ErrNoSuchUser)
	}

	udb.lock.Lock()
	defer udb.lock.Unlock()

	id, ok := udb.names[name]
	if !ok {
		return user, errors.NewCode(errors.ErrNoSuchUser)
	}

	password, ok := udb.users[name]
	if !ok {
		return user, errors.NewCode(errors.ErrNoSuchUser)
	}

	user.Id = id
	user.Password = password

	return user, nil
}

func (udb *userDbImpl) GetUsers() []uuid.UUID {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	users := make([]uuid.UUID, 0, len(udb.ids))
	for id := range udb.ids {
		users = append(users, id)
	}

	return users
}
