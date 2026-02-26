<template>
  <div ref="modalEl" class="modal fade" tabindex="-1" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg">
      <div class="modal-content">

        <!-- Header -->
        <div class="modal-header border-bottom">
          <h5 class="modal-title fw-semibold">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="currentColor" class="me-2 text-primary" viewBox="0 0 16 16">
              <path d="M2 2h2v2H2zm0 10h2v2H2zm10-10h2v2h-2zm0 10h2v2h-2z"/>
              <path fill-rule="evenodd" d="M6 0v6H0V0zm-5 1v4h4V1zM4 12H2v2h2zm-2-2h2v2H2zm4-2h2v2H6zm2 2H6v2h2zm-2 2h2v2H6zM4 4H2v2h2zM10 0v6h6V0zm-1 7v9h1V8h4v1h-3v1h1v2h-2v2h1v-1h1v1h2v-4h-1V9h1V7zm1 8v2h2v-2zm2-2v2h2v-2zm-3-3h2v2h-2z"/>
            </svg>
            QR Code
            <span class="text-muted fw-normal fs-6 ms-1">/{{ slug }}</span>
          </h5>
          <button type="button" class="btn-close" @click="hide" aria-label="Close"></button>
        </div>

        <div class="modal-body p-0">
          <div class="row g-0">

            <!-- Left: preview -->
            <div class="col-md-5 d-flex flex-column align-items-center justify-content-center p-4 border-end bg-light">
              <div
                class="qr-preview-box rounded-3 d-flex align-items-center justify-content-center mb-3 shadow-sm"
                :style="{ background: options.bg }"
              >
                <div v-if="previewLoading" class="d-flex flex-column align-items-center gap-2 text-muted">
                  <div class="spinner-border" style="width:2rem;height:2rem;"></div>
                  <span class="small">Generating…</span>
                </div>
                <img
                  v-show="!previewLoading && !previewError"
                  :src="previewUrl"
                  alt="QR Code preview"
                  class="img-fluid rounded"
                  style="max-width:200px;max-height:200px;"
                  @load="onPreviewLoad"
                  @error="onPreviewError"
                />
                <div v-if="previewError && !previewLoading" class="text-danger small text-center px-2">
                  <i class="bi bi-exclamation-triangle-fill d-block fs-4 mb-1"></i>
                  Failed to generate QR code
                </div>
              </div>

              <!-- Preset themes -->
              <div class="d-flex flex-wrap gap-2 justify-content-center">
                <button
                  v-for="preset in presets"
                  :key="preset.name"
                  class="preset-btn rounded-circle border"
                  :title="preset.name"
                  :style="{ background: `linear-gradient(135deg, ${preset.fg} 50%, ${preset.bg} 50%)`, width: '28px', height: '28px' }"
                  @click="applyPreset(preset)"
                ></button>
              </div>
              <div class="text-muted small mt-1">Color presets</div>
            </div>

            <!-- Right: controls -->
            <div class="col-md-7 p-4">
              <h6 class="text-muted text-uppercase small mb-3 fw-semibold letter-spacing">Customise</h6>

              <!-- Colors row -->
              <div class="row g-3 mb-3">
                <div class="col-6">
                  <label class="form-label small fw-semibold mb-1">
                    Foreground
                    <span class="text-muted fw-normal">(dots)</span>
                  </label>
                  <div class="input-group input-group-sm">
                    <span class="input-group-text p-1">
                      <input
                        type="color"
                        class="color-swatch"
                        v-model="options.fg"
                        @input="schedulePreviewRefresh"
                      />
                    </span>
                    <input
                      type="text"
                      class="form-control form-control-sm font-monospace"
                      v-model="options.fg"
                      maxlength="7"
                      placeholder="#000000"
                      @input="schedulePreviewRefresh"
                    />
                  </div>
                </div>
                <div class="col-6">
                  <label class="form-label small fw-semibold mb-1">
                    Background
                  </label>
                  <div class="input-group input-group-sm">
                    <span class="input-group-text p-1">
                      <input
                        type="color"
                        class="color-swatch"
                        v-model="options.bg"
                        @input="schedulePreviewRefresh"
                      />
                    </span>
                    <input
                      type="text"
                      class="form-control form-control-sm font-monospace"
                      v-model="options.bg"
                      maxlength="7"
                      placeholder="#ffffff"
                      @input="schedulePreviewRefresh"
                    />
                  </div>
                </div>
              </div>

              <!-- Size -->
              <div class="mb-3">
                <label class="form-label small fw-semibold mb-1">
                  Size — <span class="text-primary">{{ options.size }}px</span>
                </label>
                <input
                  type="range"
                  class="form-range"
                  v-model.number="options.size"
                  min="128"
                  max="512"
                  step="64"
                  @input="schedulePreviewRefresh"
                />
                <div class="d-flex justify-content-between text-muted" style="font-size:0.7rem;margin-top:-4px;">
                  <span>128px</span>
                  <span>512px</span>
                </div>
              </div>

              <!-- Error correction -->
              <div class="mb-4">
                <label class="form-label small fw-semibold mb-1">Error Correction</label>
                <div class="d-flex gap-2">
                  <button
                    v-for="ec in errorLevels"
                    :key="ec.value"
                    class="btn btn-sm flex-fill"
                    :class="options.ec === ec.value ? 'btn-primary' : 'btn-outline-secondary'"
                    @click="options.ec = ec.value; schedulePreviewRefresh()"
                  >
                    <span class="fw-bold">{{ ec.value }}</span>
                    <div class="text-opacity-75" style="font-size:0.65rem;line-height:1.1;">{{ ec.label }}</div>
                  </button>
                </div>
                <div class="text-muted mt-1" style="font-size:0.72rem;">
                  Higher correction = larger code but survives damage better.
                </div>
              </div>

              <!-- Reset -->
              <button class="btn btn-sm btn-outline-secondary w-100" @click="resetOptions">
                <i class="bi bi-arrow-counterclockwise me-1"></i>Reset to defaults
              </button>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="modal-footer border-top justify-content-between">
          <button type="button" class="btn btn-outline-secondary btn-sm" @click="hide">
            Close
          </button>
          <div class="d-flex gap-2">
            <button
              class="btn btn-sm btn-outline-primary"
              :disabled="previewLoading || previewError"
              @click="downloadPNG"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" fill="currentColor" class="me-1" viewBox="0 0 16 16">
                <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
                <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708z"/>
              </svg>
              Download PNG
            </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import linksApi from '@/api/links';

