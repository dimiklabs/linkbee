<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <RouterLink to="/dashboard/links">
          <md-icon-button>
            <span class="material-symbols-outlined">arrow_back</span>
          </md-icon-button>
        </RouterLink>
        <div class="page-header__title-group">
          <h1 class="md-title-large">Link Analytics</h1>
          <template v-if="link">
            <div class="page-header__subtitle">
              <span class="short-url-monospace">{{ link.short_url }}</span>
              <md-icon-button
                :class="copied ? 'copy-btn--success' : ''"
                @click="copyShortUrl"
                title="Copy short URL"
              >
                <span class="material-symbols-outlined" style="font-size:18px">{{ copied ? 'check' : 'content_copy' }}</span>
              </md-icon-button>
              <a :href="link.short_url" target="_blank" rel="noopener noreferrer">
                <md-icon-button title="Open short URL">
                  <span class="material-symbols-outlined" style="font-size:18px">open_in_new</span>
                </md-icon-button>
              </a>
            </div>
          </template>
        </div>
      </div>
      <div v-if="analytics" class="page-header__actions">
        <md-outlined-button @click="exportToCSV">
          <span class="material-symbols-outlined" slot="icon">download</span>
          Export CSV
        </md-outlined-button>
      </div>
    </div>

    <!-- Link Info Card -->
    <div class="m3-card m3-card--outlined link-info-card" style="margin-bottom:1.5rem;">
      <!-- Skeleton while loading -->
      <div v-if="linkLoading" style="display:flex;align-items:center;gap:1rem;padding:1rem 1.25rem;">
        <div class="skeleton" style="width:120px;height:1.1rem;border-radius:4px;"></div>
        <div class="skeleton" style="width:200px;height:0.9rem;border-radius:4px;"></div>
      </div>

      <!-- Populated -->
      <div v-else-if="link" style="padding:1rem 1.25rem;">
        <div class="link-info-body">
          <div class="link-info-main">
            <div class="link-title-row">
              <span class="md-title-medium link-info-title">{{ link.title || '/' + link.slug }}</span>
              <span class="m3-badge" :class="link.is_active ? 'm3-badge--success' : 'm3-badge--neutral'">
                {{ link.is_active ? 'Active' : 'Inactive' }}
              </span>
              <span class="m3-badge" :class="healthBadgeClass(link.health_status)">
                {{ healthLabel(link.health_status) }}
              </span>
            </div>
            <div class="link-dest-url md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:0.25rem;display:flex;align-items:center;gap:0.25rem;overflow:hidden;">
              <span class="material-symbols-outlined" style="font-size:14px;flex-shrink:0;">link</span>
              <span style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="link.destination_url">{{ link.destination_url }}</span>
            </div>
            <div v-if="link.tags && link.tags.length" style="display:flex;flex-wrap:wrap;gap:0.4rem;margin-top:0.5rem;">
              <span v-for="tag in link.tags" :key="tag" class="m3-badge m3-badge--neutral">{{ tag }}</span>
            </div>
          </div>
          <div class="link-info-stat">
            <div class="stat-number md-title-large" style="color:var(--md-sys-color-primary);">{{ link.click_count.toLocaleString() }}</div>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">All-time clicks</div>
          </div>
        </div>
      </div>

      <!-- Fallback -->
      <div v-else style="padding:1rem 1.25rem;">
        <span class="md-title-medium">Link Analytics</span>
        <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ route.params.id }}</div>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="error-banner" style="margin-bottom:1rem;">
      <span class="material-symbols-outlined" style="color:var(--md-sys-color-error);">error</span>
      <span class="md-body-medium" style="flex:1;">{{ error }}</span>
      <md-text-button @click="loadAnalytics">Retry</md-text-button>
    </div>

    <!-- Loading state -->
    <div v-if="loading && !error" style="display:flex;flex-direction:column;align-items:center;justify-content:center;padding:4rem 0;gap:1rem;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:48px" />
      <span class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);">Loading analytics data...</span>
    </div>

    <template v-if="!loading && analytics">

      <!-- Filters -->
      <div class="m3-card m3-card--elevated filters-card" style="margin-bottom:1.5rem;">
        <div class="m3-card-header">
          <div class="m3-card-header__left">
            <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">tune</span>
            <span class="md-label-large">Date Range &amp; Granularity</span>
          </div>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div class="filter-row">
            <md-outlined-text-field
              type="date"
              label="From"
              :value="filterFrom"
              @input="filterFrom = ($event.target as HTMLInputElement).value"
              style="min-width:160px;"
            />
            <md-outlined-text-field
              type="date"
              label="To"
              :value="filterTo"
              @input="filterTo = ($event.target as HTMLInputElement).value"
              style="min-width:160px;"
            />
            <AppSelect
              label="Granularity"
              :model-value="filterGranularity"
              @update:model-value="filterGranularity = $event as 'hour'|'day'|'week'|'month'"
              style="min-width:140px;"
            >
              <option value="hour">Hour</option>
              <option value="day">Day</option>
              <option value="week">Week</option>
              <option value="month">Month</option>
            </AppSelect>
            <md-filled-button :disabled="loading" @click="applyFilters">
              <span v-if="loading" slot="icon"><md-circular-progress indeterminate style="--md-circular-progress-size:18px" /></span>
              Apply
            </md-filled-button>
          </div>
        </div>
      </div>

      <!-- Period Comparison Card -->
      <div v-if="comparison" class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Period Comparison</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Current vs preceding period of equal length</div>
          </div>
          <span class="m3-badge m3-badge--neutral">
            {{ formatShortDate(comparison.previous.from) }} → {{ formatShortDate(comparison.current.to) }}
          </span>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div class="comparison-grid">
            <!-- Total Clicks -->
            <div class="comparison-metric">
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.25rem;">Total Clicks</div>
              <div style="display:flex;align-items:baseline;gap:0.5rem;flex-wrap:wrap;">
                <span class="md-headline-small">{{ comparison.current.total_clicks.toLocaleString() }}</span>
                <span class="trend-badge" :class="trendClass(comparison.clicks.trend)">
                  <span class="material-symbols-outlined" style="font-size:14px;vertical-align:middle;">{{ comparison.clicks.trend === 'up' ? 'trending_up' : comparison.clicks.trend === 'down' ? 'trending_down' : 'trending_flat' }}</span>
                  {{ Math.abs(comparison.clicks.percent_change).toFixed(1) }}%
                </span>
              </div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:0.2rem;">
                vs {{ comparison.previous.total_clicks.toLocaleString() }} prev. period
              </div>
            </div>
            <!-- Unique Clicks -->
            <div class="comparison-metric">
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.25rem;">Unique Clicks</div>
              <div style="display:flex;align-items:baseline;gap:0.5rem;flex-wrap:wrap;">
                <span class="md-headline-small">{{ comparison.current.unique_clicks.toLocaleString() }}</span>
                <span class="trend-badge" :class="trendClass(comparison.unique_clicks.trend)">
                  <span class="material-symbols-outlined" style="font-size:14px;vertical-align:middle;">{{ comparison.unique_clicks.trend === 'up' ? 'trending_up' : comparison.unique_clicks.trend === 'down' ? 'trending_down' : 'trending_flat' }}</span>
                  {{ Math.abs(comparison.unique_clicks.percent_change).toFixed(1) }}%
                </span>
              </div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:0.2rem;">
                vs {{ comparison.previous.unique_clicks.toLocaleString() }} prev. period
              </div>
            </div>
            <!-- Current period -->
            <div class="comparison-period-box comparison-period-box--current">
              <div class="md-label-large period-label">Current Period</div>
              <div class="md-body-small period-dates">{{ formatShortDate(comparison.current.from) }} – {{ formatShortDate(comparison.current.to) }}</div>
              <div style="display:flex;gap:1rem;margin-top:0.25rem;">
                <span class="md-body-small"><strong>{{ comparison.current.total_clicks.toLocaleString() }}</strong> clicks</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);"><strong>{{ comparison.current.unique_clicks.toLocaleString() }}</strong> unique</span>
              </div>
            </div>
            <!-- Previous period -->
            <div class="comparison-period-box comparison-period-box--previous">
              <div class="md-label-large period-label">Previous Period</div>
              <div class="md-body-small period-dates">{{ formatShortDate(comparison.previous.from) }} – {{ formatShortDate(comparison.previous.to) }}</div>
              <div style="display:flex;gap:1rem;margin-top:0.25rem;">
                <span class="md-body-small"><strong>{{ comparison.previous.total_clicks.toLocaleString() }}</strong> clicks</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);"><strong>{{ comparison.previous.unique_clicks.toLocaleString() }}</strong> unique</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Stat Cards -->
      <div class="stat-grid" style="margin-bottom:1.5rem;">
        <!-- Total Clicks -->
        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-card__header">
            <div style="display:flex;align-items:center;gap:0.5rem;">
              <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Total Clicks</span>
              <span v-if="isLive" class="live-badge">
                <span class="live-dot"></span>
                <span style="font-size:0.6rem;font-weight:600;letter-spacing:0.04em;color:#16a34a;">LIVE</span>
              </span>
            </div>
            <div class="stat-icon" style="background:rgba(var(--md-sys-color-primary-rgb,99,91,255),0.12);">
              <span class="material-symbols-outlined" style="color:var(--md-sys-color-primary);font-size:18px;">trending_up</span>
            </div>
          </div>
          <div class="md-headline-small stat-value">
            {{ (liveTotal !== null ? liveTotal : analytics.total_clicks).toLocaleString() }}
          </div>
        </div>

        <!-- Unique Clicks -->
        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-card__header">
            <div style="display:flex;align-items:center;gap:0.5rem;">
              <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Unique Clicks</span>
              <span v-if="isLive" class="live-badge">
                <span class="live-dot"></span>
                <span style="font-size:0.6rem;font-weight:600;letter-spacing:0.04em;color:#16a34a;">LIVE</span>
              </span>
            </div>
            <div class="stat-icon" style="background:rgba(20,184,166,0.12);">
              <span class="material-symbols-outlined" style="color:#14b8a6;font-size:18px;">group</span>
            </div>
          </div>
          <div class="md-headline-small stat-value">
            {{ (liveUnique !== null ? liveUnique : analytics.unique_clicks).toLocaleString() }}
          </div>
        </div>

        <!-- Top Referrer -->
        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-card__header">
            <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Top Referrer</span>
            <div class="stat-icon" style="background:rgba(245,158,11,0.12);">
              <span class="material-symbols-outlined" style="color:#f59e0b;font-size:18px;">link</span>
            </div>
          </div>
          <div class="md-title-medium stat-value" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="topReferrer">
            {{ topReferrer }}
          </div>
        </div>

        <!-- Top Device -->
        <div class="m3-card m3-card--elevated stat-card">
          <div class="stat-card__header">
            <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Top Device</span>
            <div class="stat-icon" style="background:rgba(239,68,68,0.12);">
              <span class="material-symbols-outlined" style="color:#ef4444;font-size:18px;">smartphone</span>
            </div>
          </div>
          <div class="md-title-medium stat-value" style="text-transform:capitalize;">
            {{ topDevice }}
          </div>
        </div>
      </div>

      <!-- Clicks Over Time Chart -->
      <div class="m3-card m3-card--elevated" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <span class="md-title-medium">Clicks Over Time</span>
        </div>
        <md-divider />
        <div style="padding:1rem;">
          <div v-if="analytics.time_series.length === 0" class="empty-state">
            No time series data available for this period.
          </div>
          <VChart
            v-else
            :option="chartOption"
            style="height:320px;"
            autoresize
          />
        </div>
      </div>

      <!-- Traffic Channels -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;margin-top:1rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Traffic Channels</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">How visitors are reaching your link — classified by referrer domain</div>
          </div>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div v-if="!analytics.referrer_categories?.length" class="empty-state">
            No channel data available yet.
          </div>
          <div v-else class="two-col-layout">
            <!-- Donut chart -->
            <div style="min-width:0;">
              <VChart :option="trafficChannelsChartOption" style="height:280px;" autoresize />
            </div>
            <!-- Channel breakdown list -->
            <div class="breakdown-list">
              <div
                v-for="ch in analytics.referrer_categories"
                :key="ch.category"
                class="breakdown-item"
              >
                <div style="display:flex;align-items:center;gap:0.75rem;margin-bottom:6px;">
                  <span
                    style="display:inline-flex;align-items:center;justify-content:center;border-radius:50%;flex-shrink:0;width:32px;height:32px;"
                    :style="{ backgroundColor: channelColor(ch.category) + '22' }"
                  >
                    <span style="font-size:0.95rem;">{{ channelIcon(ch.category) }}</span>
                  </span>
                  <div style="flex:1;min-width:0;">
                    <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
                      <span class="md-body-medium">{{ channelLabel(ch.category) }}</span>
                      <span class="md-label-large">{{ ch.count.toLocaleString() }} ({{ channelPercent(ch.count) }}%)</span>
                    </div>
                    <md-linear-progress :value="channelPercent(ch.count) / 100" style="--md-linear-progress-track-height:6px;--md-linear-progress-active-indicator-height:6px" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Two-column: Referrers + Devices -->
      <div class="two-col-grid" style="margin-bottom:1.5rem;">
        <!-- Top Referrers -->
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Top Referrers</span>
            <span class="m3-badge m3-badge--neutral">Top 10</span>
          </div>
          <md-divider />
          <div style="padding:0.5rem 0;">
            <div v-if="topReferrers.length === 0" class="m3-empty-state">
              <span class="material-symbols-outlined" style="font-size:2rem;color:var(--md-sys-color-on-surface-variant);opacity:0.5;">link_off</span>
              <p class="md-body-medium" style="margin:0.5rem 0 0;">No referrer data available yet.</p>
            </div>
            <div v-else class="m3-table-wrapper">
            <table class="m3-table">
              <thead>
                <tr>
                  <th>Referrer</th>
                  <th style="text-align:right;">Clicks</th>
                  <th style="width:120px;">Share</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="ref in topReferrers" :key="ref.referrer">
                  <td style="max-width:200px;">
                    <span style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;display:block;" :title="ref.referrer">
                      {{ ref.referrer || 'Direct / None' }}
                    </span>
                  </td>
                  <td style="text-align:right;font-weight:600;">{{ ref.count.toLocaleString() }}</td>
                  <td>
                    <div style="display:flex;align-items:center;gap:0.5rem;">
                      <md-linear-progress :value="referrerPercent(ref.count) / 100" style="flex:1;--md-linear-progress-track-height:6px;--md-linear-progress-active-indicator-height:6px" />
                      <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:32px;">{{ referrerPercent(ref.count) }}%</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            </div>
          </div>
        </div>

        <!-- Device Breakdown -->
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Device Breakdown</span>
            <span class="m3-badge m3-badge--neutral">All Devices</span>
          </div>
          <md-divider />
          <div style="padding:0.5rem 0;">
            <div v-if="analytics.devices.length === 0" class="m3-empty-state">
              <span class="material-symbols-outlined" style="font-size:2rem;color:var(--md-sys-color-on-surface-variant);opacity:0.5;">devices</span>
              <p class="md-body-medium" style="margin:0.5rem 0 0;">No device data available yet.</p>
            </div>
            <div v-else class="m3-table-wrapper">
            <table class="m3-table">
              <thead>
                <tr>
                  <th>Device</th>
                  <th style="text-align:right;">Clicks</th>
                  <th style="width:120px;">Share</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="device in analytics.devices" :key="device.device_type">
                  <td>
                    <div style="display:flex;align-items:center;gap:0.5rem;">
                      <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-primary);">
                        {{ device.device_type.toLowerCase() === 'mobile' ? 'smartphone' : device.device_type.toLowerCase() === 'tablet' ? 'tablet' : 'computer' }}
                      </span>
                      <span style="text-transform:capitalize;">{{ device.device_type || 'Unknown' }}</span>
                    </div>
                  </td>
                  <td style="text-align:right;font-weight:600;">{{ device.count.toLocaleString() }}</td>
                  <td>
                    <div style="display:flex;align-items:center;gap:0.5rem;">
                      <md-linear-progress :value="devicePercent(device.count) / 100" style="flex:1;--md-linear-progress-track-height:6px;--md-linear-progress-active-indicator-height:6px" />
                      <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:32px;">{{ devicePercent(device.count) }}%</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Two-column: Browser + OS -->
      <div class="two-col-grid" style="margin-bottom:1.5rem;">
        <!-- Browser Breakdown -->
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Browser Breakdown</span>
            <span class="m3-badge m3-badge--neutral">All Browsers</span>
          </div>
          <md-divider />
          <div style="padding:1rem 1.25rem;">
            <div v-if="!analytics.browsers || analytics.browsers.length === 0" class="empty-state">
              No browser data available.
            </div>
            <div v-else class="chart-legend-layout">
              <VChart :option="browserChartOption" style="height:200px;" autoresize />
              <div class="legend-list">
                <div v-for="(b, i) in analytics.browsers" :key="b.browser" class="legend-item">
                  <span class="legend-dot" :style="{ backgroundColor: chartColors[i % chartColors.length] }"></span>
                  <span class="md-body-small" style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="b.browser">{{ b.browser }}</span>
                  <span class="md-label-large">{{ browserPercent(b.count) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- OS Breakdown -->
        <div class="m3-card m3-card--outlined">
          <div class="card-header-row">
            <span class="md-title-medium">Operating System</span>
            <span class="m3-badge m3-badge--neutral">All OS</span>
          </div>
          <md-divider />
          <div style="padding:1rem 1.25rem;">
            <div v-if="!analytics.os_breakdown || analytics.os_breakdown.length === 0" class="empty-state">
              No OS data available.
            </div>
            <div v-else class="chart-legend-layout">
              <VChart :option="osChartOption" style="height:200px;" autoresize />
              <div class="legend-list">
                <div v-for="(o, i) in analytics.os_breakdown" :key="o.os" class="legend-item">
                  <span class="legend-dot" :style="{ backgroundColor: chartColors[i % chartColors.length] }"></span>
                  <span class="md-body-small" style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="o.os">{{ o.os }}</span>
                  <span class="md-label-large">{{ osPercent(o.count) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Source Breakdown -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Click Source</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">How visitors arrived — direct web link, QR code, or API call</div>
          </div>
          <span class="m3-badge m3-badge--neutral">All Sources</span>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div v-if="!analytics.sources || analytics.sources.length === 0" class="empty-state">
            No source data available.
          </div>
          <div v-else class="chart-legend-layout">
            <VChart :option="sourceChartOption" style="height:200px;" autoresize />
            <table class="m3-table" style="flex:1;">
              <tbody>
                <tr v-for="(s, i) in analytics.sources" :key="s.source">
                  <td>
                    <div style="display:flex;align-items:center;gap:0.5rem;">
                      <span class="legend-dot" :style="{ backgroundColor: chartColors[i % chartColors.length] }"></span>
                      <span class="md-body-small">{{ sourceIcon(s.source) }} {{ sourceLabel(s.source) }}</span>
                    </div>
                  </td>
                  <td style="text-align:right;font-weight:600;">{{ s.count.toLocaleString() }}</td>
                  <td style="text-align:right;" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ sourcePercent(s.count) }}%</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Visitor Loyalty -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Visitor Loyalty</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">First-time vs returning visitors, based on hashed IP tracking</div>
          </div>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div
            v-if="!analytics.first_time_visitors && !analytics.returning_visitors"
            class="empty-state"
          >
            No visitor loyalty data yet. Data is collected on new clicks.
          </div>
          <div v-else class="chart-legend-layout">
            <VChart :option="visitorLoyaltyChartOption" style="height:200px;" autoresize />
            <div style="flex:1;display:grid;grid-template-columns:repeat(auto-fill,minmax(140px,1fr));gap:0.75rem;align-content:center;">
              <div class="loyalty-box loyalty-box--first">
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">First-Time</div>
                <div class="md-headline-small" style="color:var(--md-sys-color-primary);">{{ analytics.first_time_visitors.toLocaleString() }}</div>
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ firstTimePercent }}% of tracked</div>
              </div>
              <div class="loyalty-box loyalty-box--return">
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Returning</div>
                <div class="md-headline-small" style="color:#14b8a6;">{{ analytics.returning_visitors.toLocaleString() }}</div>
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ returningPercent }}% of tracked</div>
              </div>
              <div class="loyalty-box loyalty-box--rate">
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Return Rate</div>
                <div class="md-headline-small" style="color:#f59e0b;">{{ returningPercent }}%</div>
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">visitors who came back</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Link Timeline Card -->
      <div v-if="link" class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Link Timeline</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Key dates and milestones for this link</div>
          </div>
        </div>
        <md-divider />
        <div style="padding:0.75rem 1.25rem;">
          <div class="timeline-list">
            <!-- Created -->
            <div class="timeline-item">
              <div class="timeline-icon timeline-icon--created">
                <span class="material-symbols-outlined" style="font-size:16px;">celebration</span>
              </div>
              <div class="timeline-body">
                <div class="md-label-large timeline-label">Created</div>
                <div class="md-body-medium timeline-date">{{ formatLinkDate(link.created_at) }}</div>
              </div>
            </div>
            <!-- Last Updated -->
            <div v-if="link.updated_at && link.updated_at !== link.created_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--updated">
                <span class="material-symbols-outlined" style="font-size:16px;">edit</span>
              </div>
              <div class="timeline-body">
                <div class="md-label-large timeline-label">Last Updated</div>
                <div class="md-body-medium timeline-date">{{ formatLinkDate(link.updated_at) }}</div>
              </div>
            </div>
            <!-- Health Last Checked -->
            <div v-if="link.health_checked_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--health">
                <span class="material-symbols-outlined" style="font-size:16px;">search</span>
              </div>
              <div class="timeline-body">
                <div style="display:flex;align-items:center;gap:0.5rem;flex-wrap:wrap;">
                  <div class="md-label-large timeline-label">Health Last Checked</div>
                  <span class="m3-badge" :class="healthBadgeClass(link.health_status)">{{ healthLabel(link.health_status) }}</span>
                </div>
                <div class="md-body-medium timeline-date">{{ formatLinkDate(link.health_checked_at) }}</div>
              </div>
            </div>
            <!-- Expires -->
            <div v-if="link.expires_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--expires">
                <span class="material-symbols-outlined" style="font-size:16px;">schedule</span>
              </div>
              <div class="timeline-body">
                <div class="md-label-large timeline-label">Expires</div>
                <div class="md-body-medium timeline-date" :class="expiryDateClass(link.expires_at)">
                  {{ formatLinkDate(link.expires_at) }}
                  <span class="md-body-small" :class="expiryDateClass(link.expires_at)">{{ expiryCountdown(link.expires_at) }}</span>
                </div>
              </div>
            </div>
            <!-- Click Limit -->
            <div v-if="link.max_clicks" class="timeline-item">
              <div class="timeline-icon timeline-icon--limit">
                <span class="material-symbols-outlined" style="font-size:16px;">flag</span>
              </div>
              <div class="timeline-body">
                <div class="md-label-large timeline-label">Click Limit</div>
                <div class="md-body-medium timeline-date">
                  Max clicks: <strong>{{ link.max_clicks.toLocaleString() }}</strong>
                  <span style="color:var(--md-sys-color-on-surface-variant);"> ({{ link.click_count.toLocaleString() }} used, {{ Math.max(0, link.max_clicks - link.click_count).toLocaleString() }} remaining)</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Security & Limits Card -->
      <div v-if="link" class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Security &amp; Limits</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Access controls and usage limits for this link</div>
          </div>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div class="security-grid">
            <div class="security-item">
              <div class="md-body-small security-item-label">Password Protected</div>
              <div>
                <span v-if="link.has_password" class="m3-badge m3-badge--warning">
                  <span class="material-symbols-outlined" style="font-size:12px;vertical-align:middle;">lock</span> Yes
                </span>
                <span v-else class="m3-badge m3-badge--neutral">No</span>
              </div>
            </div>
            <div class="security-item">
              <div class="md-body-small security-item-label">Active Status</div>
              <span class="m3-badge" :class="link.is_active ? 'm3-badge--success' : 'm3-badge--neutral'">
                {{ link.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>
            <div class="security-item">
              <div class="md-body-small security-item-label">Redirect Type</div>
              <span class="m3-badge m3-badge--neutral">
                {{ link.redirect_type === 301 ? '301 Permanent' : '302 Temporary' }}
              </span>
            </div>
            <div class="security-item">
              <div class="md-body-small security-item-label">Expiry</div>
              <span v-if="link.expires_at" :class="expiryDateClass(link.expires_at)" class="md-body-medium">
                {{ formatLinkDate(link.expires_at) }}
                <span class="md-body-small">{{ expiryCountdown(link.expires_at) }}</span>
              </span>
              <span v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Never expires</span>
            </div>
            <!-- Click Limit full width -->
            <div class="security-item security-item--full">
              <div class="md-body-small security-item-label" style="margin-bottom:0.5rem;">Click Limit</div>
              <div v-if="link.max_clicks">
                <div style="display:flex;justify-content:space-between;align-items:baseline;margin-bottom:0.25rem;">
                  <span class="md-body-medium">{{ link.click_count.toLocaleString() }} / {{ link.max_clicks.toLocaleString() }} clicks used</span>
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ clickLimitPercent(link.click_count, link.max_clicks) }}%</span>
                </div>
                <md-linear-progress
                  :value="clickLimitPercent(link.click_count, link.max_clicks) / 100"
                  style="--md-linear-progress-track-height:8px;--md-linear-progress-active-indicator-height:8px;"
                />
              </div>
              <div v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Unlimited</div>
            </div>
            <!-- Tags full width -->
            <div class="security-item security-item--full">
              <div class="md-body-small security-item-label" style="margin-bottom:0.25rem;">Tags</div>
              <div v-if="link.tags && link.tags.length" style="display:flex;flex-wrap:wrap;gap:0.4rem;">
                <span v-for="tag in link.tags" :key="tag" class="m3-badge m3-badge--neutral">{{ tag }}</span>
              </div>
              <span v-else class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">None</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Heatmap -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Click Heatmap</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Clicks by hour of day × day of week (UTC)</div>
          </div>
          <span class="m3-badge m3-badge--neutral">24 × 7</span>
        </div>
        <md-divider />
        <div style="padding:1rem;">
          <div v-if="!analytics.heatmap || analytics.heatmap.length === 0" class="empty-state">
            No heatmap data available for this period.
          </div>
          <VChart v-else :option="heatmapChartOption" style="height:260px;" autoresize />
        </div>
      </div>

      <!-- Geographic World Map -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">Geographic Distribution</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Click density by country — scroll to zoom, drag to pan</div>
          </div>
          <span class="m3-badge m3-badge--neutral">World Map</span>
        </div>
        <md-divider />
        <div style="padding:1rem;">
          <div v-if="worldMapLoading" style="display:flex;flex-direction:column;align-items:center;padding:3rem 0;gap:0.75rem;">
            <md-circular-progress indeterminate style="--md-circular-progress-size:32px" />
            <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Loading world map…</span>
          </div>
          <div v-else-if="!mapLoaded || !analytics?.countries?.length" class="empty-state">
            No geographic data available yet. Geographic data is collected on new clicks.
          </div>
          <VChart
            v-else
            :option="geoMapOption"
            style="height:420px;"
            autoresize
          />
        </div>
      </div>

      <!-- Top Countries -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <span class="md-title-medium">Top Countries</span>
          <span class="m3-badge m3-badge--neutral">Top 15</span>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div v-if="!analytics.countries || analytics.countries.length === 0" class="empty-state">
            No country data available yet. Geographic data is collected on new clicks.
          </div>
          <div v-else class="two-col-layout">
            <div style="min-width:0;">
              <VChart :option="countriesChartOption" style="height:360px;" autoresize />
            </div>
            <div style="min-width:0;">
              <table class="m3-table">
                <thead>
                  <tr>
                    <th>Country</th>
                    <th style="text-align:right;">Clicks</th>
                    <th style="width:120px;">Share</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="c in analytics.countries" :key="c.country">
                    <td>
                      <div style="display:flex;align-items:center;gap:0.5rem;">
                        <span style="font-size:1.1rem;line-height:1;">{{ countryFlag(c.country) }}</span>
                        <span class="md-body-small">{{ countryDisplayName(c.country) }}</span>
                      </div>
                    </td>
                    <td style="text-align:right;font-weight:600;">{{ c.count.toLocaleString() }}</td>
                    <td>
                      <div style="display:flex;align-items:center;gap:0.4rem;">
                        <md-linear-progress :value="countryPercent(c.count) / 100" style="flex:1;--md-linear-progress-track-height:6px;--md-linear-progress-active-indicator-height:6px" />
                        <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:32px;">{{ countryPercent(c.count) }}%</span>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Top Cities -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <span class="md-title-medium">Top Cities</span>
          <span class="m3-badge m3-badge--neutral">Top 20</span>
        </div>
        <md-divider />
        <div style="padding:0.5rem 0;">
          <div v-if="!analytics.cities || analytics.cities.length === 0" class="empty-state">
            No city data available yet. City data is collected on new clicks.
          </div>
          <table v-else class="m3-table">
            <thead>
              <tr>
                <th style="width:36px;">#</th>
                <th>City</th>
                <th>Country</th>
                <th style="text-align:right;">Clicks</th>
                <th style="width:140px;">Share</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(c, idx) in analytics.cities" :key="c.city + c.country">
                <td class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ idx + 1 }}</td>
                <td><span class="md-body-medium">{{ c.city }}</span></td>
                <td>
                  <div style="display:flex;align-items:center;gap:0.5rem;">
                    <span style="font-size:1rem;line-height:1;">{{ countryFlag(c.country) }}</span>
                    <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ countryDisplayName(c.country) }}</span>
                  </div>
                </td>
                <td style="text-align:right;font-weight:600;">{{ c.count.toLocaleString() }}</td>
                <td>
                  <div style="display:flex;align-items:center;gap:0.4rem;">
                    <md-linear-progress :value="cityPercent(c.count) / 100" style="flex:1;--md-linear-progress-track-height:6px;--md-linear-progress-active-indicator-height:6px" />
                    <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);min-width:32px;">{{ cityPercent(c.count) }}%</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- UTM Campaign Tracking -->
      <div class="m3-card m3-card--outlined" style="margin-bottom:1.5rem;">
        <div class="card-header-row">
          <div>
            <span class="md-title-medium">UTM Campaign Tracking</span>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Clicks from links with UTM parameters appended</div>
          </div>
        </div>
        <md-divider />
        <div style="padding:1rem 1.25rem;">
          <div
            v-if="!analytics.utm_sources?.length && !analytics.utm_mediums?.length && !analytics.utm_campaigns?.length && !analytics.utm_contents?.length && !analytics.utm_terms?.length"
            class="empty-state"
          >
            No UTM data yet. Append <code>?utm_source=email&amp;utm_medium=newsletter&amp;utm_campaign=launch</code> to your short link when sharing.
          </div>
          <div v-else>
            <!-- Row 1: Source, Medium, Campaign -->
            <div class="utm-grid">
              <div>
                <div class="md-label-large" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.75rem;">Source</div>
                <div v-if="!analytics.utm_sources?.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No data</div>
                <div v-else class="breakdown-list">
                  <div v-for="u in analytics.utm_sources" :key="u.value" class="breakdown-item">
                    <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
                      <span class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:140px;" :title="u.value">{{ u.value }}</span>
                      <span class="md-label-large">{{ u.count.toLocaleString() }}</span>
                    </div>
                    <md-linear-progress :value="utmPercent(u.count, analytics.utm_sources) / 100" style="--md-linear-progress-track-height:5px;--md-linear-progress-active-indicator-height:5px" />
                  </div>
                </div>
              </div>
              <div>
                <div class="md-label-large" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.75rem;">Medium</div>
                <div v-if="!analytics.utm_mediums?.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No data</div>
                <div v-else class="breakdown-list">
                  <div v-for="u in analytics.utm_mediums" :key="u.value" class="breakdown-item">
                    <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
                      <span class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:140px;" :title="u.value">{{ u.value }}</span>
                      <span class="md-label-large">{{ u.count.toLocaleString() }}</span>
                    </div>
                    <md-linear-progress :value="utmPercent(u.count, analytics.utm_mediums) / 100" style="--md-linear-progress-track-height:5px;--md-linear-progress-active-indicator-height:5px" />
                  </div>
                </div>
              </div>
              <div>
                <div class="md-label-large" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.75rem;">Campaign</div>
                <div v-if="!analytics.utm_campaigns?.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No data</div>
                <div v-else class="breakdown-list">
                  <div v-for="u in analytics.utm_campaigns" :key="u.value" class="breakdown-item">
                    <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
                      <span class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:140px;" :title="u.value">{{ u.value }}</span>
                      <span class="md-label-large">{{ u.count.toLocaleString() }}</span>
                    </div>
                    <md-linear-progress :value="utmPercent(u.count, analytics.utm_campaigns) / 100" style="--md-linear-progress-track-height:5px;--md-linear-progress-active-indicator-height:5px" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Row 2: Content + Term -->
            <div v-if="analytics.utm_contents?.length || analytics.utm_terms?.length" style="margin-top:1.5rem;">
              <md-divider style="margin-bottom:1.5rem;" />
              <div class="two-col-grid">
                <div>
                  <div class="md-label-large" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.75rem;">Content</div>
                  <div v-if="!analytics.utm_contents?.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No data</div>
                  <div v-else class="breakdown-list">
                    <div v-for="u in analytics.utm_contents" :key="u.value" class="breakdown-item">
                      <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
                        <span class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:200px;" :title="u.value">{{ u.value }}</span>
                        <span class="md-label-large">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <md-linear-progress :value="utmPercent(u.count, analytics.utm_contents) / 100" style="--md-linear-progress-track-height:5px;--md-linear-progress-active-indicator-height:5px" />
                    </div>
                  </div>
                </div>
                <div>
                  <div class="md-label-large" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:0.75rem;">Term</div>
                  <div v-if="!analytics.utm_terms?.length" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">No data</div>
                  <div v-else class="breakdown-list">
                    <div v-for="u in analytics.utm_terms" :key="u.value" class="breakdown-item">
                      <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
                        <span class="md-body-medium" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:200px;" :title="u.value">{{ u.value }}</span>
                        <span class="md-label-large">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <md-linear-progress :value="utmPercent(u.count, analytics.utm_terms) / 100" style="--md-linear-progress-track-height:5px;--md-linear-progress-active-indicator-height:5px" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { RouterLink } from 'vue-router';
