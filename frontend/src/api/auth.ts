import axios from 'axios';
import apiClient from './client';
import type {
  ApiResponse,
  ChangePasswordRequest,
  ForgotPasswordRequest,
  LoginRequest,
  LoginResponse,
  LogoutRequest,
  ProfileResponse,
  ReactivateAccountRequest,
  RefreshTokenRequest,
  ResendVerificationRequest,
  ResetPasswordRequest,
  SessionsListResponse,
  SignupRequest,
  SignupResponse,
  TOTPBackupCodesResponse,
  TOTPSetupResponse,
  TOTPStatusResponse,
  TOTPVerifyLoginRequest,
  UpdateProfileRequest,
  VerifyEmailRequest,
} from '@/types/auth';

export const authApi = {
  signup: async (data: SignupRequest): Promise<ApiResponse<SignupResponse>> => {
    const response = await apiClient.post('/auth/signup', data);
    return response.data;
  },

  login: async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    const response = await apiClient.post('/auth/login', data);
    return response.data;
  },

  refreshToken: async (data: RefreshTokenRequest): Promise<ApiResponse<LoginResponse>> => {
    const response = await apiClient.post('/auth/refresh', data);
    return response.data;
  },

  forgotPassword: async (data: ForgotPasswordRequest): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.post('/auth/forgot-password', data);
    return response.data;
  },

  resetPassword: async (data: ResetPasswordRequest): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.post('/auth/reset-password', data);
    return response.data;
  },

  verifyEmail: async (data: VerifyEmailRequest): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.post('/auth/verify-email', data);
    return response.data;
  },

  resendVerification: async (data: ResendVerificationRequest): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.post('/auth/resend-verification', data);
    return response.data;
  },

  reactivateAccount: async (data: ReactivateAccountRequest): Promise<ApiResponse<LoginResponse>> => {
    const response = await apiClient.post('/auth/reactivate', data);
    return response.data;
  },

  validateSession: async (refreshToken: string): Promise<void> => {
    await axios.post(`${apiClient.defaults.baseURL}/auth/session/validate`, {
      refresh_token: refreshToken,
    });
  },

  getProfile: async (): Promise<ApiResponse<ProfileResponse>> => {
    const response = await apiClient.get('/auth/profile');
    return response.data;
  },

  updateProfile: async (data: UpdateProfileRequest): Promise<ApiResponse<ProfileResponse>> => {
    const response = await apiClient.put('/auth/profile', data);
    return response.data;
  },

  changePassword: async (data: ChangePasswordRequest): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.put('/auth/change-password', data);
    return response.data;
  },

  logout: async (data?: LogoutRequest): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.post('/auth/logout', data);
    return response.data;
  },

  deleteAccount: async (): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.delete('/auth/account');
    return response.data;
  },

  getSessions: async (): Promise<ApiResponse<SessionsListResponse>> => {
    const response = await apiClient.get('/auth/sessions');
    return response.data;
  },

  deleteSession: async (sessionId: string): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.delete(`/auth/sessions/${sessionId}`);
    return response.data;
  },

  logoutAll: async (): Promise<ApiResponse<{ message: string }>> => {
    const response = await apiClient.post('/auth/logout-all');
    return response.data;
  },

  getTOTPStatus: async (): Promise<ApiResponse<TOTPStatusResponse>> => {
    const response = await apiClient.get('/auth/totp/status');
    return response.data;
  },

  setupTOTP: async (): Promise<ApiResponse<TOTPSetupResponse>> => {
    const response = await apiClient.get('/auth/totp/setup');
    return response.data;
  },

  confirmTOTP: async (code: string): Promise<ApiResponse<TOTPBackupCodesResponse>> => {
    const response = await apiClient.post('/auth/totp/confirm', { code });
    return response.data;
  },

  disableTOTP: async (password: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.delete('/auth/totp/disable', { data: { password } });
    return response.data;
  },

  verifyTOTPLogin: async (data: TOTPVerifyLoginRequest): Promise<ApiResponse<LoginResponse>> => {
    const response = await apiClient.post('/auth/totp/verify-login', data);
    return response.data;
  },

  downloadExport: async (): Promise<void> => {
    const response = await apiClient.get('/auth/data-export', { responseType: 'blob' });
    const blob = new Blob([response.data], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    const date = new Date().toISOString().slice(0, 10);
    a.href = url;
    a.download = `shortlink-data-${date}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  },
};

export const oauthApi = {
  getGoogleLoginUrl: () => `${apiClient.defaults.baseURL}/auth/google`,
  getGitHubLoginUrl: () => `${apiClient.defaults.baseURL}/auth/github`,
};

export default authApi;
