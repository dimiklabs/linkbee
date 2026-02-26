<template>
  <div class="auth-wrapper">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-12" style="max-width: 420px;">

          <!-- Logo -->
          <div class="text-center mb-4">
            <div class="auth-logo">🔗 Shortlink</div>
            <p class="text-muted small mt-1">Create your account</p>
          </div>

          <!-- Success state -->
          <div v-if="successState" class="card shadow-sm border-0 auth-card">
            <div class="card-body p-4 text-center">
              <div class="mb-3" style="font-size: 3rem;">✅</div>
              <h5 class="fw-semibold mb-2">Check your email!</h5>
              <p class="text-muted small mb-4">
                We've sent a verification link to <strong>{{ registeredEmail }}</strong>.
                Click the link in the email to verify your account and get started.
              </p>
              <router-link to="/auth/login" class="btn btn-primary px-4">
                Go to Sign In
              </router-link>
            </div>
          </div>

          <!-- Signup Card -->
          <div v-else class="card shadow-sm border-0 auth-card">
            <div class="card-body p-4">

              <!-- Error Alert -->
              <div v-if="errorMessage" class="alert alert-danger alert-dismissible py-2 small" role="alert">
                {{ errorMessage }}
                <button type="button" class="btn-close btn-close-sm" @click="errorMessage = ''"></button>
              </div>

              <!-- Signup Form -->
              <form @submit.prevent="handleSignup" novalidate>

                <!-- First Name + Last Name row -->
                <div class="row g-2 mb-3">
                  <div class="col-6">
                    <label for="firstName" class="form-label">First Name</label>
                    <input
                      id="firstName"
                      v-model="form.firstName"
                      type="text"
                      class="form-control"
                      :class="{ 'is-invalid': errors.firstName }"
                      placeholder="John"
                      autocomplete="given-name"
                    />
                    <div v-if="errors.firstName" class="invalid-feedback">{{ errors.firstName }}</div>
                  </div>
                  <div class="col-6">
                    <label for="lastName" class="form-label">Last Name</label>
                    <input
                      id="lastName"
                      v-model="form.lastName"
                      type="text"
                      class="form-control"
                      :class="{ 'is-invalid': errors.lastName }"
                      placeholder="Doe"
                      autocomplete="family-name"
                    />
                    <div v-if="errors.lastName" class="invalid-feedback">{{ errors.lastName }}</div>
                  </div>
                </div>

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
                  <label for="password" class="form-label">Password</label>
                  <div class="input-group">
                    <input
                      id="password"
                      v-model="form.password"
                      :type="showPassword ? 'text' : 'password'"
                      class="form-control"
                      :class="{ 'is-invalid': errors.password }"
                      placeholder="Min. 8 characters"
                      autocomplete="new-password"
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
                  <div class="form-text text-muted" style="font-size: 0.75rem;">
                    Must be at least 8 characters long.
                  </div>
                </div>

                <div class="mb-4">
                  <label for="confirmPassword" class="form-label">Confirm Password</label>
                  <div class="input-group">
                    <input
                      id="confirmPassword"
                      v-model="form.confirmPassword"
                      :type="showConfirmPassword ? 'text' : 'password'"
                      class="form-control"
                      :class="{ 'is-invalid': errors.confirmPassword }"
                      placeholder="••••••••"
                      autocomplete="new-password"
                      required
                    />
                    <button
                      type="button"
                      class="btn btn-outline-secondary"
                      @click="showConfirmPassword = !showConfirmPassword"
                      tabindex="-1"
                    >
                      {{ showConfirmPassword ? '🙈' : '👁️' }}
                    </button>
                    <div v-if="errors.confirmPassword" class="invalid-feedback">{{ errors.confirmPassword }}</div>
                  </div>
                </div>

                <button
                  type="submit"
                  class="btn btn-primary w-100"
                  :disabled="loading"
                >
                  <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                  Create Account
                </button>
              </form>

              <!-- Divider -->
              <div class="auth-divider my-4">
                <span class="text-muted small">Or sign up with</span>
              </div>

              <!-- OAuth Buttons -->
              <div class="d-flex gap-2">
                <button
                  type="button"
                  class="btn btn-outline-secondary flex-fill oauth-btn"
                  :disabled="oauthLoading"
                  @click="handleOAuth('google')"
                >
                  <span class="fw-bold" style="font-size: 0.75rem;">G</span>
                  <span class="d-none d-sm-inline ms-1">Google</span>
                </button>
                <button
                  type="button"
                  class="btn btn-outline-secondary flex-fill oauth-btn"
                  :disabled="oauthLoading"
                  @click="handleOAuth('github')"
                >
                  <span class="fw-bold" style="font-size: 0.75rem;">GH</span>
                  <span class="d-none d-sm-inline ms-1">GitHub</span>
                </button>
                <button
                  type="button"
                  class="btn btn-outline-secondary flex-fill oauth-btn"
                  :disabled="oauthLoading"
                  @click="handleOAuth('facebook')"
                >
                  <span class="fw-bold" style="font-size: 0.75rem;">FB</span>
                  <span class="d-none d-sm-inline ms-1">Facebook</span>
                </button>
              </div>

            </div>
          </div>

          <!-- Sign in link -->
          <p class="text-center mt-4 text-muted small">
            Already have an account?
            <router-link to="/auth/login" class="text-decoration-none fw-medium" style="color: #635bff;">Sign in</router-link>
          </p>

        </div>
      </div>
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
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  padding: 0.5rem 0.5rem;
}
</style>
