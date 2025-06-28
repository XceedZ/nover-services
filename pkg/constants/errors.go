package constants

const (
	ErrCodeBadRequest     = "bad_request"
	ErrCodeInternalServer = "internal_server"

	ErrCodeAuthInvalidCredentials   = "invalid_credentials"
	ErrCodeAuthInputRequired        = "input_required"
	ErrCodeAuthEmailOrUsernameTaken = "email_or_username_taken"
	ErrCodeAuthJWTSecretMissing     = "jwt_secret_missing"
	ErrCodeAuthTokenCreation        = "token_creation"

	ErrCodeUserUnauthorized = "unauthorized"
	ErrCodeUserPenNameTaken = "pen_name_taken"
	ErrCodeUserNotFound     = "not_found"
	ErrCodeUserUpdateFailed = "update_failed"
)
