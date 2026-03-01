<template>
  <div
    class="bio-page"
    :class="theme === 'dark' ? 'bio-page--dark' : 'bio-page--light'"
  >
    <!-- Loading -->
    <div v-if="loading" class="bio-loading">
      <md-circular-progress indeterminate />
    </div>

    <!-- Not found -->
    <div v-else-if="notFound" class="bio-not-found">
      <div class="bio-not-found__icon-wrap">
        <span class="material-symbols-outlined not-found-icon">link_off</span>
      </div>
      <h3 class="md-headline-small not-found-title">Page not found</h3>
      <p class="md-body-medium not-found-desc">This bio page doesn't exist or has been unpublished.</p>
      <a href="/">
        <button class="btn-filled">
          <span class="material-symbols-outlined">home</span>
          Go home
        </button>
      </a>
    </div>

    <!-- Bio page -->
    <div v-else-if="page" class="bio-content">

      <!-- Profile header -->
      <div class="bio-profile">
        <!-- Avatar -->
        <div class="bio-avatar-wrap">
          <img
            v-if="page.avatar_url && !avatarError"
            :src="page.avatar_url"
            :alt="page.title"
            class="bio-avatar"
            @error="avatarError = true"
          />
          <div v-else class="bio-avatar bio-avatar--fallback">
            {{ (page.title || page.username).charAt(0).toUpperCase() }}
          </div>
        </div>

        <h1 class="md-title-large bio-name">{{ page.title || page.username }}</h1>
        <p v-if="page.description" class="md-body-medium bio-description">
          {{ page.description }}
        </p>
      </div>

      <!-- Links -->
      <div class="bio-links">
        <a
          v-for="link in activeLinks"
          :key="link.id"
          :href="link.url"
          target="_blank"
          rel="noopener noreferrer"
          class="bio-link-card m3-card--outlined"
          :class="theme === 'dark' ? 'bio-link-card--dark' : 'bio-link-card--light'"
          @click="trackClick(link.id)"
        >
          <span class="md-label-large bio-link-title">{{ link.title }}</span>
          <span class="material-symbols-outlined bio-link-arrow">open_in_new</span>
        </a>
      </div>

      <!-- Powered by footer -->
      <div class="bio-footer">
        <a
          href="/"
          class="bio-powered-link md-body-small"
        >
          Powered by Linkbee
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

function trackClick(linkId: string) {
  const username = route.params.username as string;
  bioApi.recordLinkClick(username, linkId);
}

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

<style scoped lang="scss">
.bio-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 16px;

  &--light {
    background: var(--md-sys-color-background);
    color: var(--md-sys-color-on-background);
  }

  &--dark {
    background: #0f0f1a;
    color: #e8e8f0;
  }
}

.bio-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}

.bio-not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  text-align: center;
  gap: 12px;

  a {
    text-decoration: none;
    margin-top: 8px;
  }
}

.bio-not-found__icon-wrap {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}

.not-found-icon {
  font-size: 40px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.7;
}

.not-found-title {
  color: var(--md-sys-color-on-surface);
  margin: 0;
}

.not-found-desc {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.bio-content {
  width: 100%;
  max-width: 480px;
}

.bio-profile {
  text-align: center;
  margin-bottom: 36px;
}

.bio-avatar-wrap {
  margin-bottom: 16px;
}

.bio-avatar {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid var(--md-sys-color-outline-variant);

  &--fallback {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    font-size: 2rem;
    font-weight: 700;
    border: none;
  }
}

.bio-name {
  margin-bottom: 8px;
  color: inherit;
}

.bio-description {
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.6;
}

.bio-links {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bio-link-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-radius: 14px;
  text-decoration: none;
  transition: transform 0.18s ease, box-shadow 0.18s ease, background 0.15s;
  border: 1.5px solid var(--md-sys-color-outline-variant);

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  }

  &:active {
    transform: translateY(-1px);
  }

  &--light {
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface);

    &:hover {
      background: var(--md-sys-color-surface-container-low);
    }
  }

  &--dark {
    background: #1e1e30;
    border-color: #3d3d5c;
    color: #e8e8f0;

    &:hover {
      background: #28283d;
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
    }
  }
}

.bio-link-title {
  color: inherit;
}

.bio-link-arrow {
  font-size: 18px;
  opacity: 0.6;
  flex-shrink: 0;
}

.bio-footer {
  text-align: center;
  margin-top: 40px;
}

.bio-powered-link {
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  opacity: 0.6;

  &:hover {
    opacity: 1;
  }
}
</style>
