<template>
  <div class="auth-wrapper">
    <div class="auth-container">

      <!-- Logo -->
      <div class="auth-logo-wrap">
        <div class="logo-icon">S</div>
        <span class="logo-text">Shortlink</span>
      </div>

      <!-- Card -->
      <div class="m3-card m3-card--elevated auth-card">

        <h1 class="md-headline-small auth-heading">Sign in</h1>

        <!-- Error Banner -->
        <div v-if="errorMessage" class="error-banner">
          <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">error</span>
          <span class="md-body-medium" style="flex:1;">{{ errorMessage }}</span>
          <md-icon-button @click="errorMessage = ''">
            <span class="material-symbols-outlined">close</span>
          </md-icon-button>
        </div>

        <!-- TOTP Step -->
        <div v-if="pendingTOTPSession">
          <p class="md-body-medium totp-hint">
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
                style="width: 100%; text-align: center;"
              />
            </div>
            <md-filled-button type="submit" :disabled="loading" style="width: 100%; margin-bottom: 8px;">
              <md-circular-progress v-if="loading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
              Verify &amp; Sign In
            </md-filled-button>
            <md-text-button type="button" @click="pendingTOTPSession = ''" style="width: 100%;">
              <span class="material-symbols-outlined" style="font-size:18px; margin-right:4px;">arrow_back</span>
              Back to login
            </md-text-button>
          </form>
        </div>

        <!-- Login Form -->
        <template v-else>
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
                style="width: 100%;"
              />
            </div>

            <div class="field-wrap">
              <div class="password-label-row">
                <span class="md-label-large" style="color: var(--md-sys-color-on-surface-variant);">Password</span>
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
                  style="width: 100%;"
                >
                  <md-icon-button slot="trailing-icon" type="button" @click="showPassword = !showPassword" tabindex="-1">
                    <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
                  </md-icon-button>
                </md-outlined-text-field>
              </div>
            </div>

            <div class="remember-row">
              <label class="remember-label">
                <md-checkbox
                  :checked="form.rememberMe"
                  @change="form.rememberMe = ($event.target as HTMLInputElement).checked"
                />
                <span class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant);">Remember me</span>
              </label>
            </div>

            <md-filled-button type="submit" :disabled="loading" style="width: 100%;">
              <md-circular-progress v-if="loading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
              Sign In
            </md-filled-button>
          </form>

          <!-- Divider -->
          <div class="auth-divider">
            <md-divider />
            <span class="divider-label md-body-small">Or continue with</span>
            <md-divider />
          </div>

          <!-- OAuth Buttons -->
          <div class="oauth-row">
            <md-outlined-button
              :disabled="oauthLoading"
              @click="handleOAuth('google')"
              style="flex: 1;"
            >
              <span class="oauth-icon">G</span>
              <span class="oauth-label">Google</span>
            </md-outlined-button>
            <md-outlined-button
              :disabled="oauthLoading"
              @click="handleOAuth('github')"
              style="flex: 1;"
            >
              <span class="oauth-icon">GH</span>
              <span class="oauth-label">GitHub</span>
            </md-outlined-button>
            <md-outlined-button
              :disabled="oauthLoading"
              @click="handleOAuth('facebook')"
              style="flex: 1;"
            >
              <span class="oauth-icon">FB</span>
              <span class="oauth-label">Facebook</span>
            </md-outlined-button>
          </div>
        </template>

      </div>

      <!-- Sign up link -->
      <p class="auth-footer-text md-body-medium">
        Don't have an account?
        <router-link to="/auth/signup" class="auth-link">Sign up</router-link>
      </p>

    </div>

    <!-- Reactivation Dialog -->
    <md-dialog :open="showReactivationModal" @closed="showReactivationModal = false">
      <div slot="headline">Reactivate Your Account</div>
      <div slot="content" style="width: 400px; max-width: 100%;">
        <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 20px;">
          Your account has been deactivated. Enter your credentials below to reactivate it and restore full access.
        </p>

        <div v-if="reactivationError" class="error-banner" style="margin-bottom: 16px;">
          <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">error</span>
          <span class="md-body-medium">{{ reactivationError }}</span>
        </div>

        <div class="field-wrap">
          <md-outlined-text-field
            :value="reactivationForm.email"
            @input="reactivationForm.email = ($event.target as HTMLInputElement).value"
            label="Email address"
            type="email"
            style="width: 100%;"
          />
        </div>
        <div class="field-wrap">
          <md-outlined-text-field
            :value="reactivationForm.password"
            @input="reactivationForm.password = ($event.target as HTMLInputElement).value"
            label="Password"
            type="password"
            style="width: 100%;"
          />
        </div>
      </div>
      <div slot="actions">
        <md-text-button @click="showReactivationModal = false">Cancel</md-text-button>
        <md-filled-button :disabled="reactivationLoading" @click="handleReactivation">
          <md-circular-progress v-if="reactivationLoading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
          Reactivate Account
        </md-filled-button>
      </div>
    </md-dialog>

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
  justify-content: center;
  background: var(--md-sys-color-background);
  padding: 32px 16px;
}

.auth-container {
  width: 100%;
  max-width: 400px;
}

.auth-logo-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 24px;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-weight: 700;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
}

.logo-text {
  font-weight: 700;
  font-size: 1.25rem;
  color: var(--md-sys-color-on-surface);
}

.auth-card {
  padding: 32px;
  border-radius: 12px;
  background: var(--md-sys-color-surface-container-low);
}

.auth-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 24px;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 20px;
}

.field-wrap {
  margin-bottom: 16px;
}

.password-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.forgot-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;

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
}

.oauth-icon {
  font-weight: 700;
  font-size: 0.75rem;
}

.oauth-label {
  margin-left: 4px;
}

.totp-hint {
  color: var(--md-sys-color-on-surface-variant);
  text-align: center;
  margin-bottom: 20px;
}

.auth-footer-text {
  text-align: center;
  margin-top: 20px;
  color: var(--md-sys-color-on-surface-variant);
}

.auth-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  font-weight: 500;

  &:hover {
    text-decoration: underline;
  }
}
</style>
