package common

type ExecutionResult struct {
	ExecutionErr error
	ProcessErr   error
}

func (r ExecutionResult) Error() error {
	if r.ExecutionErr != nil {
		return r.ExecutionErr
	}
	return r.ProcessErr
}
