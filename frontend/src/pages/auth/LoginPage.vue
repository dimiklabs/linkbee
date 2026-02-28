<template>
  <div class="auth-split-page">

    <!-- LEFT PANEL -->
    <div class="auth-left-panel">
      <div class="left-panel-inner">
        <!-- Brand logo -->
        <div class="brand-logo">
          <div class="brand-logo-icon">
            <span class="material-symbols-outlined">link</span>
          </div>
          <span class="brand-logo-text">Shortlink</span>
        </div>

        <!-- Tagline -->
        <p class="brand-tagline">Shorten links. Track clicks. Grow faster.</p>

        <!-- Feature list -->
        <ul class="feature-list">
          <li class="feature-item">
            <span class="feature-check material-symbols-outlined">check_circle</span>
            <span>Create short links in seconds</span>
          </li>
          <li class="feature-item">
            <span class="feature-check material-symbols-outlined">check_circle</span>
            <span>Track every click with analytics</span>
          </li>
          <li class="feature-item">
            <span class="feature-check material-symbols-outlined">check_circle</span>
            <span>QR codes for every link</span>
          </li>
          <li class="feature-item">
            <span class="feature-check material-symbols-outlined">check_circle</span>
            <span>Password protection &amp; expiry</span>
          </li>
          <li class="feature-item">
            <span class="feature-check material-symbols-outlined">check_circle</span>
            <span>Team collaboration tools</span>
          </li>
        </ul>

        <div class="left-panel-footer">&copy; 2026 Shortlink</div>
      </div>
    </div>

    <!-- RIGHT PANEL -->
    <div class="auth-right-panel auth-page-bg">
      <div class="right-panel-inner">

        <!-- Mobile logo (shown only on mobile) -->
        <div class="mobile-logo">
          <div class="mobile-logo-icon">
            <span class="material-symbols-outlined">link</span>
          </div>
          <span class="mobile-logo-text">Shortlink</span>
        </div>

        <!-- Error Banner -->
        <div v-if="errorMessage" class="m3-error-banner error-banner-anim">
          <span class="material-symbols-outlined err-icon">error</span>
          <span class="md-body-medium err-text">{{ errorMessage }}</span>
          <button class="btn-icon" @click="errorMessage = ''">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>

        <!-- TOTP Step -->
        <div v-if="pendingTOTPSession">
          <h1 class="form-heading md-headline-small">Two-factor authentication</h1>
          <p class="form-subtext md-body-medium">
            Enter the 6-digit code from your authenticator app, or a backup code.
          </p>
          <form @submit.prevent="handleTOTPVerify" novalidate>
            <div class="field-wrap">
              <md-outlined-text-field
                :value="totpCode"
                @input="totpCode = ($event.target as HTMLInputElement).value"
                label="Authentication code"
                type="text"
                inputmode="numeric"
                maxlength="8"
                autocomplete="one-time-code"
                autofocus
                class="field-full totp-field"
              />
            </div>
            <button class="btn-filled btn-full btn-mb" type="submit" :disabled="loading" >
              <md-circular-progress v-if="loading" indeterminate class="btn-spinner" />
              Verify &amp; Sign In
            </button>
            <button class="btn-text btn-full" type="button" @click="pendingTOTPSession = ''" >
              <span class="material-symbols-outlined back-icon-sm">arrow_back</span>
              Back to login
            </button>
          </form>
        </div>

        <!-- Login Form -->
        <template v-else>
          <h1 class="form-heading md-headline-small">Welcome back</h1>

          <form @submit.prevent="handleLogin" novalidate>
            <div class="field-wrap">
              <md-outlined-text-field
                :value="form.email"
                @input="form.email = ($event.target as HTMLInputElement).value"
                label="Email address"
                type="email"
                autocomplete="email"
                :error="!!errors.email"
                :error-text="errors.email"
                class="field-full"
              />
            </div>

            <div class="field-wrap">
              <div class="password-label-row">
                <span class="md-label-large pass-label">Password</span>
                <router-link to="/forgot-password" class="forgot-link md-label-large">
                  Forgot password?
                </router-link>
              </div>
              <div class="password-field-wrap">
                <md-outlined-text-field
                  :value="form.password"
                  @input="form.password = ($event.target as HTMLInputElement).value"
                  label="Password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  :error="!!errors.password"
                  :error-text="errors.password"
                  class="field-full"
                >
                  <button class="btn-icon" slot="trailing-icon" type="button" @click="showPassword = !showPassword" tabindex="-1">
                    <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
                  </button>
                </md-outlined-text-field>
              </div>
            </div>

            <div class="remember-row">
              <label class="remember-label">
                <md-checkbox
                  :checked="form.rememberMe"
                  @change="form.rememberMe = ($event.target as HTMLInputElement).checked"
                />
                <span class="md-body-medium remember-text">Remember me</span>
              </label>
            </div>

            <button class="btn-filled btn-full" type="submit" :disabled="loading" >
              <md-circular-progress v-if="loading" indeterminate class="btn-spinner" />
              Sign In
            </button>
          </form>

          <!-- Divider -->
          <div class="auth-divider">
            <md-divider />
            <span class="divider-label md-body-small">Or continue with</span>
            <md-divider />
          </div>

          <!-- OAuth Buttons -->
          <div class="oauth-row">
            <button class="btn-outlined oauth-btn" :disabled="oauthLoading" @click="handleOAuth('google')">
              <svg class="oauth-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
              </svg>
              <span class="oauth-label">Google</span>
            </button>
            <button class="btn-outlined oauth-btn" :disabled="oauthLoading" @click="handleOAuth('github')">
              <svg class="oauth-icon oauth-icon--github" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
              </svg>
              <span class="oauth-label">GitHub</span>
            </button>
          </div>

          <!-- Sign up link -->
          <p class="auth-footer-text md-body-medium">
            Don't have an account?
            <router-link to="/signup" class="auth-link">Sign up</router-link>
          </p>
        </template>

      </div>
    </div>

    <!-- Reactivation Dialog -->
    <BaseModal v-model="showReactivationModal" size="sm">
      <template #headline>
        Reactivate Your Account
      </template>

      <div class="reactivation-body">
        <p class="md-body-medium reactivation-text">
          Your account has been deactivated. Enter your credentials below to reactivate it and restore full access.
        </p>

        <div v-if="reactivationError" class="m3-error-banner reactivation-error">
          <span class="material-symbols-outlined err-icon">error</span>
          <span class="md-body-medium">{{ reactivationError }}</span>
        </div>

        <div class="field-wrap">
          <md-outlined-text-field
            :value="reactivationForm.email"
            @input="reactivationForm.email = ($event.target as HTMLInputElement).value"
            label="Email address"
            type="email"
            class="field-full"
          />
        </div>
        <div class="field-wrap">
          <md-outlined-text-field
            :value="reactivationForm.password"
            @input="reactivationForm.password = ($event.target as HTMLInputElement).value"
            label="Password"
            type="password"
            class="field-full"
          />
        </div>
      </div>

      <template #actions>
        <button class="btn-text" @click="showReactivationModal = false">Cancel</button>
        <button class="btn-filled" :disabled="reactivationLoading" @click="handleReactivation">
          <md-circular-progress v-if="reactivationLoading" indeterminate class="btn-spinner" />
          Reactivate Account
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import authApi, { oauthApi } from '@/api/auth';
import BaseModal from '@/components/BaseModal.vue';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const OAUTH_ERROR_MESSAGES: Record<string, string> = {
  EMAIL_ALREADY_EXISTS: 'An account with this email already exists. Please sign in with your email and password.',
  OAUTH_DISABLED: 'OAuth sign-in is currently disabled. Please use email and password.',
  OAUTH_INVALID_STATE: 'The sign-in session expired or is invalid. Please try again.',
  OAUTH_TOKEN_EXCHANGE: 'Could not complete sign-in with the provider. Please try again.',
  OAUTH_PROVIDER_ERROR: 'The sign-in provider returned an error. Please try again.',
  OAUTH_EMAIL_NOT_FOUND: 'Could not retrieve your email from the provider. Please ensure your account has a verified email.',
  USER_INACTIVE: 'Your account is inactive. Please contact support.',
  FORBIDDEN: 'Access denied. Please contact support.',
};

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

