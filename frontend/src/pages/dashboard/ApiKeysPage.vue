<template>
  <UpgradeWall
    v-if="!billingStore.isPaidPlan && billingStore.loaded"
    icon="key"
    title="API Keys is a Pro feature"
    description="Upgrade to Pro to create API keys for programmatic access to your links."
  />
  <div v-else class="apikeys-page">

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">API Keys</h1>
        <p class="page-subtitle">
          Use API keys to authenticate programmatic requests with
          <code class="inline-code">X-API-Key: &lt;key&gt;</code> header.
        </p>
      </div>
      <div class="page-header__actions">
        <button class="btn-filled" @click="showCreate = true">
          <span class="material-symbols-outlined">add</span>
          New API Key
        </button>
      </div>
    </div>

    <!-- Usage Warning Banner -->
    <div v-if="usageWarning" class="alert-banner" :class="usageWarning.level === 'danger' ? 'alert-banner--error' : 'alert-banner--warning'">
      <span class="material-symbols-outlined alert-banner__icon">{{ usageWarning.level === 'danger' ? 'block' : 'warning' }}</span>
      <span class="alert-banner__msg">{{ usageWarning.msg }}</span>
      <RouterLink to="/dashboard/billing">
        <button class="btn-tonal">Upgrade</button>
      </RouterLink>
    </div>

    <!-- One-time key reveal banner -->
    <div v-if="newKey" class="reveal-banner">
      <span class="material-symbols-outlined reveal-banner__icon">lock</span>
      <div class="reveal-banner__body">
        <div class="reveal-banner__title">Save your API key — it won't be shown again</div>
        <div class="copy-field">
          <span class="copy-field__value">{{ newKey.key }}</span>
          <button class="copy-field__btn" @click="copyKey(newKey.key)">
            <span class="material-symbols-outlined">{{ copied ? 'check' : 'content_copy' }}</span>
            {{ copied ? 'Copied' : 'Copy' }}
          </button>
        </div>
      </div>
      <button class="btn-icon" @click="newKey = null">
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>

    <!-- Create form (inline) -->
    <div v-if="showCreate" class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">add_circle</span>
          </div>
          <div class="an-card-title">Create New API Key</div>
        </div>
      </div>
      <div class="an-card-body">
        <div class="create-form-fields">
          <div class="form-field form-field--grow">
            <label class="form-field__label">Name <span class="form-field__required">*</span></label>
            <input
              ref="nameInputRef"
              type="text"
              class="form-input"
              :value="form.name"
              @input="form.name = ($event.target as HTMLInputElement).value"
              placeholder="e.g. My integration, CI/CD pipeline"
              maxlength="100"
              @keydown.enter="createKey"
            />
          </div>
          <div class="form-field">
            <label class="form-field__label">Expiry <span class="form-field__optional">(optional)</span></label>
            <input
              type="date"
              class="form-input form-input--date"
              :value="form.expires_at"
              @input="form.expires_at = ($event.target as HTMLInputElement).value"
              :min="tomorrow"
            />
          </div>
          <div class="create-form-actions">
            <button class="btn-filled" @click="createKey" :disabled="!form.name.trim() || creating">
              <div v-if="creating" class="css-spinner css-spinner--sm css-spinner--white"></div>
              Create
            </button>
            <button class="btn-outlined" @click="cancelCreate">Cancel</button>
          </div>
        </div>
        <div v-if="createError" class="create-error">{{ createError }}</div>
      </div>
    </div>

    <!-- Keys table card -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">key</span>
          </div>
          <div class="an-card-title">Your API Keys</div>
        </div>
        <span class="m3-badge m3-badge--neutral">{{ keys.length }} key{{ keys.length !== 1 ? 's' : '' }}</span>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="css-spinner"></div>
      </div>

      <div v-else-if="keys.length === 0" class="m3-empty-state">
        <div class="m3-empty-state__icon">
          <span class="material-symbols-outlined">key</span>
        </div>
        <div class="m3-empty-state__title">No API keys yet</div>
        <p class="m3-empty-state__text">Create an API key to authenticate programmatic requests.</p>
        <button class="btn-filled" @click="showCreate = true">
          <span class="material-symbols-outlined">add</span>
          Create your first key
        </button>
      </div>

      <div v-else class="m3-table-wrapper">
        <table class="m3-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Key Prefix</th>
              <th>Last Used</th>
              <th>Expires</th>
              <th>Created</th>
              <th class="th-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="key in keys" :key="key.id">
              <td>
                <div class="key-name">{{ key.name }}</div>
              </td>
              <td>
                <code class="prefix-code">{{ key.key_prefix }}••••••••</code>
              </td>
              <td class="td-muted">
                <span v-if="key.last_used_at">{{ formatDate(key.last_used_at) }}</span>
                <span v-else class="m3-badge m3-badge--neutral">Never used</span>
              </td>
              <td>
                <span v-if="key.expires_at">
                  <span v-if="isExpired(key.expires_at)" class="m3-badge m3-badge--error">Expired {{ formatDate(key.expires_at) }}</span>
                  <span v-else class="td-muted">{{ formatDate(key.expires_at) }}</span>
                </span>
                <span v-else class="m3-badge m3-badge--neutral">No expiry</span>
              </td>
              <td class="td-muted">{{ formatDate(key.created_at) }}</td>
              <td class="th-right">
                <div class="key-actions">
                  <button
                    class="btn-icon btn-icon--danger"
                    title="Revoke key"
                    :disabled="revokingId === key.id"
                    @click="promptRevoke(key)"
                  >
                    <div v-if="revokingId === key.id" class="css-spinner css-spinner--sm"></div>
                    <span v-else class="material-symbols-outlined">delete</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- How to use card -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--neutral">
            <span class="material-symbols-outlined">code</span>
          </div>
          <div class="an-card-title">How to use</div>
        </div>
      </div>
      <div class="an-card-body">
        <p class="usage-intro">
          Include the key in the <code class="inline-code">X-API-Key</code> header of every request:
        </p>
        <pre class="code-block">curl https://your-domain.com/api/v1/links \
  -H "X-API-Key: sl_your_api_key_here"</pre>
        <div class="usage-note">
          <span class="material-symbols-outlined usage-note__icon">info</span>
          <p class="usage-note__text">API keys have the same access level as your account. Revoke keys you no longer need to keep your account secure.</p>
        </div>
      </div>
    </div>

    <!-- Confirm Revoke Dialog -->
    <BaseModal v-model="showRevokeConfirm" size="sm" :persistent="false">
      <template #headline>
        <span class="material-symbols-outlined revoke-icon">delete</span>
        Revoke API Key
      </template>
      <p class="revoke-msg">
        Revoke <strong>{{ keyToRevoke?.name }}</strong>? Any integrations using it will stop working immediately. This cannot be undone.
      </p>
      <template #actions>
        <button class="btn-text" @click="showRevokeConfirm = false">Cancel</button>
        <button class="btn-filled btn-danger" :disabled="revokingId !== null" @click="confirmRevoke">
          <div v-if="revokingId !== null" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Revoke Key
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import { RouterLink } from 'vue-router';
import UpgradeWall from '@/components/UpgradeWall.vue';
import { useBillingStore } from '@/stores/billing';

