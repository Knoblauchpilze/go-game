package db

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSelectQueryBuilder_SetTable(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()

	err := b.SetTable("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlTable))

	err = b.SetTable("haha")
	assert.Nil(err)
}

func TestSelectQueryBuilder_AddProp(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()

	err := b.AddProp("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlProp))

	prop := "someProp"
	err = b.AddProp(prop)
	assert.Nil(err)

	err = b.AddProp(prop)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDuplicatedSqlProp))
}

func TestSelectQueryBuilder_SetFilter(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()

	f := filterImpl{}
	err := b.SetFilter(f)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlFilter))

	f.sqlCode = "someSqlCode"
	err = b.SetFilter(f)
	assert.Nil(err)
}

func TestSelectQueryBuilder_SetFilterOverride(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()

	b.SetTable("table")
	b.AddProp("prop")
	f := filterImpl{
		sqlCode: "someSqlCode",
	}
	b.SetFilter(f)
	f.sqlCode = "someOtherCode"
	b.SetFilter(f)

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("SELECT prop FROM table WHERE someOtherCode", query.ToSql())
}

func TestSelectQueryBuilder_SetVerbose(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()

	b.SetTable("table")
	b.AddProp("prop")
	_, err := b.Build()
	assert.Nil(err)

	b.SetVerbose(true)
	_, err = b.Build()
	assert.Nil(err)
}

func TestSelectQueryBuilder_Build_NoTable(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrInvalidSqlTable))
}

func TestSelectQueryBuilder_Build_NoProp(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()
	b.SetTable("table")

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrNoPropInSqlSelectQuery))
}

func TestSelectQueryBuilder_Build(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()
	b.SetTable("table")
	b.AddProp("prop")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("SELECT prop FROM table", query.ToSql())
}

func TestSelectQueryBuilder_Build_MultiArgs(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()
	b.SetTable("table")
	b.AddProp("prop1")
	b.AddProp("prop2")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("SELECT prop1, prop2 FROM table", query.ToSql())
}

func TestSelectQueryBuilder_Build_WithFilter(t *testing.T) {
	assert := assert.New(t)

	b := NewSelectQueryBuilder()
	b.SetTable("table")
	b.AddProp("prop1")
	f := filterImpl{
		sqlCode: "someFilter",
	}
	b.SetFilter(f)

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.Equal("SELECT prop1 FROM table WHERE someFilter", query.ToSql())
}
