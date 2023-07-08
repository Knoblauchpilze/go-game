package db

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteQueryBuilder_SetTable(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()

	err := b.SetTable("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlTable))

	err = b.SetTable("haha")
	assert.Nil(err)
}

func TestDeleteQueryBuilder_SetFilter(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()

	f := filterImpl{}
	err := b.SetFilter(f)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlFilter))

	f.sqlCode = "someSqlCode"
	err = b.SetFilter(f)
	assert.Nil(err)
}

func TestDeleteQueryBuilder_SetFilterOverride(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()

	b.SetTable("table")
	f := filterImpl{
		sqlCode: "someSqlCode",
	}
	b.SetFilter(f)
	f.sqlCode = "someOtherCode"
	b.SetFilter(f)

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("DELETE FROM table WHERE someOtherCode", query.ToSql())
}

func TestDeleteQueryBuilder_SetVerbose(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()

	b.SetTable("table")
	_, err := b.Build()
	assert.Nil(err)

	b.SetVerbose(true)
	_, err = b.Build()
	assert.Nil(err)
}

func TestDeleteQueryBuilder_Build_NoTable(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrInvalidSqlTable))
}

func TestDeleteQueryBuilder_Build(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()
	b.SetTable("table")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("DELETE FROM table", query.ToSql())
}

func TestDeleteQueryBuilder_Build_WithFilter(t *testing.T) {
	assert := assert.New(t)

	b := NewDeleteQueryBuilder()
	b.SetTable("table")
	f := filterImpl{
		sqlCode: "someFilter",
	}
	b.SetFilter(f)

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("DELETE FROM table WHERE someFilter", query.ToSql())
}
