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

        <!-- Success state -->
        <div v-if="successState" class="success-state">
          <div class="success-icon-wrap">
            <span class="material-symbols-outlined success-icon">mark_email_read</span>
          </div>
          <h2 class="md-headline-small success-heading">Check your email!</h2>
          <p class="md-body-medium success-text">
            We've sent a verification link to <strong>{{ registeredEmail }}</strong>.
            Click the link in the email to verify your account and get started.
          </p>
          <router-link to="/login" class="success-login-link">
            <button class="btn-filled">Go to Sign In</button>
          </router-link>
        </div>

        <!-- Signup Form -->
        <template v-else>
          <h1 class="form-heading md-headline-small">Create your account</h1>

          <!-- Error Banner -->
          <div v-if="errorMessage" class="m3-error-banner error-banner-anim">
            <span class="material-symbols-outlined err-icon">error</span>
            <span class="md-body-medium err-text">{{ errorMessage }}</span>
            <button class="btn-icon" @click="errorMessage = ''">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <form @submit.prevent="handleSignup" novalidate>

            <!-- First Name + Last Name row -->
            <div class="name-row">
              <div class="field-wrap name-field">
                <md-outlined-text-field
                  :value="form.firstName"
                  @input="form.firstName = ($event.target as HTMLInputElement).value"
                  label="First Name"
                  type="text"
                  autocomplete="given-name"
                  :error="!!errors.firstName"
                  :error-text="errors.firstName"
                  class="field-full"
                />
              </div>
              <div class="field-wrap name-field">
                <md-outlined-text-field
                  :value="form.lastName"
                  @input="form.lastName = ($event.target as HTMLInputElement).value"
                  label="Last Name"
                  type="text"
                  autocomplete="family-name"
                  :error="!!errors.lastName"
                  :error-text="errors.lastName"
                  class="field-full"
                />
              </div>
            </div>

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
              <md-outlined-text-field
                :value="form.password"
                @input="onPasswordInput($event)"
                label="Password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                :error="!!errors.password"
                :error-text="errors.password"
                supporting-text="Must be at least 8 characters long."
                class="field-full"
              >
                <button class="btn-icon" slot="trailing-icon" type="button" @click="showPassword = !showPassword" tabindex="-1">
                  <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
                </button>
              </md-outlined-text-field>
              <!-- Password strength bar -->
              <div v-if="form.password" class="strength-bar-wrap">
                <div class="strength-bar">
                  <div
                    v-for="i in 4"
                    :key="i"
                    class="strength-segment"
                    :class="{ 'strength-segment--filled': i <= passwordStrength, [`strength-segment--level${passwordStrength}`]: i <= passwordStrength }"
                  />
                </div>
                <span class="strength-label md-body-small">{{ strengthLabel }}</span>
              </div>
            </div>

            <div class="field-wrap">
              <md-outlined-text-field
                :value="form.confirmPassword"
                @input="form.confirmPassword = ($event.target as HTMLInputElement).value"
                label="Confirm Password"
                :type="showConfirmPassword ? 'text' : 'password'"
                autocomplete="new-password"
                :error="!!errors.confirmPassword"
                :error-text="errors.confirmPassword"
                class="field-full"
              >
                <button class="btn-icon" slot="trailing-icon" type="button" @click="showConfirmPassword = !showConfirmPassword" tabindex="-1">
                  <span class="material-symbols-outlined">{{ showConfirmPassword ? 'visibility_off' : 'visibility' }}</span>
                </button>
              </md-outlined-text-field>
            </div>

            <button class="btn-filled btn-full btn-mt" type="submit" :disabled="loading" >
              <md-circular-progress v-if="loading" indeterminate class="btn-spinner" />
              Create Account
            </button>
          </form>

          <!-- Divider -->
          <div class="auth-divider">
            <md-divider />
            <span class="divider-label md-body-small">Or sign up with</span>
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

          <!-- Sign in link -->
          <p class="auth-footer-text md-body-medium">
            Already have an account?
            <router-link to="/login" class="auth-link">Sign in</router-link>
          </p>
        </template>

      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { oauthApi } from '@/api/auth';

