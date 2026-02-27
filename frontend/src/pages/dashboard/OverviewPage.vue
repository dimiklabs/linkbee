<template>
  <div class="overview-page">
    <div class="page-header mb-4">
      <h1 class="page-title">Overview</h1>
      <p class="page-subtitle">A snapshot of your account activity.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading…</span>
      </div>
    </div>

    <template v-else-if="data">

      <!-- ── Stat cards ──────────────────────────────────────────────────── -->
      <div class="stat-grid mb-4">
        <div class="stat-card card">
          <div class="stat-icon">🔗</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.total_links.toLocaleString() }}</div>
            <div class="stat-label">Total Links</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">👆</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.total_clicks.toLocaleString() }}</div>
            <div class="stat-label">Total Clicks</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">☀️</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.clicks_today.toLocaleString() }}</div>
            <div class="stat-label">Clicks Today</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">📅</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.clicks_30_days.toLocaleString() }}</div>
            <div class="stat-label">Clicks (30 days)</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">⚡</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.clicks_7_days.toLocaleString() }}</div>
            <div class="stat-label">Clicks (7 days)</div>
          </div>
        </div>
      </div>

      <!-- ── 30-Day Click Trend ───────────────────────────────────────────── -->
      <div class="card mb-4">
        <div class="card-header-row">
          <h2 class="section-title">30-Day Click Trend</h2>
          <span class="text-muted" style="font-size: 0.8125rem;">Daily clicks across all your links</span>
        </div>
        <div class="card-body px-3 pb-3 pt-2">
          <div v-if="!data.time_series_30d || data.time_series_30d.length === 0" class="text-center py-4 text-muted small">
            No click data yet. Clicks will appear here as your links are visited.
          </div>
          <VChart v-else :option="trendChartOption" style="height: 220px;" autoresize />
        </div>
      </div>

      <!-- ── Expiring & At-Limit Links ────────────────────────────────────── -->
      <div v-if="expiringLinks.length > 0 || atLimitLinks.length > 0" class="card attention-card mb-4">
        <div class="card-header-row">
          <h2 class="section-title">
            ⚠️ Links Needing Attention
            <span class="attention-badge">{{ expiringLinks.length + atLimitLinks.length }}</span>
          </h2>
          <router-link to="/dashboard/links?expiring_soon=true" class="view-all-link">View all</router-link>
        </div>
        <div class="card-body p-0">

          <!-- Expiring Soon -->
          <div v-if="expiringLinks.length > 0" class="attention-section">
            <div class="attention-section-title">Expiring Soon</div>
            <table class="link-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th>Short URL</th>
                  <th class="text-end">Expires</th>
                  <th class="text-end">Time Left</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in expiringLinks" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="link-slug">
                      /{{ link.slug }}
                    </router-link>
                  </td>
                  <td class="text-muted small">{{ link.short_url }}</td>
                  <td class="text-end text-muted small text-nowrap">
                    {{ formatDate(link.expires_at!) }}
                  </td>
                  <td class="text-end text-nowrap">
                    <span :class="daysUntil(link.expires_at!) <= 3 ? 'expiry-urgent' : 'expiry-soon'">
                      {{ daysUntil(link.expires_at!) }} days left
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Approaching Click Limit -->
          <div v-if="atLimitLinks.length > 0" class="attention-section">
            <div class="attention-section-title">Approaching Click Limit</div>
            <table class="link-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th class="text-end">Clicks Used</th>
                  <th>Progress</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in atLimitLinks" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="link-slug">
                      /{{ link.slug }}
                    </router-link>
                  </td>
                  <td class="text-end text-muted small text-nowrap">
                    {{ link.click_count.toLocaleString() }} / {{ link.max_clicks!.toLocaleString() }}
                  </td>
                  <td style="min-width: 120px;">
                    <div class="click-progress-track">
                      <div
                        class="click-progress-bar"
                        :style="{ width: Math.min(100, Math.round((link.click_count / link.max_clicks!) * 100)) + '%' }"
                      ></div>
                    </div>
                    <span class="click-progress-pct">
                      {{ Math.min(100, Math.round((link.click_count / link.max_clicks!) * 100)) }}%
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

        </div>
      </div>

      <!-- ── Two-column tables ───────────────────────────────────────────── -->
      <div class="tables-grid">

        <!-- Top links by clicks -->
        <div class="card">
          <div class="card-header-row">
            <h2 class="section-title">Top Links by Clicks</h2>
            <router-link to="/dashboard/links" class="view-all-link">View all</router-link>
          </div>
          <div class="card-body p-0">
            <div v-if="data.top_links.length === 0" class="empty-state py-4 text-center text-muted small">
              No clicks recorded yet.
            </div>
            <table v-else class="link-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th class="text-end">Clicks</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in data.top_links" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="link-slug">
                      /{{ link.slug }}
                    </router-link>
                    <div class="link-dest text-muted">{{ truncate(link.destination_url, 40) }}</div>
                  </td>
                  <td class="text-end">
                    <span class="click-badge">{{ link.click_count.toLocaleString() }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Recently created links -->
        <div class="card">
          <div class="card-header-row">
            <h2 class="section-title">Recently Created</h2>
            <router-link to="/dashboard/links" class="view-all-link">View all</router-link>
          </div>
          <div class="card-body p-0">
            <div v-if="data.recent_links.length === 0" class="empty-state py-4 text-center text-muted small">
              No links created yet.
            </div>
            <table v-else class="link-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th class="text-end">Created</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in data.recent_links" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="link-slug">
                      /{{ link.slug }}
                    </router-link>
                    <div class="link-dest text-muted">{{ truncate(link.destination_url, 40) }}</div>
                  </td>
                  <td class="text-end text-muted small text-nowrap">
                    {{ formatDate(link.created_at) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      </div>

    </template>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-danger">
      {{ error }}
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { use } from 'echarts/core';
import { LineChart } from 'echarts/charts';
import { GridComponent, TooltipComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import VChart from 'vue-echarts';
import dashboardApi from '@/api/dashboard';
import linksApi from '@/api/links';
import type { DashboardOverviewResponse } from '@/types/dashboard';
import type { LinkResponse } from '@/types/links';

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer]);

const loading = ref(true);
const data = ref<DashboardOverviewResponse | null>(null);
const error = ref('');

const expiringLinks = ref<LinkResponse[]>([]);
const atLimitLinks = ref<LinkResponse[]>([]);

async function load() {
  loading.value = true;
  error.value = '';
  try {
    const res = await dashboardApi.getOverview();
    data.value = res.data;
  } catch {
    error.value = 'Failed to load dashboard overview.';
  } finally {
    loading.value = false;
  }
}

async function loadExpiringLinks() {
  try {
    const res = await linksApi.list(1, 20, '', '', undefined, undefined, undefined, true);
    expiringLinks.value = res.data?.links ?? [];
  } catch { /* non-fatal */ }
}

function computeAtLimitLinks(links: LinkResponse[]): LinkResponse[] {
  return links.filter(
    (l) => l.max_clicks != null && l.click_count >= l.max_clicks * 0.9,
  );
}

function daysUntil(iso: string): number {
  return Math.ceil((new Date(iso).getTime() - Date.now()) / 86400000);
}

function truncate(str: string, max: number): string {
  return str.length > max ? str.slice(0, max) + '…' : str;
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

const trendChartOption = computed(() => {
  const series = data.value?.time_series_30d ?? [];
  const dates = series.map((p) => new Date(p.timestamp).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }));
  const counts = series.map((p) => Number(p.count));

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: { name: string; value: number }[]) => {
        const p = params[0];
        return `${p.name}<br/><strong>${p.value.toLocaleString()} clicks</strong>`;
      },
      backgroundColor: '#fff',
      borderColor: '#e3e8ee',
      textStyle: { color: '#1a1f36', fontSize: 12 },
    },
    grid: { top: 12, right: 16, bottom: 32, left: 48 },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: '#e3e8ee' } },
      axisTick: { show: false },
      axisLabel: { color: '#697386', fontSize: 11, interval: Math.max(0, Math.floor(dates.length / 6) - 1) },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: '#f0f2f5' } },
      axisLabel: { color: '#697386', fontSize: 11 },
      minInterval: 1,
    },
    series: [
      {
        type: 'line',
        data: counts,
        smooth: true,
        symbol: 'none',
        lineStyle: { color: '#635bff', width: 2 },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(99,91,255,0.18)' },
              { offset: 1, color: 'rgba(99,91,255,0.01)' },
            ],
          },
        },
      },
    ],
  };
});

