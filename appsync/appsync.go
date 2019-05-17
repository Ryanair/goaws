package appsync

const (
	ErrorTypeValidation ErrorType = "Validation"
	ErrorTypeConflict   ErrorType = "Conflict"
	ErrorTypeInternal   ErrorType = "Internal"
	ErrorTypeNotFound   ErrorType = "NotFound"
)

type ErrorType string

type Error struct {
	ErrorType ErrorType   `json:"errorType"`
	Message   string      `json:"message"`
	Info      interface{} `json:"info"`
}

func NewError(errType ErrorType, msg string, ops ...func(*Error)) *Error {

	err := &Error{
		ErrorType: errType,
		Message:   msg,
	}

	for _, op := range ops {
		op(err)
	}

	return err
}

func Info(info interface{}) func(*Error) {
	return func(err *Error) {
		err.Info = info
	}
}
