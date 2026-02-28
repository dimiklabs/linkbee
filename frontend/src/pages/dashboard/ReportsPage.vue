<template>
  <div class="reports-page">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Email Reports</h1>
        <p class="page-subtitle">Automatically receive analytics summaries by email on a set cadence.</p>
      </div>
      <button class="btn-filled" @click="openCreate">
        <span class="material-symbols-outlined">add</span>
        New Report
      </button>
    </div>

    <!-- Error -->
    <div v-if="error" class="error-banner">
      <span class="material-symbols-outlined">warning</span>
      {{ error }}
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <span class="css-spinner"></span>
    </div>

    <!-- Empty state -->
    <div v-else-if="!loading && reports.length === 0 && !error" class="an-card empty-state-card">
      <div class="empty-icon">
        <span class="material-symbols-outlined">email</span>
      </div>
      <div class="empty-title">No scheduled reports yet</div>
      <p class="empty-desc">Set up automated email summaries to track your link performance.</p>
      <button class="btn-filled" @click="openCreate">
        <span class="material-symbols-outlined">add</span>
        Create your first report
      </button>
    </div>

    <!-- Reports list -->
    <div v-else class="reports-list">
      <div v-for="report in reports" :key="report.id" class="report-card an-card">
        <div class="report-card__body">
          <!-- Status dot -->
          <span
            class="status-dot"
            :style="{ background: report.is_active ? '#22c55e' : 'var(--md-sys-color-outline-variant)' }"
            :title="report.is_active ? 'Active' : 'Paused'"
          ></span>

          <!-- Info -->
          <div class="report-info">
            <div class="report-info__name">{{ report.name }}</div>
            <div class="report-info__meta">
              <span class="freq-badge">{{ capitalize(report.frequency) }}</span>
              <span>{{ report.link_ids.length }} link{{ report.link_ids.length !== 1 ? 's' : '' }}</span>
              <span v-if="report.next_run_at">· Next: {{ formatDate(report.next_run_at) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="report-actions">
            <button class="btn-outlined"
              :title="report.is_active ? 'Pause report' : 'Resume report'"
              :disabled="actionLoading === report.id"
              @click="toggleActive(report)"
            >
              {{ report.is_active ? 'Pause' : 'Resume' }}
            </button>
            <button class="btn-outlined"
              title="Send report immediately"
              :disabled="actionLoading === report.id"
              @click="sendNow(report.id)"
            >
              <span v-if="actionLoading === report.id" class="css-spinner css-spinner--sm"></span>
              Send now
            </button>
            <button class="btn-icon" title="Preview report" @click="openPreview(report)">
              <span class="material-symbols-outlined">visibility</span>
            </button>
            <button class="btn-icon" title="View delivery history" @click="openDeliveries(report)">
              <span class="material-symbols-outlined">history</span>
            </button>
            <button class="btn-icon" title="Edit report" @click="openEdit(report)">
              <span class="material-symbols-outlined">edit</span>
            </button>
            <button class="btn-icon btn-sm btn-danger" title="Delete report" @click="deleteReport(report.id)">
              <span class="material-symbols-outlined">delete</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Create / Edit Dialog ──────────────────────────────────────────── -->
    <BaseModal v-model="showFormModal" size="md" @closed="closeForm">
      <template #headline>
        {{ editingReport ? 'Edit Report' : 'New Scheduled Report' }}
      </template>

      <div class="modal-form">

        <!-- Name -->
        <div class="form-field">
          <label class="form-field__label">Report Name *</label>
          <input
            class="form-input"
            :value="form.name"
            @input="form.name=($event.target as HTMLInputElement).value"
            placeholder="e.g. Weekly Top Links"
            maxlength="255"
          />
        </div>

        <!-- Links selector -->
        <div>
          <div class="form-field__label" style="margin-bottom:8px;">
            Links * <span style="color:var(--md-sys-color-on-surface-variant);font-weight:400;">(select 1–10)</span>
          </div>
          <div class="link-selector" ref="selectorRef">
            <div class="selector-trigger" @click="selectorOpen = !selectorOpen">
              <span v-if="form.link_ids.length === 0" class="selector-placeholder">Choose links…</span>
              <span v-else>{{ form.link_ids.length }} link{{ form.link_ids.length !== 1 ? 's' : '' }} selected</span>
              <span class="material-symbols-outlined selector-chevron">expand_more</span>
            </div>
            <div v-if="selectorOpen" class="selector-dropdown">
              <div style="padding:8px 8px 0;">
                <input
                  class="form-input"
                  :value="linkSearch"
                  @input="linkSearch=($event.target as HTMLInputElement).value"
                  placeholder="Search…"
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
                    <div class="selector-option__title">{{ link.title || link.slug }}</div>
                    <div class="selector-option__url">{{ link.short_url }}</div>
                  </div>
                </label>
                <div v-if="filteredAllLinks.length === 0" class="selector-empty">No links found</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Frequency -->
        <div>
          <div class="form-field__label" style="margin-bottom:8px;">Frequency *</div>
          <div class="freq-buttons">
            <button :class="['btn-outlined', form.frequency === freq.value ? 'btn-active' : '']"
              v-for="freq in frequencies"
              :key="freq.value"
              @click="form.frequency = freq.value as 'daily' | 'weekly' | 'monthly'"
              style="flex:1;"
            >
              {{ freq.label }}
            </button>
          </div>
        </div>

        <!-- Error -->
        <div v-if="formError" class="error-box">{{ formError }}</div>
      </div>
      <template #actions>
        <button class="btn-text" @click="closeForm">Cancel</button>
        <button class="btn-filled"
          :disabled="formLoading || !form.name.trim() || form.link_ids.length === 0"
          @click="submitForm"
        >
          <span v-if="formLoading" class="css-spinner css-spinner--sm css-spinner--white"></span>
          {{ editingReport ? 'Save Changes' : 'Create Report' }}
        </button>
      </template>
    </BaseModal>

    <!-- ── Delivery History Dialog ─────────────────────────────────────── -->
    <BaseModal v-model="showDeliveriesModal" size="md">
      <template #headline>
        Delivery History — {{ deliveriesReport?.name }}
      </template>

      <div class="modal-deliveries">
        <div v-if="deliveriesLoading" class="loading-state">
          <span class="css-spinner"></span>
        </div>
        <div v-else-if="deliveries.length === 0" class="deliveries-empty">No deliveries yet.</div>
        <div v-else class="deliveries-list">
          <div v-for="d in deliveries" :key="d.id" class="delivery-row">
            <span
              class="delivery-status-dot"
              :class="d.status === 'sent' ? 'delivery-status-dot--success' : 'delivery-status-dot--error'"
            >{{ d.status === 'sent' ? '✓' : '✗' }}</span>
            <div class="delivery-info">
              <div class="delivery-info__status">{{ d.status === 'sent' ? 'Sent' : 'Failed' }}</div>
              <div v-if="d.delivered_at" class="delivery-info__date">{{ formatDateTime(d.delivered_at) }}</div>
              <div v-if="d.failure_reason" class="delivery-info__error">{{ d.failure_reason }}</div>
            </div>
            <div class="delivery-created">{{ formatDateTime(d.created_at) }}</div>
          </div>
        </div>
      </div>

      <template #actions>
        <button class="btn-text" @click="showDeliveriesModal = false">Close</button>
      </template>
    </BaseModal>

    <!-- ── Preview Dialog ─────────────────────────────────────────────── -->
    <BaseModal v-model="showPreviewModal" size="lg">
      <template #headline>
        Report Preview
      </template>

      <div class="modal-preview">
        <div class="preview-report-name">{{ previewReport?.name }}</div>

        <!-- Loading skeleton -->
        <div v-if="previewLoading" class="preview-skeleton">
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
          <div class="preview-header-bar">
            <div class="preview-header-bar__label">Analytics Report</div>
            <div class="preview-header-bar__period">
              <span style="text-transform:capitalize;font-weight:500;">{{ previewReport?.frequency }}</span> summary — {{ previewPeriodLabel }}
            </div>
          </div>

          <!-- Stats -->
          <div class="preview-stats-grid">
            <div class="preview-stat-card">
              <div class="preview-stat-value">{{ previewData.total_clicks.toLocaleString() }}</div>
              <div class="preview-stat-label">Total Clicks</div>
            </div>
            <div class="preview-stat-card">
              <div class="preview-stat-value">{{ previewData.total_links.toLocaleString() }}</div>
              <div class="preview-stat-label">Active Links</div>
            </div>
            <div class="preview-stat-card">
              <div class="preview-stat-value">{{ previewData.clicks_7_days.toLocaleString() }}</div>
              <div class="preview-stat-label">Last 7 Days</div>
            </div>
          </div>

          <!-- Report links -->
          <div class="preview-links-section">
            <div class="preview-section-title">This report would include:</div>
            <div class="preview-link-chips">
              <span
                v-for="lid in previewReport?.link_ids"
                :key="lid"
                class="m3-badge m3-badge--neutral"
              >{{ lid.slice(0, 8) }}…</span>
            </div>
          </div>

          <!-- Top links -->
          <div v-if="previewData.top_links && previewData.top_links.length > 0" class="preview-top-links">
            <div class="preview-section-title">Top Performing Links (Account-wide)</div>
            <div>
              <div
                v-for="(link, idx) in previewData.top_links.slice(0, 5)"
                :key="link.id"
                class="top-link-row"
              >
                <span class="top-link-rank">{{ idx + 1 }}</span>
                <div class="top-link-info">
                  <div class="top-link-info__title">{{ link.title || link.slug }}</div>
                  <div class="top-link-info__url">{{ link.short_url }}</div>
                </div>
                <div class="top-link-clicks">
                  <div class="top-link-clicks__value">{{ (link.click_count || 0).toLocaleString() }}</div>
                  <div class="top-link-clicks__label">clicks</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Schedule note -->
          <div class="preview-schedule-note">
            <span class="preview-schedule-note__label">Email delivery schedule:</span>
            <span v-if="previewReport?.frequency === 'daily'"> This report is sent every day.</span>
            <span v-else-if="previewReport?.frequency === 'weekly'"> This report is sent every Monday.</span>
            <span v-else-if="previewReport?.frequency === 'monthly'"> This report is sent on the 1st of each month.</span>
            <span v-if="previewReport?.next_run_at"> Next delivery: <strong>{{ formatDate(previewReport.next_run_at) }}</strong>.</span>
          </div>
        </div>

        <!-- Error fallback -->
        <div v-else class="preview-error">Could not load preview data.</div>

        <div class="preview-footnote">Preview shows current account analytics data.</div>
      </div>

      <template #actions>
        <button class="btn-text" @click="showPreviewModal = false">Close</button>
        <button class="btn-filled" :disabled="sendingPreview" @click="sendNowFromPreview">
          <span v-if="sendingPreview" class="css-spinner css-spinner--sm css-spinner--white"></span>
          Send now
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, onUnmounted } from 'vue';
import BaseModal from '@/components/BaseModal.vue';
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
  return s.charAt(0).toUpperCase() + s.slice(1);
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
    const [p1, p2] = await Promise.all([
      linksApi.list(1, 100),
      linksApi.list(2, 100),
    ]);
    allLinks.value = [...(p1.data?.links ?? []), ...(p2.data?.links ?? [])];
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
/* ── Root ─────────────────────────────────────────────────────────────────── */
.reports-page {
  max-width: 900px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page header ──────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
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

/* ── Error banner ─────────────────────────────────────────────────────────── */
.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 12px;
  font-size: 0.875rem;
}

/* ── Loading ──────────────────────────────────────────────────────────────── */
.loading-state {
  display: flex;
  justify-content: center;
  padding: 48px;
}

/* ── CSS Spinner ──────────────────────────────────────────────────────────── */
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

/* ── AN Card ──────────────────────────────────────────────────────────────── */
.an-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
}

