<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">API Keys</h1>
        <p class="page-subtitle">
          Use API keys to authenticate programmatic requests with
          <code style="color:var(--md-sys-color-primary);">X-API-Key: &lt;key&gt;</code> header.
        </p>
      </div>
      <md-filled-button @click="showCreate = true">
        <span class="material-symbols-outlined" slot="icon">add</span>
        New API Key
      </md-filled-button>
    </div>

    <!-- Usage Warning Banner -->
    <div v-if="usageWarning" :class="['warning-banner', usageWarning.level === 'danger' ? 'warning-banner--error' : 'warning-banner--warning']">
      <span class="material-symbols-outlined">{{ usageWarning.level === 'danger' ? 'block' : 'warning' }}</span>
      <span style="flex:1;font-size:0.875rem;">{{ usageWarning.msg }}</span>
      <RouterLink to="/dashboard/billing">
        <md-filled-tonal-button>Upgrade</md-filled-tonal-button>
      </RouterLink>
    </div>

    <!-- One-time key reveal banner -->
    <div v-if="newKey" class="new-key-banner">
      <span class="material-symbols-outlined" style="color:#b45309;font-size:1.25rem;flex-shrink:0;margin-top:2px;">lock</span>
      <div style="flex:1;min-width:0;">
        <div style="font-weight:600;margin-bottom:8px;">Save your API key — it won't be shown again</div>
        <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
          <code class="key-display">{{ newKey.key }}</code>
          <md-filled-tonal-button @click="copyKey(newKey.key)">
            <span class="material-symbols-outlined" slot="icon">{{ copied ? 'check' : 'content_copy' }}</span>
            {{ copied ? 'Copied' : 'Copy' }}
          </md-filled-tonal-button>
        </div>
      </div>
      <md-icon-button @click="newKey = null">
        <span class="material-symbols-outlined">close</span>
      </md-icon-button>
    </div>

    <!-- Create form (inline) -->
    <div v-if="showCreate" class="m3-card m3-card--elevated create-form">
      <div class="m3-card-header">
        <div class="m3-card-header__left">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">add_circle</span>
          <span class="md-title-medium">Create New API Key</span>
        </div>
      </div>
      <md-divider />
      <div style="padding:1.25rem;">
      <div class="create-form-fields">
        <md-outlined-text-field
          :value="form.name"
          @input="form.name=($event.target as HTMLInputElement).value"
          label="Name *"
          placeholder="e.g. My integration, CI/CD pipeline"
          maxlength="100"
          style="flex:2;min-width:220px;"
          @keydown.enter="createKey"
          ref="nameInputRef"
        />
        <md-outlined-text-field
          :value="form.expires_at"
          @input="form.expires_at=($event.target as HTMLInputElement).value"
          label="Expiry (optional)"
          type="date"
          :min="tomorrow"
          style="flex:1;min-width:160px;"
        />
        <div style="display:flex;gap:8px;align-items:center;">
          <md-filled-button @click="createKey" :disabled="!form.name.trim() || creating">
            <md-circular-progress v-if="creating" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
            Create
          </md-filled-button>
          <md-outlined-button @click="cancelCreate">Cancel</md-outlined-button>
        </div>
      </div>
      <div v-if="createError" style="color:var(--md-sys-color-error);font-size:0.875rem;margin-top:8px;">{{ createError }}</div>
      </div>
    </div>

    <!-- Keys table -->
    <div class="m3-card m3-card--elevated">
      <div class="m3-card-header">
        <div class="m3-card-header__left">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">key</span>
          <span class="md-title-medium">Your API Keys</span>
        </div>
        <span class="m3-badge m3-badge--neutral">{{ keys.length }} key{{ keys.length !== 1 ? 's' : '' }}</span>
      </div>
      <md-divider />

      <div v-if="loading" style="display:flex;justify-content:center;padding:40px;">
        <md-circular-progress indeterminate style="--md-circular-progress-size:32px" />
      </div>

      <div v-else-if="keys.length === 0 && !loading" class="m3-empty-state">
        <div class="m3-empty-state__icon">
          <span class="material-symbols-outlined">key</span>
        </div>
        <div class="md-title-medium" style="margin-bottom:0.5rem;">No API keys yet</div>
        <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0 0 1rem;">Create an API key to authenticate programmatic requests.</p>
        <md-filled-button @click="showCreate = true">
          <span class="material-symbols-outlined" slot="icon">add</span>
          Create your first key
        </md-filled-button>
      </div>

      <div v-else class="m3-table-wrapper">
        <table class="m3-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Key Prefix</th>
              <th>Last used</th>
              <th>Expires</th>
              <th>Created</th>
              <th style="text-align:right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="key in keys" :key="key.id">
              <td>
                <div class="key-name">{{ key.name }}</div>
              </td>
              <td>
                <div class="key-prefix-cell">
                  <code class="prefix-code">{{ key.key_prefix }}••••••••</code>
                </div>
              </td>
              <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;">
                <span v-if="key.last_used_at">{{ formatDate(key.last_used_at) }}</span>
                <span v-else class="m3-badge m3-badge--neutral" style="font-size:0.7rem;">Never used</span>
              </td>
              <td style="font-size:0.8rem;">
                <span v-if="key.expires_at">
                  <span v-if="isExpired(key.expires_at)" class="m3-badge m3-badge--error">Expired {{ formatDate(key.expires_at) }}</span>
                  <span v-else style="color:var(--md-sys-color-on-surface-variant);">{{ formatDate(key.expires_at) }}</span>
                </span>
                <span v-else class="m3-badge m3-badge--neutral" style="font-size:0.7rem;">No expiry</span>
              </td>
              <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;">{{ formatDate(key.created_at) }}</td>
              <td style="text-align:right;">
                <div class="key-actions">
                  <md-icon-button
                    title="Revoke key"
                    :disabled="revokingId === key.id"
                    @click="revokeKey(key)"
                    class="danger-btn"
                  >
                    <md-circular-progress v-if="revokingId === key.id" indeterminate style="--md-circular-progress-size:20px" />
                    <span v-else class="material-symbols-outlined">delete</span>
                  </md-icon-button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Usage docs -->
    <div class="m3-card m3-card--elevated">
      <div class="m3-card-header">
        <div class="m3-card-header__left">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">code</span>
          <span class="md-title-medium">How to use</span>
        </div>
      </div>
      <md-divider />
      <div class="usage-card">
        <p style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;margin:0 0 12px;">Include the key in the <code style="background:var(--md-sys-color-surface-container-low);padding:2px 6px;border-radius:4px;">X-API-Key</code> header of every request:</p>
        <pre class="code-block">curl https://your-domain.com/api/v1/links \
  -H "X-API-Key: sl_your_api_key_here"</pre>
        <div class="usage-note">
          <span class="material-symbols-outlined" style="font-size:16px;flex-shrink:0;">info</span>
          <p style="margin:0;font-size:0.875rem;">API keys have the same access level as your account. Revoke keys you no longer need to keep your account secure.</p>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import { RouterLink } from 'vue-router';
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

