package constant

const (
	ErrCodeInternalServer  = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeValidationError = "VALIDATION_ERROR"
	ErrCodeConflict        = "CONFLICT"

	// User
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	ErrCodeInvalidPassword    = "INVALID_PASSWORD"
	ErrCodeUserInactive       = "USER_INACTIVE"
	ErrCodeUserDeleted        = "USER_DELETED"
	ErrCodeSamePassword       = "SAME_PASSWORD"

	// Auth
	ErrCodeInvalidCredentials     = "INVALID_CREDENTIALS"
	ErrCodeWeakPassword           = "WEAK_PASSWORD"
	ErrCodeInvalidRefreshToken    = "INVALID_REFRESH_TOKEN"
	ErrCodeInvalidResetToken      = "INVALID_RESET_TOKEN"
	ErrCodeAccountPendingDeletion = "ACCOUNT_PENDING_DELETION"

	// Session
	ErrCodeMaxSessionsExceeded = "MAX_SESSIONS_EXCEEDED"
	ErrCodeSessionNotFound     = "SESSION_NOT_FOUND"

	// Email Verification
	ErrCodeInvalidVerificationToken = "INVALID_VERIFICATION_TOKEN"
	ErrCodeEmailAlreadyVerified     = "EMAIL_ALREADY_VERIFIED"
	ErrCodeEmailSendFailed          = "EMAIL_SEND_FAILED"

	// User (extended)
	ErrCodePhoneAlreadyExists = "PHONE_ALREADY_EXISTS"

	// OAuth
	ErrCodeOAuthDisabled      = "OAUTH_DISABLED"
	ErrCodeOAuthInvalidState  = "OAUTH_INVALID_STATE"
	ErrCodeOAuthStateMismatch = "OAUTH_STATE_MISMATCH"
	ErrCodeOAuthTokenExchange = "OAUTH_TOKEN_EXCHANGE_FAILED"
	ErrCodeOAuthInvalidToken  = "OAUTH_INVALID_TOKEN"
	ErrCodeOAuthEmailNotFound = "OAUTH_EMAIL_NOT_FOUND"
	ErrCodeOAuthProviderError = "OAUTH_PROVIDER_ERROR"
	ErrCodeOAuthAccountLinked = "OAUTH_ACCOUNT_ALREADY_LINKED"
	ErrCodeOAuthNotLinked     = "OAUTH_NOT_LINKED"

	// GitHub OAuth
	ErrCodeGitHubOAuthDisabled = "GITHUB_OAUTH_DISABLED"
	ErrCodeGitHubAccountLinked = "GITHUB_ACCOUNT_ALREADY_LINKED"

	// Rate Limiting
	ErrCodeRateLimited = "RATE_LIMITED"

	// Custom Domains
	ErrCodeDomainAlreadyExists = "DOMAIN_ALREADY_EXISTS"
	ErrCodeDomainNotFound      = "DOMAIN_NOT_FOUND"
	ErrCodeDomainVerifyFailed  = "DOMAIN_VERIFY_FAILED"

	// TOTP / 2FA
	ErrCodeTOTPNotEnabled      = "TOTP_NOT_ENABLED"
	ErrCodeTOTPAlreadyEnabled  = "TOTP_ALREADY_ENABLED"
	ErrCodeTOTPInvalidCode     = "TOTP_INVALID_CODE"
	ErrCodeTOTPInvalidSession  = "TOTP_INVALID_SESSION"

	// Links
	ErrCodeLinkNotFound     = "LINK_NOT_FOUND"
	ErrCodeSlugTaken        = "SLUG_ALREADY_TAKEN"
	ErrCodeLinkExpired      = "LINK_EXPIRED"
	ErrCodeLinkDisabled     = "LINK_DISABLED"
	ErrCodeInvalidURL       = "INVALID_URL"
	ErrCodePlanLimitReached = "PLAN_LIMIT_REACHED"
)
