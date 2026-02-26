import type { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import axios from 'axios';
import type { ApiResponse } from '@/types/auth';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('access_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

// Prevent concurrent refresh attempts
let refreshPromise: Promise<string> | null = null;

// Response interceptor to handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiResponse<unknown>>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
    const requestUrl = originalRequest?.url || '';

    const skipRefreshUrls = ['/auth/login'];
    const shouldSkipRefresh = skipRefreshUrls.some(url => requestUrl.endsWith(url));

    if (error.response?.status === 401 && !originalRequest._retry && !shouldSkipRefresh) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          if (!refreshPromise) {
            refreshPromise = (async () => {
              const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
                refresh_token: refreshToken,
              });
              const { access_token, refresh_token: newRefreshToken } = response.data.data;
              localStorage.setItem('access_token', access_token);
              localStorage.setItem('refresh_token', newRefreshToken);
              const { useAuthStore } = await import('@/stores/auth');
              const authStore = useAuthStore();
              authStore.setTokens(access_token, newRefreshToken);
              return access_token;
            })();
          }

          const newAccessToken = await refreshPromise;
          refreshPromise = null;

          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          }
          return apiClient(originalRequest);
        } catch (refreshError) {
          refreshPromise = null;
          const { useAuthStore } = await import('@/stores/auth');
          const authStore = useAuthStore();
          authStore.clearTokens();
          const { default: router } = await import('@/router');
          router.push('/login');
          return Promise.reject(refreshError);
        }
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