import AppSelect from '@/components/AppSelect.vue';
import { use, registerMap } from 'echarts/core';
import { LineChart, PieChart, BarChart, HeatmapChart, MapChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent, VisualMapComponent, GeoComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import VChart from 'vue-echarts';
import linksApi from '@/api/links';
import type { AnalyticsResponse, UTMPoint, SourcePoint, PeriodComparisonResponse, LinkResponse } from '@/types/links';

use([LineChart, PieChart, BarChart, HeatmapChart, MapChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent, VisualMapComponent, GeoComponent, CanvasRenderer]);

const route = useRoute();
const router = useRouter();

const analytics = ref<AnalyticsResponse | null>(null);
const comparison = ref<PeriodComparisonResponse | null>(null);
const link = ref<LinkResponse | null>(null);
const linkLoading = ref(false);
const copied = ref(false);
const loading = ref(false);
const error = ref('');

// World map
const mapLoaded = ref(false);
const worldMapLoading = ref(false);

// ISO alpha-2 overrides where Intl.DisplayNames diverges from ECharts world map names
const isoToEchartsName: Record<string, string> = {
  US: 'United States', GB: 'United Kingdom', KR: 'South Korea',
  KP: 'North Korea', TW: 'Taiwan', VN: 'Vietnam', LA: 'Lao PDR',
  MM: 'Myanmar', CD: 'Dem. Rep. Congo', CG: 'Congo', BO: 'Bolivia',
  VE: 'Venezuela', IR: 'Iran', SY: 'Syria', RU: 'Russia',
  MK: 'Macedonia', PS: 'Palestine', XK: 'Kosovo',
};

function countryToEchartsName(iso: string): string {
  return isoToEchartsName[iso] ?? countryDisplayName(iso);
}

async function loadWorldGeoJSON() {
  if (mapLoaded.value) return;
  worldMapLoading.value = true;
  try {
    const res = await fetch(
      'https://cdn.jsdelivr.net/gh/apache/echarts-examples@main/public/data/asset/geo/World.json',
    );
    if (!res.ok) throw new Error('fetch failed');
    const geoJson = await res.json() as Parameters<typeof registerMap>[1];
    registerMap('world', geoJson);
    mapLoaded.value = true;
  } catch {
    // Fail gracefully — country table still shows
  } finally {
    worldMapLoading.value = false;
  }
}

// Live click counter (SSE)
const isLive = ref(false);
const liveTotal = ref<number | null>(null);
const liveUnique = ref<number | null>(null);
let eventSource: EventSource | null = null;

function connectLiveCount() {
  const id = route.params.id as string;
  const url = linksApi.getLiveCountUrl(id);
  eventSource = new EventSource(url);

  eventSource.addEventListener('count', (e: MessageEvent) => {
    try {
      const data = JSON.parse(e.data) as { total_clicks: number; unique_clicks: number };
      liveTotal.value = data.total_clicks;
      liveUnique.value = data.unique_clicks;
      isLive.value = true;
    } catch {
      // ignore parse errors
    }
  });

  eventSource.onerror = () => {
    isLive.value = false;
    eventSource?.close();
    eventSource = null;
  };
}

function disconnectLiveCount() {
  eventSource?.close();
  eventSource = null;
  isLive.value = false;
}

async function loadLink() {
  const id = route.params.id as string;
  linkLoading.value = true;
  try {
    const res = await linksApi.get(id);
    link.value = res.data;
  } catch {
    // non-fatal — header just won't show link metadata
  } finally {
    linkLoading.value = false;
  }
}

async function copyShortUrl() {
  if (!link.value?.short_url) return;
  try {
    await navigator.clipboard.writeText(link.value.short_url);
    copied.value = true;
    setTimeout(() => (copied.value = false), 2000);
  } catch {
    // ignore
  }
}

function healthLabel(status: string): string {
  const map: Record<string, string> = { healthy: 'Healthy', unhealthy: 'Unhealthy', timeout: 'Timeout', error: 'Error', unknown: 'Unchecked' };
  return map[status] ?? 'Unknown';
}

function healthClass(status: string): string {
  const map: Record<string, string> = { healthy: 'badge-health--ok', unhealthy: 'badge-health--bad', timeout: 'badge-health--warn', error: 'badge-health--warn', unknown: 'badge-health--gray' };
  return map[status] ?? 'badge-health--gray';
}

function healthBadgeClass(status: string): string {
  const map: Record<string, string> = { healthy: 'm3-badge--success', unhealthy: 'm3-badge--error', timeout: 'm3-badge--warning', error: 'm3-badge--warning', unknown: 'm3-badge--neutral' };
  return map[status] ?? 'm3-badge--neutral';
}

// Feature #58 / #62 helpers
function formatLinkDate(iso: string | undefined): string {
  if (!iso) return '';
  const d = new Date(iso);
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    timeZone: 'UTC',
    timeZoneName: 'short',
  });
}

