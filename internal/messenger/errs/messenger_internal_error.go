package errs

type MessengerInternalError struct {
	err error
}

func NewMessengerInternalError(err error) *MessengerInternalError {
	return &MessengerInternalError{
		err: err,
	}
}

func (e *MessengerInternalError) Error() string {
	return e.err.Error()
}

func (e *MessengerInternalError) Reason() string {
	return "messenger_internal_error"
}

func (e *MessengerInternalError) Unwrap() error {
	return e.err
}
