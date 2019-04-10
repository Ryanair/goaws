package dynamodb

import "github.com/aws/aws-sdk-go/service/dynamodb"

const (
	MarshalErrCode            = "DynamoDBMarshalErr"
	UnmarshalErrCode          = "DynamoDBUnmarshalErr"
	InvalidConditionErrCode   = "DynamoDBInvalidConditionErr"
	ValidationErrCode         = "ValidationException"
	ThrottlingErrCode         = "ThrottlingException"
	UnrecognizedClientErrCode = "UnrecognizedClientException"
)

type Error struct {
	Message string
	Code    string
}

func NewError(message, code string) Error {
	return Error{Message: message, Code: code}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) MarshallingFailed() bool {
	return anyEquals(e.Code, MarshalErrCode)
}

func (e Error) UnmarshallingFailed() bool {
	return anyEquals(e.Code, UnmarshalErrCode)
}

func (e Error) InvalidCondition() bool {
	return anyEquals(e.Code, InvalidConditionErrCode)
}

func (e Error) ValidationFailed() bool {
	return anyEquals(e.Code, ValidationErrCode)
}

func (e Error) ConditionFailed() bool {
	return anyEquals(e.Code, dynamodb.ErrCodeConditionalCheckFailedException)
}

func (e Error) BackupUnavailable() bool {
	return anyEquals(e.Code,
		dynamodb.ErrCodeContinuousBackupsUnavailableException,
		dynamodb.ErrCodePointInTimeRecoveryUnavailableException)
}

func (e Error) InternalError() bool {
	return anyEquals(e.Code, dynamodb.ErrCodeInternalServerError)
}

func (e Error) ResourceNotFound() bool {
	return anyEquals(e.Code,
		dynamodb.ErrCodeResourceNotFoundException,
		dynamodb.ErrCodeBackupNotFoundException,
		dynamodb.ErrCodeGlobalTableNotFoundException,
		dynamodb.ErrCodeIndexNotFoundException,
		dynamodb.ErrCodeReplicaNotFoundException,
		dynamodb.ErrCodeTableNotFoundException)
}

func (e Error) ResourceAlreadyExists() bool {
	return anyEquals(e.Code,
		dynamodb.ErrCodeGlobalTableAlreadyExistsException,
		dynamodb.ErrCodeReplicaAlreadyExistsException,
		dynamodb.ErrCodeTableAlreadyExistsException)
}

func (e Error) InvalidOperation() bool {
	return anyEquals(e.Code,
		dynamodb.ErrCodeIdempotentParameterMismatchException,
		dynamodb.ErrCodeInvalidRestoreTimeException,
		dynamodb.ErrCodeTransactionCanceledException,
		dynamodb.ErrCodeTransactionConflictException,
		dynamodb.ErrCodeTransactionInProgressException)
}

func (e Error) LimitExceeded() bool {
	return anyEquals(e.Code,
		dynamodb.ErrCodeItemCollectionSizeLimitExceededException,
		dynamodb.ErrCodeLimitExceededException,
		dynamodb.ErrCodeRequestLimitExceeded,
		dynamodb.ErrCodeProvisionedThroughputExceededException)
}

func (e Error) ResourceInUse() bool {
	return anyEquals(e.Code,
		dynamodb.ErrCodeBackupInUseException,
		dynamodb.ErrCodeResourceInUseException,
		dynamodb.ErrCodeTableInUseException)
}

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Programming.Errors.html
func (e Error) Retryable() bool {
	return anyEquals(e.Code,
		ThrottlingErrCode,
		UnrecognizedClientErrCode,
		dynamodb.ErrCodeItemCollectionSizeLimitExceededException,
		dynamodb.ErrCodeLimitExceededException,
		dynamodb.ErrCodeProvisionedThroughputExceededException,
		dynamodb.ErrCodeRequestLimitExceeded,
		dynamodb.ErrCodeInternalServerError)
}

func anyEquals(origCode string, errCodes ...string) bool {
	for _, ec := range errCodes {
		if origCode == ec {
			return true
		}
	}

	return false
}
