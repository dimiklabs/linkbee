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

        <!-- Loading state -->
        <div v-if="verifying" class="state-block">
          <div class="verifying-spinner">
            <md-circular-progress indeterminate class="verify-progress" />
          </div>
          <h2 class="md-headline-small state-heading">Verifying your email...</h2>
          <p class="md-body-medium state-text">Please wait a moment.</p>
        </div>

        <!-- Success state -->
        <div v-else-if="verified" class="state-block">
          <span class="material-symbols-outlined state-icon state-icon--success">check_circle</span>
          <h2 class="md-headline-small state-heading">Email Verified!</h2>
          <p class="md-body-medium state-text state-text--mb">
            Your email has been successfully verified. You can now sign in to your Shortlink account.
          </p>
          <router-link to="/login" class="state-action-link">
            <md-filled-button>Sign In to Shortlink</md-filled-button>
          </router-link>
        </div>

        <!-- Error state -->
        <div v-else class="error-content">
          <span class="material-symbols-outlined state-icon state-icon--error">cancel</span>
          <h2 class="md-headline-small state-heading">Verification Failed</h2>
          <p class="md-body-medium state-text state-text--mb-sm">
            {{ errorMessage }}
          </p>

          <!-- Resend verification section -->
          <div class="resend-section">
            <p class="md-title-small resend-title">Resend verification email</p>

            <div v-if="resendSuccess" class="success-banner">
              <span class="material-symbols-outlined resend-icon">check_circle</span>
              <span class="md-body-medium">{{ resendSuccess }}</span>
            </div>
            <div v-if="resendError" class="m3-error-banner">
              <span class="material-symbols-outlined resend-icon">error</span>
              <span class="md-body-medium">{{ resendError }}</span>
            </div>

            <div class="resend-row">
              <md-outlined-text-field
                :value="resendEmail"
                @input="resendEmail = ($event.target as HTMLInputElement).value"
                label="your@email.com"
                type="email"
                class="resend-field"
              />
              <md-filled-button
                :disabled="resending || !resendEmail"
                @click="handleResend"
                class="resend-btn"
              >
                <md-circular-progress v-if="resending" indeterminate class="btn-spinner" />
                Resend
              </md-filled-button>
            </div>
          </div>

          <div class="back-row">
            <router-link to="/login" class="back-link md-label-large">
              <span class="material-symbols-outlined back-icon">arrow_back</span>
              Back to Sign In
            </router-link>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import authApi from '@/api/auth';

const route = useRoute();

const verifying = ref(true);
const verified = ref(false);
const errorMessage = ref('');

const resendEmail = ref('');
const resending = ref(false);
const resendSuccess = ref('');
const resendError = ref('');

onMounted(async () => {
  const token = (route.query.token as string) || '';

  if (!token) {
    errorMessage.value = 'No verification token found. Please check your email link.';
    verifying.value = false;
    return;
  }

  try {
    await authApi.verifyEmail({ token });
    verified.value = true;
  } catch (err: any) {
    const data = err?.response?.data;
    errorMessage.value =
      data?.message || data?.description || 'Verification failed. The link may have expired.';
  } finally {
    verifying.value = false;
  }
});

async function handleResend() {
  if (!resendEmail.value) {
    resendError.value = 'Please enter your email address.';
    return;
  }

  resending.value = true;
  resendError.value = '';
  resendSuccess.value = '';

  try {
    await authApi.resendVerification({ email: resendEmail.value });
    resendSuccess.value = 'Verification email sent! Please check your inbox.';
    resendEmail.value = '';
  } catch (err: any) {
    const data = err?.response?.data;
    resendError.value = data?.message || data?.description || 'Failed to resend verification email.';
  } finally {
    resending.value = false;
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

/* ── State blocks ─────────────────────────────────────────────── */
.state-block {
  text-align: center;
  padding: 24px 0;
}

.verifying-spinner {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.verify-progress {
  --md-circular-progress-size: 56px;
}

.state-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;

  &--success {
    color: var(--md-sys-color-primary);
  }

  &--error {
    color: var(--md-sys-color-error);
  }
}

.state-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 12px;
}

.state-text {
  color: var(--md-sys-color-on-surface-variant);
}

.state-text--mb {
  margin-bottom: 28px;
}

.state-text--mb-sm {
  margin-bottom: 24px;
}

.state-action-link {
  text-decoration: none;
}

.error-content {
  // left-aligned for error + resend block
}

/* ── Resend section ───────────────────────────────────────────── */
.resend-section {
  background: var(--md-sys-color-surface-container);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 8px;
}

.resend-title {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 12px;
}

.resend-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  margin-top: 12px;
}

.resend-field {
  flex: 1;
}

.resend-btn {
  height: 56px;
}

.resend-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.btn-spinner {
  --md-circular-progress-size: 20px;
  margin-right: 8px;
}

.success-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
  border-radius: 8px;
  padding: 10px 14px;
  margin-bottom: 12px;
}

.m3-error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  padding: 10px 14px;
  margin-bottom: 12px;
}

/* ── Back link ────────────────────────────────────────────────── */
.back-row {
  text-align: center;
  margin-top: 20px;
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
