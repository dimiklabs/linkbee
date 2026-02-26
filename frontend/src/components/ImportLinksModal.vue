<template>
  <div ref="modalEl" class="modal fade" tabindex="-1" aria-hidden="true">
    <div class="modal-dialog modal-lg modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title fw-semibold">Bulk Import Links</h5>
          <button type="button" class="btn-close" @click="hide" aria-label="Close"></button>
        </div>

        <div class="modal-body">
          <!-- Instructions -->
          <div class="alert alert-info d-flex gap-2 align-items-start py-2 px-3 mb-4" style="font-size: 0.85rem;">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="flex-shrink-0 mt-1" viewBox="0 0 16 16">
              <path d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16m.93-9.412-1 4.705c-.07.34.029.533.304.533.194 0 .487-.07.686-.246l-.088.416c-.287.346-.92.598-1.465.598-.703 0-1.002-.422-.808-1.319l.738-3.468c.064-.293.006-.399-.287-.47l-.451-.081.082-.381 2.29-.287zM8 5.5a1 1 0 1 1 0-2 1 1 0 0 1 0 2"/>
            </svg>
            <div>
              Upload a CSV file with up to <strong>500 rows</strong>. Required column:
              <code>destination_url</code>. Optional: <code>slug</code>, <code>title</code>,
              <code>tags</code> (semicolon-separated), <code>redirect_type</code> (301 or 302),
              <code>folder_id</code>.
              <button class="btn btn-link btn-sm p-0 ms-1 text-info" style="font-size: 0.82rem; vertical-align: baseline;" @click="downloadTemplate">
                Download template
              </button>
            </div>
          </div>

          <!-- Error alert -->
          <div v-if="error" class="alert alert-danger alert-dismissible" role="alert">
            {{ error }}
            <button type="button" class="btn-close" @click="error = ''" aria-label="Close"></button>
          </div>

          <!-- File picker -->
          <div v-if="!result" class="mb-3">
            <label class="form-label fw-medium">CSV File <span class="text-danger">*</span></label>
            <div
              class="drop-zone border rounded d-flex flex-column align-items-center justify-content-center gap-2 py-4 px-3"
              :class="{ 'drop-zone-over': isDragOver }"
              @dragover.prevent="isDragOver = true"
              @dragleave.prevent="isDragOver = false"
              @drop.prevent="onDrop"
              @click="fileInput?.click()"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor" class="text-muted" viewBox="0 0 16 16">
                <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
                <path d="M7.646 1.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 2.707V11.5a.5.5 0 0 1-1 0V2.707L5.354 4.854a.5.5 0 1 1-.708-.708z"/>
              </svg>
              <div class="text-center">
                <span v-if="selectedFile" class="fw-medium text-primary">{{ selectedFile.name }}</span>
                <span v-else class="text-muted" style="font-size: 0.9rem;">
                  Drag &amp; drop a CSV file or <span class="text-primary text-decoration-underline" style="cursor:pointer;">browse</span>
                </span>
              </div>
            </div>
            <input ref="fileInput" type="file" accept=".csv,text/csv" class="d-none" @change="onFileChange" />
            <div v-if="selectedFile" class="form-text">
              {{ (selectedFile.size / 1024).toFixed(1) }} KB &mdash;
              <button class="btn btn-link btn-sm p-0 text-danger" style="font-size: 0.82rem;" @click.stop="clearFile">Remove</button>
            </div>
          </div>

          <!-- Results -->
          <div v-if="result">
            <div class="d-flex gap-3 mb-3 flex-wrap">
              <div class="stat-card border rounded p-3 text-center flex-fill">
                <div class="fs-4 fw-bold">{{ result.total }}</div>
                <div class="text-muted small">Total rows</div>
              </div>
              <div class="stat-card border rounded p-3 text-center flex-fill" :class="result.created > 0 ? 'border-success' : ''">
                <div class="fs-4 fw-bold text-success">{{ result.created }}</div>
                <div class="text-muted small">Created</div>
              </div>
              <div class="stat-card border rounded p-3 text-center flex-fill" :class="result.failed > 0 ? 'border-danger' : ''">
                <div class="fs-4 fw-bold" :class="result.failed > 0 ? 'text-danger' : 'text-muted'">{{ result.failed }}</div>
                <div class="text-muted small">Failed</div>
              </div>
            </div>

            <div v-if="result.errors && result.errors.length > 0">
              <p class="fw-medium small text-danger mb-2">Failed rows:</p>
              <div class="table-responsive" style="max-height: 260px; overflow-y: auto;">
                <table class="table table-sm table-bordered mb-0" style="font-size: 0.82rem;">
                  <thead class="table-light sticky-top">
                    <tr>
                      <th style="width: 50px;">Row</th>
                      <th>URL</th>
                      <th>Error</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="e in result.errors" :key="e.row">
                      <td>{{ e.row }}</td>
                      <td class="text-truncate" style="max-width: 220px;" :title="e.url">{{ e.url || '—' }}</td>
                      <td class="text-danger">{{ e.error }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary" @click="hide" :disabled="importing">
            {{ result ? 'Close' : 'Cancel' }}
          </button>
          <button
            v-if="!result"
            type="button"
            class="btn btn-primary"
            :disabled="!selectedFile || importing"
            @click="handleImport"
          >
            <span v-if="importing" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
            {{ importing ? 'Importing...' : 'Import' }}
          </button>
          <button v-else type="button" class="btn btn-outline-primary" @click="reset">
            Import Another File
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import linksApi from '@/api/links';
import type { ImportLinksResponse } from '@/types/links';

const emit = defineEmits<{
  imported: [result: ImportLinksResponse];
}>();

const modalEl = ref<HTMLElement | null>(null);
let modalInstance: Modal | null = null;

const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const isDragOver = ref(false);
const importing = ref(false);
const error = ref('');
const result = ref<ImportLinksResponse | null>(null);

onMounted(() => {
  if (modalEl.value) {
    modalInstance = new Modal(modalEl.value, { backdrop: 'static' });
    modalEl.value.addEventListener('hidden.bs.modal', reset);
  }
});

function show() {
  reset();
  modalInstance?.show();
}

function hide() {
  modalInstance?.hide();
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
.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}

.btn-primary:hover:not(:disabled) {
  background-color: #5249e0;
  border-color: #5249e0;
}

.drop-zone {
  cursor: pointer;
  min-height: 120px;
  transition: background 0.15s, border-color 0.15s;
  border-style: dashed !important;
  border-color: #dee2e6;
}

.drop-zone:hover,
.drop-zone-over {
  background-color: rgba(99, 91, 255, 0.04);
  border-color: #635bff !important;
}

.stat-card {
  min-width: 90px;
}
</style>
