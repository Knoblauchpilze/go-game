package users

import (
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type userDbRepo struct {
	qe db.QueryExecutor
}

const userTableName = "users"

const userIdColumnName = "id"
const userMailColumnName = "mail"
const userNameColumnName = "name"
const userPasswordColumnName = "password"
const userCreatedAtColumnName = "created_at"

var insertQueryBuilderFunc = db.NewInsertQueryBuilder
var selectQueryBuilderFunc = db.NewSelectQueryBuilder
var inFilterBuilderFunc = db.NewInFilterBuilder
var deleteQueryBuilderFunc = db.NewDeleteQueryBuilder

func NewDbRepository(qe db.QueryExecutor) Repository {
	return &userDbRepo{
		qe: qe,
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

	if err := repo.qe.RunQuery(qb); err != nil {
		return out, errors.WrapCode(err, errors.ErrUserCreationFailure)
	}

	return out, nil
}

func (repo *userDbRepo) Get(id uuid.UUID) (User, error) {
	var user User

	qb := selectQueryBuilderFunc()

	qb.SetTable(userTableName)

	qb.AddProp(userIdColumnName)
	qb.AddProp(userMailColumnName)
	qb.AddProp(userNameColumnName)
	qb.AddProp(userPasswordColumnName)
	qb.AddProp(userCreatedAtColumnName)

	fb := inFilterBuilderFunc()
	fb.SetKey(userIdColumnName)
	fb.AddValue(id)
	f, err := fb.Build()
	if err != nil {
		return user, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	qb.SetFilter(f)

	qb.SetVerbose(true)

	scanner := func(row db.Scannable) error {
		return row.Scan(&user.Id, &user.Mail, &user.Name, &user.Password, &user.CreatedAt)
	}

	if err := repo.qe.RunQueryAndScanSingleResult(qb, scanner); err != nil {
		return user, errors.WrapCode(err, errors.ErrUserGetFailure)
	}

	return user, nil
}

func (repo *userDbRepo) Patch(id uuid.UUID, patch User) (User, error) {
	return User{}, errors.NewCode(errors.ErrNotImplemented)
}

func (repo *userDbRepo) Delete(id uuid.UUID) error {
	qb := deleteQueryBuilderFunc()

	qb.SetTable(userTableName)

	fb := inFilterBuilderFunc()
	fb.SetKey(userIdColumnName)
	fb.AddValue(id)
	f, err := fb.Build()
	if err != nil {
		return errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	qb.SetFilter(f)

	qb.SetVerbose(true)

	if err := repo.qe.RunQuery(qb); err != nil {
		return errors.WrapCode(err, errors.ErrUserDeletionFailure)
	}

	return nil
}

func (repo *userDbRepo) GetAll() ([]uuid.UUID, error) {
	var users []uuid.UUID

	qb := selectQueryBuilderFunc()
	qb.SetTable("users")

	qb.AddProp(userIdColumnName)

	qb.SetVerbose(true)

	scanner := func(row db.Scannable) error {
		var id uuid.UUID
		if err := row.Scan(&id); err != nil {
			return err
		}

		users = append(users, id)
		return nil
	}

	if err := repo.qe.RunQueryAndScanAllResults(qb, scanner); err != nil {
		return users, errors.WrapCode(err, errors.ErrUserGetFailure)
	}

	return users, nil
}
