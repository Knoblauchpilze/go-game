package rest

import (
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func generateRequestWithHeader() http.Request {
	return http.Request{
		Header: make(http.Header),
	}
}

func TestGetHeaderFromHttpRequest_NoHeader(t *testing.T) {
	assert := assert.New(t)

	req := http.Request{}

	_, err := GetHeaderFromHttpRequest(&req, "foo")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchHeader))
}

func TestGetHeaderFromHttpRequest_OneValue(t *testing.T) {
	assert := assert.New(t)

	headerValues := []string{"haha"}

	req := generateRequestWithHeader()
	req.Header["foo"] = headerValues

	out, err := GetHeaderFromHttpRequest(&req, "foo")
	assert.Nil(err)
	assert.Equal(len(headerValues), len(out))

	for id, expectedHeader := range headerValues {
		assert.Equal(expectedHeader, out[id])
	}
}

func TestGetHeaderFromHttpRequest_AnotherValue(t *testing.T) {
	assert := assert.New(t)

	req := generateRequestWithHeader()
	req.Header["food"] = []string{"haha"}

	_, err := GetHeaderFromHttpRequest(&req, "foo")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchHeader))
}

func TestGetHeaderFromHttpRequest_TwoValues(t *testing.T) {
	assert := assert.New(t)

	fooHeaderValues := []string{"haha"}
	barHeaderValues := []string{"hihi"}

	req := generateRequestWithHeader()
	req.Header["foo"] = fooHeaderValues
	req.Header["bar"] = barHeaderValues

	out, err := GetHeaderFromHttpRequest(&req, "foo")
	assert.Nil(err)
	assert.Equal(len(fooHeaderValues), len(out))

	for id, expectedHeader := range fooHeaderValues {
		assert.Equal(expectedHeader, out[id])
	}
}

func TestGetSingleHeaderFromHttpRequest_NoHeader(t *testing.T) {
	assert := assert.New(t)

	req := http.Request{}

	_, err := GetSingleHeaderFromHttpRequest(&req, "foo")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchHeader))
}

func TestGetSingleHeaderFromHttpRequest_OneValue(t *testing.T) {
	assert := assert.New(t)

	headerValues := []string{"haha"}

	req := generateRequestWithHeader()
	req.Header["foo"] = headerValues

	out, err := GetSingleHeaderFromHttpRequest(&req, "foo")
	assert.Nil(err)
	assert.Equal("haha", out)
}

func TestGetSingleHeaderFromHttpRequest_MultipleValues(t *testing.T) {
	assert := assert.New(t)

	headerValues := []string{"haha", "hihi"}

	req := generateRequestWithHeader()
	req.Header["foo"] = headerValues

	_, err := GetSingleHeaderFromHttpRequest(&req, "foo")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNonUniqueHeader))
}
