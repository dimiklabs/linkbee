<template>
  <div class="container-fluid py-4">
    <!-- Page title -->
    <div class="mb-4">
      <h4 class="fw-bold mb-1">Compare Links</h4>
      <p class="text-muted small mb-0">Select 2–5 links to compare their performance side by side.</p>
    </div>

    <!-- Filters card -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-body">
        <div class="row g-3 align-items-end">
          <!-- Link selector -->
          <div class="col-12 col-md-6">
            <label class="form-label fw-medium small">Links <span class="text-muted">(select 2–5)</span></label>
            <div class="link-selector" :class="{ open: selectorOpen }" ref="selectorRef">
              <div
                class="selector-trigger form-control form-control-sm d-flex align-items-center justify-content-between"
                @click="selectorOpen = !selectorOpen"
              >
                <span v-if="selectedIds.length === 0" class="text-muted">Choose links to compare…</span>
                <span v-else class="text-truncate">{{ selectedIds.length }} link{{ selectedIds.length > 1 ? 's' : '' }} selected</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" viewBox="0 0 16 16" class="flex-shrink-0 ms-2">
                  <path fill-rule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
                </svg>
              </div>
              <div v-if="selectorOpen" class="selector-dropdown card shadow-sm border">
                <div class="selector-search px-2 pt-2">
                  <input
                    v-model="linkSearch"
                    class="form-control form-control-sm"
                    placeholder="Search links…"
                    @click.stop
                  />
                </div>
                <div class="selector-options">
                  <label
                    v-for="link in filteredLinks"
                    :key="link.id"
                    class="selector-option d-flex align-items-center gap-2 px-3 py-2"
                    :class="{ disabled: !selectedIds.includes(link.id) && selectedIds.length >= 5 }"
                    @click.stop
                  >
                    <input
                      type="checkbox"
                      :value="link.id"
                      v-model="selectedIds"
                      :disabled="!selectedIds.includes(link.id) && selectedIds.length >= 5"
                      class="form-check-input flex-shrink-0"
                    />
                    <div class="min-w-0">
                      <div class="small fw-medium text-truncate">{{ link.title || link.slug }}</div>
                      <div class="text-muted" style="font-size:0.7rem;">{{ link.short_url }}</div>
                    </div>
                  </label>
                  <div v-if="filteredLinks.length === 0" class="text-muted small text-center py-3">No links found</div>
                </div>
              </div>
            </div>
          </div>

          <!-- From date -->
          <div class="col-sm-6 col-md-2">
            <label class="form-label fw-medium small">From</label>
            <input v-model="filterFrom" type="date" class="form-control form-control-sm" />
          </div>

          <!-- To date -->
          <div class="col-sm-6 col-md-2">
            <label class="form-label fw-medium small">To</label>
            <input v-model="filterTo" type="date" class="form-control form-control-sm" />
          </div>

          <!-- Actions -->
          <div class="col-12 col-md-2">
            <div class="d-flex gap-2">
              <button
                class="btn btn-primary btn-sm flex-grow-1 d-flex align-items-center justify-content-center gap-2"
                :disabled="loading || selectedIds.length < 2"
                @click="runComparison"
              >
                <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                Compare
              </button>
              <button
                v-if="result"
                class="btn btn-outline-secondary btn-sm flex-shrink-0"
                title="Export to CSV"
                @click="exportCSV"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
                  <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Validation hint -->
        <div v-if="selectedIds.length === 1" class="mt-2">
          <small class="text-warning">Select at least one more link to run a comparison.</small>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="alert alert-danger d-flex align-items-center gap-2" role="alert">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5m.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
      </svg>
      {{ error }}
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" style="width: 2.5rem; height: 2.5rem;" role="status">
        <span class="visually-hidden">Loading comparison…</span>
      </div>
      <p class="text-muted mt-3 mb-0">Fetching comparison data…</p>
    </div>

    <template v-if="!loading && result">
      <!-- Period info -->
      <div class="mb-3 text-muted small">
        Showing data from <strong>{{ formatDate(result.from) }}</strong> to <strong>{{ formatDate(result.to) }}</strong>
        &nbsp;·&nbsp; <strong>{{ result.span_days }}</strong> day{{ result.span_days !== 1 ? 's' : '' }}
      </div>

      <!-- Bar chart — total clicks -->
      <div class="card border-0 shadow-sm mb-4">
        <div class="card-header bg-transparent border-0 pt-3 pb-0">
          <h6 class="mb-0 fw-semibold">Total Clicks Comparison</h6>
        </div>
        <div class="card-body">
          <VChart :option="barChartOption" style="height: 260px;" autoresize />
        </div>
      </div>

      <!-- Comparison table -->
      <div class="card border-0 shadow-sm">
        <div class="card-header bg-transparent border-0 pt-3 pb-0">
          <h6 class="mb-0 fw-semibold">Detailed Breakdown</h6>
        </div>
        <div class="card-body p-0">
          <div class="table-responsive">
            <table class="table table-hover align-middle mb-0 comparison-table">
              <thead class="table-light">
                <tr>
                  <th class="ps-3">Link</th>
                  <th class="text-end">Total Clicks</th>
                  <th class="text-end">Unique Clicks</th>
                  <th class="text-end">Clicks / Day</th>
                  <th>Top Referrer</th>
                  <th>Top Country</th>
                  <th>Top Browser</th>
                  <th class="pe-3">Top Device</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(link, idx) in result.links" :key="link.link_id">
                  <td class="ps-3">
                    <div class="d-flex align-items-center gap-2">
                      <span class="rank-badge" :style="{ background: chartColors[idx % chartColors.length] }">{{ idx + 1 }}</span>
                      <div class="min-w-0">
                        <div class="fw-medium small text-truncate" style="max-width: 180px;" :title="link.title || link.slug">
                          {{ link.title || link.slug }}
                        </div>
                        <div class="text-muted" style="font-size: 0.7rem;">{{ link.slug }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="text-end fw-semibold">{{ link.total_clicks.toLocaleString() }}</td>
                  <td class="text-end">{{ link.unique_clicks.toLocaleString() }}</td>
                  <td class="text-end">{{ link.clicks_per_day.toFixed(1) }}</td>
                  <td>
                    <span v-if="link.top_referrer" class="badge-pill-text">{{ link.top_referrer }}</span>
                    <span v-else class="text-muted small">—</span>
                  </td>
                  <td>
                    <span v-if="link.top_country" class="badge-pill-text">{{ link.top_country }}</span>
                    <span v-else class="text-muted small">—</span>
                  </td>
                  <td>
                    <span v-if="link.top_browser" class="badge-pill-text">{{ link.top_browser }}</span>
                    <span v-else class="text-muted small">—</span>
                  </td>
                  <td class="pe-3">
                    <span v-if="link.top_device" class="badge-pill-text">{{ link.top_device }}</span>
                    <span v-else class="text-muted small">—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </template>

    <!-- Empty state (no run yet) -->
    <div v-if="!loading && !result && !error" class="text-center py-5 text-muted">
      <div style="font-size: 3rem; line-height: 1; margin-bottom: 0.75rem;">📊</div>
      <p class="mb-0">Select 2–5 links above and click <strong>Compare</strong> to see their performance side by side.</p>
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
$primary: #635bff;

// ── Link selector ─────────────────────────────────────────────────────────────
.link-selector {
  position: relative;
}

.selector-trigger {
  cursor: pointer;
  user-select: none;
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 400;
  border-radius: 8px;
  overflow: hidden;
}

.selector-options {
  max-height: 220px;
  overflow-y: auto;
}

.selector-option {
  cursor: pointer;
  font-size: 0.8125rem;
  transition: background 0.12s;

  &:hover {
    background: #f7f9fc;
  }

  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

// ── Comparison table ──────────────────────────────────────────────────────────
.comparison-table {
  font-size: 0.8125rem;

  th {
    font-size: 0.75rem;
    font-weight: 600;
    letter-spacing: 0.02em;
    color: #697386;
    white-space: nowrap;
  }
}

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

.badge-pill-text {
  display: inline-block;
  background: #eef0ff;
  color: $primary;
  border-radius: 999px;
  padding: 0.15em 0.6em;
  font-size: 0.7rem;
  font-weight: 500;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.min-w-0 {
  min-width: 0;
}

// ── Dark mode ─────────────────────────────────────────────────────────────────
:global(.dark-mode) {
  .selector-dropdown {
    background: #161b22;
    border-color: #30363d !important;
  }

  .selector-option:hover {
    background: #21262d;
  }

  .selector-trigger {
    background: #161b22;
    border-color: #30363d;
    color: #e6edf3;
  }

  .comparison-table {
    th {
      color: #8b949e;
    }

    td {
      color: #e6edf3;
      border-color: #30363d;
    }
  }

  .badge-pill-text {
    background: rgba(99, 91, 255, 0.2);
  }

  .card {
    background: #161b22;
    border-color: #30363d !important;
  }

  .card-header {
    border-color: #30363d !important;
  }

  .table-light {
    --bs-table-bg: #21262d;
    --bs-table-color: #8b949e;
  }

  .table-hover > tbody > tr:hover > * {
    background: #21262d;
  }
}
</style>
