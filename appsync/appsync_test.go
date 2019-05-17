package appsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	errMessage = "internal server error"
	errType    = ErrorTypeInternal
)

func TestError_NewError_OK(t *testing.T) {

	// when
	err := NewError(errType, errMessage)

	// then
	assert.Equal(t, errMessage, err.Message)
	assert.Equal(t, errType, err.ErrorType)
}

func TestError_SetInfo_OK(t *testing.T) {

	// given
	errInfo := "connection error"

	// when
	err := NewError(errType, errMessage, Info(errInfo))

	// then
	infoStr, ok := err.Info.(string)
	if !ok {
		t.Errorf("invalid info type, expected: string, actual: %T", err.Info)
	}
	assert.Equal(t, errInfo, infoStr)
}
