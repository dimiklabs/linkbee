<template>
  <div class="link-comparison-page">

    <!-- Page Header -->
    <div class="page-header">
      <h1 class="md-title-large page-title">Compare Links</h1>
      <p class="page-subtitle">Select 2–5 links to compare their performance side by side.</p>
    </div>

    <!-- Filters Card -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-icon an-card-icon--primary">
          <span class="material-symbols-outlined">compare_arrows</span>
        </div>
        <span class="an-card-title">Configure Comparison</span>
      </div>
      <div class="an-card-body">
        <div class="filter-grid">
          <!-- Link selector -->
          <div class="link-selector-wrapper">
            <div class="form-field__label" style="margin-bottom:8px;">
              Links <span style="font-weight:400;color:var(--md-sys-color-on-surface-variant);">(select 2–5)</span>
            </div>
            <div class="link-selector" :class="{ open: selectorOpen }" ref="selectorRef">
              <div class="selector-trigger" @click="selectorOpen = !selectorOpen">
                <span v-if="selectedIds.length === 0" class="selector-placeholder">Choose links to compare…</span>
                <span v-else class="selector-value">{{ selectedIds.length }} link{{ selectedIds.length > 1 ? 's' : '' }} selected</span>
                <span class="material-symbols-outlined selector-chevron">expand_more</span>
              </div>
              <div v-if="selectorOpen" class="selector-dropdown">
                <div style="padding:0.5rem 0.75rem 0.25rem;">
                  <input
                    class="form-input"
                    :value="linkSearch"
                    @input="linkSearch = ($event.target as HTMLInputElement).value"
                    placeholder="Search links…"
                    @click.stop
                  />
                </div>
                <div class="selector-options">
                  <label
                    v-for="link in filteredLinks"
                    :key="link.id"
                    class="selector-option"
                    :class="{ disabled: !selectedIds.includes(link.id) && selectedIds.length >= 5 }"
                    @click.stop
                  >
                    <input
                      type="checkbox"
                      :value="link.id"
                      v-model="selectedIds"
                      :disabled="!selectedIds.includes(link.id) && selectedIds.length >= 5"
                      class="selector-checkbox"
                    />
                    <div style="min-width:0;">
                      <div class="selector-option__title">{{ link.title || link.slug }}</div>
                      <div class="selector-option__url">{{ link.short_url }}</div>
                    </div>
                  </label>
                  <div v-if="filteredLinks.length === 0" class="selector-empty">No links found</div>
                </div>
              </div>
            </div>
          </div>

          <!-- From date -->
          <label class="date-field">
            <span class="date-field__label">From</span>
            <input type="date" class="form-input form-input--date" :value="filterFrom" @input="filterFrom = ($event.target as HTMLInputElement).value" />
          </label>

          <!-- To date -->
          <label class="date-field">
            <span class="date-field__label">To</span>
            <input type="date" class="form-input form-input--date" :value="filterTo" @input="filterTo = ($event.target as HTMLInputElement).value" />
          </label>

          <!-- Actions -->
          <div class="filter-actions">
            <button class="btn-filled"
              style="flex:1;"
              :disabled="loading || selectedIds.length < 2"
              @click="runComparison"
            >
              <span v-if="loading" class="css-spinner css-spinner--sm css-spinner--white"></span>
              Compare
            </button>
            <button class="btn-icon" v-if="result" title="Export to CSV" @click="exportCSV">
              <span class="material-symbols-outlined">download</span>
            </button>
          </div>
        </div>

        <!-- Validation hint -->
        <div v-if="selectedIds.length === 1" class="validation-hint">
          <span class="material-symbols-outlined" style="font-size:14px;vertical-align:middle;">info</span>
          Select at least one more link to run a comparison.
        </div>

        <!-- Selected chips -->
        <div v-if="selectedIds.length > 0" class="selected-chips">
          <span v-for="id in selectedIds" :key="id" class="sel-chip">
            {{ getLinkLabel(id) }}
            <button class="sel-chip__remove" @click="selectedIds = selectedIds.filter(s => s !== id)">
              <span class="material-symbols-outlined">close</span>
            </button>
          </span>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="error-banner">
      <span class="material-symbols-outlined" style="color:var(--md-sys-color-error);">error</span>
      <span style="flex:1;">{{ error }}</span>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <span class="css-spinner"></span>
      <span class="loading-label">Fetching comparison data…</span>
    </div>

    <template v-if="!loading && result">
      <!-- Period info -->
      <div class="period-info">
        Showing data from <strong>{{ formatDate(result.from) }}</strong> to <strong>{{ formatDate(result.to) }}</strong>
        &nbsp;·&nbsp; <strong>{{ result.span_days }}</strong> day{{ result.span_days !== 1 ? 's' : '' }}
      </div>

      <!-- Bar chart — total clicks -->
      <div class="an-card">
        <div class="an-card-header">
          <span class="an-card-title">Total Clicks Comparison</span>
        </div>
        <div style="padding:1rem;">
          <VChart :option="barChartOption" style="height:260px;" autoresize />
        </div>
      </div>

      <!-- Comparison table -->
      <div class="an-card">
        <div class="an-card-header">
          <div>
            <span class="an-card-title">Detailed Breakdown</span>
            <div class="an-card-subtitle">Per-link metrics for the selected period</div>
          </div>
        </div>
        <div class="m3-table-wrapper">
          <table class="m3-table">
            <thead>
              <tr>
                <th>Link</th>
                <th style="text-align:right;">Total Clicks</th>
                <th style="text-align:right;">Unique Clicks</th>
                <th style="text-align:right;">Clicks / Day</th>
                <th>Top Referrer</th>
                <th>Top Country</th>
                <th>Top Browser</th>
                <th>Top Device</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(link, idx) in result.links" :key="link.link_id">
                <td>
                  <div class="link-cell">
                    <span class="rank-badge" :style="{ background: chartColors[idx % chartColors.length] }">{{ idx + 1 }}</span>
                    <div style="min-width:0;">
                      <div class="link-cell__title" :title="link.title || link.slug">{{ link.title || link.slug }}</div>
                      <div class="link-cell__slug">{{ link.slug }}</div>
                    </div>
                  </div>
                </td>
                <td style="text-align:right;font-weight:600;">{{ link.total_clicks.toLocaleString() }}</td>
                <td style="text-align:right;">{{ link.unique_clicks.toLocaleString() }}</td>
                <td style="text-align:right;">{{ link.clicks_per_day.toFixed(1) }}</td>
                <td>
                  <span v-if="link.top_referrer" class="m3-badge m3-badge--primary">{{ link.top_referrer }}</span>
                  <span v-else class="cell-muted">—</span>
                </td>
                <td>
                  <span v-if="link.top_country" class="m3-badge m3-badge--primary">{{ link.top_country }}</span>
                  <span v-else class="cell-muted">—</span>
                </td>
                <td>
                  <span v-if="link.top_browser" class="m3-badge m3-badge--primary">{{ link.top_browser }}</span>
                  <span v-else class="cell-muted">—</span>
                </td>
                <td>
                  <span v-if="link.top_device" class="m3-badge m3-badge--primary">{{ link.top_device }}</span>
                  <span v-else class="cell-muted">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Empty state -->
    <div v-if="!loading && !result && !error" class="empty-state-card">
      <div class="empty-icon">
        <span class="material-symbols-outlined">bar_chart</span>
      </div>
      <div class="empty-title">No comparison yet</div>
      <p class="empty-desc">Select 2–5 links above and click <strong>Compare</strong> to see their performance side by side.</p>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import VChart from 'vue-echarts';
