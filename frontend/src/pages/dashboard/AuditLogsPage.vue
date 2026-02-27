<template>
  <div class="page-section" style="max-width: 1100px;">

    <!-- Page Header -->
    <div class="dash-page-header">
      <div class="dash-page-header__left">
        <h1 class="dash-page-header__title">Audit Logs</h1>
        <p class="dash-page-header__subtitle">A record of all security-relevant actions in your account.</p>
      </div>
      <div class="dash-page-header__actions">
        <md-outlined-button :disabled="exporting" @click="exportLogs">
          <span class="material-symbols-outlined" style="font-size:18px;margin-right:6px;">download</span>
          <md-circular-progress v-if="exporting" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
          Export CSV
        </md-outlined-button>
      </div>
    </div>

    <!-- Filters -->
    <div class="m3-card m3-card--elevated filters-card">
      <div class="m3-card-header">
        <div style="display:flex;align-items:center;gap:8px;">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">filter_list</span>
          <span style="font-size:16px;font-weight:600;">Filter Events</span>
        </div>
      </div>
      <md-divider />
      <div class="filters-row">
        <AppSelect
          v-model="filters.action"
          label="Action"
          style="min-width:160px;"
        >
          <option value="">All actions</option>
          <option v-for="(label, key) in ACTION_LABELS" :key="key" :value="key">{{ label }}</option>
        </AppSelect>

        <AppSelect
          v-model="filters.resource_type"
          label="Resource"
          style="min-width:140px;"
        >
          <option value="">All resources</option>
          <option v-for="(label, key) in RESOURCE_LABELS" :key="key" :value="key">{{ label }}</option>
        </AppSelect>

        <md-outlined-text-field
          :value="filters.from"
          @input="filters.from=($event.target as HTMLInputElement).value"
          label="From"
          type="date"
          style="min-width:160px;"
        />

        <md-outlined-text-field
          :value="filters.to"
          @input="filters.to=($event.target as HTMLInputElement).value"
          label="To"
          type="date"
          style="min-width:160px;"
        />

        <div style="display:flex;gap:8px;align-items:center;">
          <md-filled-button @click="applyFilters">
            <span class="material-symbols-outlined" style="font-size:18px;margin-right:6px;">search</span>
            Apply
          </md-filled-button>
          <md-outlined-button @click="resetFilters">Reset</md-outlined-button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:40px" />
    </div>

    <!-- Empty -->
    <div v-else-if="logs.length === 0" class="m3-card m3-card--elevated m3-empty-state">
      <div class="m3-empty-state__icon">
        <span class="material-symbols-outlined">event_note</span>
      </div>
      <div class="md-title-medium" style="margin-bottom:8px;">No audit logs found</div>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0;">Actions you take will be recorded here for security and compliance.</p>
    </div>

    <!-- Table -->
    <div v-else>
      <div class="m3-card m3-card--elevated">
        <div class="m3-card-header">
          <div style="display:flex;align-items:center;gap:8px;">
            <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">security</span>
            <span style="font-size:16px;font-weight:600;">Audit Events</span>
          </div>
          <span class="m3-badge m3-badge--neutral">{{ total }} total</span>
        </div>
        <md-divider />
        <div class="m3-table-wrapper">
          <table class="m3-table">
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
                <td style="white-space:nowrap;color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;">{{ formatDate(log.created_at) }}</td>
                <td>
                  <span :class="['m3-badge', actionBadge(log.action)]">
                    {{ actionLabel(log.action) }}
                  </span>
                </td>
                <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;text-transform:capitalize;">{{ resourceLabel(log.resource_type) }}</td>
                <td style="font-size:0.875rem;">
                  <span v-if="log.resource_name" style="font-weight:500;">{{ log.resource_name }}</span>
                  <span v-else-if="log.resource_id" style="color:var(--md-sys-color-on-surface-variant);max-width:160px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;display:inline-block;" :title="log.resource_id">
                    {{ log.resource_id }}
                  </span>
                  <span v-else style="color:var(--md-sys-color-on-surface-variant);">—</span>
                </td>
                <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;">{{ log.ip_address || '—' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="m3-pagination">
        <p class="m3-pagination__info">
          Showing {{ (page - 1) * limit + 1 }}–{{ Math.min(page * limit, total) }} of {{ total }} events
        </p>
        <div class="m3-pagination__controls">
          <md-outlined-button :disabled="page <= 1" @click="changePage(page - 1)">
            <span class="material-symbols-outlined" style="font-size:18px;margin-right:4px;">chevron_left</span>
            Prev
          </md-outlined-button>
          <span class="m3-pagination__page-label">Page {{ page }}</span>
          <md-outlined-button :disabled="page * limit >= total" @click="changePage(page + 1)">
            Next
            <span class="material-symbols-outlined" style="font-size:18px;margin-left:4px;">chevron_right</span>
          </md-outlined-button>
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
/* page-section (global) handles padding; max-width set via style attribute on root */

/* ── Page Header ──────────────────────────────────────────────────────────── */
.dash-page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;

  &__left {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  &__title {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  &__subtitle {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.875rem;
    margin: 0;
  }

  &__actions {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-shrink: 0;
  }
}

/* ── Cards ────────────────────────────────────────────────────────────────── */
.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
  margin-bottom: 24px;

  &--outlined {
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
  }
}

/* ── M3 Card header ───────────────────────────────────────────────────────── */
.m3-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  gap: 1rem;
  flex-wrap: wrap;
}

/* ── Filters card ─────────────────────────────────────────────────────────── */
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

/* ── M3 Table wrapper ─────────────────────────────────────────────────────── */
.m3-table-wrapper {
  overflow-x: auto;
}

/* ── M3 Empty state ───────────────────────────────────────────────────────── */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 24px;
  text-align: center;

  &__icon {
    width: 72px;
    height: 72px;
    border-radius: 50%;
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
}

/* ── Table ────────────────────────────────────────────────────────────────── */
.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.m3-table thead tr {
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.m3-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  background: var(--md-sys-color-surface-container-low);
  white-space: nowrap;
}

.m3-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  color: var(--md-sys-color-on-surface);
}

.m3-table tbody tr:last-child td {
  border-bottom: none;
}

.m3-table tbody tr:hover td {
  background: var(--md-sys-color-surface-container-low);
}

/* ── M3 Badges ────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
}

.m3-badge--primary {
  background: rgba(99, 91, 255, 0.12);
  color: var(--md-sys-color-primary);
}

.m3-badge--neutral {
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  border: 1px solid var(--md-sys-color-outline-variant);
}

.m3-badge--success {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
}

.m3-badge--warning {
  background: rgba(245, 158, 11, 0.12);
  color: #b45309;
}

.m3-badge--error {
  background: rgba(220, 38, 38, 0.12);
  color: #dc2626;
}

/* ── Pagination ───────────────────────────────────────────────────────────── */
.m3-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
  margin-bottom: 8px;
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

  &__page-label {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface-variant);
    min-width: 56px;
    text-align: center;
  }
}
</style>
