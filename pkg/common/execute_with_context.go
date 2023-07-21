package common

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
)

type Process struct {
	WorkFunc          func() error
	CleanUpIfFailFunc func()
}

var errorLog = logger.ScopedErrorf
var traceLog = logger.ScopedTracef

func ExecuteWithContext(process Process, ctx context.Context, timeout time.Duration) error {
	if timeout == 0 {
		return executeWithNoTimeout(process)
	}

	return executeWithTimeout(process, ctx, timeout)
}

func executeWithNoTimeout(process Process) error {
	return process.WorkFunc()
}

func executeWithTimeout(process Process, ctx context.Context, timeout time.Duration) error {
	// https://medium.com/geekculture/timeout-context-in-go-e88af0abd08d
	decoratedCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	processDone := make(chan bool, 1)
	start := time.Now()
	var cancelled atomic.Bool

	// https://go.dev/doc/articles/race_detector
	var errP error
	go func() {
		errP = process.WorkFunc()
		processDone <- true
		if cancelled.Load() {
			if process.CleanUpIfFailFunc != nil {
				process.CleanUpIfFailFunc()
			}
		}
	}()

	var err error
	select {
	case <-decoratedCtx.Done():
		err = decoratedCtx.Err()
		cancelled.Store(true)
		errorLog(ctx, "process didn't finish after %+v (err: %+v)", timeout, err)
	case <-processDone:
		err = errP
		traceLog(ctx, "executed process after %+v", time.Since(start))
	}

	return err
}
