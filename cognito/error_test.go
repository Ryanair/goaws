package cognito

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
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
			err := NewError("", code)
			return err.AliasExists()
		}},
		{params{cognitoidentityprovider.ErrCodeCodeDeliveryFailureException, false}, func(code string) bool {
			err := NewError("", code)
			return err.AliasExists()
		}},

		{params{cognitoidentityprovider.ErrCodeCodeDeliveryFailureException, true}, func(code string) bool {
			err := NewError("", code)
			return err.CodeDeliveryFailure()
		}},
		{params{cognitoidentityprovider.ErrCodeGroupExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.CodeDeliveryFailure()
		}},

		{params{cognitoidentityprovider.ErrCodeCodeMismatchException, true}, func(code string) bool {
			err := NewError("", code)
			return err.CodeMismatch()
		}},
		{params{cognitoidentityprovider.ErrCodeUserNotConfirmedException, false}, func(code string) bool {
			err := NewError("", code)
			return err.CodeMismatch()
		}},

		{params{cognitoidentityprovider.ErrCodeExpiredCodeException, true}, func(code string) bool {
			err := NewError("", code)
			return err.CodeExpired()
		}},
		{params{cognitoidentityprovider.ErrCodeUnsupportedUserStateException, false}, func(code string) bool {
			err := NewError("", code)
			return err.CodeExpired()
		}},

		{params{cognitoidentityprovider.ErrCodeGroupExistsException, true}, func(code string) bool {
			err := NewError("", code)
			return err.GroupExists()
		}},
		{params{cognitoidentityprovider.ErrCodePasswordResetRequiredException, false}, func(code string) bool {
			err := NewError("", code)
			return err.GroupExists()
		}},

		{params{cognitoidentityprovider.ErrCodeInvalidPasswordException, true}, func(code string) bool {
			err := NewError("", code)
			return err.InvalidPassword()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.InvalidPassword()
		}},

		{params{cognitoidentityprovider.ErrCodeNotAuthorizedException, true}, func(code string) bool {
			err := NewError("", code)
			return err.NotAuthorized()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.InvalidPassword()
		}},

		{params{cognitoidentityprovider.ErrCodePasswordResetRequiredException, true}, func(code string) bool {
			err := NewError("", code)
			return err.PasswordResetRequired()
		}},
		{params{cognitoidentityprovider.ErrCodeInvalidPasswordException, false}, func(code string) bool {
			err := NewError("", code)
			return err.PasswordResetRequired()
		}},

		{params{cognitoidentityprovider.ErrCodeUserNotConfirmedException, true}, func(code string) bool {
			err := NewError("", code)
			return err.UserNotConfirmed()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.UserNotConfirmed()
		}},

		{params{cognitoidentityprovider.ErrCodeUserNotFoundException, true}, func(code string) bool {
			err := NewError("", code)
			return err.UserNotFound()
		}},
		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.UserNotFound()
		}},

		{params{cognitoidentityprovider.ErrCodeUsernameExistsException, true}, func(code string) bool {
			err := NewError("", code)
			return err.UsernameExists()
		}},
		{params{cognitoidentityprovider.ErrCodeGroupExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.UsernameExists()
		}},

		{params{cognitoidentityprovider.ErrCodeUnsupportedUserStateException, true}, func(code string) bool {
			err := NewError("", code)
			return err.UnsupportedUserState()
		}},
		{params{cognitoidentityprovider.ErrCodeGroupExistsException, false}, func(code string) bool {
			err := NewError("", code)
			return err.UnsupportedUserState()
		}},
	}

	for _, data := range testData {
		result := data.function(data.params.in)
		assert.Equal(t, data.params.out, result)
	}
}
