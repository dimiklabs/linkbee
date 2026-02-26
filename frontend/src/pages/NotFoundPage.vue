<template>
  <div class="not-found-page">
    <div class="not-found-content text-center">
      <div class="not-found-number">404</div>
      <h1 class="not-found-title">Page not found</h1>
      <p class="not-found-description text-muted">
        The page you're looking for doesn't exist or has been moved.
      </p>
      <button class="btn btn-primary px-4" @click="router.push('/')">
        Go back home
      </button>
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
  background: linear-gradient(135deg, #f7f9fc 0%, #eef0ff 100%);
  padding: 2rem 1rem;
}

.not-found-content {
  max-width: 420px;
}

.not-found-number {
  font-size: 8rem;
  font-weight: 800;
  color: #635bff;
  line-height: 1;
  margin-bottom: 0.5rem;
  letter-spacing: -0.04em;
  background: linear-gradient(135deg, #635bff, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.not-found-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1a1f36;
  margin-bottom: 1rem;
}

.not-found-description {
  font-size: 1rem;
  margin-bottom: 2rem;
  line-height: 1.6;
}
</style>
