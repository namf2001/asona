package constants

type ResponseCode struct {
	Code    string
	Message string
}

var (
	// General
	Success = ResponseCode{"00", "Success"}

	// Server Errors
	InternalServerError = ResponseCode{"ERR500", "Internal server error"}
	GenerateTokenFailed = ResponseCode{"ERR500", "Failed to generate token"}

	// Auth errors
	MissingAuthorizationHeader = ResponseCode{"ERR012", "Missing authorization header"}
	InvalidAuthorizationHeader = ResponseCode{"ERR013", "Invalid authorization header format"}
	InvalidToken               = ResponseCode{"ERR009", "Invalid token"}

	// Token errors
	TokenInvalid = ResponseCode{"ERR009", "Invalid token"}
	TokenExpired = ResponseCode{"ERR011", "Token expired"}

	// Request errors
	InvalidRequestParams = ResponseCode{"ERR001", "Invalid request parameters"}
	DecryptRSAFail       = ResponseCode{"ERR099", "Decrypt data fail"}
	RSANotInitialized    = ResponseCode{"ERR010", "RSA key pair not initialized"}
)
	