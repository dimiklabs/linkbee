<template>
  <BaseModal v-model="isOpen" size="md" @closed="onDialogClosed">
    <template #headline>
      <span class="material-symbols-outlined dialog-headline-icon">visibility</span>
      <div class="dialog-headline-text">
        <span>Link Preview</span>
        <code class="dialog-headline-slug">{{ link?.slug }}</code>
      </div>
    </template>

    <div class="dialog-content">

      <!-- Loading state -->
      <div v-if="loading" class="state-loading">
        <div class="state-loading-spinner">
          <md-circular-progress indeterminate />
        </div>
        <div class="state-loading-text">
          <span class="state-loading-label">Fetching preview</span>
          <span class="state-loading-sub">Reaching out to the destination URL…</span>
        </div>
      </div>

      <!-- Error state -->
      <div v-else-if="fetchError" class="state-error">
        <span class="material-symbols-outlined state-error-icon">wifi_off</span>
        <div class="state-error-text">
          <span class="state-error-label">Preview unavailable</span>
          <span class="state-error-sub">{{ fetchError }}</span>
        </div>
      </div>

      <!-- Main preview -->
      <template v-else-if="preview">

        <!-- OG Card -->
        <div class="og-card">
          <!-- OG image -->
          <div v-if="preview.image_url && !imgError" class="og-image-wrap">
            <img
              :src="preview.image_url"
              :alt="preview.title || 'Preview image'"
              class="og-image"
              @error="imgError = true"
            />
          </div>
          <!-- Image placeholder when no image -->
          <div v-else class="og-image-placeholder">
            <span class="material-symbols-outlined og-image-placeholder-icon">image_not_supported</span>
          </div>

          <div class="og-body">
            <!-- Site row: favicon + name -->
            <div class="og-site-row">
              <img
                v-if="preview.favicon && !faviconError"
                :src="preview.favicon"
                alt=""
                width="16"
                height="16"
                class="og-favicon"
                @error="faviconError = true"
              />
              <span v-else class="material-symbols-outlined og-favicon-fallback">public</span>
              <span class="og-site-name">{{ preview.site_name || destinationDomain }}</span>
              <a
                :href="link?.destination_url"
                target="_blank"
                rel="noopener noreferrer"
                class="og-external-link"
                title="Open destination URL"
              >
                <span class="material-symbols-outlined">open_in_new</span>
              </a>
            </div>

            <!-- Title -->
            <h3 v-if="preview.title" class="og-title">{{ preview.title }}</h3>
            <p v-else class="og-no-meta">No page title available.</p>

            <!-- Description -->
            <p v-if="preview.description" class="og-description">{{ preview.description }}</p>
          </div>
        </div>

        <!-- Link details card -->
        <div class="link-details-card">
          <div class="link-details-header">
            <span class="material-symbols-outlined link-details-header-icon">info</span>
            <span class="link-details-header-title">Link Details</span>
          </div>

          <div class="link-details-body">
            <!-- Short URL -->
            <div class="detail-row">
              <span class="detail-label">
                <span class="material-symbols-outlined detail-label-icon">bolt</span>
                Short URL
              </span>
              <a
                v-if="link?.short_url"
                :href="link.short_url"
                target="_blank"
                rel="noopener noreferrer"
                class="detail-value detail-value--primary"
              >
                {{ link.short_url }}
                <span class="material-symbols-outlined" style="font-size:12px">open_in_new</span>
              </a>
            </div>

            <!-- Destination URL -->
            <div class="detail-row">
              <span class="detail-label">
                <span class="material-symbols-outlined detail-label-icon">link</span>
                Destination
              </span>
              <a
                :href="link?.destination_url"
                target="_blank"
                rel="noopener noreferrer"
                class="detail-value detail-value--muted"
                :title="link?.destination_url"
              >
                {{ link?.destination_url }}
              </a>
            </div>

            <!-- Domain -->
            <div class="detail-row">
              <span class="detail-label">
                <span class="material-symbols-outlined detail-label-icon">public</span>
                Domain
              </span>
              <span class="detail-value detail-value--chip">{{ destinationDomain }}</span>
            </div>
          </div>
        </div>

      </template>

      <!-- Empty state -->
      <div v-else class="state-empty">
        <span class="material-symbols-outlined state-empty-icon">search_off</span>
        <span class="state-empty-label">No preview available</span>
        <span class="state-empty-sub">The destination URL may not return Open Graph metadata.</span>
      </div>

    </div>

    <template #actions>
      <button class="btn-text" @click="hide">Close</button>
      <button class="btn-outlined" v-if="link?.short_url" :href="link.short_url" target="_blank" rel="noopener noreferrer">
        <span class="material-symbols-outlined">open_in_new</span>
        Open Short Link
      </button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import previewApi from '@/api/preview';
import type { LinkPreviewData } from '@/types/preview';
import type { LinkResponse } from '@/types/links';
import BaseModal from '@/components/BaseModal.vue';

const props = defineProps<{
  modalId: string;
  link: LinkResponse | null;
}>();

const modalEl = ref<HTMLElement | null>(null);
let bsModal: any = null;

const isOpen = ref(false);
const loading = ref(false);
const fetchError = ref('');
const preview = ref<LinkPreviewData | null>(null);
const imgError = ref(false);
const faviconError = ref(false);

