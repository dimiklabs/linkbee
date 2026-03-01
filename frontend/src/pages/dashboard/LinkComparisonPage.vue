<template>
  <UpgradeWall
    v-if="!billingStore.isPaidPlan && billingStore.loaded"
    icon="compare_arrows"
    title="Link comparison is a Pro feature"
    description="Upgrade to Pro to compare analytics side-by-side across multiple links."
  />
  <div v-else class="link-comparison-page">

    <!-- Page Header -->
    <div class="page-header">
      <h1 class="md-title-large page-title">Compare Links</h1>
      <p class="page-subtitle">Select 2–5 links to compare their performance side by side.</p>
    </div>

    <!-- Filters Card — overflow:visible so the dropdown escapes the card -->
    <div class="an-card an-card--selector">
      <div class="an-card-header">
        <div class="an-card-icon an-card-icon--primary">
          <span class="material-symbols-outlined">compare_arrows</span>
        </div>
        <span class="an-card-title">Configure Comparison</span>
      </div>
      <div class="an-card-body">
        <div class="filter-grid">

          <!-- Link multi-selector -->
          <div class="link-selector-wrapper">
            <div class="form-label">
              Links <span class="form-label__hint">(select 2–5)</span>
            </div>
            <div class="link-selector" :class="{ open: selectorOpen }" ref="selectorRef">
              <div class="selector-trigger" @click="selectorOpen = !selectorOpen">
                <span v-if="selectedIds.length === 0" class="selector-placeholder">Choose links to compare…</span>
                <span v-else class="selector-value">
                  <span class="selector-count">{{ selectedIds.length }}</span>
                  link{{ selectedIds.length > 1 ? 's' : '' }} selected
                </span>
                <span class="material-symbols-outlined selector-chevron"
                  :class="{ 'selector-chevron--open': selectorOpen }">expand_more</span>
              </div>
              <div v-if="selectorOpen" class="selector-dropdown">
                <div class="selector-search">
                  <span class="material-symbols-outlined selector-search-icon">search</span>
                  <input
                    class="selector-search-input"
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
                    :class="{
                      disabled: !selectedIds.includes(link.id) && selectedIds.length >= 5,
                      selected: selectedIds.includes(link.id),
                    }"
                    @click.stop
                  >
                    <input
                      type="checkbox"
                      :value="link.id"
                      v-model="selectedIds"
                      :disabled="!selectedIds.includes(link.id) && selectedIds.length >= 5"
                      class="selector-checkbox"
                    />
                    <div class="selector-option__text">
                      <div class="selector-option__title">{{ link.title || link.slug }}</div>
                      <div class="selector-option__url">{{ link.short_url }}</div>
                    </div>
                  </label>
                  <div v-if="filteredLinks.length === 0" class="selector-empty">No links found</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Date range -->
          <label class="date-field">
            <span class="form-label">From</span>
            <input type="date" class="form-input" :value="filterFrom"
              @input="filterFrom = ($event.target as HTMLInputElement).value" />
          </label>

          <label class="date-field">
            <span class="form-label">To</span>
            <input type="date" class="form-input" :value="filterTo"
              @input="filterTo = ($event.target as HTMLInputElement).value" />
          </label>

          <!-- Actions -->
          <div class="filter-actions">
            <button class="btn-filled compare-btn"
              :disabled="loading || selectedIds.length < 2"
              @click="runComparison"
            >
              <span v-if="loading" class="css-spinner css-spinner--sm css-spinner--white"></span>
              <span class="material-symbols-outlined" v-else style="font-size:18px;">compare_arrows</span>
              Compare
            </button>
            <button v-if="result" class="btn-icon export-btn" title="Export to CSV" @click="exportCSV">
              <span class="material-symbols-outlined">download</span>
            </button>
          </div>
        </div>

        <!-- Validation hint -->
        <p v-if="selectedIds.length === 1" class="validation-hint">
          <span class="material-symbols-outlined" style="font-size:15px;vertical-align:-3px;">info</span>
          Select at least one more link to enable comparison.
        </p>

        <!-- Selected chips -->
        <div v-if="selectedIds.length > 0" class="selected-chips">
          <span v-for="id in selectedIds" :key="id" class="sel-chip">
            <span class="material-symbols-outlined sel-chip__icon">link</span>
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
      <span class="material-symbols-outlined">error</span>
      <span>{{ error }}</span>
      <button class="btn-icon" style="margin-left:auto;" @click="error = ''">
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <span class="css-spinner"></span>
      <span class="loading-label">Fetching comparison data…</span>
    </div>

    <template v-if="!loading && result">
      <!-- Period info bar -->
      <div class="period-info">
        <span class="material-symbols-outlined" style="font-size:15px;vertical-align:-3px;opacity:.6;">calendar_today</span>
        Showing <strong>{{ formatDate(result.from) }}</strong> – <strong>{{ formatDate(result.to) }}</strong>
        &nbsp;·&nbsp; <strong>{{ result.span_days }}</strong> day{{ result.span_days !== 1 ? 's' : '' }}
      </div>

      <!-- Bar chart — total clicks -->
      <div class="an-card">
        <div class="an-card-header">
          <span class="material-symbols-outlined an-card-header-icon">bar_chart</span>
          <span class="an-card-title">Total Clicks Comparison</span>
        </div>
        <div class="chart-wrap">
          <VChart :option="barChartOption" style="height:280px;" autoresize />
        </div>
      </div>

      <!-- Comparison table -->
      <div class="an-card">
        <div class="an-card-header">
          <span class="material-symbols-outlined an-card-header-icon">table_chart</span>
          <div>
            <div class="an-card-title">Detailed Breakdown</div>
            <div class="an-card-subtitle">Per-link metrics for the selected period</div>
          </div>
        </div>
        <div class="m3-table-wrapper">
          <table class="m3-table">
            <thead>
              <tr>
                <th>Link</th>
                <th class="text-right">Total Clicks</th>
                <th class="text-right">Unique Clicks</th>
                <th class="text-right">Clicks / Day</th>
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
                <td class="text-right fw-600">{{ link.total_clicks.toLocaleString() }}</td>
                <td class="text-right">{{ link.unique_clicks.toLocaleString() }}</td>
                <td class="text-right">{{ link.clicks_per_day.toFixed(1) }}</td>
                <td><span v-if="link.top_referrer" class="m3-badge m3-badge--primary">{{ link.top_referrer }}</span><span v-else class="cell-muted">—</span></td>
                <td><span v-if="link.top_country" class="m3-badge m3-badge--primary">{{ link.top_country }}</span><span v-else class="cell-muted">—</span></td>
                <td><span v-if="link.top_browser" class="m3-badge m3-badge--primary">{{ link.top_browser }}</span><span v-else class="cell-muted">—</span></td>
                <td><span v-if="link.top_device" class="m3-badge m3-badge--primary">{{ link.top_device }}</span><span v-else class="cell-muted">—</span></td>
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
import UpgradeWall from '@/components/UpgradeWall.vue';
import { useBillingStore } from '@/stores/billing';

