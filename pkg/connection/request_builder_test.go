package connection

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBuilder_InvalidReq(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()

	rw, err := rb.Build()
	assert.Nil(rw)
	assert.NotNil(err)

	rb.SetUrl("haha")
	rw, err = rb.Build()
	assert.Nil(rw)
	assert.NotNil(err)
}

func TestRequestBuilder_SetUrl(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()
	url := "haha"
	rb.SetUrl(url)
	assert.Equal(url, rb.url)

	url = "haha2"
	rb.SetUrl(url)
	assert.Equal(url, rb.url)
}

func TestRequestBuilder_SetHeaders(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()
	headers := http.Header{
		"hihi": []string{"haha"},
	}
	rb.SetHeaders(headers)
	assert.Equal(headers, rb.headers)

	headers = http.Header{
		"jiji": []string{"koko"},
	}
	rb.SetHeaders(headers)
	assert.Equal(headers, rb.headers)
}

func TestRequestBuilder_AddHeader(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()

	rb.AddHeader("haha", []string{"hihi"})
	exp := http.Header{
		"haha": []string{"hihi"},
	}
	assert.Equal(exp, rb.headers)

	rb.AddHeader("hihi", []string{"jiji"})
	exp = http.Header{
		"haha": []string{"hihi"},
		"hihi": []string{"jiji"},
	}
	assert.Equal(exp, rb.headers)

	rb.AddHeader("haha", []string{"jojo", "juju"})
	exp = http.Header{
		"haha": []string{"hihi", "jojo", "juju"},
		"hihi": []string{"jiji"},
	}
	assert.Equal(exp, rb.headers)
}

func TestRequestBuilder_SetBody(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()
	rb.SetBody("haha", 32)

	exp := http.Header{
		"Content-Type": []string{"haha"},
	}
	assert.Equal(exp, rb.headers)
	assert.Equal(32, rb.body)
}

func TestRequestBuilder_SetBodyWithHeaders(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()
	rb.SetHeaders(http.Header{
		"hihi": []string{"jojo"},
	})
	rb.SetBody("haha", 32)

	exp := http.Header{
		"hihi":         []string{"jojo"},
		"Content-Type": []string{"haha"},
	}
	assert.Equal(exp, rb.headers)
	assert.Equal(32, rb.body)

	rb.SetHeaders(http.Header{
		"Content-Type": []string{"hihi"},
	})
	rb.SetBody("koko", 26)

	exp = http.Header{
		"Content-Type": []string{"hihi", "koko"},
	}
	assert.Equal(exp, rb.headers)
	assert.Equal(26, rb.body)
}

func TestRequestBuilder_Success(t *testing.T) {
	assert := assert.New(t)

	rb := newRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpRequestBuilder(buildGetRequest)

	rw, err := rb.Build()
	assert.NotNil(rw)
	assert.Nil(err)
}
