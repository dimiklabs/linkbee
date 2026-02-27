<template>
  <div class="container-fluid py-4">
    <!-- Page header -->
    <div class="d-flex align-items-center justify-content-between mb-4 flex-wrap gap-2">
      <div>
        <h4 class="fw-bold mb-1">Scheduled Reports</h4>
        <p class="text-muted small mb-0">Automatically receive analytics summaries by email on a set cadence.</p>
      </div>
      <button class="btn btn-primary btn-sm d-flex align-items-center gap-2" @click="openCreate">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
        </svg>
        New Report
      </button>
    </div>

    <!-- Error -->
    <div v-if="error" class="alert alert-danger d-flex align-items-center gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5m.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
      </svg>
      {{ error }}
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" style="width:2.5rem;height:2.5rem;" role="status">
        <span class="visually-hidden">Loading…</span>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="!loading && reports.length === 0 && !error" class="text-center py-5 text-muted">
      <div style="font-size:3rem;line-height:1;margin-bottom:.75rem;">📧</div>
      <p class="mb-2">No scheduled reports yet.</p>
      <button class="btn btn-primary btn-sm" @click="openCreate">Create your first report</button>
    </div>

    <!-- Reports list -->
    <div v-else class="row g-3">
      <div v-for="report in reports" :key="report.id" class="col-12">
        <div class="card border-0 shadow-sm report-card">
          <div class="card-body d-flex align-items-center gap-3 flex-wrap">
            <!-- Status indicator -->
            <div class="status-dot" :class="report.is_active ? 'active' : 'inactive'" :title="report.is_active ? 'Active' : 'Paused'"></div>

            <!-- Info -->
            <div class="flex-fill min-w-0">
              <div class="fw-semibold text-truncate">{{ report.name }}</div>
              <div class="text-muted small mt-1">
                <span class="badge-freq me-2">{{ capitalize(report.frequency) }}</span>
                <span>{{ report.link_ids.length }} link{{ report.link_ids.length !== 1 ? 's' : '' }}</span>
                <span v-if="report.next_run_at" class="ms-2">· Next: {{ formatDate(report.next_run_at) }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="d-flex align-items-center gap-2 flex-shrink-0">
              <!-- Toggle active -->
              <button
                class="btn btn-sm btn-outline-secondary"
                :title="report.is_active ? 'Pause report' : 'Resume report'"
                :disabled="actionLoading === report.id"
                @click="toggleActive(report)"
              >
                {{ report.is_active ? 'Pause' : 'Resume' }}
              </button>
              <!-- Send now -->
              <button
                class="btn btn-sm btn-outline-primary"
                title="Send report immediately"
                :disabled="actionLoading === report.id"
                @click="sendNow(report.id)"
              >
                <span v-if="actionLoading === report.id" class="spinner-border spinner-border-sm me-1" role="status"></span>
                Send now
              </button>
              <!-- Preview -->
              <button
                class="btn btn-sm btn-outline-secondary"
                title="Preview report"
                @click="openPreview(report)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8M1.173 8a13 13 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5s3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5s-3.879-1.168-5.168-2.457A13 13 0 0 1 1.172 8z"/>
                  <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5M4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0"/>
                </svg>
              </button>
              <!-- Deliveries -->
              <button
                class="btn btn-sm btn-outline-secondary"
                title="View delivery history"
                @click="openDeliveries(report)"
              >
                History
              </button>
              <!-- Edit -->
              <button
                class="btn btn-sm btn-outline-secondary"
                title="Edit report"
                @click="openEdit(report)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325z"/>
                </svg>
              </button>
              <!-- Delete -->
              <button
                class="btn btn-sm btn-outline-danger"
                title="Delete report"
                @click="deleteReport(report.id)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                  <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Create / Edit Modal ──────────────────────────────────────────────── -->
    <div v-if="showFormModal" class="modal-backdrop-custom" @click.self="closeForm">
      <div class="modal-dialog-custom card shadow">
        <div class="card-header d-flex align-items-center justify-content-between border-0 pt-3 px-4 pb-0">
          <h6 class="mb-0 fw-semibold">{{ editingReport ? 'Edit Report' : 'New Scheduled Report' }}</h6>
          <button class="btn-icon" @click="closeForm">✕</button>
        </div>
        <div class="card-body px-4 pb-4">
          <!-- Name -->
          <div class="mb-3">
            <label class="form-label fw-medium small">Report Name <span class="text-danger">*</span></label>
            <input v-model="form.name" type="text" class="form-control form-control-sm" placeholder="e.g. Weekly Top Links" maxlength="255" />
          </div>

          <!-- Links -->
          <div class="mb-3">
            <label class="form-label fw-medium small">Links <span class="text-danger">*</span> <span class="text-muted">(select 1–10)</span></label>
            <div class="link-selector" ref="selectorRef">
              <div
                class="selector-trigger form-control form-control-sm d-flex align-items-center justify-content-between"
                @click="selectorOpen = !selectorOpen"
              >
                <span v-if="form.link_ids.length === 0" class="text-muted">Choose links…</span>
                <span v-else>{{ form.link_ids.length }} link{{ form.link_ids.length !== 1 ? 's' : '' }} selected</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16" class="flex-shrink-0 ms-1">
                  <path fill-rule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
                </svg>
              </div>
              <div v-if="selectorOpen" class="selector-dropdown card shadow-sm border">
                <div class="px-2 pt-2">
                  <input v-model="linkSearch" class="form-control form-control-sm" placeholder="Search…" @click.stop />
                </div>
                <div class="selector-options">
                  <label
                    v-for="link in filteredAllLinks"
                    :key="link.id"
                    class="selector-option d-flex align-items-center gap-2 px-3 py-2"
                    :class="{ disabled: !form.link_ids.includes(link.id) && form.link_ids.length >= 10 }"
                    @click.stop
                  >
                    <input
                      type="checkbox"
                      :value="link.id"
                      v-model="form.link_ids"
                      :disabled="!form.link_ids.includes(link.id) && form.link_ids.length >= 10"
                      class="form-check-input flex-shrink-0"
                    />
                    <div class="min-w-0">
                      <div class="small fw-medium text-truncate">{{ link.title || link.slug }}</div>
                      <div class="text-muted" style="font-size:.7rem;">{{ link.short_url }}</div>
                    </div>
                  </label>
                  <div v-if="filteredAllLinks.length === 0" class="text-muted small text-center py-3">No links found</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Frequency -->
          <div class="mb-4">
            <label class="form-label fw-medium small">Frequency <span class="text-danger">*</span></label>
            <div class="d-flex gap-2">
              <button
                v-for="freq in frequencies"
                :key="freq.value"
                class="btn btn-sm flex-fill"
                :class="form.frequency === freq.value ? 'btn-primary' : 'btn-outline-secondary'"
                @click="form.frequency = freq.value"
              >
                {{ freq.label }}
              </button>
            </div>
          </div>

          <!-- Error -->
          <div v-if="formError" class="alert alert-danger py-2 px-3 small mb-3">{{ formError }}</div>

          <!-- Submit -->
          <button
            class="btn btn-primary btn-sm w-100"
            :disabled="formLoading || !form.name.trim() || form.link_ids.length === 0"
            @click="submitForm"
          >
            <span v-if="formLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
            {{ editingReport ? 'Save Changes' : 'Create Report' }}
          </button>
        </div>
      </div>
    </div>

    <!-- ── Delivery History Modal ─────────────────────────────────────────── -->
    <div v-if="showDeliveriesModal" class="modal-backdrop-custom" @click.self="showDeliveriesModal = false">
      <div class="modal-dialog-custom card shadow">
        <div class="card-header d-flex align-items-center justify-content-between border-0 pt-3 px-4 pb-0">
          <h6 class="mb-0 fw-semibold">Delivery History — {{ deliveriesReport?.name }}</h6>
          <button class="btn-icon" @click="showDeliveriesModal = false">✕</button>
        </div>
        <div class="card-body px-4 pb-4">
          <div v-if="deliveriesLoading" class="text-center py-3">
            <div class="spinner-border text-primary spinner-border-sm" role="status"></div>
          </div>
          <div v-else-if="deliveries.length === 0" class="text-muted small text-center py-3">No deliveries yet.</div>
          <div v-else class="deliveries-list">
            <div v-for="d in deliveries" :key="d.id" class="delivery-row d-flex align-items-start gap-3 py-2 border-bottom">
              <span class="delivery-badge" :class="d.status === 'sent' ? 'success' : 'danger'">
                {{ d.status === 'sent' ? '✓' : '✗' }}
              </span>
              <div class="min-w-0">
                <div class="small fw-medium">{{ d.status === 'sent' ? 'Sent' : 'Failed' }}</div>
                <div v-if="d.delivered_at" class="text-muted" style="font-size:.72rem;">{{ formatDateTime(d.delivered_at) }}</div>
                <div v-if="d.failure_reason" class="text-danger" style="font-size:.72rem;">{{ d.failure_reason }}</div>
              </div>
              <div class="text-muted ms-auto" style="font-size:.72rem;white-space:nowrap;">{{ formatDateTime(d.created_at) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Report Preview Modal ───────────────────────────────────────────── -->
    <div v-if="showPreviewModal" class="modal-backdrop-custom" @click.self="showPreviewModal = false">
      <div class="modal-dialog-preview card shadow">
        <div class="card-header d-flex align-items-center justify-content-between border-0 pt-3 px-4 pb-0">
          <div>
            <h6 class="mb-0 fw-semibold">Report Preview</h6>
            <div class="text-muted small mt-1">{{ previewReport?.name }}</div>
          </div>
          <button class="btn-icon" @click="showPreviewModal = false">✕</button>
        </div>
        <div class="card-body px-4 pb-4">
          <!-- Loading skeleton -->
          <div v-if="previewLoading" class="py-4">
            <div class="skeleton-line mb-3" style="width:60%;height:14px;"></div>
            <div class="skeleton-line mb-2" style="width:100%;height:10px;"></div>
            <div class="skeleton-line mb-2" style="width:100%;height:10px;"></div>
            <div class="skeleton-line mb-2" style="width:80%;height:10px;"></div>
            <div class="skeleton-line mt-4 mb-2" style="width:40%;height:14px;"></div>
            <div class="skeleton-line mb-2" style="width:100%;height:10px;"></div>
            <div class="skeleton-line mb-2" style="width:100%;height:10px;"></div>
          </div>

          <!-- Preview content -->
          <div v-else-if="previewData" class="preview-email-body">
            <!-- Header section -->
            <div class="preview-header mb-4 p-3 rounded" style="background:#f0efff;border-left:4px solid #635bff;">
              <div class="fw-semibold" style="color:#635bff;">Analytics Report</div>
              <div class="text-muted small mt-1">
                <span class="text-capitalize fw-medium">{{ previewReport?.frequency }}</span> summary &mdash;
                {{ previewPeriodLabel }}
              </div>
            </div>

            <!-- Stats row -->
            <div class="row g-3 mb-4">
              <div class="col-4">
                <div class="preview-stat-card text-center p-3 rounded border">
                  <div class="fw-bold" style="font-size:1.4rem;color:#635bff;">{{ previewData.total_clicks.toLocaleString() }}</div>
                  <div class="text-muted small mt-1">Total Clicks</div>
                </div>
              </div>
              <div class="col-4">
                <div class="preview-stat-card text-center p-3 rounded border">
                  <div class="fw-bold" style="font-size:1.4rem;color:#635bff;">{{ previewData.total_links.toLocaleString() }}</div>
                  <div class="text-muted small mt-1">Active Links</div>
                </div>
              </div>
              <div class="col-4">
                <div class="preview-stat-card text-center p-3 rounded border">
                  <div class="fw-bold" style="font-size:1.4rem;color:#635bff;">{{ previewData.clicks_7_days.toLocaleString() }}</div>
                  <div class="text-muted small mt-1">Last 7 Days</div>
                </div>
              </div>
            </div>

            <!-- Top links in this report -->
            <div class="mb-4">
              <div class="fw-semibold small mb-2">This report would include:</div>
              <div class="d-flex flex-wrap gap-1 mb-3">
                <span
                  v-for="lid in previewReport?.link_ids"
                  :key="lid"
                  class="badge bg-light text-secondary border small"
                  style="font-size:.72rem;"
                >{{ lid.slice(0, 8) }}…</span>
              </div>
            </div>

            <!-- Top performing links from overview -->
            <div v-if="previewData.top_links && previewData.top_links.length > 0" class="mb-4">
              <div class="fw-semibold small mb-2">Top Performing Links (Account-wide)</div>
              <div class="preview-links-list">
                <div
                  v-for="(link, idx) in previewData.top_links.slice(0, 5)"
                  :key="link.id"
                  class="d-flex align-items-center gap-2 py-2 border-bottom"
                >
                  <span class="text-muted fw-bold" style="width:20px;font-size:.75rem;">{{ idx + 1 }}</span>
                  <div class="flex-fill min-w-0">
                    <div class="small fw-medium text-truncate">{{ link.title || link.slug }}</div>
                    <div class="text-muted" style="font-size:.7rem;">{{ link.short_url }}</div>
                  </div>
                  <div class="text-end flex-shrink-0">
                    <div class="small fw-semibold" style="color:#635bff;">{{ (link.click_count || 0).toLocaleString() }}</div>
                    <div class="text-muted" style="font-size:.68rem;">clicks</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Delivery schedule note -->
            <div class="preview-schedule-note small text-muted p-3 rounded" style="background:#f9fafb;border:1px solid #e5e7eb;">
              <span class="fw-medium">Email delivery schedule:</span>
              <span v-if="previewReport?.frequency === 'daily'"> This report is sent every day.</span>
              <span v-else-if="previewReport?.frequency === 'weekly'"> This report is sent every Monday.</span>
              <span v-else-if="previewReport?.frequency === 'monthly'"> This report is sent on the 1st of each month.</span>
              <span v-if="previewReport?.next_run_at"> Next delivery: <strong>{{ formatDate(previewReport.next_run_at) }}</strong>.</span>
            </div>
          </div>

          <!-- Error fallback -->
          <div v-else class="text-center text-muted py-4 small">
            Could not load preview data.
          </div>

          <!-- Footer actions -->
          <div class="d-flex align-items-center justify-content-between gap-2 mt-4">
            <div class="text-muted small">Preview shows current account analytics data.</div>
            <div class="d-flex gap-2">
              <button class="btn btn-sm btn-outline-secondary" @click="showPreviewModal = false">Close</button>
              <button
                class="btn btn-sm btn-primary"
                :disabled="sendingPreview"
                @click="sendNowFromPreview"
                style="background:#635bff;border-color:#635bff;"
              >
                <span v-if="sendingPreview" class="spinner-border spinner-border-sm me-1" role="status"></span>
                Send now
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, onUnmounted } from 'vue';
import { reportsApi } from '@/api/reports';
import { linksApi } from '@/api/links';
import { dashboardApi } from '@/api/dashboard';
import type { ReportResponse, ReportDelivery } from '@/types/reports';
import type { LinkResponse } from '@/types/links';
import type { DashboardOverviewResponse } from '@/types/dashboard';

// ── State ─────────────────────────────────────────────────────────────────────
const reports = ref<ReportResponse[]>([]);
const loading = ref(false);
const error = ref('');
const actionLoading = ref('');

const allLinks = ref<LinkResponse[]>([]);
const linkSearch = ref('');
const selectorOpen = ref(false);
const selectorRef = ref<HTMLElement | null>(null);

const showFormModal = ref(false);
const editingReport = ref<ReportResponse | null>(null);
const formLoading = ref(false);
const formError = ref('');

const form = reactive({
  name: '',
  link_ids: [] as string[],
  frequency: 'weekly' as 'daily' | 'weekly' | 'monthly',
});

const showDeliveriesModal = ref(false);
const deliveriesReport = ref<ReportResponse | null>(null);
const deliveries = ref<ReportDelivery[]>([]);
const deliveriesLoading = ref(false);

// Preview modal state
const showPreviewModal = ref(false);
const previewReport = ref<ReportResponse | null>(null);
const previewLoading = ref(false);
const previewData = ref<DashboardOverviewResponse | null>(null);
const sendingPreview = ref(false);

const frequencies = [
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
];

// ── Computed ──────────────────────────────────────────────────────────────────
const filteredAllLinks = computed(() => {
  const q = linkSearch.value.toLowerCase();
  if (!q) return allLinks.value;
  return allLinks.value.filter(
    (l) => (l.title || '').toLowerCase().includes(q) || l.slug.toLowerCase().includes(q),
  );
});

const previewPeriodLabel = computed(() => {
  const freq = previewReport.value?.frequency;
  if (freq === 'daily') return 'Last 24 hours';
  if (freq === 'weekly') return 'Last 7 days';
  if (freq === 'monthly') return 'Last 30 days';
  return 'Recent period';
});

// ── Helpers ───────────────────────────────────────────────────────────────────
function capitalize(s: string): string {
  if (!s) return s;
  return s[0].toUpperCase() + s.slice(1);
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
}

// ── Data loading ──────────────────────────────────────────────────────────────
async function loadReports() {
  loading.value = true;
  error.value = '';
  try {
    const resp = await reportsApi.list();
    reports.value = resp.data ?? [];
  } catch {
    error.value = 'Failed to load reports.';
  } finally {
    loading.value = false;
  }
}

async function loadLinks() {
  try {
    const resp = await linksApi.list(1, 200);
    allLinks.value = resp.data?.links ?? [];
  } catch {
    // non-fatal
  }
}

// ── Form ──────────────────────────────────────────────────────────────────────
function openCreate() {
  editingReport.value = null;
  form.name = '';
  form.link_ids = [];
  form.frequency = 'weekly';
  formError.value = '';
  showFormModal.value = true;
}

function openEdit(report: ReportResponse) {
  editingReport.value = report;
  form.name = report.name;
  form.link_ids = [...report.link_ids];
  form.frequency = report.frequency;
  formError.value = '';
  showFormModal.value = true;
}

function closeForm() {
  showFormModal.value = false;
  selectorOpen.value = false;
}

async function submitForm() {
  if (!form.name.trim() || form.link_ids.length === 0) return;
  formLoading.value = true;
  formError.value = '';
  try {
    if (editingReport.value) {
      const resp = await reportsApi.update(editingReport.value.id, {
        name: form.name,
        link_ids: form.link_ids,
        frequency: form.frequency,
      });
      const idx = reports.value.findIndex((r) => r.id === editingReport.value!.id);
      if (idx !== -1 && resp.data) reports.value[idx] = resp.data;
    } else {
      const resp = await reportsApi.create({
        name: form.name,
        link_ids: form.link_ids,
        frequency: form.frequency,
      });
      if (resp.data) reports.value.unshift(resp.data);
    }
    closeForm();
  } catch (e: unknown) {
    formError.value = (e as { response?: { data?: { message?: string } } })?.response?.data?.message ?? 'Failed to save report.';
  } finally {
    formLoading.value = false;
  }
}

// ── Actions ───────────────────────────────────────────────────────────────────
async function toggleActive(report: ReportResponse) {
  actionLoading.value = report.id;
  try {
    const resp = await reportsApi.update(report.id, { is_active: !report.is_active });
    const idx = reports.value.findIndex((r) => r.id === report.id);
    if (idx !== -1 && resp.data) reports.value[idx] = resp.data;
  } catch {
    // noop
  } finally {
    actionLoading.value = '';
  }
}

async function sendNow(id: string) {
  actionLoading.value = id;
  try {
    await reportsApi.sendNow(id);
  } catch {
    // noop — errors surfaced via delivery history
  } finally {
    actionLoading.value = '';
  }
}

async function deleteReport(id: string) {
  if (!confirm('Delete this report?')) return;
  try {
    await reportsApi.delete(id);
    reports.value = reports.value.filter((r) => r.id !== id);
  } catch {
    // noop
  }
}

async function openDeliveries(report: ReportResponse) {
  deliveriesReport.value = report;
  deliveries.value = [];
  deliveriesLoading.value = true;
  showDeliveriesModal.value = true;
  try {
    const resp = await reportsApi.getDeliveries(report.id);
    deliveries.value = resp.data ?? [];
  } catch {
    // noop
  } finally {
    deliveriesLoading.value = false;
  }
}

async function openPreview(report: ReportResponse) {
  previewReport.value = report;
  previewData.value = null;
  previewLoading.value = true;
  showPreviewModal.value = true;
  try {
    const resp = await dashboardApi.getOverview();
    previewData.value = resp.data ?? null;
  } catch {
    // noop — fallback shown in template
  } finally {
    previewLoading.value = false;
  }
}

async function sendNowFromPreview() {
  if (!previewReport.value) return;
  sendingPreview.value = true;
  try {
    await reportsApi.sendNow(previewReport.value.id);
    showPreviewModal.value = false;
  } catch {
    // noop
  } finally {
    sendingPreview.value = false;
  }
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────
function handleClickOutside(e: MouseEvent) {
  if (selectorRef.value && !selectorRef.value.contains(e.target as Node)) {
    selectorOpen.value = false;
  }
}

onMounted(() => {
  loadReports();
  loadLinks();
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped lang="scss">
$primary: #635bff;

// ── Report card ────────────────────────────────────────────────────────────────
.report-card {
  transition: box-shadow 0.15s;
  &:hover { box-shadow: 0 4px 16px rgba(0,0,0,.08) !important; }
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  &.active  { background: #22c55e; }
  &.inactive { background: #d1d5db; }
}

.badge-freq {
  display: inline-block;
  background: #eef0ff;
  color: $primary;
  border-radius: 999px;
  padding: .1em .55em;
  font-size: .72rem;
  font-weight: 600;
}

// ── Link selector ─────────────────────────────────────────────────────────────
.link-selector { position: relative; }

.selector-trigger {
  cursor: pointer;
  user-select: none;
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 500;
  border-radius: 8px;
  overflow: hidden;
}

.selector-options {
  max-height: 180px;
  overflow-y: auto;
}

.selector-option {
  cursor: pointer;
  font-size: .8125rem;
  transition: background .12s;
  &:hover { background: #f7f9fc; }
  &.disabled { opacity: .5; cursor: not-allowed; }
}

// ── Modal ─────────────────────────────────────────────────────────────────────
.modal-backdrop-custom {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.45);
  z-index: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-dialog-custom {
  width: 100%;
  max-width: 480px;
  border-radius: 12px;
  overflow: hidden;
}

.btn-icon {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: .8rem;
  color: #697386;
  cursor: pointer;
  &:hover { background: #f7f9fc; }
}

// ── Preview modal ─────────────────────────────────────────────────────────────
.modal-dialog-preview {
  width: 100%;
  max-width: 600px;
  border-radius: 12px;
  overflow: hidden;
  max-height: 90vh;
  overflow-y: auto;
}

.preview-stat-card {
  background: #fafafa;
  transition: box-shadow 0.12s;
  &:hover { box-shadow: 0 2px 8px rgba(0,0,0,.06); }
}

.preview-links-list .border-bottom:last-child { border-bottom: none !important; }

// Skeleton loading
.skeleton-line {
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
  border-radius: 4px;
  display: block;
}

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

// ── Delivery history ──────────────────────────────────────────────────────────
.deliveries-list { max-height: 360px; overflow-y: auto; }

.delivery-badge {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  font-size: .7rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  &.success { background: #dcfce7; color: #16a34a; }
  &.danger  { background: #fee2e2; color: #dc2626; }
}

.delivery-row:last-child { border-bottom: none !important; }

.min-w-0 { min-width: 0; }

// ── Dark mode ─────────────────────────────────────────────────────────────────
:global(.dark-mode) {
  .report-card {
    background: #161b22;
    border-color: #30363d !important;
  }

  .badge-freq {
    background: rgba(99, 91, 255, .2);
  }

  .selector-dropdown {
    background: #161b22;
    border-color: #30363d !important;
  }

  .selector-option:hover { background: #21262d; }

  .selector-trigger {
    background: #161b22;
    border-color: #30363d;
    color: #e6edf3;
  }

  .modal-dialog-custom {
    background: #161b22;
    border-color: #30363d !important;
  }

  .card-header {
    border-color: #30363d !important;
    background: #161b22;
    color: #e6edf3;
  }

  .delivery-row { border-color: #30363d !important; }

  .modal-dialog-preview {
    background: #161b22;
    border-color: #30363d !important;
  }

  .preview-stat-card {
    background: #21262d;
    border-color: #30363d !important;
  }

  .preview-header {
    background: rgba(99, 91, 255, .15) !important;
  }

  .preview-schedule-note {
    background: #21262d !important;
    border-color: #30363d !important;
  }

  .skeleton-line {
    background: linear-gradient(90deg, #21262d 25%, #30363d 50%, #21262d 75%);
    background-size: 200% 100%;
  }
}
</style>
