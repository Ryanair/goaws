package cognito

import (
	"time"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Ryanair/goaws"
)

var (
	emailDeliveryMethod = cip.DeliveryMediumTypeEmail
	smsDeliveryMethod   = cip.DeliveryMediumTypeSms

	DeliveryMediumEmpty       = DeliveryMedium([]*string{})
	DeliveryMediumEmail       = DeliveryMedium([]*string{&emailDeliveryMethod})
	DeliveryMediumSms         = DeliveryMedium([]*string{&smsDeliveryMethod})
	DeliveryMediumEmailAndSms = DeliveryMedium([]*string{&emailDeliveryMethod, &smsDeliveryMethod})

	listGroupsDefaultLimit int64 = 1000
)

type DeliveryMedium []*string

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

type Group struct {
	CreationDate     *time.Time
	Description      *string
	Name             string
	LastModifiedDate *time.Time
	Precedence       *int64
	RoleArn          *string
	UserPoolId       string
}

type ListGroupsResult struct {
	Groups    []Group
	NextToken *string
}

type provider interface {
	GetUser(*cip.GetUserInput) (*cip.GetUserOutput, error)
	AdminCreateUser(*cip.AdminCreateUserInput) (*cip.AdminCreateUserOutput, error)
	AdminInitiateAuth(*cip.AdminInitiateAuthInput) (*cip.AdminInitiateAuthOutput, error)
	ConfirmForgotPassword(*cip.ConfirmForgotPasswordInput) (*cip.ConfirmForgotPasswordOutput, error)
	AdminRespondToAuthChallenge(*cip.AdminRespondToAuthChallengeInput) (*cip.AdminRespondToAuthChallengeOutput, error)
	ChangePassword(*cip.ChangePasswordInput) (*cip.ChangePasswordOutput, error)
	AdminUserGlobalSignOut(*cip.AdminUserGlobalSignOutInput) (*cip.AdminUserGlobalSignOutOutput, error)
	AdminResetUserPassword(input *cip.AdminResetUserPasswordInput) (*cip.AdminResetUserPasswordOutput, error)
	AdminListGroupsForUser(input *cip.AdminListGroupsForUserInput) (*cip.AdminListGroupsForUserOutput, error)
}

type Adapter struct {
	poolID   string
	clientID string
	provider provider
}

func NewAdapter(cfg *goaws.Config, poolID, clientID string) *Adapter {

	provider := cip.New(cfg.Provider)

	return &Adapter{
		poolID:   poolID,
		clientID: clientID,
		provider: provider,
	}
}

func (ca *Adapter) ChangePassword(username, oldPassword, newPassword string) error {
	authFlow := cip.AuthFlowTypeAdminNoSrpAuth

	input := &cip.AdminInitiateAuthInput{
		AuthFlow: &authFlow,
		AuthParameters: map[string]*string{
			"USERNAME": &username,
			"PASSWORD": &oldPassword,
		},
		ClientId:   &ca.clientID,
		UserPoolId: &ca.poolID,
	}

	output, err := ca.provider.AdminInitiateAuth(input)
	if err != nil {
		return wrapErrWithCode(err, "error in cognito.Adapter while sending AdminInitiateAuthRequest", ErrCodeSignIn)
	}

	switch *output.ChallengeName {
	case cip.ChallengeNameTypeNewPasswordRequired:
		if _, err := ca.respondToAuthChallenge(username, newPassword, output.Session); err != nil {
			return wrapErrWithCode(err, "error in cognito.Adapter while responding to auth challenge", ErrCodeRespondToAuthChallenge)
		}
	default:
		if _, err := ca.changePassword(oldPassword, newPassword, output.AuthenticationResult.AccessToken); err != nil {
			return wrapErrWithCode(err, "error in cognito.Adapter while changing password", ErrCodeChangePasswordRequest)
		}
	}
	return nil
}

func (ca *Adapter) GetUser(accessToken string) (*GetUserResult, error) {

	getUserInput := &cip.GetUserInput{
		AccessToken: &accessToken,
	}

	output, err := ca.provider.GetUser(getUserInput)
	if err != nil {
		return nil, wrapErr(err, "error in cognito.Adapter while sending GetUserRequest")
	}

	result := &GetUserResult{
		UserAttributes: fromAttributes(output.UserAttributes),
		Username:       output.Username,
	}
	return result, nil
}

