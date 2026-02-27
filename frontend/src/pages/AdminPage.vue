<template>
  <div class="page-wrapper">

    <!-- Impersonation Banner -->
    <div v-if="isImpersonating" class="impersonation-banner">
      <span class="material-symbols-outlined" style="font-size:20px;color:#92400e;flex-shrink:0;">manage_accounts</span>
      <span class="md-body-medium" style="flex:1;color:#78350f;">You are currently impersonating a user. Actions taken will affect the impersonated account.</span>
      <md-outlined-button
        @click="stopImpersonation"
        style="--md-outlined-button-label-text-color:#92400e;--md-outlined-button-outline-color:#d97706;"
      >
        <span class="material-symbols-outlined" slot="icon">close</span>
        Stop impersonating
      </md-outlined-button>
    </div>

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Admin Dashboard</h1>
        <p class="page-subtitle">Platform overview and user management.</p>
      </div>
    </div>

    <!-- Quick stats bar -->
    <div v-if="stats" class="quick-stats-bar">
      <div class="quick-stat-item">
        <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-primary);">group</span>
        <span class="md-label-large">{{ stats.total_users.toLocaleString() }}</span>
        <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Total Users</span>
      </div>
      <div class="quick-stat-divider"></div>
      <div class="quick-stat-item">
        <span class="material-symbols-outlined" style="font-size:16px;color:#16a34a;">person_check</span>
        <span class="md-label-large" style="color:#16a34a;">{{ stats.active_users.toLocaleString() }}</span>
        <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Active</span>
      </div>
      <div class="quick-stat-divider"></div>
      <div class="quick-stat-item">
        <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-primary);">link</span>
        <span class="md-label-large">{{ stats.total_links.toLocaleString() }}</span>
        <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Links</span>
      </div>
      <div class="quick-stat-divider"></div>
      <div class="quick-stat-item">
        <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-primary);">ads_click</span>
        <span class="md-label-large">{{ stats.total_clicks.toLocaleString() }}</span>
        <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Clicks</span>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="m3-card m3-card--elevated stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:var(--md-sys-color-primary)">group</span>
        <div class="stat-value">{{ stats?.total_users ?? '—' }}</div>
        <div class="stat-label">Total Users</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:#22c55e">person_check</span>
        <div class="stat-value" style="color:#22c55e">{{ stats?.active_users ?? '—' }}</div>
        <div class="stat-label">Active Users</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:#f59e0b">person_off</span>
        <div class="stat-value" style="color:#f59e0b">{{ stats?.inactive_users ?? '—' }}</div>
        <div class="stat-label">Inactive Users</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:var(--md-sys-color-primary)">link</span>
        <div class="stat-value">{{ stats?.total_links ?? '—' }}</div>
        <div class="stat-label">Total Links</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:var(--md-sys-color-primary)">ads_click</span>
        <div class="stat-value">{{ stats != null ? stats.total_clicks.toLocaleString() : '—' }}</div>
        <div class="stat-label">Total Clicks</div>
      </div>
    </div>

    <!-- Charts Row -->
    <div v-if="stats" class="charts-row">
      <div class="m3-card m3-card--elevated chart-card chart-card--wide">
        <div class="chart-card-header">Platform Overview</div>
        <div class="chart-body">
          <VChart :option="platformBarOption" style="height: 160px;" autoresize />
        </div>
      </div>
      <div class="m3-card m3-card--elevated chart-card">
        <div class="chart-card-header">User Activity</div>
        <div class="chart-body" style="display:flex;flex-direction:column;align-items:center;justify-content:center;">
          <VChart :option="userRingOption" style="height: 160px; width: 100%;" autoresize />
          <div style="display:flex;gap:12px;font-size:0.78rem;margin-top:8px;">
            <span><span style="display:inline-block;border-radius:50%;width:10px;height:10px;background:#22c55e;margin-right:4px;"></span>Active</span>
            <span><span style="display:inline-block;border-radius:50%;width:10px;height:10px;background:#f59e0b;margin-right:4px;"></span>Inactive</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Growth Chart -->
    <div class="m3-card m3-card--elevated" style="margin-bottom:24px;">
      <div class="chart-card-header">Platform Growth (Last 30 Days)</div>
      <div class="chart-body">
        <div v-if="growthLoading" style="display:flex;justify-content:center;padding:24px;">
          <md-circular-progress indeterminate style="--md-circular-progress-size:32px" />
        </div>
        <VChart v-else :option="growthLineOption" style="height: 240px;" autoresize />
      </div>
    </div>

    <!-- Users Section -->
    <div class="m3-card m3-card--elevated">
      <!-- Users header -->
      <div class="section-header">
        <span class="section-title">Users</span>
        <div style="max-width:280px;width:100%;">
          <md-outlined-text-field
            :value="search"
            @input="search=($event.target as HTMLInputElement).value; onSearchInput()"
            label="Search by email or name"
            style="width:100%;"
          >
            <span slot="leading-icon" class="material-symbols-outlined">search</span>
          </md-outlined-text-field>
        </div>
      </div>

      <div v-if="usersLoading" style="display:flex;justify-content:center;padding:40px;">
        <md-circular-progress indeterminate style="--md-circular-progress-size:32px" />
      </div>

      <div v-else-if="users.length === 0" class="users-empty-state">
        <span class="material-symbols-outlined" style="font-size:2rem;color:var(--md-sys-color-on-surface-variant);opacity:0.5;margin-bottom:8px;">manage_accounts</span>
        <span class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);">No users found.</span>
      </div>

      <div v-else class="m3-table-wrapper">
        <table class="m3-table">
          <thead>
            <tr>
              <th>User</th>
              <th>Provider</th>
              <th>Role</th>
              <th>Status</th>
              <th>Joined</th>
              <th>Last Login</th>
              <th style="text-align:right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td>
                <div style="display:flex;align-items:center;gap:10px;">
                  <div class="user-avatar">{{ (displayName(u)).charAt(0).toUpperCase() }}</div>
                  <div style="min-width:0;">
                    <div style="font-weight:500;font-size:0.875rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ displayName(u) }}</div>
                    <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.75rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ u.email }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span class="m3-badge m3-badge--neutral" style="text-transform:capitalize;">{{ u.auth_provider }}</span>
              </td>
              <td>
                <span :class="['m3-badge', u.role === 'admin' ? 'm3-badge--primary' : 'm3-badge--neutral']">{{ u.role }}</span>
              </td>
              <td>
                <span :class="['m3-badge', statusBadge(u.status)]">{{ u.status }}</span>
              </td>
              <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;white-space:nowrap;">{{ formatDate(u.created_at) }}</td>
              <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;white-space:nowrap;">{{ u.last_login ? formatDate(u.last_login) : '—' }}</td>
              <td>
                <div style="display:flex;gap:8px;justify-content:flex-end;align-items:center;flex-wrap:wrap;">
                  <AppSelect
                    style="min-width:90px;"
                    :model-value="u.role"
                    :disabled="changingRole === u.id || u.id === authStore.profile?.id"
                    @update:model-value="setRole(u, $event)"
                  >
                    <option value="user">User</option>
                    <option value="admin">Admin</option>
                  </AppSelect>
                  <md-icon-button
                    :disabled="impersonating === u.id || u.id === authStore.profile?.id"
                    @click="startImpersonation(u)"
                    title="Impersonate user"
                  >
                    <span v-if="impersonating === u.id">
                      <md-circular-progress indeterminate style="--md-circular-progress-size:20px" />
                    </span>
                    <span v-else class="material-symbols-outlined">manage_accounts</span>
                  </md-icon-button>
                  <md-outlined-button
                    v-if="u.status !== 'active'"
                    :disabled="updatingId === u.id"
                    @click="setStatus(u, 'active')"
                    style="--md-outlined-button-label-text-color:#16a34a;--md-outlined-button-outline-color:#16a34a;"
                  >
                    <md-circular-progress v-if="updatingId === u.id" indeterminate style="--md-circular-progress-size:16px" />
                    <span v-else>Activate</span>
                  </md-outlined-button>
                  <md-outlined-button
                    v-if="u.status !== 'banned'"
                    :disabled="updatingId === u.id"
                    @click="setStatus(u, 'banned')"
                    style="--md-outlined-button-label-text-color:var(--md-sys-color-error);--md-outlined-button-outline-color:var(--md-sys-color-error);"
                  >
                    <md-circular-progress v-if="updatingId === u.id" indeterminate style="--md-circular-progress-size:16px" />
                    <span v-else>Ban</span>
                  </md-outlined-button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination-bar">
        <span style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;">
          Page {{ currentPage }} of {{ totalPages }} &mdash; {{ totalUsers }} users
        </span>
        <div style="display:flex;gap:8px;">
          <md-outlined-button :disabled="currentPage <= 1" @click="goToPage(currentPage - 1)">
            <span class="material-symbols-outlined" slot="icon">chevron_left</span>
            Previous
          </md-outlined-button>
          <md-outlined-button :disabled="currentPage >= totalPages" @click="goToPage(currentPage + 1)">
            Next
            <span class="material-symbols-outlined" slot="trailing-icon">chevron_right</span>
          </md-outlined-button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import AppSelect from '@/components/AppSelect.vue';
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
    case 'active':   return 'm3-badge--success';
    case 'inactive': return 'm3-badge--warning';
    case 'banned':   return 'm3-badge--error';
    default:         return 'm3-badge--neutral';
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

