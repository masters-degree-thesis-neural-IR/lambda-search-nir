package exception

type UnexpectedError struct {
	StatusCode int
	Message    string
}

func ThrowUnexpectedError(message string) error {

	var err error = UnexpectedError{
		StatusCode: 500,
		Message:    message,
	}
	return err
}

func (e UnexpectedError) Error() string {
	return e.Message
}
