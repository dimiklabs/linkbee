<template>
  <div class="auth-wrapper">
    <div class="auth-container">

      <!-- Logo -->
      <div class="auth-logo-wrap">
        <div class="logo-icon">S</div>
        <span class="logo-text">Shortlink</span>
      </div>

      <!-- Success state -->
      <div v-if="successState" class="m3-card m3-card--elevated auth-card text-center">
        <span class="material-symbols-outlined success-icon">check_circle</span>
        <h2 class="md-headline-small" style="color: var(--md-sys-color-on-surface); margin-bottom: 12px;">Check your email!</h2>
        <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 28px;">
          We've sent a verification link to <strong>{{ registeredEmail }}</strong>.
          Click the link in the email to verify your account and get started.
        </p>
        <router-link to="/auth/login" style="text-decoration: none;">
          <md-filled-button>Go to Sign In</md-filled-button>
        </router-link>
      </div>

      <!-- Signup Card -->
      <div v-else class="m3-card m3-card--elevated auth-card">

        <h1 class="md-headline-small auth-heading">Create your account</h1>

        <!-- Error Banner -->
        <div v-if="errorMessage" class="error-banner">
          <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">error</span>
          <span class="md-body-medium" style="flex:1;">{{ errorMessage }}</span>
          <md-icon-button @click="errorMessage = ''">
            <span class="material-symbols-outlined">close</span>
          </md-icon-button>
        </div>

        <!-- Signup Form -->
        <form @submit.prevent="handleSignup" novalidate>

          <!-- First Name + Last Name row -->
          <div class="name-row">
            <div class="field-wrap" style="flex: 1;">
              <md-outlined-text-field
                :value="form.firstName"
                @input="form.firstName = ($event.target as HTMLInputElement).value"
                label="First Name"
                type="text"
                autocomplete="given-name"
                :error="!!errors.firstName"
                :error-text="errors.firstName"
                style="width: 100%;"
              />
            </div>
            <div class="field-wrap" style="flex: 1;">
              <md-outlined-text-field
                :value="form.lastName"
                @input="form.lastName = ($event.target as HTMLInputElement).value"
                label="Last Name"
                type="text"
                autocomplete="family-name"
                :error="!!errors.lastName"
                :error-text="errors.lastName"
                style="width: 100%;"
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
              style="width: 100%;"
            />
          </div>

          <div class="field-wrap">
            <md-outlined-text-field
              :value="form.password"
              @input="form.password = ($event.target as HTMLInputElement).value"
              label="Password"
              :type="showPassword ? 'text' : 'password'"
              autocomplete="new-password"
              :error="!!errors.password"
              :error-text="errors.password"
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

          <md-filled-button type="submit" :disabled="loading" style="width: 100%; margin-top: 8px;">
            <md-circular-progress v-if="loading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
            Create Account
          </md-filled-button>
        </form>

        <!-- Divider -->
        <div class="auth-divider">
          <md-divider />
          <span class="divider-label md-body-small">Or sign up with</span>
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

      </div>

      <!-- Sign in link -->
      <p class="auth-footer-text md-body-medium">
        Already have an account?
        <router-link to="/auth/login" class="auth-link">Sign in</router-link>
      </p>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
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
    errorMessage.value = `Failed to initiate ${provider} sign up. Please try again.`;
    oauthLoading.value = false;
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

.auth-heading {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 24px;
}

.text-center {
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

.name-row {
  display: flex;
  gap: 12px;
}

.field-wrap {
  margin-bottom: 16px;
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
