<template>
  <div class="auth-wrapper">
    <div class="auth-container">

      <!-- Logo -->
      <div class="auth-logo-wrap">
        <div class="logo-icon">S</div>
        <span class="logo-text">Shortlink</span>
      </div>

      <div class="m3-card m3-card--elevated auth-card">

        <!-- Loading state -->
        <div v-if="verifying" class="state-block">
          <md-circular-progress indeterminate style="--md-circular-progress-size:56px; margin-bottom:24px;" />
          <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 8px;">Verifying your email...</h2>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant);">Please wait a moment.</p>
        </div>

        <!-- Success state -->
        <div v-else-if="verified" class="state-block">
          <span class="material-symbols-outlined state-icon state-icon--success">check_circle</span>
          <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 12px;">Email Verified!</h2>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 28px;">
            Your email has been successfully verified. You can now sign in to your Shortlink account.
          </p>
          <router-link to="/auth/login" style="text-decoration: none;">
            <md-filled-button>Sign In to Shortlink</md-filled-button>
          </router-link>
        </div>

        <!-- Error state -->
        <div v-else class="error-state">
          <span class="material-symbols-outlined state-icon state-icon--error">cancel</span>
          <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 12px;">Verification Failed</h2>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 24px;">
            {{ errorMessage }}
          </p>

          <!-- Resend verification section -->
          <div class="resend-section">
            <p class="md-title-small resend-title">Resend verification email</p>

            <div v-if="resendSuccess" class="success-banner">
              <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">check_circle</span>
              <span class="md-body-medium">{{ resendSuccess }}</span>
            </div>
            <div v-if="resendError" class="error-banner">
              <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">error</span>
              <span class="md-body-medium">{{ resendError }}</span>
            </div>

            <div class="resend-row">
              <md-outlined-text-field
                :value="resendEmail"
                @input="resendEmail = ($event.target as HTMLInputElement).value"
                label="your@email.com"
                type="email"
                style="flex: 1;"
              />
              <md-filled-button
                :disabled="resending || !resendEmail"
                @click="handleResend"
                style="height: 56px;"
              >
                <md-circular-progress v-if="resending" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
                Resend
              </md-filled-button>
            </div>
          </div>

          <div style="text-align: center; margin-top: 20px;">
            <router-link to="/auth/login" class="back-link md-label-large">
              <span class="material-symbols-outlined" style="font-size:18px; vertical-align: middle; margin-right:4px;">arrow_back</span>
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
  max-width: 420px;
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

.state-block {
  text-align: center;
}

.state-icon {
  font-size: 56px;
  display: block;
  margin-bottom: 16px;

  &--success {
    color: var(--md-sys-color-primary);
  }

  &--error {
    color: var(--md-sys-color-error);
  }
}

.error-state {
  // no special centering needed for error state — has left-aligned resend section
}

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

.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  padding: 10px 14px;
  margin-bottom: 12px;
}

.back-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}
</style>
