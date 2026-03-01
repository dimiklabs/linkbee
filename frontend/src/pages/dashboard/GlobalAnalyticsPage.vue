<template>
  <UpgradeWall
    v-if="!billingStore.isPaidPlan && billingStore.loaded"
    icon="bar_chart"
    title="Analytics is a Pro feature"
    description="Upgrade to Pro to unlock account-wide analytics, time-series charts, and geographic breakdowns."
  />
  <div v-else class="global-analytics-page">

    <!-- Page Header + Filters -->
    <div class="page-header">
      <div class="page-header__title-block">
        <h1 class="md-headline-small page-title">Analytics</h1>
        <p class="md-body-medium page-subtitle">Account-wide click performance across all your links.</p>
      </div>
      <div class="page-header__filters">
        <!-- Period presets chip group -->
        <div class="chip-group" role="group" aria-label="Period presets">
          <button
            v-for="p in presets"
            :key="p.label"
            class="chip"
            :class="{ active: activePreset === p.label }"
            @click="applyPreset(p)"
          >
            {{ p.label }}
          </button>
        </div>

        <!-- Date range inputs -->
        <div class="date-range-row">
          <label class="date-field">
            <span class="date-field__label">From</span>
            <input type="date" class="form-input" :value="filterFrom" @input="filterFrom = ($event.target as HTMLInputElement).value" />
          </label>
          <span class="date-range-sep">–</span>
          <label class="date-field">
            <span class="date-field__label">To</span>
            <input type="date" class="form-input" :value="filterTo" @input="filterTo = ($event.target as HTMLInputElement).value" />
          </label>
          <button class="btn-filled" :disabled="loading" @click="load">
            <span v-if="loading" class="css-spinner css-spinner--sm css-spinner--white"></span>
            <span v-else class="material-symbols-outlined">refresh</span>
            Apply
          </button>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" class="error-banner">
      <span class="material-symbols-outlined" style="color:var(--md-sys-color-error);">error</span>
      <span class="md-body-medium" style="flex:1;">{{ error }}</span>
    </div>

    <!-- Loading skeleton -->
    <template v-if="loading && !data">
      <div class="stat-grid">
        <div v-for="i in 4" :key="i" class="m3-card m3-card--elevated stat-skeleton-card">
          <div class="skeleton-icon-box skeleton"></div>
          <div style="flex:1">
            <div class="skeleton" style="width:55%;height:1.6rem;border-radius:4px;margin-bottom:0.4rem;"></div>
            <div class="skeleton" style="width:38%;height:0.8rem;border-radius:4px;"></div>
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="data">

      <!-- Stat Cards -->
      <div class="stat-grid">
        <div class="m3-card m3-card--elevated stat-card" aria-label="Total clicks" role="region">
          <div class="stat-card__icon stat-card__icon--primary">
            <span class="material-symbols-outlined">ads_click</span>
          </div>
          <div class="stat-card__body">
            <div class="md-headline-small stat-value">{{ data.total_clicks.toLocaleString() }}</div>
            <div class="md-body-small stat-label">Total Clicks</div>
          </div>
        </div>
        <div class="m3-card m3-card--elevated stat-card" aria-label="Unique visitors" role="region">
          <div class="stat-card__icon stat-card__icon--teal">
            <span class="material-symbols-outlined">person</span>
          </div>
          <div class="stat-card__body">
            <div class="md-headline-small stat-value">{{ data.unique_clicks.toLocaleString() }}</div>
            <div class="md-body-small stat-label">Unique Visitors</div>
          </div>
        </div>
        <div class="m3-card m3-card--elevated stat-card" aria-label="Top country" role="region">
          <div class="stat-card__icon stat-card__icon--amber">
            <span class="material-symbols-outlined">public</span>
          </div>
          <div class="stat-card__body">
            <div class="md-title-medium stat-value">{{ topCountry }}</div>
            <div class="md-body-small stat-label">Top Country</div>
          </div>
        </div>
        <div class="m3-card m3-card--elevated stat-card" aria-label="Top device" role="region">
          <div class="stat-card__icon stat-card__icon--secondary">
            <span class="material-symbols-outlined">smartphone</span>
          </div>
          <div class="stat-card__body">
            <div class="md-title-medium stat-value" style="text-transform:capitalize;">{{ topDevice }}</div>
            <div class="md-body-small stat-label">Top Device</div>
          </div>
        </div>
      </div>

      <!-- Period Comparison Card -->
      <div v-if="comparison" class="m3-card m3-card--outlined">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Period Comparison</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Current vs preceding period of equal length</div>
          </div>
          <span class="m3-badge m3-badge--neutral">
            {{ formatShortDate(comparison.previous.from) }} → {{ formatShortDate(comparison.current.to) }}
          </span>
        </div>
        <div style="padding:1rem 1.25rem;">
          <div class="comparison-grid">
            <div class="comparison-metric">
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.25rem;">Total Clicks</div>
              <div style="display:flex;align-items:baseline;gap:0.5rem;flex-wrap:wrap;">
                <span class="md-headline-small">{{ comparison.current.total_clicks.toLocaleString() }}</span>
                <span :style="{ color: comparison.clicks.trend === 'up' ? 'var(--md-sys-color-tertiary,#16a34a)' : 'var(--md-sys-color-error)' }" class="trend-indicator">
                  <span class="material-symbols-outlined" style="font-size:16px;vertical-align:middle">{{ comparison.clicks.trend === 'up' ? 'trending_up' : 'trending_down' }}</span>
                  {{ Math.abs(comparison.clicks.percent_change).toFixed(1) }}%
                </span>
              </div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:0.2rem;">vs {{ comparison.previous.total_clicks.toLocaleString() }} prev.</div>
            </div>
            <div class="comparison-metric">
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.25rem;">Unique Clicks</div>
              <div style="display:flex;align-items:baseline;gap:0.5rem;flex-wrap:wrap;">
                <span class="md-headline-small">{{ comparison.current.unique_clicks.toLocaleString() }}</span>
                <span :style="{ color: comparison.unique_clicks.trend === 'up' ? 'var(--md-sys-color-tertiary,#16a34a)' : 'var(--md-sys-color-error)' }" class="trend-indicator">
                  <span class="material-symbols-outlined" style="font-size:16px;vertical-align:middle">{{ comparison.unique_clicks.trend === 'up' ? 'trending_up' : 'trending_down' }}</span>
                  {{ Math.abs(comparison.unique_clicks.percent_change).toFixed(1) }}%
                </span>
              </div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:0.2rem;">vs {{ comparison.previous.unique_clicks.toLocaleString() }} prev.</div>
            </div>
            <div class="comparison-period-box comparison-period-box--current">
              <div class="md-label-large" style="font-size:0.7rem;text-transform:uppercase;letter-spacing:0.06em;color:var(--md-sys-color-primary);margin-bottom:0.2rem;">Current Period</div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:0.25rem;">{{ formatShortDate(comparison.current.from) }} – {{ formatShortDate(comparison.current.to) }}</div>
              <div style="display:flex;gap:1rem;">
                <span class="md-body-small"><strong>{{ comparison.current.total_clicks.toLocaleString() }}</strong> clicks</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);"><strong>{{ comparison.current.unique_clicks.toLocaleString() }}</strong> unique</span>
              </div>
            </div>
            <div class="comparison-period-box comparison-period-box--previous">
              <div class="md-label-large" style="font-size:0.7rem;text-transform:uppercase;letter-spacing:0.06em;color:var(--md-sys-color-on-surface-variant);margin-bottom:0.2rem;">Previous Period</div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:0.25rem;">{{ formatShortDate(comparison.previous.from) }} – {{ formatShortDate(comparison.previous.to) }}</div>
              <div style="display:flex;gap:1rem;">
                <span class="md-body-small"><strong>{{ comparison.previous.total_clicks.toLocaleString() }}</strong> clicks</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);"><strong>{{ comparison.previous.unique_clicks.toLocaleString() }}</strong> unique</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Trend Chart -->
      <div class="m3-card m3-card--elevated chart-section">
        <div class="card-header-row">
          <div class="card-header-row__left">
            <span class="material-symbols-outlined card-header-row__icon">show_chart</span>
            <div>
              <span class="md-title-medium">Click Trend</span>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:2px;">{{ formatDate(data.from) }} – {{ formatDate(data.to) }}</div>
            </div>
          </div>
        </div>
        <div class="chart-body-pad">
          <div v-if="!data.time_series.length" class="m3-empty-state m3-empty-state--compact">
            <span class="material-symbols-outlined m3-empty-state__icon">show_chart</span>
            <p class="md-body-medium m3-empty-state__text">No click data for this period.</p>
          </div>
          <VChart v-else :option="trendOption" style="height:280px;" autoresize />
        </div>
      </div>

      <!-- Two-column: Countries + Devices -->
      <div class="two-col-grid">
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Top Countries</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.top_countries.length" class="empty-state">No country data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(c, idx) in data.top_countries" :key="c.country" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="flag-emoji">{{ countryFlag(c.country) }}</span>
                <span class="md-body-medium" style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ countryName(c.country) }}</span>
                <div class="mini-bar"><div class="mini-bar__fill" :style="{ width: countryPct(c.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ c.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Device Breakdown</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.device_breakdown.length" class="empty-state">No device data yet.</div>
            <template v-else>
              <VChart :option="deviceOption" style="height:200px;width:100%;" autoresize />
              <div style="display:flex;flex-wrap:wrap;justify-content:center;gap:0.75rem;margin-top:0.5rem;">
                <div v-for="(d, idx) in data.device_breakdown" :key="d.device_type" style="display:flex;align-items:center;gap:0.35rem;">
                  <span class="legend-dot" :style="{ backgroundColor: DEVICE_COLORS[idx % DEVICE_COLORS.length] }"></span>
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ deviceLabel(d.device_type) }}</span>
                  <span class="md-body-small" style="font-weight:600;">{{ d.count.toLocaleString() }}</span>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- Two-column: OS + Top Cities -->
      <div class="two-col-grid">
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Operating Systems</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.os_breakdown.length" class="empty-state">No OS data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(o, idx) in data.os_breakdown" :key="o.os" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="md-body-medium" style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ o.os || 'Unknown' }}</span>
                <div class="mini-bar"><div class="mini-bar__fill" :style="{ width: osPct(o.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ o.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Top Cities</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.top_cities.length" class="empty-state">No city data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(c, idx) in data.top_cities" :key="c.city + c.country" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="flag-emoji">{{ countryFlag(c.country) }}</span>
                <span class="md-body-medium" style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ c.city }}</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ countryName(c.country) }}</span>
                <div class="mini-bar mini-bar--sm"><div class="mini-bar__fill" :style="{ width: cityPct(c.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ c.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Two-column: Browsers + Referrers -->
      <div class="two-col-grid">
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Top Browsers</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.top_browsers.length" class="empty-state">No browser data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(b, idx) in data.top_browsers" :key="b.browser" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="md-body-medium" style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ b.browser || 'Unknown' }}</span>
                <div class="mini-bar"><div class="mini-bar__fill" :style="{ width: browserPct(b.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ b.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Top Referrers</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.top_referrers.length" class="empty-state">No referrer data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(r, idx) in data.top_referrers" :key="r.referrer" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="md-body-medium" style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:260px;" :title="r.referrer">{{ r.referrer }}</span>
                <div class="mini-bar"><div class="mini-bar__fill" :style="{ width: referrerPct(r.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ r.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Two-column: Click Source + Traffic Channels -->
      <div class="two-col-grid">
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Click Source</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.source_breakdown.length" class="empty-state">No source data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(s, idx) in data.source_breakdown" :key="s.source" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="md-body-small">{{ sourceIcon(s.source) }}</span>
                <span class="md-body-medium" style="flex:1;text-transform:capitalize;">{{ s.source }}</span>
                <div class="mini-bar"><div class="mini-bar__fill" :style="{ width: sourcePct(s.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ s.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Traffic Channels</span>
          </div>
          <div style="padding:0.75rem 1.25rem;">
            <div v-if="!data.referrer_categories.length" class="empty-state">No channel data yet.</div>
            <div v-else class="breakdown-list">
              <div v-for="(rc, idx) in data.referrer_categories" :key="rc.category" class="country-item">
                <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                <span class="md-body-small">{{ channelIcon(rc.category) }}</span>
                <span class="md-body-medium" style="flex:1;text-transform:capitalize;">{{ rc.category }}</span>
                <div class="mini-bar"><div class="mini-bar__fill" :style="{ width: refCatPct(rc.count) + '%' }"></div></div>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:36px;text-align:right;">{{ rc.count.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- UTM Tracking -->
      <div v-if="hasUTMData" class="m3-card m3-card--outlined">
        <div class="card-header-row">
          <span class="md-title-medium">UTM Tracking</span>
        </div>
        <div style="padding:1rem 1.25rem;">
          <div class="utm-grid">
            <div>
              <div class="utm-section-label">UTM Source</div>
              <div v-if="!data.utm_sources.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No UTM source data.</div>
              <div v-else class="breakdown-list">
                <div v-for="(u, idx) in data.utm_sources" :key="u.value" class="country-item">
                  <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                  <span class="md-body-medium" style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="u.value">{{ u.value || '—' }}</span>
                  <div class="mini-bar mini-bar--sm"><div class="mini-bar__fill" :style="{ width: utmSrcPct(u.count) + '%' }"></div></div>
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:30px;text-align:right;">{{ u.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
            <div>
              <div class="utm-section-label">UTM Medium</div>
              <div v-if="!data.utm_mediums.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No UTM medium data.</div>
              <div v-else class="breakdown-list">
                <div v-for="(u, idx) in data.utm_mediums" :key="u.value" class="country-item">
                  <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                  <span class="md-body-medium" style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="u.value">{{ u.value || '—' }}</span>
                  <div class="mini-bar mini-bar--sm"><div class="mini-bar__fill" :style="{ width: utmMedPct(u.count) + '%' }"></div></div>
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:30px;text-align:right;">{{ u.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
            <div>
              <div class="utm-section-label">UTM Campaign</div>
              <div v-if="!data.utm_campaigns.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No UTM campaign data.</div>
              <div v-else class="breakdown-list">
                <div v-for="(u, idx) in data.utm_campaigns" :key="u.value" class="country-item">
                  <span class="rank-num md-body-small">{{ idx + 1 }}</span>
                  <span class="md-body-medium" style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="u.value">{{ u.value || '—' }}</span>
                  <div class="mini-bar mini-bar--sm"><div class="mini-bar__fill" :style="{ width: utmCamPct(u.count) + '%' }"></div></div>
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:30px;text-align:right;">{{ u.count.toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </template>

    <!-- Empty state -->
    <div v-else-if="!loading && !error" class="m3-empty-state m3-empty-state--full">
      <div class="m3-empty-state__icon-wrap">
        <span class="material-symbols-outlined m3-empty-state__icon">bar_chart</span>
      </div>
      <h2 class="md-title-large m3-empty-state__title">No analytics data yet</h2>
      <p class="md-body-medium m3-empty-state__text">Clicks will appear here as your links are visited.</p>
      <router-link to="/dashboard/links">
        <button class="btn-filled">
          <span class="material-symbols-outlined">add</span>
          Create a link
        </button>
      </router-link>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { use } from 'echarts/core';