function expiryDateClass(iso: string | undefined): string {
  if (!iso) return '';
  const now = Date.now();
  const exp = new Date(iso).getTime();
  const diffDays = (exp - now) / 86_400_000;
  if (diffDays <= 0) return 'text-danger';
  if (diffDays <= 3) return 'text-danger';
  if (diffDays <= 7) return 'text-warning';
  return '';
}

function expiryCountdown(iso: string | undefined): string {
  if (!iso) return '';
  const now = Date.now();
  const exp = new Date(iso).getTime();
  const diffMs = exp - now;
  if (diffMs <= 0) return '(expired)';
  const diffDays = Math.floor(diffMs / 86_400_000);
  if (diffDays === 0) {
    const diffHrs = Math.floor(diffMs / 3_600_000);
    return diffHrs > 0 ? `(${diffHrs}h remaining)` : '(< 1h remaining)';
  }
  return `(${diffDays}d remaining)`;
}

function clickLimitPercent(used: number, max: number): number {
  if (!max) return 0;
  return Math.min(100, Math.round((used / max) * 100));
}

function clickLimitBarColor(used: number, max: number): string {
  const pct = clickLimitPercent(used, max);
  if (pct > 90) return '#ef4444';
  if (pct > 75) return '#f59e0b';
  return '#635bff';
}

