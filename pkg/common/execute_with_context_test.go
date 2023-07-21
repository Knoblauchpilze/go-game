package common

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithNoTimeout(t *testing.T) {
	assert := assert.New(t)

	p := func() error { return nil }
	res := executeWithNoTimeout(p)
	assert.Nil(res.ExecutionErr)
	assert.Nil(res.ProcessErr)

	p = func() error { return fmt.Errorf("someError") }
	res = executeWithNoTimeout(p)
	assert.Nil(res.ExecutionErr)
	assert.Equal("someError", res.ProcessErr.Error())
}

func TestExecuteWithContext_NoTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()

	p := func() error { return nil }
	res := ExecuteWithContext(p, ctx, 0)
	assert.Nil(res.ExecutionErr)
	assert.Nil(res.ProcessErr)

	p = func() error { return fmt.Errorf("someError") }
	res = ExecuteWithContext(p, ctx, 0)
	assert.Nil(res.ExecutionErr)
	assert.Equal("someError", res.ProcessErr.Error())
}

func TestExecuteWithTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := 100 * time.Millisecond

	p := func() error {
		time.Sleep(timeout / 2)
		return nil
	}
	res := executeWithTimeout(p, ctx, timeout)
	assert.Nil(res.ExecutionErr)
	assert.Nil(res.ProcessErr)

	p = func() error {
		time.Sleep(timeout / 2)
		return fmt.Errorf("someError")
	}
	res = executeWithTimeout(p, ctx, timeout)
	assert.Nil(res.ExecutionErr)
	assert.Equal("someError", res.ProcessErr.Error())
}

func TestExecuteWithTimeout_TimeoutExceeded(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := 100 * time.Millisecond

	p := func() error {
		time.Sleep(2 * timeout)
		return nil
	}
	res := executeWithTimeout(p, ctx, timeout)
	assert.Equal(context.DeadlineExceeded, res.ExecutionErr)
	assert.Nil(res.ProcessErr)

	p = func() error {
		time.Sleep(2 * timeout)
		return fmt.Errorf("someError")
	}
	res = executeWithTimeout(p, ctx, timeout)
	assert.Equal(context.DeadlineExceeded, res.ExecutionErr)
	assert.Nil(res.ProcessErr)
}

func TestExecuteWithTimeout_TimeoutExceeded_ExpectLog(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetLogFuncs)

	ctx := context.TODO()
	timeout := 100 * time.Millisecond
	var actual string
	errorLog = func(ctx context.Context, format string, args ...interface{}) {
		actual = format
	}

	p := func() error {
		time.Sleep(2 * timeout)
		return nil
	}
	res := executeWithTimeout(p, ctx, timeout)
	// Safety to sleep at least the time it takes for the process to complete.
	time.Sleep(2 * timeout)
	assert.Equal(context.DeadlineExceeded, res.ExecutionErr)
	assert.Nil(res.ProcessErr)

	expected := "process finished after "
	assert.Greater(len(actual), len(expected))
	assert.Equal(expected, actual[:len(expected)])
}

func resetLogFuncs() {
	errorLog = logger.ScopedErrorf
	traceLog = logger.ScopedTracef
}