<style scoped lang="scss">
.page-wrapper {
  padding: 24px;
  max-width: 1100px;
}

.page-header {
  display: flex;
  align-items: center;
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

/* Stats */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.stat-card {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-icon {
  font-size: 1.5rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  line-height: 1;
}

.stat-label {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 500;
}

/* Charts */
.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  margin-bottom: 24px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  overflow: hidden;
}

.chart-card-header {
  padding: 14px 20px;
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--md-sys-color-on-surface);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart-body {
  padding: 16px 20px;
}

/* Users section */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  flex-wrap: wrap;
  gap: 12px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--md-sys-color-on-surface);
}

.users-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  text-align: center;
}

.m3-table-wrapper {
  overflow-x: auto;
}

.m3-table {
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

/* M3 Cards */
.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
    margin-bottom: 24px;
  }
}

/* Badges */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  text-transform: capitalize;

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

/* Quick stats bar */
.quick-stats-bar {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.quick-stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 20px;
  flex: 1;
  min-width: 120px;
}

.quick-stat-divider {
  width: 1px;
  height: 32px;
  background: var(--md-sys-color-outline-variant);
  flex-shrink: 0;
}

/* User avatar */
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
  flex-shrink: 0;
}

/* Impersonation banner */
.impersonation-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  margin-bottom: 20px;
  border-radius: 12px;
  background: rgba(245, 158, 11, 0.12);
  border: 1px solid rgba(217, 119, 6, 0.35);
  flex-wrap: wrap;
}

/* Pagination */
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  flex-wrap: wrap;
  gap: 8px;
}
</style>
