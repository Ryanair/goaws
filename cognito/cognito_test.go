package cognito

import (
	"testing"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	poolID   = "abc-def-pool-id"
	clientID = "ghi-jklm-client-id"
)

func NewTestAdapter(provider provider) *Adapter {
	return &Adapter{
		poolID:   poolID,
		clientID: clientID,
		provider: provider,
	}
}

var (
	username    = "John"
	email       = "john@example.com"
	oldPassword = "oldSecret"
	newPassword = "newSecret"

	session = "abc-session-id"

	emailAttr = "email"

	errCodeSignIn                 = ErrCodeSignIn
	errCodeRespondToAuthChallenge = ErrCodeRespondToAuthChallenge
	errCodeChangePasswordRequest  = ErrCodeChangePasswordRequest
)

func TestAdapter_ChangePassword(t *testing.T) {

	// given
	var testData = []struct {
		signInErr                 error
		challengeName             string
		respondToAuthChallengeErr error
		changePasswordErr         error
		expectedErrorCode         *string
	}{
		{nil, cip.ChallengeNameTypeNewPasswordRequired, nil, nil, nil},
		{errors.New("error in signIn"), "", nil, nil, &errCodeSignIn},
		{nil, cip.ChallengeNameTypeNewPasswordRequired, errors.New("error in respondToAuthChallenge"), nil, &errCodeRespondToAuthChallenge},
		{nil, "", nil, errors.New("error in changePassword"), &errCodeChangePasswordRequest},
	}

	for _, data := range testData {

		signInOutput := &cip.AdminInitiateAuthOutput{
			AuthenticationResult: &cip.AuthenticationResultType{},
			ChallengeName:        &data.challengeName,
			Session:              &session,
		}

		provider := &providerMock{
			respondToAuthChallengeOutput: &cip.AdminRespondToAuthChallengeOutput{},
			respondToAuthChallengeErr:    data.respondToAuthChallengeErr,
			changePasswordOutput:         &cip.ChangePasswordOutput{},
			changePasswordErr:            data.changePasswordErr,
			authOutput:                   signInOutput,
			authErr:                      data.signInErr,
		}
		adapter := NewTestAdapter(provider)

		// when
		err := adapter.ChangePassword(username, oldPassword, newPassword)

		// then
		if err != nil {
			awsgoError := toAwsgoError(t, err)
			if data.expectedErrorCode != nil {
				assert.Equal(t, *data.expectedErrorCode, awsgoError.Code)
			}
		}
	}
}

func TestAdapter_GetUser_ok(t *testing.T) {

	// given
	emailAttribute := cip.AttributeType{Name: &emailAttr, Value: &email}
	getUserOutput := &cip.GetUserOutput{
		Username:       &username,
		UserAttributes: []*cip.AttributeType{&emailAttribute},
	}

	provider := &providerMock{
		getUserOutput: getUserOutput,
	}

	adapter := NewTestAdapter(provider)

	// when
	user, err := adapter.GetUser("abc-access-token")

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, user.UserAttributes[emailAttr])
	assert.Equal(t, email, user.UserAttributes[emailAttr])
	assert.Equal(t, *getUserOutput.Username, *user.Username)
}

func TestAdapter_GetUser_error(t *testing.T) {

	// given
	getUserErr := errors.New("error while getting user data")
	provider := &providerMock{getUserErr: getUserErr}
	adapter := NewTestAdapter(provider)

	// when
	user, err := adapter.GetUser("abc-access-token")

	assert.Error(t, err)
	assert.Nil(t, user)

	awsgoError := toAwsgoError(t, err)
	assert.Contains(t, awsgoError.Message, getUserErr.Error())
}

func TestAdapter_CreateUser_ok(t *testing.T) {

	// given
	attrs := map[string]string{emailAttr: email}

	emailAttribute := cip.AttributeType{Name: &emailAttr, Value: &email}
	createUserOutput := &cip.AdminCreateUserOutput{
		User: &cip.UserType{
			Attributes: []*cip.AttributeType{&emailAttribute},
			Username:   &username,
		},
	}

	provider := &providerMock{
		createUserOutput: createUserOutput,
	}

	adapter := NewTestAdapter(provider)

	// when
	user, err := adapter.CreateUser(username, newPassword, attrs, DeliveryMediumEmailAndSms, false)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, email, user.Attributes[emailAttr])
	assert.Equal(t, username, *user.Username)
}

