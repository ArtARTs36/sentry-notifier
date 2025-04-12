package errs

type NetworkError struct {
	err error
}

func NewNetworkError(err error) *NetworkError {
	return &NetworkError{
		err: err,
	}
}

func (e *NetworkError) Error() string {
	return e.err.Error()
}

func (e *NetworkError) Reason() string {
	return "network_error"
}

func (e *NetworkError) Unwrap() error {
	return e.err
}
