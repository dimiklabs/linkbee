<template>
  <div class="container-fluid py-4">

    <!-- Header + filters -->
    <div class="d-flex align-items-start justify-content-between flex-wrap gap-3 mb-4">
      <div>
        <h4 class="fw-bold mb-1">Analytics</h4>
        <p class="text-muted small mb-0">Account-wide click performance across all your links.</p>
      </div>
      <div class="d-flex align-items-end gap-2 flex-wrap">
        <div>
          <label class="form-label fw-medium small mb-1">From</label>
          <input v-model="filterFrom" type="date" class="form-control form-control-sm" style="width:140px;" />
        </div>
        <div>
          <label class="form-label fw-medium small mb-1">To</label>
          <input v-model="filterTo" type="date" class="form-control form-control-sm" style="width:140px;" />
        </div>
        <button
          class="btn btn-primary btn-sm d-flex align-items-center gap-2"
          :disabled="loading"
          @click="load"
        >
          <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
          Apply
        </button>
        <!-- Quick presets -->
        <div class="btn-group btn-group-sm">
          <button
            v-for="p in presets"
            :key="p.label"
            class="btn btn-outline-secondary"
            :class="{ active: activePreset === p.label }"
            @click="applyPreset(p)"
          >{{ p.label }}</button>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" class="alert alert-danger d-flex align-items-center gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5m.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
      </svg>
      {{ error }}
    </div>

    <!-- Loading skeleton -->
    <template v-if="loading && !data">
      <div class="stat-grid mb-4">
        <div v-for="i in 4" :key="i" class="card p-4">
          <div class="skeleton mb-2" style="width:60%;height:1.8rem;"></div>
          <div class="skeleton" style="width:40%;height:0.85rem;"></div>
        </div>
      </div>
    </template>

    <template v-else-if="data">

      <!-- ── Stat cards ─────────────────────────────────────────────────── -->
      <div class="stat-grid mb-4">
        <div class="stat-card card">
          <div class="stat-icon">👆</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.total_clicks.toLocaleString() }}</div>
            <div class="stat-label">Total Clicks</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">👤</div>
          <div class="stat-body">
            <div class="stat-value">{{ data.unique_clicks.toLocaleString() }}</div>
            <div class="stat-label">Unique Visitors</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">🌍</div>
          <div class="stat-body">
            <div class="stat-value">{{ topCountry }}</div>
            <div class="stat-label">Top Country</div>
          </div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">📱</div>
          <div class="stat-body">
            <div class="stat-value">{{ topDevice }}</div>
            <div class="stat-label">Top Device</div>
          </div>
        </div>
      </div>

      <!-- ── Period Comparison card ──────────────────────────────────────── -->
      <div v-if="comparison" class="card border-0 shadow-sm mb-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Period Comparison</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">
              Current vs preceding period of equal length
            </p>
          </div>
          <span class="badge bg-light text-muted border fw-normal small">
            {{ formatShortDate(comparison.previous.from) }} → {{ formatShortDate(comparison.current.to) }}
          </span>
        </div>
        <div class="card-body">
          <div class="row g-3">
            <!-- Total Clicks comparison -->
            <div class="col-sm-6 col-lg-3">
              <div class="comparison-metric">
                <div class="comp-label text-muted small mb-1">Total Clicks</div>
                <div class="d-flex align-items-baseline gap-2 flex-wrap">
                  <span class="comp-current fw-bold fs-5">{{ comparison.current.total_clicks.toLocaleString() }}</span>
                  <span
                    class="comp-badge"
                    :class="trendClass(comparison.clicks.trend)"
                  >
                    {{ trendArrow(comparison.clicks.trend) }}{{ Math.abs(comparison.clicks.percent_change).toFixed(1) }}%
                  </span>
                </div>
                <div class="comp-prev text-muted" style="font-size: 0.75rem;">
                  vs {{ comparison.previous.total_clicks.toLocaleString() }} prev. period
                </div>
              </div>
            </div>

            <!-- Unique Clicks comparison -->
            <div class="col-sm-6 col-lg-3">
              <div class="comparison-metric">
                <div class="comp-label text-muted small mb-1">Unique Clicks</div>
                <div class="d-flex align-items-baseline gap-2 flex-wrap">
                  <span class="comp-current fw-bold fs-5">{{ comparison.current.unique_clicks.toLocaleString() }}</span>
                  <span
                    class="comp-badge"
                    :class="trendClass(comparison.unique_clicks.trend)"
                  >
                    {{ trendArrow(comparison.unique_clicks.trend) }}{{ Math.abs(comparison.unique_clicks.percent_change).toFixed(1) }}%
                  </span>
                </div>
                <div class="comp-prev text-muted" style="font-size: 0.75rem;">
                  vs {{ comparison.previous.unique_clicks.toLocaleString() }} prev. period
                </div>
              </div>
            </div>

            <!-- Current period summary -->
            <div class="col-sm-6 col-lg-3">
              <div class="comparison-period-box current-period">
                <div class="comp-period-label">Current Period</div>
                <div class="comp-period-dates">{{ formatShortDate(comparison.current.from) }} – {{ formatShortDate(comparison.current.to) }}</div>
                <div class="d-flex gap-3 mt-1">
                  <span class="small"><strong>{{ comparison.current.total_clicks.toLocaleString() }}</strong> clicks</span>
                  <span class="small text-muted"><strong>{{ comparison.current.unique_clicks.toLocaleString() }}</strong> unique</span>
                </div>
              </div>
            </div>

            <!-- Previous period summary -->
            <div class="col-sm-6 col-lg-3">
              <div class="comparison-period-box previous-period">
                <div class="comp-period-label">Previous Period</div>
                <div class="comp-period-dates">{{ formatShortDate(comparison.previous.from) }} – {{ formatShortDate(comparison.previous.to) }}</div>
                <div class="d-flex gap-3 mt-1">
                  <span class="small"><strong>{{ comparison.previous.total_clicks.toLocaleString() }}</strong> clicks</span>
                  <span class="small text-muted"><strong>{{ comparison.previous.unique_clicks.toLocaleString() }}</strong> unique</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Click Trend ────────────────────────────────────────────────── -->
      <div class="card border-0 shadow-sm mb-4">
        <div class="card-header bg-transparent border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <h6 class="mb-0 fw-semibold">Click Trend</h6>
          <span class="text-muted small">{{ formatDate(data.from) }} – {{ formatDate(data.to) }}</span>
        </div>
        <div class="card-body px-3 pb-3 pt-2">
          <div v-if="!data.time_series.length" class="text-center py-4 text-muted small">
            No click data for this period.
          </div>
          <VChart v-else :option="trendOption" style="height: 240px;" autoresize />
        </div>
      </div>

      <!-- ── Two-column: Countries + Devices ────────────────────────────── -->
      <div class="row g-4 mb-4">
        <!-- Top Countries -->
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Top Countries</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.top_countries.length" class="text-center text-muted small py-3">No country data yet.</div>
              <div v-else>
                <div
                  v-for="(c, idx) in data.top_countries"
                  :key="c.country"
                  class="country-row d-flex align-items-center gap-3 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flag-emoji">{{ countryFlag(c.country) }}</span>
                  <span class="flex-grow-1 small fw-medium">{{ countryName(c.country) }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:120px;">
                    <div
                      class="progress-bar"
                      :style="{ width: countryPct(c.count) + '%', backgroundColor: '#635bff' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ c.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Device Breakdown -->
        <div class="col-lg-5">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Device Breakdown</h6>
            </div>
            <div class="card-body d-flex flex-column align-items-center justify-content-center py-3">
              <div v-if="!data.device_breakdown.length" class="text-center text-muted small py-3">No device data yet.</div>
              <template v-else>
                <VChart :option="deviceOption" style="height: 200px;width:100%;" autoresize />
                <div class="d-flex flex-wrap justify-content-center gap-3 mt-1">
                  <div v-for="(d, idx) in data.device_breakdown" :key="d.device_type" class="d-flex align-items-center gap-1">
                    <span class="legend-dot" :style="{ backgroundColor: DEVICE_COLORS[idx % DEVICE_COLORS.length] }"></span>
                    <span class="small text-muted">{{ deviceLabel(d.device_type) }}</span>
                    <span class="small fw-medium">{{ d.count.toLocaleString() }}</span>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Two-column: OS + City ───────────────────────────────────────── -->
      <div class="row g-4 mb-4">
        <!-- OS Breakdown -->
        <div class="col-lg-5">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Operating Systems</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.os_breakdown.length" class="text-center text-muted small py-3">No OS data yet.</div>
              <div v-else>
                <div
                  v-for="(o, idx) in data.os_breakdown"
                  :key="o.os"
                  class="d-flex align-items-center gap-2 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flex-grow-1 small fw-medium">{{ o.os || 'Unknown' }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:100px;">
                    <div
                      class="progress-bar"
                      :style="{ width: osPct(o.count) + '%', backgroundColor: '#6366f1' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ o.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Top Cities -->
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Top Cities</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.top_cities.length" class="text-center text-muted small py-3">No city data yet.</div>
              <div v-else>
                <div
                  v-for="(c, idx) in data.top_cities"
                  :key="c.city + c.country"
                  class="d-flex align-items-center gap-3 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flag-emoji">{{ countryFlag(c.country) }}</span>
                  <span class="flex-grow-1 small fw-medium">{{ c.city }}</span>
                  <span class="text-muted small">{{ countryName(c.country) }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:100px;">
                    <div
                      class="progress-bar"
                      :style="{ width: cityPct(c.count) + '%', backgroundColor: '#635bff' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ c.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Two-column: Browsers + Top Referrers ───────────────────────── -->
      <div class="row g-4 mb-4">
        <!-- Top Browsers -->
        <div class="col-lg-5">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Top Browsers</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.top_browsers.length" class="text-center text-muted small py-3">No browser data yet.</div>
              <div v-else>
                <div
                  v-for="(b, idx) in data.top_browsers"
                  :key="b.browser"
                  class="d-flex align-items-center gap-2 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flex-grow-1 small fw-medium">{{ b.browser || 'Unknown' }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:100px;">
                    <div
                      class="progress-bar"
                      :style="{ width: browserPct(b.count) + '%', backgroundColor: '#14b8a6' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ b.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Top Referrers -->
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Top Referrers</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.top_referrers.length" class="text-center text-muted small py-3">No referrer data yet.</div>
              <div v-else>
                <div
                  v-for="(r, idx) in data.top_referrers"
                  :key="r.referrer"
                  class="d-flex align-items-center gap-3 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flex-grow-1 small fw-medium text-truncate" style="max-width:260px;" :title="r.referrer">{{ r.referrer }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:120px;">
                    <div
                      class="progress-bar"
                      :style="{ width: referrerPct(r.count) + '%', backgroundColor: '#f59e0b' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ r.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Two-column: Click Source + Referrer Categories ────────────── -->
      <div class="row g-4 mb-4">
        <!-- Click Source -->
        <div class="col-lg-5">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Click Source</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.source_breakdown.length" class="text-center text-muted small py-3">No source data yet.</div>
              <div v-else>
                <div
                  v-for="(s, idx) in data.source_breakdown"
                  :key="s.source"
                  class="d-flex align-items-center gap-2 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="me-1">{{ sourceIcon(s.source) }}</span>
                  <span class="flex-grow-1 small fw-medium text-capitalize">{{ s.source }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:100px;">
                    <div
                      class="progress-bar"
                      :style="{ width: sourcePct(s.count) + '%', backgroundColor: '#ec4899' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ s.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Referrer Categories -->
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-transparent border-bottom py-3 px-4">
              <h6 class="mb-0 fw-semibold">Traffic Channels</h6>
            </div>
            <div class="card-body px-4 py-3">
              <div v-if="!data.referrer_categories.length" class="text-center text-muted small py-3">No channel data yet.</div>
              <div v-else>
                <div
                  v-for="(rc, idx) in data.referrer_categories"
                  :key="rc.category"
                  class="d-flex align-items-center gap-3 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="me-1">{{ channelIcon(rc.category) }}</span>
                  <span class="flex-grow-1 small fw-medium text-capitalize">{{ rc.category }}</span>
                  <div class="progress flex-grow-1" style="height:6px;max-width:120px;">
                    <div
                      class="progress-bar"
                      :style="{ width: refCatPct(rc.count) + '%', backgroundColor: '#f59e0b' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:36px;text-align:right;">{{ rc.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── UTM Tracking ────────────────────────────────────────────────── -->
      <div
        v-if="hasUTMData"
        class="card border-0 shadow-sm mb-4"
      >
        <div class="card-header bg-transparent border-bottom py-3 px-4">
          <h6 class="mb-0 fw-semibold">UTM Tracking</h6>
        </div>
        <div class="card-body">
          <div class="row g-4">
            <!-- UTM Sources -->
            <div class="col-lg-4">
              <div class="utm-section-label text-muted small fw-semibold mb-2 text-uppercase" style="letter-spacing:0.04em;font-size:0.7rem;">UTM Source</div>
              <div v-if="!data.utm_sources.length" class="text-muted small">No UTM source data.</div>
              <div v-else>
                <div
                  v-for="(u, idx) in data.utm_sources"
                  :key="u.value"
                  class="d-flex align-items-center gap-2 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flex-grow-1 small fw-medium text-truncate" :title="u.value">{{ u.value || '—' }}</span>
                  <div class="progress flex-grow-1" style="height:5px;max-width:80px;">
                    <div
                      class="progress-bar"
                      :style="{ width: utmSrcPct(u.count) + '%', backgroundColor: '#635bff' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:30px;text-align:right;">{{ u.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>

            <!-- UTM Mediums -->
            <div class="col-lg-4">
              <div class="utm-section-label text-muted small fw-semibold mb-2 text-uppercase" style="letter-spacing:0.04em;font-size:0.7rem;">UTM Medium</div>
              <div v-if="!data.utm_mediums.length" class="text-muted small">No UTM medium data.</div>
              <div v-else>
                <div
                  v-for="(u, idx) in data.utm_mediums"
                  :key="u.value"
                  class="d-flex align-items-center gap-2 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flex-grow-1 small fw-medium text-truncate" :title="u.value">{{ u.value || '—' }}</span>
                  <div class="progress flex-grow-1" style="height:5px;max-width:80px;">
                    <div
                      class="progress-bar"
                      :style="{ width: utmMedPct(u.count) + '%', backgroundColor: '#14b8a6' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:30px;text-align:right;">{{ u.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>

            <!-- UTM Campaigns -->
            <div class="col-lg-4">
              <div class="utm-section-label text-muted small fw-semibold mb-2 text-uppercase" style="letter-spacing:0.04em;font-size:0.7rem;">UTM Campaign</div>
              <div v-if="!data.utm_campaigns.length" class="text-muted small">No UTM campaign data.</div>
              <div v-else>
                <div
                  v-for="(u, idx) in data.utm_campaigns"
                  :key="u.value"
                  class="d-flex align-items-center gap-2 mb-2"
                >
                  <span class="rank-num text-muted">{{ idx + 1 }}</span>
                  <span class="flex-grow-1 small fw-medium text-truncate" :title="u.value">{{ u.value || '—' }}</span>
                  <div class="progress flex-grow-1" style="height:5px;max-width:80px;">
                    <div
                      class="progress-bar"
                      :style="{ width: utmCamPct(u.count) + '%', backgroundColor: '#f59e0b' }"
                    ></div>
                  </div>
                  <span class="text-muted small" style="min-width:30px;text-align:right;">{{ u.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </template>

    <!-- Empty state (no data, not loading) -->
    <div v-else-if="!loading && !error" class="text-center py-5 text-muted">
      <div style="font-size:3rem;line-height:1;margin-bottom:0.75rem;">📊</div>
      <p class="mb-0">No analytics data yet. Clicks will appear here as your links are visited.</p>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { use } from 'echarts/core';
import { LineChart, PieChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import VChart from 'vue-echarts';
import dashboardApi from '@/api/dashboard';
import type { GlobalAnalyticsResponse, GlobalAnalyticsComparisonResponse } from '@/types/dashboard';

use([LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer]);

const DEVICE_COLORS = ['#635bff', '#14b8a6', '#f59e0b', '#ec4899', '#6366f1'];

// ── Date helpers ──────────────────────────────────────────────────────────────
function dateStr(d: Date) { return d.toISOString().slice(0, 10); }

const now = new Date();
const presets = [
  { label: '7d',  days: 7 },
  { label: '30d', days: 30 },
  { label: '90d', days: 90 },
];
const activePreset = ref('30d');
const filterFrom = ref(dateStr(new Date(Date.now() - 30 * 86400000)));
const filterTo   = ref(dateStr(now));

function applyPreset(p: { label: string; days: number }) {
  activePreset.value = p.label;
  filterFrom.value = dateStr(new Date(Date.now() - p.days * 86400000));
  filterTo.value   = dateStr(new Date());
  load();
}

// ── State ─────────────────────────────────────────────────────────────────────
const loading    = ref(false);
const error      = ref('');
const data       = ref<GlobalAnalyticsResponse | null>(null);
const comparison = ref<GlobalAnalyticsComparisonResponse | null>(null);

async function load() {
  loading.value = true;
  error.value   = '';
  try {
    const from = filterFrom.value ? `${filterFrom.value}T00:00:00Z` : undefined;
    const to   = filterTo.value   ? `${filterTo.value}T23:59:59Z`   : undefined;

    const [mainRes, compRes] = await Promise.allSettled([
      dashboardApi.getGlobalAnalytics(from, to),
      dashboardApi.getGlobalAnalyticsComparison(from, to),
    ]);

    if (mainRes.status === 'fulfilled') {
      data.value = mainRes.value.data ?? null;
    } else {
      error.value = 'Failed to load analytics. Please try again.';
    }

    if (compRes.status === 'fulfilled') {
      comparison.value = compRes.value.data ?? null;
    }
  } catch {
    error.value = 'Failed to load analytics. Please try again.';
  } finally {
    loading.value = false;
  }
}

onMounted(load);

// ── Computed helpers ──────────────────────────────────────────────────────────
const topCountry = computed(() => {
  const c = data.value?.top_countries[0];
  return c ? countryName(c.country) : '—';
});

const topDevice = computed(() => {
  const d = data.value?.device_breakdown[0];
  return d ? deviceLabel(d.device_type) : '—';
});

const hasUTMData = computed(() => {
  if (!data.value) return false;
  return (
    (data.value.utm_sources?.length ?? 0) > 0 ||
    (data.value.utm_mediums?.length ?? 0) > 0 ||
    (data.value.utm_campaigns?.length ?? 0) > 0
  );
});

// Max counts for progress bar calculations
const totalCountryClicks = computed(() =>
  (data.value?.top_countries ?? []).reduce((s, c) => s + c.count, 0),
);
const totalBrowserClicks = computed(() =>
  (data.value?.top_browsers ?? []).reduce((s, b) => s + b.count, 0),
);
const totalReferrerClicks = computed(() =>
  (data.value?.top_referrers ?? []).reduce((s, r) => s + r.count, 0),
);
const totalOSClicks = computed(() =>
  (data.value?.os_breakdown ?? []).reduce((s, o) => s + o.count, 0),
);
const totalCityClicks = computed(() =>
  (data.value?.top_cities ?? []).reduce((s, c) => s + c.count, 0),
);
const totalSourceClicks = computed(() =>
  (data.value?.source_breakdown ?? []).reduce((s, sc) => s + sc.count, 0),
);
const totalRefCatClicks = computed(() =>
  (data.value?.referrer_categories ?? []).reduce((s, rc) => s + rc.count, 0),
);
const totalUTMSrcClicks = computed(() =>
  (data.value?.utm_sources ?? []).reduce((s, u) => s + u.count, 0),
);
const totalUTMMedClicks = computed(() =>
  (data.value?.utm_mediums ?? []).reduce((s, u) => s + u.count, 0),
);
const totalUTMCamClicks = computed(() =>
  (data.value?.utm_campaigns ?? []).reduce((s, u) => s + u.count, 0),
);

function countryPct(n: number)  { return totalCountryClicks.value  ? Math.round((n / totalCountryClicks.value)  * 100) : 0; }
function browserPct(n: number)  { return totalBrowserClicks.value  ? Math.round((n / totalBrowserClicks.value)  * 100) : 0; }
function referrerPct(n: number) { return totalReferrerClicks.value ? Math.round((n / totalReferrerClicks.value) * 100) : 0; }
function osPct(n: number)       { return totalOSClicks.value       ? Math.round((n / totalOSClicks.value)       * 100) : 0; }
function cityPct(n: number)     { return totalCityClicks.value     ? Math.round((n / totalCityClicks.value)     * 100) : 0; }
function sourcePct(n: number)   { return totalSourceClicks.value   ? Math.round((n / totalSourceClicks.value)   * 100) : 0; }
function refCatPct(n: number)   { return totalRefCatClicks.value   ? Math.round((n / totalRefCatClicks.value)   * 100) : 0; }
function utmSrcPct(n: number)   { return totalUTMSrcClicks.value   ? Math.round((n / totalUTMSrcClicks.value)   * 100) : 0; }
function utmMedPct(n: number)   { return totalUTMMedClicks.value   ? Math.round((n / totalUTMMedClicks.value)   * 100) : 0; }
function utmCamPct(n: number)   { return totalUTMCamClicks.value   ? Math.round((n / totalUTMCamClicks.value)   * 100) : 0; }

function deviceLabel(t: string): string {
  const map: Record<string, string> = { desktop: 'Desktop', mobile: 'Mobile', tablet: 'Tablet', bot: 'Bot', unknown: 'Unknown', '': 'Unknown' };
  return map[t] ?? t;
}

function countryName(iso: string): string {
  try { return new Intl.DisplayNames(['en'], { type: 'region' }).of(iso) ?? iso; } catch { return iso; }
}

function countryFlag(iso: string): string {
  if (!iso || iso.length !== 2) return '🌐';
  return String.fromCodePoint(...iso.toUpperCase().split('').map((c) => c.charCodeAt(0) + 127397));
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

function formatShortDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

function sourceIcon(source: string): string {
  const icons: Record<string, string> = {
    qr: '📷',
    direct: '🔗',
    web: '🌐',
    api: '⚡',
  };
  return icons[source.toLowerCase()] ?? '🌐';
}

function channelIcon(category: string): string {
  const icons: Record<string, string> = {
    direct: '🔗',
    search: '🔍',
    social: '📱',
    email: '📧',
    referral: '🌐',
  };
  return icons[category.toLowerCase()] ?? '🌐';
}

// ── Period comparison helpers ─────────────────────────────────────────────────
function trendClass(trend: string): string {
  if (trend === 'up')   return 'comp-badge--up';
  if (trend === 'down') return 'comp-badge--down';
  return 'comp-badge--flat';
}

function trendArrow(trend: string): string {
  if (trend === 'up')   return '↑ ';
  if (trend === 'down') return '↓ ';
  return '→ ';
}

// ── Chart options ─────────────────────────────────────────────────────────────
const trendOption = computed(() => {
  const series = data.value?.time_series ?? [];
  const dates  = series.map((p) => new Date(p.timestamp).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }));
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
    series: [{
      type: 'line',
      data: counts,
      smooth: true,
      symbol: 'none',
      lineStyle: { color: '#635bff', width: 2 },
      areaStyle: {
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(99,91,255,0.18)' },
            { offset: 1, color: 'rgba(99,91,255,0.01)' },
          ],
        },
      },
    }],
  };
});

const deviceOption = computed(() => {
  const devices = data.value?.device_breakdown ?? [];
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { show: false },
    series: [{
      type: 'pie',
      radius: ['48%', '74%'],
      avoidLabelOverlap: false,
      label: { show: false },
      data: devices.map((d, idx) => ({
        name: deviceLabel(d.device_type),
        value: d.count,
        itemStyle: { color: DEVICE_COLORS[idx % DEVICE_COLORS.length] },
      })),
    }],
  };
});
</script>

<style scoped lang="scss">
$primary: #635bff;

.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;

  @media (max-width: 991px) { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;

  .stat-icon  { font-size: 2rem; line-height: 1; flex-shrink: 0; }
  .stat-value { font-size: 1.75rem; font-weight: 700; color: #1a1f36; line-height: 1; }
  .stat-label { font-size: 0.8125rem; color: #697386; margin-top: 0.25rem; }
}

.rank-num {
  font-size: 0.75rem;
  font-weight: 600;
  min-width: 16px;
  text-align: center;
}

.flag-emoji { font-size: 1.1rem; line-height: 1; }

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.skeleton {
  background: linear-gradient(90deg, #e3e8ee 25%, #f0f2f5 50%, #e3e8ee 75%);
  background-size: 200% 100%;
  animation: shimmer 1.2s infinite;
  border-radius: 4px;
  display: block;
}
@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

// ── Period comparison ──────────────────────────────────────────────────────────
.comparison-metric {
  padding: 0.25rem 0;
}

.comp-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.78rem;
  font-weight: 600;
  padding: 0.15rem 0.45rem;
  border-radius: 999px;

  &--up   { background-color: rgba(22, 163, 74, 0.1);  color: #16a34a; }
  &--down { background-color: rgba(220, 38,  38, 0.1);  color: #dc2626; }
  &--flat { background-color: rgba(107, 114, 128, 0.1); color: #6b7280; }
}

.comparison-period-box {
  border-radius: 8px;
  padding: 0.75rem 1rem;

  .comp-period-label {
    font-size: 0.72rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    margin-bottom: 0.2rem;
  }

  .comp-period-dates {
    font-size: 0.8rem;
    color: #697386;
    margin-bottom: 0.25rem;
  }

  &.current-period {
    background-color: rgba(99, 91, 255, 0.06);
    border: 1px solid rgba(99, 91, 255, 0.2);
    .comp-period-label { color: $primary; }
  }

  &.previous-period {
    background-color: rgba(107, 114, 128, 0.06);
    border: 1px solid rgba(107, 114, 128, 0.2);
    .comp-period-label { color: #6b7280; }
  }
}

// ── Dark mode ──────────────────────────────────────────────────────────────────
:global(.dark-mode) .stat-card .stat-value { color: #e6edf3; }
:global(.dark-mode) .skeleton { background: linear-gradient(90deg, #21262d 25%, #30363d 50%, #21262d 75%); background-size: 200% 100%; }
:global(.dark-mode) .comparison-period-box.current-period  { background-color: rgba(99, 91, 255, 0.12); border-color: rgba(99, 91, 255, 0.3); }
:global(.dark-mode) .comparison-period-box.previous-period { background-color: rgba(107, 114, 128, 0.12); border-color: rgba(107, 114, 128, 0.3); }
:global(.dark-mode) .card-header.bg-white { background-color: transparent !important; }

.btn-outline-secondary.active {
  background-color: $primary;
  border-color: $primary;
  color: #fff;
}

.btn-primary {
  background-color: $primary;
  border-color: $primary;
  &:hover:not(:disabled) { background-color: #5249e0; border-color: #5249e0; }
}
</style>
