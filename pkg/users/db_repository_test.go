package users

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultTestUser = User{Id: uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca"), Mail: "some@mail", Name: "someName", Password: "somePassword"}

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

	mqe := &mockQueryExecutor{
		runQueryAndScanSingleResultErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserGetFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetUser_FilterBuildError(t *testing.T) {
	assert := assert.New(t)
	// https://stackoverflow.com/questions/61107654/make-go-tests-independent-from-each-other-mutation-of-global-vars
	t.Cleanup(resetQueryBuilderFuncs)

	mqe := &mockQueryExecutor{}
	inFilterBuilderFunc = func() db.InFilterBuilder {
		return mockFilterBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(mqe)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

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

	mqe := &mockQueryExecutor{
		runQueryErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserDeletionFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_Delete_FilterBuildError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	mqe := &mockQueryExecutor{}
	inFilterBuilderFunc = func() db.InFilterBuilder {
		return mockFilterBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(mqe)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

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

	mqe := &mockQueryExecutor{
		runQueryAndScanAllResultsErr: fmt.Errorf("someError"),
	}
	repo := NewDbRepository(mqe)

	_, err := repo.GetAll()
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserGetFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetAll(t *testing.T) {
	assert := assert.New(t)

	mqe := &mockQueryExecutor{
		result:  2,
		scanner: &mockScannable{},
	}
	repo := NewDbRepository(mqe)

	out, err := repo.GetAll()
	assert.Nil(err)
	assert.Equal(2, len(out))
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

	result  int
	scanner *mockScannable

	queries []db.QueryBuilder
}

func (m *mockQueryExecutor) RunQuery(qb db.QueryBuilder) error {
	m.runQueryCalled++
	m.queries = append(m.queries, qb)

	return m.runQueryErr
}

func (m *mockQueryExecutor) RunQueryAndScanSingleResult(qb db.QueryBuilder, parser db.RowParser) error {
	m.runQueryAndScanSingleResultCalled++
	m.queries = append(m.queries, qb)

	if m.scanner != nil && m.result > 0 {
		for id := 0; id < m.result; id++ {
			parser.ScanRow(m.scanner)
		}
	}

	return m.runQueryAndScanSingleResultErr
}

func (m *mockQueryExecutor) RunQueryAndScanAllResults(qb db.QueryBuilder, parser db.RowParser) error {
	m.runQueryAndScanAllResultsCalled++
	m.queries = append(m.queries, qb)

	if m.scanner != nil && m.result > 0 {
		for id := 0; id < m.result; id++ {
			parser.ScanRow(m.scanner)
		}
	}

	return m.runQueryAndScanAllResultsErr
}

type mockFilterBuilder struct {
	buildErr error
}

func (m mockFilterBuilder) SetKey(key string) error {
	return nil
}

func (m mockFilterBuilder) AddValue(value interface{}) error {
	return nil
}

func (m mockFilterBuilder) Build() (db.Filter, error) {
	return nil, m.buildErr
}
