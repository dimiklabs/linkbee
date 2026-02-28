<template>
  <div class="domains-page">

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">Custom Domains</h1>
        <p class="page-subtitle">Use your own domain to serve short links.</p>
      </div>
      <button class="btn-filled" @click="showAddModal = true">
        <span class="material-symbols-outlined">add</span>
        Add Domain
      </button>
    </div>

    <!-- Snackbar / Alert -->
    <Transition name="snack">
      <div v-if="alertMsg" :class="['m3-snackbar', alertType === 'error' ? 'm3-snackbar--error' : '']">
        <span class="material-symbols-outlined snack-icon">{{ alertType === 'error' ? 'error' : 'check_circle' }}</span>
        <span class="snack-text">{{ alertMsg }}</span>
        <button class="btn-text" @click="alertMsg = ''">Dismiss</button>
      </div>
    </Transition>

    <!-- DNS instructions info card -->
    <div class="info-card">
      <span class="material-symbols-outlined info-card__icon">info</span>
      <div class="info-card__text">
        <strong>How it works:</strong> Add your domain, then create a DNS TXT record
        <code>_shortlink-verify.&lt;yourdomain.com&gt;</code> with the value shown below,
        then click <em>Verify</em>. Once verified, point your domain's CNAME to
        <code>{{ appDomain }}</code> and your short links will be served under your brand.
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-center">
      <div class="css-spinner"></div>
    </div>

    <!-- Empty state -->
    <div v-else-if="domains.length === 0" class="an-card empty-state">
      <div class="empty-icon">
        <span class="material-symbols-outlined">language</span>
      </div>
      <div class="empty-title">No custom domains yet</div>
      <p class="empty-desc">Add your first domain to start using branded short links.</p>
      <button class="btn-filled" @click="showAddModal = true">
        <span class="material-symbols-outlined">add</span>
        Add Domain
      </button>
    </div>

    <!-- Domains table -->
    <div v-else class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">language</span>
          </div>
          <span class="an-card-title">Custom Domains</span>
        </div>
        <span class="m3-badge m3-badge--neutral">{{ domains.length }} domain{{ domains.length !== 1 ? 's' : '' }}</span>
      </div>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>Domain</th>
              <th>Status</th>
              <th>Verify Token</th>
              <th>Added</th>
              <th class="th-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in domains" :key="d.id">
              <td class="cell-strong">{{ d.domain }}</td>
              <td>
                <span :class="['m3-badge', statusClass(d.status)]">{{ d.status }}</span>
              </td>
              <td>
                <div class="copy-field">
                  <code class="copy-field__value">{{ d.verify_token }}</code>
                  <button class="btn-icon" title="Copy verification token" @click="copyToken(d.verify_token)">
                    <span class="material-symbols-outlined copy-icon">content_copy</span>
                  </button>
                </div>
              </td>
              <td class="cell-muted cell-sm cell-nowrap">{{ formatDate(d.created_at) }}</td>
              <td class="td-right">
                <div class="row-actions">
                  <button class="btn-outlined"
                    v-if="d.status !== 'verified'"
                    :disabled="verifying === d.id"
                    @click="verify(d)"
                  >
                    <div v-if="verifying === d.id" class="css-spinner css-spinner--sm"></div>
                    <span v-else class="material-symbols-outlined">verified</span>
                    Verify
                  </button>
                  <button class="btn-outlined btn-danger"
                    :disabled="deleting === d.id"
                    @click="confirmDelete(d)"
                  >
                    <div v-if="deleting === d.id" class="css-spinner css-spinner--sm"></div>
                    <span v-else class="material-symbols-outlined">delete</span>
                    Remove
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Domain Dialog -->
    <BaseModal v-model="showAddModal" size="md" @closed="closeAddModal">
      <template #headline>Add Custom Domain</template>

      <div class="modal-body">
        <div v-if="addError" class="feedback-error">
          {{ addError }}
        </div>
        <label class="form-field">
          <span class="form-field__label">Domain</span>
          <input
            type="text"
            class="form-input"
            :value="newDomain"
            @input="newDomain=($event.target as HTMLInputElement).value"
            placeholder="links.yourdomain.com"
            @keydown.enter="addDomain"
          />
          <span class="form-field__hint">Enter the domain or subdomain you want to use (e.g. <code>go.acme.com</code>).</span>
        </label>
      </div>

      <template #actions>
        <button class="btn-text" @click="closeAddModal">Cancel</button>
        <button class="btn-filled" :disabled="adding" @click="addDomain">
          <div v-if="adding" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Add Domain
        </button>
      </template>
    </BaseModal>

    <!-- Confirm Delete Dialog -->
    <BaseModal v-model="showDeleteModal" size="sm" @closed="deleteTarget = null">
      <template #headline>Remove Domain</template>

      <div class="modal-body">
        <p class="modal-text">Remove <strong>{{ deleteTarget?.domain }}</strong>? This cannot be undone.</p>
      </div>

      <template #actions>
        <button class="btn-text" @click="showDeleteModal = false; deleteTarget = null">Cancel</button>
        <button class="btn-filled btn-danger" :disabled="!!deleting" @click="deleteDomain">
          <div v-if="deleting" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Remove
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import domainsApi from '@/api/domains';
import type { DomainResponse } from '@/types/domains';
import BaseModal from '@/components/BaseModal.vue';

