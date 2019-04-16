package dynamodb

import (
	"github.com/Ryanair/goaws/internal"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	MarshalErrCode            = "DynamoDBMarshalErr"
	UnmarshalErrCode          = "DynamoDBUnmarshalErr"
	InvalidConditionErrCode   = "DynamoDBInvalidConditionErr"
	ValidationErrCode         = "ValidationException"
	ThrottlingErrCode         = "ThrottlingException"
	UnrecognizedClientErrCode = "UnrecognizedClientException"
)

type Error internal.Error

func (e Error) Error() string {
	return e.Message
}

func wrapErr(err error, code, msg string) error {
	return Error(internal.WrapErr(err, code, msg))
}

func wrapOpsErr(err error, msg string) error {
	return Error(internal.WrapOpsErr(err, msg))
}

func (e Error) MarshallingFailed() bool {
	return internal.AnyEquals(e.Code, MarshalErrCode)
}

func (e Error) UnmarshallingFailed() bool {
	return internal.AnyEquals(e.Code, UnmarshalErrCode)
}

func (e Error) InvalidCondition() bool {
	return internal.AnyEquals(e.Code, InvalidConditionErrCode)
}

func (e Error) ValidationFailed() bool {
	return internal.AnyEquals(e.Code, ValidationErrCode)
}

func (e Error) ConditionFailed() bool {
	return internal.AnyEquals(e.Code, dynamodb.ErrCodeConditionalCheckFailedException)
}

func (e Error) BackupUnavailable() bool {
	return internal.AnyEquals(e.Code,
		dynamodb.ErrCodeContinuousBackupsUnavailableException,
		dynamodb.ErrCodePointInTimeRecoveryUnavailableException)
}

func (e Error) InternalError() bool {
	return internal.AnyEquals(e.Code, dynamodb.ErrCodeInternalServerError)
}

func (e Error) ResourceNotFound() bool {
	return internal.AnyEquals(e.Code,
		dynamodb.ErrCodeResourceNotFoundException,
		dynamodb.ErrCodeBackupNotFoundException,
		dynamodb.ErrCodeGlobalTableNotFoundException,
		dynamodb.ErrCodeIndexNotFoundException,
		dynamodb.ErrCodeReplicaNotFoundException,
		dynamodb.ErrCodeTableNotFoundException)
}

func (e Error) ResourceAlreadyExists() bool {
	return internal.AnyEquals(e.Code,
		dynamodb.ErrCodeGlobalTableAlreadyExistsException,
		dynamodb.ErrCodeReplicaAlreadyExistsException,
		dynamodb.ErrCodeTableAlreadyExistsException)
}

func (e Error) InvalidOperation() bool {
	return internal.AnyEquals(e.Code,
		dynamodb.ErrCodeIdempotentParameterMismatchException,
		dynamodb.ErrCodeInvalidRestoreTimeException,
		dynamodb.ErrCodeTransactionCanceledException,
		dynamodb.ErrCodeTransactionConflictException,
		dynamodb.ErrCodeTransactionInProgressException)
}

func (e Error) LimitExceeded() bool {
	return internal.AnyEquals(e.Code,
		dynamodb.ErrCodeItemCollectionSizeLimitExceededException,
		dynamodb.ErrCodeLimitExceededException,
		dynamodb.ErrCodeRequestLimitExceeded,
		dynamodb.ErrCodeProvisionedThroughputExceededException)
}

func (e Error) ResourceInUse() bool {
	return internal.AnyEquals(e.Code,
		dynamodb.ErrCodeBackupInUseException,
		dynamodb.ErrCodeResourceInUseException,
		dynamodb.ErrCodeTableInUseException)
}

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Programming.Errors.html
func (e Error) Retryable() bool {
	return internal.AnyEquals(e.Code,
		ThrottlingErrCode,
		UnrecognizedClientErrCode,
		dynamodb.ErrCodeItemCollectionSizeLimitExceededException,
		dynamodb.ErrCodeLimitExceededException,
		dynamodb.ErrCodeProvisionedThroughputExceededException,
		dynamodb.ErrCodeRequestLimitExceeded,
		dynamodb.ErrCodeInternalServerError)
}
