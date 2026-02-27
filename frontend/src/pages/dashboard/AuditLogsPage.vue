<template>
  <div class="container-fluid py-4">

    <!-- Header -->
    <div class="mb-4">
      <h1 class="h4 fw-bold mb-0">Audit Logs</h1>
      <p class="text-muted small mb-0">A record of all security-relevant actions in your account.</p>
    </div>

    <!-- Filters -->
    <div class="card mb-4">
      <div class="card-body py-3">
        <div class="row g-2 align-items-end">
          <div class="col-sm-6 col-md-3">
            <label class="form-label small mb-1">Action</label>
            <select v-model="filters.action" class="form-select form-select-sm">
              <option value="">All actions</option>
              <option v-for="(label, key) in ACTION_LABELS" :key="key" :value="key">{{ label }}</option>
            </select>
          </div>
          <div class="col-sm-6 col-md-2">
            <label class="form-label small mb-1">Resource</label>
            <select v-model="filters.resource_type" class="form-select form-select-sm">
              <option value="">All resources</option>
              <option v-for="(label, key) in RESOURCE_LABELS" :key="key" :value="key">{{ label }}</option>
            </select>
          </div>
          <div class="col-sm-6 col-md-2">
            <label class="form-label small mb-1">From</label>
            <input v-model="filters.from" type="date" class="form-control form-control-sm" />
          </div>
          <div class="col-sm-6 col-md-2">
            <label class="form-label small mb-1">To</label>
            <input v-model="filters.to" type="date" class="form-control form-control-sm" />
          </div>
          <div class="col-sm-6 col-md-3 d-flex gap-2">
            <button class="btn btn-sm btn-primary flex-fill" @click="applyFilters">Apply</button>
            <button class="btn btn-sm btn-outline-secondary" @click="resetFilters">Reset</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading…</span>
      </div>
    </div>

    <!-- Empty -->
    <div v-else-if="logs.length === 0" class="card text-center py-5">
      <div class="card-body">
        <p class="fs-1 mb-2">📋</p>
        <h5 class="mb-1">No audit logs found</h5>
        <p class="text-muted mb-0">Actions you take will be recorded here.</p>
      </div>
    </div>

    <!-- Table -->
    <div v-else>
      <div class="card">
        <div class="table-responsive">
          <table class="table table-hover mb-0 align-middle" style="font-size: 0.875rem;">
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
                <td class="text-nowrap text-muted small">{{ formatDate(log.created_at) }}</td>
                <td>
                  <span :class="['badge', actionBadge(log.action)]">
                    {{ actionLabel(log.action) }}
                  </span>
                </td>
                <td class="text-muted small text-capitalize">{{ resourceLabel(log.resource_type) }}</td>
                <td class="small">
                  <span v-if="log.resource_name" class="fw-medium">{{ log.resource_name }}</span>
                  <span v-else-if="log.resource_id" class="text-muted text-truncate d-inline-block" style="max-width:160px;" :title="log.resource_id">
                    {{ log.resource_id }}
                  </span>
                  <span v-else class="text-muted">—</span>
                </td>
                <td class="text-muted small">{{ log.ip_address || '—' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="d-flex align-items-center justify-content-between mt-3 flex-wrap gap-2">
        <p class="text-muted small mb-0">
          Showing {{ (page - 1) * limit + 1 }}–{{ Math.min(page * limit, total) }} of {{ total }} events
        </p>
        <div class="d-flex gap-1">
          <button class="btn btn-sm btn-outline-secondary" :disabled="page <= 1" @click="changePage(page - 1)">
            ← Prev
          </button>
          <button class="btn btn-sm btn-outline-secondary" :disabled="page * limit >= total" @click="changePage(page + 1)">
            Next →
          </button>
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
  return ACTION_BADGE[action] ?? 'text-bg-secondary';
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  });
}

onMounted(fetchLogs);
</script>
