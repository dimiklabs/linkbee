export interface SignupRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
  remember_me?: boolean;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  phone?: string;
  profile_picture?: string;
}

export interface ReactivateAccountRequest {
  email: string;
  password: string;
}

export interface LogoutRequest {
  refresh_token?: string;
}

export interface VerifyEmailRequest {
  token: string;
}

export interface ResendVerificationRequest {
  email: string;
}

// Response types
export interface User {
  id: string;
  email: string;
  role: string;
}

export interface SignupResponse {
  id: string;
  email: string;
  status: string;
  role: string;
  created_at: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface ProfileResponse {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  profile_picture?: string;
}

export interface SessionResponse {
  id: string;
  user_agent: string;
  ip_address: string;
  created_at: string;
  last_activity_at: string;
  is_current: boolean;
  device_name?: string;
  device_type?: string;
  browser?: string;
  os?: string;
  login_method?: string;
}

export interface SessionsListResponse {
  sessions: SessionResponse[];
  count: number;
  max_allowed: number;
}

// API Response wrapper
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  description?: string;
  error_code?: string;
}
