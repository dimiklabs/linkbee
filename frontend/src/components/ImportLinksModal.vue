<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" class="import-dialog">
    <div slot="headline" class="dialog-headline">
      <span class="material-symbols-outlined dialog-headline-icon">upload_file</span>
      Bulk Import Links
    </div>

    <div slot="content" class="dialog-content">

      <!-- Instructions banner -->
      <div class="info-banner">
        <div class="info-banner-header">
          <span class="material-symbols-outlined info-banner-icon">info</span>
          <span class="info-banner-title">CSV Format</span>
        </div>
        <div class="info-banner-body">
          <p class="info-banner-text">
            Upload a CSV with up to <strong>500 rows</strong>.
            Required column: <code>destination_url</code>.
            Optional: <code>slug</code>, <code>title</code>,
            <code>tags</code> (semicolon-separated), <code>redirect_type</code> (301 or 302),
            <code>folder_id</code>.
          </p>
          <button class="template-dl-btn" type="button" @click="downloadTemplate">
            <span class="material-symbols-outlined template-dl-icon">download</span>
            Download CSV template
          </button>
        </div>
      </div>

      <!-- Error alert -->
      <div v-if="error" class="alert-banner">
        <span class="material-symbols-outlined alert-icon">error</span>
        <span class="alert-text">{{ error }}</span>
        <md-icon-button class="alert-close" @click="error = ''">
          <span class="material-symbols-outlined">close</span>
        </md-icon-button>
      </div>

      <!-- File picker -->
      <div v-if="!result" class="file-section">
        <div
          class="drop-zone"
          :class="{
            'drop-zone--over': isDragOver,
            'drop-zone--has-file': !!selectedFile,
          }"
          @dragover.prevent="isDragOver = true"
          @dragleave.prevent="isDragOver = false"
          @drop.prevent="onDrop"
          @click="fileInput?.click()"
        >
          <!-- Idle state -->
          <template v-if="!selectedFile">
            <span class="material-symbols-outlined drop-zone-icon">cloud_upload</span>
            <div class="drop-zone-text">
              <span class="drop-zone-primary">Drag &amp; drop your CSV file here</span>
              <span class="drop-zone-secondary">
                or <span class="drop-zone-browse">browse files</span>
              </span>
            </div>
            <div class="drop-zone-hint-chips">
              <span class="drop-zone-chip">CSV only</span>
              <span class="drop-zone-chip">Max 500 rows</span>
            </div>
          </template>

          <!-- File selected state -->
          <template v-else>
            <span class="material-symbols-outlined drop-zone-icon drop-zone-icon--ready">task</span>
            <div class="drop-zone-text">
              <span class="drop-zone-primary drop-zone-primary--ready">{{ selectedFile.name }}</span>
              <span class="drop-zone-secondary">{{ (selectedFile.size / 1024).toFixed(1) }} KB · Ready to import</span>
            </div>
            <button
              class="drop-zone-remove-btn"
              type="button"
              @click.stop="clearFile"
            >
              <span class="material-symbols-outlined">close</span>
              Remove
            </button>
          </template>
        </div>

        <input ref="fileInput" type="file" accept=".csv,text/csv" style="display:none" @change="onFileChange" />

        <!-- Import progress -->
        <div v-if="importing" class="import-progress">
          <md-circular-progress indeterminate style="--md-circular-progress-size:20px;flex-shrink:0" />
          <div class="import-progress-text">
            <span class="import-progress-label">Importing your links…</span>
            <span class="import-progress-sub">This may take a moment for large files.</span>
          </div>
        </div>
      </div>

      <!-- Results section -->
      <div v-if="result" class="results-section">

        <!-- Success/partial header -->
        <div class="results-header" :class="result.failed > 0 ? 'results-header--partial' : 'results-header--success'">
          <span class="material-symbols-outlined results-header-icon">
            {{ result.failed === 0 ? 'check_circle' : 'warning' }}
          </span>
          <div class="results-header-text">
            <span class="results-header-title">
              {{ result.failed === 0 ? 'Import complete!' : 'Import finished with errors' }}
            </span>
            <span class="results-header-sub">
              {{ result.created }} of {{ result.total }} links created successfully.
            </span>
          </div>
        </div>

        <!-- Stats row -->
        <div class="stats-row">
          <div class="stat-card">
            <span class="material-symbols-outlined stat-card-icon stat-card-icon--neutral">storage</span>
            <div class="stat-value">{{ result.total }}</div>
            <div class="stat-label">Total rows</div>
          </div>
          <div class="stat-card stat-card--success">
            <span class="material-symbols-outlined stat-card-icon stat-card-icon--success">add_link</span>
            <div class="stat-value stat-value--success">{{ result.created }}</div>
            <div class="stat-label">Created</div>
          </div>
          <div class="stat-card" :class="result.failed > 0 ? 'stat-card--error' : ''">
            <span class="material-symbols-outlined stat-card-icon" :class="result.failed > 0 ? 'stat-card-icon--error' : 'stat-card-icon--neutral'">
              {{ result.failed > 0 ? 'error' : 'check' }}
            </span>
            <div class="stat-value" :class="result.failed > 0 ? 'stat-value--error' : ''">{{ result.failed }}</div>
            <div class="stat-label">Failed</div>
          </div>
        </div>

        <!-- Error rows table -->
        <div v-if="result.errors && result.errors.length > 0" class="errors-section">
          <div class="errors-section-header">
            <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-error)">error_outline</span>
            <span class="errors-section-title">Failed rows ({{ result.errors.length }})</span>
          </div>
          <div class="errors-table-wrap">
            <table class="errors-table">
              <thead>
                <tr>
                  <th class="col-row">Row</th>
                  <th class="col-url">URL</th>
                  <th class="col-error">Reason</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="e in result.errors" :key="e.row">
                  <td class="col-row">{{ e.row }}</td>
                  <td class="col-url url-cell" :title="e.url">{{ e.url || '—' }}</td>
                  <td class="col-error error-cell">{{ e.error }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

    </div>

    <div slot="actions" class="dialog-actions">
      <md-text-button @click="hide" :disabled="importing">
        {{ result ? 'Close' : 'Cancel' }}
      </md-text-button>
      <md-filled-tonal-button v-if="result" @click="reset" class="action-btn">
        <span class="material-symbols-outlined" slot="icon">upload_file</span>
        Import Another
      </md-filled-tonal-button>
      <md-filled-button
        v-else
        :disabled="!selectedFile || importing"
        @click="handleImport"
        class="action-btn"
      >
        <md-circular-progress v-if="importing" indeterminate style="--md-circular-progress-size:18px" slot="icon" />
        <span class="material-symbols-outlined" v-else slot="icon">rocket_launch</span>
        {{ importing ? 'Importing…' : 'Import Links' }}
      </md-filled-button>
    </div>
  </md-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import linksApi from '@/api/links';
import type { ImportLinksResponse } from '@/types/links';

const emit = defineEmits<{
  imported: [result: ImportLinksResponse];
}>();

const modalEl = ref<HTMLElement | null>(null);
let modalInstance: any = null;

const isOpen = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const isDragOver = ref(false);
const importing = ref(false);
const error = ref('');
const result = ref<ImportLinksResponse | null>(null);

onMounted(() => {
  // Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
});

function show() {
  reset();
  isOpen.value = true;
  modalInstance?.show();
}

function hide() {
  isOpen.value = false;
  modalInstance?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

function reset() {
  selectedFile.value = null;
  isDragOver.value = false;
  importing.value = false;
  error.value = '';
  result.value = null;
  if (fileInput.value) fileInput.value.value = '';
}

function clearFile() {
  selectedFile.value = null;
  if (fileInput.value) fileInput.value.value = '';
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    selectedFile.value = input.files[0];
  }
}

function onDrop(e: DragEvent) {
  isDragOver.value = false;
  const files = e.dataTransfer?.files;
  if (files && files[0]) {
    selectedFile.value = files[0];
  }
}

async function handleImport() {
  if (!selectedFile.value) return;
  importing.value = true;
  error.value = '';
  try {
    const response = await linksApi.importCSV(selectedFile.value);
    if (response.data) {
      result.value = response.data;
      emit('imported', response.data);
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Import failed. Please try again.';
  } finally {
    importing.value = false;
  }
}

function downloadTemplate() {
  const header = 'destination_url,slug,title,tags,redirect_type,folder_id';
  const sample1 = 'https://example.com/my-long-url,my-slug,My Link Title,tag1;tag2,302,';
  const sample2 = 'https://another.com/page,,Another Page,,301,';
  const csv = [header, sample1, sample2].join('\n');
  const blob = new Blob([csv], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'shortlink-import-template.csv';
  a.click();
  URL.revokeObjectURL(url);
}

defineExpose({ show, hide });
</script>

<style scoped lang="scss">
.import-dialog {
  --md-dialog-container-shape: 20px;
  --md-dialog-container-max-inline-size: 600px;
}

/* ── Headline ─────────────────────────────────────────── */
.dialog-headline {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.dialog-headline-icon {
  font-size: 22px;
  color: var(--md-sys-color-primary);
}

/* ── Content ──────────────────────────────────────────── */
.dialog-content {
  min-width: 520px;
  max-width: 100%;
  padding: 0 2px;
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* ── Info banner ──────────────────────────────────────── */
.info-banner {
  border: 1px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 25%, transparent);
  background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 6%, transparent);
  border-radius: 14px;
  overflow: hidden;
  margin-bottom: 18px;
}

.info-banner-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  border-bottom: 1px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 20%, transparent);
  background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 10%, transparent);
}

