<template>
  <div class="unlock-wrapper">
    <div class="unlock-container">

      <!-- Logo -->
      <div class="unlock-logo-wrap">
        <div class="logo-icon">S</div>
        <span class="logo-text">Linkbee</span>
      </div>

      <!-- Card -->
      <div class="m3-card m3-card--elevated unlock-card">

        <div class="lock-header">
          <span class="material-symbols-outlined lock-icon">lock</span>
          <h1 class="md-headline-small lock-title">Enter the password</h1>
          <p class="md-body-medium lock-desc">
            The owner of this link has set a password to protect it.
          </p>
        </div>

        <!-- Error banner -->
        <div v-if="invalidPassword" class="error-banner">
          <span class="material-symbols-outlined" style="font-size:20px; flex-shrink:0;">error</span>
          <span class="md-body-medium">Incorrect password. Please try again.</span>
        </div>

        <form @submit.prevent="unlock" novalidate>
          <div class="field-wrap">
            <md-outlined-text-field
              :value="password"
              @input="password = ($event.target as HTMLInputElement).value"
              label="Password"
              :type="showPassword ? 'text' : 'password'"
              autocomplete="current-password"
              autofocus
              :error="invalidPassword && password === ''"
              style="width: 100%;"
            >
              <button class="btn-icon" slot="trailing-icon" type="button" @click="showPassword = !showPassword" tabindex="-1">
                <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
              </button>
            </md-outlined-text-field>
          </div>

          <button class="btn-filled"
            type="submit"
            :disabled="!password.trim()"
            style="width: 100%;"
          >
            <span class="material-symbols-outlined" style="font-size:18px; margin-right:6px;">lock_open</span>
            Unlock &amp; Visit Link
          </button>
        </form>

      </div>

      <p class="powered-by md-body-small">
        Powered by <a href="/" class="powered-link">Linkbee</a>
      </p>

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

<style scoped lang="scss">
.unlock-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--md-sys-color-background);
  padding: 32px 16px;
}

.unlock-container {
  width: 100%;
  max-width: 400px;
}

.unlock-logo-wrap {
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

.unlock-card {
  padding: 32px;
  border-radius: 16px;
  background: var(--md-sys-color-surface-container-low);
}

.lock-header {
  text-align: center;
  margin-bottom: 24px;
}

.lock-icon {
  font-size: 48px;
  color: var(--md-sys-color-primary);
  display: block;
  margin-bottom: 12px;
}

.lock-title {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.lock-desc {
  color: var(--md-sys-color-on-surface-variant);
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
  margin-bottom: 20px;
}

.powered-by {
  text-align: center;
  margin-top: 20px;
  color: var(--md-sys-color-on-surface-variant);
}

.powered-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}
</style>
