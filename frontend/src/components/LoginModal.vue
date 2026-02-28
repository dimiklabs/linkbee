<template>
  <BaseModal v-model="isOpen" size="sm" @closed="onClosed">
    <template #headline>
      <span class="material-symbols-outlined dialog-headline-icon">lock</span>
      <span>Sign in to Linkbee</span>
    </template>

    <div class="lm-body">

      <!-- Error banner -->
      <div v-if="errorMessage" class="lm-error">
        <span class="material-symbols-outlined lm-error-icon">error</span>
        <span class="lm-error-text">{{ errorMessage }}</span>
        <button class="btn-icon lm-error-close" @click="errorMessage = ''">
          <span class="material-symbols-outlined">close</span>
        </button>
      </div>

      <!-- TOTP step -->
      <template v-if="pendingTOTPSession">
        <p class="lm-subtext">Enter the 6-digit code from your authenticator app.</p>
        <form @submit.prevent="handleTOTPVerify" novalidate>
          <md-outlined-text-field
            :value="totpCode"
            @input="totpCode = ($event.target as HTMLInputElement).value"
            label="Authentication code"
            type="text"
            inputmode="numeric"
            maxlength="8"
            autocomplete="one-time-code"
            class="field-full lm-field"
          />
          <button class="btn-filled btn-full lm-submit" type="submit" :disabled="loading">
            <md-circular-progress v-if="loading" indeterminate class="lm-spinner" />
            <span v-else>Verify &amp; Sign In</span>
          </button>
          <button class="btn-text btn-full lm-back" type="button" @click="pendingTOTPSession = ''">
            <span class="material-symbols-outlined lm-back-icon">arrow_back</span>
            Back to login
          </button>
        </form>
      </template>

      <!-- Login form -->
      <template v-else>
        <form @submit.prevent="handleLogin" novalidate>
          <md-outlined-text-field
            :value="form.email"
            @input="form.email = ($event.target as HTMLInputElement).value"
            label="Email address"
            type="email"
            autocomplete="email"
            :error="!!errors.email"
            :error-text="errors.email"
            class="field-full lm-field"
          >
            <span class="material-symbols-outlined" slot="leading-icon">mail</span>
          </md-outlined-text-field>

          <md-outlined-text-field
            :value="form.password"
            @input="form.password = ($event.target as HTMLInputElement).value"
            label="Password"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="current-password"
            :error="!!errors.password"
            :error-text="errors.password"
            class="field-full lm-field"
          >
            <span class="material-symbols-outlined" slot="leading-icon">lock</span>
            <button class="btn-icon" slot="trailing-icon" type="button" tabindex="-1"
              @click="showPassword = !showPassword">
              <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
            </button>
          </md-outlined-text-field>

          <button class="btn-filled btn-full lm-submit" type="submit" :disabled="loading">
            <md-circular-progress v-if="loading" indeterminate class="lm-spinner" />
            <span v-else>Sign In</span>
          </button>
        </form>

        <!-- OAuth -->
        <div class="lm-divider">
          <md-divider />
          <span class="lm-divider-label">Or continue with</span>
          <md-divider />
        </div>

        <div class="lm-oauth-row">
          <button class="btn-outlined lm-oauth-btn" :disabled="oauthLoading" @click="handleOAuth('google')">
            <svg class="lm-oauth-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            Google
          </button>
          <button class="btn-outlined lm-oauth-btn" :disabled="oauthLoading" @click="handleOAuth('github')">
            <svg class="lm-oauth-icon lm-oauth-icon--github" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
            </svg>
            GitHub
          </button>
        </div>

        <p class="lm-signup-text">
          Don't have an account?
          <router-link to="/signup" class="lm-signup-link" @click="isOpen = false">Sign up free</router-link>
        </p>
      </template>

    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { oauthApi } from '@/api/auth';
import BaseModal from '@/components/BaseModal.vue';

const props = defineProps<{
  /** Route to push to after a successful login. Defaults to /dashboard/links?create=1 */
  redirectTo?: string;
}>();

const emit = defineEmits<{
  (e: 'logged-in'): void;
}>();

const router = useRouter();
const authStore = useAuthStore();

const isOpen = ref(false);
const loading = ref(false);
const oauthLoading = ref(false);
const errorMessage = ref('');
const showPassword = ref(false);

