<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <div class="not-found-number md-display-large">404</div>
      <h1 class="md-headline-medium not-found-title">Page not found</h1>
      <p class="md-body-large not-found-description">
        The page you're looking for doesn't exist or has been moved.
      </p>
      <md-filled-button @click="router.push('/')">
        <span class="material-symbols-outlined" style="font-size:18px; margin-right:6px;">home</span>
        Go back home
      </md-filled-button>
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
  max-width: 420px;
  text-align: center;
}

.not-found-number {
  color: var(--md-sys-color-primary);
  line-height: 1;
  margin-bottom: 16px;
  letter-spacing: -0.04em;
}

.not-found-title {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 16px;
}

.not-found-description {
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 32px;
  line-height: 1.6;
}
</style>
