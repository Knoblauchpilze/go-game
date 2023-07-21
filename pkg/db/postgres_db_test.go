package db

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	queryDelay time.Duration
	rows       sqlRows
	queryError error

	sqlQueriesReceived []string

	execDelay time.Duration
	tag       pgx.CommandTag
	execError error

	sqlExecuteReceived []string

	closeCalled int
}

func (m *mockPgxDbFacade) Close() {
	m.closeCalled++
}

func (m *mockPgxDbFacade) Query(sql string, args ...interface{}) (sqlRows, error) {
	m.sqlQueriesReceived = append(m.sqlQueriesReceived, sql)
	if m.queryDelay > 0 {
		time.Sleep(m.queryDelay)
	}
	return m.rows, m.queryError
}

func (m *mockPgxDbFacade) Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error) {
	m.sqlExecuteReceived = append(m.sqlExecuteReceived, sql)
	if m.execDelay > 0 {
		time.Sleep(m.execDelay)
	}
	return m.tag, m.execError
}

var defaultSleep = 100 * time.Millisecond

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

func TestPostgresDatabase_Connect_Timeout(t *testing.T) {
	assert := assert.New(t)

	timeout := 100 * time.Millisecond
	config := testConfig
	config.DbConnectionTimeout = timeout
	config.creationFunc = func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
		time.Sleep(2 * timeout)
		return &mockPgxDbFacade{}, nil
	}

	db := NewPostgresDatabase(config)

	err := db.Connect(context.TODO())
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbConnectionTimeout))
	cause := errors.Unwrap(err)
	assert.Equal(context.DeadlineExceeded, cause)
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

func TestPostgresDatabase_Query_Timeout(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	config.DbQueryTimeout = defaultSleep / 2
	mockRows := &mockSqlRows{}
	mockDb := mockPgxDbFacade{
		queryDelay: defaultSleep,
		rows:       mockRows,
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
	// Wait for query to finish with some margin.
	time.Sleep(2 * defaultSleep)
	err := rows.Err()
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestTimeout))
	cause := errors.Unwrap(err)
	assert.Equal(context.DeadlineExceeded, cause)
	assert.Equal(int32(1), mockRows.closeCalls.Load())
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

func TestPostgresDatabase_Execute_Timeout(t *testing.T) {
	assert := assert.New(t)

	config := testConfig
	config.DbQueryTimeout = defaultSleep / 2
	mockDb := mockPgxDbFacade{
		execDelay: defaultSleep,
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
	// Wait for query to finish with some margin.
	time.Sleep(2 * defaultSleep)
	assert.True(errors.IsErrorWithCode(result.Err, errors.ErrDbRequestTimeout))
	cause := errors.Unwrap(result.Err)
	assert.Equal(context.DeadlineExceeded, cause)
}
