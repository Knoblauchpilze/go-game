package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/common"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var dummyId = uuid.MustParse("eb10f542-c2a8-11ed-befe-18c04d0e6a41")

func TestResponseBuilder_New(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetIdUnwrapperFunc)

	out := newResponseBuilder(context.TODO())
	assert.Equal(http.StatusOK, out.code)
	assert.Equal(statusSuccess, out.status)
	assert.Nil(out.details)
}

func TestResponseBuilder_New_UseIdFromContext(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetIdUnwrapperFunc)

	idUnwrapper = func(ctx context.Context) (uuid.UUID, bool) {
		return dummyId, true
	}

	b := newResponseBuilder(context.TODO())
	assert.Equal(dummyId, b.id)
}

func TestResponseBuilder_New_GenerateId(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetIdUnwrapperFunc)

	idUnwrapper = func(ctx context.Context) (uuid.UUID, bool) {
		return dummyId, false
	}

	b := newResponseBuilder(context.TODO())
	assert.NotEqual(dummyId, b.id)
}

func TestResponseBuilder_SetCode(t *testing.T) {
	assert := assert.New(t)

	b := newResponseBuilder(context.TODO())

	b.setCode(http.StatusBadGateway)
	assert.Equal(http.StatusBadGateway, b.code)
	assert.Equal(statusFailure, b.status)

	b.setCode(http.StatusOK)
	assert.Equal(http.StatusOK, b.code)
	assert.Equal(statusSuccess, b.status)
}

type unmarshallableContent struct{}

func (uc unmarshallableContent) MarshalJSON() ([]byte, error) {
	return []byte{}, fmt.Errorf("some error")
}

func TestResponseBuilder_SetDetails(t *testing.T) {
	assert := assert.New(t)

	b := newResponseBuilder(context.TODO())

	b.setDetails(32)
	assert.Equal("32", string(b.details))
}

func TestResponseBuilder_SetDetails_Unmarshallable(t *testing.T) {
	assert := assert.New(t)

	b := newResponseBuilder(context.TODO())

	b.setDetails(unmarshallableContent{})
	assert.True(common.IsInterfaceNil(b.details))
}

func TestResponseBuilder_Serialize(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetIdUnwrapperFunc)

	idUnwrapper = func(ctx context.Context) (uuid.UUID, bool) {
		return dummyId, true
	}
	b := newResponseBuilder(context.TODO())

	out, err := b.serialize()
	assert.Nil(err)
	expected := "{\"RequestId\":\"eb10f542-c2a8-11ed-befe-18c04d0e6a41\",\"Status\":\"SUCCESS\"}"
	assert.Equal(expected, string(out))
}

func TestResponseBuilder_Serialize_WithDetails(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetIdUnwrapperFunc)

	idUnwrapper = func(ctx context.Context) (uuid.UUID, bool) {
		return dummyId, true
	}
	b := newResponseBuilder(context.TODO())
	b.setDetails(32)

	out, err := b.serialize()
	assert.Nil(err)
	expected := "{\"RequestId\":\"eb10f542-c2a8-11ed-befe-18c04d0e6a41\",\"Status\":\"SUCCESS\",\"Details\":32}"
	assert.Equal(expected, string(out))
}

func TestResponseBuilder_Write(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetIdUnwrapperFunc)

	idUnwrapper = func(ctx context.Context) (uuid.UUID, bool) {
		return dummyId, true
	}
	b := newResponseBuilder(context.TODO())
	b.setCode(http.StatusConflict)
	b.setDetails(32)

	m := mockHttpResponseWriter{}
	b.Write(&m)
	expected := "{\"RequestId\":\"eb10f542-c2a8-11ed-befe-18c04d0e6a41\",\"Status\":\"ERROR\",\"Details\":32}"
	assert.Equal(expected, string(m.written))
	assert.Equal(http.StatusConflict, m.code)
}

func resetIdUnwrapperFunc() {
	idUnwrapper = logger.UnwrapIdFromContext
}

type mockHttpResponseWriter struct {
	code    int
	written []byte
}

func (m *mockHttpResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockHttpResponseWriter) Write(in []byte) (int, error) {
	m.written = in
	return 0, nil
}

func (m *mockHttpResponseWriter) WriteHeader(statusCode int) {
	m.code = statusCode
}
