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

    <!-- Page title -->
    <div class="mb-4">
      <h4 class="fw-bold mb-1">Link Analytics</h4>
      <p class="text-muted small mb-0">{{ route.params.id }}</p>
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
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { RouterLink } from 'vue-router';
import { use } from 'echarts/core';
import { LineChart, PieChart, BarChart, HeatmapChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent, VisualMapComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import VChart from 'vue-echarts';
import linksApi from '@/api/links';
import type { AnalyticsResponse } from '@/types/links';

use([LineChart, PieChart, BarChart, HeatmapChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent, VisualMapComponent, CanvasRenderer]);

const route = useRoute();
const router = useRouter();

const analytics = ref<AnalyticsResponse | null>(null);
const loading = ref(false);
const error = ref('');

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

function applyFilters() {
  loadAnalytics();
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
  rows.push('Total Clicks,Unique Clicks');
  rows.push(`${a.total_clicks},${a.unique_clicks}`);
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

  // ── Click heatmap ───────────────────────────────────────────────────────────
  rows.push('CLICK HEATMAP (UTC)');
  rows.push('Day,Hour,Clicks');
  a.heatmap.forEach((h) => {
    const hour = `${String(h.hour).padStart(2, '0')}:00`;
    rows.push(`${dayNames[h.day_of_week]},${hour},${h.count}`);
  });

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
  loadAnalytics();
  connectLiveCount();
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
</style>
