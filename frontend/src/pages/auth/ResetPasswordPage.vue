<template>
  <div class="auth-wrapper">
    <div class="auth-container">

      <!-- Logo -->
      <div class="auth-logo-wrap">
        <div class="logo-icon">S</div>
        <span class="logo-text">Shortlink</span>
      </div>

      <div class="m3-card m3-card--elevated auth-card">

        <!-- No token error -->
        <div v-if="!token && !loading" class="state-block">
          <span class="material-symbols-outlined state-icon state-icon--error">error</span>
          <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 12px;">Invalid reset link</h2>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 24px;">
            The link you followed is invalid or has expired.
          </p>
          <router-link to="/auth/forgot-password" style="text-decoration: none;">
            <md-filled-button>Request a new reset link</md-filled-button>
          </router-link>
        </div>

        <!-- Success state -->
        <div v-else-if="successState" class="state-block">
          <span class="material-symbols-outlined state-icon state-icon--success">check_circle</span>
          <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 12px;">Password updated!</h2>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 4px;">
            Your password has been changed successfully.
          </p>
          <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 28px;">
            Redirecting to sign in in <strong>{{ countdown }}</strong> second{{ countdown !== 1 ? 's' : '' }}...
          </p>
          <router-link to="/auth/login" style="text-decoration: none;">
            <md-filled-button>Sign In Now</md-filled-button>
          </router-link>
        </div>

        <!-- Form -->
        <template v-else-if="token">
          <h1 class="md-headline-small auth-heading">Set new password</h1>
          <p class="md-body-medium auth-subtext">Enter and confirm your new password below.</p>

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
                :value="form.newPassword"
                @input="form.newPassword = ($event.target as HTMLInputElement).value"
                label="New Password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                :error="!!errors.newPassword"
                :error-text="errors.newPassword"
                supporting-text="Must be at least 8 characters long."
                style="width: 100%;"
              >
                <md-icon-button slot="trailing-icon" type="button" @click="showPassword = !showPassword" tabindex="-1">
                  <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
                </md-icon-button>
              </md-outlined-text-field>
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
                style="width: 100%;"
              >
                <md-icon-button slot="trailing-icon" type="button" @click="showConfirmPassword = !showConfirmPassword" tabindex="-1">
                  <span class="material-symbols-outlined">{{ showConfirmPassword ? 'visibility_off' : 'visibility' }}</span>
                </md-icon-button>
              </md-outlined-text-field>
            </div>

            <md-filled-button type="submit" :disabled="loading" style="width: 100%; margin-bottom: 16px;">
              <md-circular-progress v-if="loading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
              Update Password
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
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { authApi } from '@/api/auth';

const router = useRouter();
const route = useRoute();

const token = ref('');
const loading = ref(false);
const errorMessage = ref('');
const successState = ref(false);
const countdown = ref(3);
const showPassword = ref(false);
const showConfirmPassword = ref(false);

const form = reactive({
  newPassword: '',
  confirmPassword: '',
});

const errors = reactive({
  newPassword: '',
  confirmPassword: '',
});

onMounted(() => {
  token.value = (route.query.token as string) || '';
});

function validate(): boolean {
  errors.newPassword = '';
  errors.confirmPassword = '';

  if (!form.newPassword) {
    errors.newPassword = 'Password is required.';
  } else if (form.newPassword.length < 8) {
    errors.newPassword = 'Password must be at least 8 characters.';
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = 'Please confirm your password.';
  } else if (form.newPassword !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match.';
  }

  return !errors.newPassword && !errors.confirmPassword;
}

function startCountdown() {
  const timer = setInterval(() => {
    countdown.value -= 1;
    if (countdown.value <= 0) {
      clearInterval(timer);
      router.push('/auth/login');
    }
  }, 1000);
}

async function handleSubmit() {
  if (!validate()) return;

  loading.value = true;
  errorMessage.value = '';

  try {
    await authApi.resetPassword({
      token: token.value,
      new_password: form.newPassword,
    });
    successState.value = true;
    startCountdown();
  } catch (err: any) {
    const data = err?.response?.data;
    errorMessage.value =
      data?.message || data?.description || 'Failed to reset password. The link may have expired.';
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
