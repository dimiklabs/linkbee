<template>
  <div class="container-fluid py-4">
    <!-- Back navigation -->
    <div class="mb-4">
      <RouterLink
        to="/dashboard/links"
        class="btn btn-sm btn-outline-secondary d-inline-flex align-items-center gap-2"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
          <path fill-rule="evenodd" d="M15 8a.5.5 0 0 0-.5-.5H2.707l3.147-3.146a.5.5 0 1 0-.708-.708l-4 4a.5.5 0 0 0 0 .708l4 4a.5.5 0 0 0 .708-.708L2.707 8.5H14.5A.5.5 0 0 0 15 8"/>
        </svg>
        Back to Links
      </RouterLink>
    </div>

    <!-- Link Info Card -->
    <div class="card border-0 shadow-sm mb-4 link-info-card">
      <!-- Skeleton while loading -->
      <div v-if="linkLoading" class="card-body d-flex align-items-center gap-3 py-3">
        <div class="skeleton" style="width:120px;height:1.1rem;border-radius:4px;"></div>
        <div class="skeleton" style="width:200px;height:0.9rem;border-radius:4px;"></div>
      </div>

      <!-- Populated -->
      <div v-else-if="link" class="card-body py-3 px-4">
        <div class="d-flex align-items-start justify-content-between flex-wrap gap-3">
          <!-- Left: title + short URL + destination -->
          <div class="min-w-0 flex-grow-1">
            <div class="d-flex align-items-center gap-2 flex-wrap mb-1">
              <h5 class="mb-0 fw-bold link-info-title">{{ link.title || '/' + link.slug }}</h5>
              <!-- Active/Inactive badge -->
              <span class="badge rounded-pill px-2 py-1" :class="link.is_active ? 'text-bg-success' : 'bg-secondary text-white'" style="font-size:0.7rem;">
                {{ link.is_active ? 'Active' : 'Inactive' }}
              </span>
              <!-- Health badge -->
              <span class="badge rounded-pill px-2 py-1 badge-health" :class="healthClass(link.health_status)" style="font-size:0.7rem;">
                {{ healthLabel(link.health_status) }}
              </span>
            </div>

            <!-- Short URL row -->
            <div class="d-flex align-items-center gap-2 mb-1 flex-wrap">
              <span class="link-short-url fw-medium">{{ link.short_url }}</span>
              <button
                class="btn btn-outline-secondary btn-xs copy-btn"
                :class="{ 'btn-success text-white border-success': copied }"
                @click="copyShortUrl"
                title="Copy short URL"
              >
                <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
                  <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0z"/>
                </svg>
                {{ copied ? 'Copied!' : 'Copy' }}
              </button>
              <a :href="link.short_url" target="_blank" rel="noopener noreferrer" class="btn btn-outline-secondary btn-xs" title="Open short URL">
                <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16">
                  <path fill-rule="evenodd" d="M8.636 3.5a.5.5 0 0 0-.5-.5H1.5A1.5 1.5 0 0 0 0 4.5v10A1.5 1.5 0 0 0 1.5 16h10a1.5 1.5 0 0 0 1.5-1.5V7.864a.5.5 0 0 0-1 0V14.5a.5.5 0 0 1-.5.5h-10a.5.5 0 0 1-.5-.5v-10a.5.5 0 0 1 .5-.5h6.636a.5.5 0 0 0 .5-.5z"/>
                  <path fill-rule="evenodd" d="M16 .5a.5.5 0 0 0-.5-.5h-5a.5.5 0 0 0 0 1h3.793L6.146 9.146a.5.5 0 1 0 .708.708L15 1.707V5.5a.5.5 0 0 0 1 0v-5z"/>
                </svg>
                Open
              </a>
            </div>

            <!-- Destination URL -->
            <div class="text-muted link-dest-url">
              <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16" class="me-1 flex-shrink-0">
                <path d="M4.715 6.542 3.343 7.914a3 3 0 1 0 4.243 4.243l1.828-1.829A3 3 0 0 0 8.586 5.5L8 6.086a1 1 0 0 0-.154.199 2 2 0 0 1 .861 3.337L6.88 11.45a2 2 0 1 1-2.83-2.83l.793-.792a4 4 0 0 1-.128-1.287z"/>
                <path d="M6.586 4.672A3 3 0 0 0 7.414 9.5l.775-.776a2 2 0 0 1-.896-3.346L9.12 3.55a2 2 0 1 1 2.83 2.83l-.793.792c.112.42.155.855.128 1.287l1.372-1.372a3 3 0 1 0-4.243-4.243z"/>
              </svg>
              <span class="text-truncate" :title="link.destination_url">{{ link.destination_url }}</span>
            </div>

            <!-- Tags -->
            <div v-if="link.tags && link.tags.length" class="mt-2 d-flex flex-wrap gap-1">
              <span v-for="tag in link.tags" :key="tag" class="badge rounded-pill text-bg-light border tag-badge">{{ tag }}</span>
            </div>
          </div>

          <!-- Right: quick stats -->
          <div class="d-flex gap-3 flex-shrink-0 align-items-start pt-1">
            <div class="text-center">
              <div class="fw-bold" style="font-size:1.25rem;color:#1a1f36;">{{ link.click_count.toLocaleString() }}</div>
              <div class="text-muted" style="font-size:0.7rem;">All-time clicks</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Fallback: still show minimal heading if link load failed -->
      <div v-else class="card-body py-3">
        <h5 class="mb-0 fw-bold">Link Analytics</h5>
        <p class="text-muted small mb-0">{{ route.params.id }}</p>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="alert alert-danger d-flex align-items-center gap-2" role="alert">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5m.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
      </svg>
      {{ error }}
      <button class="btn btn-sm btn-outline-danger ms-auto" @click="loadAnalytics">Retry</button>
    </div>

    <!-- Loading state -->
    <div v-if="loading && !error" class="text-center py-5">
      <div class="spinner-border text-primary" style="width: 2.5rem; height: 2.5rem;" role="status">
        <span class="visually-hidden">Loading analytics...</span>
      </div>
      <p class="text-muted mt-3 mb-0">Loading analytics data...</p>
    </div>

    <template v-if="!loading && analytics">
      <!-- Date range filter -->
      <div class="card border-0 shadow-sm mb-4">
        <div class="card-body">
          <div class="row g-3 align-items-end">
            <div class="col-sm-6 col-md-3">
              <label class="form-label fw-medium small">From</label>
              <input v-model="filterFrom" type="date" class="form-control form-control-sm" />
            </div>
            <div class="col-sm-6 col-md-3">
              <label class="form-label fw-medium small">To</label>
              <input v-model="filterTo" type="date" class="form-control form-control-sm" />
            </div>
            <div class="col-sm-6 col-md-3">
              <label class="form-label fw-medium small">Granularity</label>
              <select v-model="filterGranularity" class="form-select form-select-sm">
                <option value="hour">Hour</option>
                <option value="day">Day</option>
                <option value="week">Week</option>
                <option value="month">Month</option>
              </select>
            </div>
            <div class="col-sm-6 col-md-3">
              <div class="d-flex gap-2">
                <button
                  class="btn btn-primary btn-sm flex-grow-1 d-flex align-items-center justify-content-center gap-2"
                  :disabled="loading"
                  @click="applyFilters"
                >
                  <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                  Apply
                </button>
                <button
                  class="btn btn-outline-secondary btn-sm d-inline-flex align-items-center gap-1 flex-shrink-0"
                  title="Export all analytics to CSV"
                  @click="exportToCSV"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                    <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
                    <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708z"/>
                  </svg>
                  <span class="d-none d-lg-inline">CSV</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Period Comparison Card -->
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

      <!-- Stat cards -->
      <div class="row g-3 mb-4">
        <!-- Total Clicks -->
        <div class="col-sm-6 col-xl-3">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-body">
              <div class="d-flex align-items-center justify-content-between mb-2">
                <div class="d-flex align-items-center gap-2">
                  <span class="text-muted small fw-medium">Total Clicks</span>
                  <span v-if="isLive" class="live-badge d-inline-flex align-items-center gap-1">
                    <span class="live-dot"></span>
                    <span style="font-size: 0.6rem; font-weight: 600; letter-spacing: 0.04em; color: #16a34a;">LIVE</span>
                  </span>
                </div>
                <div class="stat-icon bg-primary-soft rounded-circle d-flex align-items-center justify-content-center" style="width: 36px; height: 36px;">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="#635bff" viewBox="0 0 16 16">
                    <path d="M0 0h1v15h15v1H0zm14.817 3.113a.5.5 0 0 1 .07.704l-4.5 5.5a.5.5 0 0 1-.74.037L7.06 6.767l-3.656 5.027a.5.5 0 0 1-.808-.588l4-5.5a.5.5 0 0 1 .758-.06l2.609 2.61 4.15-5.073a.5.5 0 0 1 .704-.07"/>
                  </svg>
                </div>
              </div>
              <div class="fw-bold fs-4 mb-0">
                {{ (liveTotal !== null ? liveTotal : analytics.total_clicks).toLocaleString() }}
              </div>
            </div>
          </div>
        </div>

        <!-- Unique Clicks -->
        <div class="col-sm-6 col-xl-3">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-body">
              <div class="d-flex align-items-center justify-content-between mb-2">
                <div class="d-flex align-items-center gap-2">
                  <span class="text-muted small fw-medium">Unique Clicks</span>
                  <span v-if="isLive" class="live-badge d-inline-flex align-items-center gap-1">
                    <span class="live-dot"></span>
                    <span style="font-size: 0.6rem; font-weight: 600; letter-spacing: 0.04em; color: #16a34a;">LIVE</span>
                  </span>
                </div>
                <div class="stat-icon rounded-circle d-flex align-items-center justify-content-center" style="width: 36px; height: 36px; background-color: rgba(20, 184, 166, 0.12);">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="#14b8a6" viewBox="0 0 16 16">
                    <path d="M15 14s1 0 1-1-1-4-5-4-5 3-5 4 1 1 1 1zm-7.978-1L7 12.996c.001-.264.167-1.03.76-1.72C8.312 10.629 9.282 10 11 10c1.717 0 2.687.63 3.24 1.276.593.69.758 1.457.76 1.72l-.008.002-.014.002zM11 7a2 2 0 1 0 0-4 2 2 0 0 0 0 4m3-2a3 3 0 1 1-6 0 3 3 0 0 1 6 0M6.936 9.28a6 6 0 0 0-1.23-.247A7 7 0 0 0 5 9c-4 0-5 3-5 4q0 1 1 1h4.216A2.24 2.24 0 0 1 5 13c0-1.01.377-2.042 1.09-2.904.243-.294.526-.569.846-.816M4.92 10A5.5 5.5 0 0 0 4 13H1c0-.26.164-1.03.76-1.724.545-.636 1.492-1.256 3.16-1.275ZM1.5 5.5a3 3 0 1 1 6 0 3 3 0 0 1-6 0m3-2a2 2 0 1 0 0 4 2 2 0 0 0 0-4"/>
                  </svg>
                </div>
              </div>
              <div class="fw-bold fs-4 mb-0">
                {{ (liveUnique !== null ? liveUnique : analytics.unique_clicks).toLocaleString() }}
              </div>
            </div>
          </div>
        </div>

        <!-- Top Referrer -->
        <div class="col-sm-6 col-xl-3">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-body">
              <div class="d-flex align-items-center justify-content-between mb-2">
                <span class="text-muted small fw-medium">Top Referrer</span>
                <div class="stat-icon rounded-circle d-flex align-items-center justify-content-center" style="width: 36px; height: 36px; background-color: rgba(245, 158, 11, 0.12);">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="#f59e0b" viewBox="0 0 16 16">
                    <path d="M4.715 6.542 3.343 7.914a3 3 0 1 0 4.243 4.243l1.828-1.829A3 3 0 0 0 8.586 5.5L8 6.086a1 1 0 0 0-.154.199 2 2 0 0 1 .861 3.337L6.88 11.45a2 2 0 1 1-2.83-2.83l.793-.792a4 4 0 0 1-.128-1.287z"/>
                    <path d="M6.586 4.672A3 3 0 0 0 7.414 9.5l.775-.776a2 2 0 0 1-.896-3.346L9.12 3.55a2 2 0 1 1 2.83 2.83l-.793.792c.112.42.155.855.128 1.287l1.372-1.372a3 3 0 1 0-4.243-4.243z"/>
                  </svg>
                </div>
              </div>
              <div class="fw-bold fs-6 mb-0 text-truncate" :title="topReferrer">
                {{ topReferrer }}
              </div>
            </div>
          </div>
        </div>

        <!-- Top Device -->
        <div class="col-sm-6 col-xl-3">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-body">
              <div class="d-flex align-items-center justify-content-between mb-2">
                <span class="text-muted small fw-medium">Top Device</span>
                <div class="stat-icon rounded-circle d-flex align-items-center justify-content-center" style="width: 36px; height: 36px; background-color: rgba(239, 68, 68, 0.12);">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="#ef4444" viewBox="0 0 16 16">
                    <path d="M11 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM5 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
                    <path d="M8 14a1 1 0 1 0 0-2 1 1 0 0 0 0 2"/>
                  </svg>
                </div>
              </div>
              <div class="fw-bold fs-6 mb-0 text-capitalize">
                {{ topDevice }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Clicks over time chart -->
      <div class="card border-0 shadow-sm mb-4">
        <div class="card-header bg-white border-bottom py-3 px-4">
          <h6 class="mb-0 fw-semibold">Clicks Over Time</h6>
        </div>
        <div class="card-body p-4">
          <div v-if="analytics.time_series.length === 0" class="text-center py-4 text-muted">
            No time series data available for this period.
          </div>
          <VChart
            v-else
            :option="chartOption"
            style="height: 320px;"
            autoresize
          />
        </div>
      </div>

      <!-- Traffic Channels -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Traffic Channels</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">How visitors are reaching your link — classified by referrer domain</p>
          </div>
        </div>
        <div class="card-body">
          <div v-if="!analytics.referrer_categories?.length" class="text-center py-4 text-muted small">
            No channel data available yet.
          </div>
          <div v-else class="row g-4 align-items-center">
            <!-- Donut chart -->
            <div class="col-lg-5">
              <VChart :option="trafficChannelsChartOption" style="height: 280px;" autoresize />
            </div>
            <!-- Channel breakdown list -->
            <div class="col-lg-7">
              <div class="d-flex flex-column gap-3">
                <div
                  v-for="ch in analytics.referrer_categories"
                  :key="ch.category"
                  class="d-flex align-items-center gap-3"
                >
                  <span
                    class="d-inline-flex align-items-center justify-content-center rounded-circle flex-shrink-0"
                    :style="{ width: '32px', height: '32px', backgroundColor: channelColor(ch.category) + '22' }"
                  >
                    <span style="font-size: 0.95rem;">{{ channelIcon(ch.category) }}</span>
                  </span>
                  <div class="flex-grow-1">
                    <div class="d-flex justify-content-between mb-1">
                      <span class="small fw-semibold">{{ channelLabel(ch.category) }}</span>
                      <span class="small text-muted">{{ ch.count.toLocaleString() }} ({{ channelPercent(ch.count) }}%)</span>
                    </div>
                    <div class="progress" style="height: 6px;">
                      <div
                        class="progress-bar"
                        role="progressbar"
                        :style="{ width: channelPercent(ch.count) + '%', backgroundColor: channelColor(ch.category) }"
                        :aria-valuenow="channelPercent(ch.count)"
                        aria-valuemin="0"
                        aria-valuemax="100"
                      ></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Bottom tables row -->
      <div class="row g-4">
        <!-- Referrers table -->
        <div class="col-lg-6">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
              <h6 class="mb-0 fw-semibold">Top Referrers</h6>
              <span class="badge bg-light text-muted border fw-normal">Top 10</span>
            </div>
            <div class="card-body p-0">
              <div v-if="topReferrers.length === 0" class="text-center py-4 text-muted small">
                No referrer data available.
              </div>
              <div v-else class="table-responsive">
                <table class="table table-sm align-middle mb-0">
                  <thead class="table-light">
                    <tr>
                      <th class="ps-4 py-2 fw-semibold text-muted" style="font-size: 0.75rem;">Referrer</th>
                      <th class="py-2 fw-semibold text-muted text-end" style="font-size: 0.75rem; width: 60px;">Clicks</th>
                      <th class="pe-4 py-2 fw-semibold text-muted" style="font-size: 0.75rem; width: 120px;">Share</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="ref in topReferrers" :key="ref.referrer">
                      <td class="ps-4 py-2" style="max-width: 180px;">
                        <span class="text-truncate d-block small" :title="ref.referrer">
                          {{ ref.referrer || 'Direct / None' }}
                        </span>
                      </td>
                      <td class="py-2 text-end small fw-semibold">{{ ref.count.toLocaleString() }}</td>
                      <td class="pe-4 py-2">
                        <div class="d-flex align-items-center gap-2">
                          <div class="progress flex-grow-1" style="height: 6px;">
                            <div
                              class="progress-bar"
                              role="progressbar"
                              :style="{ width: referrerPercent(ref.count) + '%', backgroundColor: '#635bff' }"
                              :aria-valuenow="referrerPercent(ref.count)"
                              aria-valuemin="0"
                              aria-valuemax="100"
                            ></div>
                          </div>
                          <span class="text-muted" style="font-size: 0.7rem; min-width: 32px;">{{ referrerPercent(ref.count) }}%</span>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <!-- Devices table -->
        <div class="col-lg-6">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
              <h6 class="mb-0 fw-semibold">Device Breakdown</h6>
              <span class="badge bg-light text-muted border fw-normal">All Devices</span>
            </div>
            <div class="card-body p-0">
              <div v-if="analytics.devices.length === 0" class="text-center py-4 text-muted small">
                No device data available.
              </div>
              <div v-else class="table-responsive">
                <table class="table table-sm align-middle mb-0">
                  <thead class="table-light">
                    <tr>
                      <th class="ps-4 py-2 fw-semibold text-muted" style="font-size: 0.75rem;">Device</th>
                      <th class="py-2 fw-semibold text-muted text-end" style="font-size: 0.75rem; width: 60px;">Clicks</th>
                      <th class="pe-4 py-2 fw-semibold text-muted" style="font-size: 0.75rem; width: 120px;">Share</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="device in analytics.devices" :key="device.device_type">
                      <td class="ps-4 py-2">
                        <div class="d-flex align-items-center gap-2">
                          <span class="device-icon">
                            <svg v-if="device.device_type.toLowerCase() === 'mobile'" xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="#635bff" viewBox="0 0 16 16">
                              <path d="M11 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM5 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
                              <path d="M8 14a1 1 0 1 0 0-2 1 1 0 0 0 0 2"/>
                            </svg>
                            <svg v-else-if="device.device_type.toLowerCase() === 'tablet'" xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="#635bff" viewBox="0 0 16 16">
                              <path d="M12 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM4 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
                              <path d="M8 14a1 1 0 1 0 0-2 1 1 0 0 0 0 2"/>
                            </svg>
                            <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="#635bff" viewBox="0 0 16 16">
                              <path d="M0 4s0-2 2-2h12s2 0 2 2v6s0 2-2 2h-4q0 1 .25 1.5H11a.5.5 0 0 1 0 1H5a.5.5 0 0 1 0-1h.75Q6 13 6 12H2s-2 0-2-2zm1.398 0a.53.53 0 0 0-.398.49v5.62a.53.53 0 0 0 .398.49h13.204a.53.53 0 0 0 .398-.49V4.49a.53.53 0 0 0-.398-.49z"/>
                            </svg>
                          </span>
                          <span class="small text-capitalize">{{ device.device_type || 'Unknown' }}</span>
                        </div>
                      </td>
                      <td class="py-2 text-end small fw-semibold">{{ device.count.toLocaleString() }}</td>
                      <td class="pe-4 py-2">
                        <div class="d-flex align-items-center gap-2">
                          <div class="progress flex-grow-1" style="height: 6px;">
                            <div
                              class="progress-bar"
                              role="progressbar"
                              :style="{ width: devicePercent(device.count) + '%', backgroundColor: '#635bff' }"
                              :aria-valuenow="devicePercent(device.count)"
                              aria-valuemin="0"
                              aria-valuemax="100"
                            ></div>
                          </div>
                          <span class="text-muted" style="font-size: 0.7rem; min-width: 32px;">{{ devicePercent(device.count) }}%</span>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Browser & OS breakdown -->
      <div class="row g-4 mt-0">
        <!-- Browser breakdown -->
        <div class="col-lg-6">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
              <h6 class="mb-0 fw-semibold">Browser Breakdown</h6>
              <span class="badge bg-light text-muted border fw-normal">All Browsers</span>
            </div>
            <div class="card-body">
              <div v-if="!analytics.browsers || analytics.browsers.length === 0" class="text-center py-4 text-muted small">
                No browser data available.
              </div>
              <div v-else class="row align-items-center g-0">
                <div class="col-6">
                  <VChart :option="browserChartOption" style="height: 200px;" autoresize />
                </div>
                <div class="col-6">
                  <table class="table table-sm align-middle mb-0">
                    <tbody>
                      <tr v-for="(b, i) in analytics.browsers" :key="b.browser">
                        <td class="ps-0 py-1 border-0">
                          <div class="d-flex align-items-center gap-2">
                            <span
                              class="rounded-circle d-inline-block flex-shrink-0"
                              :style="{ width: '8px', height: '8px', backgroundColor: chartColors[i % chartColors.length] }"
                            ></span>
                            <span class="small text-truncate" style="max-width: 110px;" :title="b.browser">{{ b.browser }}</span>
                          </div>
                        </td>
                        <td class="py-1 text-end small fw-semibold pe-0 border-0">{{ browserPercent(b.count) }}%</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- OS breakdown -->
        <div class="col-lg-6">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
              <h6 class="mb-0 fw-semibold">Operating System</h6>
              <span class="badge bg-light text-muted border fw-normal">All OS</span>
            </div>
            <div class="card-body">
              <div v-if="!analytics.os_breakdown || analytics.os_breakdown.length === 0" class="text-center py-4 text-muted small">
                No OS data available.
              </div>
              <div v-else class="row align-items-center g-0">
                <div class="col-6">
                  <VChart :option="osChartOption" style="height: 200px;" autoresize />
                </div>
                <div class="col-6">
                  <table class="table table-sm align-middle mb-0">
                    <tbody>
                      <tr v-for="(o, i) in analytics.os_breakdown" :key="o.os">
                        <td class="ps-0 py-1 border-0">
                          <div class="d-flex align-items-center gap-2">
                            <span
                              class="rounded-circle d-inline-block flex-shrink-0"
                              :style="{ width: '8px', height: '8px', backgroundColor: chartColors[i % chartColors.length] }"
                            ></span>
                            <span class="small text-truncate" style="max-width: 110px;" :title="o.os">{{ o.os }}</span>
                          </div>
                        </td>
                        <td class="py-1 text-end small fw-semibold pe-0 border-0">{{ osPercent(o.count) }}%</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Source Breakdown -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Click Source</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">How visitors arrived — direct web link, QR code, or API call</p>
          </div>
          <span class="badge bg-light text-muted border fw-normal">All Sources</span>
        </div>
        <div class="card-body">
          <div v-if="!analytics.sources || analytics.sources.length === 0" class="text-center py-4 text-muted small">
            No source data available.
          </div>
          <div v-else class="row align-items-center g-0">
            <div class="col-lg-6">
              <VChart :option="sourceChartOption" style="height: 200px;" autoresize />
            </div>
            <div class="col-lg-6">
              <table class="table table-sm align-middle mb-0">
                <tbody>
                  <tr v-for="(s, i) in analytics.sources" :key="s.source">
                    <td class="ps-0 py-2 border-0">
                      <div class="d-flex align-items-center gap-2">
                        <span
                          class="rounded-circle d-inline-block flex-shrink-0"
                          :style="{ width: '8px', height: '8px', backgroundColor: chartColors[i % chartColors.length] }"
                        ></span>
                        <span class="small">{{ sourceIcon(s.source) }} {{ sourceLabel(s.source) }}</span>
                      </div>
                    </td>
                    <td class="py-2 text-end small fw-semibold border-0">{{ s.count.toLocaleString() }}</td>
                    <td class="py-2 text-end border-0" style="width: 80px;">
                      <span class="text-muted" style="font-size: 0.7rem;">{{ sourcePercent(s.count) }}%</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Visitor Loyalty -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Visitor Loyalty</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">First-time vs returning visitors, based on hashed IP tracking</p>
          </div>
        </div>
        <div class="card-body">
          <div
            v-if="!analytics.first_time_visitors && !analytics.returning_visitors"
            class="text-center py-4 text-muted small"
          >
            No visitor loyalty data yet. Data is collected on new clicks.
          </div>
          <div v-else class="row align-items-center g-0">
            <!-- Donut chart -->
            <div class="col-lg-4">
              <VChart :option="visitorLoyaltyChartOption" style="height: 200px;" autoresize />
            </div>
            <!-- Stat breakdown -->
            <div class="col-lg-8">
              <div class="row g-3">
                <!-- First-time -->
                <div class="col-sm-4">
                  <div class="p-3 rounded-3" style="background: rgba(99, 91, 255, 0.06);">
                    <div class="small text-muted mb-1 fw-medium">First-Time</div>
                    <div class="fw-bold fs-5 mb-0" style="color: #635bff;">
                      {{ analytics.first_time_visitors.toLocaleString() }}
                    </div>
                    <div class="text-muted" style="font-size: 0.72rem;">
                      {{ firstTimePercent }}% of tracked visitors
                    </div>
                  </div>
                </div>
                <!-- Returning -->
                <div class="col-sm-4">
                  <div class="p-3 rounded-3" style="background: rgba(20, 184, 166, 0.06);">
                    <div class="small text-muted mb-1 fw-medium">Returning</div>
                    <div class="fw-bold fs-5 mb-0" style="color: #14b8a6;">
                      {{ analytics.returning_visitors.toLocaleString() }}
                    </div>
                    <div class="text-muted" style="font-size: 0.72rem;">
                      {{ returningPercent }}% of tracked visitors
                    </div>
                  </div>
                </div>
                <!-- Return rate -->
                <div class="col-sm-4">
                  <div class="p-3 rounded-3" style="background: rgba(245, 158, 11, 0.06);">
                    <div class="small text-muted mb-1 fw-medium">Return Rate</div>
                    <div class="fw-bold fs-5 mb-0" style="color: #f59e0b;">
                      {{ returningPercent }}%
                    </div>
                    <div class="text-muted" style="font-size: 0.72rem;">
                      visitors who came back
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Feature #58: Link Key Dates Timeline Card -->
      <div v-if="link" class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Link Timeline</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">Key dates and milestones for this link</p>
          </div>
        </div>
        <div class="card-body py-3 px-4">
          <div class="timeline-list">
            <!-- Created — always shown -->
            <div class="timeline-item">
              <div class="timeline-icon timeline-icon--created">🎉</div>
              <div class="timeline-body">
                <div class="timeline-label">Created</div>
                <div class="timeline-date">{{ formatLinkDate(link.created_at) }}</div>
              </div>
            </div>

            <!-- Last Updated — only if different from created_at -->
            <div v-if="link.updated_at && link.updated_at !== link.created_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--updated">✏️</div>
              <div class="timeline-body">
                <div class="timeline-label">Last Updated</div>
                <div class="timeline-date">{{ formatLinkDate(link.updated_at) }}</div>
              </div>
            </div>

            <!-- Health Last Checked -->
            <div v-if="link.health_checked_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--health">🔍</div>
              <div class="timeline-body">
                <div class="d-flex align-items-center gap-2 flex-wrap">
                  <div class="timeline-label">Health Last Checked</div>
                  <span class="badge rounded-pill px-2 py-1 badge-health" :class="healthClass(link.health_status)" style="font-size:0.68rem;">
                    {{ healthLabel(link.health_status) }}
                  </span>
                </div>
                <div class="timeline-date">{{ formatLinkDate(link.health_checked_at) }}</div>
              </div>
            </div>

            <!-- Expires -->
            <div v-if="link.expires_at" class="timeline-item">
              <div class="timeline-icon timeline-icon--expires">⏰</div>
              <div class="timeline-body">
                <div class="timeline-label">Expires</div>
                <div class="timeline-date" :class="expiryDateClass(link.expires_at)">
                  {{ formatLinkDate(link.expires_at) }}
                  <span class="ms-1 small fw-semibold" :class="expiryDateClass(link.expires_at)">
                    {{ expiryCountdown(link.expires_at) }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Click Limit -->
            <div v-if="link.max_clicks" class="timeline-item">
              <div class="timeline-icon timeline-icon--limit">🎯</div>
              <div class="timeline-body">
                <div class="timeline-label">Click Limit</div>
                <div class="timeline-date">
                  Max clicks: <strong>{{ link.max_clicks.toLocaleString() }}</strong>
                  <span class="text-muted ms-1">({{ link.click_count.toLocaleString() }} used, {{ Math.max(0, link.max_clicks - link.click_count).toLocaleString() }} remaining)</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Feature #62: Security & Limits Card -->
      <div v-if="link" class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Security &amp; Limits</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">Access controls and usage limits for this link</p>
          </div>
        </div>
        <div class="card-body">
          <div class="security-grid">
            <!-- Password Protected -->
            <div class="security-item">
              <div class="security-item-label text-muted small">Password Protected</div>
              <div class="security-item-value">
                <span v-if="link.has_password" class="badge rounded-pill text-bg-warning" style="font-size:0.75rem;">🔒 Yes</span>
                <span v-else class="badge rounded-pill bg-light text-muted border" style="font-size:0.75rem;">No</span>
              </div>
            </div>

            <!-- Active Status -->
            <div class="security-item">
              <div class="security-item-label text-muted small">Active Status</div>
              <div class="security-item-value">
                <span class="badge rounded-pill px-2 py-1" :class="link.is_active ? 'text-bg-success' : 'bg-secondary text-white'" style="font-size:0.75rem;">
                  {{ link.is_active ? 'Active' : 'Inactive' }}
                </span>
              </div>
            </div>

            <!-- Redirect Type -->
            <div class="security-item">
              <div class="security-item-label text-muted small">Redirect Type</div>
              <div class="security-item-value">
                <span class="badge rounded-pill bg-light text-muted border" style="font-size:0.75rem;">
                  {{ link.redirect_type === 301 ? '301 Permanent' : '302 Temporary' }}
                </span>
              </div>
            </div>

            <!-- Expiry -->
            <div class="security-item">
              <div class="security-item-label text-muted small">Expiry</div>
              <div class="security-item-value">
                <span v-if="link.expires_at" :class="expiryDateClass(link.expires_at)" style="font-size:0.82rem;font-weight:500;">
                  {{ formatLinkDate(link.expires_at) }}
                  <span class="ms-1 small">{{ expiryCountdown(link.expires_at) }}</span>
                </span>
                <span v-else class="text-muted small">Never expires</span>
              </div>
            </div>

            <!-- Click Limit with progress bar — full width row -->
            <div class="security-item security-item--full">
              <div class="security-item-label text-muted small mb-2">Click Limit</div>
              <div v-if="link.max_clicks">
                <div class="d-flex justify-content-between align-items-baseline mb-1">
                  <span class="small fw-semibold">
                    {{ link.click_count.toLocaleString() }} / {{ link.max_clicks.toLocaleString() }} clicks used
                  </span>
                  <span class="small text-muted">{{ clickLimitPercent(link.click_count, link.max_clicks) }}%</span>
                </div>
                <div class="progress" style="height: 8px;">
                  <div
                    class="progress-bar"
                    role="progressbar"
                    :style="{ width: clickLimitPercent(link.click_count, link.max_clicks) + '%', backgroundColor: clickLimitBarColor(link.click_count, link.max_clicks) }"
                    :aria-valuenow="clickLimitPercent(link.click_count, link.max_clicks)"
                    aria-valuemin="0"
                    aria-valuemax="100"
                  ></div>
                </div>
              </div>
              <div v-else class="text-muted small">Unlimited</div>
            </div>

            <!-- Tags — full width row -->
            <div class="security-item security-item--full">
              <div class="security-item-label text-muted small mb-1">Tags</div>
              <div v-if="link.tags && link.tags.length" class="d-flex flex-wrap gap-1">
                <span v-for="tag in link.tags" :key="tag" class="badge rounded-pill text-bg-light border tag-badge">{{ tag }}</span>
              </div>
              <span v-else class="text-muted small">None</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Click Heatmap -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Click Heatmap</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">Clicks by hour of day × day of week (UTC)</p>
          </div>
          <span class="badge bg-light text-muted border fw-normal">24 × 7</span>
        </div>
        <div class="card-body pb-2">
          <div v-if="!analytics.heatmap || analytics.heatmap.length === 0" class="text-center py-4 text-muted small">
            No heatmap data available for this period.
          </div>
          <VChart v-else :option="heatmapChartOption" style="height: 260px;" autoresize />
        </div>
      </div>

      <!-- Geographic World Map -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">Geographic Distribution</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">Click density by country — scroll to zoom, drag to pan</p>
          </div>
          <span class="badge bg-light text-muted border fw-normal">World Map</span>
        </div>
        <div class="card-body p-3">
          <div v-if="worldMapLoading" class="text-center py-5">
            <div class="spinner-border spinner-border-sm text-muted" role="status">
              <span class="visually-hidden">Loading map...</span>
            </div>
            <p class="text-muted small mt-2 mb-0">Loading world map…</p>
          </div>
          <div
            v-else-if="!mapLoaded || !analytics?.countries?.length"
            class="text-center py-5 text-muted small"
          >
            No geographic data available yet. Geographic data is collected on new clicks.
          </div>
          <VChart
            v-else
            :option="geoMapOption"
            style="height: 420px;"
            autoresize
          />
        </div>
      </div>

      <!-- Top Countries -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <h6 class="mb-0 fw-semibold">Top Countries</h6>
          <span class="badge bg-light text-muted border fw-normal">Top 15</span>
        </div>
        <div class="card-body">
          <div v-if="!analytics.countries || analytics.countries.length === 0" class="text-center py-4 text-muted small">
            No country data available yet. Geographic data is collected on new clicks.
          </div>
          <div v-else class="row g-4">
            <!-- Horizontal bar chart -->
            <div class="col-lg-7">
              <VChart :option="countriesChartOption" style="height: 360px;" autoresize />
            </div>
            <!-- Country table -->
            <div class="col-lg-5">
              <div class="table-responsive">
                <table class="table table-sm align-middle mb-0">
                  <thead class="table-light">
                    <tr>
                      <th class="ps-3 py-2 fw-semibold text-muted" style="font-size: 0.75rem;">Country</th>
                      <th class="py-2 fw-semibold text-muted text-end" style="font-size: 0.75rem; width: 60px;">Clicks</th>
                      <th class="pe-3 py-2 fw-semibold text-muted" style="font-size: 0.75rem; width: 120px;">Share</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="c in analytics.countries" :key="c.country">
                      <td class="ps-3 py-2">
                        <div class="d-flex align-items-center gap-2">
                          <span style="font-size: 1.1rem; line-height: 1;">{{ countryFlag(c.country) }}</span>
                          <span class="small">{{ countryDisplayName(c.country) }}</span>
                        </div>
                      </td>
                      <td class="py-2 text-end small fw-semibold">{{ c.count.toLocaleString() }}</td>
                      <td class="pe-3 py-2">
                        <div class="d-flex align-items-center gap-2">
                          <div class="progress flex-grow-1" style="height: 6px;">
                            <div
                              class="progress-bar"
                              role="progressbar"
                              :style="{ width: countryPercent(c.count) + '%', backgroundColor: '#635bff' }"
                              :aria-valuenow="countryPercent(c.count)"
                              aria-valuemin="0"
                              aria-valuemax="100"
                            ></div>
                          </div>
                          <span class="text-muted" style="font-size: 0.7rem; min-width: 32px;">{{ countryPercent(c.count) }}%</span>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- Top Cities -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <h6 class="mb-0 fw-semibold">Top Cities</h6>
          <span class="badge bg-light text-muted border fw-normal">Top 20</span>
        </div>
        <div class="card-body">
          <div v-if="!analytics.cities || analytics.cities.length === 0" class="text-center py-4 text-muted small">
            No city data available yet. City data is collected on new clicks.
          </div>
          <div v-else class="table-responsive">
            <table class="table table-sm align-middle mb-0">
              <thead class="table-light">
                <tr>
                  <th class="ps-3 py-2 fw-semibold text-muted" style="font-size: 0.75rem;">#</th>
                  <th class="py-2 fw-semibold text-muted" style="font-size: 0.75rem;">City</th>
                  <th class="py-2 fw-semibold text-muted" style="font-size: 0.75rem;">Country</th>
                  <th class="py-2 fw-semibold text-muted text-end" style="font-size: 0.75rem; width: 70px;">Clicks</th>
                  <th class="pe-3 py-2 fw-semibold text-muted" style="font-size: 0.75rem; width: 140px;">Share</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(c, idx) in analytics.cities" :key="c.city + c.country">
                  <td class="ps-3 py-2 text-muted small">{{ idx + 1 }}</td>
                  <td class="py-2">
                    <span class="small fw-medium">{{ c.city }}</span>
                  </td>
                  <td class="py-2">
                    <div class="d-flex align-items-center gap-2">
                      <span style="font-size: 1rem; line-height: 1;">{{ countryFlag(c.country) }}</span>
                      <span class="small text-muted">{{ countryDisplayName(c.country) }}</span>
                    </div>
                  </td>
                  <td class="py-2 text-end small fw-semibold">{{ c.count.toLocaleString() }}</td>
                  <td class="pe-3 py-2">
                    <div class="d-flex align-items-center gap-2">
                      <div class="progress flex-grow-1" style="height: 6px;">
                        <div
                          class="progress-bar"
                          role="progressbar"
                          :style="{ width: cityPercent(c.count) + '%', backgroundColor: '#635bff' }"
                          :aria-valuenow="cityPercent(c.count)"
                          aria-valuemin="0"
                          aria-valuemax="100"
                        ></div>
                      </div>
                      <span class="text-muted" style="font-size: 0.7rem; min-width: 32px;">{{ cityPercent(c.count) }}%</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- UTM Campaign Tracking -->
      <div class="card border-0 shadow-sm mt-4">
        <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
          <div>
            <h6 class="mb-0 fw-semibold">UTM Campaign Tracking</h6>
            <p class="text-muted mb-0" style="font-size: 0.72rem;">Clicks from links with UTM parameters appended (e.g., <code>?utm_source=email</code>)</p>
          </div>
        </div>
        <div class="card-body">
          <div
            v-if="!analytics.utm_sources?.length && !analytics.utm_mediums?.length && !analytics.utm_campaigns?.length && !analytics.utm_contents?.length && !analytics.utm_terms?.length"
            class="text-center py-4 text-muted small"
          >
            No UTM data yet. Append <code>?utm_source=email&amp;utm_medium=newsletter&amp;utm_campaign=launch</code> to your short link when sharing.
          </div>
          <div v-else>
            <!-- Row 1: Source · Medium · Campaign -->
            <div class="row g-4 mb-4">
              <!-- UTM Sources -->
              <div class="col-lg-4">
                <h6 class="small fw-semibold text-muted text-uppercase mb-2" style="letter-spacing: 0.05em;">Source</h6>
                <div v-if="!analytics.utm_sources?.length" class="text-muted small py-1">No data</div>
                <div v-else class="d-flex flex-column gap-1">
                  <div v-for="u in analytics.utm_sources" :key="u.value" class="d-flex align-items-center gap-2">
                    <div class="flex-grow-1">
                      <div class="d-flex justify-content-between small mb-1">
                        <span class="text-truncate fw-medium" style="max-width: 140px;" :title="u.value">{{ u.value }}</span>
                        <span class="text-muted">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="progress" style="height: 5px;">
                        <div
                          class="progress-bar"
                          :style="{ width: utmPercent(u.count, analytics.utm_sources) + '%', backgroundColor: '#635bff' }"
                        ></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- UTM Mediums -->
              <div class="col-lg-4">
                <h6 class="small fw-semibold text-muted text-uppercase mb-2" style="letter-spacing: 0.05em;">Medium</h6>
                <div v-if="!analytics.utm_mediums?.length" class="text-muted small py-1">No data</div>
                <div v-else class="d-flex flex-column gap-1">
                  <div v-for="u in analytics.utm_mediums" :key="u.value" class="d-flex align-items-center gap-2">
                    <div class="flex-grow-1">
                      <div class="d-flex justify-content-between small mb-1">
                        <span class="text-truncate fw-medium" style="max-width: 140px;" :title="u.value">{{ u.value }}</span>
                        <span class="text-muted">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="progress" style="height: 5px;">
                        <div
                          class="progress-bar"
                          :style="{ width: utmPercent(u.count, analytics.utm_mediums) + '%', backgroundColor: '#14b8a6' }"
                        ></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- UTM Campaigns -->
              <div class="col-lg-4">
                <h6 class="small fw-semibold text-muted text-uppercase mb-2" style="letter-spacing: 0.05em;">Campaign</h6>
                <div v-if="!analytics.utm_campaigns?.length" class="text-muted small py-1">No data</div>
                <div v-else class="d-flex flex-column gap-1">
                  <div v-for="u in analytics.utm_campaigns" :key="u.value" class="d-flex align-items-center gap-2">
                    <div class="flex-grow-1">
                      <div class="d-flex justify-content-between small mb-1">
                        <span class="text-truncate fw-medium" style="max-width: 140px;" :title="u.value">{{ u.value }}</span>
                        <span class="text-muted">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="progress" style="height: 5px;">
                        <div
                          class="progress-bar"
                          :style="{ width: utmPercent(u.count, analytics.utm_campaigns) + '%', backgroundColor: '#f59e0b' }"
                        ></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Row 2: Content · Term -->
            <div v-if="analytics.utm_contents?.length || analytics.utm_terms?.length" class="row g-4 pt-3 border-top">
              <!-- UTM Contents -->
              <div class="col-lg-6">
                <h6 class="small fw-semibold text-muted text-uppercase mb-2" style="letter-spacing: 0.05em;">Content</h6>
                <div v-if="!analytics.utm_contents?.length" class="text-muted small py-1">No data</div>
                <div v-else class="d-flex flex-column gap-1">
                  <div v-for="u in analytics.utm_contents" :key="u.value" class="d-flex align-items-center gap-2">
                    <div class="flex-grow-1">
                      <div class="d-flex justify-content-between small mb-1">
                        <span class="text-truncate fw-medium" style="max-width: 200px;" :title="u.value">{{ u.value }}</span>
                        <span class="text-muted">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="progress" style="height: 5px;">
                        <div
                          class="progress-bar"
                          :style="{ width: utmPercent(u.count, analytics.utm_contents) + '%', backgroundColor: '#ec4899' }"
                        ></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- UTM Terms -->
              <div class="col-lg-6">
                <h6 class="small fw-semibold text-muted text-uppercase mb-2" style="letter-spacing: 0.05em;">Term</h6>
                <div v-if="!analytics.utm_terms?.length" class="text-muted small py-1">No data</div>
                <div v-else class="d-flex flex-column gap-1">
                  <div v-for="u in analytics.utm_terms" :key="u.value" class="d-flex align-items-center gap-2">
                    <div class="flex-grow-1">
                      <div class="d-flex justify-content-between small mb-1">
                        <span class="text-truncate fw-medium" style="max-width: 200px;" :title="u.value">{{ u.value }}</span>
                        <span class="text-muted">{{ u.count.toLocaleString() }}</span>
                      </div>
                      <div class="progress" style="height: 5px;">
                        <div
                          class="progress-bar"
                          :style="{ width: utmPercent(u.count, analytics.utm_terms) + '%', backgroundColor: '#6366f1' }"
                        ></div>
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

<style scoped>
/* ── Link info card ───────────────────────────────────────────────────────── */
.link-info-card {
  border-left: 3px solid #635bff !important;
}

.link-info-title {
  font-size: 1.05rem;
  color: #1a1f36;
}

.link-short-url {
  font-size: 0.875rem;
  color: #635bff;
  font-family: monospace;
}

.link-dest-url {
  font-size: 0.78rem;
  display: flex;
  align-items: center;
  max-width: 520px;
  overflow: hidden;
  white-space: nowrap;
}

.btn-xs {
  padding: 0.15rem 0.5rem;
  font-size: 0.75rem;
  line-height: 1.4;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.copy-btn {
  transition: background-color 0.15s, border-color 0.15s, color 0.15s;
}

.tag-badge {
  font-size: 0.7rem;
  color: #697386 !important;
}

/* Health status badges */
.badge-health--ok   { background: rgba(34,197,94,0.12);  color: #16a34a; }
.badge-health--bad  { background: rgba(239,68,68,0.12);  color: #dc2626; }
.badge-health--warn { background: rgba(245,158,11,0.12); color: #d97706; }
.badge-health--gray { background: rgba(107,114,128,0.1); color: #6b7280; }

/* Loading skeleton */
.skeleton {
  background: linear-gradient(90deg, #e3e8ee 25%, #f0f2f5 50%, #e3e8ee 75%);
  background-size: 200% 100%;
  animation: shimmer 1.2s infinite;
  display: inline-block;
}
@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

:global(.dark-mode) .link-info-title { color: #e6edf3; }
:global(.dark-mode) .link-short-url  { color: #a89fff; }
:global(.dark-mode) .tag-badge       { color: #8b949e !important; background: #21262d !important; border-color: #30363d !important; }
:global(.dark-mode) .skeleton        { background: linear-gradient(90deg, #21262d 25%, #30363d 50%, #21262d 75%); background-size: 200% 100%; }

.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}

.btn-primary:hover:not(:disabled) {
  background-color: #5249e0;
  border-color: #5249e0;
}

.text-primary {
  color: #635bff !important;
}

.stat-icon {
  flex-shrink: 0;
}

.bg-primary-soft {
  background-color: rgba(99, 91, 255, 0.12);
}

.progress {
  border-radius: 100px;
}

.progress-bar {
  border-radius: 100px;
}

/* Live indicator */
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

/* ── Period Comparison Card ───────────────────────────────────────────────── */
.comparison-metric {
  padding: 0.5rem 0;
}

.comp-label {
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.comp-current {
  color: #1a1f36;
}

.comp-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.8rem;
  font-weight: 600;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
}

.comp-badge--up {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
}

.comp-badge--down {
  background: rgba(220, 38, 38, 0.12);
  color: #dc2626;
}

.comp-badge--stable {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.comp-prev {
  margin-top: 0.2rem;
}

.comparison-period-box {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  height: 100%;
}

.current-period {
  background: rgba(99, 91, 255, 0.06);
  border: 1px solid rgba(99, 91, 255, 0.2);
}

.previous-period {
  background: rgba(107, 114, 128, 0.06);
  border: 1px solid rgba(107, 114, 128, 0.15);
}

.comp-period-label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: #697386;
  margin-bottom: 0.25rem;
}

.comp-period-dates {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #1a1f36;
}

/* ── Feature #58: Link Timeline ─────────────────────────────────────────── */
.timeline-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.timeline-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.65rem 0;
  border-bottom: 1px solid #f1f3f5;
}

.timeline-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.timeline-item:first-child {
  padding-top: 0;
}

.timeline-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  font-size: 0.9rem;
  flex-shrink: 0;
  margin-top: 1px;
}

.timeline-icon--created { background: rgba(99, 91, 255, 0.1); }
.timeline-icon--updated { background: rgba(20, 184, 166, 0.1); }
.timeline-icon--health  { background: rgba(245, 158, 11, 0.1); }
.timeline-icon--expires { background: rgba(239, 68, 68, 0.1); }
.timeline-icon--limit   { background: rgba(99, 91, 255, 0.08); }

.timeline-body {
  flex: 1;
  min-width: 0;
}

.timeline-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #697386;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: 0.15rem;
}

.timeline-date {
  font-size: 0.85rem;
  color: #1a1f36;
  font-weight: 500;
}

:global(.dark-mode) .timeline-item { border-bottom-color: #21262d; }
:global(.dark-mode) .timeline-label { color: #8b949e; }
:global(.dark-mode) .timeline-date  { color: #e6edf3; }

/* ── Feature #62: Security & Limits ─────────────────────────────────────── */
.security-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem 1.5rem;
}

.security-item {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.security-item--full {
  grid-column: 1 / -1;
}

.security-item-label {
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.security-item-value {
  font-size: 0.85rem;
  font-weight: 500;
  color: #1a1f36;
}

:global(.dark-mode) .security-item-value { color: #e6edf3; }
</style>
