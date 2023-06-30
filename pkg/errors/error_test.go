package errors

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errSomeError = fmt.Errorf("some error")

func TestError_New(t *testing.T) {
	assert := assert.New(t)

	err := New("haha")

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("haha", impl.Message)
	assert.Nil(impl.Cause)
	assert.False(impl.hasCode)
}

func TestError_NewFromCode(t *testing.T) {
	assert := assert.New(t)

	err := NewCode(ErrInvalidUserName)
	errMsg := errorsCodeToMessage[ErrInvalidUserName]

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal(errMsg, impl.Message)
	assert.Nil(impl.Cause)
	assert.True(impl.hasCode)
	assert.Equal(ErrInvalidUserName, impl.Value)

	err = NewCode(lastErrorCode)

	impl, ok = err.(errorImpl)
	assert.True(ok)
	assert.Equal(defaultErrorMessage, impl.Message)
	assert.Nil(impl.Cause)
	assert.True(impl.hasCode)
	assert.Equal(lastErrorCode, impl.Value)
}

func TestError_Newf(t *testing.T) {
	assert := assert.New(t)

	err := Newf("haha %d", 22)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("haha 22", impl.Message)
	assert.Nil(impl.Cause)
}

func TestError_Wrap(t *testing.T) {
	assert := assert.New(t)

	err := Wrap(errSomeError, "context")

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("context", impl.Message)
	assert.Equal(errSomeError, impl.Cause)
}

func TestError_WrapCode(t *testing.T) {
	assert := assert.New(t)

	err := WrapCode(errSomeError, ErrInvalidUserName)
	errMsg := errorsCodeToMessage[ErrInvalidUserName]

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal(errMsg, impl.Message)
	assert.Equal(errSomeError, impl.Cause)
	assert.True(impl.hasCode)
	assert.Equal(ErrInvalidUserName, impl.Value)
}

func TestError_WrapCodeInvalidCode(t *testing.T) {
	assert := assert.New(t)

	err := WrapCode(errSomeError, lastErrorCode)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal(defaultErrorMessage, impl.Message)
	assert.Equal(errSomeError, impl.Cause)
	assert.True(impl.hasCode)
	assert.Equal(lastErrorCode, impl.Value)
}

func TestError_Wrapf(t *testing.T) {
	assert := assert.New(t)

	err := Wrapf(errSomeError, "context %d", -44)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("context -44", impl.Message)
	assert.Equal(errSomeError, impl.Cause)
}

func TestError_Unwrap(t *testing.T) {
	assert := assert.New(t)

	err := Unwrap(nil)
	assert.Nil(err)

	err = Unwrap(errSomeError)
	assert.Nil(err)

	err = New("haha")
	cause := Unwrap(err)
	assert.Nil(cause)

	err = Wrap(errSomeError, "haha")
	cause = Unwrap(err)
	assert.Equal(errSomeError, cause)

	causeOfCause := Unwrap(cause)
	assert.Nil(causeOfCause)
}

func TestError_Error(t *testing.T) {
	assert := assert.New(t)

	err := Wrapf(errSomeError, "context %d", -44)

	expected := "context -44 (cause: some error)"
	assert.Equal(expected, err.Error())

	err = WrapCode(errSomeError, ErrInvalidUserName)

	errMsg := errorsCodeToMessage[ErrInvalidUserName]
	expected = fmt.Sprintf("(%d) %s (cause: some error)", ErrInvalidUserName, errMsg)
	assert.Equal(expected, err.Error())
}

func TestError_Code(t *testing.T) {
	assert := assert.New(t)

	err := NewCode(ErrInvalidUserName)

	impl, ok := err.(ErrorWithCode)
	assert.True(ok)
	assert.Equal(ErrInvalidUserName, impl.Code())
}

func TestError_MarshalJSON(t *testing.T) {
	assert := assert.New(t)

	err := New("haha")
	out, mErr := json.Marshal(err)

	expected := "{\"Message\":\"haha\"}"
	assert.Nil(mErr)
	assert.Equal(expected, string(out))

	err = NewCode(ErrInvalidUserName)
	out, mErr = json.Marshal(err)

	errMsg := errorsCodeToMessage[ErrInvalidUserName]
	expected = fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}", ErrInvalidUserName, errMsg)
	assert.Nil(mErr)
	assert.Equal(expected, string(out))

	err = Wrap(errSomeError, "hihi")
	out, mErr = json.Marshal(err)

	expected = "{\"Message\":\"hihi\",\"Cause\":\"some error\"}"
	assert.Nil(mErr)
	assert.Equal(expected, string(out))

	err = Wrap(New("haha"), "hihi")
	out, mErr = json.Marshal(err)

	expected = "{\"Message\":\"hihi\",\"Cause\":{\"Message\":\"haha\"}}"
	assert.Nil(mErr)
	assert.Equal(expected, string(out))
}

func TestError_IsErrorWithCode(t *testing.T) {
	assert := assert.New(t)

	assert.False(IsErrorWithCode(nil, ErrInvalidUserName))
	assert.False(IsErrorWithCode(errSomeError, ErrInvalidUserName))
	assert.True(IsErrorWithCode(NewCode(ErrInvalidUserName), ErrInvalidUserName))
	assert.False(IsErrorWithCode(NewCode(ErrInvalidPassword), ErrInvalidUserName))
	assert.False(IsErrorWithCode(New("haha"), ErrInvalidUserName))
}
