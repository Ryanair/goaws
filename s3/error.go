package s3

const (
	SigningURLErrCode = "SigningURLErr"
)

type Error struct {
	Message string
	Code    string
}

func newError(message, code string) Error {
	return Error{Message: message, Code: code}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) SigningFailed() bool {
	return e.Code == SigningURLErrCode
}
