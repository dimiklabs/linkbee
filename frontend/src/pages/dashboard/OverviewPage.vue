<template>
  <div class="overview-page">

    <!-- ── Page Header ──────────────────────────────────────────────────── -->
    <div class="page-header">
      <div class="page-header-main">
        <h1 class="page-title">Good {{ timeOfDay }}, {{ authStore.userName }}!</h1>
        <p class="page-subtitle">{{ todayDisplay }}</p>
      </div>
      <router-link to="/dashboard/links?create=1" class="header-cta">
        <span class="material-symbols-outlined">add</span>
        New Link
      </router-link>
    </div>

    <!-- ── Loading Skeleton ─────────────────────────────────────────────── -->
    <div v-if="loading" class="stat-grid" aria-busy="true" aria-label="Loading statistics">
      <div v-for="i in 5" :key="i" class="stat-card stat-card--skeleton">
        <div class="skeleton skel-icon"></div>
        <div class="skel-body">
          <div class="skeleton skel-value"></div>
          <div class="skeleton skel-label"></div>
        </div>
      </div>
    </div>

    <template v-else-if="data">

      <!-- ── Stat Cards ────────────────────────────────────────────────── -->
      <div class="stat-grid">

        <div class="stat-card stat-card--purple" aria-label="Total links">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--primary">
              <span class="material-symbols-outlined">link</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ data.total_links.toLocaleString() }}</div>
              <div class="stat-label">Total Links</div>
            </div>
          </div>
        </div>

        <div class="stat-card stat-card--teal" aria-label="Total clicks">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--secondary">
              <span class="material-symbols-outlined">ads_click</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ data.total_clicks.toLocaleString() }}</div>
              <div class="stat-label">Total Clicks</div>
            </div>
          </div>
        </div>

        <div class="stat-card stat-card--amber" aria-label="Clicks today">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--tertiary">
              <span class="material-symbols-outlined">today</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ data.clicks_today.toLocaleString() }}</div>
              <div class="stat-label">Clicks Today</div>
            </div>
          </div>
        </div>

        <div class="stat-card stat-card--violet" aria-label="Clicks in the last 30 days">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--amber">
              <span class="material-symbols-outlined">calendar_month</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ data.clicks_30_days.toLocaleString() }}</div>
              <div class="stat-label">Last 30 Days</div>
            </div>
          </div>
        </div>

        <div class="stat-card stat-card--emerald" aria-label="Clicks in the last 7 days">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--teal">
              <span class="material-symbols-outlined">bolt</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ data.clicks_7_days.toLocaleString() }}</div>
              <div class="stat-label">Last 7 Days</div>
            </div>
          </div>
        </div>

      </div>

      <!-- ── 30-Day Click Trend ─────────────────────────────────────────── -->
      <div class="overview-card">
        <div class="card-header">
          <div class="card-header-left">
            <div class="card-header-icon">
              <span class="material-symbols-outlined">show_chart</span>
            </div>
            <div>
              <div class="card-title">30-Day Click Trend</div>
              <div class="card-subtitle">Daily clicks across all your links</div>
            </div>
          </div>
          <router-link to="/dashboard/analytics" class="card-action-link">View analytics</router-link>
        </div>

        <div class="chart-body">
          <div v-if="!data.time_series_30d || data.time_series_30d.length === 0" class="empty-state">
            <div class="empty-icon-wrap">
              <span class="material-symbols-outlined empty-icon">show_chart</span>
            </div>
            <p class="empty-title">No click data yet</p>
            <p class="empty-sub">Clicks will appear here as your links are visited.</p>
            <router-link to="/dashboard/links">
              <button class="btn-filled">
                <span class="material-symbols-outlined">add_link</span>
                Create a link
              </button>
            </router-link>
          </div>
          <VChart v-else :option="trendChartOption" style="height: 280px;" autoresize />
        </div>
      </div>

      <!-- ── Links Needing Attention ───────────────────────────────────── -->
      <div
        v-if="expiringLinks.length > 0 || atLimitLinks.length > 0"
        class="overview-card overview-card--warning"
      >
        <div class="card-header">
          <div class="card-header-left">
            <div class="card-header-icon card-header-icon--warning">
              <span class="material-symbols-outlined">warning</span>
            </div>
            <div>
              <div class="card-title">Links Needing Attention</div>
              <div class="card-subtitle">
                {{ expiringLinks.length + atLimitLinks.length }} link{{ expiringLinks.length + atLimitLinks.length !== 1 ? 's' : '' }} require action
              </div>
            </div>
          </div>
          <router-link to="/dashboard/links" class="card-action-link">View all</router-link>
        </div>

        <!-- Expiring Soon -->
        <div v-if="expiringLinks.length > 0" class="attention-section">
          <div class="attention-section-label">Expiring Soon</div>
          <div class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th>Short URL</th>
                  <th class="col-end">Expires</th>
                  <th class="col-end">Time Left</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in expiringLinks" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="row-link">
                      /{{ link.slug }}
                    </router-link>
                  </td>
                  <td class="row-muted">{{ link.short_url }}</td>
                  <td class="col-end row-muted row-nowrap">{{ formatDate(link.expires_at!) }}</td>
                  <td class="col-end row-nowrap">
                    <span class="status-chip" :class="daysUntil(link.expires_at!) <= 3 ? 'status-chip--error' : 'status-chip--warn'">
                      {{ daysUntil(link.expires_at!) }}d left
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="expiringLinks.length > 0 && atLimitLinks.length > 0" class="section-divider" />

        <!-- Approaching Click Limit -->
        <div v-if="atLimitLinks.length > 0" class="attention-section">
          <div class="attention-section-label">Approaching Click Limit</div>
          <div class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th class="col-end">Clicks Used</th>
                  <th>Progress</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in atLimitLinks" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="row-link">
                      /{{ link.slug }}
                    </router-link>
                  </td>
                  <td class="col-end row-muted row-nowrap">
                    {{ link.click_count.toLocaleString() }} / {{ link.max_clicks!.toLocaleString() }}
                  </td>
                  <td>
                    <div class="progress-row">
                      <div class="native-progress">
                        <div
                          class="native-progress-fill"
                          :class="{ 'native-progress-fill--danger': link.click_count >= link.max_clicks! }"
                          :style="{ width: Math.min(100, Math.round((link.click_count / link.max_clicks!) * 100)) + '%' }"
                        />
                      </div>
                      <span class="progress-pct">{{ Math.min(100, Math.round((link.click_count / link.max_clicks!) * 100)) }}%</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- ── Two-column table grid ──────────────────────────────────────── -->
      <div class="tables-grid">

        <!-- Top links by clicks -->
        <div class="overview-card">
          <div class="card-header">
            <div class="card-header-left">
              <div class="card-header-icon">
                <span class="material-symbols-outlined">leaderboard</span>
              </div>
              <div class="card-title">Top Links by Clicks</div>
            </div>
            <router-link to="/dashboard/links" class="card-action-link">View all</router-link>
          </div>

          <div v-if="data.top_links.length === 0" class="empty-state empty-state--compact">
            <div class="empty-icon-wrap">
              <span class="material-symbols-outlined empty-icon">link_off</span>
            </div>
            <p class="empty-sub">No clicks recorded yet.</p>
            <router-link to="/dashboard/links">
              <button class="btn-outlined">Create a link</button>
            </router-link>
          </div>
          <div v-else class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th class="col-end">Clicks</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in data.top_links" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="row-link">
                      /{{ link.slug }}
                    </router-link>
                    <div class="row-dest">{{ truncate(link.destination_url, 40) }}</div>
                  </td>
                  <td class="col-end">
                    <span class="clicks-badge">{{ link.click_count.toLocaleString() }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Recently created links -->
        <div class="overview-card">
          <div class="card-header">
            <div class="card-header-left">
              <div class="card-header-icon">
                <span class="material-symbols-outlined">history</span>
              </div>
              <div class="card-title">Recently Created</div>
            </div>
            <router-link to="/dashboard/links" class="card-action-link">View all</router-link>
          </div>

          <div v-if="data.recent_links.length === 0" class="empty-state empty-state--compact">
            <div class="empty-icon-wrap">
              <span class="material-symbols-outlined empty-icon">add_link</span>
            </div>
            <p class="empty-sub">No links created yet.</p>
            <router-link to="/dashboard/links">
              <button class="btn-filled">
                <span class="material-symbols-outlined">add_link</span>
                Create your first link
              </button>
            </router-link>
          </div>
          <div v-else class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th class="col-end">Created</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="link in data.recent_links" :key="link.id">
                  <td>
                    <router-link :to="`/dashboard/links/${link.id}`" class="row-link">
                      /{{ link.slug }}
                    </router-link>
                    <div class="row-dest">{{ truncate(link.destination_url, 40) }}</div>
                  </td>
                  <td class="col-end row-muted row-nowrap">{{ formatDate(link.created_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      </div>

    </template>

    <!-- ── Error ─────────────────────────────────────────────────────────── -->
    <div v-else-if="error" class="error-banner">
      <span class="material-symbols-outlined">error</span>
      <span>{{ error }}</span>
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
import { useUiStore } from '@/stores/ui';

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer]);

const authStore = useAuthStore();
const uiStore = useUiStore();

// Dark-mode-aware chart color helpers
const surfaceColor = computed(() => uiStore.darkMode ? '#1A1F36' : '#ffffff');
const onSurface = computed(() => uiStore.darkMode ? '#E2E0EA' : '#1a1f36');
const muted = computed(() => uiStore.darkMode ? '#CAC4D0' : '#697386');
const gridLine = computed(() => uiStore.darkMode ? '#2C283F' : '#f0f2f5');

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
    if (res.data) data.value = res.data;
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

  const bg = surfaceColor.value;
  const textColor = onSurface.value;
  const mutedColor = muted.value;
  const gridColor = gridLine.value;
  const primaryColor = uiStore.darkMode ? '#C5C1FF' : '#635BFF';

  return {
    backgroundColor: bg,
    textStyle: { color: textColor },
    tooltip: {
      trigger: 'axis',
      formatter: (params: { name: string; value: number }[]) => {
        const p = params[0];
        if (!p) return '';
        return `${p.name}<br/><strong>${p.value.toLocaleString()} clicks</strong>`;
      },
      backgroundColor: uiStore.darkMode ? '#2C283F' : '#ffffff',
      borderColor: uiStore.darkMode ? '#49454F' : '#CAC4D0',
      textStyle: { color: textColor, fontSize: 12 },
    },
    grid: { top: 16, right: 16, bottom: 36, left: 52 },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: gridColor } },
      axisTick: { show: false },
      axisLabel: { color: mutedColor, fontSize: 11, interval: Math.max(0, Math.floor(dates.length / 6) - 1) },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: gridColor, type: 'dashed' } },
      axisLabel: { color: mutedColor, fontSize: 11 },
      minInterval: 1,
    },
    series: [
      {
        type: 'line',
        data: counts,
        smooth: true,
        symbol: 'circle',
        symbolSize: 4,
        lineStyle: { color: primaryColor, width: 2.5 },
        itemStyle: { color: primaryColor },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: uiStore.darkMode ? 'rgba(197,193,255,0.20)' : 'rgba(99,91,255,0.20)' },
              { offset: 1, color: uiStore.darkMode ? 'rgba(197,193,255,0.02)' : 'rgba(99,91,255,0.02)' },
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
  const overviewLinks: LinkResponse[] = [
    ...(data.value?.top_links ?? []),
    ...(data.value?.recent_links ?? []),
  ];
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
$ease: cubic-bezier(0.2, 0, 0, 1);

/* ── Page ────────────────────────────────────────────────────────────────── */
.overview-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  // No padding — DashboardLayout .page-content already provides 24px
}

