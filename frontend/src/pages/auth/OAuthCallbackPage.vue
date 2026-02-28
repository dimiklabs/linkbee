<template>
  <div class="oauth-callback">
    <div class="oauth-callback__card">
      <div v-if="error" class="oauth-callback__error">
        <div class="oauth-callback__icon oauth-callback__icon--error">
          <span class="material-symbols-outlined">error</span>
        </div>
        <h2 class="oauth-callback__title">Sign-in failed</h2>
        <p class="oauth-callback__message">{{ errorMessage }}</p>
        <div class="oauth-callback__actions">
          <router-link to="/login" class="btn-filled">Back to Login</router-link>
          <router-link v-if="showSignupHint" to="/signup" class="btn-outlined">Create Account</router-link>
        </div>
      </div>

      <div v-else class="oauth-callback__loading">
        <span class="css-spinner css-spinner--lg"></span>
        <p class="oauth-callback__loading-text">Signing you in…</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const error = ref(false);
const errorMessage = ref('');
const showSignupHint = ref(false);

const ERROR_MESSAGES: Record<string, string> = {
  OAUTH_NOT_LINKED: 'This account is not linked to any existing account. Please sign up first or log in with your email and password.',
  OAUTH_DISABLED: 'OAuth sign-in is currently disabled. Please use email and password.',
  OAUTH_INVALID_STATE: 'The sign-in session expired or is invalid. Please try again.',
  OAUTH_TOKEN_EXCHANGE: 'Could not complete sign-in with the provider. Please try again.',
  OAUTH_PROVIDER_ERROR: 'The sign-in provider returned an error. Please try again.',
  OAUTH_EMAIL_NOT_FOUND: 'Could not retrieve your email from the provider. Please ensure your account has a verified email.',
  USER_INACTIVE: 'Your account is inactive. Please contact support.',
  FORBIDDEN: 'Access denied. Please contact support.',
};

onMounted(async () => {
  const accessToken = route.query.access_token as string | undefined;
  const refreshToken = route.query.refresh_token as string | undefined;
  const oauthError = route.query.oauth_error as string | undefined;
  const oauthMessage = route.query.oauth_message as string | undefined;

  if (oauthError) {
    error.value = true;
    errorMessage.value = ERROR_MESSAGES[oauthError] ?? oauthMessage ?? 'Sign-in failed. Please try again.';
    showSignupHint.value = oauthError === 'OAUTH_NOT_LINKED';
    return;
  }

  if (!accessToken || !refreshToken) {
    error.value = true;
    errorMessage.value = 'Sign-in response was incomplete. Please try again.';
    return;
  }

  authStore.setTokens(accessToken, refreshToken);
  await authStore.fetchProfile();
  router.replace('/dashboard/overview');
});
</script>

<style scoped lang="scss">
.oauth-callback {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--md-sys-color-surface);
  padding: 24px;
}

.oauth-callback__card {
  width: 100%;
  max-width: 420px;
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 20px;
  padding: 48px 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

/* ── Loading state ──────────────────────────────────────────────────────── */
.oauth-callback__loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.oauth-callback__loading-text {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.9375rem;
  margin: 0;
}

/* ── Error state ────────────────────────────────────────────────────────── */
.oauth-callback__error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  width: 100%;
}

.oauth-callback__icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;

  .material-symbols-outlined { font-size: 28px; }

  &--error {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }
}

.oauth-callback__title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  margin: 0;
}

.oauth-callback__message {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.6;
  margin: 0;
}

.oauth-callback__actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  margin-top: 8px;

  a {
    text-align: center;
    text-decoration: none;
  }
}

/* ── CSS Spinner ────────────────────────────────────────────────────────── */
.css-spinner {
  display: inline-block;
  border: 3px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;

  &--lg {
    width: 40px;
    height: 40px;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