interface Props {
  linkId: string;
  slug: string;
}

const props = defineProps<Props>();

// ── Defaults ──────────────────────────────────────────────────────────────────
const DEFAULT_OPTIONS = { fg: '#000000', bg: '#ffffff', size: 256, ec: 'M' };

const options = reactive({ ...DEFAULT_OPTIONS });

// ── Color presets ─────────────────────────────────────────────────────────────
const presets = [
  { name: 'Classic',    fg: '#000000', bg: '#ffffff' },
  { name: 'Stripe',     fg: '#635bff', bg: '#ffffff' },
  { name: 'Ocean',      fg: '#0d6efd', bg: '#e8f4fd' },
  { name: 'Forest',     fg: '#198754', bg: '#f0fff4' },
  { name: 'Sunset',     fg: '#dc3545', bg: '#fff5f5' },
  { name: 'Dark',       fg: '#f8f9fa', bg: '#212529' },
  { name: 'Amber',      fg: '#fd7e14', bg: '#fff8f0' },
  { name: 'Slate',      fg: '#6c757d', bg: '#f8f9fa' },
];

// ── Error correction levels ───────────────────────────────────────────────────
const errorLevels = [
  { value: 'L', label: '~7%' },
  { value: 'M', label: '~15%' },
  { value: 'Q', label: '~25%' },
  { value: 'H', label: '~30%' },
];

// ── Modal state ───────────────────────────────────────────────────────────────
const modalEl = ref<HTMLElement | null>(null);
let modalInstance: Modal | null = null;

const previewLoading = ref(true);
const previewError = ref(false);
let debounceTimer: ReturnType<typeof setTimeout> | null = null;

// ── QR URL (regenerated on each refresh) ─────────────────────────────────────
const previewUrl = ref('');

function buildUrl() {
  return linksApi.getQRUrl(props.linkId, {
    fg: options.fg,
    bg: options.bg,
    size: options.size,
    ec: options.ec,
  });
}

function refreshPreview() {
  previewLoading.value = true;
  previewError.value = false;
  previewUrl.value = buildUrl();
}

function schedulePreviewRefresh() {
  if (debounceTimer) clearTimeout(debounceTimer);
  debounceTimer = setTimeout(refreshPreview, 400);
}

function onPreviewLoad() {
  previewLoading.value = false;
  previewError.value = false;
}

function onPreviewError() {
  previewLoading.value = false;
  previewError.value = true;
}

// ── Preset / reset ────────────────────────────────────────────────────────────
function applyPreset(preset: { fg: string; bg: string }) {
  options.fg = preset.fg;
  options.bg = preset.bg;
  refreshPreview();
}

function resetOptions() {
  Object.assign(options, DEFAULT_OPTIONS);
  refreshPreview();
}

// ── Download ──────────────────────────────────────────────────────────────────
function downloadPNG() {
  const url = buildUrl();
  const a = document.createElement('a');
  a.href = url;
  a.download = `qr-${props.slug}.png`;
  a.click();
}

// ── Bootstrap modal ───────────────────────────────────────────────────────────
onMounted(() => {
  if (modalEl.value) {
    modalInstance = new Modal(modalEl.value);
    modalEl.value.addEventListener('show.bs.modal', () => {
      Object.assign(options, DEFAULT_OPTIONS);
      refreshPreview();
    });
  }
});

function show() {
  Object.assign(options, DEFAULT_OPTIONS);
  refreshPreview();
  modalInstance?.show();
}

function hide() {
  modalInstance?.hide();
}

defineExpose({ show, hide });
</script>

<style scoped>
.qr-preview-box {
  width: 240px;
  height: 240px;
}

.color-swatch {
  width: 28px;
  height: 28px;
  padding: 0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  background: none;
}

.preset-btn {
  cursor: pointer;
  padding: 0;
  transition: transform 0.15s, box-shadow 0.15s;
}

.preset-btn:hover {
  transform: scale(1.2);
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
}

.letter-spacing {
  letter-spacing: 0.05em;
}

.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}

.btn-primary:hover:not(:disabled) {
  background-color: #5249e0;
  border-color: #5249e0;
}

.text-primary {
  color: #635bff !important;
}

.btn-outline-primary {
  color: #635bff;
  border-color: #635bff;
}

.btn-outline-primary:hover {
  background-color: #635bff;
  border-color: #635bff;
}
</style>
