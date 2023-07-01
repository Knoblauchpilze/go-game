package users

import (
	"sync"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type userMemoryRepo struct {
	lock  sync.Mutex
	users map[uuid.UUID]User
}

func NewMemoryRepository() Repository {
	return &userMemoryRepo{
		users: make(map[uuid.UUID]User),
	}
}

func (repo *userMemoryRepo) Create(user User) (uuid.UUID, error) {
	if err := user.validate(); err != nil {
		return uuid.UUID{}, err
	}

	repo.lock.Lock()
	defer repo.lock.Unlock()

	if _, ok := repo.users[user.Id]; ok {
		return uuid.UUID{}, errors.NewCode(errors.ErrUserAlreadyExists)
	}

	repo.users[user.Id] = user

	return user.Id, nil
}

func (repo *userMemoryRepo) Get(id uuid.UUID) (User, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	user, ok := repo.users[id]
	if !ok {
		return User{}, errors.NewCode(errors.ErrNoSuchUser)
	}

	return user, nil
}

func (repo *userMemoryRepo) Patch(id uuid.UUID, patch User) (User, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	user, ok := repo.users[id]
	if !ok {
		return User{}, errors.NewCode(errors.ErrNoSuchUser)
	}

	if len(patch.Mail) > 0 {
		user.Mail = patch.Mail
	}
	if len(patch.Name) > 0 {
		user.Name = patch.Name
	}
	if len(patch.Password) > 0 {
		user.Password = patch.Password
	}

	repo.users[id] = user

	return user, nil
}

func (repo *userMemoryRepo) Delete(id uuid.UUID) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	_, ok := repo.users[id]
	delete(repo.users, id)

	if !ok {
		return errors.NewCode(errors.ErrNoSuchUser)
	}

	return nil
}

func (repo *userMemoryRepo) GetAll() ([]uuid.UUID, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	users := make([]uuid.UUID, 0, len(repo.users))
	for id := range repo.users {
		users = append(users, id)
	}

	return users, nil
}
