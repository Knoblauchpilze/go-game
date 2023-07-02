package db

import (
	"fmt"
	"strings"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestScriptQueryBuilder_SetHasReturnValue(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()

	b.SetScript("script")
	b.AddArg("arg")
	query, err := b.Build()
	assert.Nil(err)
	assert.Equal("SELECT script('arg')", query.ToSql())

	b.SetHasReturnValue(true)
	query, err = b.Build()
	assert.Nil(err)
	assert.Equal("SELECT * FROM script('arg')", query.ToSql())
}

func TestScriptQueryBuilder_SetScript(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()

	err := b.SetScript("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlScript))

	err = b.SetScript("script")
	assert.Nil(err)
}

func TestScriptQueryBuilder_AddArg(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()

	err := b.AddArg(nil)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlScriptArg))

	err = b.AddArg("arg")
	assert.Nil(err)
}

func TestScriptQueryBuilder_SetVerbose(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()

	b.SetScript("script")
	b.AddArg("arg")
	query, err := b.Build()
	assert.Nil(err)
	assert.False(query.Verbose())

	b.SetVerbose(true)
	query, err = b.Build()
	assert.Nil(err)
	assert.True(query.Verbose())
}

func TestScriptQueryBuilder_Build_NoScript(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, errors.ErrInvalidSqlScript))
}

func TestScriptQueryBuilder_Build_NoArg(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.False(query.Verbose())
	assert.Equal("SELECT script()", query.ToSql())
}

func TestScriptQueryBuilder_Build_StringArg(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")
	b.AddArg("arg")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.False(query.Verbose())
	assert.Equal("SELECT script('arg')", query.ToSql())
}

type mockConvertible struct {
	value int
}

func (mc mockConvertible) Convert() interface{} {
	return mc.value
}

func TestScriptQueryBuilder_Build_ConvertibleArg(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")
	b.AddArg(mockConvertible{value: 32})

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.False(query.Verbose())
	assert.Equal("SELECT script('32')", query.ToSql())
}

type mockComplexArg struct {
	Value int
	Name  string
}

func TestScriptQueryBuilder_Build_ComplexArg(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")
	b.AddArg(mockComplexArg{Value: 26, Name: "someName"})

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.False(query.Verbose())
	assert.Equal("SELECT script('{\"Value\":26,\"Name\":\"someName\"}')", query.ToSql())
}

type mockUnmarshalable struct{}

func (mu mockUnmarshalable) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("someError")
}

func TestScriptQueryBuilder_Build_Unmarshallable(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")
	b.AddArg(mockUnmarshalable{})

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	fmt.Printf("err: %+v\n", cause)
	assert.True(strings.Contains(cause.Error(), "someError"))
}

type mockUnmarshalableConvertible struct{}

func (mc mockUnmarshalableConvertible) Convert() interface{} {
	return mockUnmarshalable{}
}

func TestScriptQueryBuilder_Build_UnmarshallableConvertible(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")
	b.AddArg(mockUnmarshalableConvertible{})

	_, err := b.Build()
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlTranslationFailed))
	cause := errors.Unwrap(err)
	assert.True(strings.Contains(cause.Error(), "someError"))
}

func TestScriptQueryBuilder_Build_MultiArgs(t *testing.T) {
	assert := assert.New(t)

	b := NewScriptQueryBuilder()
	b.SetScript("script")
	b.AddArg("arg1")
	b.AddArg("arg2")

	query, err := b.Build()
	assert.Nil(err)
	assert.True(query.Valid())
	assert.False(query.Verbose())
	assert.Equal("SELECT script('arg1', 'arg2')", query.ToSql())
}
