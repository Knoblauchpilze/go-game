package db

import (
	"context"
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
	err := db.Connect(context.TODO())
	assert.Nil(err)
}

func TestPostgresDatabase_Connect_Fail(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return nil, fmt.Errorf("someError")
	}

	db := NewPostgresDatabase(config)

	err := db.Connect(context.TODO())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbConnectionFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
}

func TestPostgresDatabase_Disconnect_NotConnected(t *testing.T) {
	assert := assert.New(t)

	db := NewPostgresDatabase(testConfig)

	err := db.Disconnect(context.TODO())
	assert.Nil(err)
}

func TestPostgresDatabase_Disconnect(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	err := db.Disconnect(ctx)
	assert.Nil(err)
	assert.Equal(1, mockDb.closeCalled)
}

func TestPostgresDatabase_Query_NotConnected(t *testing.T) {
	assert := assert.New(t)

	db := NewPostgresDatabase(testConfig)

	q := queryImpl{}
	rows := db.Query(context.TODO(), q)
	assert.True(errors.IsErrorWithCode(rows.Err(), errors.ErrDbConnectionInvalid))
}

func TestPostgresDatabase_Query_Invalid(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{}
	rows := db.Query(ctx, q)
	assert.True(errors.IsErrorWithCode(rows.Err(), errors.ErrInvalidQuery))
}

func TestPostgresDatabase_Query(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{
		sqlCode: "someSqlCode",
	}
	rows := db.Query(ctx, q)
	assert.Nil(rows.Err())
	assert.Equal(1, len(mockDb.sqlQueriesReceived))
	assert.Equal("someSqlCode", mockDb.sqlQueriesReceived[0])
}

func TestPostgresDatabase_Query_Verbose(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{
		sqlCode: "someSqlCode",
		verbose: true,
	}
	rows := db.Query(ctx, q)
	assert.Nil(rows.Err())
	assert.Equal(1, len(mockDb.sqlQueriesReceived))
	assert.Equal("someSqlCode", mockDb.sqlQueriesReceived[0])
}

func TestPostgresDatabase_Query_Fail(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{
		queryError: fmt.Errorf("someError"),
	}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{
		sqlCode: "someSqlCode",
	}
	rows := db.Query(ctx, q)
	assert.True(errors.IsErrorWithCode(rows.Err(), errors.ErrDbRequestFailed))
	cause := errors.Unwrap(rows.Err())
	assert.Equal("someError", cause.Error())
}

func TestPostgresDatabase_Execute_NotConnected(t *testing.T) {
	assert := assert.New(t)

	db := NewPostgresDatabase(testConfig)

	q := queryImpl{}
	rows := db.Execute(context.TODO(), q)
	assert.True(errors.IsErrorWithCode(rows.Err, errors.ErrDbConnectionInvalid))
}

func TestPostgresDatabase_Execute_Invalid(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{}
	rows := db.Execute(ctx, q)
	assert.True(errors.IsErrorWithCode(rows.Err, errors.ErrInvalidQuery))
}

func TestPostgresDatabase_Execute(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{
		sqlCode: "someSqlCode",
	}
	result := db.Execute(ctx, q)
	assert.Nil(result.Err)
	assert.Equal(1, len(mockDb.sqlExecuteReceived))
	assert.Equal("someSqlCode", mockDb.sqlExecuteReceived[0])
}

func TestPostgresDatabase_Execute_Verbose(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{
		sqlCode: "someSqlCode",
		verbose: true,
	}
	result := db.Execute(ctx, q)
	assert.Nil(result.Err)
	assert.Equal(1, len(mockDb.sqlExecuteReceived))
	assert.Equal("someSqlCode", mockDb.sqlExecuteReceived[0])
}

func TestPostgresDatabase_Execute_Fail(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	mockDb := mockPgxDbFacade{
		execError: fmt.Errorf("someError"),
	}
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		return &mockDb, nil
	}
	ctx := context.TODO()

	db := NewPostgresDatabase(config)
	db.Connect(ctx)

	q := queryImpl{
		sqlCode: "someSqlCode",
	}
	result := db.Execute(ctx, q)
	assert.True(errors.IsErrorWithCode(result.Err, errors.ErrDbRequestFailed))
	cause := errors.Unwrap(result.Err)
	assert.Equal("someError", cause.Error())
}