import { use } from 'echarts/core';
import { BarChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import { linksApi } from '@/api/links';
import type { LinkResponse, MultiLinkComparisonResponse } from '@/types/links';

use([CanvasRenderer, BarChart, GridComponent, TooltipComponent, LegendComponent]);

// ── State ─────────────────────────────────────────────────────────────────────
const allLinks = ref<LinkResponse[]>([]);
const selectedIds = ref<string[]>([]);
const linkSearch = ref('');
const selectorOpen = ref(false);
const selectorRef = ref<HTMLElement | null>(null);

const filterFrom = ref(dateString(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)));
const filterTo = ref(dateString(new Date()));

const loading = ref(false);
const error = ref('');
const result = ref<MultiLinkComparisonResponse | null>(null);

const chartColors = ['#635bff', '#22c55e', '#f59e0b', '#ef4444', '#06b6d4'];

// ── Helpers ───────────────────────────────────────────────────────────────────
function dateString(d: Date): string {
  return d.toISOString().slice(0, 10);
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

function getLinkLabel(id: string): string {
  const link = allLinks.value.find((l) => l.id === id);
  return link ? (link.title || link.slug) : id.slice(0, 8);
}

// ── Computed ──────────────────────────────────────────────────────────────────
const filteredLinks = computed(() => {
  const q = linkSearch.value.toLowerCase();
  if (!q) return allLinks.value;
  return allLinks.value.filter(
    (l) => (l.title || '').toLowerCase().includes(q) || l.slug.toLowerCase().includes(q),
  );
});

const barChartOption = computed(() => {
  if (!result.value) return {};
  const labels = result.value.links.map((l) => l.title || l.slug);
  const totalData = result.value.links.map((l) => l.total_clicks);
  const uniqueData = result.value.links.map((l) => l.unique_clicks);

  return {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { data: ['Total Clicks', 'Unique Clicks'], bottom: 0 },
    grid: { left: '3%', right: '4%', bottom: 40, top: 20, containLabel: true },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: {
        overflow: 'truncate',
        width: 80,
        interval: 0,
      },
    },
    yAxis: { type: 'value' },
    series: [
      {
        name: 'Total Clicks',
        type: 'bar',
        data: totalData,
        itemStyle: { color: '#635bff', borderRadius: [4, 4, 0, 0] },
      },
      {
        name: 'Unique Clicks',
        type: 'bar',
        data: uniqueData,
        itemStyle: { color: '#22c55e', borderRadius: [4, 4, 0, 0] },
      },
    ],
  };
});

