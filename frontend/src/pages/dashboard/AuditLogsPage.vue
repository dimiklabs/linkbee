<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Audit Logs</h1>
        <p class="page-subtitle">A record of all security-relevant actions in your account.</p>
      </div>
      <md-outlined-button :disabled="exporting" @click="exportLogs">
        <span class="material-symbols-outlined" slot="icon">download</span>
        <md-circular-progress v-if="exporting" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
        Export CSV
      </md-outlined-button>
    </div>

    <!-- Filters -->
    <div class="m3-card m3-card--outlined filters-card">
      <div class="filters-row">
        <md-outlined-select
          :value="filters.action"
          @change="filters.action=($event.target as HTMLSelectElement).value"
          label="Action"
          style="min-width:160px;"
        >
          <md-select-option value=""><div slot="headline">All actions</div></md-select-option>
          <md-select-option v-for="(label, key) in ACTION_LABELS" :key="key" :value="key">
            <div slot="headline">{{ label }}</div>
          </md-select-option>
        </md-outlined-select>

        <md-outlined-select
          :value="filters.resource_type"
          @change="filters.resource_type=($event.target as HTMLSelectElement).value"
          label="Resource"
          style="min-width:140px;"
        >
          <md-select-option value=""><div slot="headline">All resources</div></md-select-option>
          <md-select-option v-for="(label, key) in RESOURCE_LABELS" :key="key" :value="key">
            <div slot="headline">{{ label }}</div>
          </md-select-option>
        </md-outlined-select>

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
          <md-filled-button @click="applyFilters">Apply</md-filled-button>
          <md-outlined-button @click="resetFilters">Reset</md-outlined-button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:40px" />
    </div>

    <!-- Empty -->
    <div v-else-if="logs.length === 0" class="m3-card m3-card--outlined empty-state">
      <span class="material-symbols-outlined" style="font-size:3rem;color:var(--md-sys-color-on-surface-variant);">event_note</span>
      <h5 style="margin:12px 0 4px;font-weight:600;">No audit logs found</h5>
      <p style="color:var(--md-sys-color-on-surface-variant);margin:0;font-size:0.875rem;">Actions you take will be recorded here.</p>
    </div>

    <!-- Table -->
    <div v-else>
      <div class="m3-card m3-card--outlined">
        <div class="table-container">
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
      <div class="pagination-bar">
        <p style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;margin:0;">
          Showing {{ (page - 1) * limit + 1 }}–{{ Math.min(page * limit, total) }} of {{ total }} events
        </p>
        <div style="display:flex;gap:8px;">
          <md-outlined-button :disabled="page <= 1" @click="changePage(page - 1)">
            <span class="material-symbols-outlined" slot="icon">chevron_left</span>
            Prev
          </md-outlined-button>
          <md-outlined-button :disabled="page * limit >= total" @click="changePage(page + 1)">
            Next
            <span class="material-symbols-outlined" slot="trailing-icon">chevron_right</span>
          </md-outlined-button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
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

<style scoped>
.page-wrapper {
  padding: 24px;
  max-width: 1100px;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 4px;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0;
}

.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
}

.m3-card--outlined {
  border: 1px solid var(--md-sys-color-outline-variant);
  margin-bottom: 24px;
}

.filters-card {
  padding: 16px 20px;
}

.filters-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: flex-end;
}

.table-container {
  overflow-x: auto;
}

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
  background: var(--md-sys-color-primary-container, #e8def8);
  color: var(--md-sys-color-on-primary-container, #21005d);
}

.m3-badge--neutral {
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  border: 1px solid var(--md-sys-color-outline-variant);
}

.m3-badge--success {
  background: #dcfce7;
  color: #16a34a;
}

.m3-badge--warning {
  background: #fef3c7;
  color: #b45309;
}

.m3-badge--error {
  background: #fee2e2;
  color: #dc2626;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 24px;
  text-align: center;
}

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
