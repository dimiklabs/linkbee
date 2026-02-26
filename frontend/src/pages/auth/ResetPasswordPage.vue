<template>
  <div class="auth-wrapper">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-12" style="max-width: 420px;">

          <!-- Logo -->
          <div class="text-center mb-4">
            <div class="auth-logo">🔗 Shortlink</div>
          </div>

          <div class="card shadow-sm border-0 auth-card">
            <div class="card-body p-4">

              <h2 class="h5 fw-semibold mb-1">Set new password</h2>
              <p class="text-muted small mb-4">Enter and confirm your new password below.</p>

              <!-- No token error -->
              <div v-if="!token && !loading" class="alert alert-danger" role="alert">
                <strong>Invalid reset link.</strong>
                The link you followed is invalid or has expired.
                <div class="mt-2">
                  <router-link to="/auth/forgot-password" class="alert-link">Request a new reset link</router-link>
                </div>
              </div>

              <!-- Success state -->
              <div v-else-if="successState" class="text-center py-2">
                <div class="mb-3" style="font-size: 2.5rem;">✅</div>
                <h6 class="fw-semibold mb-2">Password updated!</h6>
                <p class="text-muted small mb-1">
                  Your password has been changed successfully.
                </p>
                <p class="text-muted small mb-4">
                  Redirecting to sign in in <strong>{{ countdown }}</strong> second{{ countdown !== 1 ? 's' : '' }}...
                </p>
                <router-link to="/auth/login" class="btn btn-primary btn-sm px-4">
                  Sign In Now
                </router-link>
              </div>

              <!-- Form -->
              <template v-else-if="token">
                <!-- Error Alert -->
                <div v-if="errorMessage" class="alert alert-danger alert-dismissible py-2 small" role="alert">
                  {{ errorMessage }}
                  <button type="button" class="btn-close btn-close-sm" @click="errorMessage = ''"></button>
                </div>

                <form @submit.prevent="handleSubmit" novalidate>
                  <div class="mb-3">
                    <label for="newPassword" class="form-label">New Password</label>
                    <div class="input-group">
                      <input
                        id="newPassword"
                        v-model="form.newPassword"
                        :type="showPassword ? 'text' : 'password'"
                        class="form-control"
                        :class="{ 'is-invalid': errors.newPassword }"
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
                      <div v-if="errors.newPassword" class="invalid-feedback">{{ errors.newPassword }}</div>
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
                    Update Password
                  </button>
                </form>

                <div class="text-center mt-4">
                  <router-link to="/auth/login" class="text-decoration-none small" style="color: #635bff;">
                    ← Back to Sign In
                  </router-link>
                </div>
              </template>

            </div>
          </div>

        </div>
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
</style>
