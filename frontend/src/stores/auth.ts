import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import authApi from '@/api/auth';
import type { LoginRequest, ProfileResponse, SignupRequest, User } from '@/types/auth';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const profile = ref<ProfileResponse | null>(null);
  const accessToken = ref<string | null>(localStorage.getItem('access_token'));
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'));
  const loading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = computed(() => !!accessToken.value);

  const userInitials = computed(() => {
    const firstName = profile.value?.first_name;
    const lastName = profile.value?.last_name;
    const email = profile.value?.email;
    if (firstName && lastName) return `${firstName[0]}${lastName[0]}`.toUpperCase();
    if (email) return email.charAt(0).toUpperCase();
    return 'U';
  });

  const userName = computed(() => {
    const firstName = profile.value?.first_name;
    const lastName = profile.value?.last_name;
    if (firstName && lastName) return `${firstName} ${lastName}`;
    return profile.value?.email || 'User';
  });

  const setTokens = (access: string, refresh: string) => {
    accessToken.value = access;
    refreshToken.value = refresh;
    localStorage.setItem('access_token', access);
    localStorage.setItem('refresh_token', refresh);
  };

  const clearTokens = () => {
    accessToken.value = null;
    refreshToken.value = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  };

  const login = async (data: LoginRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await authApi.login(data);
      if (response.data) {
        setTokens(response.data.access_token, response.data.refresh_token);
        user.value = response.data.user;
        await fetchProfile();
        return { success: true };
      }
      throw new Error(response.message || 'Login failed');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Login failed';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const signup = async (data: SignupRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await authApi.signup(data);
      if (response.data) return response.data;
      throw new Error(response.message || 'Signup failed');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Signup failed';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const logout = async () => {
    loading.value = true;
    try {
      if (refreshToken.value) {
        await authApi.logout({ refresh_token: refreshToken.value });
      }
    } catch {
      // Ignore logout errors
    } finally {
      clearTokens();
      user.value = null;
      profile.value = null;
      loading.value = false;
    }
  };

  const fetchProfile = async () => {
    if (!accessToken.value) return;
    try {
      const response = await authApi.getProfile();
      if (response.data) profile.value = response.data;
    } catch (err) {
      console.error('Failed to fetch profile:', err);
    }
  };

  const updateProfile = async (data: { first_name?: string; last_name?: string; phone?: string }) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await authApi.updateProfile(data);
      if (response.data) {
        profile.value = response.data;
        return response.data;
      }
      throw new Error(response.message || 'Failed to update profile');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update profile';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const changePassword = async (oldPassword: string, newPassword: string) => {
    loading.value = true;
    error.value = null;
    try {
      await authApi.changePassword({ old_password: oldPassword, new_password: newPassword });
      return { success: true };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to change password';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteAccount = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await authApi.deleteAccount();
      if (response.success) {
        clearTokens();
        user.value = null;
        profile.value = null;
        return { success: true };
      }
      throw new Error(response.message || 'Failed to delete account');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete account';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const init = async () => {
    if (accessToken.value) await fetchProfile();
  };

  return {
    user, profile, accessToken, refreshToken, loading, error,
    isAuthenticated, userInitials, userName,
    login, signup, logout, fetchProfile, updateProfile, changePassword, deleteAccount, init,
    setTokens, clearTokens,
  };
});