// Filter state — default to last 30 days
const now = new Date();
const thirtyDaysAgo = new Date(now);
thirtyDaysAgo.setDate(now.getDate() - 30);

const filterFrom = ref(thirtyDaysAgo.toISOString().slice(0, 10));
const filterTo = ref(now.toISOString().slice(0, 10));
const filterGranularity = ref<'hour' | 'day' | 'week' | 'month'>('day');

async function loadAnalytics() {
  const id = route.params.id as string;
  loading.value = true;
  error.value = '';
  try {
    analytics.value = await linksApi.getAnalytics(
      id,
      filterFrom.value ? new Date(filterFrom.value).toISOString() : undefined,
      filterTo.value ? new Date(filterTo.value).toISOString() : undefined,
      filterGranularity.value,
    );
  } catch (err: unknown) {
    if (err instanceof Error) {
      error.value = err.message;
    } else {
      error.value = 'Failed to load analytics. Please try again.';
    }
  } finally {
    loading.value = false;
  }
}

async function loadComparison() {
  const id = route.params.id as string;
  try {
    const res = await linksApi.getComparison(
      id,
      filterFrom.value ? new Date(filterFrom.value).toISOString() : undefined,
      filterTo.value ? new Date(filterTo.value).toISOString() : undefined,
    );
    comparison.value = res.data;
  } catch {
    comparison.value = null;
  }
}