onMounted(async () => {
  await Promise.all([
    load(),
    loadExpiringLinks(),
  ]);
  // Derive at-limit links from the recently-created / top links already in the overview
  const overviewLinks: LinkResponse[] = [
    ...(data.value?.top_links ?? []),
    ...(data.value?.recent_links ?? []),
  ];
  // Deduplicate by id before filtering
  const seen = new Set<string>();
  const unique = overviewLinks.filter((l) => {
    if (seen.has(l.id)) return false;
    seen.add(l.id);
    return true;
  });
  atLimitLinks.value = computeAtLimitLinks(unique);
});
</script>

<style scoped lang="scss">
$primary: #635bff;

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: #1a1f36;
  margin: 0;
}

.page-subtitle {
  font-size: 0.875rem;
  color: #697386;
  margin: 0.25rem 0 0;
}

// ── Stat grid ─────────────────────────────────────────────────────────────────
.stat-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1rem;

  @media (max-width: 1199px) {
    grid-template-columns: repeat(3, 1fr);
  }

  @media (max-width: 767px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 575px) {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;

  .stat-icon {
    font-size: 2rem;
    line-height: 1;
    flex-shrink: 0;
  }

  .stat-value {
    font-size: 1.75rem;
    font-weight: 700;
    color: #1a1f36;
    line-height: 1;
  }

  .stat-label {
    font-size: 0.8125rem;
    color: #697386;
    margin-top: 0.25rem;
  }
}

