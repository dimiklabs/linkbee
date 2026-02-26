<template>
  <div
    class="min-vh-100 d-flex flex-column align-items-center py-5 px-3"
    :class="theme === 'dark' ? 'bg-dark text-white' : 'bg-light text-dark'"
  >
    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border" :class="theme === 'dark' ? 'text-light' : 'text-primary'" role="status"></div>
    </div>

    <!-- Not found -->
    <div v-else-if="notFound" class="text-center py-5">
      <h3 class="fw-bold mb-2">Page not found</h3>
      <p class="text-muted mb-4">This bio page doesn't exist or has been unpublished.</p>
      <a href="/" class="btn btn-primary">Go home</a>
    </div>

    <!-- Bio page -->
    <div v-else-if="page" class="w-100" style="max-width: 480px;">

      <!-- Profile header -->
      <div class="text-center mb-5">
        <!-- Avatar -->
        <div class="mb-3">
          <img
            v-if="page.avatar_url"
            :src="page.avatar_url"
            :alt="page.title"
            class="rounded-circle border"
            :class="theme === 'dark' ? 'border-secondary' : 'border-light'"
            style="width: 88px; height: 88px; object-fit: cover;"
            @error="avatarError = true"
          />
          <div
            v-else
            class="rounded-circle d-inline-flex align-items-center justify-content-center fw-bold"
            :class="theme === 'dark' ? 'bg-secondary text-white' : 'bg-primary text-white'"
            style="width: 88px; height: 88px; font-size: 2rem;"
          >
            {{ (page.title || page.username).charAt(0).toUpperCase() }}
          </div>
        </div>

        <h1 class="fw-bold mb-1" style="font-size: 1.4rem;">{{ page.title || page.username }}</h1>
        <p
          v-if="page.description"
          class="mb-0"
          :class="theme === 'dark' ? 'text-light opacity-75' : 'text-muted'"
          style="font-size: 0.92rem;"
        >
          {{ page.description }}
        </p>
      </div>

      <!-- Links -->
      <div class="d-flex flex-column gap-3">
        <a
          v-for="link in activeLinks"
          :key="link.id"
          :href="link.url"
          target="_blank"
          rel="noopener noreferrer"
          class="bio-link-btn text-decoration-none text-center px-4 py-3 rounded-3 fw-semibold"
          :class="theme === 'dark' ? 'bio-link-dark' : 'bio-link-light'"
        >
          {{ link.title }}
        </a>
      </div>

      <!-- Powered by footer -->
      <div class="text-center mt-5">
        <a
          href="/"
          :class="theme === 'dark' ? 'text-light opacity-50' : 'text-muted opacity-75'"
          style="font-size: 0.78rem; text-decoration: none;"
        >
          Powered by Shortlink
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import bioApi from '@/api/bio';
import type { BioPage } from '@/types/bio';

const route = useRoute();
const loading = ref(true);
const notFound = ref(false);
const page = ref<BioPage | null>(null);
const avatarError = ref(false);

const theme = computed(() => page.value?.theme ?? 'light');

const activeLinks = computed(() =>
  (page.value?.links ?? []).filter(l => l.is_active)
);

onMounted(async () => {
  const username = route.params.username as string;
  try {
    const res = await bioApi.getPublic(username);
    if (res.data) {
      page.value = res.data;
    } else {
      notFound.value = true;
    }
  } catch {
    notFound.value = true;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.bio-link-light {
  background-color: #ffffff;
  border: 1.5px solid #e2e2e2;
  color: #1a1a2e;
  transition: background-color 0.15s, transform 0.1s;
}
.bio-link-light:hover {
  background-color: #f4f4f6;
  transform: translateY(-1px);
  color: #1a1a2e;
}

.bio-link-dark {
  background-color: #2a2a3e;
  border: 1.5px solid #3d3d5c;
  color: #e8e8f0;
  transition: background-color 0.15s, transform 0.1s;
}
.bio-link-dark:hover {
  background-color: #35354f;
  transform: translateY(-1px);
  color: #ffffff;
}

.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}
</style>
