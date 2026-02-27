<template>
  <div class="overview-page">

    <!-- Page header -->
    <div class="page-header">
      <h1 class="md-headline-small page-title">Overview</h1>
      <p class="md-body-medium page-subtitle">A snapshot of your account activity.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <md-circular-progress indeterminate style="--md-circular-progress-size: 48px" />
    </div>

    <template v-else-if="data">

      <!-- ── Stat cards ──────────────────────────────────────────────────── -->
      <div class="stat-grid">

        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-icon-wrap" style="background: var(--md-sys-color-primary-container)">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-on-primary-container)">link</span>
          </div>
          <div class="stat-value md-display-small">{{ data.total_links.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Total Links</div>
        </div>

        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-icon-wrap" style="background: var(--md-sys-color-secondary-container)">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-on-secondary-container)">ads_click</span>
          </div>
          <div class="stat-value md-display-small">{{ data.total_clicks.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Total Clicks</div>
        </div>

        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-icon-wrap" style="background: var(--md-sys-color-tertiary-container, color-mix(in srgb, var(--md-sys-color-primary) 15%, transparent))">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-on-tertiary-container, var(--md-sys-color-primary))">today</span>
          </div>
          <div class="stat-value md-display-small">{{ data.clicks_today.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Clicks Today</div>
        </div>

        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-icon-wrap" style="background: color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent)">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-primary)">calendar_month</span>
          </div>
          <div class="stat-value md-display-small">{{ data.clicks_30_days.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Clicks (30 days)</div>
        </div>

        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-icon-wrap" style="background: color-mix(in srgb, var(--md-sys-color-secondary) 12%, transparent)">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-secondary, var(--md-sys-color-primary))">bolt</span>
          </div>
          <div class="stat-value md-display-small">{{ data.clicks_7_days.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Clicks (7 days)</div>
        </div>

      </div>

      <!-- ── 30-Day Click Trend ───────────────────────────────────────────── -->
      <div class="m3-card m3-card--elevated chart-card">
        <div class="card-header-row">
          <span class="md-title-medium card-section-title">30-Day Click Trend</span>
          <span class="md-body-small" style="color: var(--md-sys-color-on-surface-variant)">Daily clicks across all your links</span>
        </div>
        <md-divider />
        <div class="chart-body">
          <div
            v-if="!data.time_series_30d || data.time_series_30d.length === 0"
            class="empty-chart"
          >
            <span class="material-symbols-outlined empty-icon">show_chart</span>
            <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant)">
              No click data yet. Clicks will appear here as your links are visited.
            </p>
          </div>
          <VChart v-else :option="trendChartOption" style="height: 220px;" autoresize />
        </div>
      </div>

      <!-- ── Expiring & At-Limit Links ────────────────────────────────────── -->
      <div
        v-if="expiringLinks.length > 0 || atLimitLinks.length > 0"
        class="m3-card m3-card--outlined attention-card"
      >
        <div class="card-header-row">
          <div style="display: flex; align-items: center; gap: 8px">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-error); font-size: 20px">warning</span>
            <span class="md-title-medium card-section-title">Links Needing Attention</span>
            <span class="m3-badge m3-badge--warning">{{ expiringLinks.length + atLimitLinks.length }}</span>
          </div>
          <router-link to="/dashboard/links?expiring_soon=true" class="view-all-link md-label-large">
            View all
          </router-link>
        </div>
        <md-divider />

        <!-- Expiring Soon -->
        <div v-if="expiringLinks.length > 0" class="attention-section">
          <div class="attention-section-title md-label-large">Expiring Soon</div>
          <table class="m3-table">
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
                <td class="text-muted md-body-small">{{ link.short_url }}</td>
                <td class="text-end text-muted md-body-small text-nowrap">
                  {{ formatDate(link.expires_at!) }}
                </td>
                <td class="text-end text-nowrap">
                  <span
                    class="m3-badge"
                    :class="daysUntil(link.expires_at!) <= 3 ? 'm3-badge--error' : 'm3-badge--warning'"
                  >
                    {{ daysUntil(link.expires_at!) }} days left
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <md-divider v-if="expiringLinks.length > 0 && atLimitLinks.length > 0" />

        <!-- Approaching Click Limit -->
        <div v-if="atLimitLinks.length > 0" class="attention-section">
          <div class="attention-section-title md-label-large">Approaching Click Limit</div>
          <table class="m3-table">
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
                <td class="text-end text-muted md-body-small text-nowrap">
                  {{ link.click_count.toLocaleString() }} / {{ link.max_clicks!.toLocaleString() }}
                </td>
                <td style="min-width: 140px">
                  <div class="progress-row">
                    <md-linear-progress
                      :value="Math.min(1, link.click_count / link.max_clicks!)"
                      style="flex: 1; --md-linear-progress-active-indicator-color: var(--md-sys-color-error)"
                    />
                    <span class="md-body-small progress-pct">
                      {{ Math.min(100, Math.round((link.click_count / link.max_clicks!) * 100)) }}%
                    </span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ── Two-column tables ───────────────────────────────────────────── -->
      <div class="tables-grid">

        <!-- Top links by clicks -->
        <div class="m3-card m3-card--elevated">
          <div class="card-header-row">
            <span class="md-title-medium card-section-title">Top Links by Clicks</span>
            <router-link to="/dashboard/links" class="view-all-link md-label-large">View all</router-link>
          </div>
          <md-divider />
          <div v-if="data.top_links.length === 0" class="empty-state">
            <span class="material-symbols-outlined empty-icon">link_off</span>
            <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant)">No clicks recorded yet.</p>
          </div>
          <table v-else class="m3-table">
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
                  <div class="link-dest md-body-small">{{ truncate(link.destination_url, 40) }}</div>
                </td>
                <td class="text-end">
                  <span class="m3-badge m3-badge--primary">{{ link.click_count.toLocaleString() }}</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Recently created links -->
        <div class="m3-card m3-card--elevated">
          <div class="card-header-row">
            <span class="md-title-medium card-section-title">Recently Created</span>
            <router-link to="/dashboard/links" class="view-all-link md-label-large">View all</router-link>
          </div>
          <md-divider />
          <div v-if="data.recent_links.length === 0" class="empty-state">
            <span class="material-symbols-outlined empty-icon">add_link</span>
            <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant)">No links created yet.</p>
          </div>
          <table v-else class="m3-table">
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
                  <div class="link-dest md-body-small">{{ truncate(link.destination_url, 40) }}</div>
                </td>
                <td class="text-end text-muted md-body-small text-nowrap">
                  {{ formatDate(link.created_at) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

      </div>

    </template>

    <!-- Error -->
    <div v-else-if="error" class="m3-card m3-card--outlined error-card">
      <span class="material-symbols-outlined" style="color: var(--md-sys-color-error)">error</span>
      <span class="md-body-medium" style="color: var(--md-sys-color-error)">{{ error }}</span>
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
.overview-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 24px;

  @media (max-width: 575px) {
    padding: 16px;
    gap: 16px;
  }
}

/* ── Page header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  margin: 0;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Loading ─────────────────────────────────────────────────────────────── */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 64px 0;
}

/* ── Stat grid ───────────────────────────────────────────────────────────── */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;

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
  flex-direction: column;
  padding: 20px;
  gap: 4px;
}

.stat-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.stat-value {
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  line-height: 1.1;
}

.stat-label {
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
}

/* ── M3 Cards ────────────────────────────────────────────────────────────── */
.m3-card {
  border-radius: 12px;
  overflow: hidden;
}

.m3-card--elevated {
  background: var(--md-sys-color-surface);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.08), 0 1px 4px rgba(0, 0, 0, 0.06);
}

.m3-card--outlined {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
}

/* ── Attention card ──────────────────────────────────────────────────────── */
.attention-card {
  border-color: var(--md-sys-color-error) !important;
  border-left-width: 3px !important;
}

/* ── Card header row ─────────────────────────────────────────────────────── */
.card-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  gap: 8px;
}

.card-section-title {
  color: var(--md-sys-color-on-surface);
  font-weight: 600;
}

.view-all-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  white-space: nowrap;

  &:hover {
    text-decoration: underline;
  }
}

/* ── Chart card ──────────────────────────────────────────────────────────── */
.chart-body {
  padding: 16px 16px 20px;
}

.empty-chart {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 32px 0;
  text-align: center;
}

/* ── Tables ──────────────────────────────────────────────────────────────── */
.tables-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;

  @media (max-width: 767px) {
    grid-template-columns: 1fr;
  }
}