function applyFilters() {
  loadAnalytics();
  loadComparison();
}

function trendClass(trend: string): string {
  if (trend === 'up') return 'comp-badge--up';
  if (trend === 'down') return 'comp-badge--down';
  return 'comp-badge--stable';
}

function trendArrow(trend: string): string {
  if (trend === 'up') return '↑ ';
  if (trend === 'down') return '↓ ';
  return '→ ';
}

function formatShortDate(iso: string): string {
  if (!iso) return '';
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: '2-digit' });
}

function exportToCSV() {
  const a = analytics.value;
  if (!a) return;

  const esc = (v: string | number) => {
    const s = String(v);
    return s.includes(',') || s.includes('"') || s.includes('\n') ? `"${s.replace(/"/g, '""')}"` : s;
  };

  const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
  const rows: string[] = [];

  // ── Report header ───────────────────────────────────────────────────────────
  rows.push('LINK ANALYTICS REPORT');
  rows.push(`Link ID,${esc(a.link_id)}`);
  rows.push(`Exported,${new Date().toISOString().slice(0, 10)}`);
  rows.push(`Period,${filterFrom.value} to ${filterTo.value}`);
  rows.push(`Granularity,${filterGranularity.value}`);
  rows.push('');

  // ── Summary ─────────────────────────────────────────────────────────────────
  rows.push('SUMMARY');
  rows.push('Total Clicks,Unique Clicks,First-Time Visitors,Returning Visitors,Return Rate');
  const returnRate = totalTrackedVisitors.value
    ? Math.round((a.returning_visitors / totalTrackedVisitors.value) * 100)
    : 0;
  rows.push(`${a.total_clicks},${a.unique_clicks},${a.first_time_visitors},${a.returning_visitors},${returnRate}%`);
  rows.push('');

  // ── Visitor loyalty ──────────────────────────────────────────────────────────
  rows.push('VISITOR LOYALTY');
  rows.push('Type,Count,Percentage');
  rows.push(`First-Time,${a.first_time_visitors},${firstTimePercent.value}%`);
  rows.push(`Returning,${a.returning_visitors},${returningPercent.value}%`);
  rows.push('');

  // ── Clicks over time ────────────────────────────────────────────────────────
  rows.push('CLICKS OVER TIME');
  rows.push('Date,Clicks');
  a.time_series.forEach((p) => rows.push(`${p.timestamp.substring(0, 10)},${p.count}`));
  rows.push('');

  // ── Top referrers ───────────────────────────────────────────────────────────
  rows.push('TOP REFERRERS');
  rows.push('Referrer,Clicks');
  a.referrers.forEach((r) => rows.push(`${esc(r.referrer || 'Direct / None')},${r.count}`));
  rows.push('');

  // ── Traffic channels ────────────────────────────────────────────────────────
  if (a.referrer_categories?.length) {
    rows.push('TRAFFIC CHANNELS');
    rows.push('Channel,Clicks');
    a.referrer_categories.forEach((c) => rows.push(`${esc(channelLabel(c.category))},${c.count}`));
    rows.push('');
  }

  // ── Device breakdown ────────────────────────────────────────────────────────
  rows.push('DEVICE BREAKDOWN');
  rows.push('Device,Clicks');
  a.devices.forEach((d) => rows.push(`${esc(d.device_type || 'Unknown')},${d.count}`));
  rows.push('');

  // ── Browser breakdown ───────────────────────────────────────────────────────
  rows.push('BROWSER BREAKDOWN');
  rows.push('Browser,Clicks');
  a.browsers.forEach((b) => rows.push(`${esc(b.browser || 'Unknown')},${b.count}`));
  rows.push('');

  // ── OS breakdown ────────────────────────────────────────────────────────────
  rows.push('OS BREAKDOWN');
  rows.push('OS,Clicks');
  a.os_breakdown.forEach((o) => rows.push(`${esc(o.os || 'Unknown')},${o.count}`));
  rows.push('');

  // ── Top countries ───────────────────────────────────────────────────────────
  rows.push('TOP COUNTRIES');
  rows.push('Country Code,Country,Clicks');
  a.countries.forEach((c) => {
    const name = countryDisplayName(c.country);
    rows.push(`${esc(c.country)},${esc(name)},${c.count}`);
  });
  rows.push('');

  // ── Top cities ───────────────────────────────────────────────────────────────
  if (a.cities?.length) {
    rows.push('TOP CITIES');
    rows.push('City,Country Code,Country,Clicks');
    a.cities.forEach((c) => {
      const name = countryDisplayName(c.country);
      rows.push(`${esc(c.city)},${esc(c.country)},${esc(name)},${c.count}`);
    });
    rows.push('');
  }

  // ── Click source breakdown ───────────────────────────────────────────────────
  if (a.sources?.length) {
    rows.push('CLICK SOURCE');
    rows.push('Source,Clicks');
    a.sources.forEach((s) => rows.push(`${esc(sourceLabel(s.source))},${s.count}`));
    rows.push('');
  }

  // ── Click heatmap ───────────────────────────────────────────────────────────
  rows.push('CLICK HEATMAP (UTC)');
  rows.push('Day,Hour,Clicks');
  a.heatmap.forEach((h) => {
    const hour = `${String(h.hour).padStart(2, '0')}:00`;
    rows.push(`${dayNames[h.day_of_week]},${hour},${h.count}`);
  });
  rows.push('');

  // ── UTM campaign tracking ────────────────────────────────────────────────────
  if (a.utm_sources?.length || a.utm_mediums?.length || a.utm_campaigns?.length || a.utm_contents?.length || a.utm_terms?.length) {
    rows.push('UTM SOURCES');
    rows.push('Source,Clicks');
    (a.utm_sources ?? []).forEach((u) => rows.push(`${esc(u.value)},${u.count}`));
    rows.push('');
    rows.push('UTM MEDIUMS');
    rows.push('Medium,Clicks');
    (a.utm_mediums ?? []).forEach((u) => rows.push(`${esc(u.value)},${u.count}`));
    rows.push('');
    rows.push('UTM CAMPAIGNS');
    rows.push('Campaign,Clicks');
    (a.utm_campaigns ?? []).forEach((u) => rows.push(`${esc(u.value)},${u.count}`));
    if (a.utm_contents?.length) {
      rows.push('');
      rows.push('UTM CONTENTS');
      rows.push('Content,Clicks');
      a.utm_contents.forEach((u) => rows.push(`${esc(u.value)},${u.count}`));
    }
    if (a.utm_terms?.length) {
      rows.push('');
      rows.push('UTM TERMS');
      rows.push('Term,Clicks');
      a.utm_terms.forEach((u) => rows.push(`${esc(u.value)},${u.count}`));
    }
  }

  // ── Trigger download ────────────────────────────────────────────────────────
  const csv = rows.join('\r\n');
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `analytics-${a.link_id.slice(0, 8)}-${filterFrom.value}-to-${filterTo.value}.csv`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

onMounted(() => {
  loadLink();
  loadAnalytics();
  loadComparison();
  connectLiveCount();
  loadWorldGeoJSON();
});

onUnmounted(() => {
  disconnectLiveCount();
});

// Computed helpers
const topReferrer = computed(() => {
  if (!analytics.value || analytics.value.referrers.length === 0) return 'N/A';
  const top = analytics.value.referrers.reduce((a, b) => (a.count >= b.count ? a : b));
  return top.referrer || 'Direct / None';
});

const topDevice = computed(() => {
  if (!analytics.value || analytics.value.devices.length === 0) return 'N/A';
  const top = analytics.value.devices.reduce((a, b) => (a.count >= b.count ? a : b));
  return top.device_type || 'Unknown';
});

const topReferrers = computed(() => {
  if (!analytics.value) return [];
  return [...analytics.value.referrers]
    .sort((a, b) => b.count - a.count)
    .slice(0, 10);
});

const totalReferrerClicks = computed(() =>
  topReferrers.value.reduce((sum, r) => sum + r.count, 0)
);

const totalDeviceClicks = computed(() => {
  if (!analytics.value) return 0;
  return analytics.value.devices.reduce((sum, d) => sum + d.count, 0);
});

function referrerPercent(count: number): number {
  if (totalReferrerClicks.value === 0) return 0;
  return Math.round((count / totalReferrerClicks.value) * 100);
}

function devicePercent(count: number): number {
  if (totalDeviceClicks.value === 0) return 0;
  return Math.round((count / totalDeviceClicks.value) * 100);
}

// Country helpers
function countryFlag(code: string): string {
  if (!code || code.length !== 2 || code === 'Unknown') return '🌐';
  try {
    return String.fromCodePoint(
      ...code.toUpperCase().split('').map((c: string) => c.charCodeAt(0) + 127397),
    );
  } catch {
    return '🌐';
  }
}

function countryDisplayName(code: string): string {
  if (!code || code === 'Unknown') return 'Unknown';
  try {
    return new Intl.DisplayNames(['en'], { type: 'region' }).of(code) ?? code;
  } catch {
    return code;
  }
}

const totalCountryClicks = computed(() =>
  (analytics.value?.countries ?? []).reduce((sum, c) => sum + c.count, 0),
);

function countryPercent(count: number): number {
  if (totalCountryClicks.value === 0) return 0;
  return Math.round((count / totalCountryClicks.value) * 100);
}

// City helpers
const totalCityClicks = computed(() =>
  (analytics.value?.cities ?? []).reduce((sum, c) => sum + c.count, 0),
);

function cityPercent(count: number): number {
  if (totalCityClicks.value === 0) return 0;
  return Math.round((count / totalCityClicks.value) * 100);
}

// Visitor Loyalty helpers
const totalTrackedVisitors = computed(() =>
  (analytics.value?.first_time_visitors ?? 0) + (analytics.value?.returning_visitors ?? 0),
);
const firstTimePercent = computed(() => {
  if (!totalTrackedVisitors.value) return 0;
  return Math.round(((analytics.value?.first_time_visitors ?? 0) / totalTrackedVisitors.value) * 100);
});
const returningPercent = computed(() => {
  if (!totalTrackedVisitors.value) return 0;
  return Math.round(((analytics.value?.returning_visitors ?? 0) / totalTrackedVisitors.value) * 100);
});

const visitorLoyaltyChartOption = computed(() => {
  const ft = analytics.value?.first_time_visitors ?? 0;
  const rv = analytics.value?.returning_visitors ?? 0;
  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e9ecef',
      borderWidth: 1,
      textStyle: { color: '#333', fontSize: 12 },
      formatter: (p: { name: string; value: number; percent: number }) =>
        `${p.name}<br/><b>${p.value.toLocaleString()}</b> visitors (${p.percent}%)`,
    },
    legend: { show: false },
    series: [{
      type: 'pie',
      radius: ['48%', '74%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: true,
      label: { show: false },
      emphasis: {
        label: { show: true, fontSize: 13, fontWeight: 'bold' },
        itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.1)' },
      },
      data: [
        { name: 'First-Time',  value: ft, itemStyle: { color: '#635bff' } },
        { name: 'Returning',   value: rv, itemStyle: { color: '#14b8a6' } },
      ],
    }],
  };
});

