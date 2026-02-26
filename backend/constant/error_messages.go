package constant

const (
	ErrMsgInternalServer  = "An internal server error occurred"
	ErrMsgBadRequest      = "Invalid request"
	ErrMsgUnauthorized    = "Unauthorized access"
	ErrMsgForbidden       = "Access forbidden"
	ErrMsgNotFound        = "Resource not found"
	ErrMsgValidationError = "Validation failed"

	// User (extended)
	ErrMsgPhoneAlreadyExists = "Phone number is already in use"

	// User
	ErrMsgUserNotFound       = "User not found"
	ErrMsgEmailAlreadyExists = "Email already exists"
	ErrMsgInvalidPassword    = "Invalid password"
	ErrMsgUserInactive       = "User account is inactive"
	ErrMsgUserDeleted        = "User account has been deleted"
	ErrMsgSamePassword       = "New password cannot be the same as old password"

	// Auth
	ErrMsgInvalidCredentials     = "Invalid email or password"
	ErrMsgWeakPassword           = "Password must be at least 8 characters and meet at least 3 of: uppercase letter, lowercase letter, number, special character"
	ErrMsgInvalidRefreshToken    = "Invalid or expired refresh token"
	ErrMsgInvalidResetToken      = "Invalid or expired password reset token"
	ErrMsgAccountPendingDeletion = "Account is scheduled for deletion. Please confirm reactivation to continue."

	// Session
	ErrMsgMaxSessionsExceeded = "Maximum number of active sessions reached"
	ErrMsgSessionNotFound     = "Session not found"

	// Email Verification
	ErrMsgInvalidVerificationToken = "Invalid or expired verification token"
	ErrMsgEmailAlreadyVerified     = "Email is already verified"
	ErrMsgEmailSendFailed          = "Failed to send verification email"

	// OAuth
	ErrMsgOAuthDisabled      = "Google authentication is not enabled"
	ErrMsgOAuthInvalidState  = "Invalid or expired OAuth state"
	ErrMsgOAuthStateMismatch = "OAuth state validation failed"
	ErrMsgOAuthTokenExchange = "Failed to exchange authorization code"
	ErrMsgOAuthInvalidToken  = "Invalid OAuth token"
	ErrMsgOAuthEmailNotFound = "Email not provided by OAuth provider"
	ErrMsgOAuthProviderError = "OAuth provider returned an error"
	ErrMsgOAuthAccountLinked    = "This Google account is already linked to another user"
	ErrMsgOAuthLoginNotLinked   = "This OAuth account is not linked to any existing account. Please sign up first."

	// GitHub OAuth
	ErrMsgGitHubOAuthDisabled = "GitHub authentication is not enabled"
	ErrMsgGitHubAccountLinked = "This GitHub account is already linked to another user"

	// Facebook OAuth
	ErrMsgFacebookOAuthDisabled = "Facebook authentication is not enabled"
	ErrMsgFacebookAccountLinked = "This Facebook account is already linked to another user"

	// Rate Limiting
	ErrMsgAccountRateLimited = "Too many failed login attempts for this account. Please try again later."
	ErrMsgIPRateLimited      = "Too many failed login attempts from this IP address. Please try again later."

	// Links
	ErrMsgLinkNotFound     = "Link not found or you do not have access"
	ErrMsgSlugTaken        = "This custom slug is already taken. Please try another."
	ErrMsgLinkExpired      = "This link has expired"
	ErrMsgLinkDisabled     = "This link has been disabled"
	ErrMsgInvalidURL       = "The destination URL is invalid. Please provide a valid URL."
	ErrMsgPlanLimitReached = "You have reached your plan's link limit. Please upgrade to create more links."
)
