<template>
  <div class="container-fluid py-4" style="max-width: 1100px;">
    <!-- Impersonation Banner -->
    <div v-if="isImpersonating" class="alert alert-warning d-flex align-items-center gap-3 mb-4">
      <i class="bi bi-person-fill-gear fs-5"></i>
      <span class="flex-grow-1">You are currently impersonating a user.</span>
      <button class="btn btn-sm btn-warning" @click="stopImpersonation">Stop Impersonation</button>
    </div>

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
      <div class="col-6 col-lg-3">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body px-4 py-3">
            <div class="text-muted small mb-1">&#128070; Total Clicks</div>
            <div class="fw-bold fs-4">{{ stats != null ? stats.total_clicks.toLocaleString() : '—' }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Platform Overview charts ───────────────────────────────────────── -->
    <div class="row g-3 mb-4" v-if="stats">
      <!-- Horizontal bar: Platform metrics comparison -->
      <div class="col-12 col-lg-8">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-header bg-white border-bottom py-3 px-4">
            <span class="fw-semibold">Platform Overview</span>
          </div>
          <div class="card-body px-4 py-3">
            <VChart :option="platformBarOption" style="height: 160px;" autoresize />
          </div>
        </div>
      </div>

      <!-- Donut ring: Active vs Inactive users -->
      <div class="col-12 col-lg-4">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-header bg-white border-bottom py-3 px-4">
            <span class="fw-semibold">User Activity</span>
          </div>
          <div class="card-body px-4 py-3 d-flex flex-column align-items-center justify-content-center">
            <VChart :option="userRingOption" style="height: 160px; width: 100%;" autoresize />
            <div class="d-flex gap-3 mt-2" style="font-size: 0.78rem;">
              <span><span class="d-inline-block rounded-circle me-1" style="width:10px;height:10px;background:#22c55e;"></span>Active</span>
              <span><span class="d-inline-block rounded-circle me-1" style="width:10px;height:10px;background:#f59e0b;"></span>Inactive</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Platform Growth Chart ──────────────────────────────────────────── -->
    <div class="row g-3 mb-4">
      <div class="col-12">
        <div class="card border-0 shadow-sm">
          <div class="card-header bg-white border-bottom py-3 px-4">
            <span class="fw-semibold">Platform Growth (Last 30 Days)</span>
          </div>
          <div class="card-body px-4 py-3">
            <div v-if="growthLoading" class="text-center py-4">
              <div class="spinner-border spinner-border-sm text-primary"></div>
            </div>
            <VChart v-else :option="growthLineOption" style="height: 240px;" autoresize />
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
                <div class="d-flex gap-2 justify-content-end align-items-center">
                  <!-- Role change -->
                  <select
                    class="form-select form-select-sm"
                    style="width: auto; min-width: 90px;"
                    :value="u.role"
                    :disabled="changingRole === u.id || u.id === authStore.profile?.id"
                    @change="setRole(u, ($event.target as HTMLSelectElement).value)"
                  >
                    <option value="user">User</option>
                    <option value="admin">Admin</option>
                  </select>
                  <!-- Impersonate -->
                  <button
                    class="btn btn-sm btn-outline-secondary"
                    :disabled="impersonating === u.id || u.id === authStore.profile?.id"
                    @click="startImpersonation(u)"
                    title="Impersonate user"
                  >
                    <span v-if="impersonating === u.id" class="spinner-border spinner-border-sm"></span>
                    <i v-else class="bi bi-person-fill-exclamation"></i>
                  </button>
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
import { useRouter } from 'vue-router';
import { use } from 'echarts/core';
import { BarChart, PieChart, LineChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import VChart from 'vue-echarts';
import adminApi from '@/api/admin';
import type { AdminStats, AdminUser, GrowthTimeSeriesPoint } from '@/types/admin';
import { useAuthStore } from '@/stores/auth';

use([BarChart, PieChart, LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer]);

const LIMIT = 20;
const ORIGINAL_TOKEN_KEY = 'impersonation_original_token';

const authStore = useAuthStore();
const router = useRouter();

const stats = ref<AdminStats | null>(null);
const users = ref<AdminUser[]>([]);
const totalUsers = ref(0);
const currentPage = ref(1);
const usersLoading = ref(false);
const search = ref('');
const updatingId = ref<string | null>(null);
const growthLoading = ref(true);
const growthUsers = ref<GrowthTimeSeriesPoint[]>([]);
const growthLinks = ref<GrowthTimeSeriesPoint[]>([]);
const impersonating = ref<string | null>(null);
const changingRole = ref<string | null>(null);

let searchTimer: ReturnType<typeof setTimeout> | null = null;

const totalPages = computed(() => Math.max(1, Math.ceil(totalUsers.value / LIMIT)));

const platformBarOption = computed(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { top: 8, right: 16, bottom: 8, left: 100, containLabel: false },
  xAxis: { type: 'value', axisLabel: { color: '#697386', fontSize: 11 } },
  yAxis: {
    type: 'category',
    data: ['Total Links', 'Inactive Users', 'Active Users', 'Total Users'],
    axisLabel: { color: '#697386', fontSize: 11 },
  },
  series: [{
    type: 'bar',
    data: [
      { value: stats.value?.total_links ?? 0, itemStyle: { color: '#635bff' } },
      { value: stats.value?.inactive_users ?? 0, itemStyle: { color: '#f59e0b' } },
      { value: stats.value?.active_users ?? 0, itemStyle: { color: '#22c55e' } },
      { value: stats.value?.total_users ?? 0, itemStyle: { color: '#635bff', opacity: 0.5 } },
    ],
    barMaxWidth: 28,
    borderRadius: [0, 4, 4, 0],
  }],
}));