/* ── Empty state ──────────────────────────────────────────────────────────── */
.empty-state-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 24px;
  text-align: center;
}

.empty-icon {
  width: 72px;
  height: 72px;
  border-radius: 20px;
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

.empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.empty-desc {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0 0 20px;
}

/* ── Reports list ─────────────────────────────────────────────────────────── */
.reports-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ── Report card ──────────────────────────────────────────────────────────── */
.report-card {
  transition: box-shadow 0.15s;

  &:hover {
    box-shadow: 0 4px 16px rgba(0,0,0,0.10);
  }

  &__body {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
    padding: 16px 20px;
  }
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.report-info {
  flex: 1;
  min-width: 0;

  &__name {
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--md-sys-color-on-surface);
  }

  &__meta {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.8rem;
    margin-top: 4px;
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }
}

.report-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

/* ── Frequency badge ──────────────────────────────────────────────────────── */
.freq-badge {
  display: inline-block;
  background: rgba(99, 91, 255, 0.10);
  color: var(--md-sys-color-primary, #635bff);
  border-radius: 6px;
  padding: 0.1em 0.55em;
  font-size: 0.72rem;
  font-weight: 600;
}

/* ── M3 Badges ────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Error box ────────────────────────────────────────────────────────────── */
.error-box {
  padding: 10px 14px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  font-size: 0.875rem;
}

/* ── Form field ───────────────────────────────────────────────────────────── */
.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;

  &__label {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
  }
}

.form-input {
  height: 40px;
  padding: 0 12px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9375rem;
  font-family: inherit;
  transition: border-color 0.15s, box-shadow 0.15s;
  width: 100%;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }
}

