package cognito

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/Ryanair/goaws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

type deliveryMediums []*string

var (
	emailDeliveryMethod = cognitoidentityprovider.DeliveryMediumTypeEmail
	smsDeliveryMethod   = cognitoidentityprovider.DeliveryMediumTypeSms

	DeliveryMediumsEmpty       = deliveryMediums([]*string{})
	DeliveryMediumsEmail       = deliveryMediums([]*string{&emailDeliveryMethod})
	DeliveryMediumsSms         = deliveryMediums([]*string{&smsDeliveryMethod})
	DeliveryMediumsEmailAndSms = deliveryMediums([]*string{&emailDeliveryMethod, &smsDeliveryMethod})
)

type AuthenticationResult struct {
	AccessToken  *string
	ExpiresIn    *int64
	IDToken      *string
	RefreshToken *string
	TokenType    *string
}

type SignInResult struct {
	AuthenticationResult *AuthenticationResult
	ChallengeName        *string
	ChallengeParameters  map[string]*string
	Session              *string
}

type CreateUserResult struct {
	Attributes       map[string]string
	Enabled          *bool
	CreateDate       *time.Time
	LastModifiedDate *time.Time
	UserStatus       *string
	Username         *string
}

type GetUserResult struct {
	UserAttributes map[string]string
	Username       *string
}

type Adapter struct {
	poolID          string
	clientID        string
	clientSecret    string
	createAlias     bool
	provider        *cognitoidentityprovider.CognitoIdentityProvider
	deliveryMediums deliveryMediums
}

func NewAdapter(cfg *goaws.Config, poolID, clientID, clientSecret string, deliveryMediums deliveryMediums, options ...func(*Adapter) *Adapter) *Adapter {

	provider := cognitoidentityprovider.New(cfg.Provider)

	adapter := &Adapter{
		poolID:          poolID,
		clientID:        clientID,
		clientSecret:    clientSecret,
		provider:        provider,
		deliveryMediums: deliveryMediums,
	}

	for _, opt := range options {
		adapter = opt(adapter)
	}

	return adapter
}

func ForceAliasCreation(adapter *Adapter) *Adapter {
	adapter.createAlias = true
	return adapter
}

func (ca *Adapter) ChangePassword(username, oldPassword, newPassword string) error {

	adminAuthResponse, err := ca.SignIn(username, oldPassword)
	if err != nil {
		return err
	}

	switch *adminAuthResponse.ChallengeName {
	case cognitoidentityprovider.ChallengeNameTypeNewPasswordRequired:
		_, err := ca.respondToAuthChallenge(username, newPassword, adminAuthResponse.Session)
		return err
	default:
		_, err := ca.changePassword(oldPassword, newPassword, adminAuthResponse.AuthenticationResult.AccessToken)
		return err
	}
}

func (ca *Adapter) GetUser(accessToken string) (*GetUserResult, error) {

	getUserInput := &cognitoidentityprovider.GetUserInput{
		AccessToken: &accessToken,
	}

	getUserRequest, getUserOutput := ca.provider.GetUserRequest(getUserInput)
	if err := getUserRequest.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending GetUserRequest")
	}

	result := &GetUserResult{
		UserAttributes: fromAttributes(getUserOutput.UserAttributes),
		Username:       getUserOutput.Username,
	}
	return result, nil
}

func (ca *Adapter) CreateUser(username, password string, attributesMap map[string]string) (*CreateUserResult, error) {

	deliveryMediums := make([]*string, 0)
	for _, medium := range ca.deliveryMediums {
		deliveryMediums = append(deliveryMediums, medium)
	}

	adminCreateUserInput := &cognitoidentityprovider.AdminCreateUserInput{
		ForceAliasCreation:     &ca.createAlias,
		UserAttributes:         toAttributes(attributesMap),
		DesiredDeliveryMediums: deliveryMediums,
		TemporaryPassword:      &password,
		UserPoolId:             &ca.poolID,
		Username:               &username,
	}

	adminCreateUserRequest, adminCreateUserOutput := ca.provider.AdminCreateUserRequest(adminCreateUserInput)
	if err := adminCreateUserRequest.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending AdminCreateUserRequest")
	}

	user := adminCreateUserOutput.User
	return &CreateUserResult{
		Attributes:       fromAttributes(user.Attributes),
		Enabled:          user.Enabled,
		CreateDate:       user.UserCreateDate,
		LastModifiedDate: user.UserLastModifiedDate,
		UserStatus:       user.UserStatus,
		Username:         user.Username,
	}, nil
}

func (ca *Adapter) SignIn(username, password string) (*SignInResult, error) {

	secretHash, err := ca.generateSecretHash(username)
	if err != nil {
		return nil, wrapErr(err, ErrSecretHashEncoding, "Cannot encode secret hash")
	}
	authFlow := cognitoidentityprovider.AuthFlowTypeAdminNoSrpAuth

	adminInitiateAuthInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: &authFlow,
		AuthParameters: map[string]*string{
			"USERNAME":    &username,
			"PASSWORD":    &password,
			"SECRET_HASH": &secretHash,
		},
		ClientId:   &ca.clientID,
		UserPoolId: &ca.poolID,
	}

	adminInitiateAuthRequest, adminInitiateAuthOutput := ca.provider.AdminInitiateAuthRequest(adminInitiateAuthInput)
	if err := adminInitiateAuthRequest.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending AdminInitiateAuthRequest")
	}

	return &SignInResult{
		AuthenticationResult: &AuthenticationResult{
			AccessToken:  adminInitiateAuthOutput.AuthenticationResult.AccessToken,
			ExpiresIn:    adminInitiateAuthOutput.AuthenticationResult.ExpiresIn,
			IDToken:      adminInitiateAuthOutput.AuthenticationResult.IdToken,
			RefreshToken: adminInitiateAuthOutput.AuthenticationResult.RefreshToken,
			TokenType:    adminInitiateAuthOutput.AuthenticationResult.TokenType,
		},
		ChallengeName:       adminInitiateAuthOutput.ChallengeName,
		ChallengeParameters: adminInitiateAuthOutput.ChallengeParameters,
		Session:             adminInitiateAuthOutput.Session,
	}, nil
}

