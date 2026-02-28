<template>
  <div class="audit-page">

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">Audit Logs</h1>
        <p class="page-subtitle">A record of all security-relevant actions in your account.</p>
      </div>
      <div class="page-header__actions">
        <button class="btn-outlined" :disabled="exporting" @click="exportLogs">
          <div v-if="exporting" class="css-spinner css-spinner--sm"></div>
          <span v-else class="material-symbols-outlined btn-icon-left">download</span>
          Export CSV
        </button>
      </div>
    </div>

    <!-- Filters Card -->
    <div class="an-card filters-card">
      <div class="an-card-header">
        <div class="an-card-icon an-card-icon--primary">
          <span class="material-symbols-outlined">filter_list</span>
        </div>
        <span class="an-card-title">Filter Events</span>
      </div>
      <div class="filters-row">
        <AppSelect v-model="filters.action" label="Action">
          <option value="">All actions</option>
          <option v-for="(label, key) in ACTION_LABELS" :key="key" :value="key">{{ label }}</option>
        </AppSelect>

        <AppSelect v-model="filters.resource_type" label="Resource">
          <option value="">All resources</option>
          <option v-for="(label, key) in RESOURCE_LABELS" :key="key" :value="key">{{ label }}</option>
        </AppSelect>

        <label class="form-field">
          <span class="form-field__label">From</span>
          <input
            type="date"
            class="form-input"
            :value="filters.from"
            @input="filters.from=($event.target as HTMLInputElement).value"
          />
        </label>

        <label class="form-field">
          <span class="form-field__label">To</span>
          <input
            type="date"
            class="form-input"
            :value="filters.to"
            @input="filters.to=($event.target as HTMLInputElement).value"
          />
        </label>

        <div class="filter-actions">
          <button class="btn-filled" @click="applyFilters">
            <span class="material-symbols-outlined btn-icon-left">search</span>
            Apply
          </button>
          <button class="btn-outlined" @click="resetFilters">Reset</button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-center">
      <div class="css-spinner"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="logs.length === 0" class="an-card empty-state">
      <div class="empty-icon">
        <span class="material-symbols-outlined">event_note</span>
      </div>
      <div class="empty-title">No audit logs found</div>
      <p class="empty-desc">Actions you take will be recorded here for security and compliance.</p>
    </div>

    <!-- Table -->
    <div v-else>
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">security</span>
            </div>
            <span class="an-card-title">Audit Events</span>
          </div>
          <span class="m3-badge m3-badge--neutral">{{ total }} total</span>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Time</th>
                <th>Action</th>
                <th>Resource</th>
                <th>Name / ID</th>
                <th>IP Address</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="log in logs" :key="log.id">
                <td class="cell-muted cell-nowrap cell-sm">{{ formatDate(log.created_at) }}</td>
                <td>
                  <span :class="['m3-badge', actionBadge(log.action)]">
                    {{ actionLabel(log.action) }}
                  </span>
                </td>
                <td class="cell-muted cell-sm cell-capitalize">{{ resourceLabel(log.resource_type) }}</td>
                <td class="cell-sm">
                  <span v-if="log.resource_name" class="cell-strong">{{ log.resource_name }}</span>
                  <span v-else-if="log.resource_id" class="cell-truncate cell-muted" :title="log.resource_id">{{ log.resource_id }}</span>
                  <span v-else class="cell-muted">—</span>
                </td>
                <td class="cell-muted cell-sm">{{ log.ip_address || '—' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination">
        <p class="pagination__info">
          Showing {{ (page - 1) * limit + 1 }}–{{ Math.min(page * limit, total) }} of {{ total }} events
        </p>
        <div class="pagination__controls">
          <button class="btn-outlined" :disabled="page <= 1" @click="changePage(page - 1)">
            <span class="material-symbols-outlined">chevron_left</span>
            Prev
          </button>
          <span class="pagination__label">Page {{ page }}</span>
          <button class="btn-outlined" :disabled="page * limit >= total" @click="changePage(page + 1)">
            Next
            <span class="material-symbols-outlined">chevron_right</span>
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import AppSelect from '@/components/AppSelect.vue';
import auditApi from '@/api/audit';
import type { AuditLog } from '@/types/audit';
import { ACTION_LABELS, RESOURCE_LABELS, ACTION_BADGE } from '@/types/audit';

const logs = ref<AuditLog[]>([]);
const total = ref(0);
const page = ref(1);
const limit = 20;
const loading = ref(true);
const exporting = ref(false);

const filters = reactive({
  action: '',
  resource_type: '',
  from: '',
  to: '',
});

async function fetchLogs() {
  loading.value = true;
  try {
    const params: Record<string, unknown> = { page: page.value, limit };
    if (filters.action) params.action = filters.action;
    if (filters.resource_type) params.resource_type = filters.resource_type;
    if (filters.from) params.from = new Date(filters.from).toISOString();
    if (filters.to) {
      const d = new Date(filters.to);
      d.setHours(23, 59, 59, 999);
      params.to = d.toISOString();
    }
    const res = await auditApi.list(params as Parameters<typeof auditApi.list>[0]);
    const data = res.data.data!;
    logs.value = data.logs ?? [];
    total.value = data.total ?? 0;
  } catch {
    logs.value = [];
  } finally {
    loading.value = false;
  }
}

function applyFilters() {
  page.value = 1;
  fetchLogs();
}

function resetFilters() {
  filters.action = '';
  filters.resource_type = '';
  filters.from = '';
  filters.to = '';
  page.value = 1;
  fetchLogs();
}

function changePage(n: number) {
  page.value = n;
  fetchLogs();
}

function actionLabel(action: string) {
  return ACTION_LABELS[action] ?? action.replace(/_/g, ' ');
}

function resourceLabel(type: string) {
  return RESOURCE_LABELS[type] ?? type;
}

function actionBadge(action: string) {
  const raw = ACTION_BADGE[action] ?? '';
  if (raw.includes('success')) return 'm3-badge--success';
  if (raw.includes('danger')) return 'm3-badge--error';
  if (raw.includes('warning')) return 'm3-badge--warning';
  if (raw.includes('primary')) return 'm3-badge--primary';
  return 'm3-badge--neutral';
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  });
}

