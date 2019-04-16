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

func newError(message, code string) Error {
	return Error{Message: message, Code: code}
}

func WrapOpsErr(err error, msg string) Error {
	wrappedErrMsg := errors.Wrap(err, msg).Error()
	if awsErr, ok := err.(awserr.Error); ok {
		return newError(wrappedErrMsg, awsErr.Code())
	}

	return newError(wrappedErrMsg, UnknownOpsErrCode)
}

func WrapErr(err error, code, msg string) Error {
	wrappedErrMsg := errors.Wrap(err, msg).Error()
	return newError(wrappedErrMsg, code)
}

func AnyEquals(origCode string, errCodes ...string) bool {
	for _, ec := range errCodes {
		if origCode == ec {
			return true
		}
	}

	return false
}
