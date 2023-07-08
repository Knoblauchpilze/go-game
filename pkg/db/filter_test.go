package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Valid(t *testing.T) {
	assert := assert.New(t)

	f := filterImpl{}
	assert.False(f.Valid())

	f.sqlCode = "someSqlCode"
	assert.True(f.Valid())
}

func TestFilter_ToSql(t *testing.T) {
	assert := assert.New(t)

	f := filterImpl{}
	assert.Equal("", f.ToSql())

	f.sqlCode = "someSqlCode"
	assert.Equal("someSqlCode", f.ToSql())
}
