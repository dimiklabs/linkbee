package request

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password_complexity"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me,omitempty"` // Optional: for extended session duration
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password_complexity"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password_complexity"`
}

type UpdateProfileRequest struct {
	FirstName      string `json:"first_name" binding:"max=100"`
	LastName       string `json:"last_name" binding:"max=100"`
	Phone          string `json:"phone" binding:"max=20"`
	ProfilePicture string `json:"profile_picture" binding:"max=2048"`
}

type ReactivateAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type TOTPConfirmRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type TOTPDisableRequest struct {
	Password string `json:"password" binding:"required"`
}

type TOTPVerifyLoginRequest struct {
	TOTPSession string `json:"totp_session" binding:"required"`
	Code        string `json:"code" binding:"required"`
}
