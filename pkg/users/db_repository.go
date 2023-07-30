package users

import (
	"context"

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

func (repo *userDbRepo) Create(ctx context.Context, user User) (uuid.UUID, error) {
	out := user.Id
	if err := user.validate(); err != nil {
		return out, err
	}

	qb := insertQueryBuilderFunc()

	qb.SetTable(userTableName)

	qb.AddElement(userIdColumnName, user.Id)
	qb.AddElement(userMailColumnName, user.Mail)
	qb.AddElement(userNameColumnName, user.Name)
	qb.AddElement(userPasswordColumnName, user.Password)

	qb.SetVerbose(true)

	if err := repo.qe.ExecuteQuery(ctx, qb); err != nil {
		return out, errors.WrapCode(err, errors.ErrUserCreationFailure)
	}

	return out, nil
}

func (repo *userDbRepo) Get(ctx context.Context, id uuid.UUID) (User, error) {
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
		return User{}, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	qb.SetFilter(f)

	qb.SetVerbose(true)

	scanner := &userRowParser{}
	if err := repo.qe.RunQueryAndScanSingleResult(ctx, qb, scanner); err != nil {
		return User{}, errors.WrapCode(err, errors.ErrUserGetFailure)
	}

	return scanner.user, nil
}

func (repo *userDbRepo) Delete(ctx context.Context, id uuid.UUID) error {
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

	if err := repo.qe.ExecuteQuery(ctx, qb); err != nil {
		return errors.WrapCode(err, errors.ErrUserDeletionFailure)
	}

	return nil
}

func (repo *userDbRepo) GetAll(ctx context.Context) ([]uuid.UUID, error) {
	qb := selectQueryBuilderFunc()
	qb.SetTable("users")

	qb.AddProp(userIdColumnName)

	qb.SetVerbose(true)

	scanner := &userIdsParser{}
	if err := repo.qe.RunQueryAndScanAllResults(ctx, qb, scanner); err != nil {
		return []uuid.UUID{}, errors.WrapCode(err, errors.ErrUserGetFailure)
	}

	return scanner.ids, nil
}
