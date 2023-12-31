package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArgToStr_String(t *testing.T) {
	assert := assert.New(t)

	arg := "hello"

	out, err := argToStr(arg)
	assert.Nil(err)
	assert.Equal("hello", out)
}

func TestArgToStr_Uuid(t *testing.T) {
	assert := assert.New(t)

	arg := uuid.New()

	out, err := argToStr(arg)
	assert.Nil(err)
	assert.Equal(arg.String(), out)
}

type mockConvertible struct {
	value int
}

func (mc mockConvertible) Convert() interface{} {
	return mc.value
}

func TestArgToStr_Convertible(t *testing.T) {
	assert := assert.New(t)

	arg := mockConvertible{value: 32}

	out, err := argToStr(arg)
	assert.Nil(err)
	assert.Equal("32", out)
}

type mockComplexArg struct {
	Value int
	Name  string
}

func TestArgToStr_ComplexArg(t *testing.T) {
	assert := assert.New(t)

	arg := mockComplexArg{Value: 26, Name: "someName"}

	out, err := argToStr(arg)
	assert.Nil(err)
	assert.Equal("{\"Value\":26,\"Name\":\"someName\"}", out)
}

type mockUnmarshalable struct{}

func (mu mockUnmarshalable) MarshalJSON() ([]byte, error) {
	return nil, errDefault
}

func TestArgToStr_Unmarshalable(t *testing.T) {
	assert := assert.New(t)

	arg := mockUnmarshalable{}

	_, err := argToStr(arg)
	assert.Contains(err.Error(), errDefault.Error())
}

type mockUnmarshalableConvertible struct{}

func (mc mockUnmarshalableConvertible) Convert() interface{} {
	return mockUnmarshalable{}
}

func TestArgToStr_UnmarshalableConvertible(t *testing.T) {
	assert := assert.New(t)

	arg := mockUnmarshalableConvertible{}

	_, err := argToStr(arg)
	assert.Contains(err.Error(), errDefault.Error())
}

func TestSqlPropAsUpdateToStr(t *testing.T) {
	assert := assert.New(t)

	update := sqlProp{
		column: "column",
		value:  32,
	}

	out, err := sqlPropAsUpdateToStr(update)
	assert.Nil(err)
	assert.Equal("column = '32'", out)
}

func TestSqlPropAsUpdateToStr_Unmarshalable(t *testing.T) {
	assert := assert.New(t)

	update := sqlProp{
		column: "column",
		value:  mockUnmarshalable{},
	}

	_, err := sqlPropAsUpdateToStr(update)
	assert.Contains(err.Error(), errDefault.Error())
}
