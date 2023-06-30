package connection

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestWrapper interface {
	Perform() (*http.Response, error)
}