.info-banner-icon {
  font-size: 16px;
  color: var(--md-sys-color-primary);
}

.info-banner-title {
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--md-sys-color-primary);
}

.info-banner-body {
  padding: 10px 14px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-banner-text {
  margin: 0;
  font-size: 0.84rem;
  line-height: 1.5;
  color: var(--md-sys-color-on-surface);

  code {
    font-family: monospace;
    font-size: 0.82rem;
    background: var(--md-sys-color-surface-container-highest);
    padding: 1px 5px;
    border-radius: 4px;
    color: var(--md-sys-color-on-surface-variant);
  }
}

.template-dl-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border: 1.5px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 40%, transparent);
  border-radius: 10px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-primary);
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
  align-self: flex-start;

  &:hover {
    background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 8%, transparent);
    border-color: var(--md-sys-color-primary);
  }
}

.template-dl-icon {
  font-size: 16px;
}

/* ── Alert banner ─────────────────────────────────────── */
.alert-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 12px;
  margin-bottom: 16px;
  background: var(--md-sys-color-error-container, #ffdad6);
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 30%, transparent);
  color: var(--md-sys-color-on-error-container, #410002);
  font-size: 0.875rem;
}

.alert-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.alert-text {
  flex: 1;
}

.alert-close {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

/* ── File section ─────────────────────────────────────── */
.file-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 4px;
}