const userRingOption = computed(() => ({
  tooltip: { trigger: 'item' },
  series: [{
    type: 'pie',
    radius: ['55%', '80%'],
    label: { show: false },
    data: [
      { name: 'Active', value: stats.value?.active_users ?? 0, itemStyle: { color: '#22c55e' } },
      { name: 'Inactive', value: stats.value?.inactive_users ?? 0, itemStyle: { color: '#f59e0b' } },
    ],
  }],
}));

const growthLineOption = computed(() => {
  const formatLabel = (ts: string) => {
    const d = new Date(ts);
    return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
  };

  const userDates = growthUsers.value.map(p => formatLabel(p.timestamp));
  const linkDates = growthLinks.value.map(p => formatLabel(p.timestamp));
  // Merge and deduplicate x-axis dates
  const allDates = Array.from(new Set([...userDates, ...linkDates])).sort((a, b) =>
    new Date(a).getTime() - new Date(b).getTime()
  );

  // Build lookup maps for fast access
  const userMap = Object.fromEntries(growthUsers.value.map(p => [formatLabel(p.timestamp), p.count]));
  const linkMap = Object.fromEntries(growthLinks.value.map(p => [formatLabel(p.timestamp), p.count]));

  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['New Users', 'New Links'], bottom: 0, textStyle: { color: '#697386', fontSize: 12 } },
    grid: { top: 16, right: 16, bottom: 36, left: 40 },
    xAxis: {
      type: 'category',
      data: allDates,
      axisLabel: { color: '#697386', fontSize: 11, rotate: 30 },
      axisLine: { lineStyle: { color: '#e2e8f0' } },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: '#697386', fontSize: 11 },
      splitLine: { lineStyle: { color: '#f1f5f9' } },
    },
    series: [
      {
        name: 'New Users',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 5,
        data: allDates.map(d => userMap[d] ?? 0),
        lineStyle: { color: '#635bff', width: 2 },
        itemStyle: { color: '#635bff' },
        areaStyle: { color: 'rgba(99,91,255,0.08)' },
      },
      {
        name: 'New Links',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 5,
        data: allDates.map(d => linkMap[d] ?? 0),
        lineStyle: { color: '#14b8a6', width: 2 },
        itemStyle: { color: '#14b8a6' },
        areaStyle: { color: 'rgba(20,184,166,0.08)' },
      },
    ],
  };
});

async function loadStats() {
  try {
    const res = await adminApi.getStats();
    stats.value = res.data?.data ?? null;
  } catch {
    // non-critical
  }
}

async function loadGrowth() {
  growthLoading.value = true;
  try {
    const res = await adminApi.getGrowthTimeSeries();
    growthUsers.value = res.data?.data?.users ?? [];
    growthLinks.value = res.data?.data?.links ?? [];
  } catch {
    // non-critical
  } finally {
    growthLoading.value = false;
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

async function setRole(user: AdminUser, role: string) {
  changingRole.value = user.id;
  try {
    await adminApi.updateUserRole(user.id, role);
    user.role = role;
  } catch {
    // silently fail — user will see no change
  } finally {
    changingRole.value = null;
  }
}

async function startImpersonation(user: AdminUser) {
  impersonating.value = user.id;
  try {
    const res = await adminApi.impersonateUser(user.id);
    const impersonationToken = res.data.data.access_token;
    const currentToken = localStorage.getItem('access_token');
    if (currentToken) {
      localStorage.setItem(ORIGINAL_TOKEN_KEY, currentToken);
    }
    localStorage.setItem('access_token', impersonationToken);
    window.location.href = '/dashboard/overview';
  } catch {
    impersonating.value = null;
  }
}

const isImpersonating = computed(() => !!localStorage.getItem(ORIGINAL_TOKEN_KEY));

function stopImpersonation() {
  const originalToken = localStorage.getItem(ORIGINAL_TOKEN_KEY);
  if (originalToken) {
    localStorage.setItem('access_token', originalToken);
    localStorage.removeItem(ORIGINAL_TOKEN_KEY);
    window.location.href = '/admin';
  }
}

onMounted(() => {
  loadStats();
  loadUsers();
  loadGrowth();
});
</script>

<style scoped>
.btn-outline-success { --bs-btn-hover-color: #fff; }
</style>
