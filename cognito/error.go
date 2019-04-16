package cognito

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	ErrSecretHashEncoding = "SecretHashEncoding"
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

func (e Error) AliasExists() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeAliasExistsException)
}

func (e Error) CodeDeliveryFailure() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeCodeDeliveryFailureException)
}

func (e Error) CodeMismatch() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeCodeMismatchException)
}

func (e Error) CodeExpired() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeExpiredCodeException)
}

func (e Error) GroupExists() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeGroupExistsException)
}

func (e Error) InvalidPassword() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeInvalidPasswordException)
}

func (e Error) NotAuthorized() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeNotAuthorizedException)
}

func (e Error) PasswordResetRequired() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodePasswordResetRequiredException)
}

func (e Error) UserNotConfirmed() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeUserNotConfirmedException)
}

func (e Error) UserNotFound() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeUserNotFoundException)
}

func (e Error) UsernameExists() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeUsernameExistsException)
}

func (e Error) UnsupportedUserState() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeUnsupportedUserStateException)
}

func (e Error) InternalError() bool {
	return anyEquals(e.Code, cognitoidentityprovider.ErrCodeConcurrentModificationException,
		cognitoidentityprovider.ErrCodeDuplicateProviderException,
		cognitoidentityprovider.ErrCodeEnableSoftwareTokenMFAException,
		cognitoidentityprovider.ErrCodeInternalErrorException,
		cognitoidentityprovider.ErrCodeInvalidEmailRoleAccessPolicyException,
		cognitoidentityprovider.ErrCodeInvalidLambdaResponseException,
		cognitoidentityprovider.ErrCodeInvalidOAuthFlowException,
		cognitoidentityprovider.ErrCodeInvalidParameterException,
		cognitoidentityprovider.ErrCodeInvalidSmsRoleAccessPolicyException,
		cognitoidentityprovider.ErrCodeInvalidSmsRoleTrustRelationshipException,
		cognitoidentityprovider.ErrCodeInvalidUserPoolConfigurationException,
		cognitoidentityprovider.ErrCodeLimitExceededException,
		cognitoidentityprovider.ErrCodeMFAMethodNotFoundException,
		cognitoidentityprovider.ErrCodePreconditionNotMetException,
		cognitoidentityprovider.ErrCodeResourceNotFoundException,
		cognitoidentityprovider.ErrCodeScopeDoesNotExistException,
		cognitoidentityprovider.ErrCodeSoftwareTokenMFANotFoundException,
		cognitoidentityprovider.ErrCodeTooManyFailedAttemptsException,
		cognitoidentityprovider.ErrCodeTooManyRequestsException,
		cognitoidentityprovider.ErrCodeUnexpectedLambdaException,
		cognitoidentityprovider.ErrCodeUnsupportedIdentityProviderException,
		cognitoidentityprovider.ErrCodeUserImportInProgressException,
		cognitoidentityprovider.ErrCodeUserLambdaValidationException,
		cognitoidentityprovider.ErrCodeUserPoolAddOnNotEnabledException,
		cognitoidentityprovider.ErrCodeUserPoolTaggingException,
	)
}

func anyEquals(origCode string, errCodes ...string) bool {
	for _, ec := range errCodes {
		if origCode == ec {
			return true
		}
	}
	return false
}
