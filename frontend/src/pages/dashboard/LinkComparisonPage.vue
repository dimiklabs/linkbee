<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="md-title-large">Compare Links</h1>
        <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin-top:0.25rem;">Select 2–5 links to compare their performance side by side.</p>
      </div>
    </div>

    <!-- Filters Card -->
    <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;padding:1.25rem;">
      <div class="filter-grid">
        <!-- Link selector -->
        <div class="link-selector-wrapper">
          <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:0.5rem;font-weight:500;">
            Links <span style="font-weight:400;">(select 2–5)</span>
          </div>
          <div class="link-selector" :class="{ open: selectorOpen }" ref="selectorRef">
            <div
              class="selector-trigger"
              @click="selectorOpen = !selectorOpen"
            >
              <span v-if="selectedIds.length === 0" style="color:var(--md-sys-color-on-surface-variant);" class="md-body-medium">Choose links to compare…</span>
              <span v-else class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ selectedIds.length }} link{{ selectedIds.length > 1 ? 's' : '' }} selected</span>
              <span class="material-symbols-outlined" style="font-size:18px;flex-shrink:0;color:var(--md-sys-color-on-surface-variant);">expand_more</span>
            </div>
            <div v-if="selectorOpen" class="selector-dropdown">
              <div style="padding:0.5rem 0.75rem 0.25rem;">
                <md-outlined-text-field
                  :value="linkSearch"
                  @input="linkSearch = ($event.target as HTMLInputElement).value"
                  label="Search links…"
                  style="width:100%;"
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
                    <div class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ link.title || link.slug }}</div>
                    <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ link.short_url }}</div>
                  </div>
                </label>
                <div v-if="filteredLinks.length === 0" class="md-body-small" style="text-align:center;color:var(--md-sys-color-on-surface-variant);padding:1rem;">No links found</div>
              </div>
            </div>
          </div>
        </div>

        <!-- From date -->
        <md-outlined-text-field
          type="date"
          label="From"
          :value="filterFrom"
          @input="filterFrom = ($event.target as HTMLInputElement).value"
          style="min-width:150px;"
        />

        <!-- To date -->
        <md-outlined-text-field
          type="date"
          label="To"
          :value="filterTo"
          @input="filterTo = ($event.target as HTMLInputElement).value"
          style="min-width:150px;"
        />

        <!-- Actions -->
        <div style="display:flex;gap:0.5rem;align-items:flex-end;padding-bottom:0.125rem;">
          <md-filled-button
            style="flex:1;"
            :disabled="loading || selectedIds.length < 2"
            @click="runComparison"
          >
            <span v-if="loading" slot="icon"><md-circular-progress indeterminate style="--md-circular-progress-size:18px" /></span>
            Compare
          </md-filled-button>
          <md-icon-button
            v-if="result"
            title="Export to CSV"
            @click="exportCSV"
          >
            <span class="material-symbols-outlined">download</span>
          </md-icon-button>
        </div>
      </div>

      <!-- Validation hint -->
      <div v-if="selectedIds.length === 1" style="margin-top:0.75rem;">
        <span class="md-body-small" style="color:#d97706;">
          <span class="material-symbols-outlined" style="font-size:14px;vertical-align:middle;">info</span>
          Select at least one more link to run a comparison.
        </span>
      </div>

      <!-- Selected chips -->
      <div v-if="selectedIds.length > 0" style="margin-top:0.75rem;display:flex;flex-wrap:wrap;gap:0.5rem;">
        <md-chip-set>
          <md-input-chip
            v-for="id in selectedIds"
            :key="id"
            :label="getLinkLabel(id)"
            @remove="selectedIds = selectedIds.filter(s => s !== id)"
          />
        </md-chip-set>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="error-banner" style="margin-bottom:1rem;">
      <span class="material-symbols-outlined" style="color:var(--md-sys-color-error);">error</span>
      <span class="md-body-medium" style="flex:1;">{{ error }}</span>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;flex-direction:column;align-items:center;justify-content:center;padding:4rem 0;gap:1rem;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:48px" />
      <span class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);">Fetching comparison data…</span>
    </div>

    <template v-if="!loading && result">
      <!-- Period info -->
      <div class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:1rem;">
        Showing data from <strong>{{ formatDate(result.from) }}</strong> to <strong>{{ formatDate(result.to) }}</strong>
        &nbsp;·&nbsp; <strong>{{ result.span_days }}</strong> day{{ result.span_days !== 1 ? 's' : '' }}
      </div>

      <!-- Bar chart — total clicks -->
      <div class="m3-card m3-card--elevated" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <span class="md-title-medium">Total Clicks Comparison</span>
        </div>
        <md-divider />
        <div style="padding:1rem;">
          <VChart :option="barChartOption" style="height:260px;" autoresize />
        </div>
      </div>

      <!-- Comparison table -->
      <div class="m3-card m3-card--outlined">
        <div class="card-header-row">
          <span class="md-title-medium">Detailed Breakdown</span>
        </div>
        <md-divider />
        <div style="overflow-x:auto;">
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
                  <div style="display:flex;align-items:center;gap:0.5rem;">
                    <span class="rank-badge" :style="{ background: chartColors[idx % chartColors.length] }">{{ idx + 1 }}</span>
                    <div style="min-width:0;">
                      <div class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:180px;" :title="link.title || link.slug">
                        {{ link.title || link.slug }}
                      </div>
                      <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ link.slug }}</div>
                    </div>
                  </div>
                </td>
                <td style="text-align:right;font-weight:600;">{{ link.total_clicks.toLocaleString() }}</td>
                <td style="text-align:right;">{{ link.unique_clicks.toLocaleString() }}</td>
                <td style="text-align:right;">{{ link.clicks_per_day.toFixed(1) }}</td>
                <td>
                  <span v-if="link.top_referrer" class="m3-badge m3-badge--primary">{{ link.top_referrer }}</span>
                  <span v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">—</span>
                </td>
                <td>
                  <span v-if="link.top_country" class="m3-badge m3-badge--primary">{{ link.top_country }}</span>
                  <span v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">—</span>
                </td>
                <td>
                  <span v-if="link.top_browser" class="m3-badge m3-badge--primary">{{ link.top_browser }}</span>
                  <span v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">—</span>
                </td>
                <td>
                  <span v-if="link.top_device" class="m3-badge m3-badge--primary">{{ link.top_device }}</span>
                  <span v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Empty state -->
    <div v-if="!loading && !result && !error" class="empty-state" style="padding:4rem 2rem;">
      <span class="material-symbols-outlined" style="font-size:3rem;display:block;margin-bottom:0.75rem;color:var(--md-sys-color-on-surface-variant);">bar_chart</span>
      <p class="md-body-medium">Select 2–5 links above and click <strong>Compare</strong> to see their performance side by side.</p>
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
.page-wrapper {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page header */
.page-header {
  margin-bottom: 1.5rem;
}

/* Card header row */
.card-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1.25rem;
  gap: 1rem;
  flex-wrap: wrap;
}

