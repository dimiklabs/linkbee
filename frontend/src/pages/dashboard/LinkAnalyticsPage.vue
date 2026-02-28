<template>
  <div class="analytics-page">

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <RouterLink to="/dashboard/links">
          <button class="btn-icon">
            <span class="material-symbols-outlined">arrow_back</span>
          </button>
        </RouterLink>
        <div class="page-header__title-group">
          <h1 class="page-title">Link Analytics</h1>
          <template v-if="link">
            <div class="page-header__subtitle">
              <span class="short-url-monospace">{{ link.short_url }}</span>
              <button class="btn-icon" :class="copied ? 'copy-btn--success' : ''" @click="copyShortUrl" title="Copy short URL">
                <span class="material-symbols-outlined icon-sm">{{ copied ? 'check' : 'content_copy' }}</span>
              </button>
              <a :href="link.short_url" target="_blank" rel="noopener noreferrer">
                <button class="btn-icon" title="Open short URL">
                  <span class="material-symbols-outlined icon-sm">open_in_new</span>
                </button>
              </a>
            </div>
          </template>
        </div>
      </div>
      <div v-if="analytics" class="page-header__actions">
        <button class="btn-outlined" @click="exportToCSV">
          <span class="material-symbols-outlined">download</span>
          Export CSV
        </button>
      </div>
    </div>

    <!-- Link Info Card -->
    <div class="an-card">
      <div v-if="linkLoading" class="link-info-skeleton">
        <div class="skeleton skeleton--120"></div>
        <div class="skeleton skeleton--200"></div>
      </div>
      <div v-else-if="link" class="an-card-body">
        <div class="link-info-body">
          <div class="link-info-main">
            <div class="link-title-row">
              <span class="link-info-title">{{ link.title || '/' + link.slug }}</span>
              <span class="m3-badge" :class="link.is_active ? 'm3-badge--success' : 'm3-badge--neutral'">
                {{ link.is_active ? 'Active' : 'Inactive' }}
              </span>
              <span class="m3-badge" :class="healthBadgeClass(link.health_status)">
                {{ healthLabel(link.health_status) }}
              </span>
            </div>
            <div class="link-dest-url">
              <span class="material-symbols-outlined icon-xs">link</span>
              <span class="link-dest-url__text" :title="link.destination_url">{{ link.destination_url }}</span>
            </div>
            <div v-if="link.tags && link.tags.length" class="link-tags">
              <span v-for="tag in link.tags" :key="tag" class="m3-badge m3-badge--neutral">{{ tag }}</span>
            </div>
          </div>
          <div class="link-info-stat">
            <div class="link-info-stat__value">{{ link.click_count.toLocaleString() }}</div>
            <div class="link-info-stat__label">All-time clicks</div>
          </div>
        </div>
      </div>
      <div v-else class="an-card-body">
        <span class="link-info-title">Link Analytics</span>
        <div class="text-muted">{{ route.params.id }}</div>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="error-banner">
      <span class="material-symbols-outlined error-icon">error</span>
      <span class="error-msg">{{ error }}</span>
      <button class="btn-text" @click="loadAnalytics">Retry</button>
    </div>

    <!-- Loading state -->
    <div v-if="loading && !error" class="loading-state">
      <div class="css-spinner"></div>
      <span class="text-muted">Loading analytics data...</span>
    </div>

    <template v-if="!loading && analytics">

      <!-- Filters -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">tune</span>
            </div>
            <div class="an-card-title">Date Range &amp; Granularity</div>
          </div>
        </div>
        <div class="an-card-body">
          <div class="filter-row">
            <div class="date-field">
              <label class="date-field__label">From</label>
              <input type="date" class="date-input" :value="filterFrom" @input="filterFrom = ($event.target as HTMLInputElement).value" />
            </div>
            <div class="date-field">
              <label class="date-field__label">To</label>
              <input type="date" class="date-input" :value="filterTo" @input="filterTo = ($event.target as HTMLInputElement).value" />
            </div>
            <AppSelect
              label="Granularity"
              :model-value="filterGranularity"
              @update:model-value="filterGranularity = $event as 'hour'|'day'|'week'|'month'"
            >
              <option value="hour">Hour</option>
              <option value="day">Day</option>
              <option value="week">Week</option>
              <option value="month">Month</option>
            </AppSelect>
            <button class="btn-filled" :disabled="loading" @click="applyFilters">
              <div v-if="loading" class="css-spinner css-spinner--sm css-spinner--white"></div>
              Apply
            </button>
          </div>
        </div>
      </div>

      <!-- Period Comparison Card -->
      <div v-if="comparison" class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--neutral">
              <span class="material-symbols-outlined">compare_arrows</span>
            </div>
            <div>
              <div class="an-card-title">Period Comparison</div>
              <div class="an-card-subtitle">Current vs preceding period of equal length</div>
            </div>
          </div>
          <span class="m3-badge m3-badge--neutral">
            {{ formatShortDate(comparison.previous.from) }} → {{ formatShortDate(comparison.current.to) }}
          </span>
        </div>
        <div class="an-card-body">
          <div class="comparison-grid">
            <div class="comparison-metric">
              <div class="comparison-metric__label">Total Clicks</div>
              <div class="comparison-metric__row">
                <span class="comparison-metric__value">{{ comparison.current.total_clicks.toLocaleString() }}</span>
                <span class="trend-badge" :class="trendClass(comparison.clicks.trend)">
                  <span class="material-symbols-outlined icon-xs">{{ comparison.clicks.trend === 'up' ? 'trending_up' : comparison.clicks.trend === 'down' ? 'trending_down' : 'trending_flat' }}</span>
                  {{ Math.abs(comparison.clicks.percent_change).toFixed(1) }}%
                </span>
              </div>
              <div class="comparison-metric__prev">vs {{ comparison.previous.total_clicks.toLocaleString() }} prev. period</div>
            </div>
            <div class="comparison-metric">
              <div class="comparison-metric__label">Unique Clicks</div>
              <div class="comparison-metric__row">
                <span class="comparison-metric__value">{{ comparison.current.unique_clicks.toLocaleString() }}</span>
                <span class="trend-badge" :class="trendClass(comparison.unique_clicks.trend)">
                  <span class="material-symbols-outlined icon-xs">{{ comparison.unique_clicks.trend === 'up' ? 'trending_up' : comparison.unique_clicks.trend === 'down' ? 'trending_down' : 'trending_flat' }}</span>
                  {{ Math.abs(comparison.unique_clicks.percent_change).toFixed(1) }}%
                </span>
              </div>
              <div class="comparison-metric__prev">vs {{ comparison.previous.unique_clicks.toLocaleString() }} prev. period</div>
            </div>
            <div class="comparison-period-box comparison-period-box--current">
              <div class="period-label">Current Period</div>
              <div class="period-dates">{{ formatShortDate(comparison.current.from) }} – {{ formatShortDate(comparison.current.to) }}</div>
              <div class="period-stats">
                <span><strong>{{ comparison.current.total_clicks.toLocaleString() }}</strong> clicks</span>
                <span class="text-muted"><strong>{{ comparison.current.unique_clicks.toLocaleString() }}</strong> unique</span>
              </div>
            </div>
            <div class="comparison-period-box comparison-period-box--previous">
              <div class="period-label">Previous Period</div>
              <div class="period-dates">{{ formatShortDate(comparison.previous.from) }} – {{ formatShortDate(comparison.previous.to) }}</div>
              <div class="period-stats">
                <span><strong>{{ comparison.previous.total_clicks.toLocaleString() }}</strong> clicks</span>
                <span class="text-muted"><strong>{{ comparison.previous.unique_clicks.toLocaleString() }}</strong> unique</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Stat Cards -->
      <div class="stat-grid">
        <div class="an-card stat-card">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--primary">
              <span class="material-symbols-outlined">trending_up</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ (liveTotal !== null ? liveTotal : analytics.total_clicks).toLocaleString() }}</div>
              <div class="stat-label">
                Total Clicks
                <span v-if="isLive" class="live-badge">
                  <span class="live-dot"></span>
                  <span class="live-text">LIVE</span>
                </span>
              </div>
            </div>
          </div>
        </div>
        <div class="an-card stat-card">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--teal">
              <span class="material-symbols-outlined">group</span>
            </div>
            <div class="stat-body">
              <div class="stat-value">{{ (liveUnique !== null ? liveUnique : analytics.unique_clicks).toLocaleString() }}</div>
              <div class="stat-label">
                Unique Clicks
                <span v-if="isLive" class="live-badge">
                  <span class="live-dot"></span>
                  <span class="live-text">LIVE</span>
                </span>
              </div>
            </div>
          </div>
        </div>
        <div class="an-card stat-card">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--amber">
              <span class="material-symbols-outlined">link</span>
            </div>
            <div class="stat-body">
              <div class="stat-value stat-value--truncate" :title="topReferrer">{{ topReferrer }}</div>
              <div class="stat-label">Top Referrer</div>
            </div>
          </div>
        </div>
        <div class="an-card stat-card">
          <div class="stat-card-inner">
            <div class="stat-icon-wrap stat-icon-wrap--red">
              <span class="material-symbols-outlined">smartphone</span>
            </div>
            <div class="stat-body">
              <div class="stat-value stat-value--capitalize">{{ topDevice }}</div>
              <div class="stat-label">Top Device</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Clicks Over Time Chart -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">show_chart</span>
            </div>
            <div class="an-card-title">Clicks Over Time</div>
          </div>
        </div>
        <div class="an-card-body">
          <div v-if="analytics.time_series.length === 0" class="empty-state">
            No time series data available for this period.
          </div>
          <VChart v-else :option="chartOption" class="chart-tall" autoresize />
        </div>
      </div>

      <!-- Traffic Channels -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">account_tree</span>
            </div>
            <div>
              <div class="an-card-title">Traffic Channels</div>
              <div class="an-card-subtitle">How visitors are reaching your link — classified by referrer domain</div>
            </div>
          </div>
        </div>
        <div class="an-card-body">
          <div v-if="!analytics.referrer_categories?.length" class="empty-state">
            No channel data available yet.
          </div>
          <div v-else class="two-col-layout">
            <div class="col-min0">
              <VChart :option="trafficChannelsChartOption" class="chart-medium" autoresize />
            </div>
            <div class="breakdown-list">
              <div v-for="ch in analytics.referrer_categories" :key="ch.category" class="breakdown-item">
                <div class="breakdown-item__row">
                  <span class="channel-icon-wrap" :style="{ backgroundColor: channelColor(ch.category) + '22' }">
                    <span class="channel-icon-text">{{ channelIcon(ch.category) }}</span>
                  </span>
                  <div class="breakdown-item__content">
                    <div class="breakdown-item__meta">
                      <span class="breakdown-item__name">{{ channelLabel(ch.category) }}</span>
                      <span class="breakdown-item__count">{{ ch.count.toLocaleString() }} ({{ channelPercent(ch.count) }}%)</span>
                    </div>
                    <div class="prog-bar">
                      <div class="prog-fill" :style="{ width: channelPercent(ch.count) + '%' }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Two-column: Referrers + Devices -->
      <div class="two-col-grid">
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-header__left">
              <div class="an-card-icon an-card-icon--neutral">
                <span class="material-symbols-outlined">link</span>
              </div>
              <div class="an-card-title">Top Referrers</div>
            </div>
            <span class="m3-badge m3-badge--neutral">Top 10</span>
          </div>
          <div v-if="topReferrers.length === 0" class="m3-empty-state">
            <div class="m3-empty-state__icon">
              <span class="material-symbols-outlined">link_off</span>
            </div>
            <p class="m3-empty-state__text">No referrer data available yet.</p>
          </div>
          <div v-else class="m3-table-wrapper">
            <table class="m3-table">
              <thead>
                <tr>
                  <th>Referrer</th>
                  <th class="th-right">Clicks</th>
                  <th class="th-share">Share</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="ref in topReferrers" :key="ref.referrer">
                  <td class="td-truncate">
                    <span class="truncate-text" :title="ref.referrer">{{ ref.referrer || 'Direct / None' }}</span>
                  </td>
                  <td class="td-right td-bold">{{ ref.count.toLocaleString() }}</td>
                  <td>
                    <div class="prog-with-pct">
                      <div class="prog-bar prog-bar--flex">
                        <div class="prog-fill" :style="{ width: referrerPercent(ref.count) + '%' }"></div>
                      </div>
                      <span class="prog-pct">{{ referrerPercent(ref.count) }}%</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-header__left">
              <div class="an-card-icon an-card-icon--neutral">
                <span class="material-symbols-outlined">devices</span>
              </div>
              <div class="an-card-title">Device Breakdown</div>
            </div>
            <span class="m3-badge m3-badge--neutral">All Devices</span>
          </div>
          <div v-if="analytics.devices.length === 0" class="m3-empty-state">
            <div class="m3-empty-state__icon">
              <span class="material-symbols-outlined">devices</span>
            </div>
            <p class="m3-empty-state__text">No device data available yet.</p>
          </div>
          <div v-else class="m3-table-wrapper">
            <table class="m3-table">
              <thead>
                <tr>
                  <th>Device</th>
                  <th class="th-right">Clicks</th>
                  <th class="th-share">Share</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="device in analytics.devices" :key="device.device_type">
                  <td>
                    <div class="device-cell">
                      <span class="material-symbols-outlined device-icon">
                        {{ device.device_type.toLowerCase() === 'mobile' ? 'smartphone' : device.device_type.toLowerCase() === 'tablet' ? 'tablet' : 'computer' }}
                      </span>
                      <span class="device-label">{{ device.device_type || 'Unknown' }}</span>
                    </div>
                  </td>
                  <td class="td-right td-bold">{{ device.count.toLocaleString() }}</td>
                  <td>
                    <div class="prog-with-pct">
                      <div class="prog-bar prog-bar--flex">
                        <div class="prog-fill" :style="{ width: devicePercent(device.count) + '%' }"></div>
                      </div>
                      <span class="prog-pct">{{ devicePercent(device.count) }}%</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Two-column: Browser + OS -->
      <div class="two-col-grid">
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-header__left">
              <div class="an-card-icon an-card-icon--neutral">
                <span class="material-symbols-outlined">globe</span>
              </div>
              <div class="an-card-title">Browser Breakdown</div>
            </div>
            <span class="m3-badge m3-badge--neutral">All Browsers</span>
          </div>
          <div class="an-card-body">
            <div v-if="!analytics.browsers || analytics.browsers.length === 0" class="empty-state">
              No browser data available.
            </div>
            <div v-else class="chart-legend-layout">
              <VChart :option="browserChartOption" class="chart-small" autoresize />
              <div class="legend-list">
                <div v-for="(b, i) in analytics.browsers" :key="b.browser" class="legend-item">
                  <span class="legend-dot" :style="{ backgroundColor: chartColors[i % chartColors.length] }"></span>
                  <span class="legend-name" :title="b.browser">{{ b.browser }}</span>
                  <span class="legend-pct">{{ browserPercent(b.count) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-header__left">
              <div class="an-card-icon an-card-icon--neutral">
                <span class="material-symbols-outlined">computer</span>
              </div>
              <div class="an-card-title">Operating System</div>
            </div>
            <span class="m3-badge m3-badge--neutral">All OS</span>
          </div>
          <div class="an-card-body">
            <div v-if="!analytics.os_breakdown || analytics.os_breakdown.length === 0" class="empty-state">
              No OS data available.
            </div>
            <div v-else class="chart-legend-layout">
              <VChart :option="osChartOption" class="chart-small" autoresize />
              <div class="legend-list">
                <div v-for="(o, i) in analytics.os_breakdown" :key="o.os" class="legend-item">
                  <span class="legend-dot" :style="{ backgroundColor: chartColors[i % chartColors.length] }"></span>
                  <span class="legend-name" :title="o.os">{{ o.os }}</span>
                  <span class="legend-pct">{{ osPercent(o.count) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Source Breakdown -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--neutral">
              <span class="material-symbols-outlined">share</span>
            </div>
            <div>
              <div class="an-card-title">Click Source</div>
              <div class="an-card-subtitle">How visitors arrived — direct web link, QR code, or API call</div>
            </div>
          </div>
          <span class="m3-badge m3-badge--neutral">All Sources</span>
        </div>
        <div class="an-card-body">
          <div v-if="!analytics.sources || analytics.sources.length === 0" class="empty-state">
            No source data available.
          </div>
          <div v-else class="chart-legend-layout">
            <VChart :option="sourceChartOption" class="chart-small" autoresize />
            <table class="m3-table source-table">
              <tbody>
                <tr v-for="(s, i) in analytics.sources" :key="s.source">
                  <td>
                    <div class="source-cell">
                      <span class="legend-dot" :style="{ backgroundColor: chartColors[i % chartColors.length] }"></span>
                      <span>{{ sourceIcon(s.source) }} {{ sourceLabel(s.source) }}</span>
                    </div>
                  </td>
                  <td class="td-right td-bold">{{ s.count.toLocaleString() }}</td>
                  <td class="td-right text-muted">{{ sourcePercent(s.count) }}%</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Visitor Loyalty -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--teal">
              <span class="material-symbols-outlined">favorite</span>
            </div>
            <div>
              <div class="an-card-title">Visitor Loyalty</div>
              <div class="an-card-subtitle">First-time vs returning visitors, based on hashed IP tracking</div>
            </div>
          </div>
        </div>
        <div class="an-card-body">
          <div v-if="!analytics.first_time_visitors && !analytics.returning_visitors" class="empty-state">
            No visitor loyalty data yet. Data is collected on new clicks.
          </div>
          <div v-else class="chart-legend-layout">
            <VChart :option="visitorLoyaltyChartOption" class="chart-small" autoresize />
            <div class="loyalty-boxes">
              <div class="loyalty-box loyalty-box--first">
                <div class="loyalty-box__label">First-Time</div>
                <div class="loyalty-box__value loyalty-box__value--primary">{{ analytics.first_time_visitors.toLocaleString() }}</div>
                <div class="loyalty-box__sub">{{ firstTimePercent }}% of tracked</div>
              </div>
              <div class="loyalty-box loyalty-box--return">
                <div class="loyalty-box__label">Returning</div>
                <div class="loyalty-box__value loyalty-box__value--teal">{{ analytics.returning_visitors.toLocaleString() }}</div>
                <div class="loyalty-box__sub">{{ returningPercent }}% of tracked</div>
              </div>
              <div class="loyalty-box loyalty-box--rate">
                <div class="loyalty-box__label">Return Rate</div>
                <div class="loyalty-box__value loyalty-box__value--amber">{{ returningPercent }}%</div>
                <div class="loyalty-box__sub">visitors who came back</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Link Timeline Card -->
      <div v-if="link" class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">timeline</span>
            </div>
            <div>
              <div class="an-card-title">Link Timeline</div>
              <div class="an-card-subtitle">Key dates and milestones for this link</div>
            </div>
          </div>
        </div>
        <div class="an-card-body an-card-body--tight">
          <div class="timeline-list">
            <div class="timeline-item">
              <div class="timeline-icon timeline-icon--created">
                <span class="material-symbols-outlined icon-sm">celebration</span>
              </div>
              <div class="timeline-body">
                <div class="timeline-label">Created</div>
                <div class="timeline-date">{{ formatLinkDate(link.created_at) }}</div>
              </div>
            </div>
            <div v-if="link.updated_at && link.updated_at !== link.created_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--updated">
                <span class="material-symbols-outlined icon-sm">edit</span>
              </div>
              <div class="timeline-body">
                <div class="timeline-label">Last Updated</div>
                <div class="timeline-date">{{ formatLinkDate(link.updated_at) }}</div>
              </div>
            </div>
            <div v-if="link.health_checked_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--health">
                <span class="material-symbols-outlined icon-sm">search</span>
              </div>
              <div class="timeline-body">
                <div class="timeline-label-row">
                  <div class="timeline-label">Health Last Checked</div>
                  <span class="m3-badge" :class="healthBadgeClass(link.health_status)">{{ healthLabel(link.health_status) }}</span>
                </div>
                <div class="timeline-date">{{ formatLinkDate(link.health_checked_at) }}</div>
              </div>
            </div>
            <div v-if="link.expires_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--expires">
                <span class="material-symbols-outlined icon-sm">schedule</span>
              </div>
              <div class="timeline-body">
                <div class="timeline-label">Expires</div>
                <div class="timeline-date" :class="expiryDateClass(link.expires_at)">
                  {{ formatLinkDate(link.expires_at) }}
                  <span class="timeline-date__sub" :class="expiryDateClass(link.expires_at)">{{ expiryCountdown(link.expires_at) }}</span>
                </div>
              </div>
            </div>
            <div v-if="link.max_clicks" class="timeline-item">
              <div class="timeline-icon timeline-icon--limit">
                <span class="material-symbols-outlined icon-sm">flag</span>
              </div>
              <div class="timeline-body">
                <div class="timeline-label">Click Limit</div>
                <div class="timeline-date">
                  Max clicks: <strong>{{ link.max_clicks.toLocaleString() }}</strong>
                  <span class="text-muted"> ({{ link.click_count.toLocaleString() }} used, {{ Math.max(0, link.max_clicks - link.click_count).toLocaleString() }} remaining)</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Security & Limits Card -->
      <div v-if="link" class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--neutral">
              <span class="material-symbols-outlined">shield</span>
            </div>
            <div>
              <div class="an-card-title">Security &amp; Limits</div>
              <div class="an-card-subtitle">Access controls and usage limits for this link</div>
            </div>
          </div>
        </div>
        <div class="an-card-body">
          <div class="security-grid">
            <div class="security-item">
              <div class="security-item__label">Password Protected</div>
              <div>
                <span v-if="link.has_password" class="m3-badge m3-badge--warning">
                  <span class="material-symbols-outlined icon-xs">lock</span> Yes
                </span>
                <span v-else class="m3-badge m3-badge--neutral">No</span>
              </div>
            </div>
            <div class="security-item">
              <div class="security-item__label">Active Status</div>
              <span class="m3-badge" :class="link.is_active ? 'm3-badge--success' : 'm3-badge--neutral'">
                {{ link.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>
            <div class="security-item">
              <div class="security-item__label">Redirect Type</div>
              <span class="m3-badge m3-badge--neutral">
                {{ link.redirect_type === 301 ? '301 Permanent' : '302 Temporary' }}
              </span>
            </div>
            <div class="security-item">
              <div class="security-item__label">Expiry</div>
              <span v-if="link.expires_at" :class="expiryDateClass(link.expires_at)" class="security-item__value">
                {{ formatLinkDate(link.expires_at) }}
                <span class="security-item__sub">{{ expiryCountdown(link.expires_at) }}</span>
              </span>
              <span v-else class="text-muted">Never expires</span>
            </div>
            <div class="security-item security-item--full">
              <div class="security-item__label">Click Limit</div>
              <div v-if="link.max_clicks">
                <div class="click-limit-row">
                  <span>{{ link.click_count.toLocaleString() }} / {{ link.max_clicks.toLocaleString() }} clicks used</span>
                  <span class="text-muted">{{ clickLimitPercent(link.click_count, link.max_clicks) }}%</span>
                </div>
                <div class="prog-bar prog-bar--mt">
                  <div class="prog-fill" :style="{ width: clickLimitPercent(link.click_count, link.max_clicks) + '%', backgroundColor: clickLimitBarColor(link.click_count, link.max_clicks) }"></div>
                </div>
              </div>
              <div v-else class="text-muted">Unlimited</div>
            </div>
            <div class="security-item security-item--full">
              <div class="security-item__label">Tags</div>
              <div v-if="link.tags && link.tags.length" class="link-tags">
                <span v-for="tag in link.tags" :key="tag" class="m3-badge m3-badge--neutral">{{ tag }}</span>
              </div>
              <span v-else class="text-muted">None</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Heatmap -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--amber">
              <span class="material-symbols-outlined">grid_view</span>
            </div>
            <div>
              <div class="an-card-title">Click Heatmap</div>
              <div class="an-card-subtitle">Clicks by hour of day × day of week (UTC)</div>
            </div>
          </div>
          <span class="m3-badge m3-badge--neutral">24 × 7</span>
        </div>
        <div class="an-card-body">
          <div v-if="!analytics.heatmap || analytics.heatmap.length === 0" class="empty-state">
            No heatmap data available for this period.
          </div>
          <VChart v-else :option="heatmapChartOption" class="chart-heatmap" autoresize />
        </div>
      </div>

      <!-- Geographic World Map -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">public</span>
            </div>
            <div>
              <div class="an-card-title">Geographic Distribution</div>
              <div class="an-card-subtitle">Click density by country — scroll to zoom, drag to pan</div>
            </div>
          </div>
          <span class="m3-badge m3-badge--neutral">World Map</span>
        </div>
        <div class="an-card-body">
          <div v-if="worldMapLoading" class="loading-state loading-state--inline">
            <div class="css-spinner"></div>
            <span class="text-muted">Loading world map…</span>
          </div>
          <div v-else-if="!mapLoaded || !analytics?.countries?.length" class="empty-state">
            No geographic data available yet. Geographic data is collected on new clicks.
          </div>
          <VChart v-else :option="geoMapOption" class="chart-map" autoresize />
        </div>
      </div>

      <!-- Top Countries -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">flag</span>
            </div>
            <div class="an-card-title">Top Countries</div>
          </div>
          <span class="m3-badge m3-badge--neutral">Top 15</span>
        </div>
        <div class="an-card-body">
          <div v-if="!analytics.countries || analytics.countries.length === 0" class="empty-state">
            No country data available yet. Geographic data is collected on new clicks.
          </div>
          <div v-else class="two-col-layout">
            <div class="col-min0">
              <VChart :option="countriesChartOption" class="chart-countries" autoresize />
            </div>
            <div class="col-min0">
              <table class="m3-table">
                <thead>
                  <tr>
                    <th>Country</th>
                    <th class="th-right">Clicks</th>
                    <th class="th-share">Share</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="c in analytics.countries" :key="c.country">
                    <td>
                      <div class="country-cell">
                        <span class="country-flag">{{ countryFlag(c.country) }}</span>
                        <span>{{ countryDisplayName(c.country) }}</span>
                      </div>
                    </td>
                    <td class="td-right td-bold">{{ c.count.toLocaleString() }}</td>
                    <td>
                      <div class="prog-with-pct">
                        <div class="prog-bar prog-bar--flex">
                          <div class="prog-fill" :style="{ width: countryPercent(c.count) + '%' }"></div>
                        </div>
                        <span class="prog-pct">{{ countryPercent(c.count) }}%</span>
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
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--neutral">
              <span class="material-symbols-outlined">location_city</span>
            </div>
            <div class="an-card-title">Top Cities</div>
          </div>
          <span class="m3-badge m3-badge--neutral">Top 20</span>
        </div>
        <div v-if="!analytics.cities || analytics.cities.length === 0" class="empty-state">
          No city data available yet. City data is collected on new clicks.
        </div>
        <div v-else class="m3-table-wrapper">
          <table class="m3-table">
            <thead>
              <tr>
                <th class="th-num">#</th>
                <th>City</th>
                <th>Country</th>
                <th class="th-right">Clicks</th>
                <th class="th-share">Share</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(c, idx) in analytics.cities" :key="c.city + c.country">
                <td class="td-num text-muted">{{ idx + 1 }}</td>
                <td>{{ c.city }}</td>
                <td>
                  <div class="country-cell">
                    <span class="country-flag">{{ countryFlag(c.country) }}</span>
                    <span class="text-muted">{{ countryDisplayName(c.country) }}</span>
                  </div>
                </td>
                <td class="td-right td-bold">{{ c.count.toLocaleString() }}</td>
                <td>
                  <div class="prog-with-pct">
                    <div class="prog-bar prog-bar--flex">
                      <div class="prog-fill" :style="{ width: cityPercent(c.count) + '%' }"></div>
                    </div>
                    <span class="prog-pct">{{ cityPercent(c.count) }}%</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- UTM Campaign Tracking -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-header__left">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">campaign</span>
            </div>
            <div>
              <div class="an-card-title">UTM Campaign Tracking</div>
              <div class="an-card-subtitle">Clicks from links with UTM parameters appended</div>
            </div>
          </div>
        </div>
        <div class="an-card-body">
          <div
            v-if="!analytics.utm_sources?.length && !analytics.utm_mediums?.length && !analytics.utm_campaigns?.length && !analytics.utm_contents?.length && !analytics.utm_terms?.length"
            class="empty-state"
          >
            No UTM data yet. Append <code>?utm_source=email&amp;utm_medium=newsletter&amp;utm_campaign=launch</code> to your short link when sharing.
          </div>
          <div v-else>
            <div class="utm-grid">
              <div>
                <div class="utm-group-label">Source</div>
                <div v-if="!analytics.utm_sources?.length" class="text-muted">No data</div>
                <div v-else class="breakdown-list">
                  <div v-for="u in analytics.utm_sources" :key="u.value" class="breakdown-item">
                    <div class="breakdown-item__meta">
                      <span class="breakdown-item__name" :title="u.value">{{ u.value }}</span>
                      <span class="breakdown-item__count">{{ u.count.toLocaleString() }}</span>
                    </div>
                    <div class="prog-bar">
                      <div class="prog-fill" :style="{ width: utmPercent(u.count, analytics.utm_sources) + '%' }"></div>
                    </div>
                  </div>
                </div>
              </div>
              <div>
                <div class="utm-group-label">Medium</div>
                <div v-if="!analytics.utm_mediums?.length" class="text-muted">No data</div>
                <div v-else class="breakdown-list">
                  <div v-for="u in analytics.utm_mediums" :key="u.value" class="breakdown-item">
                    <div class="breakdown-item__meta">
                      <span class="breakdown-item__name" :title="u.value">{{ u.value }}</span>
                      <span class="breakdown-item__count">{{ u.count.toLocaleString() }}</span>
                    </div>
                    <div class="prog-bar">
                      <div class="prog-fill" :style="{ width: utmPercent(u.count, analytics.utm_mediums) + '%' }"></div>
                    </div>
                  </div>
                </div>
              </div>
              <div>
                <div class="utm-group-label">Campaign</div>
                <div v-if="!analytics.utm_campaigns?.length" class="text-muted">No data</div>
                <div v-else class="breakdown-list">
                  <div v-for="u in analytics.utm_campaigns" :key="u.value" class="breakdown-item">
                    <div class="breakdown-item__meta">
                      <span class="breakdown-item__name" :title="u.value">{{ u.value }}</span>
                      <span class="breakdown-item__count">{{ u.count.toLocaleString() }}</span>
                    </div>
                    <div class="prog-bar">
                      <div class="prog-fill" :style="{ width: utmPercent(u.count, analytics.utm_campaigns) + '%' }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="analytics.utm_contents?.length || analytics.utm_terms?.length" class="utm-section2">
              <div class="utm-divider"></div>
              <div class="two-col-grid">
                <div>
                  <div class="utm-group-label">Content</div>
                  <div v-if="!analytics.utm_contents?.length" class="text-muted">No data</div>
                  <div v-else class="breakdown-list">
                    <div v-for="u in analytics.utm_contents" :key="u.value" class="breakdown-item">
                      <div class="breakdown-item__meta">
                        <span class="breakdown-item__name" :title="u.value">{{ u.value }}</span>
                        <span class="breakdown-item__count">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="prog-bar">
                        <div class="prog-fill" :style="{ width: utmPercent(u.count, analytics.utm_contents) + '%' }"></div>
                      </div>
                    </div>
                  </div>
                </div>
                <div>
                  <div class="utm-group-label">Term</div>
                  <div v-if="!analytics.utm_terms?.length" class="text-muted">No data</div>
                  <div v-else class="breakdown-list">
                    <div v-for="u in analytics.utm_terms" :key="u.value" class="breakdown-item">
                      <div class="breakdown-item__meta">
                        <span class="breakdown-item__name" :title="u.value">{{ u.value }}</span>
                        <span class="breakdown-item__count">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="prog-bar">
                        <div class="prog-fill" :style="{ width: utmPercent(u.count, analytics.utm_terms) + '%' }"></div>
                      </div>
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
    if (res.data) link.value = res.data;
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
    analytics.value = (await linksApi.getAnalytics(
      id,
      filterFrom.value ? new Date(filterFrom.value).toISOString() : undefined,
      filterTo.value ? new Date(filterTo.value).toISOString() : undefined,
      filterGranularity.value,
    )).data ?? null;
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
    comparison.value = res.data ?? null;
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
/* ── Root ─────────────────────────────────────────────────────────── */
.analytics-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0;
}

/* ── Page header ─────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
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

.page-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0;
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

.icon-sm { font-size: 18px; }
.icon-xs { font-size: 14px; }

/* ── Cards ───────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-radius: 14px;
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  gap: 0.75rem;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);

  &__left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.an-card-subtitle {
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 1px;
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
    background: rgba(99, 91, 255, 0.1);
    color: var(--md-sys-color-primary, #635bff);
  }
  &--teal {
    background: rgba(20, 184, 166, 0.1);
    color: #14b8a6;
  }
  &--amber {
    background: rgba(245, 158, 11, 0.1);
    color: #f59e0b;
  }
  &--neutral {
    background: var(--md-sys-color-surface-container, #f3f4f6);
    color: var(--md-sys-color-on-surface-variant);
  }
}

.an-card-body {
  padding: 16px 20px;

  &--tight {
    padding: 8px 20px;
  }
}

/* ── Link info card ──────────────────────────────────────────────── */
.link-info-skeleton {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 16px 20px;
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

.link-info-stat__value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--md-sys-color-primary);
  letter-spacing: -0.02em;
  line-height: 1;
}

.link-info-stat__label {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 4px;
}

.link-title-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.25rem;
}

.link-info-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.link-dest-url {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.8125rem;
  margin-top: 0.25rem;
  overflow: hidden;

  &__text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.link-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-top: 0.5rem;
}

/* ── Error banner ────────────────────────────────────────────────── */
.error-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  background: rgba(220, 38, 38, 0.08);
  border: 1px solid rgba(220, 38, 38, 0.3);
}

.error-icon {
  color: var(--md-sys-color-error);
  flex-shrink: 0;
}

.error-msg {
  flex: 1;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

/* ── Loading state ───────────────────────────────────────────────── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 0;
  gap: 1rem;

  &--inline {
    padding: 3rem 0;
    flex-direction: row;
    gap: 0.75rem;
  }
}

/* ── CSS spinner ─────────────────────────────────────────────────── */
.css-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-top-color: var(--md-sys-color-primary, #635bff);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm {
    width: 16px;
    height: 16px;
    border-width: 2px;
  }

  &--white {
    border-color: rgba(255, 255, 255, 0.35);
    border-top-color: #fff;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Date input ──────────────────────────────────────────────────── */
.date-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.date-field__label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
}

.date-input {
  height: 40px;
  padding: 0 12px;
  border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.875rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.15s;
  min-width: 150px;

  &:focus {
    border-color: var(--md-sys-color-primary, #635bff);
  }
}

/* ── Filter row ──────────────────────────────────────────────────── */
.filter-row {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  flex-wrap: wrap;
}

/* ── Stat grid ───────────────────────────────────────────────────── */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;

  @media (max-width: 991px) { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.stat-card {
  transition: box-shadow 0.15s, border-color 0.15s;

  &:hover {
    border-color: var(--md-sys-color-outline, #ccc);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
}

.stat-card-inner {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
}

.stat-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 22px; }

  &--primary { background: rgba(99, 91, 255, 0.12); color: var(--md-sys-color-primary, #635bff); }
  &--teal    { background: rgba(20, 184, 166, 0.12); color: #14b8a6; }
  &--amber   { background: rgba(245, 158, 11, 0.12); color: #f59e0b; }
  &--red     { background: rgba(239, 68, 68, 0.12);  color: #ef4444; }
}

.stat-body {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  letter-spacing: -0.02em;
  line-height: 1.1;

  &--truncate {
    font-size: 1rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &--capitalize {
    font-size: 1rem;
    text-transform: capitalize;
  }
}

.stat-label {
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* ── Period comparison ───────────────────────────────────────────── */
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

.comparison-metric__label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 0.25rem;
}

.comparison-metric__row {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.comparison-metric__value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  letter-spacing: -0.02em;
}

.comparison-metric__prev {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 0.2rem;
}

.trend-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.15rem;
  font-size: 0.8rem;
  font-weight: 600;
  padding: 0.15rem 0.45rem;
  border-radius: 999px;

  &.comp-badge--up     { background: rgba(22, 163, 74, 0.12);  color: #16a34a; }
  &.comp-badge--down   { background: rgba(220, 38, 38, 0.12);  color: #dc2626; }
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

.period-stats {
  display: flex;
  gap: 1rem;
  font-size: 0.8125rem;
}

/* ── Chart dimensions ────────────────────────────────────────────── */
.chart-tall     { height: 320px; display: block; }
.chart-medium   { height: 280px; display: block; }
.chart-small    { height: 200px; display: block; }
.chart-heatmap  { height: 260px; display: block; }
.chart-map      { height: 420px; display: block; }
.chart-countries { height: 360px; display: block; }

/* ── Two column layouts ──────────────────────────────────────────── */
.two-col-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  @media (max-width: 767px) { grid-template-columns: 1fr; }
}

.two-col-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  align-items: start;

  @media (max-width: 767px) { grid-template-columns: 1fr; }
}

.col-min0 { min-width: 0; }

/* ── Progress bars ───────────────────────────────────────────────── */
.prog-bar {
  height: 6px;
  background: var(--md-sys-color-surface-container-high, #e9ecef);
  border-radius: 999px;
  overflow: hidden;

  &--flex { flex: 1; }
  &--mt   { margin-top: 6px; }
}

.prog-fill {
  height: 100%;
  background: var(--md-sys-color-primary, #635bff);
  border-radius: 999px;
  transition: width 0.4s ease;
}

.prog-with-pct {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.prog-pct {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  min-width: 30px;
  text-align: right;
}

/* ── Breakdown list ──────────────────────────────────────────────── */
.breakdown-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.breakdown-item {
  &__row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  &__content {
    flex: 1;
    min-width: 0;
  }

  &__meta {
    display: flex;
    justify-content: space-between;
    margin-bottom: 4px;
  }

  &__name {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 180px;
  }

  &__count {
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface);
    white-space: nowrap;
  }
}

.channel-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
}

.channel-icon-text { font-size: 0.95rem; }

/* ── Chart + legend ──────────────────────────────────────────────── */
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

.legend-name {
  flex: 1;
  font-size: 0.8125rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--md-sys-color-on-surface);
}

.legend-pct {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

/* ── Table helpers ───────────────────────────────────────────────── */
.th-right  { text-align: right; }
.th-share  { width: 120px; }
.th-num    { width: 36px; }
.td-right  { text-align: right; }
.td-bold   { font-weight: 600; }
.td-num    { font-size: 0.8125rem; }
.td-truncate { max-width: 200px; }

.truncate-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}

.device-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.device-icon {
  font-size: 16px;
  color: var(--md-sys-color-primary);
}

.device-label { text-transform: capitalize; }

.country-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.country-flag {
  font-size: 1.1rem;
  line-height: 1;
}

.source-table { flex: 1; }

.source-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* ── Visitor loyalty ─────────────────────────────────────────────── */
.loyalty-boxes {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 0.75rem;
  align-content: center;
}

.loyalty-box {
  padding: 0.75rem 1rem;
  border-radius: 12px;

  &--first  { background: rgba(99, 91, 255, 0.06); }
  &--return { background: rgba(20, 184, 166, 0.06); }
  &--rate   { background: rgba(245, 158, 11, 0.06); }

  &__label {
    font-size: 0.8rem;
    color: var(--md-sys-color-on-surface-variant);
    margin-bottom: 0.25rem;
  }

  &__value {
    font-size: 1.5rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    line-height: 1.1;

    &--primary { color: var(--md-sys-color-primary, #635bff); }
    &--teal    { color: #14b8a6; }
    &--amber   { color: #f59e0b; }
  }

  &__sub {
    font-size: 0.75rem;
    color: var(--md-sys-color-on-surface-variant);
    margin-top: 2px;
  }
}

/* ── Timeline ────────────────────────────────────────────────────── */
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

  &:last-child { border-bottom: none; padding-bottom: 0; }
  &:first-child { padding-top: 0; }
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

  &--created { background: rgba(99, 91, 255, 0.1);  color: var(--md-sys-color-primary); }
  &--updated { background: rgba(20, 184, 166, 0.1); color: #14b8a6; }
  &--health  { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
  &--expires { background: rgba(239, 68, 68, 0.1);  color: #ef4444; }
  &--limit   { background: rgba(99, 91, 255, 0.08); color: var(--md-sys-color-primary); }
}

.timeline-body {
  flex: 1;
  min-width: 0;
}

.timeline-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 0.15rem;
}

.timeline-label-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.15rem;
}

.timeline-date {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);

  &__sub {
    font-size: 0.8rem;
    margin-left: 4px;
  }
}

/* ── Security grid ───────────────────────────────────────────────── */
.security-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem 1.5rem;
}

.security-item {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;

  &--full { grid-column: 1 / -1; }

  &__label {
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--md-sys-color-on-surface-variant);
  }

  &__value {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface);
  }

  &__sub {
    font-size: 0.8rem;
    margin-left: 4px;
  }
}

.click-limit-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

/* ── UTM ─────────────────────────────────────────────────────────── */
.utm-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;

  @media (max-width: 991px) { grid-template-columns: 1fr 1fr; }
  @media (max-width: 575px) { grid-template-columns: 1fr; }
}

.utm-group-label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 0.75rem;
}

.utm-section2 { margin-top: 1.5rem; }

.utm-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  margin-bottom: 1.5rem;
}

/* ── Badges ──────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.2rem 0.6rem;
  border-radius: 6px;

  &--primary { background: rgba(99, 91, 255, 0.12);  color: var(--md-sys-color-primary, #635bff); }
  &--success { background: rgba(22, 163, 74, 0.12);   color: #16a34a; }
  &--error   { background: rgba(220, 38, 38, 0.12);   color: #dc2626; }
  &--warning { background: rgba(245, 158, 11, 0.12);  color: #d97706; }
  &--neutral { background: rgba(107, 114, 128, 0.1);  color: #6b7280; }
}

/* ── M3 Table ────────────────────────────────────────────────────── */
.m3-table-wrapper { overflow-x: auto; }

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
    background: var(--md-sys-color-surface-container-low, #f8f9fa);
    border-bottom: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
    white-space: nowrap;
  }

  td {
    padding: 8px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant, #f0f2f5);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td { border-bottom: none; }
  tbody tr:hover td { background: var(--md-sys-color-surface-container-low, #f8f9fa); }
}

/* ── Empty states ────────────────────────────────────────────────── */
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
  gap: 0.5rem;

  &__icon {
    width: 48px;
    height: 48px;
    border-radius: 14px;
    background: var(--md-sys-color-surface-container);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 0.25rem;

    .material-symbols-outlined {
      font-size: 24px;
      color: var(--md-sys-color-on-surface-variant);
      opacity: 0.6;
    }
  }

  &__text {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface-variant);
    margin: 0;
  }
}

/* ── Skeleton ────────────────────────────────────────────────────── */
.skeleton {
  background: linear-gradient(90deg, var(--md-sys-color-outline-variant, #e3e8ee) 25%, #f0f2f5 50%, var(--md-sys-color-outline-variant, #e3e8ee) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.2s infinite;
  border-radius: 4px;
  display: inline-block;

  &--120 { width: 120px; height: 1.1rem; }
  &--200 { width: 200px; height: 0.9rem; }
}

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ── Live indicator ──────────────────────────────────────────────── */
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

.live-text {
  font-size: 0.6rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  color: #16a34a;
}

@keyframes live-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50%      { opacity: 0.4; transform: scale(0.75); }
}

/* ── Utility ─────────────────────────────────────────────────────── */
.text-muted   { color: var(--md-sys-color-on-surface-variant); font-size: 0.8125rem; }
.text-danger  { color: var(--md-sys-color-error); }
.text-warning { color: #d97706; }

.copy-btn--success { color: #16a34a !important; }
</style>