import UpgradeWall from '@/components/UpgradeWall.vue';
import { useBillingStore } from '@/stores/billing';

const billingStore = useBillingStore();
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

onMounted(() => {
  billingStore.fetchPlan();
  load();
});

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
        if (!p) return '';
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
.global-analytics-page {
  padding: 0;
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 20px;

  @media (max-width: 575px) {
    gap: 16px;
  }
}

/* ── Page header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
}

.page-header__title-block {
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
  margin: 0;
  color: var(--md-sys-color-on-surface-variant);
}

.page-header__filters {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.date-range-row {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.date-range-sep {
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 600;
  padding-bottom: 8px;
}

/* ── Date field ──────────────────────────────────────────────────────────── */
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

.form-input {
  height: 36px;
  padding: 0 10px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.875rem;
  font-family: inherit;
  transition: border-color 0.15s, box-shadow 0.15s;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }
}

/* ── CSS Spinner ─────────────────────────────────────────────────────────── */
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

/* ── Period chip group ───────────────────────────────────────────────────── */
.chip-group {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.chip {
  display: inline-flex;
  align-items: center;
  padding: 5px 14px;
  border-radius: 999px;
  border: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface);
  cursor: pointer;
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  transition: background 0.15s, color 0.15s, border-color 0.15s;
  white-space: nowrap;
  line-height: 1.4;

  &:hover {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface);
    border-color: var(--md-sys-color-outline);
  }

  &.active {
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    border-color: var(--md-sys-color-primary);

    &:hover {
      background: var(--md-sys-color-primary);
      color: var(--md-sys-color-on-primary);
    }
  }
}

