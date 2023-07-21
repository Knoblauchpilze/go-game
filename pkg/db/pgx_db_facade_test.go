package db

import (
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestPgxDbFacade_New(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetPgxConnFunc)

	pgxConnectionFunc = func(config pgx.ConnPoolConfig) (*pgx.ConnPool, error) {
		return nil, errDefault
	}

	_, err := newPgxDbFacadeImpl(pgx.ConnPoolConfig{})
	assert.Equal(errDefault, err)
}

func TestPgxDbFacade_Close(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConn{}
	f := pgxDbFacadeImpl{
		pool: m,
	}

	f.Close()
	assert.Equal(1, m.closeCalled)
}

func TestPgxDbFacade_Query(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConn{
		queryErr: errDefault,
	}
	f := pgxDbFacadeImpl{
		pool: m,
	}

	rows, err := f.Query("someSql")
	assert.Nil(rows)
	assert.Equal(errDefault, err)
}

func TestPgxDbFacade_Exec(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConn{
		execError: errDefault,
	}
	f := pgxDbFacadeImpl{
		pool: m,
	}

	tag, err := f.Exec("someSql")
	assert.Equal(pgx.CommandTag(""), tag)
	assert.Equal(errDefault, err)
}

type mockPgxDbConn struct {
	closeCalled int
	queryErr    error
	execError   error
}

func (m *mockPgxDbConn) Close() {
	m.closeCalled++
}

func (m *mockPgxDbConn) Query(sql string, args ...interface{}) (*pgx.Rows, error) {
	return nil, m.queryErr
}

func (m *mockPgxDbConn) Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error) {
	return "", m.execError
}

func resetPgxConnFunc() {
	pgxConnectionFunc = pgx.NewConnPool
}
