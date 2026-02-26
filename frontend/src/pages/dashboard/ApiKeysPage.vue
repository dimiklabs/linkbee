<template>
  <div class="container-fluid py-4" style="max-width: 860px;">

    <!-- Page header -->
    <div class="d-flex align-items-center justify-content-between mb-4">
      <div>
        <h4 class="mb-1 fw-semibold">API Keys</h4>
        <p class="text-muted mb-0 small">
          Use API keys to authenticate programmatic requests with
          <code class="text-primary">X-API-Key: &lt;key&gt;</code> header.
        </p>
      </div>
      <button class="btn btn-primary btn-sm" @click="showCreate = true">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" class="me-1" viewBox="0 0 16 16">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
        </svg>
        New API Key
      </button>
    </div>

    <!-- One-time key reveal banner -->
    <div v-if="newKey" class="alert border-0 mb-4 rounded-3 new-key-banner" role="alert">
      <div class="d-flex align-items-start gap-3">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" class="flex-shrink-0 mt-1 text-warning" viewBox="0 0 16 16">
          <path d="M8 1a2 2 0 0 1 2 2v4H6V3a2 2 0 0 1 2-2m3 6V3a3 3 0 0 0-6 0v4a2 2 0 0 0-2 2v5a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2"/>
        </svg>
        <div class="flex-grow-1">
          <div class="fw-semibold mb-1">Save your API key — it won't be shown again</div>
          <div class="d-flex align-items-center gap-2 flex-wrap">
            <code class="key-display px-3 py-2 rounded flex-grow-1" style="word-break:break-all;">{{ newKey.key }}</code>
            <button class="btn btn-sm btn-warning flex-shrink-0" @click="copyKey(newKey.key)" :title="copied ? 'Copied!' : 'Copy'">
              <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1z"/>
                <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0z"/>
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                <path d="M10.97 4.97a.75.75 0 0 1 1.07 1.05l-3.99 4.99a.75.75 0 0 1-1.08.02L4.324 8.384a.75.75 0 1 1 1.06-1.06l2.094 2.093 3.473-4.425z"/>
              </svg>
              {{ copied ? 'Copied' : 'Copy' }}
            </button>
          </div>
        </div>
        <button type="button" class="btn-close" @click="newKey = null"></button>
      </div>
    </div>

    <!-- Create form (inline) -->
    <div v-if="showCreate" class="card border-0 shadow-sm mb-4">
      <div class="card-body">
        <h6 class="mb-3 fw-semibold">Create new API key</h6>
        <div class="row g-3 align-items-end">
          <div class="col-md-5">
            <label class="form-label small fw-semibold">Name <span class="text-danger">*</span></label>
            <input
              ref="nameInputRef"
              v-model="form.name"
              type="text"
              class="form-control form-control-sm"
              placeholder="e.g. My integration, CI/CD pipeline"
              maxlength="100"
              @keydown.enter="createKey"
            />
          </div>
          <div class="col-md-4">
            <label class="form-label small fw-semibold">
              Expiry
              <span class="text-muted fw-normal">(optional)</span>
            </label>
            <input
              v-model="form.expires_at"
              type="date"
              class="form-control form-control-sm"
              :min="tomorrow"
            />
          </div>
          <div class="col-md-3 d-flex gap-2">
            <button
              class="btn btn-primary btn-sm flex-fill"
              @click="createKey"
              :disabled="!form.name.trim() || creating"
            >
              <span v-if="creating" class="spinner-border spinner-border-sm me-1"></span>
              Create
            </button>
            <button class="btn btn-outline-secondary btn-sm" @click="cancelCreate">Cancel</button>
          </div>
        </div>
        <div v-if="createError" class="text-danger small mt-2">{{ createError }}</div>
      </div>
    </div>

    <!-- Keys table -->
    <div class="card border-0 shadow-sm">
      <div class="card-body p-0">

        <div v-if="loading" class="text-center py-5 text-muted">
          <div class="spinner-border spinner-border-sm me-2"></div>Loading API keys…
        </div>

        <div v-else-if="keys.length === 0 && !loading" class="text-center py-5">
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" fill="currentColor" class="text-muted mb-3 d-block mx-auto opacity-50" viewBox="0 0 16 16">
            <path d="M0 8a4 4 0 0 1 7.465-2H14a.5.5 0 0 1 .354.146l1.5 1.5a.5.5 0 0 1 0 .708l-1.5 1.5a.5.5 0 0 1-.708 0L13 9.207l-.646.647a.5.5 0 0 1-.708 0L11 9.207l-.646.647a.5.5 0 0 1-.354.146H7.465A4 4 0 0 1 0 8m4-3a3 3 0 1 0 0 6 3 3 0 0 0 0-6m0 1a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
          </svg>
          <div class="text-muted">No API keys yet.</div>
          <button class="btn btn-sm btn-outline-primary mt-2" @click="showCreate = true">Create your first key</button>
        </div>

        <table v-else class="table table-hover mb-0 align-middle">
          <thead class="table-light">
            <tr>
              <th class="ps-4 py-3 small text-muted fw-semibold">Name</th>
              <th class="py-3 small text-muted fw-semibold">Prefix</th>
              <th class="py-3 small text-muted fw-semibold">Last used</th>
              <th class="py-3 small text-muted fw-semibold">Expires</th>
              <th class="py-3 small text-muted fw-semibold">Created</th>
              <th class="pe-4 py-3 small text-muted fw-semibold text-end">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="key in keys" :key="key.id">
              <td class="ps-4 py-3 fw-semibold">{{ key.name }}</td>
              <td class="py-3">
                <code class="bg-light px-2 py-1 rounded small text-muted">{{ key.key_prefix }}…</code>
              </td>
              <td class="py-3 text-muted small">
                <span v-if="key.last_used_at">{{ formatDate(key.last_used_at) }}</span>
                <span v-else class="text-muted">Never</span>
              </td>
              <td class="py-3 small">
                <span v-if="key.expires_at" :class="isExpired(key.expires_at) ? 'text-danger' : 'text-muted'">
                  {{ isExpired(key.expires_at) ? 'Expired ' : '' }}{{ formatDate(key.expires_at) }}
                </span>
                <span v-else class="text-muted">No expiry</span>
              </td>
              <td class="py-3 text-muted small">{{ formatDate(key.created_at) }}</td>
              <td class="pe-4 py-3 text-end">
                <button
                  class="btn btn-sm btn-outline-danger border-0"
                  title="Revoke"
                  :disabled="revokingId === key.id"
                  @click="revokeKey(key)"
                >
                  <span v-if="revokingId === key.id" class="spinner-border spinner-border-sm"></span>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" width="15" height="15" fill="currentColor" viewBox="0 0 16 16">
                    <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                    <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                  </svg>
                  Revoke
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Usage docs -->
    <div class="card border-0 shadow-sm mt-4">
      <div class="card-body">
        <h6 class="fw-semibold mb-3">How to use</h6>
        <p class="text-muted small mb-2">Include the key in the <code>X-API-Key</code> header of every request:</p>
        <pre class="bg-light rounded p-3 small mb-2">curl https://your-domain.com/api/v1/links \
  -H "X-API-Key: sl_your_api_key_here"</pre>
        <p class="text-muted small mb-0">API keys have the same access as your account. Revoke keys you no longer need.</p>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import apiKeysApi from '@/api/apikeys';
