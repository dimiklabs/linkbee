<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Email Reports</h1>
        <p class="page-subtitle">Automatically receive analytics summaries by email on a set cadence.</p>
      </div>
      <md-filled-button @click="openCreate">
        <span class="material-symbols-outlined" slot="icon">add</span>
        New Report
      </md-filled-button>
    </div>

    <!-- Error -->
    <div v-if="error" style="display:flex;align-items:center;gap:10px;padding:12px 16px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:12px;margin-bottom:16px;font-size:0.875rem;">
      <span class="material-symbols-outlined">warning</span>
      {{ error }}
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:40px" />
    </div>

    <!-- Empty state -->
    <div v-else-if="!loading && reports.length === 0 && !error" class="m3-card m3-card--elevated m3-empty-state">
      <div class="m3-empty-state__icon">
        <span class="material-symbols-outlined">email</span>
      </div>
      <div class="md-title-medium" style="margin-bottom:8px;">No scheduled reports yet</div>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0 0 20px;">Set up automated email summaries to track your link performance.</p>
      <md-filled-button @click="openCreate">
        <span class="material-symbols-outlined" slot="icon">add</span>
        Create your first report
      </md-filled-button>
    </div>

    <!-- Reports list -->
    <div v-else class="reports-list">
      <div v-for="report in reports" :key="report.id" class="m3-card m3-card--elevated report-card">
        <div style="display:flex;align-items:center;gap:12px;flex-wrap:wrap;padding:16px 20px;">
          <!-- Status dot -->
          <span
            style="width:10px;height:10px;border-radius:50%;flex-shrink:0;"
            :style="{ background: report.is_active ? '#22c55e' : 'var(--md-sys-color-outline-variant)' }"
            :title="report.is_active ? 'Active' : 'Paused'"
          ></span>

          <!-- Info -->
          <div style="flex:1;min-width:0;">
            <div style="font-weight:600;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ report.name }}</div>
            <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;margin-top:4px;display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
              <span class="freq-badge">{{ capitalize(report.frequency) }}</span>
              <span>{{ report.link_ids.length }} link{{ report.link_ids.length !== 1 ? 's' : '' }}</span>
              <span v-if="report.next_run_at">· Next: {{ formatDate(report.next_run_at) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div style="display:flex;align-items:center;gap:6px;flex-shrink:0;flex-wrap:wrap;">
            <md-outlined-button
              :title="report.is_active ? 'Pause report' : 'Resume report'"
              :disabled="actionLoading === report.id"
              @click="toggleActive(report)"
            >
              {{ report.is_active ? 'Pause' : 'Resume' }}
            </md-outlined-button>
            <md-outlined-button
              title="Send report immediately"
              :disabled="actionLoading === report.id"
              @click="sendNow(report.id)"
            >
              <md-circular-progress v-if="actionLoading === report.id" indeterminate style="--md-circular-progress-size:16px;margin-right:4px;" />
              Send now
            </md-outlined-button>
            <md-icon-button title="Preview report" @click="openPreview(report)">
              <span class="material-symbols-outlined">visibility</span>
            </md-icon-button>
            <md-icon-button title="View delivery history" @click="openDeliveries(report)">
              <span class="material-symbols-outlined">history</span>
            </md-icon-button>
            <md-icon-button title="Edit report" @click="openEdit(report)">
              <span class="material-symbols-outlined">edit</span>
            </md-icon-button>
            <md-icon-button title="Delete report" @click="deleteReport(report.id)" style="--md-icon-button-icon-color:var(--md-sys-color-error);">
              <span class="material-symbols-outlined">delete</span>
            </md-icon-button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Create / Edit Dialog ──────────────────────────────────────────── -->
    <md-dialog :open="showFormModal" @closed="closeForm">
      <div slot="headline">{{ editingReport ? 'Edit Report' : 'New Scheduled Report' }}</div>
      <div slot="content" style="min-width:480px;max-width:100%;display:flex;flex-direction:column;gap:16px;">

        <!-- Name -->
        <md-outlined-text-field
          :value="form.name"
          @input="form.name=($event.target as HTMLInputElement).value"
          label="Report Name *"
          placeholder="e.g. Weekly Top Links"
          maxlength="255"
          style="width:100%;"
        />

        <!-- Links selector -->
        <div>
          <div style="font-size:0.875rem;font-weight:500;margin-bottom:8px;color:var(--md-sys-color-on-surface);">Links * <span style="color:var(--md-sys-color-on-surface-variant);font-weight:400;">(select 1–10)</span></div>
          <div class="link-selector" ref="selectorRef">
            <div
              class="selector-trigger"
              @click="selectorOpen = !selectorOpen"
            >
              <span v-if="form.link_ids.length === 0" style="color:var(--md-sys-color-on-surface-variant);">Choose links…</span>
              <span v-else>{{ form.link_ids.length }} link{{ form.link_ids.length !== 1 ? 's' : '' }} selected</span>
              <span class="material-symbols-outlined" style="font-size:18px;flex-shrink:0;">expand_more</span>
            </div>
            <div v-if="selectorOpen" class="selector-dropdown m3-card m3-card--outlined">
              <div style="padding:8px 8px 0;">
                <md-outlined-text-field
                  :value="linkSearch"
                  @input="linkSearch=($event.target as HTMLInputElement).value"
                  label="Search"
                  style="width:100%;"
                  @click.stop
                />
              </div>
              <div class="selector-options">
                <label
                  v-for="link in filteredAllLinks"
                  :key="link.id"
                  class="selector-option"
                  :class="{ disabled: !form.link_ids.includes(link.id) && form.link_ids.length >= 10 }"
                  @click.stop
                >
                  <input
                    type="checkbox"
                    :value="link.id"
                    v-model="form.link_ids"
                    :disabled="!form.link_ids.includes(link.id) && form.link_ids.length >= 10"
                    style="accent-color:var(--md-sys-color-primary);flex-shrink:0;"
                  />
                  <div style="min-width:0;">
                    <div style="font-size:0.8125rem;font-weight:500;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ link.title || link.slug }}</div>
                    <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.7rem;">{{ link.short_url }}</div>
                  </div>
                </label>
                <div v-if="filteredAllLinks.length === 0" style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;text-align:center;padding:16px;">No links found</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Frequency -->
        <div>
          <div style="font-size:0.875rem;font-weight:500;margin-bottom:8px;color:var(--md-sys-color-on-surface);">Frequency *</div>
          <div style="display:flex;gap:8px;">
            <md-outlined-button
              v-for="freq in frequencies"
              :key="freq.value"
              @click="form.frequency = freq.value"
              :style="form.frequency === freq.value ? 'background:var(--md-sys-color-primary-container);--md-outlined-button-label-text-color:var(--md-sys-color-on-primary-container);--md-outlined-button-outline-color:var(--md-sys-color-primary);' : ''"
              style="flex:1;"
            >
              {{ freq.label }}
            </md-outlined-button>
          </div>
        </div>

        <!-- Error -->
        <div v-if="formError" style="padding:10px 14px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:8px;font-size:0.875rem;">
          {{ formError }}
        </div>
      </div>
      <div slot="actions">
        <md-text-button @click="closeForm">Cancel</md-text-button>
        <md-filled-button
          :disabled="formLoading || !form.name.trim() || form.link_ids.length === 0"
          @click="submitForm"
        >
          <md-circular-progress v-if="formLoading" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
          {{ editingReport ? 'Save Changes' : 'Create Report' }}
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- ── Delivery History Dialog ─────────────────────────────────────── -->
    <md-dialog :open="showDeliveriesModal" @closed="showDeliveriesModal = false">
      <div slot="headline">Delivery History — {{ deliveriesReport?.name }}</div>
      <div slot="content" style="min-width:480px;max-width:100%;min-height:200px;">
        <div v-if="deliveriesLoading" style="display:flex;justify-content:center;padding:24px;">
          <md-circular-progress indeterminate style="--md-circular-progress-size:32px" />
        </div>
        <div v-else-if="deliveries.length === 0" style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;text-align:center;padding:24px;">No deliveries yet.</div>
        <div v-else style="max-height:360px;overflow-y:auto;">
          <div
            v-for="d in deliveries"
            :key="d.id"
            style="display:flex;align-items:flex-start;gap:12px;padding:10px 0;border-bottom:1px solid var(--md-sys-color-outline-variant);"
          >
            <span
              style="width:22px;height:22px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.7rem;font-weight:700;flex-shrink:0;"
              :style="d.status === 'sent' ? 'background:#dcfce7;color:#16a34a;' : 'background:#fee2e2;color:#dc2626;'"
            >{{ d.status === 'sent' ? '✓' : '✗' }}</span>
            <div style="min-width:0;flex:1;">
              <div style="font-size:0.875rem;font-weight:500;">{{ d.status === 'sent' ? 'Sent' : 'Failed' }}</div>
              <div v-if="d.delivered_at" style="color:var(--md-sys-color-on-surface-variant);font-size:0.72rem;">{{ formatDateTime(d.delivered_at) }}</div>
              <div v-if="d.failure_reason" style="color:var(--md-sys-color-error);font-size:0.72rem;">{{ d.failure_reason }}</div>
            </div>
            <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.72rem;white-space:nowrap;flex-shrink:0;">{{ formatDateTime(d.created_at) }}</div>
          </div>
        </div>
      </div>
      <div slot="actions">
        <md-text-button @click="showDeliveriesModal = false">Close</md-text-button>
      </div>
    </md-dialog>

    <!-- ── Preview Dialog ─────────────────────────────────────────────── -->
    <md-dialog :open="showPreviewModal" @closed="showPreviewModal = false">
      <div slot="headline">Report Preview</div>
      <div slot="content" style="min-width:540px;max-width:100%;max-height:70vh;overflow-y:auto;">
        <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;margin-bottom:16px;">{{ previewReport?.name }}</div>

        <!-- Loading skeleton -->
        <div v-if="previewLoading" style="padding:16px 0;">
          <div class="skeleton-line" style="width:60%;height:14px;margin-bottom:12px;"></div>
          <div class="skeleton-line" style="width:100%;height:10px;margin-bottom:8px;"></div>
          <div class="skeleton-line" style="width:100%;height:10px;margin-bottom:8px;"></div>
          <div class="skeleton-line" style="width:80%;height:10px;margin-bottom:16px;"></div>
          <div class="skeleton-line" style="width:40%;height:14px;margin-bottom:8px;"></div>
          <div class="skeleton-line" style="width:100%;height:10px;margin-bottom:8px;"></div>
          <div class="skeleton-line" style="width:100%;height:10px;"></div>
        </div>

        <!-- Preview content -->
        <div v-else-if="previewData">
          <!-- Header -->
          <div style="background:#f0efff;border-left:4px solid var(--md-sys-color-primary);padding:12px 16px;border-radius:0 8px 8px 0;margin-bottom:20px;">
            <div style="font-weight:600;color:var(--md-sys-color-primary);">Analytics Report</div>
            <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;margin-top:4px;">
              <span style="text-transform:capitalize;font-weight:500;">{{ previewReport?.frequency }}</span> summary — {{ previewPeriodLabel }}
            </div>
          </div>

          <!-- Stats -->
          <div class="preview-stats-grid">
            <div class="m3-card m3-card--outlined preview-stat-card">
              <div style="font-size:1.4rem;font-weight:700;color:var(--md-sys-color-primary);">{{ previewData.total_clicks.toLocaleString() }}</div>
              <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;margin-top:4px;">Total Clicks</div>
            </div>
            <div class="m3-card m3-card--outlined preview-stat-card">
              <div style="font-size:1.4rem;font-weight:700;color:var(--md-sys-color-primary);">{{ previewData.total_links.toLocaleString() }}</div>
              <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;margin-top:4px;">Active Links</div>
            </div>
            <div class="m3-card m3-card--outlined preview-stat-card">
              <div style="font-size:1.4rem;font-weight:700;color:var(--md-sys-color-primary);">{{ previewData.clicks_7_days.toLocaleString() }}</div>
              <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;margin-top:4px;">Last 7 Days</div>
            </div>
          </div>

          <!-- Report links -->
          <div style="margin-bottom:16px;">
            <div style="font-weight:600;font-size:0.875rem;margin-bottom:8px;">This report would include:</div>
            <div style="display:flex;flex-wrap:wrap;gap:4px;">
              <span
                v-for="lid in previewReport?.link_ids"
                :key="lid"
                class="m3-badge m3-badge--neutral"
                style="font-size:0.72rem;"
              >{{ lid.slice(0, 8) }}…</span>
            </div>
          </div>

          <!-- Top links -->
          <div v-if="previewData.top_links && previewData.top_links.length > 0" style="margin-bottom:16px;">
            <div style="font-weight:600;font-size:0.875rem;margin-bottom:8px;">Top Performing Links (Account-wide)</div>
            <div>
              <div
                v-for="(link, idx) in previewData.top_links.slice(0, 5)"
                :key="link.id"
                style="display:flex;align-items:center;gap:8px;padding:8px 0;border-bottom:1px solid var(--md-sys-color-outline-variant);"
              >
                <span style="width:20px;font-size:0.75rem;font-weight:700;color:var(--md-sys-color-on-surface-variant);">{{ idx + 1 }}</span>
                <div style="flex:1;min-width:0;">
                  <div style="font-size:0.875rem;font-weight:500;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ link.title || link.slug }}</div>
                  <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.7rem;">{{ link.short_url }}</div>
                </div>
                <div style="text-align:right;flex-shrink:0;">
                  <div style="font-size:0.875rem;font-weight:600;color:var(--md-sys-color-primary);">{{ (link.click_count || 0).toLocaleString() }}</div>
                  <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.68rem;">clicks</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Schedule note -->
          <div style="background:var(--md-sys-color-surface-container-low);border:1px solid var(--md-sys-color-outline-variant);padding:12px 16px;border-radius:8px;font-size:0.875rem;color:var(--md-sys-color-on-surface-variant);">
            <span style="font-weight:500;color:var(--md-sys-color-on-surface);">Email delivery schedule:</span>
            <span v-if="previewReport?.frequency === 'daily'"> This report is sent every day.</span>
            <span v-else-if="previewReport?.frequency === 'weekly'"> This report is sent every Monday.</span>
            <span v-else-if="previewReport?.frequency === 'monthly'"> This report is sent on the 1st of each month.</span>
            <span v-if="previewReport?.next_run_at"> Next delivery: <strong>{{ formatDate(previewReport.next_run_at) }}</strong>.</span>
          </div>
        </div>

        <!-- Error fallback -->
        <div v-else style="text-align:center;color:var(--md-sys-color-on-surface-variant);padding:32px;font-size:0.875rem;">
          Could not load preview data.
        </div>

        <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;margin-top:16px;">Preview shows current account analytics data.</div>
      </div>
      <div slot="actions">
        <md-text-button @click="showPreviewModal = false">Close</md-text-button>
        <md-filled-button :disabled="sendingPreview" @click="sendNowFromPreview">
          <md-circular-progress v-if="sendingPreview" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
          Send now
        </md-filled-button>
      </div>
    </md-dialog>

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
.page-wrapper {
  padding: 24px;
  max-width: 900px;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 4px;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0;
}

/* Cards */
.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
  }

  &--outlined {
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* Badges */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* Reports list */
.reports-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Report card */
.report-card {
  transition: box-shadow 0.15s;

  &:hover {
    box-shadow: 0 4px 16px rgba(0,0,0,0.10);
  }
}

/* Frequency badge */
.freq-badge {
  display: inline-block;
  background: rgba(99, 91, 255, 0.10);
  color: var(--md-sys-color-primary, #635bff);
  border-radius: 999px;
  padding: 0.1em 0.55em;
  font-size: 0.72rem;
  font-weight: 600;
}

/* Empty state */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 24px;
  text-align: center;

  &__icon {
    width: 72px;
    height: 72px;
    border-radius: 50%;
    background: var(--md-sys-color-surface-container-low);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 16px;

    .material-symbols-outlined {
      font-size: 2rem;
      color: var(--md-sys-color-on-surface-variant);
      opacity: 0.6;
    }
  }
}

/* Link selector */
.link-selector {
  position: relative;
}

.selector-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border: 1px solid var(--md-sys-color-outline);
  border-radius: 4px;
  cursor: pointer;
  user-select: none;
  font-size: 0.9375rem;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  transition: background 0.12s;

  &:hover {
    background: var(--md-sys-color-surface-container-low);
  }
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 500;
  border-radius: 8px !important;
  overflow: hidden;
}

.selector-options {
  max-height: 200px;
  overflow-y: auto;
}

.selector-option {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  font-size: 0.8125rem;
  transition: background 0.12s;

  &:hover {
    background: var(--md-sys-color-surface-container-low);
  }

  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

/* Preview stats */
.preview-stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.preview-stat-card {
  padding: 16px;
  text-align: center;
}

/* Skeleton loading */
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
</style>
