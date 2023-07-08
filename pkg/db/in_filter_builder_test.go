package db

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInFilterBuilder_SetKey(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()

	err := b.SetKey("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlComparisonKey))

	err = b.SetKey("key")
	assert.Nil(err)
}

func TestInFilterBuilder_AddArg(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()

	err := b.AddValue(nil)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlComparisonValue))

	err = b.AddValue("value")
	assert.Nil(err)
}

func TestInFilterBuilder_Build_NoKey(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrInvalidSqlComparisonKey))
}

func TestInFilterBuilder_Build_NoValue(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()
	b.SetKey("key")

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrNoValuesInSqlComparison))
}

func TestInFilterBuilder_Build_StringValue(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()
	b.SetKey("key")
	b.AddValue("value")

	filter, err := b.Build()
	assert.Nil(err)
	assert.True(filter.Valid())
	assert.Equal("key in ('value')", filter.ToSql())
}

func TestInFilterBuilder_Build_TimeValue(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()
	b.SetKey("key")
	someTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	b.AddValue(someTime)

	filter, err := b.Build()
	assert.Nil(err)
	assert.True(filter.Valid())
	assert.Equal("key in ('2009-11-17T20:34:58Z')", filter.ToSql())
}

func TestInFilterBuilder_Build_MultiArgs(t *testing.T) {
	assert := assert.New(t)

	b := NewInFilterBuilder()
	b.SetKey("key")
	b.AddValue("value1")
	b.AddValue("value2")

	filter, err := b.Build()
	assert.Nil(err)
	assert.True(filter.Valid())
	assert.Equal("key in ('value1', 'value2')", filter.ToSql())
}