const billingStore = useBillingStore();
import BaseModal from '@/components/BaseModal.vue';
import apiKeysApi from '@/api/apikeys';
import type { APIKey, CreateAPIKeyResponse } from '@/types/apikeys';
import billingApi from '@/api/billing';
import type { UsageCounts, PlanInfo } from '@/types/billing';

// ── Billing / Usage ───────────────────────────────────────────────────────────
const usage = ref<UsageCounts | null>(null);
const plan = ref<PlanInfo | null>(null);

const usageWarning = computed(() => {
  if (!usage.value || !plan.value) return null;
  const used = usage.value.api_keys;
  const max = plan.value.max_api_keys;
  if (max === -1) return null;
  const pct = used / max;
  if (pct >= 1) return { level: 'danger', msg: `You've reached your limit of ${max} API keys. Upgrade to add more.` };
  if (pct >= 0.8) return { level: 'warning', msg: `You've used ${used} of ${max} API keys (${Math.round(pct * 100)}%). Consider upgrading.` };
  return null;
});

// ── State ──────────────────────────────────────────────────────────────────────
const keys = ref<APIKey[]>([]);
const loading = ref(false);
const showCreate = ref(false);
const creating = ref(false);
const createError = ref('');
const revokingId = ref<string | null>(null);
const newKey = ref<CreateAPIKeyResponse | null>(null);
const copied = ref(false);

