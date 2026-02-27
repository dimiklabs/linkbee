<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" class="qr-dialog">
    <div slot="headline" class="dialog-headline">
      <span class="material-symbols-outlined dialog-headline-icon">qr_code_2</span>
      <div class="dialog-headline-text">
        <span>QR Code</span>
        <span class="dialog-headline-sub">/{{ slug }}</span>
      </div>
    </div>

    <div slot="content" class="dialog-content">
      <div class="qr-layout">

        <!-- Left: QR preview -->
        <div class="qr-preview-panel">
          <!-- Framed QR box -->
          <div class="qr-frame-outer">
            <div class="qr-frame-inner" :style="{ background: options.bg }">
              <div v-if="previewLoading" class="qr-loading">
                <md-circular-progress indeterminate style="--md-circular-progress-size:44px" />
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
                <span class="material-symbols-outlined" style="font-size:36px">error_outline</span>
                <span>Failed to generate</span>
              </div>
            </div>

            <!-- Corner accents -->
            <span class="qr-corner qr-corner--tl"></span>
            <span class="qr-corner qr-corner--tr"></span>
            <span class="qr-corner qr-corner--bl"></span>
            <span class="qr-corner qr-corner--br"></span>
          </div>

          <!-- Slug label below QR -->
          <div class="qr-slug-label">
            <span class="material-symbols-outlined" style="font-size:14px">link</span>
            <span class="qr-slug-text">sl.ink/<strong>{{ slug }}</strong></span>
          </div>

          <!-- Color presets -->
          <div class="preset-section">
            <div class="preset-title">Color presets</div>
            <div class="preset-grid">
              <button
                v-for="preset in presets"
                :key="preset.name"
                class="preset-btn"
                :title="preset.name"
                :style="{ '--preset-fg': preset.fg, '--preset-bg': preset.bg }"
                @click="applyPreset(preset)"
              >
                <span class="preset-btn-inner" :style="{ background: `linear-gradient(135deg, ${preset.fg} 50%, ${preset.bg} 50%)` }"></span>
                <span class="preset-btn-label">{{ preset.name }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Right: controls -->
        <div class="qr-controls-panel">

          <!-- Colors -->
          <div class="ctrl-section">
            <div class="ctrl-section-header">
              <span class="material-symbols-outlined ctrl-section-icon">palette</span>
              <span class="ctrl-section-title">Colors</span>
            </div>
            <div class="color-row">
              <div class="color-field">
                <div class="color-field-label">
                  Foreground <span class="color-field-sub">(dots)</span>
                </div>
                <div class="color-input-group">
                  <label class="color-swatch-label">
                    <input type="color" class="color-swatch" v-model="options.fg" @input="schedulePreviewRefresh" />
                    <span class="color-swatch-preview" :style="{ background: options.fg }"></span>
                  </label>
                  <md-outlined-text-field
                    :value="options.fg"
                    @input="options.fg = ($event.target as HTMLInputElement).value; schedulePreviewRefresh()"
                    label="Hex"
                    maxlength="7"
                    class="color-hex-field"
                  />
                </div>
              </div>
              <div class="color-field">
                <div class="color-field-label">Background</div>
                <div class="color-input-group">
                  <label class="color-swatch-label">
                    <input type="color" class="color-swatch" v-model="options.bg" @input="schedulePreviewRefresh" />
                    <span class="color-swatch-preview" :style="{ background: options.bg }"></span>
                  </label>
                  <md-outlined-text-field
                    :value="options.bg"
                    @input="options.bg = ($event.target as HTMLInputElement).value; schedulePreviewRefresh()"
                    label="Hex"
                    maxlength="7"
                    class="color-hex-field"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Size -->
          <div class="ctrl-section">
            <div class="ctrl-section-header">
              <span class="material-symbols-outlined ctrl-section-icon">photo_size_select_large</span>
              <span class="ctrl-section-title">Size</span>
              <span class="ctrl-section-value">{{ options.size }}px</span>
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
              <span>256px</span>
              <span>384px</span>
              <span>512px</span>
            </div>
          </div>

          <!-- Error correction -->
          <div class="ctrl-section">
            <div class="ctrl-section-header">
              <span class="material-symbols-outlined ctrl-section-icon">shield</span>
              <span class="ctrl-section-title">Error Correction</span>
            </div>
            <div class="ec-buttons">
              <button
                v-for="ec in errorLevels"
                :key="ec.value"
                class="ec-btn"
                :class="{ 'ec-btn--active': options.ec === ec.value }"
                @click="options.ec = ec.value; schedulePreviewRefresh()"
              >
                <span class="ec-value">{{ ec.value }}</span>
                <span class="ec-pct">{{ ec.label }}</span>
              </button>
            </div>
            <p class="ctrl-hint">Higher correction allows the QR to remain readable even if partially damaged.</p>
          </div>

          <!-- Reset -->
          <md-outlined-button @click="resetOptions" class="reset-btn">
            <span class="material-symbols-outlined" slot="icon">restart_alt</span>
            Reset to defaults
          </md-outlined-button>

        </div>
      </div>
    </div>

    <div slot="actions" class="dialog-actions">
      <md-text-button @click="hide">Close</md-text-button>
      <div class="download-btns">
        <md-outlined-button :disabled="previewLoading || previewError" @click="downloadSVG" class="dl-btn">
          <span class="material-symbols-outlined" slot="icon">download</span>
          SVG
        </md-outlined-button>
        <md-filled-tonal-button :disabled="previewLoading || previewError" @click="downloadPNG" class="dl-btn">
          <span class="material-symbols-outlined" slot="icon">image</span>
          PNG
        </md-filled-tonal-button>
      </div>
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

<style scoped lang="scss">
.qr-dialog {
  --md-dialog-container-shape: 20px;
  --md-dialog-container-max-inline-size: 680px;
}

/* ── Headline ─────────────────────────────────────────── */
.dialog-headline {
  display: flex;
  align-items: center;
  gap: 10px;
}

.dialog-headline-icon {
  font-size: 24px;
  color: var(--md-sys-color-primary);
}

.dialog-headline-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.dialog-headline-sub {
  font-size: 0.8rem;
  font-weight: 400;
  color: var(--md-sys-color-on-surface-variant);
  font-family: monospace;
}

/* ── Content ──────────────────────────────────────────── */
.dialog-content {
  min-width: 580px;
  max-width: 100%;
  padding: 0 2px;
}

/* ── Layout ───────────────────────────────────────────── */
.qr-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 28px;
  align-items: start;
}

