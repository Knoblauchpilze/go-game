package users

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type dbMock struct {
	connectErr    error
	disconnectErr error

	queryCalls int
	rows       db.Rows

	executeCalls int
	result       db.Result
}

func (m *dbMock) Connect() error {
	return m.connectErr
}

func (m *dbMock) Disconnect() error {
	return m.disconnectErr
}

func (m *dbMock) Query(query db.Query) db.Rows {
	m.queryCalls++
	return m.rows
}

func (m *dbMock) Execute(query db.Query) db.Result {
	m.executeCalls++
	return m.result
}

type rowsMock struct {
	err         error
	closeCalled int
	empty       bool

	singleValueCalled int
	getSingleValueErr error

	allCalled int
	getAllErr error
}

func (m *rowsMock) Err() error {
	return m.err
}

func (m *rowsMock) Close() {
	m.closeCalled++
}

func (m *rowsMock) Empty() bool {
	return m.empty
}

func (m *rowsMock) GetSingleValue(scanner db.ScanRow) error {
	m.singleValueCalled++
	return m.getSingleValueErr
}

func (m *rowsMock) GetAll(scanner db.ScanRow) error {
	m.allCalled++
	return m.getAllErr
}

func TestDbRepository_GetUser_DbQueryError(t *testing.T) {
	assert := assert.New(t)

	r := &rowsMock{
		getSingleValueErr: fmt.Errorf("someError"),
	}
	m := &dbMock{
		rows: r,
	}
	repo := NewDbRepository(m)

	someId := uuid.New()
	_, err := repo.Get(someId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbCorruptedData))
	assert.Equal(1, m.queryCalls)
	assert.Equal(1, r.singleValueCalled)
}
