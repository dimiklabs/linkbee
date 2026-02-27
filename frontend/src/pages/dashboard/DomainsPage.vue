<template>
  <div class="container-fluid py-4">

    <!-- Page header -->
    <div class="d-flex align-items-center justify-content-between mb-4 flex-wrap gap-2">
      <div>
        <h1 class="h4 fw-bold mb-0">Custom Domains</h1>
        <p class="text-muted small mb-0">Use your own domain to serve short links.</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        + Add Domain
      </button>
    </div>

    <!-- Alert -->
    <div v-if="alertMsg" :class="['alert', alertType === 'error' ? 'alert-danger' : 'alert-success', 'alert-dismissible']" role="alert">
      {{ alertMsg }}
      <button type="button" class="btn-close" @click="alertMsg = ''"></button>
    </div>

    <!-- DNS instructions banner -->
    <div class="alert alert-info mb-4 small">
      <strong>How it works:</strong> Add your domain, then create a DNS TXT record
      <code>_shortlink-verify.&lt;yourdomain.com&gt;</code> with the value shown below,
      then click <em>Verify</em>. Once verified, point your domain's CNAME to
      <code>{{ appDomain }}</code> and your short links will be served under your brand.
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading…</span>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="domains.length === 0" class="card text-center py-5">
      <div class="card-body">
        <p class="fs-1 mb-2">🌐</p>
        <h5 class="mb-1">No custom domains yet</h5>
        <p class="text-muted mb-3">Add your first domain to start using branded short links.</p>
        <button class="btn btn-primary" @click="showAddModal = true">Add Domain</button>
      </div>
    </div>

    <!-- Domains table -->
    <div v-else class="card">
      <div class="table-responsive">
        <table class="table table-hover mb-0 align-middle">
          <thead>
            <tr>
              <th>Domain</th>
              <th>Status</th>
              <th>Verify Token</th>
              <th>Added</th>
              <th class="text-end">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in domains" :key="d.id">
              <td class="fw-medium">{{ d.domain }}</td>
              <td>
                <span :class="['badge', statusClass(d.status)]">{{ d.status }}</span>
              </td>
              <td>
                <code class="small text-break">{{ d.verify_token }}</code>
                <button
                  class="btn btn-sm btn-link p-0 ms-1 align-middle"
                  title="Copy token"
                  @click="copyToken(d.verify_token)"
                >📋</button>
              </td>
              <td class="text-muted small">{{ formatDate(d.created_at) }}</td>
              <td class="text-end">
                <button
                  v-if="d.status !== 'verified'"
                  class="btn btn-sm btn-outline-primary me-1"
                  :disabled="verifying === d.id"
                  @click="verify(d)"
                >
                  <span v-if="verifying === d.id" class="spinner-border spinner-border-sm me-1"></span>
                  Verify
                </button>
                <button
                  class="btn btn-sm btn-outline-danger"
                  :disabled="deleting === d.id"
                  @click="confirmDelete(d)"
                >
                  <span v-if="deleting === d.id" class="spinner-border spinner-border-sm me-1"></span>
                  Remove
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Domain Modal -->
    <div v-if="showAddModal" class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,.5);">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Custom Domain</h5>
            <button type="button" class="btn-close" @click="closeAddModal"></button>
          </div>
          <form @submit.prevent="addDomain">
            <div class="modal-body">
              <div v-if="addError" class="alert alert-danger small py-2">{{ addError }}</div>
              <label class="form-label" for="domainInput">Domain</label>
              <input
                id="domainInput"
                v-model="newDomain"
                type="text"
                class="form-control"
                placeholder="links.yourdomain.com"
                required
                autocomplete="off"
              />
              <div class="form-text">Enter the domain or subdomain you want to use (e.g. <code>go.acme.com</code>).</div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" @click="closeAddModal">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="adding">
                <span v-if="adding" class="spinner-border spinner-border-sm me-1"></span>
                Add Domain
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Confirm Delete Modal -->
    <div v-if="deleteTarget" class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,.5);">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Remove Domain</h5>
            <button type="button" class="btn-close" @click="deleteTarget = null"></button>
          </div>
          <div class="modal-body">
            <p>Remove <strong>{{ deleteTarget.domain }}</strong>? This cannot be undone.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="deleteTarget = null">Cancel</button>
            <button type="button" class="btn btn-danger" :disabled="!!deleting" @click="deleteDomain">
              <span v-if="deleting" class="spinner-border spinner-border-sm me-1"></span>
              Remove
            </button>
          </div>
        </div>
      </div>
    </div>

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
  return {
    'text-bg-success': status === 'verified',
    'text-bg-warning': status === 'pending',
    'text-bg-danger': status === 'failed',
  };
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
