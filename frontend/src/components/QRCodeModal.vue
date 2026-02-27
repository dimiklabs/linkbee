<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" style="--md-dialog-container-shape:16px">
    <div slot="headline">
      QR Code
      <span style="font-size:0.85rem;font-weight:400;color:var(--md-sys-color-on-surface-variant);margin-left:6px">/{{ slug }}</span>
    </div>

    <div slot="content" style="min-width:560px;max-width:100%;padding:0 4px">
      <div class="qr-layout">

        <!-- Left: preview -->
        <div class="qr-preview-panel">
          <div class="qr-preview-box" :style="{ background: options.bg }">
            <div v-if="previewLoading" class="qr-loading">
              <md-circular-progress indeterminate style="--md-circular-progress-size:40px" />
              <span class="qr-loading-text">Generating…</span>
            </div>
            <img
              v-show="!previewLoading && !previewError"
              :src="previewUrl"
              alt="QR Code preview"
              class="qr-img"
              @load="onPreviewLoad"
              @error="onPreviewError"
            />
            <div v-if="previewError && !previewLoading" class="qr-error">
              <span class="material-symbols-outlined" style="font-size:32px;color:var(--md-sys-color-error)">error</span>
              <span style="font-size:0.82rem;text-align:center">Failed to generate QR code</span>
            </div>
          </div>

          <!-- Preset themes -->
          <div class="qr-presets">
            <button
              v-for="preset in presets"
              :key="preset.name"
              class="preset-btn"
              :title="preset.name"
              :style="{ background: `linear-gradient(135deg, ${preset.fg} 50%, ${preset.bg} 50%)` }"
              @click="applyPreset(preset)"
            ></button>
          </div>
          <div class="qr-presets-label">Color presets</div>
        </div>

        <!-- Right: controls -->
        <div class="qr-controls-panel">
          <div class="qr-section-label">Customise</div>

          <!-- Colors row -->
          <div class="color-row">
            <div class="color-field">
              <div class="color-label">Foreground <span style="color:var(--md-sys-color-on-surface-variant)">(dots)</span></div>
              <div class="color-input-group">
                <input type="color" class="color-swatch" v-model="options.fg" @input="schedulePreviewRefresh" />
                <md-outlined-text-field
                  :value="options.fg"
                  @input="options.fg = ($event.target as HTMLInputElement).value; schedulePreviewRefresh()"
                  label="Hex"
                  maxlength="7"
                  style="flex:1;font-family:monospace"
                />
              </div>
            </div>
            <div class="color-field">
              <div class="color-label">Background</div>
              <div class="color-input-group">
                <input type="color" class="color-swatch" v-model="options.bg" @input="schedulePreviewRefresh" />
                <md-outlined-text-field
                  :value="options.bg"
                  @input="options.bg = ($event.target as HTMLInputElement).value; schedulePreviewRefresh()"
                  label="Hex"
                  maxlength="7"
                  style="flex:1;font-family:monospace"
                />
              </div>
            </div>
          </div>

          <!-- Size -->
          <div class="size-section">
            <div class="size-label">
              Size — <span style="color:var(--md-sys-color-primary)">{{ options.size }}px</span>
            </div>
            <input
              type="range"
              class="size-range"
              v-model.number="options.size"
              min="128"
              max="512"
              step="64"
              @input="schedulePreviewRefresh"
            />
            <div class="size-range-labels">
              <span>128px</span>
              <span>512px</span>
            </div>
          </div>

          <!-- Error correction -->
          <div class="ec-section">
            <div class="size-label">Error Correction</div>
            <div class="ec-buttons">
              <button
                v-for="ec in errorLevels"
                :key="ec.value"
                class="ec-btn"
                :class="{ 'ec-btn--active': options.ec === ec.value }"
                @click="options.ec = ec.value; schedulePreviewRefresh()"
              >
                <span class="ec-value">{{ ec.value }}</span>
                <span class="ec-label">{{ ec.label }}</span>
              </button>
            </div>
            <div class="ec-hint">
              Higher correction = larger code but survives damage better.
            </div>
          </div>

          <!-- Reset -->
          <md-outlined-button @click="resetOptions" style="width:100%">
            <span class="material-symbols-outlined" slot="icon">restart_alt</span>
            Reset to defaults
          </md-outlined-button>
        </div>
      </div>
    </div>

    <div slot="actions">
      <md-text-button @click="hide">Close</md-text-button>
      <md-outlined-button :disabled="previewLoading || previewError" @click="downloadSVG">
        <span class="material-symbols-outlined" slot="icon">download</span>
        SVG
      </md-outlined-button>
      <md-filled-button :disabled="previewLoading || previewError" @click="downloadPNG">
        <span class="material-symbols-outlined" slot="icon">download</span>
        PNG
      </md-filled-button>
    </div>
  </md-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
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
let modalInstance: any = null;

const isOpen = ref(false);
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

function downloadSVG() {
  const url = linksApi.getQRUrl(props.linkId, {
    fg: options.fg,
    bg: options.bg,
    size: options.size,
    ec: options.ec,
    format: 'svg',
  });
  const a = document.createElement('a');
  a.href = url;
  a.download = `qr-${props.slug}.svg`;
  a.click();
}

// Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
onMounted(() => {
  // noop
});

function show() {
  Object.assign(options, DEFAULT_OPTIONS);
  isOpen.value = true;
  refreshPreview();
  modalInstance?.show();
}

function hide() {
  isOpen.value = false;
  modalInstance?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

defineExpose({ show, hide });
</script>

<style scoped>
.qr-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 24px;
}

/* Preview panel */
.qr-preview-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.qr-preview-box {
  width: 220px;
  height: 220px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12);
  overflow: hidden;
}

.qr-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--md-sys-color-on-surface-variant);
}

.qr-loading-text {
  font-size: 0.82rem;
}

.qr-img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 8px;
}

.qr-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px;
  color: var(--md-sys-color-error);
}

.qr-presets {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
}

.preset-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid var(--md-sys-color-outline-variant);
  cursor: pointer;
  padding: 0;
  transition: transform 0.15s, box-shadow 0.15s;
}

.preset-btn:hover {
  transform: scale(1.2);
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
}

.qr-presets-label {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* Controls panel */
.qr-controls-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.qr-section-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--md-sys-color-on-surface-variant);
}

.color-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.color-label {
  font-size: 0.82rem;
  font-weight: 500;
  margin-bottom: 6px;
  color: var(--md-sys-color-on-surface);
}

.color-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-swatch {
  width: 32px;
  height: 32px;
  padding: 0;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 6px;
  cursor: pointer;
  background: none;
  flex-shrink: 0;
}

/* Size slider */
.size-label {
  font-size: 0.82rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 6px;
}

.size-range {
  width: 100%;
  accent-color: var(--md-sys-color-primary);
}

.size-range-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.7rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: -2px;
}

/* Error correction */
.ec-buttons {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
}

.ec-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 4px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}

.ec-btn:hover {
  background: var(--md-sys-color-surface-container-low);
}

.ec-btn--active {
  background: var(--md-sys-color-primary-container, rgba(99,91,255,0.12));
  border-color: var(--md-sys-color-primary);
}

.ec-value {
  font-weight: 700;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.ec-label {
  font-size: 0.65rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.2;
}

.ec-hint {
  font-size: 0.72rem;
  color: var(--md-sys-color-on-surface-variant);
}
</style>
