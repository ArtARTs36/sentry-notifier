package errs

type Error interface {
	error
	Reason() string
	Unwrap() error
}