// Traffic Channel helpers
const CHANNEL_META: Record<string, { label: string; color: string; icon: string }> = {
  direct:   { label: 'Direct / None',    color: '#635bff', icon: '🔗' },
  search:   { label: 'Search Engines',   color: '#f59e0b', icon: '🔍' },
  social:   { label: 'Social Media',     color: '#14b8a6', icon: '👥' },
  email:    { label: 'Email',            color: '#ec4899', icon: '✉️' },
  referral: { label: 'Referral / Other', color: '#6366f1', icon: '🌐' },
};

function channelLabel(cat: string): string  { return CHANNEL_META[cat]?.label ?? cat; }
function channelColor(cat: string): string  { return CHANNEL_META[cat]?.color ?? '#aaa'; }
function channelIcon(cat: string): string   { return CHANNEL_META[cat]?.icon ?? '?'; }

const totalChannelClicks = computed(() =>
  (analytics.value?.referrer_categories ?? []).reduce((s, c) => s + c.count, 0),
);
function channelPercent(count: number): number {
  if (totalChannelClicks.value === 0) return 0;
  return Math.round((count / totalChannelClicks.value) * 100);
}

const trafficChannelsChartOption = computed(() => {
  const cats = analytics.value?.referrer_categories ?? [];
  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e9ecef',
      borderWidth: 1,
      textStyle: { color: '#333', fontSize: 12 },
      formatter: (p: { name: string; value: number; percent: number }) =>
        `${p.name}<br/><b>${p.value.toLocaleString()}</b> clicks (${p.percent}%)`,
    },
    legend: { show: false },
    series: [{
      type: 'pie',
      radius: ['48%', '78%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: true,
      label: { show: false },
      emphasis: {
        label: { show: true, fontSize: 13, fontWeight: 'bold' },
        itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.1)' },
      },
      data: cats.map((c) => ({
        name: channelLabel(c.category),
        value: c.count,
        itemStyle: { color: channelColor(c.category) },
      })),
    }],
  };
});

