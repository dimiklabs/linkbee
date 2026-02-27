<template>
  <div class="overview-page">

    <!-- Page header with greeting -->
    <div class="page-header">
      <div class="page-header-content">
        <h1 class="md-headline-small page-title">Good {{ timeOfDay }}, {{ authStore.userName }}!</h1>
        <p class="md-body-medium page-subtitle">{{ todayDisplay }} &mdash; Here's a snapshot of your account activity.</p>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="stat-grid" aria-busy="true" aria-label="Loading statistics">
      <div v-for="i in 5" :key="i" class="m3-card m3-card--elevated stat-card skeleton-card">
        <div class="skeleton skeleton-icon"></div>
        <div class="skeleton skeleton-value"></div>
        <div class="skeleton skeleton-label"></div>
      </div>
    </div>

    <template v-else-if="data">

      <!-- Stat cards -->
      <div class="stat-grid">

        <div
          class="m3-card m3-card--elevated stat-card"
          aria-label="Total links"
          role="region"
        >
          <div class="stat-icon-wrap stat-icon-wrap--primary">
            <span class="material-symbols-outlined">link</span>
          </div>
          <div class="stat-value md-display-small">{{ data.total_links.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Total Links</div>
        </div>

        <div
          class="m3-card m3-card--elevated stat-card"
          aria-label="Total clicks"
          role="region"
        >
          <div class="stat-icon-wrap stat-icon-wrap--secondary">
            <span class="material-symbols-outlined">ads_click</span>
          </div>
          <div class="stat-value md-display-small">{{ data.total_clicks.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Total Clicks</div>
        </div>

        <div
          class="m3-card m3-card--elevated stat-card"
          aria-label="Clicks today"
          role="region"
        >
          <div class="stat-icon-wrap stat-icon-wrap--tertiary">
            <span class="material-symbols-outlined">today</span>
          </div>
          <div class="stat-value md-display-small">{{ data.clicks_today.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Clicks Today</div>
        </div>

        <div
          class="m3-card m3-card--elevated stat-card"
          aria-label="Clicks in the last 30 days"
          role="region"
        >
          <div class="stat-icon-wrap stat-icon-wrap--amber">
            <span class="material-symbols-outlined">calendar_month</span>
          </div>
          <div class="stat-value md-display-small">{{ data.clicks_30_days.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Clicks (30 days)</div>
        </div>

        <div
          class="m3-card m3-card--elevated stat-card"
          aria-label="Clicks in the last 7 days"
          role="region"
        >
          <div class="stat-icon-wrap stat-icon-wrap--teal">
            <span class="material-symbols-outlined">bolt</span>
          </div>
          <div class="stat-value md-display-small">{{ data.clicks_7_days.toLocaleString() }}</div>
          <div class="md-body-medium stat-label">Clicks (7 days)</div>
        </div>

      </div>

      <!-- 30-Day Click Trend -->
      <div class="m3-card m3-card--elevated chart-card">
        <div class="m3-card-header">
          <div class="m3-card-header__left">
            <span class="material-symbols-outlined m3-card-header__icon">show_chart</span>
            <div>
              <div class="md-title-medium m3-card-header__title">30-Day Click Trend</div>
              <div class="md-body-small m3-card-header__subtitle">Daily clicks across all your links</div>
            </div>
          </div>
          <router-link to="/dashboard/analytics" class="view-all-link md-label-large">View analytics</router-link>
        </div>
        <md-divider />
        <div class="chart-body">
          <div
            v-if="!data.time_series_30d || data.time_series_30d.length === 0"
            class="m3-empty-state"
          >
            <span class="material-symbols-outlined m3-empty-state__icon">show_chart</span>
            <p class="md-title-small m3-empty-state__title">No click data yet</p>
            <p class="md-body-medium m3-empty-state__subtitle">
              Clicks will appear here as your links are visited.
            </p>
            <router-link to="/dashboard/links">
              <md-filled-button>Create a link</md-filled-button>
            </router-link>
          </div>
          <VChart v-else :option="trendChartOption" style="height: 280px;" autoresize />
        </div>
      </div>

      <!-- Expiring & At-Limit Links -->
      <div
        v-if="expiringLinks.length > 0 || atLimitLinks.length > 0"
        class="m3-card m3-card--outlined attention-card"
      >
        <div class="m3-card-header">
          <div class="m3-card-header__left">
            <span class="material-symbols-outlined" style="color: var(--md-sys-color-error); font-size: 20px">warning</span>
            <div>
              <div class="md-title-medium m3-card-header__title">Links Needing Attention</div>
            </div>
            <span class="m3-badge m3-badge--error">{{ expiringLinks.length + atLimitLinks.length }}</span>
          </div>
          <router-link to="/dashboard/links?expiring_soon=true" class="view-all-link md-label-large">
            View all
          </router-link>
        </div>
        <md-divider />

        <!-- Expiring Soon -->
        <div v-if="expiringLinks.length > 0" class="attention-section">
          <div class="attention-section-title md-label-large">Expiring Soon</div>
          <div class="m3-table-wrapper">
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
        </div>

        <md-divider v-if="expiringLinks.length > 0 && atLimitLinks.length > 0" />

        <!-- Approaching Click Limit -->
        <div v-if="atLimitLinks.length > 0" class="attention-section">
          <div class="attention-section-title md-label-large">Approaching Click Limit</div>
          <div class="m3-table-wrapper">
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
      </div>

      <!-- Two-column tables -->
      <div class="tables-grid">

        <!-- Top links by clicks -->
        <div class="m3-card m3-card--elevated">
          <div class="m3-card-header">
            <div class="m3-card-header__left">
              <span class="material-symbols-outlined m3-card-header__icon">leaderboard</span>
              <span class="md-title-medium m3-card-header__title">Top Links by Clicks</span>
            </div>
            <router-link to="/dashboard/links" class="view-all-link md-label-large">View all</router-link>
          </div>
          <md-divider />
          <div v-if="data.top_links.length === 0" class="m3-empty-state m3-empty-state--compact">
            <span class="material-symbols-outlined m3-empty-state__icon">link_off</span>
            <p class="md-body-medium m3-empty-state__subtitle">No clicks recorded yet.</p>
            <router-link to="/dashboard/links">
              <md-outlined-button>Create a link</md-outlined-button>
            </router-link>
          </div>
          <div v-else class="m3-table-wrapper">
            <table class="m3-table">
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
        </div>

        <!-- Recently created links -->
        <div class="m3-card m3-card--elevated">
          <div class="m3-card-header">
            <div class="m3-card-header__left">
              <span class="material-symbols-outlined m3-card-header__icon">history</span>
              <span class="md-title-medium m3-card-header__title">Recently Created</span>
            </div>
            <router-link to="/dashboard/links" class="view-all-link md-label-large">View all</router-link>
          </div>
          <md-divider />
          <div v-if="data.recent_links.length === 0" class="m3-empty-state m3-empty-state--compact">
            <span class="material-symbols-outlined m3-empty-state__icon">add_link</span>
            <p class="md-body-medium m3-empty-state__subtitle">No links created yet.</p>
            <router-link to="/dashboard/links">
              <md-filled-button>Create your first link</md-filled-button>
            </router-link>
          </div>
          <div v-else class="m3-table-wrapper">
            <table class="m3-table">
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
import { useAuthStore } from '@/stores/auth';

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer]);

const authStore = useAuthStore();

const loading = ref(true);
const data = ref<DashboardOverviewResponse | null>(null);
const error = ref('');

const expiringLinks = ref<LinkResponse[]>([]);
const atLimitLinks = ref<LinkResponse[]>([]);

// Greeting helpers
const timeOfDay = computed(() => {
  const h = new Date().getHours();
  if (h < 12) return 'morning';
  if (h < 17) return 'afternoon';
  return 'evening';
});

const todayDisplay = computed(() =>
  new Date().toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })
);

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
      backgroundColor: 'var(--md-sys-color-surface-container-high)',
      borderColor: 'var(--md-sys-color-outline-variant)',
      textStyle: { color: 'var(--md-sys-color-on-surface)', fontSize: 12 },
    },
    grid: { top: 16, right: 16, bottom: 36, left: 52 },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: 'var(--md-sys-color-outline-variant)' } },
      axisTick: { show: false },
      axisLabel: { color: 'var(--md-sys-color-on-surface-variant)', fontSize: 11, interval: Math.max(0, Math.floor(dates.length / 6) - 1) },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: 'var(--md-sys-color-outline-variant)', type: 'dashed' } },
      axisLabel: { color: 'var(--md-sys-color-on-surface-variant)', fontSize: 11 },
      minInterval: 1,
    },
    series: [
      {
        type: 'line',
        data: counts,
        smooth: true,
        symbol: 'circle',
        symbolSize: 4,
        lineStyle: { color: 'var(--md-sys-color-primary)', width: 2.5 },
        itemStyle: { color: 'var(--md-sys-color-primary)' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'color-mix(in srgb, var(--md-sys-color-primary) 20%, transparent)' },
              { offset: 1, color: 'color-mix(in srgb, var(--md-sys-color-primary) 2%, transparent)' },
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
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.page-header-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  margin: 0;
  color: var(--md-sys-color-on-surface);
  font-weight: 700;
}

.page-subtitle {
  margin: 0;
  color: var(--md-sys-color-on-surface-variant);
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

  @media (max-width: 400px) {
    grid-template-columns: 1fr;
  }
}

/* ── Stat card ───────────────────────────────────────────────────────────── */
.stat-card {
  display: flex;
  flex-direction: column;
  padding: 20px;
  gap: 6px;
  border-radius: 16px;
  transition: box-shadow 0.2s, transform 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }
}

