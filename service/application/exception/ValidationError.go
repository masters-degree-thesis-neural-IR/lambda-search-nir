package exception

type ValidationError struct {
	StatusCode int
	Message    string
}

func ThrowValidationError(message string) error {

	var err error = ValidationError{
		StatusCode: 400,
		Message:    message,
	}
	return err
}

func (e ValidationError) Error() string {
	return e.Message
}
