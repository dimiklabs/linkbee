<template>
  <div class="auth-wrapper">
    <div class="auth-container">

      <!-- Logo -->
      <div class="auth-logo-wrap">
        <div class="logo-icon">S</div>
        <span class="logo-text">Shortlink</span>
      </div>

      <div class="m3-card m3-card--elevated auth-card">

        <!-- Success state -->
        <div v-if="successState" class="success-state">
          <span class="material-symbols-outlined success-icon">mark_email_read</span>
          <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 12px;">Check your inbox</h2>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 28px;">
            We've sent a password reset link to <strong>{{ submittedEmail }}</strong>.
            The link expires in 1 hour.
          </p>
          <router-link to="/auth/login" style="text-decoration: none;">
            <md-outlined-button>
              <span class="material-symbols-outlined" style="font-size:18px; margin-right:4px;">arrow_back</span>
              Back to Sign In
            </md-outlined-button>
          </router-link>
        </div>

        <!-- Form -->
        <template v-else>
          <h1 class="md-headline-small auth-heading">Reset your password</h1>
          <p class="md-body-medium auth-subtext">
            Enter your email address and we'll send you a link to reset your password.
          </p>

          <!-- Error Banner -->
          <div v-if="errorMessage" class="error-banner">
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

          <div style="text-align: center;">
            <router-link to="/auth/login" class="back-link md-label-large">
              <span class="material-symbols-outlined" style="font-size:18px; vertical-align: middle; margin-right:4px;">arrow_back</span>
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
  margin-bottom: 8px;
}

.auth-subtext {
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 24px;
}

.success-state {
  text-align: center;
}

.success-icon {
  font-size: 56px;
  color: var(--md-sys-color-primary);
  display: block;
  margin-bottom: 16px;
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

.back-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}
</style>