/* ── Card header row ─────────────────────────────────────────────────────── */
.card-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1.25rem;
  gap: 1rem;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.card-header-row__left {
  display: flex;
  align-items: center;
  gap: 0.625rem;
}

.card-header-row__icon {
  font-size: 20px;
  color: var(--md-sys-color-primary);
  flex-shrink: 0;
}

/* ── Error banner ────────────────────────────────────────────────────────── */
.error-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  background: var(--md-sys-color-error-container);
  border: 1px solid var(--md-sys-color-error);
}

/* ── Stat grid ───────────────────────────────────────────────────────────── */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;

  @media (max-width: 991px) { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.stat-skeleton-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
}

.skeleton-icon-box {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  flex-shrink: 0;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  border-radius: 16px;
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
}

.stat-card__icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined {
    font-size: 22px;
  }

  &--primary {
    background: var(--md-sys-color-primary-container);
    .material-symbols-outlined { color: var(--md-sys-color-on-primary-container); }
  }

  &--secondary {
    background: var(--md-sys-color-secondary-container);
    .material-symbols-outlined { color: var(--md-sys-color-on-secondary-container); }
  }

  &--teal {
    background: color-mix(in srgb, #14b8a6 16%, transparent);
    .material-symbols-outlined { color: #0f766e; }
  }

  &--amber {
    background: color-mix(in srgb, #f59e0b 16%, transparent);
    .material-symbols-outlined { color: #b45309; }
  }
}

.stat-card__body {
  min-width: 0;
}

.stat-value {
  color: var(--md-sys-color-on-surface);
  line-height: 1.2;
  margin-bottom: 0.15rem;
  font-weight: 600;
}

.stat-label {
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Chart section ───────────────────────────────────────────────────────── */
.chart-section {
  border-radius: 16px;
}

.chart-body-pad {
  padding: 0.75rem 1rem;
}

/* ── Trend indicator ─────────────────────────────────────────────────────── */
.trend-indicator {
  display: inline-flex;
  align-items: center;
  gap: 0.15rem;
  font-size: 0.8rem;
  font-weight: 600;
}

/* ── Comparison grid ─────────────────────────────────────────────────────── */
.comparison-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;

  @media (max-width: 991px) { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.comparison-metric {
  padding: 0.5rem 0;
}

.comparison-period-box {
  border-radius: 12px;
  padding: 0.75rem 1rem;

  &--current {
    background: color-mix(in srgb, var(--md-sys-color-primary) 8%, transparent);
    border: 1px solid color-mix(in srgb, var(--md-sys-color-primary) 25%, transparent);
  }

  &--previous {
    background: var(--md-sys-color-surface-container-low);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Two column ──────────────────────────────────────────────────────────── */
.two-col-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;

  @media (max-width: 767px) { grid-template-columns: 1fr; }
}

/* ── Breakdown list ──────────────────────────────────────────────────────── */
.breakdown-list {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.country-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.rank-num {
  min-width: 20px;
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 700;
  font-size: 0.78rem;
  flex-shrink: 0;
}

.flag-emoji {
  font-size: 1.1rem;
  line-height: 1;
  flex-shrink: 0;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* ── Mini progress bar ───────────────────────────────────────────────────── */
.mini-bar {
  width: 100px;
  height: 6px;
  background: var(--md-sys-color-surface-container-high);
  border-radius: 999px;
  overflow: hidden;
  flex-shrink: 0;

  &--sm { width: 80px; }
}

.mini-bar__fill {
  height: 100%;
  background: var(--md-sys-color-primary);
  border-radius: 999px;
  transition: width 0.3s ease;
}

/* ── UTM grid ────────────────────────────────────────────────────────────── */
.utm-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;

  @media (max-width: 991px) { grid-template-columns: 1fr 1fr; }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.utm-section-label {
  color: var(--md-sys-color-on-surface-variant);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-size: 0.7rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
}

/* ── M3 Cards ────────────────────────────────────────────────────────────── */
.m3-card {
  border-radius: 16px;
  overflow: hidden;

  &--elevated {
    background: var(--md-sys-color-surface);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.08);
  }

  &--outlined {
    background: var(--md-sys-color-surface);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── M3 Badges ───────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 0.2rem 0.65rem;
  border-radius: 6px;

  &--primary {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }
}

/* ── Empty state (inline) ────────────────────────────────────────────────── */
.empty-state {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  padding: 8px 0;
}

/* ── Skeleton ────────────────────────────────────────────────────────────── */
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

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ── M3 Empty State ──────────────────────────────────────────────────────── */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  text-align: center;

  &--compact {
    padding: 2.5rem 1.5rem;
  }

  &--full {
    padding: 5rem 2rem;
  }
}

.m3-empty-state__icon-wrap {
  width: 76px;
  height: 76px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.m3-empty-state__icon {
  font-size: 40px;
  color: var(--md-sys-color-primary);
  opacity: 0.65;
}

.m3-empty-state__title {
  margin: 0 0 0.25rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.m3-empty-state__text {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 1rem;
  max-width: 340px;
}
</style>