func (ca *Adapter) CreateUser(username, password string, attributesMap map[string]string, deliveryMedium DeliveryMedium,
	createAlias bool) (*CreateUserResult, error) {

	deliveryMediums := make([]*string, 0)
	for _, medium := range deliveryMedium {
		deliveryMediums = append(deliveryMediums, medium)
	}

	input := &cip.AdminCreateUserInput{
		ForceAliasCreation:     &createAlias,
		UserAttributes:         toAttributes(attributesMap),
		DesiredDeliveryMediums: deliveryMediums,
		TemporaryPassword:      &password,
		UserPoolId:             &ca.poolID,
		Username:               &username,
	}

	output, err := ca.provider.AdminCreateUser(input)
	if err != nil {
		return nil, wrapErr(err, "error in cognito.Adapter while sending AdminCreateUserRequest")
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

	authFlow := cip.AuthFlowTypeAdminNoSrpAuth

	input := &cip.AdminInitiateAuthInput{
		AuthFlow: &authFlow,
		AuthParameters: map[string]*string{
			"USERNAME": &username,
			"PASSWORD": &password,
		},
		ClientId:   &ca.clientID,
		UserPoolId: &ca.poolID,
	}

	output, err := ca.provider.AdminInitiateAuth(input)
	if err != nil {
		return nil, wrapErr(err, "error in cognito.Adapter while sending AdminInitiateAuthRequest")
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

	input := &cip.AdminUserGlobalSignOutInput{
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	_, err := ca.provider.AdminUserGlobalSignOut(input)
	if err != nil {
		return wrapErr(err, "error in cognito.Adapter while sending AdminUserGlobalSignOutRequest")
	}
	return nil
}

func (ca *Adapter) ResetUserPassword(username string) error {

	input := &cip.AdminResetUserPasswordInput{
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	_, err := ca.provider.AdminResetUserPassword(input)
	if err != nil {
		return wrapErr(err, "error in cognito.Adapter while sending AdminResetUserPasswordRequest")
	}
	return nil
}

func (ca *Adapter) ConfirmForgotPassword(username, newPassword, confirmationCode string) error {

	input := &cip.ConfirmForgotPasswordInput{
		ClientId:         &ca.clientID,
		ConfirmationCode: &confirmationCode,
		Password:         &newPassword,
		Username:         &username,
	}

	_, err := ca.provider.ConfirmForgotPassword(input)
	if err != nil {
		return wrapErr(err, "error in cognito.Adapter while sending ConfirmForgotPasswordRequest")
	}
	return nil
}

func (ca *Adapter) ListUserGroups(username string, limit *int64, nextToken *string) (*ListGroupsResult, error) {

	if limit == nil {
		limit = &listGroupsDefaultLimit
	}

	input := &cip.AdminListGroupsForUserInput{
		Limit:      limit,
		NextToken:  nextToken,
		UserPoolId: &ca.poolID,
		Username:   &username,
	}

	output, err := ca.provider.AdminListGroupsForUser(input)
	if err != nil {
		return nil, wrapErr(err, "error in cognito.Adapter while sending ConfirmForgotPasswordRequest")
	}

	groups := make([]Group, 0)
	for _, cg := range output.Groups {
		group := Group{
			CreationDate:     cg.CreationDate,
			Description:      cg.Description,
			Name:             *cg.GroupName,
			LastModifiedDate: cg.LastModifiedDate,
			Precedence:       cg.Precedence,
			RoleArn:          cg.RoleArn,
			UserPoolId:       *cg.UserPoolId,
		}
		groups = append(groups, group)
	}

	return &ListGroupsResult{
		Groups:    groups,
		NextToken: output.NextToken,
	}, nil
}

func (ca *Adapter) respondToAuthChallenge(username, password string,
	session *string) (*cip.AdminRespondToAuthChallengeOutput, error) {

	challengeName := cip.ChallengeNameTypeNewPasswordRequired

	input := &cip.AdminRespondToAuthChallengeInput{
		ChallengeName: &challengeName,
		ChallengeResponses: map[string]*string{
			"USERNAME":     &username,
			"NEW_PASSWORD": &password,
		},
		ClientId:   &ca.clientID,
		UserPoolId: &ca.poolID,
		Session:    session,
	}

	output, err := ca.provider.AdminRespondToAuthChallenge(input)
	if err != nil {
		return nil, wrapErr(err, "error in cognito.Adapter while sending AdminRespondToAuthChallengeRequest")
	}
	return output, nil
}

func (ca *Adapter) changePassword(oldPassword, newPassword string, token *string) (*cip.ChangePasswordOutput, error) {

	input := &cip.ChangePasswordInput{
		AccessToken:      token,
		PreviousPassword: &oldPassword,
		ProposedPassword: &newPassword,
	}

	output, err := ca.provider.ChangePassword(input)
	if err != nil {
		return nil, wrapErr(err, "error in cognito.Adapter while sending ChangePasswordRequest")
	}
	return output, nil
}

func fromAttributes(attrs []*cip.AttributeType) map[string]string {
	attributesMap := make(map[string]string)
	for _, attr := range attrs {
		if attr.Name == nil {
			continue
		}
		if attr.Value != nil {
			attributesMap[*attr.Name] = *attr.Value
		}
	}
	return attributesMap
}

func toAttributes(attributesMap map[string]string) []*cip.AttributeType {
	attributes := make([]*cip.AttributeType, 0)
	for attrName, attrValue := range attributesMap {
		attr := attribute(attrName, attrValue)
		attributes = append(attributes, attr)
	}
	return attributes
}

func attribute(name, value string) *cip.AttributeType {
	return &cip.AttributeType{
		Name:  &name,
		Value: &value,
	}
}
