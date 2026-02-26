<template>
  <div class="modal fade" :id="modalId" tabindex="-1" ref="modalEl">
    <div class="modal-dialog modal-md">
      <div class="modal-content border-0 shadow">
        <div class="modal-header border-0 pb-0">
          <div>
            <h5 class="modal-title fw-bold">Link Preview</h5>
            <p class="text-muted small mb-0">
              <code class="small">{{ link?.slug }}</code>
            </p>
          </div>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>

        <div class="modal-body pt-3">
          <!-- Loading -->
          <div v-if="loading" class="text-center py-4">
            <div class="spinner-border spinner-border-sm text-primary me-2"></div>
            <span class="text-muted small">Fetching preview&hellip;</span>
          </div>

          <!-- Error -->
          <div v-else-if="fetchError" class="alert alert-warning py-2 small">
            {{ fetchError }}
          </div>

          <!-- Preview card -->
          <div v-else-if="preview" class="preview-card rounded-3 overflow-hidden border">

            <!-- OG image -->
            <div v-if="preview.image_url && !imgError" class="preview-image-wrap bg-light">
              <img
                :src="preview.image_url"
                :alt="preview.title || 'Preview image'"
                class="w-100"
                style="max-height: 220px; object-fit: cover; display: block;"
                @error="imgError = true"
              />
            </div>

            <div class="p-3">
              <!-- Favicon + site name -->
              <div class="d-flex align-items-center gap-2 mb-2">
                <img
                  v-if="preview.favicon && !faviconError"
                  :src="preview.favicon"
                  alt=""
                  width="16"
                  height="16"
                  class="rounded-1 flex-shrink-0"
                  @error="faviconError = true"
                />
                <span class="text-muted small text-truncate">
                  {{ preview.site_name || destinationDomain }}
                </span>
              </div>

              <!-- Title -->
              <h6 v-if="preview.title" class="fw-semibold mb-1" style="line-height: 1.3;">
                {{ preview.title }}
              </h6>
              <p v-if="!preview.title && !preview.description" class="text-muted small mb-0">
                No metadata available for this URL.
              </p>

              <!-- Description -->
              <p
                v-if="preview.description"
                class="text-muted small mb-0"
                style="display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden;"
              >
                {{ preview.description }}
              </p>
            </div>

            <!-- Destination URL -->
            <div class="px-3 pb-3">
              <a
                :href="link?.destination_url"
                target="_blank"
                rel="noopener noreferrer"
                class="text-muted text-decoration-none small text-truncate d-block"
                style="font-size: 0.78rem;"
              >
                <i class="bi bi-box-arrow-up-right me-1"></i>
                {{ link?.destination_url }}
              </a>
            </div>
          </div>

          <!-- Empty state (no OG at all) -->
          <div v-else class="text-center text-muted py-4 small">
            No preview available.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import type { Modal } from 'bootstrap';
import previewApi from '@/api/preview';
import type { LinkPreviewData } from '@/types/preview';
import type { LinkResponse } from '@/types/links';

const props = defineProps<{
  modalId: string;
  link: LinkResponse | null;
}>();

const modalEl = ref<HTMLElement | null>(null);
let bsModal: Modal | null = null;

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
  if (modalEl.value) {
    import('bootstrap').then(({ Modal }) => {
      bsModal = new Modal(modalEl.value!);
      modalEl.value!.addEventListener('shown.bs.modal', loadPreview);
    });
  }
});

onBeforeUnmount(() => {
  bsModal?.dispose();
});

defineExpose({ show: () => bsModal?.show(), hide: () => bsModal?.hide() });
</script>

<style scoped>
.preview-card {
  background: #fff;
}
</style>
