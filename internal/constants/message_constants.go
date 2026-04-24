package constants

type ResponseCode struct {
	Code    string
	Message string
}

var (
	// Success
	Success                  = ResponseCode{"00", "Success"}
	LoginSuccess             = ResponseCode{"INF003", "Login success"}
	RegisterUserSuccess      = ResponseCode{"INF001", "Register user success"}
	LogoutSuccess            = ResponseCode{"INF004", "Logout success"}
	SendEmailRegisterSuccess = ResponseCode{"INF009", "Send email register success"}
	EmailVerifiedSuccess     = ResponseCode{"INF061", "Email verified successfully"}
	OnboardSuccess           = ResponseCode{"INF010", "Onboarding completed successfully"}

	// General Errors
	InternalServerError  = ResponseCode{"ERR500", "Internal server error"}
	InvalidRequestParams = ResponseCode{"ERR001", "Invalid request parameters"}

	// Auth errors
	MissingAuthorizationHeader = ResponseCode{"ERR012", "Missing authorization header"}
	InvalidAuthorizationHeader = ResponseCode{"ERR013", "Invalid authorization header format"}
	InvalidToken               = ResponseCode{"ERR009", "Invalid token"}
	TokenExpired               = ResponseCode{"ERR011", "Token expired"}
	UserNotFound               = ResponseCode{"ERR004", "User not found"}
	EmailExists                = ResponseCode{"ERR008", "This email is already registered"}
	UsernameExists             = ResponseCode{"ERR015", "This username is already taken"}
	PasswordIncorrect          = ResponseCode{"ERR005", "The password is incorrect"}
	LoginFail                  = ResponseCode{"ERR006", "Login fail, something wrong"}
	RegisterUserFail           = ResponseCode{"ERR007", "Register user fail, something wrong"}
	VerifyCodeExpired          = ResponseCode{"ERR014", "Verify code expired"}
	PermissionDenied           = ResponseCode{"ERR030", "Permission denied, you can't perform this action"}

	// Chat / Room / Channel errors
	ChannelNotFound   = ResponseCode{"ERR045", "Channel not found"}
	CreateChannelFail = ResponseCode{"ERR023", "Create channel fail"}
	CreateMessageFail = ResponseCode{"ERR051", "Create message fail"}
	GetMessageFail    = ResponseCode{"ERR052", "Get message fail"}
	SendMessageWSFail = ResponseCode{"ERR050", "Send message to websocket fail"}

	// Organization errors
	OrganizationNotFound   = ResponseCode{"ERR024", "Organization not found"}
	CreateOrganizationFail = ResponseCode{"ERR023", "Create organization fail"}

	// Workplace errors
	WorkplaceNotFound   = ResponseCode{"ERR025", "Workplace not found"}
	CreateWorkplaceFail = ResponseCode{"ERR026", "Create workplace fail"}

	// RSA errors
	EncryptRSAFail    = ResponseCode{"ERR099", "Encrypt data fail"}
	DecryptRSAFail    = ResponseCode{"ERR099", "Decrypt data fail"}
	RSANotInitialized = ResponseCode{"ERR010", "RSA key pair not initialized"}
)