import type { APIKey, CreateAPIKeyResponse } from '@/types/apikeys';

// ── State ──────────────────────────────────────────────────────────────────────
const keys = ref<APIKey[]>([]);
const loading = ref(false);
const showCreate = ref(false);
const creating = ref(false);
const createError = ref('');
const revokingId = ref<string | null>(null);
const newKey = ref<CreateAPIKeyResponse | null>(null);
const copied = ref(false);

const form = ref({ name: '', expires_at: '' });
const nameInputRef = ref<HTMLInputElement | null>(null);

// ── Helpers ────────────────────────────────────────────────────────────────────
const tomorrow = computed(() => {
  const d = new Date();
  d.setDate(d.getDate() + 1);
  return d.toISOString().split('T')[0];
});

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
}

function isExpired(iso: string): boolean {
  return new Date(iso) < new Date();
}

// ── Load ───────────────────────────────────────────────────────────────────────
async function loadKeys() {
  loading.value = true;
  try {
    const res = await apiKeysApi.list();
    keys.value = res.data ?? [];
  } catch {
    keys.value = [];
  } finally {
    loading.value = false;
  }
}

// ── Create ─────────────────────────────────────────────────────────────────────
async function createKey() {
  if (!form.value.name.trim()) return;
  createError.value = '';
  creating.value = true;
  try {
    const payload: { name: string; expires_at?: string } = { name: form.value.name.trim() };
    if (form.value.expires_at) {
      // Convert date-only → RFC3339
      payload.expires_at = new Date(form.value.expires_at + 'T23:59:59Z').toISOString();
    }
    const res = await apiKeysApi.create(payload);
    newKey.value = res.data;
    copied.value = false;
    keys.value.unshift(res.data);
    cancelCreate();
  } catch (err: any) {
    createError.value = err?.response?.data?.description ?? 'Failed to create API key';
  } finally {
    creating.value = false;
  }
}

function cancelCreate() {
  showCreate.value = false;
  form.value = { name: '', expires_at: '' };
  createError.value = '';
}

// ── Copy ───────────────────────────────────────────────────────────────────────
async function copyKey(key: string) {
  try {
    await navigator.clipboard.writeText(key);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 2000);
  } catch {
    // fallback
  }
}

// ── Revoke ─────────────────────────────────────────────────────────────────────
async function revokeKey(key: APIKey) {
  if (!confirm(`Revoke key "${key.name}"? Any integrations using it will stop working.`)) return;
  revokingId.value = key.id;
  try {
    await apiKeysApi.revoke(key.id);
    keys.value = keys.value.filter(k => k.id !== key.id);
    if (newKey.value?.id === key.id) newKey.value = null;
  } catch {
    // ignore
  } finally {
    revokingId.value = null;
  }
}

// ── Focus name input when form opens ──────────────────────────────────────────
import { watch } from 'vue';
watch(showCreate, async (v) => {
  if (v) {
    await nextTick();
    nameInputRef.value?.focus();
  }
});

onMounted(loadKeys);
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
.btn-outline-primary {
  color: #635bff;
  border-color: #635bff;
}
.btn-outline-primary:hover {
  background-color: #635bff;
  border-color: #635bff;
}

.new-key-banner {
  background-color: #fffbeb;
  border: 1px solid #fde68a !important;
}

.key-display {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 0.78rem;
  background: #f8f9fa;
  color: #212529;
  display: inline-block;
}
</style>