// Confirm revoke dialog
const showRevokeConfirm = ref(false);
const keyToRevoke = ref<APIKey | null>(null);

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
    newKey.value = res.data ?? null;
    copied.value = false;
    if (res.data) keys.value.unshift(res.data);
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
function promptRevoke(key: APIKey) {
  keyToRevoke.value = key;
  showRevokeConfirm.value = true;
}

async function confirmRevoke() {
  if (!keyToRevoke.value) return;
  const key = keyToRevoke.value;
  revokingId.value = key.id;
  try {
    await apiKeysApi.revoke(key.id);
    keys.value = keys.value.filter(k => k.id !== key.id);
    if (newKey.value?.id === key.id) newKey.value = null;
    showRevokeConfirm.value = false;
    keyToRevoke.value = null;
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

onMounted(async () => {
  billingStore.fetchPlan();
  await loadKeys();
  try {
    const res = await billingApi.getUsage();
    if (res.data.data) usage.value = res.data.data;
  } catch {}
  try {
    const res = await billingApi.getSubscription();
    if (res.data.data) plan.value = res.data.data.plan;
  } catch {}
});
</script>

<style scoped lang="scss">
/* ── Root ─────────────────────────────────────────────────────────── */
.apikeys-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0;
  max-width: 900px;
}

/* ── Page header ─────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.page-header__left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0;
}

.page-header__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

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

  &--neutral {
    background: var(--md-sys-color-surface-container, #f3f4f6);
    color: var(--md-sys-color-on-surface-variant);
  }
}

.an-card-body {
  padding: 16px 20px;
}

/* ── Alert banners ───────────────────────────────────────────────── */
.alert-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  flex-wrap: wrap;

  &--error {
    background: var(--md-sys-color-error-container, #ffdad6);
    color: var(--md-sys-color-on-error-container, #410002);
  }

  &--warning {
    background: #fef3c7;
    color: #92400e;
  }

  &__icon { flex-shrink: 0; }
  &__msg  { flex: 1; font-size: 0.875rem; }
}

/* ── New key reveal banner ───────────────────────────────────────── */
.reveal-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 14px;

  &__icon {
    color: #b45309;
    font-size: 1.25rem;
    flex-shrink: 0;
    margin-top: 2px;
  }

  &__body {
    flex: 1;
    min-width: 0;
  }

  &__title {
    font-weight: 600;
    margin-bottom: 8px;
    font-size: 0.9rem;
    color: #92400e;
  }
}

/* ── Copy field ──────────────────────────────────────────────────── */
.copy-field {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  padding: 6px 8px 6px 12px;
  flex-wrap: wrap;

  &__value {
    font-family: 'SFMono-Regular', Consolas, monospace;
    font-size: 0.78rem;
    color: var(--md-sys-color-on-surface);
    word-break: break-all;
    flex: 1;
    min-width: 0;
  }

  &__btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 10px;
    border-radius: 6px;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.8rem;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    transition: background 0.15s, color 0.15s;
    flex-shrink: 0;

    &:hover {
      background: var(--md-sys-color-surface-container);
      color: var(--md-sys-color-on-surface);
    }

    .material-symbols-outlined { font-size: 16px; }
  }
}

