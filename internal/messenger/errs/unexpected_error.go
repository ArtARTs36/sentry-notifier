package errs

type UnexpectedError struct {
	err error
}

func NewUnexpectedError(err error) *UnexpectedError {
	return &UnexpectedError{
		err: err,
	}
}

func (e *UnexpectedError) Error() string {
	return e.err.Error()
}

func (e *UnexpectedError) Reason() string {
	return "unexpected_error"
}

func (e *UnexpectedError) Unwrap() error {
	return e.err
}