/* ── Drop zone ────────────────────────────────────────── */
.drop-zone {
  position: relative;
  cursor: pointer;
  min-height: 140px;
  border: 2px dashed var(--md-sys-color-outline-variant);
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 28px 20px;
  transition: background 0.15s, border-color 0.15s;
  text-align: center;

  &:hover,
  &--over {
    background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 5%, transparent);
    border-color: var(--md-sys-color-primary);
    border-style: solid;
  }

  &--has-file {
    border-style: solid;
    border-color: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 50%, transparent);
    background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 4%, transparent);
  }
}

.drop-zone-icon {
  font-size: 44px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.5;

  &--ready {
    color: var(--md-sys-color-primary);
    opacity: 1;
  }
}

.drop-zone-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
}

.drop-zone-primary {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);

  &--ready {
    color: var(--md-sys-color-primary);
  }
}

.drop-zone-secondary {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.drop-zone-browse {
  color: var(--md-sys-color-primary);
  font-weight: 600;
  text-decoration: underline;
}

.drop-zone-hint-chips {
  display: flex;
  gap: 6px;
  margin-top: 2px;
}

.drop-zone-chip {
  padding: 3px 10px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.72rem;
  font-weight: 500;
}

.drop-zone-remove-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 40%, transparent);
  border-radius: 10px;
  background: transparent;
  color: var(--md-sys-color-error);
  font-size: 0.78rem;
  font-weight: 500;
  cursor: pointer;
  margin-top: 4px;
  transition: background 0.15s;

  &:hover {
    background: var(--md-sys-color-error-container, #ffdad6);
  }

  .material-symbols-outlined {
    font-size: 14px;
  }
}

