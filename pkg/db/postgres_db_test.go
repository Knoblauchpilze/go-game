package db

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

var testConfig = Config{
	DbHost:                "host",
	DbPort:                652,
	DbName:                "database",
	DbUser:                "user",
	DbPassword:            "password",
	DbConnectionsPoolSize: 2,
}

type mockPgxDbFacade struct {
	rows       *pgx.Rows
	queryError error

	sqlQueriesReceived []string

	tag       pgx.CommandTag
	execError error

	sqlExecuteReceived []string

	closeCalled int
}

func (m *mockPgxDbFacade) Close() {
	m.closeCalled++
}

func (m *mockPgxDbFacade) Query(sql string, args ...interface{}) (*pgx.Rows, error) {
	m.sqlQueriesReceived = append(m.sqlQueriesReceived, sql)
	return m.rows, m.queryError
}

func (m *mockPgxDbFacade) Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error) {
	m.sqlExecuteReceived = append(m.sqlExecuteReceived, sql)
	return m.tag, m.execError
}

func TestPostgresDatabase_Connect(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockPgxDbFacade{}, nil
	}

	db := NewPostgresDatabase(config)
	err := db.Connect()
	assert.Nil(err)
}

func TestPostgresDatabase_Connect_Fail(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return nil, fmt.Errorf("someError")
	}

	db := NewPostgresDatabase(config)

	err := db.Connect()
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbConnectionFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestPostgresDatabase_Disconnect_NotConnected(t *testing.T) {
	assert := assert.New(t)

	db := NewPostgresDatabase(testConfig)

	err := db.Disconnect()
	assert.Nil(err)
}

func TestPostgresDatabase_Disconnect(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}

	db := NewPostgresDatabase(config)
	db.Connect()

	err := db.Disconnect()
	assert.Nil(err)
	assert.Equal(1, mockDb.closeCalled)
}

func TestPostgresDatabase_Query_NotConnected(t *testing.T) {
	assert := assert.New(t)

	db := NewPostgresDatabase(testConfig)

	q := queryImpl{}
	rows := db.Query(q)
	assert.True(errors.IsErrorWithCode(rows.Err(), errors.ErrDbConnectionInvalid))
}

func TestPostgresDatabase_Query_Invalid(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}

	db := NewPostgresDatabase(config)
	db.Connect()

	q := queryImpl{}
	rows := db.Query(q)
	assert.True(errors.IsErrorWithCode(rows.Err(), errors.ErrInvalidQuery))
}

func TestPostgresDatabase_Query(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}

	db := NewPostgresDatabase(config)
	db.Connect()

	q := queryImpl{
		sqlCode: "someSqlCode",
	}
	rows := db.Query(q)
	assert.Nil(rows.Err())
	assert.Equal(1, len(mockDb.sqlQueriesReceived))
	assert.Equal("someSqlCode", mockDb.sqlQueriesReceived[0])
}

func TestPostgresDatabase_Execute_NotConnected(t *testing.T) {
	assert := assert.New(t)

	db := NewPostgresDatabase(testConfig)

	q := queryImpl{}
	rows := db.Execute(q)
	assert.True(errors.IsErrorWithCode(rows.Err, errors.ErrDbConnectionInvalid))
}

func TestPostgresDatabase_Execute_Invalid(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}

	db := NewPostgresDatabase(config)
	db.Connect()

	q := queryImpl{}
	rows := db.Execute(q)
	assert.True(errors.IsErrorWithCode(rows.Err, errors.ErrInvalidQuery))
}

func TestPostgresDatabase_Execute(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}

	db := NewPostgresDatabase(config)
	db.Connect()

	q := queryImpl{
		sqlCode: "someSqlCode",
	}
	result := db.Execute(q)
	assert.Nil(result.Err)
	assert.Equal(1, len(mockDb.sqlExecuteReceived))
	assert.Equal("someSqlCode", mockDb.sqlExecuteReceived[0])
}
