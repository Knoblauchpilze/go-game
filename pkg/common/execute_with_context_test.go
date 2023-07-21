package common

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("someError")
var defaultSleep = 100 * time.Millisecond

func TestExecuteWithNoTimeout(t *testing.T) {
	assert := assert.New(t)

	err := executeWithNoTimeout(newProcess())
	assert.Nil(err)

	err = executeWithNoTimeout(newProcessWithError())
	assert.Equal(errDefault, err)
}

func TestExecuteWithContext_NoTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()

	err := ExecuteWithContext(newProcess(), ctx, 0)
	assert.Nil(err)

	err = ExecuteWithContext(newProcessWithError(), ctx, 0)
	assert.Equal(errDefault, err)
}

func TestExecuteWithTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := 2 * defaultSleep

	err := executeWithTimeout(newProcessWithSleep(), ctx, timeout)
	assert.Nil(err)

	err = executeWithTimeout(newProcessWithSleepAndError(), ctx, timeout)
	assert.Equal(errDefault, err)
}

func TestExecuteWithTimeout_TimeoutExceeded(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := defaultSleep / 2

	err := executeWithTimeout(newProcessWithSleep(), ctx, timeout)
	assert.Equal(context.DeadlineExceeded, err)

	err = executeWithTimeout(newProcessWithSleepAndError(), ctx, timeout)
	assert.Equal(context.DeadlineExceeded, err)
}

func TestExecuteWithTimeout_TimeoutExceeded_ExpectLog(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetLogFuncs)

	ctx := context.TODO()
	timeout := defaultSleep / 2
	var actual string
	errorLog = func(ctx context.Context, format string, args ...interface{}) {
		actual = format
	}

	err := executeWithTimeout(newProcessWithSleep(), ctx, timeout)
	// Safety to sleep at least the time it takes for the process to complete.
	time.Sleep(2 * timeout)
	assert.Equal(context.DeadlineExceeded, err)

	expected := "process didn't finish after "
	assert.Greater(len(actual), len(expected))
	assert.Equal(expected, actual[:len(expected)])
}

func TestExecuteWithTimeout_TimeoutExceeded_ExpectCleanUp(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	timeout := defaultSleep / 2

	p := newProcessWithSleep()
	var cleanUpCalled atomic.Int32
	p.CleanUpIfFailFunc = func() {
		cleanUpCalled.Add(1)
	}
	err := executeWithTimeout(p, ctx, timeout)
	// Safety to sleep at least the time it takes for the process to complete.
	time.Sleep(2 * timeout)
	assert.Equal(context.DeadlineExceeded, err)
	assert.Equal(int32(1), cleanUpCalled.Load())
}

func resetLogFuncs() {
	errorLog = logger.ScopedErrorf
	traceLog = logger.ScopedTracef
}

func newProcess() Process {
	return Process{
		WorkFunc: func() error {
			return nil
		},
	}
}

func newProcessWithError() Process {
	return Process{
		WorkFunc: func() error {
			return errDefault
		},
	}
}

func newProcessWithSleep() Process {
	return Process{
		WorkFunc: func() error {
			time.Sleep(defaultSleep)
			return nil
		},
	}
}

func newProcessWithSleepAndError() Process {
	return Process{
		WorkFunc: func() error {
			time.Sleep(defaultSleep)
			return errDefault
		},
	}
}
