package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_Valid(t *testing.T) {
	assert := assert.New(t)

	q := queryImpl{}
	assert.False(q.Valid())

	q.sqlCode = "someSqlCode"
	assert.True(q.Valid())
}

func TestQuery_Verbose(t *testing.T) {
	assert := assert.New(t)

	q := queryImpl{}
	assert.False(q.Verbose())

	q.verbose = true
	assert.True(q.Verbose())
}

func TestQuery_ToSql(t *testing.T) {
	assert := assert.New(t)

	q := queryImpl{}
	assert.Equal("", q.ToSql())

	q.sqlCode = "someSqlCode"
	assert.Equal("someSqlCode", q.ToSql())
}
