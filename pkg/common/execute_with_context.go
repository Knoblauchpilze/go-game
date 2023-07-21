package common

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
)

type Process func() error

var errorLog = logger.ScopedErrorf
var traceLog = logger.ScopedTracef

func ExecuteWithContext(process Process, ctx context.Context, timeout time.Duration) ExecutionResult {
	if timeout == 0 {
		return executeWithNoTimeout(process)
	}

	return executeWithTimeout(process, ctx, timeout)
}

func executeWithNoTimeout(process Process) ExecutionResult {
	var out ExecutionResult
	out.ProcessErr = process()
	return out
}

func executeWithTimeout(process Process, ctx context.Context, timeout time.Duration) ExecutionResult {
	var out ExecutionResult

	decoratedCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	processDone := make(chan bool, 1)
	shouldLog := false
	start := time.Now()

	go func() {
		out.ProcessErr = process()
		processDone <- true
		if shouldLog {
			errorLog(ctx, "process finished after %+v (greater than %+v available)", time.Since(start), timeout)
		}
	}()

	select {
	case <-decoratedCtx.Done():
		out.ExecutionErr = decoratedCtx.Err()
		shouldLog = true
		errorLog(ctx, "processed didn't finish after %+v (err: %+v)", timeout, out.ExecutionErr)
	case <-processDone:
		traceLog(ctx, "executed process after %+v", time.Since(start))
	}

	return out
}
