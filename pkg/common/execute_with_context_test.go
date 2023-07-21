package common

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("someError")

func TestExecuteWithNoTimeout(t *testing.T) {
	assert := assert.New(t)

	p := func() error { return nil }
	err := executeWithNoTimeout(p)
	assert.Nil(err)

	p = func() error { return errDefault }
	err = executeWithNoTimeout(p)
	assert.Equal(errDefault, err)
}

func TestExecuteWithContext_NoTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()

	p := func() error { return nil }
	err := ExecuteWithContext(p, ctx, 0)
	assert.Nil(err)

	p = func() error { return errDefault }
	err = ExecuteWithContext(p, ctx, 0)
	assert.Equal(errDefault, err)
}

func TestExecuteWithTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := 100 * time.Millisecond

	p := func() error {
		time.Sleep(timeout / 2)
		return nil
	}
	err := executeWithTimeout(p, ctx, timeout)
	assert.Nil(err)

	p = func() error {
		time.Sleep(timeout / 2)
		return fmt.Errorf("someError")
	}
	err = executeWithTimeout(p, ctx, timeout)
	assert.Equal(errDefault, err)
}

func TestExecuteWithTimeout_TimeoutExceeded(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := 100 * time.Millisecond

	p := func() error {
		time.Sleep(2 * timeout)
		return nil
	}
	err := executeWithTimeout(p, ctx, timeout)
	assert.Equal(context.DeadlineExceeded, err)

	p = func() error {
		time.Sleep(2 * timeout)
		return fmt.Errorf("someError")
	}
	err = executeWithTimeout(p, ctx, timeout)
	assert.Equal(context.DeadlineExceeded, err)
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
	err := executeWithTimeout(p, ctx, timeout)
	// Safety to sleep at least the time it takes for the process to complete.
	time.Sleep(2 * timeout)
	assert.Equal(context.DeadlineExceeded, err)

	expected := "process didn't finish after "
	assert.Greater(len(actual), len(expected))
	assert.Equal(expected, actual[:len(expected)])
}

func resetLogFuncs() {
	errorLog = logger.ScopedErrorf
	traceLog = logger.ScopedTracef
}
