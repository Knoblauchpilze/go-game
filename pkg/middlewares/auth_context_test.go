package middlewares

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseAuthenticationHeader_InvalidHeader(t *testing.T) {
	assert := assert.New(t)

	genErr := errors.New(genErrMsg)

	authData := ""
	_, err := parseAuthenticationHeader(authData)
	assert.Equal(err, genErr)

	authData = "toto tata"
	_, err = parseAuthenticationHeader(authData)
	assert.Equal(err, genErr)
}

func TestParseAuthenticationHeader_InvalidAuthenticationScheme(t *testing.T) {
	assert := assert.New(t)

	genErr := errors.New(genErrMsg)

	authData := "wrong-scheme haha hihi"
	_, err := parseAuthenticationHeader(authData)
	assert.Equal(err, genErr)
}

func TestParseAuthenticationHeader_InvalidPropSyntax(t *testing.T) {
	assert := assert.New(t)

	id := "74b02b94-ca64-11ed-ba48-18c04d0e6a41"

	authData := "bearer haha hihi"
	_, err := parseAuthenticationHeader(authData)
	expectedErr := fmt.Sprintf("%s (cause: Ill-formed prop in Authorization header: \"haha\")", genErrMsg)
	assert.Equal(err.Error(), expectedErr)

	authData = "bearer haha= hihi"
	_, err = parseAuthenticationHeader(authData)
	expectedErr = fmt.Sprintf("%s (cause: Ill-formed prop in Authorization header: \"haha=\")", genErrMsg)
	assert.Equal(err.Error(), expectedErr)

	authData = "bearer =haha hihi"
	_, err = parseAuthenticationHeader(authData)
	expectedErr = fmt.Sprintf("%s (cause: Ill-formed prop in Authorization header: \"=haha\")", genErrMsg)
	assert.Equal(err.Error(), expectedErr)

	authData = "bearer haha=22 hihi=45"
	_, err = parseAuthenticationHeader(authData)
	expectedErr = fmt.Sprintf("%s (cause: no \"%s\" key in Authorization header)", genErrMsg, authenticationUserKey)
	assert.Equal(err.Error(), expectedErr)

	authData = "bearer user=22 hihi=45"
	_, err = parseAuthenticationHeader(authData)
	expectedErr = fmt.Sprintf("%s (cause: failed to parse user id \"22\" (cause: invalid UUID length: 2))", genErrMsg)
	assert.Equal(err.Error(), expectedErr)

	authData = fmt.Sprintf("bearer user=%s hihi=45", id)
	_, err = parseAuthenticationHeader(authData)
	expectedErr = fmt.Sprintf("%s (cause: no \"%s\" key in Authorization header)", genErrMsg, authenticationTokenKey)
	assert.Equal(err.Error(), expectedErr)
}

func TestParseAuthenticationHeader_Success(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse("74b02b94-ca64-11ed-ba48-18c04d0e6a41")
	val := "dummy-token"

	authData := fmt.Sprintf("bearer user=%s token=%s", id, val)
	token, err := parseAuthenticationHeader(authData)
	assert.Nil(err)
	assert.Equal(token.User, id)
	assert.Equal(token.Value, val)
}

func TestGenerateAuthenticationContext_NoRequestContext(t *testing.T) {
	assert := assert.New(t)

	s := newMockServer("haha")
	s.withRequestCtx = false
	s.call()
	assert.Equal(http.StatusInternalServerError, s.mrw.code)
}

func TestGenerateAuthenticationContext_NoHeader(t *testing.T) {
	assert := assert.New(t)

	s := newMockServer("haha")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)
}

func TestGenerateAuthenticationContext_InvalidHeader(t *testing.T) {
	assert := assert.New(t)

	s := newMockServer("haha")

	s.withAuthorization("")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("haha")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("haha hihi")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("wrong-scheme haha hihi")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("bearer haha hihi")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("bearer haha= hihi")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("bearer =haha hihi")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("bearer haha=22 hihi=45")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	s.withAuthorization("bearer user=22 hihi=45")
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)

	id := "74b02b94-ca64-11ed-ba48-18c04d0e6a41"
	hdr := fmt.Sprintf("bearer user=%s hihi=45", id)
	s.withAuthorization(hdr)
	s.call()
	assert.Equal(http.StatusBadRequest, s.mrw.code)
}

func TestGenerateAuthenticationContext_NoTokenForUser(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse("74b02b94-ca64-11ed-ba48-18c04d0e6a41")
	s := newMockServer("haha")
	s.auth.isError = true

	hdr := fmt.Sprintf("bearer user=%s token=45", id)
	s.withAuthorization(hdr)
	s.call()
	assert.Equal(http.StatusUnauthorized, s.mrw.code)
	resp, err := unmarshalExpectedResponseBody(s.mrw.data)
	assert.Nil(err)
	assert.Equal("ERROR", resp.Status)
	assert.Contains(string(resp.Details), "authentication failure")
}

func TestGenerateAuthenticationContext_InvalidToken(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse("74b02b94-ca64-11ed-ba48-18c04d0e6a41")
	pwd := "hoho"
	s := newMockServer("haha")
	token, err := s.auth.GenerateToken(id, pwd)
	assert.Nil(err)

	hdr := fmt.Sprintf("bearer user=%s token=%s1", id, token.Value)
	s.withAuthorization(hdr)
	s.call()
	assert.Equal(http.StatusUnauthorized, s.mrw.code)
	resp, err := unmarshalExpectedResponseBody(s.mrw.data)
	assert.Nil(err)
	assert.Equal("ERROR", resp.Status)
	assert.Contains(string(resp.Details), "authentication failure")
}

func TestGenerateAuthenticationContext_ExpiredToken(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse("74b02b94-ca64-11ed-ba48-18c04d0e6a41")
	pwd := "hoho"
	s := newMockServer("haha")
	token, err := s.auth.GenerateToken(id, pwd)
	assert.Nil(err)

	hdr := fmt.Sprintf("bearer user=%s token=%s", id, token.Value)
	s.withAuthorization(hdr)
	s.call()
	assert.Equal(http.StatusUnauthorized, s.mrw.code)
	resp, err := unmarshalExpectedResponseBody(s.mrw.data)
	assert.Nil(err)
	assert.Equal("ERROR", resp.Status)
	assert.Contains(string(resp.Details), "authentication expired")

	fmt.Printf("%s\n", string(resp.Details))
}

func TestGenerateAuthenticationContext_Success(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse("74b02b94-ca64-11ed-ba48-18c04d0e6a41")
	pwd := "hoho"
	s := newMockServer("haha")
	s.auth.expiration = time.Now().Add(1 * time.Minute)
	token, err := s.auth.GenerateToken(id, pwd)
	assert.Nil(err)

	hdr := fmt.Sprintf("bearer user=%s token=%s", id, token.Value)
	s.withAuthorization(hdr)
	s.call()
	assert.Equal(http.StatusOK, s.mrw.code)
}
