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
          <h2 class="md-headline-small success-heading">Check your inbox</h2>
          <p class="md-body-medium success-text">
            We've sent a password reset link to <strong>{{ submittedEmail }}</strong>.
            The link expires in 1 hour.
          </p>
          <router-link to="/login" style="text-decoration: none;">
            <md-outlined-button>
              <span class="material-symbols-outlined" style="font-size:18px; margin-right:4px;">arrow_back</span>
              Back to Sign In
            </md-outlined-button>
          </router-link>
        </div>

        <!-- Form -->
        <template v-else>
          <h1 class="form-heading md-headline-small">Reset your password</h1>
          <p class="form-subtext md-body-medium">
            Enter your email address and we'll send you a link to reset your password.
          </p>

          <!-- Error Banner -->
          <div v-if="errorMessage" class="m3-error-banner error-banner-anim">
            <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">error</span>
            <span class="md-body-medium" style="flex:1;">{{ errorMessage }}</span>
            <md-icon-button @click="errorMessage = ''">
              <span class="material-symbols-outlined">close</span>
            </md-icon-button>
          </div>

          <form @submit.prevent="handleSubmit" novalidate>
            <div class="field-wrap">
              <md-outlined-text-field
                :value="email"
                @input="email = ($event.target as HTMLInputElement).value"
                label="Email address"
                type="email"
                autocomplete="email"
                :error="!!emailError"
                :error-text="emailError"
                style="width: 100%;"
              />
            </div>

            <md-filled-button type="submit" :disabled="loading" style="width: 100%; margin-bottom: 16px;">
              <md-circular-progress v-if="loading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
              Send Reset Link
            </md-filled-button>
          </form>

          <div class="back-row">
            <router-link to="/login" class="back-link md-label-large">
              <span class="material-symbols-outlined back-icon">arrow_back</span>
              Back to Sign In
            </router-link>
          </div>
        </template>

      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { authApi } from '@/api/auth';

const email = ref('');
const emailError = ref('');
const loading = ref(false);
const errorMessage = ref('');
const successState = ref(false);
const submittedEmail = ref('');

function validate(): boolean {
  emailError.value = '';
  if (!email.value) {
    emailError.value = 'Email is required.';
  } else if (!/\S+@\S+\.\S+/.test(email.value)) {
    emailError.value = 'Please enter a valid email address.';
  }
  return !emailError.value;
}

async function handleSubmit() {
  if (!validate()) return;

  loading.value = true;
  errorMessage.value = '';

  try {
    await authApi.forgotPassword({ email: email.value });
    submittedEmail.value = email.value;
    successState.value = true;
  } catch (err: any) {
    const data = err?.response?.data;
    errorMessage.value = data?.message || data?.description || 'Something went wrong. Please try again.';
  } finally {
    loading.value = false;
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
  max-width: 420px;
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

/* ── Form elements ────────────────────────────────────────────── */
.form-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.form-subtext {
  color: var(--md-sys-color-on-surface-variant);
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

.field-wrap {
  margin-bottom: 16px;
}

.back-row {
  text-align: center;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--md-sys-color-primary);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;

  &:hover {
    text-decoration: underline;
  }
}

.back-icon {
  font-size: 18px;
  vertical-align: middle;
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
    padding: 32px 20px;
    align-items: flex-start;
    padding-top: 48px;
  }

  .right-panel-inner {
    margin: 0 auto;
  }

  .mobile-logo {
    display: flex;
  }
}
</style>
