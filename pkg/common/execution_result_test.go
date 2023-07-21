package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("someError")

func TestExecutionResult_Error(t *testing.T) {
	assert := assert.New(t)

	res := ExecutionResult{}
	assert.Nil(res.Error())

	res.ExecutionErr = errDefault
	assert.Equal(errDefault, res.Error())

	someOtherError := fmt.Errorf("error2")
	res.ExecutionErr = nil
	res.ProcessErr = someOtherError
	assert.Equal(someOtherError, res.Error())
}