.m3-table {
  width: 100%;
  border-collapse: collapse;

  th {
    font-size: 11px;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    padding: 10px 20px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    white-space: nowrap;
  }

  td {
    padding: 12px 20px;
    border-bottom: 1px solid color-mix(in srgb, var(--md-sys-color-outline-variant) 50%, transparent);
    vertical-align: middle;
    color: var(--md-sys-color-on-surface);
  }

  tr:last-child td {
    border-bottom: none;
  }
}

.text-end {
  text-align: right;
}

.text-nowrap {
  white-space: nowrap;
}

.text-muted {
  color: var(--md-sys-color-on-surface-variant);
}

.link-slug {
  font-size: 14px;
  font-weight: 600;
  color: var(--md-sys-color-primary);
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

.link-dest {
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 220px;
  margin-top: 2px;
}

/* ── M3 Badges ───────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 999px;
  white-space: nowrap;
}

.m3-badge--primary {
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);
}

.m3-badge--error {
  background: var(--md-sys-color-error);
  color: #fff;
}

.m3-badge--warning {
  background: #f59e0b;
  color: #fff;
}

/* ── Attention sections ──────────────────────────────────────────────────── */
.attention-section {
  padding: 0 0 4px;
}

.attention-section-title {
  padding: 12px 20px 4px;
  color: var(--md-sys-color-on-surface-variant);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-size: 11px;
}

/* ── Progress row ────────────────────────────────────────────────────────── */
.progress-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-pct {
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  min-width: 32px;
  text-align: right;
}

/* ── Error card ──────────────────────────────────────────────────────────── */
.error-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
}

/* ── Empty state ─────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 32px 16px;
  text-align: center;
}

.empty-icon {
  font-size: 40px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.6;
}
</style>
