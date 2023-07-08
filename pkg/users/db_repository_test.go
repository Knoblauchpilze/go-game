package users

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDbRepository_CreateUser_InvalidMail(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	_, err := repo.Create(User{})
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserMail))
}

func TestDbRepository_CreateUser_InvalidName(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	user := User{Mail: "some@mail"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserName))
}

func TestDbRepository_CreateUser_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	user := User{Mail: "some@mail", Name: "someName"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidPassword))
}

func TestDbRepository_CreateUser_QueryExecutorError(t *testing.T) {
	assert := assert.New(t)
	// https://stackoverflow.com/questions/61107654/make-go-tests-independent-from-each-other-mutation-of-global-vars
	t.Cleanup(resetQueryBuilderFuncs)

	mqe := &mockQueryExecutor{
		runQueryErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	_, err := repo.Create(defaultTestUser)
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserCreationFailure))
	cause := errors.Unwrap(err)
	assert.Contains(cause.Error(), "someError")
}

func TestDbRepository_CreateUser(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	out, err := repo.Create(defaultTestUser)
	assert.Nil(err)
	assert.Equal(defaultTestUser.Id, out)
	assert.Equal(1, mqe.runQueryCalled)
	assert.Equal(1, len(mqe.queries))

	q, err := mqe.queries[0].Build()
	assert.Nil(err)
	expectedQuery := "INSERT INTO users (id, mail, name, password) VALUES ('08ce96a3-3430-48a8-a3b2-b1c987a207ca', 'some@mail', 'someName', 'somePassword')"
	assert.Equal(expectedQuery, q.ToSql())
}

func TestDbRepository_GetUser_QueryExecutorError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	mqe := &mockQueryExecutor{
		runQueryAndScanSingleResultErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserGetFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

// func TestDbRepository_GetUser_FilterBuildError(t *testing.T) {
// 	assert := assert.New(t)
// 	t.Cleanup(resetQueryBuilderFuncs)

// 	m := &mockDb{
// 		rows: &mockRows{},
// 	}
// 	inFilterBuilderFunc = func() db.InFilterBuilder {
// 		return mockFilterBuilder{
// 			buildErr: fmt.Errorf("someError"),
// 		}
// 	}
// 	repo := NewDbRepository(m)

// 	_, err := repo.Get(uuid.New())
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
// 	cause := errors.Unwrap(err)
// 	assert.Equal("someError", cause.Error())
// }

func TestDbRepository_GetUser(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	_, err := repo.Get(defaultTestUser.Id)
	assert.Nil(err)
	assert.Equal(1, mqe.runQueryAndScanSingleResultCalled)
	assert.Equal(1, len(mqe.queries))

	q, err := mqe.queries[0].Build()
	assert.Nil(err)
	expectedQuery := "SELECT id, mail, name, password, created_at FROM users WHERE id in ('08ce96a3-3430-48a8-a3b2-b1c987a207ca')"
	assert.Equal(expectedQuery, q.ToSql())
}

func TestDbRepository_Delete_QueryExecutorError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	mqe := &mockQueryExecutor{
		runQueryErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserDeletionFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

// func TestDbRepository_Delete_FilterBuildError(t *testing.T) {
// 	assert := assert.New(t)
// 	t.Cleanup(resetQueryBuilderFuncs)

// 	m := &mockDb{
// 		rows: &mockRows{},
// 	}
// 	inFilterBuilderFunc = func() db.InFilterBuilder {
// 		return mockFilterBuilder{
// 			buildErr: fmt.Errorf("someError"),
// 		}
// 	}
// 	repo := NewDbRepository(m)

// 	err := repo.Delete(uuid.New())
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
// 	cause := errors.Unwrap(err)
// 	assert.Equal("someError", cause.Error())
// }

func TestDbRepository_Delete(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	err := repo.Delete(defaultTestUser.Id)
	assert.Nil(err)
	assert.Equal(1, mqe.runQueryCalled)
	assert.Equal(1, len(mqe.queries))

	q, err := mqe.queries[0].Build()
	assert.Nil(err)
	expectedQuery := "DELETE FROM users WHERE id in ('08ce96a3-3430-48a8-a3b2-b1c987a207ca')"
	assert.Equal(expectedQuery, q.ToSql())
}

func TestDbRepository_GetAll_QueryExecutorError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	mqe := &mockQueryExecutor{
		runQueryErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserDeletionFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetAll(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{}
	repo := NewDbRepository(mqe)

	_, err := repo.GetAll()
	assert.Nil(err)
	assert.Equal(1, mqe.runQueryAndScanAllResultsCalled)
	assert.Equal(1, len(mqe.queries))

	q, err := mqe.queries[0].Build()
	assert.Nil(err)
	expectedQuery := "SELECT id FROM users"
	assert.Equal(expectedQuery, q.ToSql())
}

func resetQueryBuilderFuncs() {
	insertQueryBuilderFunc = db.NewInsertQueryBuilder
	selectQueryBuilderFunc = db.NewSelectQueryBuilder
	inFilterBuilderFunc = db.NewInFilterBuilder
	deleteQueryBuilderFunc = db.NewDeleteQueryBuilder
}

type mockQueryExecutor struct {
	runQueryCalled int
	runQueryErr    error

	runQueryAndScanSingleResultCalled int
	runQueryAndScanSingleResultErr    error

	runQueryAndScanAllResultsCalled int
	runQueryAndScanAllResultsErr    error

	queries []db.QueryBuilder
}

func (m *mockQueryExecutor) RunQuery(qb db.QueryBuilder) error {
	m.runQueryCalled++
	m.queries = append(m.queries, qb)
	return m.runQueryErr
}

func (m *mockQueryExecutor) RunQueryAndScanSingleResult(qb db.QueryBuilder, scan db.ScanRow) error {
	m.runQueryAndScanSingleResultCalled++
	m.queries = append(m.queries, qb)
	return m.runQueryAndScanSingleResultErr
}

func (m *mockQueryExecutor) RunQueryAndScanAllResults(qb db.QueryBuilder, scan db.ScanRow) error {
	m.runQueryAndScanAllResultsCalled++
	m.queries = append(m.queries, qb)
	return m.runQueryAndScanAllResultsErr
}