// ── Actions ───────────────────────────────────────────────────────────────────
async function loadLinks() {
  try {
    // Fetch up to 200 links for the selector
    const resp = await linksApi.list(1, 200);
    allLinks.value = resp.data?.links ?? [];
  } catch {
    // non-fatal; selector just stays empty
  }
}

async function runComparison() {
  if (selectedIds.value.length < 2) return;
  error.value = '';
  loading.value = true;
  result.value = null;
  selectorOpen.value = false;
  try {
    const from = filterFrom.value ? `${filterFrom.value}T00:00:00Z` : undefined;
    const to = filterTo.value ? `${filterTo.value}T23:59:59Z` : undefined;
    const resp = await linksApi.compareLinks(selectedIds.value, from, to);
    result.value = resp.data ?? null;
  } catch (e: unknown) {
    error.value = (e as { response?: { data?: { message?: string } } })?.response?.data?.message ?? 'Failed to load comparison.';
  } finally {
    loading.value = false;
  }
}

function exportCSV() {
  if (!result.value) return;
  const headers = ['Slug', 'Title', 'Total Clicks', 'Unique Clicks', 'Clicks/Day', 'Top Referrer', 'Top Country', 'Top Browser', 'Top Device'];
  const rows = result.value.links.map((l) => [
    l.slug,
    l.title ?? '',
    l.total_clicks,
    l.unique_clicks,
    l.clicks_per_day.toFixed(2),
    l.top_referrer ?? '',
    l.top_country ?? '',
    l.top_browser ?? '',
    l.top_device ?? '',
  ]);
  const csv = [headers, ...rows].map((r) => r.map((v) => `"${String(v).replace(/"/g, '""')}"`).join(',')).join('\n');
  const blob = new Blob([csv], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `link-comparison-${dateString(new Date())}.csv`;
  a.click();
  URL.revokeObjectURL(url);
}

// Close dropdown when clicking outside
function handleClickOutside(e: MouseEvent) {
  if (selectorRef.value && !selectorRef.value.contains(e.target as Node)) {
    selectorOpen.value = false;
  }
}

onMounted(() => {
  loadLinks();
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped lang="scss">
/* ── Root ─────────────────────────────────────────────────────────────────── */
.link-comparison-page {
  max-width: 1400px;
  padding: 0;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page header ──────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  margin: 0;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0;
}

/* ── AN Card ──────────────────────────────────────────────────────────────── */
.an-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);
}

.an-card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 18px; }

  &--primary {
    background: var(--md-sys-color-primary-container);
    .material-symbols-outlined { color: var(--md-sys-color-on-primary-container); }
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  flex: 1;
}

.an-card-subtitle {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
}

.an-card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ── CSS Spinner ──────────────────────────────────────────────────────────── */
.css-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2.5px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm {
    width: 16px;
    height: 16px;
    border-width: 2px;
  }

  &--white {
    border-color: rgba(255,255,255,0.35);
    border-top-color: #fff;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Form field ───────────────────────────────────────────────────────────── */
.form-field__label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
}

.form-input {
  height: 40px;
  padding: 0 12px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9375rem;
  font-family: inherit;
  transition: border-color 0.15s, box-shadow 0.15s;
  width: 100%;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }

  &--date {
    width: auto;
    min-width: 150px;
  }
}

