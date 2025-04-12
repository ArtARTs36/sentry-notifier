package errs

type InvalidCredentialsError struct {
	err error
}

func NewInvalidCredentialsError(err error) *InvalidCredentialsError {
	return &InvalidCredentialsError{
		err: err,
	}
}

func (e *InvalidCredentialsError) Error() string {
	return e.err.Error()
}

func (e *InvalidCredentialsError) Reason() string {
	return "invalid_credentials"
}

func (e *InvalidCredentialsError) Unwrap() error {
	return e.err
}