onMounted(() => {
  const oauthError = route.query.oauth_error as string | undefined;
  const oauthMessage = route.query.oauth_message as string | undefined;
  if (oauthError) {
    errorMessage.value = OAUTH_ERROR_MESSAGES[oauthError] ?? oauthMessage ?? 'Sign-in failed. Please try again.';
  }
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

async function handleOAuth(provider: 'google' | 'github') {
  oauthLoading.value = true;
  try {
    const url = provider === 'google'
      ? oauthApi.getGoogleLoginUrl()
      : oauthApi.getGitHubLoginUrl();
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
      authStore.setTokens(response.data.access_token ?? '', response.data.refresh_token ?? '');
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
/* ── Split-page shell ─────────────────────────────────────────── */
.auth-split-page {
  display: flex;
  min-height: 100vh;
}

/* ── Left panel ───────────────────────────────────────────────── */
.auth-left-panel {
  width: 45%;
  background: linear-gradient(135deg, #635BFF 0%, #8B5CF6 60%, #14B8A6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  position: relative;
  overflow: hidden;

  /* Decorative blobs */
  &::before {
    content: '';
    position: absolute;
    width: 340px;
    height: 340px;
    background: rgba(255, 255, 255, 0.07);
    border-radius: 50%;
    top: -80px;
    right: -80px;
  }

  &::after {
    content: '';
    position: absolute;
    width: 240px;
    height: 240px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 50%;
    bottom: -60px;
    left: -60px;
  }
}

.left-panel-inner {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 340px;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 32px;
}

.brand-logo-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(8px);

  .material-symbols-outlined {
    color: #fff;
    font-size: 24px;
  }
}

.brand-logo-text {
  font-weight: 800;
  font-size: 1.5rem;
  color: #fff;
  letter-spacing: -0.02em;
}

.brand-tagline {
  font-size: 1.2rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.92);
  line-height: 1.4;
  margin-bottom: 40px;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0 0 48px 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.95rem;
  font-weight: 500;
}

.feature-check {
  color: rgba(255, 255, 255, 0.95);
  font-size: 20px;
  flex-shrink: 0;
}

.left-panel-footer {
  color: rgba(255, 255, 255, 0.55);
  font-size: 0.8rem;
}

/* ── Right panel ──────────────────────────────────────────────── */
.auth-right-panel {
  width: 55%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  background: var(--md-sys-color-background);
  overflow-y: auto;
}

.right-panel-inner {
  width: 100%;
  max-width: 420px;
}

/* ── Mobile logo (hidden on desktop) ─────────────────────────── */
.mobile-logo {
  display: none;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
}

.mobile-logo-icon {
  width: 36px;
  height: 36px;
  background: var(--md-sys-color-primary);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;

  .material-symbols-outlined {
    color: var(--md-sys-color-on-primary);
    font-size: 20px;
  }
}

.mobile-logo-text {
  font-weight: 700;
  font-size: 1.25rem;
  color: var(--md-sys-color-on-surface);
}

/* ── Form elements ────────────────────────────────────────────── */
.form-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.form-subtext {
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 28px;
}

.error-banner-anim {
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.m3-error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 20px;
}

.err-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.err-text {
  flex: 1;
}

.field-wrap {
  margin-bottom: 16px;
}

.field-full {
  width: 100%;
}

.totp-field {
  text-align: center;
}

.password-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.pass-label {
  color: var(--md-sys-color-on-surface-variant);
}

.forgot-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  font-size: 0.875rem;

  &:hover {
    text-decoration: underline;
  }
}

.password-field-wrap {
  position: relative;
}

.remember-row {
  margin-bottom: 20px;
}

.remember-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.remember-text {
  color: var(--md-sys-color-on-surface-variant);
}

.btn-full {
  width: 100%;
}

.btn-mb {
  margin-bottom: 8px;
}

.btn-spinner {
    margin-right: 8px;
}

.back-icon-sm {
  font-size: 18px;
  margin-right: 4px;
}

.auth-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 24px 0;

  md-divider {
    flex: 1;
  }
}

.divider-label {
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
}

.oauth-row {
  display: flex;
  gap: 8px;
  margin-bottom: 4px;
}

.oauth-btn {
  flex: 1;
  border-radius: 999px;
}

.oauth-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;

  &--github {
    fill: var(--md-sys-color-on-surface);
  }
}

.oauth-label {
  margin-left: 6px;
  font-size: 0.85rem;
}

.auth-footer-text {
  text-align: center;
  margin-top: 24px;
  color: var(--md-sys-color-on-surface-variant);
}

.auth-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  font-weight: 600;

  &:hover {
    text-decoration: underline;
  }
}

/* ── Reactivation modal body ─────────────────────────────────── */
.reactivation-body {
  width: 400px;
  max-width: 100%;
}

.reactivation-text {
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 20px;
}

.reactivation-error {
  margin-bottom: 16px;
}

/* ── Responsive: mobile ───────────────────────────────────────── */
@media (max-width: 1023px) {
  .auth-split-page {
    flex-direction: column;
  }

  .auth-left-panel {
    display: none;
  }

  .auth-right-panel {
    width: 100%;
    min-height: 100vh;
    padding: 48px 20px 32px;
    align-items: flex-start;
  }

  .right-panel-inner {
    margin: 0 auto;
  }

  .mobile-logo {
    display: flex;
  }
}
</style>
