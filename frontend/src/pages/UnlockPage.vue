<template>
  <div class="unlock-wrapper">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-12" style="max-width: 420px;">

          <!-- Logo -->
          <div class="text-center mb-4">
            <div class="unlock-logo">🔗 Shortlink</div>
            <p class="text-muted small mt-1">This link is password protected</p>
          </div>

          <!-- Card -->
          <div class="card shadow-sm border-0 unlock-card">
            <div class="card-body p-4">

              <div class="text-center mb-4">
                <div class="lock-icon">🔒</div>
                <h5 class="fw-semibold mt-2 mb-1">Enter the password</h5>
                <p class="text-muted small mb-0">
                  The owner of this link has set a password to protect it.
                </p>
              </div>

              <!-- Error alert -->
              <div v-if="invalidPassword" class="alert alert-danger py-2 small mb-3" role="alert">
                Incorrect password. Please try again.
              </div>

              <form @submit.prevent="unlock" novalidate>
                <div class="mb-3">
                  <label for="pwd" class="form-label">Password</label>
                  <div class="input-group">
                    <input
                      id="pwd"
                      v-model="password"
                      :type="showPassword ? 'text' : 'password'"
                      class="form-control"
                      :class="{ 'is-invalid': invalidPassword && password === '' }"
                      placeholder="••••••••"
                      autocomplete="current-password"
                      autofocus
                      required
                    />
                    <button
                      type="button"
                      class="btn btn-outline-secondary"
                      tabindex="-1"
                      @click="showPassword = !showPassword"
                    >
                      <svg v-if="showPassword" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                        <path d="M13.359 11.238C15.06 9.72 16 8 16 8s-3-5.5-8-5.5a7 7 0 0 0-2.79.588l.77.771A6 6 0 0 1 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755q-.247.248-.517.486z"/>
                        <path d="M11.297 9.176a3.5 3.5 0 0 0-4.474-4.474l.823.823a2.5 2.5 0 0 1 2.829 2.829zm-2.943 1.299.822.822a3.5 3.5 0 0 1-4.474-4.474l.823.823a2.5 2.5 0 0 0 2.829 2.829"/>
                        <path d="M3.35 5.47q-.27.24-.518.487A13 13 0 0 0 1.172 8l.195.288c.335.48.83 1.12 1.465 1.755C4.121 11.332 5.881 12.5 8 12.5c.716 0 1.39-.133 2.02-.36l.77.772A7 7 0 0 1 8 13.5C3 13.5 0 8 0 8s.939-1.721 2.641-3.238l.708.708zm10.296 8.884-12-12 .708-.708 12 12z"/>
                      </svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                        <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8M1.173 8a13 13 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5s3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5s-3.879-1.168-5.168-2.457A13 13 0 0 1 1.172 8z"/>
                        <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5M4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0"/>
                      </svg>
                    </button>
                  </div>
                </div>

                <button
                  type="submit"
                  class="btn btn-primary w-100"
                  :disabled="!password.trim()"
                >
                  Unlock &amp; Visit Link
                </button>
              </form>

            </div>
          </div>

          <p class="text-center text-muted small mt-4">
            Powered by <a href="/" class="text-decoration-none" style="color: #635bff;">Shortlink</a>
          </p>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const slug = computed(() => route.params.slug as string);
const password = ref('');
const showPassword = ref(false);
const invalidPassword = ref(false);

onMounted(() => {
  invalidPassword.value = route.query.error === 'invalid';
});

function unlock() {
  if (!password.value.trim()) return;
  // Navigate to /:slug?pwd=<password> — the backend will verify and either
  // redirect to the destination or bounce back to /unlock/:slug?error=invalid.
  window.location.href = `/${slug.value}?pwd=${encodeURIComponent(password.value)}`;
}
</script>

<style scoped>
.unlock-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  background: #f7f9fc;
  padding: 2rem 1rem;
}

.unlock-logo {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1f36;
}

.lock-icon {
  font-size: 2.5rem;
  line-height: 1;
}

.unlock-card {
  border-radius: 16px;
}

.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}

.btn-primary:hover:not(:disabled) {
  background-color: #5249e0;
  border-color: #5249e0;
}

.btn-primary:disabled {
  opacity: 0.6;
}

.form-control:focus {
  border-color: #635bff;
  box-shadow: 0 0 0 0.2rem rgba(99, 91, 255, 0.2);
}
</style>