const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const oauthLoading = ref(false);
const errorMessage = ref('');
const successState = ref(false);
const registeredEmail = ref('');
const showPassword = ref(false);
const showConfirmPassword = ref(false);

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: '',
});

const errors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: '',
});

/* ── Password strength ─────────────────────────────────────────── */
const passwordStrength = computed((): number => {
  const p = form.password;
  if (!p) return 0;
  let score = 0;
  if (p.length >= 8) score++;
  if (p.length >= 12) score++;
  if (/[A-Z]/.test(p) && /[a-z]/.test(p)) score++;
  if (/[0-9]/.test(p) && /[^A-Za-z0-9]/.test(p)) score++;
  return Math.max(1, Math.min(4, score));
});

const strengthLabel = computed((): string => {
  const labels: Record<number, string> = { 1: 'Weak', 2: 'Fair', 3: 'Good', 4: 'Strong' };
  return labels[passwordStrength.value] ?? '';
});

function onPasswordInput(event: Event) {
  form.password = (event.target as HTMLInputElement).value;
}

function validateForm(): boolean {
  errors.firstName = '';
  errors.lastName = '';
  errors.email = '';
  errors.password = '';
  errors.confirmPassword = '';

  if (!form.firstName.trim()) errors.firstName = 'First name is required.';
  if (!form.lastName.trim()) errors.lastName = 'Last name is required.';

  if (!form.email) {
    errors.email = 'Email is required.';
  } else if (!/\S+@\S+\.\S+/.test(form.email)) {
    errors.email = 'Please enter a valid email address.';
  }

  if (!form.password) {
    errors.password = 'Password is required.';
  } else if (form.password.length < 8) {
    errors.password = 'Password must be at least 8 characters.';
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = 'Please confirm your password.';
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match.';
  }

  return (
    !errors.firstName &&
    !errors.lastName &&
    !errors.email &&
    !errors.password &&
    !errors.confirmPassword
  );
}

async function handleSignup() {
  if (!validateForm()) return;

  loading.value = true;
  errorMessage.value = '';

  try {
    await authStore.signup({
      first_name: form.firstName,
      last_name: form.lastName,
      email: form.email,
      password: form.password,
    });
    registeredEmail.value = form.email;
    successState.value = true;
  } catch (err: any) {
    const data = err?.response?.data;
    errorMessage.value = data?.message || data?.description || 'Failed to create account. Please try again.';
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
    errorMessage.value = `Failed to initiate ${provider} sign up. Please try again.`;
    oauthLoading.value = false;
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
  max-width: 440px;
}

/* ── Mobile logo ──────────────────────────────────────────────── */
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

/* ── Success state ────────────────────────────────────────────── */
.success-state {
  text-align: center;
  padding: 24px 0;
}

.success-icon-wrap {
  margin-bottom: 16px;
}

.success-icon {
  font-size: 64px;
  color: var(--md-sys-color-primary);
}

.success-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 12px;
}

.success-text {
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 28px;
}

.success-login-link {
  text-decoration: none;
}

/* ── Form elements ────────────────────────────────────────────── */
.form-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 24px;
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

.name-row {
  display: flex;
  gap: 12px;
}

.name-field {
  flex: 1;
}

.field-wrap {
  margin-bottom: 16px;
}

.field-full {
  width: 100%;
}

.btn-full {
  width: 100%;
}

.btn-mt {
  margin-top: 8px;
}

.btn-spinner {
    margin-right: 8px;
}

/* ── Password strength bar ────────────────────────────────────── */
.strength-bar-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
  padding: 0 4px;
}

.strength-bar {
  display: flex;
  gap: 4px;
  flex: 1;
}

.strength-segment {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: var(--md-sys-color-outline-variant);
  transition: background 0.25s ease;

  &--filled {
    &.strength-segment--level1 { background: #EF4444; }
    &.strength-segment--level2 { background: #F97316; }
    &.strength-segment--level3 { background: #EAB308; }
    &.strength-segment--level4 { background: #22C55E; }
  }
}

.strength-label {
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  font-size: 0.75rem;
  min-width: 40px;
}

/* ── Divider + OAuth ──────────────────────────────────────────── */
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
