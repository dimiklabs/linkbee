<template>
  <div class="container-fluid py-4" style="max-width: 1100px;">
    <div class="mb-4">
      <h4 class="mb-1 fw-bold">Admin Dashboard</h4>
      <p class="text-muted small mb-0">Platform overview and user management.</p>
    </div>

    <!-- ── Stats ──────────────────────────────────────────────────────────── -->
    <div class="row g-3 mb-4">
      <div class="col-6 col-lg-3">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body px-4 py-3">
            <div class="text-muted small mb-1">Total Users</div>
            <div class="fw-bold fs-4">{{ stats?.total_users ?? '—' }}</div>
          </div>
        </div>
      </div>
      <div class="col-6 col-lg-3">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body px-4 py-3">
            <div class="text-muted small mb-1">Active Users</div>
            <div class="fw-bold fs-4 text-success">{{ stats?.active_users ?? '—' }}</div>
          </div>
        </div>
      </div>
      <div class="col-6 col-lg-3">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body px-4 py-3">
            <div class="text-muted small mb-1">Inactive Users</div>
            <div class="fw-bold fs-4 text-warning">{{ stats?.inactive_users ?? '—' }}</div>
          </div>
        </div>
      </div>
      <div class="col-6 col-lg-3">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body px-4 py-3">
            <div class="text-muted small mb-1">Total Links</div>
            <div class="fw-bold fs-4">{{ stats?.total_links ?? '—' }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Users Table ────────────────────────────────────────────────────── -->
    <div class="card border-0 shadow-sm">
      <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between gap-3 flex-wrap">
        <span class="fw-semibold">Users</span>
        <div style="max-width: 280px; width: 100%;">
          <input
            v-model="search"
            type="search"
            class="form-control form-control-sm"
            placeholder="Search by email or name…"
            @input="onSearchInput"
          />
        </div>
      </div>

      <div v-if="usersLoading" class="text-center py-5">
        <div class="spinner-border text-primary spinner-border-sm"></div>
      </div>

      <div v-else-if="users.length === 0" class="text-center text-muted py-5 small">
        No users found.
      </div>

      <div v-else class="table-responsive">
        <table class="table table-hover mb-0 align-middle">
          <thead class="table-light">
            <tr>
              <th class="px-4 py-3 fw-medium small text-muted">User</th>
              <th class="py-3 fw-medium small text-muted">Provider</th>
              <th class="py-3 fw-medium small text-muted">Role</th>
              <th class="py-3 fw-medium small text-muted">Status</th>
              <th class="py-3 fw-medium small text-muted">Joined</th>
              <th class="py-3 fw-medium small text-muted">Last Login</th>
              <th class="py-3 fw-medium small text-muted text-end px-4">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td class="px-4 py-3">
                <div class="fw-medium small">{{ displayName(u) }}</div>
                <div class="text-muted" style="font-size: 0.78rem;">{{ u.email }}</div>
              </td>
              <td class="py-3">
                <span class="badge bg-light text-secondary border text-capitalize" style="font-size: 0.72rem;">
                  {{ u.auth_provider }}
                </span>
              </td>
              <td class="py-3">
                <span
                  class="badge"
                  :class="u.role === 'admin' ? 'text-bg-primary' : 'bg-light text-secondary border'"
                  style="font-size: 0.72rem;"
                >{{ u.role }}</span>
              </td>
              <td class="py-3">
                <span
                  class="badge"
                  :class="statusBadge(u.status)"
                  style="font-size: 0.72rem;"
                >{{ u.status }}</span>
              </td>
              <td class="py-3 small text-muted">{{ formatDate(u.created_at) }}</td>
              <td class="py-3 small text-muted">{{ u.last_login ? formatDate(u.last_login) : '—' }}</td>
              <td class="py-3 px-4 text-end">
                <div class="d-flex gap-2 justify-content-end">
                  <button
                    v-if="u.status !== 'active'"
                    class="btn btn-sm btn-outline-success"
                    :disabled="updatingId === u.id"
                    @click="setStatus(u, 'active')"
                  >
                    <span v-if="updatingId === u.id" class="spinner-border spinner-border-sm"></span>
                    <span v-else>Activate</span>
                  </button>
                  <button
                    v-if="u.status !== 'banned'"
                    class="btn btn-sm btn-outline-danger"
                    :disabled="updatingId === u.id"
                    @click="setStatus(u, 'banned')"
                  >
                    <span v-if="updatingId === u.id" class="spinner-border spinner-border-sm"></span>
                    <span v-else>Ban</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="card-footer bg-white border-top px-4 py-3 d-flex align-items-center justify-content-between">
        <span class="text-muted small">
          Page {{ currentPage }} of {{ totalPages }} &mdash; {{ totalUsers }} users
        </span>
        <div class="d-flex gap-2">
          <button
            class="btn btn-sm btn-outline-secondary"
            :disabled="currentPage <= 1"
            @click="goToPage(currentPage - 1)"
          >Previous</button>
          <button
            class="btn btn-sm btn-outline-secondary"
            :disabled="currentPage >= totalPages"
            @click="goToPage(currentPage + 1)"
          >Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import adminApi from '@/api/admin';
import type { AdminStats, AdminUser } from '@/types/admin';

const LIMIT = 20;

const stats = ref<AdminStats | null>(null);
const users = ref<AdminUser[]>([]);
const totalUsers = ref(0);
const currentPage = ref(1);
const usersLoading = ref(false);
const search = ref('');
const updatingId = ref<string | null>(null);

let searchTimer: ReturnType<typeof setTimeout> | null = null;

const totalPages = computed(() => Math.max(1, Math.ceil(totalUsers.value / LIMIT)));

async function loadStats() {
  try {
    const res = await adminApi.getStats();
    stats.value = res.data?.data ?? null;
  } catch {
    // non-critical
  }
}

async function loadUsers() {
  usersLoading.value = true;
  try {
    const res = await adminApi.listUsers({ page: currentPage.value, limit: LIMIT, search: search.value || undefined });
    const data = res.data?.data;
    users.value = data?.users ?? [];
    totalUsers.value = data?.total ?? 0;
  } finally {
    usersLoading.value = false;
  }
}

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    currentPage.value = 1;
    loadUsers();
  }, 400);
}

function goToPage(page: number) {
  currentPage.value = page;
  loadUsers();
}

async function setStatus(user: AdminUser, status: string) {
  updatingId.value = user.id;
  try {
    await adminApi.updateUserStatus(user.id, status);
    user.status = status;
    await loadStats();
  } finally {
    updatingId.value = null;
  }
}

function displayName(u: AdminUser) {
  if (u.first_name || u.last_name) return [u.first_name, u.last_name].filter(Boolean).join(' ');
  return u.email;
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
}

function statusBadge(status: string) {
  switch (status) {
    case 'active':   return 'text-bg-success';
    case 'inactive': return 'text-bg-warning';
    case 'banned':   return 'text-bg-danger';
    default:         return 'bg-light text-secondary border';
  }
}

onMounted(() => {
  loadStats();
  loadUsers();
});
</script>

<style scoped>
.btn-outline-success { --bs-btn-hover-color: #fff; }
</style>
