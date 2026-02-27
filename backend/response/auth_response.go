package response

import "time"

type SignupResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	TokenType    string       `json:"token_type,omitempty"`
	ExpiresIn    int          `json:"expires_in,omitempty"`
	User         UserResponse `json:"user,omitempty"`
	// TOTP challenge — set when 2FA is required; access_token is absent
	RequiresTOTP bool   `json:"requires_totp,omitempty"`
	TOTPSession  string `json:"totp_session,omitempty"`
}

type TOTPSetupResponse struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
}

type TOTPStatusResponse struct {
	Enabled bool `json:"enabled"`
}

type TOTPBackupCodesResponse struct {
	BackupCodes []string `json:"backup_codes"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ProfileResponse struct {
	ID                   string `json:"id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	FirstName            string `json:"first_name,omitempty"`
	LastName             string `json:"last_name,omitempty"`
	Phone                string `json:"phone,omitempty"`
	ProfilePicture       string `json:"profile_picture,omitempty"`
	ProfilePictureSource string `json:"profile_picture_source,omitempty"`
}

type SessionResponse struct {
	ID             string    `json:"id"`
	UserAgent      string    `json:"user_agent"`
	IPAddress      string    `json:"ip_address"`
	CreatedAt      time.Time `json:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	IsCurrent      bool      `json:"is_current"`
	RememberMe     bool      `json:"remember_me"`
	DeviceName     string    `json:"device_name,omitempty"`
	DeviceType     string    `json:"device_type,omitempty"`
	Browser        string    `json:"browser,omitempty"`
	OS             string    `json:"os,omitempty"`
	Location       string    `json:"location,omitempty"`
	LoginMethod    string    `json:"login_method,omitempty"`
	LastActivityIP string    `json:"last_activity_ip,omitempty"`
	ActivityCount  int64     `json:"activity_count"`
}

type SessionsListResponse struct {
	Sessions   []SessionResponse `json:"sessions"`
	Count      int               `json:"count"`
	MaxAllowed int               `json:"max_allowed"`
}

type RateLimitStatusResponse struct {
	AccountAttempts          int  `json:"account_attempts"`
	AccountMaxAttempts       int  `json:"account_max_attempts"`
	AccountRemainingAttempts int  `json:"account_remaining_attempts"`
	AccountLocked            bool `json:"account_locked"`
	AccountLockoutRemaining  int  `json:"account_lockout_remaining_seconds"`
	IPAttempts               int  `json:"ip_attempts"`
	IPMaxAttempts            int  `json:"ip_max_attempts"`
	IPRemainingAttempts      int  `json:"ip_remaining_attempts"`
	IPLocked                 bool `json:"ip_locked"`
	IPLockoutRemaining       int  `json:"ip_lockout_remaining_seconds"`
	Violations               int  `json:"violations"`
}

type NewSessionNotification struct {
	SessionID   string    `json:"session_id"`
	DeviceName  string    `json:"device_name,omitempty"`
	DeviceType  string    `json:"device_type,omitempty"`
	Browser     string    `json:"browser,omitempty"`
	OS          string    `json:"os,omitempty"`
	IPAddress   string    `json:"ip_address"`
	Location    string    `json:"location,omitempty"`
	LoginMethod string    `json:"login_method,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ConcurrentSessionsAlert struct {
	Message          string                   `json:"message"`
	NewSession       NewSessionNotification   `json:"new_session"`
	ExistingSessions []NewSessionNotification `json:"existing_sessions,omitempty"`
	TotalSessions    int                      `json:"total_sessions"`
}