/* ── Preview panel ────────────────────────────────────── */
.qr-preview-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

/* QR frame with corner accents */
.qr-frame-outer {
  position: relative;
  padding: 6px;
}

.qr-frame-inner {
  width: 224px;
  height: 224px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 4px 16px rgba(0, 0, 0, 0.14),
    0 1px 4px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  position: relative;
}

/* Corner accent marks */
.qr-corner {
  position: absolute;
  width: 16px;
  height: 16px;
  border-color: var(--md-sys-color-primary);
  border-style: solid;
  border-width: 0;

  &--tl {
    top: 0;
    left: 0;
    border-top-width: 2.5px;
    border-left-width: 2.5px;
    border-top-left-radius: 4px;
  }

  &--tr {
    top: 0;
    right: 0;
    border-top-width: 2.5px;
    border-right-width: 2.5px;
    border-top-right-radius: 4px;
  }

  &--bl {
    bottom: 0;
    left: 0;
    border-bottom-width: 2.5px;
    border-left-width: 2.5px;
    border-bottom-left-radius: 4px;
  }

  &--br {
    bottom: 0;
    right: 0;
    border-bottom-width: 2.5px;
    border-right-width: 2.5px;
    border-bottom-right-radius: 4px;
  }
}

.qr-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: var(--md-sys-color-on-surface-variant);
}

.qr-loading-text {
  font-size: 0.8rem;
}

.qr-img {
  width: 192px;
  height: 192px;
  object-fit: contain;
  border-radius: 6px;
}

.qr-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px;
  color: var(--md-sys-color-error);
  font-size: 0.8rem;
  text-align: center;
}

/* Slug label */
.qr-slug-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
}

.qr-slug-text {
  font-family: monospace;

  strong {
    color: var(--md-sys-color-primary);
  }
}

/* Preset section */
.preset-section {
  width: 100%;
}

.preset-title {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.07em;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 8px;
  text-align: center;
}

.preset-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
}

.preset-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 4px 2px;
  border: 1px solid transparent;
  border-radius: 10px;
  background: transparent;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s, transform 0.1s;

  &:hover {
    border-color: var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container-low);
    transform: scale(1.05);
  }
}

.preset-btn-inner {
  display: block;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1.5px solid var(--md-sys-color-outline-variant);
}

.preset-btn-label {
  font-size: 0.6rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1;
}

/* ── Controls panel ───────────────────────────────────── */
.qr-controls-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ctrl-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ctrl-section-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 2px;
}

.ctrl-section-icon {
  font-size: 16px;
  color: var(--md-sys-color-primary);
}

.ctrl-section-title {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  flex: 1;
}

.ctrl-section-value {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--md-sys-color-primary);
  font-family: monospace;
}

/* Colors */
.color-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.color-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.color-field-label {
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
}

.color-field-sub {
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 400;
}

.color-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-swatch-label {
  position: relative;
  cursor: pointer;
}

.color-swatch {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}

.color-swatch-preview {
  display: block;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 2px solid var(--md-sys-color-outline-variant);
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.08);
  transition: border-color 0.15s;
  cursor: pointer;

  &:hover {
    border-color: var(--md-sys-color-primary);
  }
}

.color-hex-field {
  flex: 1;
  font-family: monospace;
}

/* Size slider */
.size-range {
  width: 100%;
  accent-color: var(--md-sys-color-primary);
  height: 4px;
}

.size-range-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.68rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: -2px;
}

/* Error correction */
.ec-buttons {
  display: flex;
  gap: 6px;
}

.ec-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 4px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  background: transparent;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;

  &:hover {
    background: var(--md-sys-color-surface-container-low);
    border-color: var(--md-sys-color-outline);
  }

  &--active {
    background: var(--md-sys-color-primary-container, rgba(99, 91, 255, 0.12));
    border-color: var(--md-sys-color-primary);
  }
}

.ec-value {
  font-weight: 700;
  font-size: 1rem;
  color: var(--md-sys-color-on-surface);
  line-height: 1;
}

.ec-pct {
  font-size: 0.63rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.2;
  margin-top: 2px;
}

.ctrl-hint {
  font-size: 0.72rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  line-height: 1.4;
}

.reset-btn {
  width: 100%;
}

/* ── Actions ──────────────────────────────────────────── */
.dialog-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.download-btns {
  display: flex;
  gap: 8px;
  align-items: center;
}

.dl-btn {
  min-width: 100px;
}
</style>