func TestAdapter_CreateUser_error(t *testing.T) {

	// given
	provider := &providerMock{
		createUserErr: errors.New("error while creating user"),
	}

	adapter := NewTestAdapter(provider)

	// when
	user, err := adapter.CreateUser(username, newPassword, nil, DeliveryMediumEmailAndSms, false)

	// then
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestAdapter_SignIn_ok(t *testing.T) {

	// given
	provider := &providerMock{
		authOutput: &cip.AdminInitiateAuthOutput{
			AuthenticationResult: &cip.AuthenticationResultType{},
		},
	}

	adapter := NewTestAdapter(provider)

	// when
	result, err := adapter.SignIn(username, newPassword)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAdapter_SignIn_error(t *testing.T) {

	// given
	provider := &providerMock{
		authErr: errors.New("error while sending auth request"),
	}

	adapter := NewTestAdapter(provider)

	// when
	result, err := adapter.SignIn(username, newPassword)

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAdapter_SignOut_ok(t *testing.T) {

	// given
	provider := &providerMock{
		signOutOutput: &cip.AdminUserGlobalSignOutOutput{},
	}

	adapter := NewTestAdapter(provider)

	// when
	err := adapter.SignOut(username)

	// then
	assert.NoError(t, err)
}

func TestAdapter_SignOut_error(t *testing.T) {

	// given
	provider := &providerMock{
		signOutErr: errors.New("error while signing out"),
	}

	adapter := NewTestAdapter(provider)

	// when
	err := adapter.SignOut(username)

	// then
	assert.Error(t, err)
}

func TestAdapter_ResetUserPassword_ok(t *testing.T) {

	// given
	provider := &providerMock{
		resetPassOutput: &cip.AdminResetUserPasswordOutput{},
	}

	adapter := NewTestAdapter(provider)

	// when
	err := adapter.ResetUserPassword(username)

	// then
	assert.NoError(t, err)
}

func TestAdapter_ResetUserPassword_error(t *testing.T) {

	// given
	provider := &providerMock{
		resetPassErr: errors.New("error while reseting password"),
	}

	adapter := NewTestAdapter(provider)

	// when
	err := adapter.ResetUserPassword(username)

	// then
	assert.Error(t, err)
}

func TestAdapter_ConfirmForgotPassword_ok(t *testing.T) {

	// given
	provider := &providerMock{
		forgetPassOutput: &cip.ConfirmForgotPasswordOutput{},
	}

	adapter := NewTestAdapter(provider)

	// when
	err := adapter.ConfirmForgotPassword(username, newPassword, "23649")

	// then
	assert.NoError(t, err)
}

func TestAdapter_ConfirmForgotPassword_error(t *testing.T) {

	// given
	provider := &providerMock{
		forgetPassErr: errors.New("error while confirming password"),
	}

	adapter := NewTestAdapter(provider)

	// when
	err := adapter.ConfirmForgotPassword(username, newPassword, "23649")

	// then
	assert.Error(t, err)
}

func toAwsgoError(t *testing.T, err error) Error {
	awsgoError, ok := err.(Error)
	if !ok {
		t.Errorf("invalid error type, expected Error")
	}
	return awsgoError
}

type providerMock struct {
	getUserOutput                *cip.GetUserOutput
	getUserErr                   error
	createUserOutput             *cip.AdminCreateUserOutput
	createUserErr                error
	authOutput                   *cip.AdminInitiateAuthOutput
	authErr                      error
	forgetPassOutput             *cip.ConfirmForgotPasswordOutput
	forgetPassErr                error
	respondToAuthChallengeOutput *cip.AdminRespondToAuthChallengeOutput
	respondToAuthChallengeErr    error
	changePasswordOutput         *cip.ChangePasswordOutput
	changePasswordErr            error
	signOutOutput                *cip.AdminUserGlobalSignOutOutput
	signOutErr                   error
	resetPassOutput              *cip.AdminResetUserPasswordOutput
	resetPassErr                 error
}

func (pm *providerMock) GetUser(*cip.GetUserInput) (*cip.GetUserOutput, error) {
	return pm.getUserOutput, pm.getUserErr
}

func (pm *providerMock) AdminCreateUser(*cip.AdminCreateUserInput) (*cip.AdminCreateUserOutput, error) {
	return pm.createUserOutput, pm.createUserErr
}

func (pm *providerMock) AdminInitiateAuth(*cip.AdminInitiateAuthInput) (*cip.AdminInitiateAuthOutput, error) {
	return pm.authOutput, pm.authErr
}

func (pm *providerMock) ConfirmForgotPassword(*cip.ConfirmForgotPasswordInput) (*cip.ConfirmForgotPasswordOutput, error) {
	return pm.forgetPassOutput, pm.forgetPassErr
}

func (pm *providerMock) AdminRespondToAuthChallenge(*cip.AdminRespondToAuthChallengeInput) (*cip.AdminRespondToAuthChallengeOutput, error) {
	return pm.respondToAuthChallengeOutput, pm.respondToAuthChallengeErr
}

func (pm *providerMock) ChangePassword(*cip.ChangePasswordInput) (*cip.ChangePasswordOutput, error) {
	return pm.changePasswordOutput, pm.changePasswordErr
}

func (pm *providerMock) AdminUserGlobalSignOut(*cip.AdminUserGlobalSignOutInput) (*cip.AdminUserGlobalSignOutOutput, error) {
	return pm.signOutOutput, pm.signOutErr
}

func (pm *providerMock) AdminResetUserPassword(input *cip.AdminResetUserPasswordInput) (*cip.AdminResetUserPasswordOutput, error) {
	return pm.resetPassOutput, pm.resetPassErr
}

func TestProviderMockImplementsProvider(t *testing.T) {
	var _ provider = &providerMock{}
}
