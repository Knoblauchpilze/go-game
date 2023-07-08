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

	m := &mockDb{
		rows: &mockRows{},
	}
	repo := NewDbRepository(m)

	_, err := repo.Create(User{})
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserMail))
}

func TestDbRepository_CreateUser_InvalidName(t *testing.T) {
	assert := assert.New(t)

	m := &mockDb{
		rows: &mockRows{},
	}
	repo := NewDbRepository(m)

	user := User{Mail: "some@mail"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserName))
}

func TestDbRepository_CreateUser_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	m := &mockDb{
		rows: &mockRows{},
	}
	repo := NewDbRepository(m)

	user := User{Mail: "some@mail", Name: "someName"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidPassword))
}

func TestDbRepository_CreateUser_QueryBuildError(t *testing.T) {
	assert := assert.New(t)
	// https://stackoverflow.com/questions/61107654/make-go-tests-independent-from-each-other-mutation-of-global-vars
	t.Cleanup(resetQueryBuilderFuncs)

	m := &mockDb{
		rows: &mockRows{},
	}
	insertQueryBuilderFunc = func() db.InsertQueryBuilder {
		return mockInsertQueryBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(m)

	_, err := repo.Create(defaultTestUser)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_CreateUser_DbQueryError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDb{
		rows: &mockRows{
			err: fmt.Errorf("someError"),
		},
	}
	repo := NewDbRepository(m)

	_, err := repo.Create(defaultTestUser)
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserCreationFailure))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_CreateUser(t *testing.T) {
	assert := assert.New(t)

	r := &mockRows{}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	out, err := repo.Create(defaultTestUser)
	assert.Nil(err)

	assert.Equal(1, m.queryCalls)
	assert.True(m.queries[0].Valid())
	testUserInsertQuery := "INSERT INTO users (id, mail, name, password) VALUES (\"08ce96a3-3430-48a8-a3b2-b1c987a207ca\", some@mail, someName, somePassword)"
	assert.Equal(testUserInsertQuery, m.queries[0].ToSql())

	assert.Equal(1, r.closeCalled)

	assert.Equal(defaultTestUser.Id, out)
}