const destinationDomain = computed(() => {
  if (!props.link?.destination_url) return '';
  try {
    return new URL(props.link.destination_url).hostname;
  } catch {
    return props.link.destination_url;
  }
});

async function loadPreview() {
  if (!props.link) return;
  loading.value = true;
  fetchError.value = '';
  preview.value = null;
  imgError.value = false;
  faviconError.value = false;
  try {
    const res = await previewApi.get(props.link.id);
    preview.value = res.data ?? null;
  } catch {
    fetchError.value = 'Could not fetch preview. The destination URL may be inaccessible.';
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  // Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
});

onBeforeUnmount(() => {
  bsModal?.dispose();
});

function show() {
  isOpen.value = true;
  loadPreview();
  bsModal?.show();
}

function hide() {
  isOpen.value = false;
  bsModal?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

defineExpose({ show, hide });
</script>

<style scoped lang="scss">
/* ── Headline ─────────────────────────────────────────── */
.dialog-headline-icon {
  font-size: 22px;
  color: var(--md-sys-color-primary);
}

.dialog-headline-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.dialog-headline-slug {
  font-family: monospace;
  font-size: 0.78rem;
  font-weight: 400;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Content ──────────────────────────────────────────── */
.dialog-content {
  min-width: 420px;
  max-width: 100%;
  padding: 0 2px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ── Loading state ────────────────────────────────────── */
.state-loading {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 28px 20px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 16px;
}

.state-loading-spinner {
  flex-shrink: 0;
}

.state-loading-text {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.state-loading-label {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
}

.state-loading-sub {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Error state ──────────────────────────────────────── */
.state-error {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px;
  background: color-mix(in srgb, #f4a100 10%, transparent);
  border: 1px solid color-mix(in srgb, #f4a100 35%, transparent);
  border-radius: 14px;
}

.state-error-icon {
  font-size: 28px;
  color: #f4a100;
  flex-shrink: 0;
}

.state-error-text {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.state-error-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.state-error-sub {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.4;
}

/* ── OG Card ──────────────────────────────────────────── */
.og-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 16px;
  overflow: hidden;
  background: var(--md-sys-color-surface);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.og-image-wrap {
  background: var(--md-sys-color-surface-container-low);
  line-height: 0;
}

.og-image {
  width: 100%;
  max-height: 200px;
  object-fit: cover;
  display: block;
}

.og-image-placeholder {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--md-sys-color-surface-container);
}

.og-image-placeholder-icon {
  font-size: 32px;
  color: var(--md-sys-color-outline-variant);
}

.og-body {
  padding: 14px 16px;
}

.og-site-row {
  display: flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 8px;
}

.og-favicon {
  border-radius: 2px;
  flex-shrink: 0;
}

.og-favicon-fallback {
  font-size: 16px;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
}

.og-site-name {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.og-external-link {
  display: flex;
  align-items: center;
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  flex-shrink: 0;
  opacity: 0.6;
  transition: opacity 0.15s;

  .material-symbols-outlined {
    font-size: 15px;
  }

  &:hover {
    opacity: 1;
    color: var(--md-sys-color-primary);
  }
}

.og-title {
  font-size: 0.95rem;
  font-weight: 600;
  margin: 0 0 6px;
  line-height: 1.35;
  color: var(--md-sys-color-on-surface);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.og-no-meta {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  font-style: italic;
  margin: 0 0 6px;
}

.og-description {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* ── Link details card ────────────────────────────────── */
.link-details-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
}

.link-details-header {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 14px;
  background: var(--md-sys-color-surface-container-low);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.link-details-header-icon {
  font-size: 15px;
  color: var(--md-sys-color-primary);
}

.link-details-header-title {
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.07em;
  color: var(--md-sys-color-on-surface-variant);
}

.link-details-body {
  padding: 4px 0;
}

.detail-row {
  display: flex;
  align-items: baseline;
  gap: 10px;
  padding: 9px 14px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &:last-child {
    border-bottom: none;
  }
}

.detail-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  flex-shrink: 0;
  min-width: 80px;
}

.detail-label-icon {
  font-size: 13px;
}

.detail-value {
  font-size: 0.82rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 4px;

  &--primary {
    color: var(--md-sys-color-primary);
    font-weight: 600;
    text-decoration: none;
    font-family: monospace;

    &:hover {
      text-decoration: underline;
    }
  }

  &--muted {
    color: var(--md-sys-color-on-surface-variant);
    text-decoration: none;
    font-family: monospace;
    font-size: 0.78rem;

    &:hover {
      color: var(--md-sys-color-primary);
    }
  }

  &--chip {
    display: inline-flex;
    padding: 2px 10px;
    background: var(--md-sys-color-surface-container-highest);
    border-radius: 20px;
    font-family: monospace;
    font-size: 0.78rem;
    color: var(--md-sys-color-on-surface-variant);
    flex: unset;
  }
}

/* ── Empty state ──────────────────────────────────────── */
.state-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 36px 20px;
  color: var(--md-sys-color-on-surface-variant);
  text-align: center;
}

.state-empty-icon {
  font-size: 44px;
  opacity: 0.35;
}

.state-empty-label {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
}

.state-empty-sub {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  max-width: 280px;
  line-height: 1.4;
}
</style>
