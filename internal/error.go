package internal

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
)

const (
	UnknownOpsErrCode = "UnknownOpsErr"
)

type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(message, code string) error {
	return &Error{Message: message, Code: code}
}

// goErr -> awsErr -> Error -> daveErr -> Error

func WrapOpsErr(err error, msg string) error {
	wrappedErrMsg := errors.Wrap(err, msg).Error()

	switch err := errors.Cause(err).(type) {
	case awserr.Error:
		return NewError(wrappedErrMsg, err.Code())
	case *Error:
		return NewError(wrappedErrMsg, err.Code)
	default:
		return NewError(wrappedErrMsg, UnknownOpsErrCode)
	}
}

func WrapErr(err error, code, msg string) error {
	wrappedErrMsg := errors.Wrap(err, msg).Error()
	return NewError(wrappedErrMsg, code)
}

func AnyEquals(origCode string, errCodes ...string) bool {
	for _, ec := range errCodes {
		if origCode == ec {
			return true
		}
	}

	return false
}