func (ca *Adapter) SignOut(username string) error {

	adminUserGlobalSignOutInput := &cognitoidentityprovider.AdminUserGlobalSignOutInput{
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	adminUserGlobalSignOutRequest, _ := ca.provider.AdminUserGlobalSignOutRequest(adminUserGlobalSignOutInput)
	if err := adminUserGlobalSignOutRequest.Send(); err != nil {
		return wrapOpsErr(err, "error in cognito.Adapter while sending AdminUserGlobalSignOutRequest")
	}
	return nil
}

func (ca *Adapter) ResetUserPassword(username string) error {

	adminResetUserPasswordInput := &cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	adminResetUserPasswordRequest, _ := ca.provider.AdminResetUserPasswordRequest(adminResetUserPasswordInput)
	if err := adminResetUserPasswordRequest.Send(); err != nil {
		return wrapOpsErr(err, "error in cognito.Adapter while sending AdminResetUserPasswordRequest")
	}
	return nil
}

func (ca *Adapter) ConfirmForgotPassword(username, newPassword, confirmationCode string) error {

	secretHash, err := ca.generateSecretHash(username)
	if err != nil {
		return wrapErr(err, ErrSecretHashEncoding, "Cannot encode secret hash")
	}

	confirmForgotPasswordInput := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         &ca.clientID,
		ConfirmationCode: &confirmationCode,
		Password:         &newPassword,
		SecretHash:       &secretHash,
		Username:         &username,
	}

	confirmForgotPasswordRequest, _ := ca.provider.ConfirmForgotPasswordRequest(confirmForgotPasswordInput)
	if err := confirmForgotPasswordRequest.Send(); err != nil {
		return wrapOpsErr(err, "error in cognito.Adapter while sending ConfirmForgotPasswordRequest")
	}
	return nil
}

func (ca *Adapter) respondToAuthChallenge(username, password string,
	session *string) (*cognitoidentityprovider.AdminRespondToAuthChallengeOutput, error) {

	secretHash, err := ca.generateSecretHash(username)
	if err != nil {
		return nil, wrapErr(err, ErrSecretHashEncoding, "Cannot encode secret hash")
	}

	challengeName := cognitoidentityprovider.ChallengeNameTypeNewPasswordRequired

	adminRespondToAuthChallengeInput := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: &challengeName,
		ChallengeResponses: map[string]*string{
			"USERNAME":     &username,
			"NEW_PASSWORD": &password,
			"SECRET_HASH":  &secretHash,
		},
		ClientId:   &ca.clientID,
		UserPoolId: &ca.poolID,
		Session:    session,
	}

	adminRespondToAuthChallengeRequest, adminRespondToAuthChallengeOutput := ca.provider.AdminRespondToAuthChallengeRequest(adminRespondToAuthChallengeInput)
	if err := adminRespondToAuthChallengeRequest.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending AdminRespondToAuthChallengeRequest")
	}
	return adminRespondToAuthChallengeOutput, nil
}

func (ca *Adapter) changePassword(oldPassword, newPassword string, token *string) (*cognitoidentityprovider.ChangePasswordOutput, error) {

	changePasswordInput := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      token,
		PreviousPassword: &oldPassword,
		ProposedPassword: &newPassword,
	}

	changePasswordRequest, changePasswordOutput := ca.provider.ChangePasswordRequest(changePasswordInput)
	if err := changePasswordRequest.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending ChangePasswordRequest")
	}
	return changePasswordOutput, nil
}

func (ca *Adapter) generateSecretHash(username string) (string, error) {
	mac := hmac.New(sha256.New, []byte(ca.clientSecret))
	_, err := mac.Write([]byte(username + ca.clientID))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func fromAttributes(attrs []*cognitoidentityprovider.AttributeType) map[string]string {
	attributesMap := make(map[string]string)
	for _, attr := range attrs {
		attributesMap[*attr.Name] = *attr.Value
	}
	return attributesMap
}

func toAttributes(attributesMap map[string]string) []*cognitoidentityprovider.AttributeType {
	attributes := make([]*cognitoidentityprovider.AttributeType, 0)
	for attrName, attrValue := range attributesMap {
		attr := attribute(attrName, attrValue)
		attributes = append(attributes, attr)
	}
	return attributes
}

func attribute(name, value string) *cognitoidentityprovider.AttributeType {
	return &cognitoidentityprovider.AttributeType{
		Name:  &name,
		Value: &value,
	}
}

func wrapOpsErr(err error, msg string) error {
	wrappedErr := errors.Wrap(err, msg)
	if awsErr, ok := err.(awserr.Error); ok {
		return NewError(wrappedErr.Error(), awsErr.Code())
	}
	return wrappedErr
}

func wrapErr(err error, code, msg string) error {
	wrappedErrMsg := errors.Wrap(err, msg).Error()
	return NewError(wrappedErrMsg, code)
}