/* ── Date field ───────────────────────────────────────────────────────────── */
.date-field {
  display: flex;
  flex-direction: column;
  gap: 4px;

  &__label {
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }
}

/* ── Filter grid ──────────────────────────────────────────────────────────── */
.filter-grid {
  display: grid;
  grid-template-columns: minmax(200px, 1fr) auto auto auto;
  gap: 12px;
  align-items: end;

  @media (max-width: 767px) {
    grid-template-columns: 1fr 1fr;
  }

  @media (max-width: 575px) {
    grid-template-columns: 1fr;
  }
}

.filter-actions {
  display: flex;
  gap: 8px;
  align-items: flex-end;
  padding-bottom: 2px;
}

/* ── Validation hint ──────────────────────────────────────────────────────── */
.validation-hint {
  font-size: 0.8rem;
  color: #d97706;
}

/* ── Selected chips ───────────────────────────────────────────────────────── */
.selected-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.sel-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px 4px 10px;
  border-radius: 20px;
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);
  font-size: 0.8rem;
  font-weight: 500;
}

.sel-chip__remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--md-sys-color-on-secondary-container);
  padding: 0;
  flex-shrink: 0;

  &:hover { background: rgba(0,0,0,0.08); }

  .material-symbols-outlined { font-size: 14px; }
}

/* ── Link selector ────────────────────────────────────────────────────────── */
.link-selector-wrapper {
  min-width: 0;
}

.link-selector {
  position: relative;
}

.selector-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  height: 40px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  cursor: pointer;
  user-select: none;
  background: var(--md-sys-color-surface);
  transition: border-color 0.15s;

  &:hover { border-color: var(--md-sys-color-on-surface); }
}

.selector-placeholder {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.9375rem;
}

.selector-value {
  font-size: 0.9375rem;
  color: var(--md-sys-color-on-surface);
}

.selector-chevron {
  font-size: 18px;
  flex-shrink: 0;
  color: var(--md-sys-color-on-surface-variant);
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 400;
  border-radius: 12px;
  overflow: hidden;
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.selector-options {
  max-height: 240px;
  overflow-y: auto;
  padding: 4px 0;
}

.selector-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.12s;

  &:hover { background: var(--md-sys-color-surface-container-low); }

  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &__title {
    font-size: 0.8125rem;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--md-sys-color-on-surface);
  }

  &__url {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.7rem;
  }
}

.selector-checkbox {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  accent-color: var(--md-sys-color-primary);
  cursor: pointer;
}

.selector-empty {
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  padding: 1rem;
}

/* ── Error banner ─────────────────────────────────────────────────────────── */
.error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  background: rgba(176, 0, 32, 0.08);
  border: 1px solid var(--md-sys-color-error);
  font-size: 0.875rem;
}

/* ── Loading ──────────────────────────────────────────────────────────────── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 0;
  gap: 1rem;
}

.loading-label {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}

/* ── Period info ──────────────────────────────────────────────────────────── */
.period-info {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}

/* ── Table ────────────────────────────────────────────────────────────────── */
.m3-table-wrapper {
  overflow-x: auto;
}

.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;

  th {
    padding: 10px 16px;
    text-align: left;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    white-space: nowrap;
  }

  td {
    padding: 10px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td { border-bottom: none; }

  tbody tr:hover td { background: var(--md-sys-color-surface-container-low); }
}

.cell-muted {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.8rem;
}

/* ── Link cell ────────────────────────────────────────────────────────────── */
.link-cell {
  display: flex;
  align-items: center;
  gap: 8px;

  &__title {
    font-size: 0.875rem;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 180px;
  }

  &__slug {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.75rem;
  }
}

/* ── Rank badge ───────────────────────────────────────────────────────────── */
.rank-badge {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* ── M3 Badges ────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.7rem;
  font-weight: 500;
  padding: 0.15rem 0.6rem;
  border-radius: 6px;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;

  &--primary {
    background: rgba(99, 91, 255, 0.1);
    color: var(--md-sys-color-primary);
  }

  &--neutral {
    background: rgba(107, 114, 128, 0.1);
    color: #6b7280;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Empty state ──────────────────────────────────────────────────────────── */
.empty-state-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.empty-icon {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;

  .material-symbols-outlined {
    font-size: 2rem;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.6;
  }
}

.empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 0.5rem;
}

.empty-desc {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  max-width: 400px;
  margin: 0;
}
</style>
