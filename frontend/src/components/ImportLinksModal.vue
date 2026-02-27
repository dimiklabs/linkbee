<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" style="--md-dialog-container-shape:16px">
    <div slot="headline">Bulk Import Links</div>

    <div slot="content" style="min-width:540px;max-width:100%;padding:0 4px">
      <!-- Instructions -->
      <div class="import-info-banner">
        <span class="material-symbols-outlined" style="font-size:18px;flex-shrink:0;margin-top:1px">info</span>
        <div class="import-info-text">
          Upload a CSV file with up to <strong>500 rows</strong>. Required column:
          <code>destination_url</code>. Optional: <code>slug</code>, <code>title</code>,
          <code>tags</code> (semicolon-separated), <code>redirect_type</code> (301 or 302),
          <code>folder_id</code>.
          <md-text-button @click="downloadTemplate" style="--md-text-button-label-text-size:0.82rem;vertical-align:baseline">
            Download template
          </md-text-button>
        </div>
      </div>

      <!-- Error alert -->
      <div v-if="error" class="import-error-banner">
        <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-error)">error</span>
        <span style="flex:1;font-size:0.875rem">{{ error }}</span>
        <md-icon-button @click="error = ''" style="width:32px;height:32px">
          <span class="material-symbols-outlined" style="font-size:18px">close</span>
        </md-icon-button>
      </div>

      <!-- File picker -->
      <div v-if="!result" class="file-section">
        <div
          class="drop-zone"
          :class="{ 'drop-zone--over': isDragOver }"
          @dragover.prevent="isDragOver = true"
          @dragleave.prevent="isDragOver = false"
          @drop.prevent="onDrop"
          @click="fileInput?.click()"
        >
          <span class="material-symbols-outlined drop-zone-icon">upload_file</span>
          <div class="drop-zone-text">
            <span v-if="selectedFile" class="drop-zone-filename">{{ selectedFile.name }}</span>
            <span v-else class="drop-zone-hint">
              Drag &amp; drop a CSV file or <span style="color:var(--md-sys-color-primary);text-decoration:underline;cursor:pointer">browse</span>
            </span>
          </div>
        </div>
        <input ref="fileInput" type="file" accept=".csv,text/csv" style="display:none" @change="onFileChange" />
        <div v-if="selectedFile" class="file-meta">
          <span>{{ (selectedFile.size / 1024).toFixed(1) }} KB</span>
          <md-text-button @click.stop="clearFile" style="--md-text-button-label-text-color:var(--md-sys-color-error);--md-text-button-label-text-size:0.82rem">
            Remove
          </md-text-button>
        </div>
      </div>

      <!-- Results -->
      <div v-if="result" class="results-section">
        <div class="results-stats">
          <div class="stat-card">
            <div class="stat-value">{{ result.total }}</div>
            <div class="stat-label">Total rows</div>
          </div>
          <div class="stat-card stat-card--success">
            <div class="stat-value stat-value--success">{{ result.created }}</div>
            <div class="stat-label">Created</div>
          </div>
          <div class="stat-card" :class="result.failed > 0 ? 'stat-card--error' : ''">
            <div class="stat-value" :class="result.failed > 0 ? 'stat-value--error' : ''">{{ result.failed }}</div>
            <div class="stat-label">Failed</div>
          </div>
        </div>

        <div v-if="result.errors && result.errors.length > 0" class="results-errors">
          <p class="results-errors-title">Failed rows:</p>
          <div class="results-errors-scroll">
            <table class="m3-table" style="font-size:0.82rem">
              <thead>
                <tr>
                  <th style="width:50px">Row</th>
                  <th>URL</th>
                  <th>Error</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="e in result.errors" :key="e.row">
                  <td>{{ e.row }}</td>
                  <td style="max-width:220px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap" :title="e.url">{{ e.url || '—' }}</td>
                  <td style="color:var(--md-sys-color-error)">{{ e.error }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <div slot="actions">
      <md-text-button @click="hide" :disabled="importing">
        {{ result ? 'Close' : 'Cancel' }}
      </md-text-button>
      <md-filled-tonal-button v-if="result" @click="reset">
        Import Another File
      </md-filled-tonal-button>
      <md-filled-button
        v-else
        :disabled="!selectedFile || importing"
        @click="handleImport"
      >
        <md-circular-progress v-if="importing" indeterminate style="--md-circular-progress-size:18px" slot="icon" />
        {{ importing ? 'Importing...' : 'Import' }}
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

<style scoped>
.import-info-banner {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  padding: 12px 14px;
  background: rgba(99, 91, 255, 0.06);
  border: 1px solid rgba(99, 91, 255, 0.2);
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface);
}

.import-info-text {
  flex: 1;
  line-height: 1.5;
}

.import-error-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  margin-bottom: 16px;
  background: var(--md-sys-color-error-container, #FFDAD6);
  color: var(--md-sys-color-on-error-container, #410002);
  border-radius: 8px;
}

.file-section {
  margin-bottom: 8px;
}

.drop-zone {
  cursor: pointer;
  min-height: 120px;
  border: 2px dashed var(--md-sys-color-outline-variant);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 24px 16px;
  transition: background 0.15s, border-color 0.15s;
}

.drop-zone:hover,
.drop-zone--over {
  background: rgba(99, 91, 255, 0.04);
  border-color: var(--md-sys-color-primary);
}

.drop-zone-icon {
  font-size: 36px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.6;
}

.drop-zone-filename {
  font-weight: 500;
  color: var(--md-sys-color-primary);
  font-size: 0.9rem;
}

.drop-zone-hint {
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
  text-align: center;
}

.file-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* Results */
.results-section {
  padding-top: 4px;
}

.results-stats {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.stat-card {
  flex: 1;
  min-width: 80px;
  text-align: center;
  padding: 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
}

.stat-card--success {
  border-color: #1AA563;
}

.stat-card--error {
  border-color: var(--md-sys-color-error);
}

.stat-value {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
}

.stat-value--success {
  color: #1AA563;
}

.stat-value--error {
  color: var(--md-sys-color-error);
}

.stat-label {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
}

.results-errors-title {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--md-sys-color-error);
  margin: 0 0 8px;
}

.results-errors-scroll {
  max-height: 260px;
  overflow-y: auto;
}
</style>
