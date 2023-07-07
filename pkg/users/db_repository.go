package users

import (
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type userDbRepo struct {
	db db.Database
}

const userTableName = "users"

const userIdColumnName = "id"
const userMailColumnName = "mail"
const userNameColumnName = "name"
const userPasswordColumnName = "password"
const userCreatedAtColumnName = "created_at"

var insertQueryBuilderFunc = db.NewInsertQueryBuilder

func NewDbRepository(db db.Database) Repository {
	return &userDbRepo{
		db: db,
	}
}

func (repo *userDbRepo) Create(user User) (uuid.UUID, error) {
	if err := user.validate(); err != nil {
		return uuid.UUID{}, err
	}
	out := user.Id

	qb := insertQueryBuilderFunc()

	qb.SetTable(userTableName)

	qb.AddElement(userIdColumnName, user.Id)
	qb.AddElement(userMailColumnName, user.Mail)
	qb.AddElement(userNameColumnName, user.Name)
	qb.AddElement(userPasswordColumnName, user.Password)

	qb.SetVerbose(true)

	query, err := qb.Build()
	if err != nil {
		return out, errors.WrapCode(err, errors.ErrUserCreationFailure)
	}

	rows := repo.db.Query(query)
	if err := rows.Err(); err != nil {
		return out, errors.WrapCode(err, errors.ErrUserCreationFailure)
	}
	rows.Close()

	return out, nil
}

func (repo *userDbRepo) Get(id uuid.UUID) (User, error) {
	var user User

	qb := db.NewSelectQueryBuilder()
	qb.SetTable(userTableName)

	qb.AddProp(userIdColumnName)
	qb.AddProp(userMailColumnName)
	qb.AddProp(userNameColumnName)
	qb.AddProp(userPasswordColumnName)
	qb.AddProp(userCreatedAtColumnName)

	qb.SetVerbose(true)

	query, err := qb.Build()
	if err != nil {
		return user, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	rows := repo.db.Query(query)
	if err := rows.Err(); err != nil {
		return user, err
	}

	defer rows.Close()

	scanner := func(row db.Scannable) error {
		return row.Scan(&user.Id, &user.Mail, &user.Name, &user.Password, &user.CreatedAt)
	}

	if err := rows.GetSingleValue(scanner); err != nil {
		return user, errors.WrapCode(err, errors.ErrDbCorruptedData)
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

	rows := repo.db.Query(query)
	if err := rows.Err(); err != nil {
		return users, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	defer rows.Close()

	scanner := func(row db.Scannable) error {
		var id uuid.UUID
		if err := row.Scan(&id); err != nil {
			return err
		}

		users = append(users, id)
		return nil
	}

	if err := rows.GetAll(scanner); err != nil {
		return users, errors.WrapCode(err, errors.ErrDbCorruptedData)
	}

	return users, nil
}