async function exportLogs() {
  exporting.value = true;
  try {
    const fromIso = filters.from ? new Date(filters.from).toISOString() : undefined;
    let toIso: string | undefined;
    if (filters.to) {
      const d = new Date(filters.to);
      d.setHours(23, 59, 59, 999);
      toIso = d.toISOString();
    }
    await auditApi.exportAuditLogs(filters.action || undefined, fromIso, toIso);
  } catch {
    // silently ignore export errors
  } finally {
    exporting.value = false;
  }
}

onMounted(fetchLogs);
</script>

<style scoped lang="scss">
.audit-page {
  max-width: 1100px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page Header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-header__left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-header__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  margin: 0;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── CSS Spinner ─────────────────────────────────────────────────────────── */
.css-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm {
    width: 16px;
    height: 16px;
    border-width: 2px;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Cards ───────────────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);

  &__left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
  }
}

.an-card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 18px; }

  &--primary {
    background: color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
    color: var(--md-sys-color-primary);
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

/* ── Filters ─────────────────────────────────────────────────────────────── */
.filters-card {
  overflow: visible;
}

.filters-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: flex-end;
  padding: 16px 20px;
}

.filter-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

/* ── Form field ──────────────────────────────────────────────────────────── */
.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field__label {
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
}

.form-input {
  height: 40px;
  padding: 0 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9rem;
  outline: none;
  min-width: 148px;

  &:focus {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
  }
}

/* ── Loading ─────────────────────────────────────────────────────────────── */
.loading-center {
  display: flex;
  justify-content: center;
  padding: 48px;
}

/* ── Empty state ─────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 24px;
  text-align: center;
}

.empty-icon {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;

  .material-symbols-outlined {
    font-size: 2rem;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.6;
  }
}

.empty-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--md-sys-color-on-surface);
}

.empty-desc {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Table ───────────────────────────────────────────────────────────────── */
.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;

  thead tr {
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  th {
    padding: 12px 16px;
    text-align: left;
    font-weight: 600;
    font-size: 0.8rem;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    white-space: nowrap;
  }

  td {
    padding: 12px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  tbody tr:hover td {
    background: var(--md-sys-color-surface-container-low);
  }
}

.cell-muted { color: var(--md-sys-color-on-surface-variant); }
.cell-nowrap { white-space: nowrap; }
.cell-sm { font-size: 0.8rem; }
.cell-capitalize { text-transform: capitalize; }
.cell-strong { font-weight: 500; }
.cell-truncate {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}

/* ── Badges ──────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--primary {
    background: rgba(99, 91, 255, 0.12);
    color: var(--md-sys-color-primary);
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  &--success {
    background: rgba(22, 163, 74, 0.12);
    color: #16a34a;
  }

  &--warning {
    background: rgba(245, 158, 11, 0.12);
    color: #b45309;
  }

  &--error {
    background: rgba(220, 38, 38, 0.12);
    color: #dc2626;
  }
}

/* ── Pagination ──────────────────────────────────────────────────────────── */
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;

  &__info {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.875rem;
    margin: 0;
  }

  &__controls {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  &__label {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface-variant);
    min-width: 56px;
    text-align: center;
  }
}

/* ── Button helper ───────────────────────────────────────────────────────── */
.btn-icon-left {
  font-size: 18px;
  margin-right: 4px;
}
</style>