/* ── Page Header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  padding-bottom: 8px;
}

.page-header-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  letter-spacing: -0.025em;
  line-height: 1.2;
}

.page-subtitle {
  margin: 0;
  font-size: 0.825rem;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 400;
}

.header-cta {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 38px;
  padding: 0 18px;
  border-radius: 9px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-size: 0.875rem;
  font-weight: 600;
  text-decoration: none;
  white-space: nowrap;
  flex-shrink: 0;
  transition: opacity 0.15s $ease, box-shadow 0.15s $ease;
  box-shadow: 0 1px 3px rgba(99,91,255,0.30), 0 1px 2px rgba(0,0,0,0.12);

  .material-symbols-outlined { font-size: 18px; }
  &:hover {
    opacity: 0.92;
    box-shadow: 0 3px 8px rgba(99,91,255,0.35), 0 1px 3px rgba(0,0,0,0.15);
    text-decoration: none;
  }
}

/* ── Stat Grid ───────────────────────────────────────────────────────────── */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;

  @media (max-width: 1199px) { grid-template-columns: repeat(3, 1fr); }
  @media (max-width: 767px)  { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 400px)  { grid-template-columns: 1fr; }
}

/* ── Stat Card ───────────────────────────────────────────────────────────── */
.stat-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
  position: relative;
  transition: box-shadow 0.18s $ease, border-color 0.18s $ease, transform 0.18s $ease;

  &::before {
    content: '';
    position: absolute;
    top: 0; left: 0; right: 0;
    height: 3px;
    background: var(--md-sys-color-primary);
  }

  &--purple::before  { background: #635BFF; }
  &--teal::before    { background: #0EA5A0; }
  &--amber::before   { background: #F59E0B; }
  &--violet::before  { background: #8B5CF6; }
  &--emerald::before { background: #10B981; }

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.10), 0 2px 4px rgba(0,0,0,0.06);
    border-color: var(--md-sys-color-outline);
    transform: translateY(-2px);
  }

  &--skeleton {
    pointer-events: none;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.stat-card-inner {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
}

.stat-icon-wrap {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 22px; }

  &--primary {
    background: color-mix(in srgb, #635BFF 12%, transparent);
    .material-symbols-outlined { color: #635BFF; }
  }
  &--secondary {
    background: color-mix(in srgb, #0EA5A0 12%, transparent);
    .material-symbols-outlined { color: #0EA5A0; }
  }
  &--tertiary {
    background: color-mix(in srgb, #F59E0B 12%, transparent);
    .material-symbols-outlined { color: #D97706; }
  }
  &--amber {
    background: color-mix(in srgb, #8B5CF6 12%, transparent);
    .material-symbols-outlined { color: #7C3AED; }
  }
  &--teal {
    background: color-mix(in srgb, #10B981 12%, transparent);
    .material-symbols-outlined { color: #059669; }
  }
}

.stat-body {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  line-height: 1.1;
  letter-spacing: -0.03em;
}

.stat-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.01em;
}

/* ── Skeleton ────────────────────────────────────────────────────────────── */
.skel-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 7px;
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

.skel-icon  { width: 44px; height: 44px; border-radius: 12px; flex-shrink: 0; }
.skel-value { height: 1.6rem; width: 60%; border-radius: 4px; }
.skel-label { height: 0.75rem; width: 80%; border-radius: 4px; }

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ── Overview Card ───────────────────────────────────────────────────────── */
.overview-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06), 0 1px 2px rgba(0,0,0,0.04);
  transition: box-shadow 0.18s $ease;

  &:hover {
    box-shadow: 0 2px 6px rgba(0,0,0,0.08), 0 1px 3px rgba(0,0,0,0.05);
  }

  &--warning {
    border-color: color-mix(in srgb, #f59e0b 40%, var(--md-sys-color-outline-variant));
    background: color-mix(in srgb, #FEF3C7 10%, var(--md-sys-color-surface));
  }
}

/* ── Card Header ─────────────────────────────────────────────────────────── */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  gap: 8px;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.card-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.card-header-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: var(--md-sys-color-surface-container);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 18px; color: var(--md-sys-color-primary); }

  &--warning {
    background: color-mix(in srgb, #f59e0b 12%, transparent);
    .material-symbols-outlined { color: #b45309; }
  }
}

.card-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0;
}

.card-subtitle {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 2px 0 0;
}

.card-action-link {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--md-sys-color-primary);
  text-decoration: none;
  white-space: nowrap;
  flex-shrink: 0;
  &:hover { text-decoration: underline; }
}

/* ── Chart Body ──────────────────────────────────────────────────────────── */
.chart-body {
  padding: 16px 16px 20px;
}

/* ── Tables grid ─────────────────────────────────────────────────────────── */
.tables-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;

  @media (max-width: 767px) { grid-template-columns: 1fr; }
}

/* ── Data Table ──────────────────────────────────────────────────────────── */
.table-wrap {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.data-table {
  width: 100%;
  border-collapse: collapse;

  th {
    font-size: 0.68rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    text-transform: uppercase;
    letter-spacing: 0.07em;
    padding: 8px 20px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container-low);
    white-space: nowrap;
  }

  td {
    padding: 11px 20px;
    border-bottom: 1px solid color-mix(in srgb, var(--md-sys-color-outline-variant) 50%, transparent);
    vertical-align: middle;
    color: var(--md-sys-color-on-surface);
    font-size: 0.875rem;
  }

  tr:last-child td { border-bottom: none; }

  tbody tr {
    transition: background 0.12s;
    &:hover td { background: var(--md-sys-color-surface-container-low); }
  }
}

.col-end { text-align: right; }
.row-muted { color: var(--md-sys-color-on-surface-variant); font-size: 0.8rem; }
.row-nowrap { white-space: nowrap; }

.row-link {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-primary);
  text-decoration: none;
  &:hover { text-decoration: underline; }
}

.row-dest {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 220px;
  margin-top: 2px;
}

/* ── Status Chips ────────────────────────────────────────────────────────── */
.status-chip {
  display: inline-flex;
  align-items: center;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  white-space: nowrap;

  &--error {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }
  &--warn {
    background: color-mix(in srgb, #f59e0b 15%, transparent);
    color: #92400e;
  }
}

.clicks-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.78rem;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 6px;
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);
  white-space: nowrap;
}

/* ── Attention Section ───────────────────────────────────────────────────── */
.attention-section {
  padding: 0 0 4px;
}

.attention-section-label {
  padding: 10px 20px 6px;
  font-size: 0.68rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  text-transform: uppercase;
  letter-spacing: 0.07em;
}

.section-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  margin: 4px 0;
}

/* ── Native Progress Bar ─────────────────────────────────────────────────── */
.progress-row {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 140px;
}

.native-progress {
  flex: 1;
  height: 6px;
  background: var(--md-sys-color-surface-container-high);
  border-radius: 999px;
  overflow: hidden;
}

.native-progress-fill {
  height: 100%;
  border-radius: 999px;
  background: var(--md-sys-color-primary);
  transition: width 0.4s $ease;

  &--danger { background: var(--md-sys-color-error); }
}

.progress-pct {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  min-width: 32px;
  text-align: right;
}

/* ── Empty State ─────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 24px;
  text-align: center;
  gap: 8px;

  &--compact { padding: 32px 16px; }
}

.empty-icon-wrap {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.empty-icon {
  font-size: 36px;
  color: var(--md-sys-color-primary);
  opacity: 0.55;
}

.empty-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0;
}

.empty-sub {
  font-size: 0.825rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 8px;
  max-width: 300px;
}

/* ── Error Banner ────────────────────────────────────────────────────────── */
.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 18px;
  background: var(--md-sys-color-error-container);
  border-radius: 12px;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-error-container);

  .material-symbols-outlined { color: var(--md-sys-color-error); font-size: 20px; flex-shrink: 0; }
}
</style>
