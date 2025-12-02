// Package constants where all constants like error messages are stored
package constants

const (
	ErrInternalServer = "Internal server error"
	ErrUnauthorized   = "Unauthorized"
	ErrForbidden      = "Forbidden"
	ErrNotFound       = "Resource not found"
	ErrInvalidInput   = "Invalid input"
	ErrConflict       = "Resource conflict"
)

// Auth errors
const (
	ErrEmailAlreadyExists   = "Email already exists"
	ErrInvalidCredentials   = "Invalid email or password"
	ErrTokenExpired         = "Token expired"
	ErrTokenInvalid         = "Invalid token"
	ErrRefreshTokenRequired = "Refresh token required"
	ErrInvalidAuthHeader    = "Invalid Authorization header"
	ErrInvalidRefreshToken  = "Invalid refresh token"
	ErrFailedTokenGen       = "Failed to Generate new tokens"
)

// User errors
const (
	ErrUserNotFound      = "User not found"
	ErrPasswordTooWeak   = "Password does not meet complexity requirements"
	ErrInvalidLogin      = "Invalid Login details"
	ErrAccountNotApporve = "Admin account not approved, contact super admin"
)

// Validation errors
const (
	ErrValidationFailed = "Validation failed"
	ErrFailedHashPass   = "failed to hash password"
)
