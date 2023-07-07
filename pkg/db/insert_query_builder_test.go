package db

import (
	"strings"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInsertQueryBuilder_SetTable(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()

	err := b.SetTable("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlTable))

	err = b.SetTable("haha")
	assert.Nil(err)
}

func TestInsertQueryBuilder_AddElement(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()

	err := b.AddElement("", 32)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlColumn))

	column := "column"
	prop := "someProp"
	err = b.AddElement(column, prop)
	assert.Nil(err)

	err = b.AddElement(column, prop)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDuplicatedSqlColumn))
}

func TestInsertQueryBuilder_SetVerbose(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()

	b.SetTable("table")
	b.AddElement("column", "prop")
	_, err := b.Build()
	assert.Nil(err)

	b.SetVerbose(true)
	_, err = b.Build()
	assert.Nil(err)
}

func TestInsertQueryBuilder_Build_NoTable(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrInvalidSqlTable))
}

func TestInsertQueryBuilder_Build_NoColumn(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()
	b.SetTable("table")

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrNoColumnInSqlInsertQuery))
}

func TestInsertQueryBuilder_Build(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()
	b.SetTable("table")
	b.AddElement("column", "prop")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("INSERT INTO table (column) VALUES (prop)", query.ToSql())
}

func TestInsertQueryBuilder_Build_MultiColumns(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()
	b.SetTable("table")
	b.AddElement("column1", "prop1")
	b.AddElement("column2", "prop2")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("INSERT INTO table (column1, column2) VALUES (prop1, prop2)", query.ToSql())
}

func TestInsertQueryBuilder_Build_ArgWithError(t *testing.T) {
	assert := assert.New(t)

	b := NewInsertQueryBuilder()
	b.SetTable("table")
	b.AddElement("column", mockUnmarshalable{})

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(strings.Contains(cause.Error(), "someError"))
}
