package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerData_Host_NoReq(t *testing.T) {
	assert := assert.New(t)

	cd := callerData{}
	assert.Equal("N/A", cd.host())
}

func TestCallerData_Host(t *testing.T) {
	assert := assert.New(t)

	cd := callerData{
		req: &http.Request{
			Host:       "someHost",
			RequestURI: "someUri",
		},
	}
	assert.Equal("someHostsomeUri", cd.host())
}
