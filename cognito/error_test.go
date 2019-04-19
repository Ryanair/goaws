// +build local ci

package cognito

import (
	"testing"

	"github.com/Ryanair/goaws/internal"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type errorFunction func(code string) bool
type params struct {
	in  string
	out bool
}

func TestErrorBehaviour(t *testing.T) {

	var testData = []struct {
		params   params
		function errorFunction
	}{
		{params{cognitoidentityprovider.ErrCodeAliasExistsException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("alias exists")))
			return err.AliasExists()
		}},
		{params{cognitoidentityprovider.ErrCodeCodeDeliveryFailureException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("code delivery failure")))
			return err.AliasExists()
		}},

		{params{cognitoidentityprovider.ErrCodeCodeDeliveryFailureException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("code delivery failure")))
			return err.CodeDeliveryFailure()
		}},
		{params{cognitoidentityprovider.ErrCodeGroupExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("group exists")))
			return err.CodeDeliveryFailure()
		}},

		{params{cognitoidentityprovider.ErrCodeCodeMismatchException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("code mismatch")))
			return err.CodeMismatch()
		}},
		{params{cognitoidentityprovider.ErrCodeUserNotConfirmedException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("user not confirmed")))
			return err.CodeMismatch()
		}},

		{params{cognitoidentityprovider.ErrCodeExpiredCodeException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("expired code")))
			return err.CodeExpired()
		}},
		{params{cognitoidentityprovider.ErrCodeUnsupportedUserStateException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("unsupported user")))
			return err.CodeExpired()
		}},

		{params{cognitoidentityprovider.ErrCodeGroupExistsException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("group exists")))
			return err.GroupExists()
		}},
		{params{cognitoidentityprovider.ErrCodePasswordResetRequiredException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("password reset")))
			return err.GroupExists()
		}},

		{params{cognitoidentityprovider.ErrCodeInvalidPasswordException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("invalid password")))
			return err.InvalidPassword()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("username exists")))
			return err.InvalidPassword()
		}},

		{params{cognitoidentityprovider.ErrCodeNotAuthorizedException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("code not authorized")))
			return err.NotAuthorized()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("username exists")))
			return err.InvalidPassword()
		}},

		{params{cognitoidentityprovider.ErrCodePasswordResetRequiredException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("password reset exception")))
			return err.PasswordResetRequired()
		}},
		{params{cognitoidentityprovider.ErrCodeInvalidPasswordException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("invalid password")))
			return err.PasswordResetRequired()
		}},

		{params{cognitoidentityprovider.ErrCodeUserNotConfirmedException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("user not confirmed")))
			return err.UserNotConfirmed()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("username exists")))
			return err.UserNotConfirmed()
		}},

		{params{cognitoidentityprovider.ErrCodeUserNotFoundException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("user not found")))
			return err.UserNotFound()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("username exists")))
			return err.UserNotFound()
		}},

		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("username exists")))
			return err.UsernameExists()
		}},
		{params{cognitoidentityprovider.ErrCodeGroupExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("group exists")))
			return err.UsernameExists()
		}},

		{params{cognitoidentityprovider.ErrCodeUnsupportedUserStateException, true}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("unsupported state")))
			return err.UnsupportedUserState()
		}},
		{params{cognitoidentityprovider.ErrCodeGroupExistsException, false}, func(code string) bool {
			err := Error(internal.NewError("", code, errors.New("group exists")))
			return err.UnsupportedUserState()
		}},
	}

	for _, data := range testData {
		result := data.function(data.params.in)
		assert.Equal(t, data.params.out, result)
	}
}

func TestErrorBehaviour_ChangePassword_generateHashFailed(t *testing.T) {
	rootErr := errors.New("generating secret hash failed")
	firstWrapErr := wrapErrWithCode(rootErr, "Cannot encode secret hash", ErrSecretHashEncoding)
	finalErr := wrapErrWithCode(firstWrapErr, "error in cognito.Adapter while signing in before changing password", ErrCodeSignIn)

	isSignInFailed := func(err error) bool {
		type signingFailed interface {
			SignInFailed() bool
		}
		e, ok := err.(signingFailed)
		return ok && e.SignInFailed()
	}

	originErr := errors.Cause(finalErr)

	type causer interface {
		Cause() error
	}
	if cause, ok := finalErr.(causer); ok {
		assert.EqualError(t, cause.Cause(), "Cannot encode secret hash: generating secret hash failed")
	}

	assert.True(t, isSignInFailed(finalErr))
	assert.EqualError(t, originErr, "generating secret hash failed")
	assert.EqualError(t, finalErr, "error in cognito.Adapter while signing in before changing password: Cannot encode secret hash: generating secret hash failed")
}

func TestErrorBehaviour_ChangePassword_awsError(t *testing.T) {
	rootErr := awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "user not found", errors.New("root error"))
	firstWrapErr := wrapErr(rootErr, "error in cognito.Adapter while sending ChangePasswordRequest")
	finalErr := wrapErrWithCode(firstWrapErr, "error in cognito.Adapter while changing password", ErrCodeChangePasswordRequest)

	isChangePasswordFailed := func(err error) bool {
		type changePasswordFailed interface {
			ChangePasswordRequestFailed() bool
		}
		e, ok := err.(changePasswordFailed)
		return ok && e.ChangePasswordRequestFailed()
	}

	originErr := errors.Cause(finalErr)

	type causer interface {
		Cause() error
	}
	if cause, ok := finalErr.(causer); ok {
		assert.EqualError(t, cause.Cause(), "error in cognito.Adapter while sending ChangePasswordRequest: UserNotFoundException: user not found\ncaused by: root error")
	}

	assert.True(t, isChangePasswordFailed(finalErr))
	assert.EqualError(t, originErr, "UserNotFoundException: user not found\ncaused by: root error")
	assert.EqualError(t, finalErr,"error in cognito.Adapter while changing password: error in cognito.Adapter while sending ChangePasswordRequest: UserNotFoundException: user not found\ncaused by: root error")
}