/* ── Import progress ──────────────────────────────────── */
.import-progress {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 16px;
  background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 6%, transparent);
  border: 1px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 20%, transparent);
  border-radius: 12px;
}

.import-progress-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.import-progress-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-primary);
}

.import-progress-sub {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Results ──────────────────────────────────────────── */
.results-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-top: 2px;
}

.results-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 14px;

  &--success {
    background: color-mix(in srgb, #1aa563 10%, transparent);
    border: 1px solid color-mix(in srgb, #1aa563 30%, transparent);
  }

  &--partial {
    background: color-mix(in srgb, #f4a100 10%, transparent);
    border: 1px solid color-mix(in srgb, #f4a100 30%, transparent);
  }
}

.results-header-icon {
  font-size: 28px;

  .results-header--success & {
    color: #1aa563;
  }

  .results-header--partial & {
    color: #f4a100;
  }
}

.results-header-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.results-header-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.results-header-sub {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* Stats row */
.stats-row {
  display: flex;
  gap: 10px;
}

.stat-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 14px 8px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  text-align: center;
  background: var(--md-sys-color-surface-container-lowest, transparent);

  &--success {
    border-color: color-mix(in srgb, #1aa563 40%, transparent);
    background: color-mix(in srgb, #1aa563 5%, transparent);
  }

  &--error {
    border-color: color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 40%, transparent);
    background: color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 5%, transparent);
  }
}

.stat-card-icon {
  font-size: 20px;

  &--neutral {
    color: var(--md-sys-color-on-surface-variant);
  }

  &--success {
    color: #1aa563;
  }

  &--error {
    color: var(--md-sys-color-error);
  }
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  line-height: 1;

  &--success {
    color: #1aa563;
  }

  &--error {
    color: var(--md-sys-color-error);
  }
}

.stat-label {
  font-size: 0.72rem;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 500;
}

/* Errors table */
.errors-section {
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 30%, transparent);
  border-radius: 12px;
  overflow: hidden;
}

.errors-section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 6%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 20%, transparent);
}

.errors-section-title {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--md-sys-color-error);
}

.errors-table-wrap {
  max-height: 240px;
  overflow-y: auto;
}

.errors-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8rem;

  th {
    padding: 8px 12px;
    text-align: left;
    font-size: 0.68rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    position: sticky;
    top: 0;
  }

  td {
    padding: 8px 12px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    vertical-align: middle;
    color: var(--md-sys-color-on-surface);

    &:last-child {
      border-bottom-color: transparent;
    }
  }

  tr:last-child td {
    border-bottom: none;
  }
}

.col-row {
  width: 52px;
}

.col-url {
  max-width: 200px;
}

.col-error {
  color: var(--md-sys-color-error) !important;
}

.url-cell {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;
  font-size: 0.78rem;
}

.error-cell {
  color: var(--md-sys-color-error);
}

/* ── Actions ──────────────────────────────────────────── */
.dialog-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.action-btn {
  min-width: 140px;
}
</style>