onMounted(async () => {
  await loadKeys();
  try {
    const res = await billingApi.getUsage();
    usage.value = res.data.data;
  } catch {}
  try {
    const res = await billingApi.getSubscription();
    plan.value = res.data.data.plan;
  } catch {}
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

.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
  margin-bottom: 20px;

  &--outlined {
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
  }
}

/* M3 Card header */
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

/* M3 Table wrapper */
.m3-table-wrapper {
  overflow-x: auto;
}

/* M3 Badges */
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

  &--error {
    background: #fee2e2;
    color: #dc2626;
  }
}

/* M3 Empty state */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;

  &__icon {
    width: 64px;
    height: 64px;
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

/* Key name */
.key-name {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

/* Key prefix cell */
.key-prefix-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Key actions */
.key-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 4px;
}

/* Danger button */
.danger-btn {
  --md-icon-button-icon-color: var(--md-sys-color-error);
}

/* Warning banner */
.warning-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.warning-banner--error {
  background: var(--md-sys-color-error-container, #ffdad6);
  color: var(--md-sys-color-on-error-container, #410002);
}

.warning-banner--warning {
  background: #fef3c7;
  color: #92400e;
}

/* New key banner */
.new-key-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 12px;
  margin-bottom: 20px;
}

.key-display {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 0.78rem;
  background: #f8f9fa;
  color: #212529;
  padding: 8px 12px;
  border-radius: 8px;
  word-break: break-all;
  flex: 1;
  display: inline-block;
}

/* Create form */
.create-form {
  // card header adds its own padding
}

.create-form-fields {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: flex-end;
}

/* Table */
.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.m3-table thead tr {
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.m3-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  background: var(--md-sys-color-surface-container-low);
  white-space: nowrap;
}

.m3-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  color: var(--md-sys-color-on-surface);
}

.m3-table tbody tr:last-child td {
  border-bottom: none;
}

.m3-table tbody tr:hover td {
  background: var(--md-sys-color-surface-container-low);
}

.prefix-code {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 0.8rem;
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  padding: 2px 8px;
  border-radius: 4px;
}

/* Usage docs */
.usage-card {
  padding: 20px;
}

/* Usage note */
.usage-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 12px;
  padding: 10px 14px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 8px;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;

  .material-symbols-outlined {
    color: var(--md-sys-color-primary);
    margin-top: 1px;
  }
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
</style>
