package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgToStr_String(t *testing.T) {
	assert := assert.New(t)

	arg := "hello"

	out, err := argToStr(arg)
	assert.Nil(err)
	assert.Equal("hello", out)
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
	return nil, fmt.Errorf("someError")
}

func TestArgToStr_Unmarshalable(t *testing.T) {
	assert := assert.New(t)

	arg := mockUnmarshalable{}

	_, err := argToStr(arg)
	assert.Contains(err.Error(), "someError")
}

type mockUnmarshalableConvertible struct{}

func (mc mockUnmarshalableConvertible) Convert() interface{} {
	return mockUnmarshalable{}
}

func TestArgToStr_UnmarshalableConvertible(t *testing.T) {
	assert := assert.New(t)

	arg := mockUnmarshalableConvertible{}

	_, err := argToStr(arg)
	assert.Contains(err.Error(), "someError")
}
