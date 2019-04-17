package cognito

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Ryanair/goaws"
)

var (
	emailDeliveryMethod = cognitoidentityprovider.DeliveryMediumTypeEmail
	smsDeliveryMethod   = cognitoidentityprovider.DeliveryMediumTypeSms

	DeliveryMediumEmpty       = deliveryMedium([]*string{})
	DeliveryMediumEmail       = deliveryMedium([]*string{&emailDeliveryMethod})
	DeliveryMediumSms         = deliveryMedium([]*string{&smsDeliveryMethod})
	DeliveryMediumEmailAndSms = deliveryMedium([]*string{&emailDeliveryMethod, &smsDeliveryMethod})
)

type deliveryMedium []*string

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
	poolID       string
	clientID     string
	clientSecret string
	provider     *cognitoidentityprovider.CognitoIdentityProvider
}

func NewAdapter(cfg *goaws.Config, poolID, clientID, clientSecret string) *Adapter {

	provider := cognitoidentityprovider.New(cfg.Provider)

	return &Adapter{
		poolID:       poolID,
		clientID:     clientID,
		clientSecret: clientSecret,
		provider:     provider,
	}
}

func (ca *Adapter) ChangePassword(username, oldPassword, newPassword string) error {

	adminAuthResponse, err := ca.SignIn(username, oldPassword)
	if err != nil {
		return wrapOpsErr(err, "error in cognito.Adapter while signing in before changing password")
	}

	switch *adminAuthResponse.ChallengeName {
	case cognitoidentityprovider.ChallengeNameTypeNewPasswordRequired:
		_, err := ca.respondToAuthChallenge(username, newPassword, adminAuthResponse.Session)
		return wrapOpsErr(err, "error in cognito.Adapter while responding to auth challenge")
	default:
		_, err := ca.changePassword(oldPassword, newPassword, adminAuthResponse.AuthenticationResult.AccessToken)
		return wrapOpsErr(err, "error in cognito.Adapter while changing password")
	}
}

func (ca *Adapter) GetUser(accessToken string) (*GetUserResult, error) {

	getUserInput := &cognitoidentityprovider.GetUserInput{
		AccessToken: &accessToken,
	}

	request, output := ca.provider.GetUserRequest(getUserInput)
	if err := request.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending GetUserRequest")
	}

	result := &GetUserResult{
		UserAttributes: fromAttributes(output.UserAttributes),
		Username:       output.Username,
	}
	return result, nil
}

func (ca *Adapter) CreateUser(username, password string, attributesMap map[string]string, deliveryMedium deliveryMedium,
	createAlias bool) (*CreateUserResult, error) {

	deliveryMediums := make([]*string, 0)
	for _, medium := range deliveryMedium {
		deliveryMediums = append(deliveryMediums, medium)
	}

	input := &cognitoidentityprovider.AdminCreateUserInput{
		ForceAliasCreation:     &createAlias,
		UserAttributes:         toAttributes(attributesMap),
		DesiredDeliveryMediums: deliveryMediums,
		TemporaryPassword:      &password,
		UserPoolId:             &ca.poolID,
		Username:               &username,
	}

	request, output := ca.provider.AdminCreateUserRequest(input)
	if err := request.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending AdminCreateUserRequest")
	}

	user := output.User
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

	input := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: &authFlow,
		AuthParameters: map[string]*string{
			"USERNAME":    &username,
			"PASSWORD":    &password,
			"SECRET_HASH": &secretHash,
		},
		ClientId:   &ca.clientID,
		UserPoolId: &ca.poolID,
	}

	request, output := ca.provider.AdminInitiateAuthRequest(input)
	if err := request.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending AdminInitiateAuthRequest")
	}

	return &SignInResult{
		AuthenticationResult: &AuthenticationResult{
			AccessToken:  output.AuthenticationResult.AccessToken,
			ExpiresIn:    output.AuthenticationResult.ExpiresIn,
			IDToken:      output.AuthenticationResult.IdToken,
			RefreshToken: output.AuthenticationResult.RefreshToken,
			TokenType:    output.AuthenticationResult.TokenType,
		},
		ChallengeName:       output.ChallengeName,
		ChallengeParameters: output.ChallengeParameters,
		Session:             output.Session,
	}, nil
}

func (ca *Adapter) SignOut(username string) error {

	input := &cognitoidentityprovider.AdminUserGlobalSignOutInput{
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	request, _ := ca.provider.AdminUserGlobalSignOutRequest(input)
	if err := request.Send(); err != nil {
		return wrapOpsErr(err, "error in cognito.Adapter while sending AdminUserGlobalSignOutRequest")
	}
	return nil
}

func (ca *Adapter) ResetUserPassword(username string) error {

	input := &cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	request, _ := ca.provider.AdminResetUserPasswordRequest(input)
	if err := request.Send(); err != nil {
		return wrapOpsErr(err, "error in cognito.Adapter while sending AdminResetUserPasswordRequest")
	}
	return nil
}

func (ca *Adapter) ConfirmForgotPassword(username, newPassword, confirmationCode string) error {

	secretHash, err := ca.generateSecretHash(username)
	if err != nil {
		return wrapErr(err, ErrSecretHashEncoding, "Cannot encode secret hash")
	}

	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         &ca.clientID,
		ConfirmationCode: &confirmationCode,
		Password:         &newPassword,
		SecretHash:       &secretHash,
		Username:         &username,
	}

	request, _ := ca.provider.ConfirmForgotPasswordRequest(input)
	if err := request.Send(); err != nil {
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

	input := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
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

	request, output := ca.provider.AdminRespondToAuthChallengeRequest(input)
	if err := request.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending AdminRespondToAuthChallengeRequest")
	}
	return output, nil
}

func (ca *Adapter) changePassword(oldPassword, newPassword string, token *string) (*cognitoidentityprovider.ChangePasswordOutput, error) {

	input := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      token,
		PreviousPassword: &oldPassword,
		ProposedPassword: &newPassword,
	}

	request, output := ca.provider.ChangePasswordRequest(input)
	if err := request.Send(); err != nil {
		return nil, wrapOpsErr(err, "error in cognito.Adapter while sending ChangePasswordRequest")
	}
	return output, nil
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
		if attr.Name != nil {
			continue
		}
		value := ""
		if attr.Value == nil {
			value = *attr.Value
		}
		attributesMap[*attr.Name] = value
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