func TestDbRepository_GetUser_QueryBuildError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	m := &mockDb{
		rows: &mockRows{},
	}
	selectQueryBuilderFunc = func() db.SelectQueryBuilder {
		return mockSelectQueryBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(m)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetUser_FilterBuildError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	m := &mockDb{
		rows: &mockRows{},
	}
	inFilterBuilderFunc = func() db.InFilterBuilder {
		return mockFilterBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(m)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetUser_DbQueryError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDb{
		rows: &mockRows{
			err: fmt.Errorf("someError"),
		},
	}
	repo := NewDbRepository(m)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetUser_RowsFailure(t *testing.T) {
	assert := assert.New(t)

	r := &mockRows{
		getSingleValueErr: fmt.Errorf("someError"),
	}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	_, err := repo.Get(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbCorruptedData))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetUser_Scanner(t *testing.T) {
	assert := assert.New(t)

	s := &mockScannable{}
	r := &mockRows{
		getSingleValueScannable: s,
	}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	_, err := repo.Get(uuid.New())
	assert.Nil(err)

	assert.Equal(1, s.scanCalled)
}

func TestDbRepository_GetUser(t *testing.T) {
	assert := assert.New(t)

	r := &mockRows{}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	_, err := repo.Get(defaultTestUser.Id)
	assert.Nil(err)

	assert.Equal(1, m.queryCalls)
	assert.Equal(1, r.singleValueCalled)
	assert.True(m.queries[0].Valid())
	testUserSelectQuery := "SELECT id, mail, name, password, created_at FROM users WHERE id in ('08ce96a3-3430-48a8-a3b2-b1c987a207ca')"
	assert.Equal(testUserSelectQuery, m.queries[0].ToSql())

	assert.Equal(1, r.closeCalled)
}

func TestDbRepository_Delete_QueryBuildError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	m := &mockDb{
		rows: &mockRows{},
	}
	deleteQueryBuilderFunc = func() db.DeleteQueryBuilder {
		return mockSelectQueryBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(m)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_Delete_FilterBuildError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	m := &mockDb{
		rows: &mockRows{},
	}
	inFilterBuilderFunc = func() db.InFilterBuilder {
		return mockFilterBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(m)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_Delete_DbQueryError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDb{
		rows: &mockRows{
			err: fmt.Errorf("someError"),
		},
	}
	repo := NewDbRepository(m)

	err := repo.Delete(uuid.New())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_Delete(t *testing.T) {
	assert := assert.New(t)

	r := &mockRows{}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	err := repo.Delete(defaultTestUser.Id)
	assert.Nil(err)

	assert.Equal(1, m.queryCalls)
	assert.True(m.queries[0].Valid())
	testUserDeleteQuery := "DELETE FROM users WHERE id in ('08ce96a3-3430-48a8-a3b2-b1c987a207ca')"
	assert.Equal(testUserDeleteQuery, m.queries[0].ToSql())

	assert.Equal(1, r.closeCalled)
}

func TestDbRepository_GetAll_QueryBuildError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetQueryBuilderFuncs)

	m := &mockDb{
		rows: &mockRows{},
	}
	selectQueryBuilderFunc = func() db.SelectQueryBuilder {
		return mockSelectQueryBuilder{
			buildErr: fmt.Errorf("someError"),
		}
	}
	repo := NewDbRepository(m)

	_, err := repo.GetAll()
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetAll_DbQueryError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDb{
		rows: &mockRows{
			err: fmt.Errorf("someError"),
		},
	}
	repo := NewDbRepository(m)

	_, err := repo.GetAll()
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetAll_RowsFailure(t *testing.T) {
	assert := assert.New(t)

	r := &mockRows{
		getAllErr: fmt.Errorf("someError"),
	}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	_, err := repo.GetAll()
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbCorruptedData))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestDbRepository_GetAll_Scanner(t *testing.T) {
	assert := assert.New(t)

	s := &mockScannable{}
	r := &mockRows{
		getAllScannable: s,
	}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	_, err := repo.GetAll()
	assert.Nil(err)

	assert.Equal(1, s.scanCalled)
}

func TestDbRepository_GetAll_ScannerError(t *testing.T) {
	assert := assert.New(t)

	s := &mockScannable{
		scanErr: fmt.Errorf("someError"),
	}
	r := &mockRows{
		getAllScannable: s,
	}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	repo.GetAll()

	assert.Equal("someError", r.allScanErr.Error())
	assert.Equal(1, s.scanCalled)
}

func TestDbRepository_GetAll(t *testing.T) {
	assert := assert.New(t)

	r := &mockRows{}
	m := &mockDb{
		rows: r,
	}
	repo := NewDbRepository(m)

	_, err := repo.GetAll()
	assert.Nil(err)

	assert.Equal(1, m.queryCalls)
	assert.Equal(1, r.allCalled)
	assert.True(m.queries[0].Valid())
	testUsersSelectQuery := "SELECT id FROM users"
	assert.Equal(testUsersSelectQuery, m.queries[0].ToSql())

	assert.Equal(1, r.closeCalled)
}

func resetQueryBuilderFuncs() {
	insertQueryBuilderFunc = db.NewInsertQueryBuilder
	selectQueryBuilderFunc = db.NewSelectQueryBuilder
	inFilterBuilderFunc = db.NewInFilterBuilder
	deleteQueryBuilderFunc = db.NewDeleteQueryBuilder
}

type mockDb struct {
	connectErr    error
	disconnectErr error

	queryCalls int
	queries    []db.Query
	rows       db.Rows

	executeCalls int
	executions   []db.Query
	result       db.Result
}

func (m *mockDb) Connect() error {
	return m.connectErr
}

func (m *mockDb) Disconnect() error {
	return m.disconnectErr
}

func (m *mockDb) Query(query db.Query) db.Rows {
	m.queries = append(m.queries, query)
	m.queryCalls++
	return m.rows
}

func (m *mockDb) Execute(query db.Query) db.Result {
	m.executions = append(m.executions, query)
	m.executeCalls++
	return m.result
}

type mockRows struct {
	err         error
	closeCalled int
	empty       bool

	singleValueCalled       int
	getSingleValueScannable db.Scannable
	singleValueScanErr      error
	getSingleValueErr       error

	allCalled       int
	getAllScannable db.Scannable
	allScanErr      error
	getAllErr       error
}

func (m *mockRows) Err() error {
	return m.err
}

func (m *mockRows) Close() {
	m.closeCalled++
}

func (m *mockRows) Empty() bool {
	return m.empty
}

func (m *mockRows) GetSingleValue(scanner db.ScanRow) error {
	m.singleValueCalled++
	if m.getSingleValueScannable != nil {
		m.singleValueScanErr = scanner(m.getSingleValueScannable)
	}
	return m.getSingleValueErr
}

func (m *mockRows) GetAll(scanner db.ScanRow) error {
	m.allCalled++
	if m.getAllScannable != nil {
		m.allScanErr = scanner(m.getAllScannable)
	}
	return m.getAllErr
}

type mockInsertQueryBuilder struct {
	buildErr error
}

func (m mockInsertQueryBuilder) SetTable(table string) error {
	return nil
}

func (m mockInsertQueryBuilder) AddElement(column string, value interface{}) error {
	return nil
}

func (m mockInsertQueryBuilder) SetFilter(filter db.Filter) error {
	return nil
}

func (m mockInsertQueryBuilder) SetVerbose(verbose bool) {}

func (m mockInsertQueryBuilder) Build() (db.Query, error) {
	return nil, m.buildErr
}

type mockSelectQueryBuilder struct {
	buildErr error
}

func (m mockSelectQueryBuilder) SetTable(table string) error {
	return nil
}

func (m mockSelectQueryBuilder) AddProp(prop string) error {
	return nil
}

func (m mockSelectQueryBuilder) SetFilter(filter db.Filter) error {
	return nil
}

func (m mockSelectQueryBuilder) SetVerbose(verbose bool) {}

func (m mockSelectQueryBuilder) Build() (db.Query, error) {
	return nil, m.buildErr
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

type mockScannable struct {
	scanCalled int
	scanErr    error
}

func (m *mockScannable) Scan(dest ...interface{}) error {
	m.scanCalled++
	return m.scanErr
}
