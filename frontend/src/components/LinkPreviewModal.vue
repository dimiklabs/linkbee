<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" style="--md-dialog-container-shape:16px">
    <div slot="headline">
      Link Preview
      <span style="font-size:0.82rem;font-weight:400;color:var(--md-sys-color-on-surface-variant);display:block;margin-top:2px">
        <code>{{ link?.slug }}</code>
      </span>
    </div>

    <div slot="content" style="min-width:420px;max-width:100%;padding:0 4px">
      <!-- Loading -->
      <div v-if="loading" class="preview-loading">
        <md-circular-progress indeterminate style="--md-circular-progress-size:36px" />
        <span class="preview-loading-text">Fetching preview…</span>
      </div>

      <!-- Error -->
      <div v-else-if="fetchError" class="preview-error-banner">
        <span class="material-symbols-outlined" style="font-size:18px;color:#F4A100">warning</span>
        <span>{{ fetchError }}</span>
      </div>

      <!-- Preview card -->
      <div v-else-if="preview" class="preview-card">
        <!-- OG image -->
        <div v-if="preview.image_url && !imgError" class="preview-image-wrap">
          <img
            :src="preview.image_url"
            :alt="preview.title || 'Preview image'"
            class="preview-image"
            @error="imgError = true"
          />
        </div>

        <div class="preview-body">
          <!-- Favicon + site name -->
          <div class="preview-site-row">
            <img
              v-if="preview.favicon && !faviconError"
              :src="preview.favicon"
              alt=""
              width="16"
              height="16"
              class="preview-favicon"
              @error="faviconError = true"
            />
            <span class="preview-site-name">
              {{ preview.site_name || destinationDomain }}
            </span>
          </div>

          <!-- Title -->
          <h6 v-if="preview.title" class="preview-title">{{ preview.title }}</h6>
          <p v-if="!preview.title && !preview.description" class="preview-empty">
            No metadata available for this URL.
          </p>

          <!-- Description -->
          <p v-if="preview.description" class="preview-description">
            {{ preview.description }}
          </p>
        </div>

        <!-- Destination URL -->
        <div class="preview-dest-url">
          <span class="material-symbols-outlined" style="font-size:14px;flex-shrink:0">open_in_new</span>
          <a
            :href="link?.destination_url"
            target="_blank"
            rel="noopener noreferrer"
            class="preview-dest-link"
          >
            {{ link?.destination_url }}
          </a>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="preview-empty-state">
        <span class="material-symbols-outlined" style="font-size:40px;opacity:0.4">link_off</span>
        <span>No preview available.</span>
      </div>
    </div>

    <div slot="actions">
      <md-text-button @click="hide">Close</md-text-button>
    </div>
  </md-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import previewApi from '@/api/preview';
import type { LinkPreviewData } from '@/types/preview';
import type { LinkResponse } from '@/types/links';

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

<style scoped>
.preview-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 32px 16px;
}

.preview-loading-text {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
}

.preview-error-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: rgba(244, 161, 0, 0.1);
  border: 1px solid rgba(244, 161, 0, 0.4);
  border-radius: 8px;
  font-size: 0.875rem;
}

.preview-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  overflow: hidden;
  background: var(--md-sys-color-surface);
}

.preview-image-wrap {
  background: var(--md-sys-color-surface-container-low);
}

.preview-image {
  width: 100%;
  max-height: 220px;
  object-fit: cover;
  display: block;
}

.preview-body {
  padding: 12px 16px;
}

.preview-site-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.preview-favicon {
  border-radius: 2px;
  flex-shrink: 0;
}

.preview-site-name {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-title {
  font-weight: 600;
  font-size: 0.95rem;
  margin: 0 0 6px;
  line-height: 1.3;
  color: var(--md-sys-color-on-surface);
}

.preview-empty {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.preview-description {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.preview-dest-url {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.78rem;
  overflow: hidden;
}

.preview-dest-link {
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.preview-dest-link:hover {
  color: var(--md-sys-color-primary);
}

.preview-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 32px 16px;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}
</style>