const billingStore = useBillingStore();
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
    grid: { left: '3%', right: '4%', bottom: 48, top: 20, containLabel: true },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: { overflow: 'truncate', width: 100, interval: 0 },
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
    // Backend max limit is 100; fetch two pages to cover up to 200 links
    const [p1, p2] = await Promise.all([
      linksApi.list(1, 100),
      linksApi.list(2, 100),
    ]);
    const page1 = p1.data?.links ?? [];
    const page2 = p2.data?.links ?? [];
    allLinks.value = [...page1, ...page2];
  } catch {
    // non-fatal; selector stays empty
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
    l.slug, l.title ?? '', l.total_clicks, l.unique_clicks,
    l.clicks_per_day.toFixed(2), l.top_referrer ?? '',
    l.top_country ?? '', l.top_browser ?? '', l.top_device ?? '',
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

function handleClickOutside(e: MouseEvent) {
  if (selectorRef.value && !selectorRef.value.contains(e.target as Node)) {
    selectorOpen.value = false;
  }
}

onMounted(() => {
  billingStore.fetchPlan();
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
.page-header { display: flex; flex-direction: column; gap: 4px; }

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

/* Filter card needs overflow:visible so the dropdown can escape */
.an-card--selector {
  overflow: visible;

  > .an-card-header {
    border-radius: 13px 13px 0 0;
  }
}

.an-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);
}

.an-card-header-icon {
  font-size: 18px;
  color: var(--md-sys-color-primary);
  flex-shrink: 0;
}

.an-card-icon {
  width: 34px;
  height: 34px;
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
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
}

.an-card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.chart-wrap {
  padding: 16px 20px 8px;
}

/* ── CSS Spinner ──────────────────────────────────────────────────────────── */
.css-spinner {
  display: inline-block;
  width: 20px; height: 20px;
  border: 2.5px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm  { width: 16px; height: 16px; border-width: 2px; }
  &--white { border-color: rgba(255,255,255,0.35); border-top-color: #fff; }
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ── Form labels ──────────────────────────────────────────────────────────── */
.form-label {
  display: block;
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.05em;

  &__hint {
    font-weight: 400;
    text-transform: none;
    letter-spacing: 0;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.7;
  }
}

/* ── Form input ───────────────────────────────────────────────────────────── */
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
}

/* ── Date field ───────────────────────────────────────────────────────────── */
.date-field {
  display: flex;
  flex-direction: column;

  .form-input { width: auto; min-width: 148px; }
}

/* ── Filter grid ──────────────────────────────────────────────────────────── */
.filter-grid {
  display: grid;
  grid-template-columns: 1fr auto auto auto;
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
}

.compare-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.export-btn {
  height: 40px;
  width: 40px;
  flex-shrink: 0;
}

/* ── Validation hint ──────────────────────────────────────────────────────── */
.validation-hint {
  margin: 0;
  font-size: 0.8rem;
  color: #d97706;
  display: flex;
  align-items: center;
  gap: 5px;
}

/* ── Selected chips ───────────────────────────────────────────────────────── */
.selected-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.sel-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 6px 4px 8px;
  border-radius: 20px;
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);
  font-size: 0.8rem;
  font-weight: 500;

  &__icon { font-size: 13px; opacity: 0.7; }

  &__remove {
    display: flex; align-items: center; justify-content: center;
    width: 18px; height: 18px;
    border-radius: 50%;
    background: transparent; border: none;
    cursor: pointer;
    color: var(--md-sys-color-on-secondary-container);
    padding: 0; flex-shrink: 0;

    &:hover { background: rgba(0,0,0,0.1); }
    .material-symbols-outlined { font-size: 13px; }
  }
}