/* Error banner */
.error-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  background: rgba(176, 0, 32, 0.08);
  border: 1px solid var(--md-sys-color-error);
}

/* Filter grid */
.filter-grid {
  display: grid;
  grid-template-columns: minmax(200px, 1fr) auto auto auto;
  gap: 0.75rem;
  align-items: end;

  @media (max-width: 767px) {
    grid-template-columns: 1fr 1fr;
  }

  @media (max-width: 575px) {
    grid-template-columns: 1fr;
  }
}

/* Link selector */
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
  height: 56px;
  border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-radius: 4px;
  cursor: pointer;
  user-select: none;
  background: var(--md-sys-color-surface, #fff);
  transition: border-color 0.15s;

  &:hover {
    border-color: var(--md-sys-color-on-surface, #1a1f36);
  }
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 400;
  border-radius: 12px;
  overflow: hidden;
  background: var(--md-sys-color-surface, #fff);
  border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.selector-options {
  max-height: 240px;
  overflow-y: auto;
  padding: 0.25rem 0;
}

.selector-option {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.625rem 1rem;
  cursor: pointer;
  transition: background 0.12s;

  &:hover {
    background: var(--md-sys-color-surface-container-low, #f7f9fc);
  }

  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.selector-checkbox {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  accent-color: var(--md-sys-color-primary, #635bff);
  cursor: pointer;
}

/* Rank badge */
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

/* M3 Cards */
.m3-card {
  border-radius: 12px;
  overflow: hidden;

  &--elevated {
    background: var(--md-sys-color-surface-container-low, #fff);
    box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.08);
  }

  &--outlined {
    background: var(--md-sys-color-surface, #fff);
    border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  }
}

/* M3 Badges */
.m3-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.7rem;
  font-weight: 500;
  padding: 0.15rem 0.6rem;
  border-radius: 999px;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;

  &--primary {
    background: rgba(99, 91, 255, 0.1);
    color: var(--md-sys-color-primary, #635bff);
  }

  &--neutral {
    background: rgba(107, 114, 128, 0.1);
    color: #6b7280;
  }
}

/* M3 Table */
.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;

  th {
    padding: 0.625rem 1rem;
    text-align: left;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low, #f8f9fa);
    border-bottom: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
    white-space: nowrap;
  }

  td {
    padding: 0.625rem 1rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant, #f0f2f5);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  tbody tr:hover td {
    background: var(--md-sys-color-surface-container-low, #f8f9fa);
  }
}

/* Empty state */
.empty-state {
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}
</style>