/* ── Modal form ───────────────────────────────────────────────────────────── */
.modal-form {
  min-width: 480px;
  max-width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ── Frequency buttons ────────────────────────────────────────────────────── */
.freq-buttons {
  display: flex;
  gap: 8px;
}

/* ── Link selector ────────────────────────────────────────────────────────── */
.link-selector {
  position: relative;
}

.selector-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
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

.selector-placeholder {
  color: var(--md-sys-color-on-surface-variant);
}

.selector-chevron {
  font-size: 18px;
  flex-shrink: 0;
  color: var(--md-sys-color-on-surface-variant);
}

.selector-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 500;
  border-radius: 8px;
  overflow: hidden;
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
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

  &:hover { background: var(--md-sys-color-surface-container-low); }

  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &__title {
    font-size: 0.8125rem;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__url {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.7rem;
  }
}

.selector-empty {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  text-align: center;
  padding: 16px;
}

/* ── Deliveries modal ─────────────────────────────────────────────────────── */
.modal-deliveries {
  min-width: 480px;
  max-width: 100%;
  min-height: 200px;
}

.deliveries-empty {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  text-align: center;
  padding: 24px;
}

.deliveries-list {
  max-height: 360px;
  overflow-y: auto;
}

.delivery-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &:last-child { border-bottom: none; }
}

