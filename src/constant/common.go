package constant

// constant for service startup failure
const (
	RetrieveAppConfigError                = -1
	RetrieveMockupSessionErrorCode        = -2
	RetrieveMockupSessionErrorMsg         = "retrieve mockup session fail"
	RetrieveMockupOrgErrorCode            = -3
	RetrieveMockupOrgErrorMsg             = "retrieve mockup organization fail"
	RetrieveMockupDataUnmarshalErrorMsg   = "unmarshal mockup session error"
	WriteMockupOrgErrorCode               = -4
	WriteMockupOrgErrorMsg                = "marshal mockup org error"
	ApplyRuleOrganizationNotExistErrorMsg = "apply rule failure, organization not exists"
	ApplyRuleProjectNotExistErrorMsg      = "apply rule failure, project not exists"
	GetOrganizationByIdNotExistErrorCode  = -5
	GetOrganizationByIdNotExistErrorMsg   = "get organization fail, organization not exists"
	ApplyRuleWriteMockupDataErrorCode     = -6
	ApplyRuleWriteMockupDataErrorMsg      = "apply rule fail, fail to write mockup data"

	DataMarshalErrorCode   = -101
	DataUnmarshalErrorCode = -102
)

// constant for service authentication
const (
	ServiceStatusOK                   = 0
	UserAuthStatusOK                  = ServiceStatusOK
	UserNotPermitted                  = -1
	UserRoleAdministrator             = "administrator"
	UserRoleNonAdministrator          = "non-administrator"
	TokenEmptyErrorCode               = -2
	TokenEmptyErrorMsg                = "authentication failure,token not found in request"
	UserAuthErrorCode                 = -3
	UserAuthErrorMsg                  = "authentication failure,fail to verify token"
	UserAuthTokenNotExistCode         = -4
	UserAuthTokenNotExistMsg          = "authentication failure, token not exist in session"
	UserAuthErrorNonAdministratorCode = -5
	UserAuthErrorNonAdministratorMsg  = "authentication failure, role of user's not administrator"

	ServiceInternalErrorCode = -6
	ServiceInternalErrorMsg  = "internal error"
)

const (
	KeyContextId = "context-id"
	KeyUserId    = "userId"
	KeyCode      = "code"
	KeyError     = "error"
	KeyUserRole  = "userRole"
	KeyMessage   = "msg"

	KeyToken            = "token"
	KeyUserOrganization = "organization"

	ContextKeyUserInfo = "userInfo"

	LogKeyContextId = "context-id"
	LogKeyCode      = "code"
	LogKeyError     = "error"
	LogKeyMessage   = "msg"
	LogKeyOrgId     = "org-id"
	LogKeyProjId    = "proj-id"
	LogKeyPath      = "path"
	LogKeyUserId    = "user-id"
	LogKeyLintId    = "lint-id"

	HTTPRespKeyCode     = "code"
	HTTPRespKeyMessage  = "msg"
	HTTPRespKeyUserId   = "user-id"
	HTTPRespKeyUserRole = "user-role"
	HTTPRespKeyUserOrg  = "user-organization"
	HTTPRespKeyUserFile = "file-content"

	HTTPReqKeyUserInPath   = "userid"
	HTTPReqKeyLintInPath   = "lintid"
	HTTPReqKeyProjIdInForm = "projId"
	HTTPReqKeyContent      = "content"

	HTTPRespMessageParamMissing  = "param missing : [%s]"
	HTTPRespCodeParamMissing     = -10001
	HTTPRespMessageInternalError = "internal error"
	HTTPRespCodeInternalError    = -10002
	HTTPRespMessageProjNotFound  = "proj not found"
	HTTPRespCodeProjNotFound     = -10003
)
