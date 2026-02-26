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

              <h2 class="h5 fw-semibold mb-1">Reset your password</h2>
              <p class="text-muted small mb-4">
                Enter your email address and we'll send you a link to reset your password.
              </p>

              <!-- Success state -->
              <div v-if="successState" class="text-center py-2">
                <div class="mb-3" style="font-size: 2.5rem;">📧</div>
                <h6 class="fw-semibold mb-2">Check your inbox</h6>
                <p class="text-muted small mb-4">
                  We've sent a password reset link to <strong>{{ submittedEmail }}</strong>.
                  The link expires in 1 hour.
                </p>
                <router-link to="/auth/login" class="btn btn-outline-secondary btn-sm">
                  ← Back to Sign In
                </router-link>
              </div>

              <!-- Form -->
              <template v-else>
                <!-- Error Alert -->
                <div v-if="errorMessage" class="alert alert-danger alert-dismissible py-2 small" role="alert">
                  {{ errorMessage }}
                  <button type="button" class="btn-close btn-close-sm" @click="errorMessage = ''"></button>
                </div>

                <form @submit.prevent="handleSubmit" novalidate>
                  <div class="mb-4">
                    <label for="email" class="form-label">Email address</label>
                    <input
                      id="email"
                      v-model="email"
                      type="email"
                      class="form-control"
                      :class="{ 'is-invalid': emailError }"
                      placeholder="you@example.com"
                      autocomplete="email"
                      required
                    />
                    <div v-if="emailError" class="invalid-feedback">{{ emailError }}</div>
                  </div>

                  <button
                    type="submit"
                    class="btn btn-primary w-100"
                    :disabled="loading"
                  >
                    <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                    Send Reset Link
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