.delivery-status-dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
  flex-shrink: 0;

  &--success { background: #dcfce7; color: #16a34a; }
  &--error   { background: #fee2e2; color: #dc2626; }
}

.delivery-info {
  min-width: 0;
  flex: 1;

  &__status { font-size: 0.875rem; font-weight: 500; }
  &__date   { color: var(--md-sys-color-on-surface-variant); font-size: 0.72rem; }
  &__error  { color: var(--md-sys-color-error); font-size: 0.72rem; }
}

.delivery-created {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.72rem;
  white-space: nowrap;
  flex-shrink: 0;
}

/* ── Preview modal ────────────────────────────────────────────────────────── */
.modal-preview {
  min-width: 540px;
  max-width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.preview-report-name {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}

.preview-skeleton {
  padding: 16px 0;
}

.preview-header-bar {
  background: #f0efff;
  border-left: 4px solid var(--md-sys-color-primary);
  padding: 12px 16px;
  border-radius: 0 8px 8px 0;

  &__label {
    font-weight: 600;
    color: var(--md-sys-color-primary);
  }

  &__period {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.875rem;
    margin-top: 4px;
  }
}

.preview-stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.preview-stat-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  padding: 16px;
  text-align: center;
}

.preview-stat-value {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--md-sys-color-primary);
}

.preview-stat-label {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.8rem;
  margin-top: 4px;
}

.preview-links-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preview-link-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.preview-section-title {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.preview-top-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.top-link-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &:last-child { border-bottom: none; }
}

.top-link-rank {
  width: 20px;
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
}

.top-link-info {
  flex: 1;
  min-width: 0;

  &__title {
    font-size: 0.875rem;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__url {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.7rem;
  }
}

.top-link-clicks {
  text-align: right;
  flex-shrink: 0;

  &__value {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--md-sys-color-primary);
  }

  &__label {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.68rem;
  }
}

.preview-schedule-note {
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);

  &__label {
    font-weight: 500;
    color: var(--md-sys-color-on-surface);
  }
}

.preview-error {
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  padding: 32px;
  font-size: 0.875rem;
}

.preview-footnote {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.8rem;
}

/* ── Skeleton loading ─────────────────────────────────────────────────────── */
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
