<template>
  <div class="auth-wrapper">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-12" style="max-width: 420px;">

          <!-- Logo -->
          <div class="text-center mb-4">
            <div class="auth-logo">🔗 Shortlink</div>
            <p class="text-muted small mt-1">Sign in to your account</p>
          </div>

          <!-- Card -->
          <div class="card shadow-sm border-0 auth-card">
            <div class="card-body p-4">

              <!-- Error Alert -->
              <div v-if="errorMessage" class="alert alert-danger alert-dismissible py-2 small" role="alert">
                {{ errorMessage }}
                <button type="button" class="btn-close btn-close-sm" @click="errorMessage = ''"></button>
              </div>

              <!-- TOTP Step -->
              <div v-if="pendingTOTPSession">
                <p class="text-center text-muted small mb-3">
                  Enter the 6-digit code from your authenticator app, or a backup code.
                </p>
                <form @submit.prevent="handleTOTPVerify" novalidate>
                  <div class="mb-4">
                    <label for="totpCode" class="form-label">Authentication code</label>
                    <input
                      id="totpCode"
                      v-model="totpCode"
                      type="text"
                      class="form-control form-control-lg text-center"
                      placeholder="000000"
                      maxlength="8"
                      autocomplete="one-time-code"
                      autofocus
                    />
                  </div>
                  <button type="submit" class="btn btn-primary w-100" :disabled="loading">
                    <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                    Verify &amp; Sign In
                  </button>
                  <button type="button" class="btn btn-link w-100 mt-2 small text-muted" @click="pendingTOTPSession = ''">
                    ← Back to login
                  </button>
                </form>
              </div>

              <!-- Login Form -->
              <template v-else>
                <form @submit.prevent="handleLogin" novalidate>
                  <div class="mb-3">
                    <label for="email" class="form-label">Email address</label>
                    <input
                      id="email"
                      v-model="form.email"
                      type="email"
                      class="form-control"
                      :class="{ 'is-invalid': errors.email }"
                      placeholder="you@example.com"
                      autocomplete="email"
                      required
                    />
                    <div v-if="errors.email" class="invalid-feedback">{{ errors.email }}</div>
                  </div>

                  <div class="mb-3">
                    <div class="d-flex justify-content-between align-items-center mb-1">
                      <label for="password" class="form-label mb-0">Password</label>
                      <router-link to="/forgot-password" class="text-decoration-none small" style="color: #635bff;">
                        Forgot password?
                      </router-link>
                    </div>
                    <div class="input-group">
                      <input
                        id="password"
                        v-model="form.password"
                        :type="showPassword ? 'text' : 'password'"
                        class="form-control"
                        :class="{ 'is-invalid': errors.password }"
                        placeholder="••••••••"
                        autocomplete="current-password"
                        required
                      />
                      <button
                        type="button"
                        class="btn btn-outline-secondary"
                        @click="showPassword = !showPassword"
                        tabindex="-1"
                      >
                        {{ showPassword ? '🙈' : '👁️' }}
                      </button>
                      <div v-if="errors.password" class="invalid-feedback">{{ errors.password }}</div>
                    </div>
                  </div>

                  <div class="mb-4">
                    <div class="form-check">
                      <input
                        id="rememberMe"
                        v-model="form.rememberMe"
                        type="checkbox"
                        class="form-check-input"
                      />
                      <label for="rememberMe" class="form-check-label small text-muted">Remember me</label>
                    </div>
                  </div>

                  <button
                    type="submit"
                    class="btn btn-primary w-100"
                    :disabled="loading"
                  >
                    <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                    Sign In
                  </button>
                </form>

                <!-- Divider -->
                <div class="auth-divider my-4">
                  <span class="text-muted small">Or continue with</span>
                </div>

                <!-- OAuth Buttons -->
                <div class="d-flex gap-2">
                  <button
                    type="button"
                    class="btn btn-outline-secondary flex-fill oauth-btn"
                    :disabled="oauthLoading"
                    @click="handleOAuth('google')"
                  >
                    <span class="oauth-icon fw-bold">G</span>
                    <span class="d-none d-sm-inline">Google</span>
                  </button>
                  <button
                    type="button"
                    class="btn btn-outline-secondary flex-fill oauth-btn"
                    :disabled="oauthLoading"
                    @click="handleOAuth('github')"
                  >
                    <span class="oauth-icon fw-bold">GH</span>
                    <span class="d-none d-sm-inline">GitHub</span>
                  </button>
                  <button
                    type="button"
                    class="btn btn-outline-secondary flex-fill oauth-btn"
                    :disabled="oauthLoading"
                    @click="handleOAuth('facebook')"
                  >
                    <span class="oauth-icon fw-bold">FB</span>
                    <span class="d-none d-sm-inline">Facebook</span>
                  </button>
                </div>
              </template>

            </div>
          </div>

          <!-- Sign up link -->
          <p class="text-center mt-4 text-muted small">
            Don't have an account?
            <router-link to="/auth/signup" class="text-decoration-none fw-medium" style="color: #635bff;">Sign up</router-link>
          </p>

        </div>
      </div>
    </div>

    <!-- Reactivation Modal -->
    <div
      v-if="showReactivationModal"
      class="modal d-block"
      tabindex="-1"
      style="background: rgba(0,0,0,0.5);"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header border-0 pb-0">
            <h5 class="modal-title fw-semibold">Reactivate Your Account</h5>
            <button type="button" class="btn-close" @click="showReactivationModal = false"></button>
          </div>
          <div class="modal-body">
            <p class="text-muted small mb-4">
              Your account has been deactivated. Enter your credentials below to reactivate it and restore full access.
            </p>

            <div v-if="reactivationError" class="alert alert-danger py-2 small">{{ reactivationError }}</div>

            <div class="mb-3">
              <label class="form-label">Email address</label>
              <input
                v-model="reactivationForm.email"
                type="email"
                class="form-control"
                placeholder="you@example.com"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Password</label>
              <input
                v-model="reactivationForm.password"
                type="password"
                class="form-control"
                placeholder="••••••••"
              />
            </div>
          </div>
          <div class="modal-footer border-0 pt-0">
            <button
              type="button"
              class="btn btn-outline-secondary"
              @click="showReactivationModal = false"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-primary"
              :disabled="reactivationLoading"
              @click="handleReactivation"
            >
              <span
                v-if="reactivationLoading"
                class="spinner-border spinner-border-sm me-2"
                role="status"
                aria-hidden="true"
              ></span>
              Reactivate Account
            </button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import authApi, { oauthApi } from '@/api/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const loading = ref(false);