const pendingTOTPSession = ref('');
const totpCode = ref('');

const form = reactive({ email: '', password: '' });
const errors = reactive({ email: '', password: '' });

function show() {
  errorMessage.value = '';
  pendingTOTPSession.value = '';
  totpCode.value = '';
  form.email = '';
  form.password = '';
  errors.email = '';
  errors.password = '';
  showPassword.value = false;
  isOpen.value = true;
}

function hide() {
  isOpen.value = false;
}

function onClosed() {
  pendingTOTPSession.value = '';
}

function validate(): boolean {
  errors.email = '';
  errors.password = '';
  if (!form.email) {
    errors.email = 'Email is required.';
  } else if (!/\S+@\S+\.\S+/.test(form.email)) {
    errors.email = 'Please enter a valid email address.';
  }
  if (!form.password) errors.password = 'Password is required.';
  return !errors.email && !errors.password;
}

async function handleLogin() {
  if (!validate()) return;
  loading.value = true;
  errorMessage.value = '';
  try {
    const result = await authStore.login({ email: form.email, password: form.password, remember_me: false });
    if (result?.requiresTOTP) {
      pendingTOTPSession.value = result.totpSession ?? '';
      totpCode.value = '';
      return;
    }
    onSuccess();
  } catch (err: unknown) {
    const anyErr = err as { response?: { data?: { message?: string; description?: string } } };
    const data = anyErr?.response?.data;
    errorMessage.value = data?.message || data?.description || 'Invalid email or password.';
  } finally {
    loading.value = false;
  }
}

async function handleTOTPVerify() {
  if (!totpCode.value.trim()) return;
  loading.value = true;
  errorMessage.value = '';
  try {
    await authStore.completeTOTPLogin(pendingTOTPSession.value, totpCode.value.trim());
    onSuccess();
  } catch (err: unknown) {
    const anyErr = err as { response?: { data?: { message?: string; description?: string } } };
    const data = anyErr?.response?.data;
    errorMessage.value = data?.message || data?.description || 'Invalid authentication code.';
  } finally {
    loading.value = false;
  }
}

function handleOAuth(provider: 'google' | 'github') {
  oauthLoading.value = true;
  const url = provider === 'google' ? oauthApi.getGoogleLoginUrl() : oauthApi.getGitHubLoginUrl();
  if (url) window.location.href = url;
}

function onSuccess() {
  isOpen.value = false;
  emit('logged-in');
  router.push(props.redirectTo ?? '/dashboard/links?create=1');
}

defineExpose({ show, hide });
</script>

<style scoped>
.lm-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 0 8px;
}

.lm-error {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--md-sys-color-error-container, #fde8e8);
  color: var(--md-sys-color-on-error-container, #410002);
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 13px;
}
.lm-error-icon { font-size: 18px; flex-shrink: 0; }
.lm-error-text { flex: 1; }
.lm-error-close { margin-left: auto; }

.lm-subtext {
  margin: 0 0 4px;
  font-size: 13px;
  color: var(--md-sys-color-on-surface-variant, #49454f);
}

.lm-field { margin-bottom: 4px; }

.lm-submit {
  width: 100%;
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 40px;
}
.lm-spinner { --md-circular-progress-size: 20px; }

.lm-back {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 4px;
}
.lm-back-icon { font-size: 16px; }

.lm-divider {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 4px 0;
}
.lm-divider md-divider { flex: 1; }
.lm-divider-label {
  font-size: 12px;
  white-space: nowrap;
  color: var(--md-sys-color-on-surface-variant, #49454f);
}

.lm-oauth-row {
  display: flex;
  gap: 8px;
}
.lm-oauth-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  min-height: 40px;
}
.lm-oauth-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}
.lm-oauth-icon--github {
  fill: var(--md-sys-color-on-surface, #1c1b1f);
}

.lm-signup-text {
  margin: 4px 0 0;
  text-align: center;
  font-size: 13px;
  color: var(--md-sys-color-on-surface-variant, #49454f);
}
.lm-signup-link {
  color: var(--md-sys-color-primary, #635bff);
  font-weight: 500;
  text-decoration: none;
}
.lm-signup-link:hover { text-decoration: underline; }

.btn-full { width: 100%; }
</style>
