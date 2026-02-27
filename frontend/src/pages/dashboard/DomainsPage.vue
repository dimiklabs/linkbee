<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Custom Domains</h1>
        <p class="page-subtitle">Use your own domain to serve short links.</p>
      </div>
      <md-filled-button @click="showAddModal = true">
        <span class="material-symbols-outlined" slot="icon">add</span>
        Add Domain
      </md-filled-button>
    </div>

    <!-- Snackbar / Alert -->
    <Transition name="snack">
      <div v-if="alertMsg" class="m3-snackbar">
        <span class="material-symbols-outlined" style="font-size:20px;">{{ alertType === 'error' ? 'error' : 'check_circle' }}</span>
        <span style="flex:1">{{ alertMsg }}</span>
        <md-text-button @click="alertMsg = ''" style="--md-text-button-label-text-color:#CFBCFF">Dismiss</md-text-button>
      </div>
    </Transition>

    <!-- DNS instructions info card -->
    <div class="m3-card m3-card--outlined info-card">
      <span class="material-symbols-outlined" style="color:var(--md-sys-color-primary);flex-shrink:0;">info</span>
      <div style="font-size:0.875rem;color:var(--md-sys-color-on-surface);">
        <strong>How it works:</strong> Add your domain, then create a DNS TXT record
        <code>_shortlink-verify.&lt;yourdomain.com&gt;</code> with the value shown below,
        then click <em>Verify</em>. Once verified, point your domain's CNAME to
        <code>{{ appDomain }}</code> and your short links will be served under your brand.
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:40px" />
    </div>

    <!-- Empty state -->
    <div v-else-if="domains.length === 0" class="m3-card m3-card--elevated m3-empty-state">
      <div class="m3-empty-state__icon">
        <span class="material-symbols-outlined">language</span>
      </div>
      <div class="md-title-medium" style="margin-bottom:8px;">No custom domains yet</div>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0 0 20px;">Add your first domain to start using branded short links.</p>
      <md-filled-button @click="showAddModal = true">
        <span class="material-symbols-outlined" slot="icon">add</span>
        Add Domain
      </md-filled-button>
    </div>

    <!-- Domains table -->
    <div v-else class="m3-card m3-card--elevated">
      <div class="m3-card-header">
        <div class="m3-card-header__left">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">language</span>
          <span class="md-title-medium">Custom Domains</span>
        </div>
        <span class="m3-badge m3-badge--neutral">{{ domains.length }} domain{{ domains.length !== 1 ? 's' : '' }}</span>
      </div>
      <md-divider />
      <div class="m3-table-wrapper">
        <table class="m3-table">
          <thead>
            <tr>
              <th>Domain</th>
              <th>Status</th>
              <th>Verify Token</th>
              <th>Added</th>
              <th style="text-align:right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in domains" :key="d.id">
              <td style="font-weight:500;">{{ d.domain }}</td>
              <td>
                <span :class="['m3-badge', statusClass(d.status)]">{{ d.status }}</span>
              </td>
              <td>
                <div style="display:flex;align-items:center;gap:6px;">
                  <code style="font-size:0.75rem;word-break:break-all;max-width:200px;">{{ d.verify_token }}</code>
                  <md-icon-button title="Copy token" @click="copyToken(d.verify_token)" style="--md-icon-button-state-layer-size:32px;">
                    <span class="material-symbols-outlined" style="font-size:18px;">content_copy</span>
                  </md-icon-button>
                </div>
              </td>
              <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;white-space:nowrap;">{{ formatDate(d.created_at) }}</td>
              <td style="text-align:right;">
                <div style="display:flex;justify-content:flex-end;gap:8px;align-items:center;">
                  <md-outlined-button
                    v-if="d.status !== 'verified'"
                    :disabled="verifying === d.id"
                    @click="verify(d)"
                  >
                    <md-circular-progress v-if="verifying === d.id" indeterminate style="--md-circular-progress-size:16px;margin-right:6px;" />
                    <span v-else class="material-symbols-outlined" slot="icon">verified</span>
                    Verify
                  </md-outlined-button>
                  <md-outlined-button
                    :disabled="deleting === d.id"
                    @click="confirmDelete(d)"
                    style="--md-outlined-button-label-text-color:var(--md-sys-color-error);--md-outlined-button-outline-color:var(--md-sys-color-error);"
                  >
                    <md-circular-progress v-if="deleting === d.id" indeterminate style="--md-circular-progress-size:16px;margin-right:6px;" />
                    <span v-else class="material-symbols-outlined" slot="icon">delete</span>
                    Remove
                  </md-outlined-button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Domain Dialog -->
    <md-dialog :open="showAddModal" @closed="closeAddModal">
      <div slot="headline">Add Custom Domain</div>
      <div slot="content" style="min-width:480px;max-width:100%;">
        <div v-if="addError" style="margin-bottom:16px;padding:10px 14px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:8px;font-size:0.875rem;">
          {{ addError }}
        </div>
        <md-outlined-text-field
          :value="newDomain"
          @input="newDomain=($event.target as HTMLInputElement).value"
          label="Domain"
          placeholder="links.yourdomain.com"
          style="width:100%;"
          @keydown.enter="addDomain"
        />
        <div style="font-size:0.75rem;color:var(--md-sys-color-on-surface-variant);margin-top:8px;">
          Enter the domain or subdomain you want to use (e.g. <code>go.acme.com</code>).
        </div>
      </div>
      <div slot="actions">
        <md-text-button @click="closeAddModal">Cancel</md-text-button>
        <md-filled-button :disabled="adding" @click="addDomain">
          <md-circular-progress v-if="adding" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
          Add Domain
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- Confirm Delete Dialog -->
    <md-dialog :open="!!deleteTarget" @closed="deleteTarget = null">
      <div slot="headline">Remove Domain</div>
      <div slot="content" style="min-width:400px;max-width:100%;">
        <p style="margin:0;font-size:0.9375rem;">Remove <strong>{{ deleteTarget?.domain }}</strong>? This cannot be undone.</p>
      </div>
      <div slot="actions">
        <md-text-button @click="deleteTarget = null">Cancel</md-text-button>
        <md-filled-button :disabled="!!deleting" @click="deleteDomain" style="--md-filled-button-container-color:var(--md-sys-color-error);--md-filled-button-label-text-color:var(--md-sys-color-on-error);">
          <md-circular-progress v-if="deleting" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
          Remove
        </md-filled-button>
      </div>
    </md-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import domainsApi from '@/api/domains';
import type { DomainResponse } from '@/types/domains';

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
}

async function deleteDomain() {
  if (!deleteTarget.value) return;
  const d = deleteTarget.value;
  deleting.value = d.id;
  try {
    await domainsApi.remove(d.id);
    domains.value = domains.value.filter((x) => x.id !== d.id);
    deleteTarget.value = null;
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
.page-wrapper {
  padding: 24px;
  max-width: 1000px;
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

.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
  margin-bottom: 20px;

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
  }

  &--outlined {
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

.m3-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  gap: 1rem;
  flex-wrap: wrap;

  &__left {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.info-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
}

.m3-table-wrapper {
  overflow-x: auto;
}

.m3-table {
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

.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
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

/* Snackbar */
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
}

.snack-enter-active, .snack-leave-active { transition: all .25s; }
.snack-enter-from, .snack-leave-to { transform: translate(-50%, 80px); opacity: 0; }
</style>
