package messenger

type PingError struct {
	Reason string
	Err    error
}

func (e *PingError) Error() string {
	return e.Err.Error()
}

func invalidCredentialsPingError(err error) error {
	return &PingError{
		Reason: "invalid_credentials",
		Err:    err,
	}
}

func chatNotFoundPingError(err error) error {
	return &PingError{
		Reason: "chat_not_found",
		Err:    err,
	}
}
