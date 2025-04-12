package errs

type ChatNotFoundError struct {
	err    error
	reason string
}

func NewChatNotFoundError(err error) *ChatNotFoundError {
	return NewChatNotFoundErrorWithReason(err, "chat_not_found")
}

func NewChatNotFoundErrorWithReason(err error, reason string) *ChatNotFoundError {
	return &ChatNotFoundError{
		err:    err,
		reason: reason,
	}
}

func (e *ChatNotFoundError) Error() string {
	return e.err.Error()
}

func (e *ChatNotFoundError) Reason() string {
	return e.reason
}

func (e *ChatNotFoundError) Unwrap() error {
	return e.err
}