const countriesChartOption = computed(() => {
  const countries = analytics.value?.countries ?? [];
  const reversed = [...countries].reverse();
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e9ecef',
      borderWidth: 1,
      textStyle: { color: '#333', fontSize: 12 },
    },
    grid: { left: '3%', right: '8%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'value',
      axisLabel: { color: '#6c757d', fontSize: 11 },
      splitLine: { lineStyle: { color: '#f1f3f5', type: 'dashed' } },
    },
    yAxis: {
      type: 'category',
      data: reversed.map((c) => countryDisplayName(c.country)),
      axisLabel: { color: '#6c757d', fontSize: 11 },
      axisLine: { lineStyle: { color: '#dee2e6' } },
      axisTick: { show: false },
    },
    series: [
      {
        name: 'Clicks',
        type: 'bar',
        data: reversed.map((c) => c.count),
        barMaxWidth: 20,
        itemStyle: { color: '#635bff', borderRadius: [0, 4, 4, 0] },
      },
    ],
  };
});

const geoMapOption = computed(() => {
  const countries = analytics.value?.countries ?? [];
  const maxVal = Math.max(...countries.map((c) => c.count), 1);
  const data = countries.map((c) => ({
    name: countryToEchartsName(c.country),
    value: c.count,
  }));
  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e9ecef',
      borderWidth: 1,
      textStyle: { color: '#333', fontSize: 12 },
      formatter: (params: { name: string; value: number | string }) =>
        params.value != null && params.value !== '-'
          ? `${params.name}: <strong>${Number(params.value).toLocaleString()}</strong> click${Number(params.value) !== 1 ? 's' : ''}`
          : `${params.name}: No data`,
    },
    visualMap: {
      min: 0,
      max: maxVal,
      show: true,
      orient: 'vertical',
      left: 10,
      top: 'center',
      inRange: { color: ['#f0eeff', '#a89cff', '#635bff'] },
      text: ['High', 'Low'],
      textStyle: { color: '#6c757d', fontSize: 10 },
      itemWidth: 14,
      itemHeight: 80,
      calculable: false,
    },
    series: [
      {
        name: 'Clicks',
        type: 'map',
        map: 'world',
        roam: true,
        emphasis: {
          label: { show: true, fontSize: 11, color: '#fff' },
          itemStyle: { areaColor: '#5249e0' },
        },
        select: { disabled: true },
        itemStyle: {
          areaColor: '#f0eeff',
          borderColor: '#c9c4f7',
          borderWidth: 0.5,
        },
        data,
      },
    ],
  };
});

// Browser & OS breakdown
const chartColors = ['#635bff', '#14b8a6', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4', '#84cc16', '#f97316'];

const totalBrowserClicks = computed(() =>
  (analytics.value?.browsers ?? []).reduce((sum, b) => sum + b.count, 0),
);

function browserPercent(count: number): number {
  if (totalBrowserClicks.value === 0) return 0;
  return Math.round((count / totalBrowserClicks.value) * 100);
}

const totalOSClicks = computed(() =>
  (analytics.value?.os_breakdown ?? []).reduce((sum, o) => sum + o.count, 0),
);

function osPercent(count: number): number {
  if (totalOSClicks.value === 0) return 0;
  return Math.round((count / totalOSClicks.value) * 100);
}

function utmPercent(count: number, series: UTMPoint[] | undefined): number {
  if (!series || series.length === 0) return 0;
  const total = series.reduce((sum, u) => sum + u.count, 0);
  if (total === 0) return 0;
  return Math.round((count / total) * 100);
}

const sourceLabels: Record<string, string> = {
  web: 'Direct / Web',
  qr: 'QR Code',
  api: 'API',
};

function sourceLabel(source: string): string {
  return sourceLabels[source] ?? source;
}

function sourceIcon(source: string): string {
  if (source === 'qr') return '▣';
  if (source === 'api') return '⚙';
  return '🔗';
}

function sourcePercent(count: number): number {
  const total = (analytics.value?.sources ?? []).reduce((s, p) => s + p.count, 0);
  if (total === 0) return 0;
  return Math.round((count / total) * 100);
}

const sourceChartOption = computed(() => {
  const src = analytics.value?.sources ?? [];
  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e9ecef',
      borderWidth: 1,
      textStyle: { color: '#333', fontSize: 12 },
      formatter: (p: { name: string; value: number; percent: number }) =>
        `${p.name}: <strong>${p.value.toLocaleString()}</strong> (${p.percent}%)`,
    },
    series: [
      {
        name: 'Source',
        type: 'pie',
        radius: ['45%', '72%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        label: { show: false },
        emphasis: { label: { show: true, fontSize: 12, fontWeight: 'bold' } },
        data: src.map((s: SourcePoint, i: number) => ({
          name: sourceLabel(s.source),
          value: s.count,
          itemStyle: { color: chartColors[i % chartColors.length] },
        })),
      },
    ],
  };
});

const browserChartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    backgroundColor: 'rgba(255,255,255,0.95)',
    borderColor: '#e9ecef',
    borderWidth: 1,
    textStyle: { color: '#333', fontSize: 12 },
    formatter: '{b}: {c} ({d}%)',
  },
  series: [
    {
      name: 'Browser',
      type: 'pie',
      radius: ['45%', '72%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: false,
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 12, fontWeight: 'bold' } },
      data: (analytics.value?.browsers ?? []).map((b, i) => ({
        name: b.browser,
        value: b.count,
        itemStyle: { color: chartColors[i % chartColors.length] },
      })),
    },
  ],
}));

const osChartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    backgroundColor: 'rgba(255,255,255,0.95)',
    borderColor: '#e9ecef',
    borderWidth: 1,
    textStyle: { color: '#333', fontSize: 12 },
    formatter: '{b}: {c} ({d}%)',
  },
  series: [
    {
      name: 'OS',
      type: 'pie',
      radius: ['45%', '72%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: false,
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 12, fontWeight: 'bold' } },
      data: (analytics.value?.os_breakdown ?? []).map((o, i) => ({
        name: o.os,
        value: o.count,
        itemStyle: { color: chartColors[i % chartColors.length] },
      })),
    },
  ],
}));

