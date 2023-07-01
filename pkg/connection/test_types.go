package connection

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var errSomeError = fmt.Errorf("some error")

// https://forum.golangbridge.org/t/invalid-url-escape-while-parsing-url/8306
var invalidUrl = "http://test.com/Segment%%2815197306101420000%29.ts"

type mockHttpClient struct {
	inReq         *http.Request
	expectedResp  *http.Response
	expectedError error
}

func (mc *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	mc.inReq = req
	return mc.expectedResp, mc.expectedError
}

func generateHttpResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusAlreadyReported,
		Header: http.Header{
			"haha": []string{"gigi"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("some data"))),
	}
}

func nilRequestHttpBuilder(ri *requestImpl) (*http.Request, error) {
	return nil, nil
}

func errorHttpBuilder(ri *requestImpl) (*http.Request, error) {
	return &http.Request{}, errSomeError
}

type unmarshallableContent struct{}

func (uc unmarshallableContent) MarshalJSON() ([]byte, error) {
	return []byte{}, errSomeError
}
