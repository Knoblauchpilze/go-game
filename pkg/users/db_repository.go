package users

import (
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type userDbRepo struct {
	db db.Database
}

const userIdColumnName = "id"
const userMailColumnName = "mail"
const userNameColumnName = "name"
const userPasswordColumnName = "password"

func NewDbRepository(db db.Database) Repository {
	return &userDbRepo{
		db: db,
	}
}

func (repo *userDbRepo) Create(user User) (uuid.UUID, error) {
	return uuid.UUID{}, errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) Get(id uuid.UUID) (User, error) {
	var user User

	qb := db.NewSelectQueryBuilder()
	qb.SetTable("users")

	qb.AddProp(userIdColumnName)
	qb.AddProp(userMailColumnName)
	qb.AddProp(userNameColumnName)
	qb.AddProp(userPasswordColumnName)

	qb.SetVerbose(true)

	query, err := qb.Build()
	if err != nil {
		return user, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	result := repo.db.Query(query)
	if err := result.Err(); err != nil {
		return user, err
	}

	defer result.Close()
	var idStr string
	result.Scan(&idStr, &user.Mail, &user.Name, &user.Password)
	if id, err := uuid.Parse(idStr); err != nil {
		return user, errors.WrapCode(err, errors.ErrDbCorruptedData)
	} else {
		user.Id = id
	}

	if result.Next() {
		return user, errors.NewCode(errors.ErrMultiValuedDbElement)
	}

	return user, nil
}

func (repo *userDbRepo) Patch(id uuid.UUID, patch User) (User, error) {
	return User{}, errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) Delete(id uuid.UUID) error {
	return errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) GetAll() ([]uuid.UUID, error) {
	var users []uuid.UUID

	qb := db.NewSelectQueryBuilder()
	qb.SetTable("users")

	qb.AddProp(userIdColumnName)

	qb.SetVerbose(true)

	query, err := qb.Build()
	if err != nil {
		return users, err
	}

	result := repo.db.Query(query)
	if err := result.Err(); err != nil {
		return users, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	defer result.Close()
	for result.Next() {
		var idStr string
		result.Scan(&idStr)
		if id, err := uuid.Parse(idStr); err != nil {
			return users, errors.WrapCode(err, errors.ErrDbCorruptedData)
		} else {
			users = append(users, id)
		}
	}

	return users, nil
}