const domains = ref<DomainResponse[]>([]);
const loading = ref(true);
const alertMsg = ref('');
const alertType = ref<'success' | 'error'>('success');

const showAddModal = ref(false);
const newDomain = ref('');
const adding = ref(false);
const addError = ref('');

const verifying = ref<string | null>(null);
const deleting = ref<string | null>(null);
const deleteTarget = ref<DomainResponse | null>(null);
const showDeleteModal = ref(false);

// Show the app's main domain so users know where to point their CNAME
const appDomain = window.location.hostname;

async function fetchDomains() {
  loading.value = true;
  try {
    const res = await domainsApi.list();
    domains.value = res.data.data ?? [];
  } catch {
    showAlert('Failed to load domains.', 'error');
  } finally {
    loading.value = false;
  }
}

async function addDomain() {
  addError.value = '';
  adding.value = true;
  try {
    const res = await domainsApi.add(newDomain.value.trim());
    domains.value.unshift(res.data.data!);
    closeAddModal();
    showAlert('Domain added. Add the DNS TXT record and click Verify.', 'success');
  } catch (err: unknown) {
    const msg = (err as { response?: { data?: { description?: string } } })?.response?.data?.description;
    addError.value = msg || 'Failed to add domain.';
  } finally {
    adding.value = false;
  }
}

async function verify(d: DomainResponse) {
  verifying.value = d.id;
  try {
    const res = await domainsApi.verify(d.id);
    const updated = res.data.data!;
    const idx = domains.value.findIndex((x) => x.id === d.id);
    if (idx !== -1) domains.value[idx] = updated;
    if (updated.status === 'verified') {
      showAlert(`${d.domain} verified successfully!`, 'success');
    } else {
      showAlert(`Verification failed — make sure the TXT record exists and try again.`, 'error');
    }
  } catch {
    showAlert('Verification request failed.', 'error');
  } finally {
    verifying.value = null;
  }
}

function confirmDelete(d: DomainResponse) {
  deleteTarget.value = d;
  showDeleteModal.value = true;
}

async function deleteDomain() {
  if (!deleteTarget.value) return;
  const d = deleteTarget.value;
  deleting.value = d.id;
  try {
    await domainsApi.remove(d.id);
    domains.value = domains.value.filter((x) => x.id !== d.id);
    deleteTarget.value = null;
    showDeleteModal.value = false;
    showAlert('Domain removed.', 'success');
  } catch {
    showAlert('Failed to remove domain.', 'error');
  } finally {
    deleting.value = null;
  }
}

function closeAddModal() {
  showAddModal.value = false;
  newDomain.value = '';
  addError.value = '';
}

function statusClass(status: string) {
  if (status === 'verified') return 'm3-badge--success';
  if (status === 'pending') return 'm3-badge--warning';
  if (status === 'failed') return 'm3-badge--error';
  return 'm3-badge--neutral';
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
}

async function copyToken(token: string) {
  try {
    await navigator.clipboard.writeText(token);
    showAlert('Token copied to clipboard.', 'success');
  } catch {
    /* ignore */
  }
}

function showAlert(msg: string, type: 'success' | 'error') {
  alertMsg.value = msg;
  alertType.value = type;
  setTimeout(() => { alertMsg.value = ''; }, 4000);
}

onMounted(fetchDomains);
</script>

