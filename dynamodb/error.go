package dynamodb

import (
	"github.com/Ryanair/goaws/internal"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	ErrCodeMarshal            = "DynamoDBMarshalErr"
	ErrCodeUnmarshal          = "DynamoDBUnmarshalErr"
	ErrCodeInvalidCondition   = "DynamoDBInvalidConditionErr"
	ErrCodeValidation         = "ValidationException"
	ErrCodeThrottling         = "ThrottlingException"
	ErrCodeUnrecognizedClient = "UnrecognizedClientException"
)

type Error internal.Error

func (e Error) Error() string {
	return e.Message
}

func (e Error) Cause() error {
	return e.Causer
}

func wrapErr(err error, msg string) error {
	return Error(internal.WrapErr(err, msg))
}

func wrapErrWithCode(err error, msg, code string) error {
	return Error(internal.WrapErrWithCode(err, msg, code))
}

func (e Error) MarshallingFailed() bool {
	return internal.AnyEquals(e.Code, ErrCodeMarshal)
}

func (e Error) UnmarshallingFailed() bool {
	return internal.AnyEquals(e.Code, ErrCodeUnmarshal)
}

func (e Error) InvalidCondition() bool {
	return internal.AnyEquals(e.Code, ErrCodeInvalidCondition)
}

func (e Error) ValidationFailed() bool {
	return internal.AnyEquals(e.Code, ErrCodeValidation)
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
		ErrCodeThrottling,
		ErrCodeUnrecognizedClient,
		dynamodb.ErrCodeItemCollectionSizeLimitExceededException,
		dynamodb.ErrCodeLimitExceededException,
		dynamodb.ErrCodeProvisionedThroughputExceededException,
		dynamodb.ErrCodeRequestLimitExceeded,
		dynamodb.ErrCodeInternalServerError)
}