.stat-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  flex-shrink: 0;

  .material-symbols-outlined {
    font-size: 22px;
  }

  &--primary {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);

    .material-symbols-outlined { color: var(--md-sys-color-on-primary-container); }
  }

  &--secondary {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);

    .material-symbols-outlined { color: var(--md-sys-color-on-secondary-container); }
  }

  &--tertiary {
    background: var(--md-sys-color-tertiary-container, color-mix(in srgb, var(--md-sys-color-primary) 15%, transparent));

    .material-symbols-outlined { color: var(--md-sys-color-on-tertiary-container, var(--md-sys-color-primary)); }
  }

  &--amber {
    background: color-mix(in srgb, #f59e0b 15%, transparent);

    .material-symbols-outlined { color: #b45309; }
  }

  &--teal {
    background: color-mix(in srgb, #14b8a6 15%, transparent);

    .material-symbols-outlined { color: #0f766e; }
  }
}

.stat-value {
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  line-height: 1.1;
}

.stat-label {
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
  font-size: 0.82rem;
}

/* ── Skeleton ────────────────────────────────────────────────────────────── */
.skeleton-card {
  pointer-events: none;
}

.skeleton {
  background: linear-gradient(
    90deg,
    var(--md-sys-color-surface-container) 25%,
    var(--md-sys-color-surface-container-high) 50%,
    var(--md-sys-color-surface-container) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
  border-radius: 6px;
  display: block;
}

.skeleton-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  margin-bottom: 8px;
}

.skeleton-value {
  width: 70%;
  height: 2rem;
  margin-bottom: 4px;
}

.skeleton-label {
  width: 50%;
  height: 0.9rem;
}

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ── M3 Cards ────────────────────────────────────────────────────────────── */
.m3-card {
  border-radius: 16px;
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
  border-left: 3px solid var(--md-sys-color-error) !important;
}

/* ── M3 Card Header ──────────────────────────────────────────────────────── */
.m3-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  gap: 8px;
  flex-wrap: wrap;
}

.m3-card-header__left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.m3-card-header__icon {
  font-size: 20px;
  color: var(--md-sys-color-primary);
  flex-shrink: 0;
}

.m3-card-header__title {
  color: var(--md-sys-color-on-surface);
  font-weight: 600;
  margin: 0;
}

.m3-card-header__subtitle {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.view-all-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  white-space: nowrap;
  font-size: 0.85rem;
  font-weight: 600;

  &:hover {
    text-decoration: underline;
  }
}

/* ── Chart card ──────────────────────────────────────────────────────────── */
.chart-card {
  overflow: visible;
}

.chart-body {
  padding: 16px 16px 20px;
}

/* ── Tables grid ─────────────────────────────────────────────────────────── */
.tables-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;

  @media (max-width: 767px) {
    grid-template-columns: 1fr;
  }
}

/* ── M3 Table Wrapper ────────────────────────────────────────────────────── */
.m3-table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
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

  tbody tr:hover td {
    background: var(--md-sys-color-surface-container-low);
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
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container, var(--md-sys-color-error));
}

.m3-badge--warning {
  background: color-mix(in srgb, #f59e0b 15%, transparent);
  color: #92400e;
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

/* ── M3 Empty State ──────────────────────────────────────────────────────── */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 48px 24px;
  text-align: center;

  &--compact {
    padding: 32px 16px;
  }
}

.m3-empty-state__icon {
  font-size: 48px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.5;
  margin-bottom: 4px;
}

.m3-empty-state__title {
  color: var(--md-sys-color-on-surface);
  margin: 0;
  font-weight: 600;
}

.m3-empty-state__subtitle {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 8px;
  max-width: 320px;
}
</style>
