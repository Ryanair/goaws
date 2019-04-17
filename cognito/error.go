package cognito

import (
	"github.com/Ryanair/goaws/internal"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	ErrSecretHashEncoding = "SecretHashEncoding"
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

func (e Error) AliasExists() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeAliasExistsException)
}

func (e Error) CodeDeliveryFailure() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeCodeDeliveryFailureException)
}

func (e Error) CodeMismatch() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeCodeMismatchException)
}

func (e Error) CodeExpired() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeExpiredCodeException)
}

func (e Error) GroupExists() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeGroupExistsException)
}

func (e Error) InvalidPassword() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeInvalidPasswordException)
}

func (e Error) NotAuthorized() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeNotAuthorizedException)
}

func (e Error) PasswordResetRequired() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodePasswordResetRequiredException)
}

func (e Error) UserNotConfirmed() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeUserNotConfirmedException)
}

func (e Error) UserNotFound() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeUserNotFoundException)
}

func (e Error) UsernameExists() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeUsernameExistsException)
}

func (e Error) UnsupportedUserState() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeUnsupportedUserStateException)
}

func (e Error) InternalError() bool {
	return internal.AnyEquals(e.Code, cognitoidentityprovider.ErrCodeConcurrentModificationException,
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
