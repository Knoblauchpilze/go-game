package db

import (
	"strings"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUpdateQueryBuilder_SetTable(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()

	err := b.SetTable("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlTable))

	err = b.SetTable("haha")
	assert.Nil(err)
}

func TestUpdateQueryBuilder_AddUpdate(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()

	err := b.AddUpdate("", 32)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlColumn))

	column := "column"
	prop := "someProp"
	err = b.AddUpdate(column, prop)
	assert.Nil(err)

	err = b.AddUpdate(column, prop)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDuplicatedSqlColumn))
}

func TestUpdateQueryBuilder_SetFilter(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()

	f := filterImpl{}
	err := b.SetFilter(f)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlFilter))

	f.sqlCode = "someSqlCode"
	err = b.SetFilter(f)
	assert.Nil(err)
}

func TestUpdateQueryBuilder_SetVerbose(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()

	b.SetTable("table")
	b.AddUpdate("column", "prop")
	_, err := b.Build()
	assert.Nil(err)

	b.SetVerbose(true)
	_, err = b.Build()
	assert.Nil(err)
}

func TestUpdateQueryBuilder_Build_NoTable(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrInvalidSqlTable))
}

func TestUpdateQueryBuilder_Build_NoColumn(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()
	b.SetTable("table")

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrNoColumnInSqlUpdateQuery))
}

func TestUpdateQueryBuilder_Build(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()
	b.SetTable("table")
	b.AddUpdate("column", "prop")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("UPDATE table SET column = 'prop'", query.ToSql())
}

func TestUpdateQueryBuilder_Build_MultiColumns(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()
	b.SetTable("table")
	b.AddUpdate("column1", "prop1")
	b.AddUpdate("column2", "prop2")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("UPDATE table SET column1 = 'prop1', column2 = 'prop2'", query.ToSql())
}

func TestUpdateQueryBuilder_Build_WithFilter(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()
	b.SetTable("table")
	b.AddUpdate("column", "prop")
	f := filterImpl{
		sqlCode: "someFilter",
	}
	b.SetFilter(f)

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("UPDATE table SET column = 'prop' WHERE someFilter", query.ToSql())
}

func TestUpdateQueryBuilder_Build_ArgWithError(t *testing.T) {
	assert := assert.New(t)

	b := NewUpdateQueryBuilder()
	b.SetTable("table")
	b.AddUpdate("column", mockUnmarshalable{})

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(strings.Contains(cause.Error(), "someError"))
}