<style scoped lang="scss">
.domains-page {
  max-width: 1000px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page Header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-header__left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  margin: 0;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── CSS Spinner ─────────────────────────────────────────────────────────── */
.css-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--md-sys-color-outline-variant);
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

/* ── Loading ─────────────────────────────────────────────────────────────── */
.loading-center {
  display: flex;
  justify-content: center;
  padding: 48px;
}

/* ── Cards ───────────────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);

  &__left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
  }
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
    background: color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
    color: var(--md-sys-color-primary);
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

/* ── Info card ───────────────────────────────────────────────────────────── */
.info-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;

  &__icon {
    color: var(--md-sys-color-primary);
    flex-shrink: 0;
    font-size: 20px;
    margin-top: 1px;
  }

  &__text {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface);
    line-height: 1.5;

    code {
      font-family: 'Courier New', monospace;
      font-size: 0.8em;
      background: var(--md-sys-color-surface-container-low);
      padding: 1px 4px;
      border-radius: 4px;
    }
  }
}

/* ── Empty state ─────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 24px;
  text-align: center;
  gap: 0;
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
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--md-sys-color-on-surface);
}

.empty-desc {
  margin: 0 0 20px;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Table ───────────────────────────────────────────────────────────────── */
.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;

  thead tr {
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  th {
    padding: 12px 16px;
    text-align: left;
    font-weight: 600;
    font-size: 0.8rem;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    white-space: nowrap;
  }

  td {
    padding: 12px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
    vertical-align: middle;
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  tbody tr:hover td {
    background: var(--md-sys-color-surface-container-low);
  }
}

.th-right { text-align: right; }
.td-right { text-align: right; }
.cell-strong { font-weight: 500; }
.cell-muted { color: var(--md-sys-color-on-surface-variant); }
.cell-sm { font-size: 0.8rem; }
.cell-nowrap { white-space: nowrap; }

.row-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}

/* ── Badges ──────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  text-transform: capitalize;

  &--success { background: rgba(22, 163, 74, 0.12); color: #16a34a; }
  &--warning { background: rgba(245, 158, 11, 0.12); color: #b45309; }
  &--error   { background: rgba(220, 38, 38, 0.12); color: #dc2626; }
  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Copy field ──────────────────────────────────────────────────────────── */
.copy-field {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px 4px 10px;
  border-radius: 8px;
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  max-width: 240px;
  overflow: hidden;

  &__value {
    font-family: 'Courier New', monospace;
    font-size: 0.72rem;
    color: var(--md-sys-color-on-surface);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    min-width: 0;
  }
}

.copy-icon {
  font-size: 16px !important;
}

/* ── Modal ───────────────────────────────────────────────────────────────── */
.modal-body {
  min-width: 400px;
  max-width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-text {
  margin: 0;
  font-size: 0.9375rem;
  color: var(--md-sys-color-on-surface);
}

/* ── Form field ──────────────────────────────────────────────────────────── */
.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field__label {
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
}

.form-field__hint {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);

  code {
    font-family: 'Courier New', monospace;
    font-size: 0.85em;
    background: var(--md-sys-color-surface-container-low);
    padding: 1px 4px;
    border-radius: 3px;
  }
}

.form-input {
  height: 40px;
  padding: 0 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9rem;
  outline: none;
  width: 100%;
  box-sizing: border-box;

  &:focus {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
  }
}

/* ── Feedback error ──────────────────────────────────────────────────────── */
.feedback-error {
  padding: 10px 14px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  font-size: 0.875rem;
}

/* ── Danger button ───────────────────────────────────────────────────────── */
.btn-danger {
  --btn-danger: var(--md-sys-color-error, #dc2626);
  border-color: var(--btn-danger) !important;
  color: var(--btn-danger) !important;
}

/* ── Snackbar ────────────────────────────────────────────────────────────── */
.m3-snackbar {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  background: #313033;
  color: #fff;
  border-radius: 4px;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 280px;
  max-width: 560px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0,0,0,0.24);

  &--error { background: var(--md-sys-color-error-container); color: var(--md-sys-color-on-error-container); }
}

.snack-icon { font-size: 20px; flex-shrink: 0; }
.snack-text { flex: 1; font-size: 0.875rem; }

.snack-enter-active, .snack-leave-active { transition: all .25s; }
.snack-enter-from, .snack-leave-to { transform: translate(-50%, 80px); opacity: 0; }
</style>