/* ── Form fields ─────────────────────────────────────────────────── */
.create-form-fields {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: flex-end;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 4px;

  &--grow {
    flex: 2;
    min-width: 220px;
  }

  &__label {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--md-sys-color-on-surface-variant);
  }

  &__required { color: var(--md-sys-color-error); }
  &__optional { font-weight: 400; }
}

.form-input {
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

  &::placeholder { color: var(--md-sys-color-on-surface-variant); opacity: 0.7; }

  &:focus { border-color: var(--md-sys-color-primary, #635bff); }

  &--date { min-width: 160px; }
}

.create-form-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  padding-top: 20px;
}

.create-error {
  color: var(--md-sys-color-error);
  font-size: 0.875rem;
  margin-top: 8px;
}

/* ── Loading state ───────────────────────────────────────────────── */
.loading-state {
  display: flex;
  justify-content: center;
  padding: 40px;
}

/* ── CSS spinner ─────────────────────────────────────────────────── */
.css-spinner {
  width: 28px;
  height: 28px;
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

/* ── Empty state ─────────────────────────────────────────────────── */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  gap: 4px;

  &__icon {
    width: 64px;
    height: 64px;
    border-radius: 20px;
    background: var(--md-sys-color-surface-container-low);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 12px;

    .material-symbols-outlined {
      font-size: 2rem;
      color: var(--md-sys-color-on-surface-variant);
      opacity: 0.6;
    }
  }

  &__title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface);
    margin-bottom: 4px;
  }

  &__text {
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface-variant);
    margin: 0 0 16px;
  }
}

/* ── Table ───────────────────────────────────────────────────────── */
.m3-table-wrapper { overflow-x: auto; }

.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;

  th {
    padding: 10px 16px;
    text-align: left;
    font-weight: 600;
    font-size: 0.75rem;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    white-space: nowrap;
  }

  td {
    padding: 12px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td { border-bottom: none; }
  tbody tr:hover td { background: var(--md-sys-color-surface-container-low); }
}

.th-right { text-align: right; }
.td-muted { color: var(--md-sys-color-on-surface-variant); font-size: 0.8rem; }

/* ── Key cells ───────────────────────────────────────────────────── */
.key-name {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.prefix-code {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 0.8rem;
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  padding: 2px 8px;
  border-radius: 4px;
}

.key-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 4px;
}

.btn-icon--danger {
  color: var(--md-sys-color-error);

  &:hover {
    background: rgba(220, 38, 38, 0.08);
  }
}

/* ── Badges ──────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  &--error {
    background: rgba(220, 38, 38, 0.12);
    color: #dc2626;
  }

  &--success {
    background: rgba(22, 163, 74, 0.12);
    color: #16a34a;
  }
}

/* ── Usage / How to use card ─────────────────────────────────────── */
.usage-intro {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0 0 12px;
}

.inline-code {
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-primary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 0.85em;
}

.code-block {
  background: var(--md-sys-color-surface-container-low);
  border-radius: 8px;
  padding: 12px 16px;
  font-size: 0.8rem;
  font-family: 'SFMono-Regular', Consolas, monospace;
  margin: 0;
  overflow-x: auto;
  color: var(--md-sys-color-on-surface);
}

.usage-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 12px;
  padding: 10px 14px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 8px;
  color: var(--md-sys-color-on-surface-variant);

  &__icon {
    color: var(--md-sys-color-primary);
    font-size: 16px;
    flex-shrink: 0;
    margin-top: 1px;
  }

  &__text {
    margin: 0;
    font-size: 0.875rem;
  }
}

/* ── Revoke dialog ───────────────────────────────────────────────── */
.revoke-icon { color: var(--md-sys-color-error); }

.revoke-msg {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  font-size: 0.9rem;
}
</style>