/* ── Link selector ────────────────────────────────────────────────────────── */
.link-selector-wrapper { min-width: 0; }

.link-selector { position: relative; }

.selector-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  height: 40px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  cursor: pointer;
  user-select: none;
  background: var(--md-sys-color-surface);
  transition: border-color 0.15s, box-shadow 0.15s;

  &:hover { border-color: var(--md-sys-color-on-surface); }

  .link-selector.open & {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }
}

.selector-placeholder {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.9375rem;
}

.selector-value {
  font-size: 0.9375rem;
  color: var(--md-sys-color-on-surface);
  display: flex;
  align-items: center;
  gap: 6px;
}

.selector-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  border-radius: 10px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-size: 0.72rem;
  font-weight: 700;
  padding: 0 5px;
}

.selector-chevron {
  font-size: 18px;
  flex-shrink: 0;
  color: var(--md-sys-color-on-surface-variant);
  transition: transform 0.2s ease;

  &--open { transform: rotate(180deg); }
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 500;
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  box-shadow: 0 8px 24px rgba(0,0,0,0.14);
  overflow: hidden;
}

.selector-search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px 8px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.selector-search-icon {
  font-size: 16px;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
}

.selector-search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  color: var(--md-sys-color-on-surface);
  font-size: 0.875rem;
  font-family: inherit;
}

.selector-options {
  max-height: 240px;
  overflow-y: auto;
  padding: 4px 0;
}

.selector-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 14px;
  cursor: pointer;
  transition: background 0.1s;

  &:hover { background: var(--md-sys-color-surface-container-low); }

  &.selected { background: color-mix(in srgb, var(--md-sys-color-primary) 6%, transparent); }

  &.disabled { opacity: 0.45; cursor: not-allowed; }

  &__text { min-width: 0; }

  &__title {
    font-size: 0.825rem;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--md-sys-color-on-surface);
  }

  &__url {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.72rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.selector-checkbox {
  flex-shrink: 0;
  width: 16px; height: 16px;
  accent-color: var(--md-sys-color-primary);
  cursor: pointer;
}

.selector-empty {
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  padding: 1.25rem;
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
  color: var(--md-sys-color-error);
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
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* ── Table ────────────────────────────────────────────────────────────────── */
.m3-table-wrapper { overflow-x: auto; }

.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;

  th {
    padding: 10px 16px;
    text-align: left;
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    white-space: nowrap;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  td {
    padding: 10px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td { border-bottom: none; }
  tbody tr:hover td { background: var(--md-sys-color-surface-container-low); }
}

.text-right { text-align: right !important; }
.fw-600 { font-weight: 600; }
.cell-muted { color: var(--md-sys-color-on-surface-variant); font-size: 0.8rem; }

/* ── Link cell ────────────────────────────────────────────────────────────── */
.link-cell {
  display: flex; align-items: center; gap: 10px;

  &__title {
    font-size: 0.875rem; font-weight: 500;
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
    max-width: 180px;
  }

  &__slug { color: var(--md-sys-color-on-surface-variant); font-size: 0.72rem; }
}

/* ── Rank badge ───────────────────────────────────────────────────────────── */
.rank-badge {
  width: 24px; height: 24px;
  border-radius: 50%;
  color: #fff;
  font-size: 0.7rem; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}

/* ── M3 Badges ────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex; align-items: center;
  font-size: 0.7rem; font-weight: 500;
  padding: 0.15rem 0.55rem;
  border-radius: 6px;
  max-width: 120px;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  vertical-align: middle;

  &--primary { background: rgba(99, 91, 255, 0.1); color: var(--md-sys-color-primary); }
}

/* ── Empty state ──────────────────────────────────────────────────────────── */
.empty-state-card {
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  padding: 5rem 2rem; text-align: center;
}

.empty-icon {
  width: 72px; height: 72px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-low);
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 1.25rem;

  .material-symbols-outlined { font-size: 2rem; color: var(--md-sys-color-on-surface-variant); opacity: 0.5; }
}

.empty-title {
  font-size: 1.125rem; font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 0.5rem;
}

.empty-desc {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  max-width: 380px; margin: 0;
  line-height: 1.6;
}
</style>