const oauthLoading = ref(false);
const errorMessage = ref('');
const showPassword = ref(false);
const showReactivationModal = ref(false);
const reactivationLoading = ref(false);
const reactivationError = ref('');

// TOTP step
const pendingTOTPSession = ref('');
const totpCode = ref('');

const form = reactive({
  email: '',
  password: '',
  rememberMe: false,
});

const errors = reactive({
  email: '',
  password: '',
});

const reactivationForm = reactive({
  email: '',
  password: '',
});

function validateForm(): boolean {
  errors.email = '';
  errors.password = '';

  if (!form.email) {
    errors.email = 'Email is required.';
  } else if (!/\S+@\S+\.\S+/.test(form.email)) {
    errors.email = 'Please enter a valid email address.';
  }

  if (!form.password) {
    errors.password = 'Password is required.';
  }

  return !errors.email && !errors.password;
}

async function handleLogin() {
  if (!validateForm()) return;

  loading.value = true;
  errorMessage.value = '';

  try {
    const result = await authStore.login({
      email: form.email,
      password: form.password,
      remember_me: form.rememberMe,
    });

    if (result?.requiresTOTP) {
      pendingTOTPSession.value = result.totpSession ?? '';
      totpCode.value = '';
      return;
    }

    const redirect = (route.query.redirect as string) || '/dashboard/links';
    router.push(redirect);
  } catch (err: any) {
    const data = err?.response?.data;
    const errorCode = data?.error_code;
    if (errorCode === 'ACCOUNT_DEACTIVATED') {
      reactivationForm.email = form.email;
      reactivationForm.password = form.password;
      showReactivationModal.value = true;
    } else {
      errorMessage.value = data?.message || data?.description || 'Invalid email or password.';
    }
  } finally {
    loading.value = false;
  }
}

async function handleTOTPVerify() {
  if (!totpCode.value.trim()) return;

  loading.value = true;
  errorMessage.value = '';

  try {
    await authStore.completeTOTPLogin(pendingTOTPSession.value, totpCode.value.trim());
    const redirect = (route.query.redirect as string) || '/dashboard/links';
    router.push(redirect);
  } catch (err: any) {
    const data = err?.response?.data;
    errorMessage.value = data?.message || data?.description || 'Invalid authentication code.';
  } finally {
    loading.value = false;
  }
}

async function handleOAuth(provider: 'google' | 'github' | 'facebook') {
  oauthLoading.value = true;
  try {
    let url = '';
    if (provider === 'google') {
      url = oauthApi.getGoogleLoginUrl();
    } else if (provider === 'github') {
      url = oauthApi.getGitHubLoginUrl();
    } else if (provider === 'facebook') {
      url = oauthApi.getFacebookLoginUrl();
    }
    if (url) {
      window.location.href = url;
    }
  } catch (err: any) {
    errorMessage.value = `Failed to initiate ${provider} login. Please try again.`;
    oauthLoading.value = false;
  }
}

async function handleReactivation() {
  if (!reactivationForm.email || !reactivationForm.password) {
    reactivationError.value = 'Email and password are required.';
    return;
  }

  reactivationLoading.value = true;
  reactivationError.value = '';

  try {
    const response = await authApi.reactivateAccount({
      email: reactivationForm.email,
      password: reactivationForm.password,
    });
    if (response.data) {
      authStore.setTokens(response.data.access_token, response.data.refresh_token);
      await authStore.fetchProfile();
    }
    showReactivationModal.value = false;
    router.push('/dashboard/links');
  } catch (err: any) {
    const data = err?.response?.data;
    reactivationError.value = data?.message || data?.description || 'Reactivation failed. Please try again.';
  } finally {
    reactivationLoading.value = false;
  }
}
</script>

<style scoped lang="scss">
.auth-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  background: linear-gradient(135deg, #f7f9fc 0%, #eef0ff 100%);
  padding: 2rem 1rem;
}

.auth-logo {
  font-size: 1.75rem;
  font-weight: 700;
  color: #635bff;
  letter-spacing: -0.5px;
}

.auth-card {
  border-radius: 12px;
}

.auth-divider {
  position: relative;
  text-align: center;

  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 0;
    right: 0;
    height: 1px;
    background-color: #dee2e6;
  }

  span {
    position: relative;
    background: #fff;
    padding: 0 1rem;
  }
}

.oauth-btn {
  font-size: 0.8rem;
  padding: 0.5rem 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.oauth-icon {
  font-size: 0.75rem;
}
</style>
