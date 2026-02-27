<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <!-- Icon circle -->
      <div class="not-found-icon-wrap">
        <span class="material-symbols-outlined not-found-icon">search_off</span>
      </div>

      <!-- 404 number -->
      <div class="not-found-number">404</div>

      <h1 class="md-headline-medium not-found-title">Page not found</h1>
      <p class="md-body-large not-found-description">
        The page you're looking for doesn't exist or has been moved.
      </p>

      <div class="not-found-actions">
        <md-filled-button @click="router.push('/')">
          <span class="material-symbols-outlined" slot="icon">home</span>
          Go back home
        </md-filled-button>
        <md-outlined-button @click="router.back()">
          <span class="material-symbols-outlined" slot="icon">arrow_back</span>
          Go back
        </md-outlined-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();
</script>

<style scoped lang="scss">
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--md-sys-color-background);
  padding: 32px 16px;
}

.not-found-content {
  max-width: 440px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0;
}

.not-found-icon-wrap {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  background: var(--md-sys-color-primary-container);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}

.not-found-icon {
  font-size: 40px;
  color: var(--md-sys-color-on-primary-container);
}

.not-found-number {
  font-size: 7rem;
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.04em;
  color: var(--md-sys-color-primary);
  margin-bottom: 16px;
  background: linear-gradient(
    135deg,
    var(--md-sys-color-primary) 0%,
    color-mix(in srgb, var(--md-sys-color-primary) 60%, var(--md-sys-color-tertiary, #14b8a6)) 100%
  );
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.not-found-title {
  color: var(--md-sys-color-on-surface);
  margin: 0 0 12px;
  font-weight: 700;
}

.not-found-description {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 32px;
  line-height: 1.6;
}

.not-found-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}
</style>
