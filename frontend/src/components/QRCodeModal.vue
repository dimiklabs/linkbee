<template>
  <div ref="modalEl" class="modal fade" tabindex="-1" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-semibold">QR Code</h5>
          <button type="button" class="btn-close" @click="hide" aria-label="Close"></button>
        </div>

        <div class="modal-body text-center py-4">
          <p class="text-muted small mb-3">
            Scan the code or download it as a PNG image.
          </p>

          <!-- Slug badge -->
          <div class="mb-3">
            <span class="badge rounded-pill bg-light text-dark border px-3 py-2 fw-normal">
              /{{ slug }}
            </span>
          </div>

          <!-- QR image area -->
          <div class="qr-wrapper d-flex align-items-center justify-content-center mx-auto mb-4">
            <!-- Loading spinner -->
            <div v-if="imageLoading" class="d-flex flex-column align-items-center gap-2">
              <div class="spinner-border text-primary" role="status" style="width: 2.5rem; height: 2.5rem;">
                <span class="visually-hidden">Loading QR code...</span>
              </div>
              <span class="text-muted small">Generating QR code...</span>
            </div>

            <!-- QR image -->
            <img
              v-show="!imageLoading"
              :src="qrUrl"
              :alt="'QR code for /' + slug"
              class="img-fluid rounded border"
              style="max-width: 220px; max-height: 220px;"
              @load="onImageLoad"
              @error="onImageError"
            />
          </div>

          <!-- Error state -->
          <div v-if="imageError" class="alert alert-danger d-inline-flex align-items-center gap-2 py-2 px-3" role="alert">
            <i class="bi bi-exclamation-triangle-fill"></i>
            Failed to load QR code. Please try again.
          </div>
        </div>

        <div class="modal-footer border-top-0 pt-0 justify-content-center gap-2">
          <button type="button" class="btn btn-outline-secondary" @click="hide">
            Close
          </button>
          <a
            :href="qrUrl"
            :download="'qr-' + slug + '.png'"
            class="btn btn-primary"
            :class="{ disabled: imageLoading || imageError }"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="me-2" viewBox="0 0 16 16">
              <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
              <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708z"/>
            </svg>
            Download PNG
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { Modal } from 'bootstrap';
import linksApi from '@/api/links';

interface Props {
  linkId: string;
  slug: string;
}

const props = defineProps<Props>();

const modalEl = ref<HTMLElement | null>(null);
let modalInstance: Modal | null = null;

const imageLoading = ref(true);
const imageError = ref(false);

const qrUrl = computed(() => linksApi.getQRUrl(props.linkId));

function onImageLoad() {
  imageLoading.value = false;
  imageError.value = false;
}

function onImageError() {
  imageLoading.value = false;
  imageError.value = true;
}

function resetState() {
  imageLoading.value = true;
  imageError.value = false;
}

import { onMounted } from 'vue';

onMounted(() => {
  if (modalEl.value) {
    modalInstance = new Modal(modalEl.value);

    modalEl.value.addEventListener('show.bs.modal', () => {
      resetState();
    });
  }
});

function show() {
  resetState();
  modalInstance?.show();
}

function hide() {
  modalInstance?.hide();
}

defineExpose({ show, hide });
</script>

<style scoped>
.qr-wrapper {
  min-height: 240px;
  width: 240px;
}

.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}

.btn-primary:hover:not(.disabled) {
  background-color: #5249e0;
  border-color: #5249e0;
}

.text-primary {
  color: #635bff !important;
}
</style>
