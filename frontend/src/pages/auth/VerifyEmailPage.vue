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
            <div class="card-body p-4 text-center">

              <!-- Loading state -->
              <div v-if="verifying" class="py-4">
                <div class="spinner-border mb-3" style="width: 2.5rem; height: 2.5rem; color: #635bff;" role="status">
                  <span class="visually-hidden">Verifying...</span>
                </div>
                <h5 class="fw-semibold mb-2">Verifying your email...</h5>
                <p class="text-muted small">Please wait a moment.</p>
              </div>

              <!-- Success state -->
              <div v-else-if="verified" class="py-2">
                <div class="mb-3" style="font-size: 3rem;">✅</div>
                <h5 class="fw-semibold mb-2">Email Verified!</h5>
                <p class="text-muted small mb-4">
                  Your email has been successfully verified. You can now sign in to your Shortlink account.
                </p>
                <router-link to="/auth/login" class="btn btn-primary px-4">
                  Sign In to Shortlink
                </router-link>
              </div>

              <!-- Error state -->
              <div v-else class="py-2">
                <div class="mb-3" style="font-size: 3rem;">❌</div>
                <h5 class="fw-semibold mb-2">Verification Failed</h5>
                <p class="text-muted small mb-4">{{ errorMessage }}</p>

                <!-- Resend verification section -->
                <div class="text-start bg-light rounded p-3 mb-4">
                  <p class="small fw-medium mb-2">Resend verification email</p>

                  <div v-if="resendSuccess" class="alert alert-success py-2 small mb-2">{{ resendSuccess }}</div>
                  <div v-if="resendError" class="alert alert-danger py-2 small mb-2">{{ resendError }}</div>

                  <div class="input-group">
                    <input
                      v-model="resendEmail"
                      type="email"
                      class="form-control form-control-sm"
                      placeholder="your@email.com"
                    />
                    <button
                      class="btn btn-sm btn-primary"
                      :disabled="resending || !resendEmail"
                      @click="handleResend"
                    >
                      <span v-if="resending" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                      Resend
                    </button>
                  </div>
                </div>

                <router-link to="/auth/login" class="text-decoration-none small" style="color: #635bff;">
                  ← Back to Sign In
                </router-link>
              </div>

            </div>
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