// ── Two-column layout ─────────────────────────────────────────────────────────
.tables-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;

  @media (max-width: 767px) {
    grid-template-columns: 1fr;
  }
}

.card-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #e3e8ee;
}

.section-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1a1f36;
  margin: 0;
}

.view-all-link {
  font-size: 0.8125rem;
  color: $primary;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

// ── Link table ────────────────────────────────────────────────────────────────
.link-table {
  width: 100%;
  border-collapse: collapse;

  th {
    font-size: 0.75rem;
    font-weight: 600;
    color: #697386;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    padding: 0.625rem 1.25rem;
    border-bottom: 1px solid #e3e8ee;
    white-space: nowrap;
  }

  td {
    padding: 0.75rem 1.25rem;
    border-bottom: 1px solid #f0f2f5;
    vertical-align: middle;

    &:last-child {
      border-bottom: none;
    }
  }

  tr:last-child td {
    border-bottom: none;
  }
}

.link-slug {
  font-size: 0.875rem;
  font-weight: 600;
  color: $primary;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

.link-dest {
  font-size: 0.75rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 220px;
  margin-top: 0.125rem;
}

.click-badge {
  display: inline-block;
  background: #eef0ff;
  color: $primary;
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
}

// ── Attention card ────────────────────────────────────────────────────────────
.attention-card {
  border-left: 3px solid #f59e0b !important;
}

.attention-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #f59e0b;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 700;
  border-radius: 999px;
  min-width: 1.25rem;
  height: 1.25rem;
  padding: 0 0.35rem;
  margin-left: 0.5rem;
  vertical-align: middle;
}

.attention-section {
  &:not(:last-child) {
    border-bottom: 1px solid #e3e8ee;
  }
}

.attention-section-title {
  font-size: 0.75rem;
  font-weight: 700;
  color: #697386;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0.625rem 1.25rem 0;
}

.expiry-urgent {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #ef4444;
}

.expiry-soon {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #f59e0b;
}

.click-progress-track {
  background: #e3e8ee;
  border-radius: 999px;
  height: 6px;
  width: 100%;
  overflow: hidden;
  display: inline-block;
  vertical-align: middle;
}

.click-progress-bar {
  background: #f59e0b;
  height: 100%;
  border-radius: 999px;
  transition: width 0.3s ease;
}

.click-progress-pct {
  font-size: 0.75rem;
  color: #697386;
  margin-left: 0.375rem;
  vertical-align: middle;
}

// ── Dark mode ─────────────────────────────────────────────────────────────────
:global(.dark-mode) .page-title { color: #e6edf3; }
:global(.dark-mode) .stat-card .stat-value { color: #e6edf3; }
:global(.dark-mode) .section-title { color: #e6edf3; }
:global(.dark-mode) .card-header-row { border-bottom-color: #30363d; }
:global(.dark-mode) .link-table th { color: #8b949e; border-bottom-color: #30363d; }
:global(.dark-mode) .link-table td { border-bottom-color: #21262d; }
:global(.dark-mode) .click-badge { background: rgba(99, 91, 255, 0.2); }
</style>