// Click heatmap (hour × day-of-week)
const heatmapChartOption = computed(() => {
  const heatmap = analytics.value?.heatmap ?? [];
  const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const hours = Array.from({ length: 24 }, (_, i) => `${String(i).padStart(2, '0')}:00`);
  const maxCount = heatmap.reduce((m, p) => Math.max(m, Number(p.count)), 0);

  return {
    tooltip: {
      position: 'top',
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e9ecef',
      borderWidth: 1,
      textStyle: { color: '#333', fontSize: 12 },
      formatter: (params: { data: [number, number, number] }) => {
        const [hour, day, count] = params.data;
        return `${days[day]} ${hours[hour]}: <strong>${count}</strong> click${count !== 1 ? 's' : ''}`;
      },
    },
    grid: { left: '4%', right: '4%', bottom: '18%', top: '4%', containLabel: true },
    xAxis: {
      type: 'category',
      data: hours,
      splitArea: { show: true },
      axisLabel: { color: '#6c757d', fontSize: 9, interval: 1 },
      axisLine: { lineStyle: { color: '#dee2e6' } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'category',
      data: days,
      splitArea: { show: true },
      axisLabel: { color: '#6c757d', fontSize: 11 },
      axisLine: { lineStyle: { color: '#dee2e6' } },
      axisTick: { show: false },
    },
    visualMap: {
      min: 0,
      max: maxCount || 1,
      calculable: false,
      orient: 'horizontal',
      left: 'center',
      bottom: '2%',
      inRange: { color: ['#f0eeff', '#a89cff', '#635bff'] },
      textStyle: { color: '#6c757d', fontSize: 10 },
      itemWidth: 14,
      itemHeight: 80,
    },
    series: [
      {
        name: 'Clicks',
        type: 'heatmap',
        data: heatmap.map((p) => [p.hour, p.day_of_week, p.count]),
        label: { show: false },
        emphasis: {
          itemStyle: { shadowBlur: 8, shadowColor: 'rgba(99,91,255,0.4)' },
        },
      },
    ],
  };
});

// ECharts option
const chartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(255,255,255,0.95)',
    borderColor: '#e9ecef',
    borderWidth: 1,
    textStyle: { color: '#333', fontSize: 12 },
    axisPointer: { type: 'cross', label: { backgroundColor: '#635bff' } },
  },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: analytics.value?.time_series.map((p) => p.timestamp.substring(0, 10)) ?? [],
    axisLine: { lineStyle: { color: '#dee2e6' } },
    axisTick: { show: false },
    axisLabel: { color: '#6c757d', fontSize: 11 },
  },
  yAxis: {
    type: 'value',
    splitLine: { lineStyle: { color: '#f1f3f5', type: 'dashed' } },
    axisLabel: { color: '#6c757d', fontSize: 11 },
  },
  series: [
    {
      name: 'Clicks',
      type: 'line',
      data: analytics.value?.time_series.map((p) => p.count) ?? [],
      smooth: true,
      areaStyle: { color: 'rgba(99,91,255,0.1)' },
      itemStyle: { color: '#635bff' },
      lineStyle: { color: '#635bff', width: 2 },
      symbol: 'circle',
      symbolSize: 6,
    },
  ],
}));
</script>

<style scoped lang="scss">
.page-wrapper {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page header */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}

.page-header__left {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.page-header__title-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.page-header__subtitle {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.page-header__actions {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.short-url-monospace {
  font-family: monospace;
  font-size: 0.875rem;
  color: var(--md-sys-color-primary);
}

/* Link info card */
.link-info-card {
  border-left: 3px solid var(--md-sys-color-primary) !important;
}

.link-info-body {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.link-info-main {
  flex: 1;
  min-width: 0;
}

.link-info-stat {
  text-align: center;
  flex-shrink: 0;
}

.link-title-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.25rem;
}

.link-info-title {
  color: var(--md-sys-color-on-surface);
}

.stat-number {
  line-height: 1;
}

/* Error banner */
.error-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  background: rgba(var(--md-sys-color-error-rgb, 176, 0, 32), 0.08);
  border: 1px solid var(--md-sys-color-error);
}

/* Filters */
.filter-row {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  flex-wrap: wrap;
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

/* Stat grid */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;

  @media (max-width: 991px) { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.stat-card {
  padding: 1.25rem;
}

.stat-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-value {
  color: var(--md-sys-color-on-surface);
}

/* Comparison */
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

.trend-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.15rem;
  font-size: 0.8rem;
  font-weight: 600;
  padding: 0.15rem 0.45rem;
  border-radius: 999px;

  &.comp-badge--up   { background: rgba(22, 163, 74, 0.12);  color: #16a34a; }
  &.comp-badge--down { background: rgba(220, 38, 38, 0.12);  color: #dc2626; }
  &.comp-badge--stable { background: rgba(107, 114, 128, 0.1); color: #6b7280; }
}

.comparison-period-box {
  border-radius: 12px;
  padding: 0.75rem 1rem;

  &--current {
    background: rgba(99, 91, 255, 0.06);
    border: 1px solid rgba(99, 91, 255, 0.2);
  }

  &--previous {
    background: rgba(107, 114, 128, 0.06);
    border: 1px solid rgba(107, 114, 128, 0.15);
  }
}

.period-label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 0.2rem;
}

.period-dates {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 0.25rem;
}

/* Two column layouts */
.two-col-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;

  @media (max-width: 767px) { grid-template-columns: 1fr; }
}

.two-col-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  align-items: start;

  @media (max-width: 767px) { grid-template-columns: 1fr; }
}

/* Breakdown list */
.breakdown-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.breakdown-item {
  // each item has its own internal layout
}

/* Chart + legend */
.chart-legend-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  align-items: center;

  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.legend-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* Visitor loyalty boxes */
.loyalty-box {
  padding: 0.75rem 1rem;
  border-radius: 12px;

  &--first {
    background: rgba(99, 91, 255, 0.06);
  }
  &--return {
    background: rgba(20, 184, 166, 0.06);
  }
  &--rate {
    background: rgba(245, 158, 11, 0.06);
  }
}

/* Timeline */
.timeline-list {
  display: flex;
  flex-direction: column;
}

.timeline-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.65rem 0;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }

  &:first-child {
    padding-top: 0;
  }
}

.timeline-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 1px;

  &--created { background: rgba(99, 91, 255, 0.1); color: var(--md-sys-color-primary); }
  &--updated { background: rgba(20, 184, 166, 0.1); color: #14b8a6; }
  &--health  { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
  &--expires { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
  &--limit   { background: rgba(99, 91, 255, 0.08); color: var(--md-sys-color-primary); }
}

.timeline-body {
  flex: 1;
  min-width: 0;
}

.timeline-label {
  color: var(--md-sys-color-on-surface-variant);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: 0.15rem;
}

.timeline-date {
  color: var(--md-sys-color-on-surface);
}

/* Security grid */
.security-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem 1.5rem;
}

.security-item {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;

  &--full {
    grid-column: 1 / -1;
  }
}

.security-item-label {
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--md-sys-color-on-surface-variant);
}

/* UTM grid */
.utm-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;

  @media (max-width: 991px) { grid-template-columns: 1fr 1fr; }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

/* M3 Card header */
.m3-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1.25rem;
  gap: 1rem;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);

  &__left {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  &__right {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
}

/* Filters card */
.filters-card {
  border-left: 3px solid var(--md-sys-color-primary) !important;
}

/* M3 Table wrapper */
.m3-table-wrapper {
  overflow-x: auto;
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
  gap: 0.2rem;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;

  &--primary   { background: rgba(99, 91, 255, 0.12); color: var(--md-sys-color-primary, #635bff); }
  &--success   { background: rgba(22, 163, 74, 0.12);  color: #16a34a; }
  &--error     { background: rgba(220, 38, 38, 0.12);  color: #dc2626; }
  &--warning   { background: rgba(245, 158, 11, 0.12); color: #d97706; }
  &--neutral   { background: rgba(107, 114, 128, 0.1); color: #6b7280; }
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
    padding: 0.5rem 1rem;
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

/* Skeleton */
.skeleton {
  background: linear-gradient(90deg, var(--md-sys-color-outline-variant, #e3e8ee) 25%, #f0f2f5 50%, var(--md-sys-color-outline-variant, #e3e8ee) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.2s infinite;
  border-radius: 4px;
  display: inline-block;
}

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* Empty state */
.empty-state {
  text-align: center;
  padding: 2rem;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}

.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2.5rem 1.5rem;
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  gap: 0.25rem;
}

/* Live indicator */
.live-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.live-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #16a34a;
  animation: live-pulse 1.4s ease-in-out infinite;
}

@keyframes live-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.4; transform: scale(0.75); }
}

/* Expiry colors */
.text-danger { color: var(--md-sys-color-error); }
.text-warning { color: #d97706; }

/* copy btn success state */
.copy-btn--success {
  color: #16a34a !important;
}
</style>
