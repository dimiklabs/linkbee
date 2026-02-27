<template>
  <div class="page-wrapper">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Webhooks</h1>
        <p class="page-subtitle">Receive real-time HTTP notifications when events happen.</p>
      </div>
      <md-filled-button @click="showCreate = !showCreate">
        <span class="material-symbols-outlined" slot="icon">add</span>
        Add Webhook
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

    <!-- Create form -->
    <div v-if="showCreate" class="m3-card m3-card--outlined create-form">
      <h6 style="margin:0 0 16px;font-weight:600;">New Webhook</h6>
      <div style="margin-bottom:16px;">
        <md-outlined-text-field
          :value="newURL"
          @input="newURL=($event.target as HTMLInputElement).value"
          label="Endpoint URL"
          type="url"
          placeholder="https://your-server.com/webhook"
          style="width:100%;"
          ref="urlInputRef"
        />
      </div>
      <div style="margin-bottom:16px;">
        <div style="font-size:0.875rem;font-weight:500;color:var(--md-sys-color-on-surface);margin-bottom:8px;">Events to subscribe</div>
        <div style="display:flex;flex-direction:column;gap:8px;">
          <label v-for="ev in WEBHOOK_EVENTS" :key="ev.value" style="display:flex;align-items:flex-start;gap:10px;cursor:pointer;">
            <input
              type="checkbox"
              :value="ev.value"
              v-model="newEvents"
              style="margin-top:2px;accent-color:var(--md-sys-color-primary);"
            />
            <div>
              <span style="font-weight:500;font-size:0.875rem;">{{ ev.label }}</span>
              <span style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;"> — {{ ev.description }}</span>
            </div>
          </label>
        </div>
        <div v-if="newEvents.length === 0" style="color:var(--md-sys-color-error);font-size:0.8rem;margin-top:4px;">
          Select at least one event.
        </div>
      </div>
      <div style="display:flex;gap:8px;align-items:center;">
        <md-filled-button :disabled="creating || !newURL || newEvents.length === 0" @click="createWebhook">
          <md-circular-progress v-if="creating" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
          Create
        </md-filled-button>
        <md-outlined-button @click="cancelCreate">Cancel</md-outlined-button>
      </div>
      <div v-if="createError" style="margin-top:12px;padding:10px 14px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:8px;font-size:0.875rem;">
        {{ createError }}
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate style="--md-circular-progress-size:40px" />
    </div>

    <!-- Empty state -->
    <div v-else-if="webhooks.length === 0 && !showCreate" class="m3-card m3-card--outlined empty-state">
      <span class="material-symbols-outlined" style="font-size:2.5rem;color:var(--md-sys-color-on-surface-variant);">notifications</span>
      <h6 style="font-weight:600;margin:12px 0 4px;">No webhooks yet</h6>
      <p style="color:var(--md-sys-color-on-surface-variant);font-size:0.875rem;margin:0 0 16px;">Add a webhook endpoint to start receiving event notifications.</p>
      <md-filled-button @click="showCreate = true">Add Webhook</md-filled-button>
    </div>

    <!-- Webhooks table -->
    <div v-else-if="webhooks.length > 0" class="m3-card m3-card--outlined">
      <div class="table-container">
        <table class="m3-table">
          <thead>
            <tr>
              <th>Endpoint</th>
              <th>Events</th>
              <th>Status</th>
              <th>Created</th>
              <th style="text-align:right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="wh in webhooks" :key="wh.id">
              <tr>
                <!-- Endpoint -->
                <td>
                  <div v-if="editingId === wh.id">
                    <md-outlined-text-field
                      :value="editURL"
                      @input="editURL=($event.target as HTMLInputElement).value"
                      type="url"
                      label="URL"
                      style="min-width:240px;"
                    />
                  </div>
                  <code v-else style="font-size:0.8rem;word-break:break-all;">{{ wh.url }}</code>
                </td>

                <!-- Events -->
                <td>
                  <div v-if="editingId === wh.id" style="display:flex;flex-direction:column;gap:6px;">
                    <label v-for="ev in WEBHOOK_EVENTS" :key="ev.value" style="display:flex;align-items:center;gap:8px;cursor:pointer;font-size:0.8rem;">
                      <input
                        type="checkbox"
                        :value="ev.value"
                        v-model="editEvents"
                        style="accent-color:var(--md-sys-color-primary);"
                      />
                      {{ ev.label }}
                    </label>
                  </div>
                  <div v-else style="display:flex;flex-wrap:wrap;gap:4px;">
                    <span
                      v-for="ev in wh.events"
                      :key="ev"
                      :class="['m3-badge', eventBadgeClass(ev)]"
                    >{{ ev }}</span>
                  </div>
                </td>

                <!-- Status -->
                <td>
                  <div v-if="editingId === wh.id" style="display:flex;align-items:center;gap:8px;">
                    <input type="checkbox" v-model="editIsActive" style="accent-color:var(--md-sys-color-primary);" />
                    <span style="font-size:0.875rem;">Active</span>
                  </div>
                  <span v-else :class="['m3-badge', wh.is_active ? 'm3-badge--success' : 'm3-badge--neutral']">
                    {{ wh.is_active ? 'Active' : 'Paused' }}
                  </span>
                </td>

                <!-- Created -->
                <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;white-space:nowrap;">{{ formatDate(wh.created_at) }}</td>

                <!-- Actions -->
                <td>
                  <div v-if="editingId === wh.id" style="display:flex;justify-content:flex-end;gap:8px;">
                    <md-filled-button :disabled="saving" @click="saveEdit(wh.id)">
                      <md-circular-progress v-if="saving" indeterminate style="--md-circular-progress-size:18px;margin-right:6px;" />
                      Save
                    </md-filled-button>
                    <md-outlined-button @click="cancelEdit">Cancel</md-outlined-button>
                  </div>
                  <div v-else style="display:flex;justify-content:flex-end;gap:4px;flex-wrap:wrap;">
                    <md-icon-button @click="sendTest(wh.id)" :disabled="testingId === wh.id" title="Send test event">
                      <md-circular-progress v-if="testingId === wh.id" indeterminate style="--md-circular-progress-size:20px" />
                      <span v-else class="material-symbols-outlined">send</span>
                    </md-icon-button>
                    <md-icon-button @click="toggleDeliveries(wh.id)" title="Delivery logs">
                      <span class="material-symbols-outlined">history</span>
                    </md-icon-button>
                    <md-icon-button @click="startEdit(wh)" title="Edit">
                      <span class="material-symbols-outlined">edit</span>
                    </md-icon-button>
                    <md-icon-button @click="deleteWebhook(wh.id)" title="Delete" style="--md-icon-button-icon-color:var(--md-sys-color-error);">
                      <span class="material-symbols-outlined">delete</span>
                    </md-icon-button>
                  </div>
                </td>
              </tr>

              <!-- Delivery logs expandable row -->
              <tr v-if="openDeliveriesId === wh.id">
                <td colspan="5" style="padding:0;">
                  <div class="deliveries-panel">

                    <!-- Secret reveal -->
                    <div style="display:flex;align-items:center;gap:8px;margin-bottom:16px;flex-wrap:wrap;">
                      <span style="font-size:0.875rem;font-weight:500;color:var(--md-sys-color-on-surface-variant);">Signing secret:</span>
                      <code v-if="revealedSecret[wh.id]" style="font-size:0.8rem;">{{ revealedSecret[wh.id] }}</code>
                      <span v-else style="font-size:0.8rem;font-family:monospace;color:var(--md-sys-color-on-surface-variant);">••••••••••••••••</span>
                      <md-text-button @click="toggleSecret(wh.id)">
                        {{ revealedSecret[wh.id] ? 'Hide' : 'Reveal' }}
                      </md-text-button>
                      <md-text-button v-if="revealedSecret[wh.id]" @click="copyText(revealedSecret[wh.id])">Copy</md-text-button>
                    </div>

                    <!-- Test result banner -->
                    <div v-if="testResults[wh.id]" style="margin-bottom:16px;">
                      <div :style="`padding:10px 14px;border-radius:8px;font-size:0.875rem;background:${testResults[wh.id].success ? '#dcfce7' : '#fee2e2'};color:${testResults[wh.id].success ? '#16a34a' : '#dc2626'};`">
                        <strong>Test delivery:</strong>
                        {{ testResults[wh.id].success ? `✓ HTTP ${testResults[wh.id].response_code}` : `✗ ${testResults[wh.id].error_message || 'Failed'}` }}
                        <span style="opacity:0.7;margin-left:8px;">({{ testResults[wh.id].duration_ms }}ms)</span>
                      </div>
                    </div>

                    <!-- Deliveries table -->
                    <div style="font-size:0.875rem;font-weight:600;margin-bottom:8px;">Recent Deliveries</div>
                    <div v-if="deliveriesLoading[wh.id]" style="display:flex;justify-content:center;padding:16px;">
                      <md-circular-progress indeterminate style="--md-circular-progress-size:24px" />
                    </div>
                    <div v-else-if="!deliveries[wh.id] || deliveries[wh.id].length === 0" style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;padding:8px 0;">
                      No deliveries recorded yet. Send a test event to get started.
                    </div>
                    <div v-else>
                      <table class="m3-table" style="font-size:0.8rem;">
                        <thead>
                          <tr>
                            <th>Time</th>
                            <th>Event</th>
                            <th>Status</th>
                            <th>Duration</th>
                            <th></th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr v-for="d in deliveries[wh.id]" :key="d.id">
                            <td style="white-space:nowrap;color:var(--md-sys-color-on-surface-variant);">{{ formatDateTime(d.created_at) }}</td>
                            <td><code>{{ d.event }}</code></td>
                            <td>
                              <span :class="['m3-badge', d.success ? 'm3-badge--success' : 'm3-badge--error']">{{ d.success ? d.response_code : (d.response_code || 'Error') }}</span>
                            </td>
                            <td style="color:var(--md-sys-color-on-surface-variant);">{{ d.duration_ms }}ms</td>
                            <td style="text-align:right;">
                              <md-text-button
                                @click="resendDelivery(wh.id, d.id)"
                                :disabled="resendingId === d.id"
                                style="font-size:0.75rem;"
                              >
                                <md-circular-progress v-if="resendingId === d.id" indeterminate style="--md-circular-progress-size:16px" />
                                <span v-else>Resend</span>
                              </md-text-button>
                            </td>
                          </tr>
                        </tbody>
                      </table>

                      <!-- Pagination -->
                      <div v-if="deliveryTotals[wh.id] > 20" style="display:flex;justify-content:space-between;align-items:center;margin-top:8px;">
                        <span style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;">{{ deliveryTotals[wh.id] }} total</span>
                        <div style="display:flex;gap:4px;">
                          <md-icon-button :disabled="deliveryPages[wh.id] <= 1" @click="loadDeliveries(wh.id, deliveryPages[wh.id] - 1)">
                            <span class="material-symbols-outlined">chevron_left</span>
                          </md-icon-button>
                          <md-icon-button :disabled="deliveryPages[wh.id] * 20 >= deliveryTotals[wh.id]" @click="loadDeliveries(wh.id, deliveryPages[wh.id] + 1)">
                            <span class="material-symbols-outlined">chevron_right</span>
                          </md-icon-button>
                        </div>
                      </div>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Signing info card -->
    <div class="m3-card m3-card--outlined info-card">
      <h6 style="font-weight:600;margin:0 0 8px;">Verifying webhook signatures</h6>
      <p style="font-size:0.875rem;color:var(--md-sys-color-on-surface-variant);margin:0 0 8px;">
        Every delivery includes an <code>X-Webhook-Signature</code> header with an HMAC-SHA256 signature of the raw request body, prefixed with <code>sha256=</code>.
        Your signing secret was generated when you created the webhook (stored server-side only).
        Verify on your server to ensure the request is genuine:
      </p>
      <pre class="code-block"><code>import hmac, hashlib

def verify(secret: str, body: bytes, header: str) -> bool:
    expected = "sha256=" + hmac.new(
        secret.encode(), body, hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, header)</code></pre>
    </div>

    <!-- Example payload card -->
    <div class="m3-card m3-card--outlined info-card">
      <h6 style="font-weight:600;margin:0 0 8px;">Example payload</h6>
      <pre class="code-block"><code>{
  "event": "link.created",
  "timestamp": "2025-01-01T12:00:00Z",
  "data": {
    "id": "...",
    "slug": "abc12",
    "short_url": "https://sl.ink/abc12",
    "destination_url": "https://example.com"
  }
}</code></pre>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue';
import { RouterLink } from 'vue-router';
import webhooksApi from '@/api/webhooks';
import { WEBHOOK_EVENTS } from '@/types/webhooks';
import type { Webhook, WebhookDelivery } from '@/types/webhooks';
import billingApi from '@/api/billing';
import type { UsageCounts, PlanInfo } from '@/types/billing';

// ── Billing / Usage ───────────────────────────────────────────────────────────
const usage = ref<UsageCounts | null>(null);
const plan = ref<PlanInfo | null>(null);

const usageWarning = computed(() => {
  if (!usage.value || !plan.value) return null;
  const used = usage.value.webhooks;
  const max = plan.value.max_webhooks;
  if (max === -1) return null;
  const pct = used / max;
  if (pct >= 1) return { level: 'danger', msg: `You've reached your limit of ${max} webhooks. Upgrade to add more.` };
  if (pct >= 0.8) return { level: 'warning', msg: `You've used ${used} of ${max} webhooks (${Math.round(pct * 100)}%). Consider upgrading.` };
  return null;
});

const webhooks = ref<Webhook[]>([]);
const loading = ref(false);
const showCreate = ref(false);

// Create form
const newURL = ref('');
const newEvents = ref<string[]>([]);
const creating = ref(false);
const createError = ref('');
const urlInputRef = ref<HTMLInputElement | null>(null);

// Edit state
const editingId = ref<string | null>(null);
const editURL = ref('');
const editEvents = ref<string[]>([]);
const editIsActive = ref(true);
const saving = ref(false);

// Delivery logs state
const openDeliveriesId = ref<string | null>(null);
const deliveries = ref<Record<string, WebhookDelivery[]>>({});
const deliveriesLoading = ref<Record<string, boolean>>({});
const deliveryTotals = ref<Record<string, number>>({});
const deliveryPages = ref<Record<string, number>>({});
const testingId = ref<string | null>(null);
const testResults = ref<Record<string, WebhookDelivery>>({});
const revealedSecret = ref<Record<string, string>>({});
const resendingId = ref<string | null>(null);

async function fetchWebhooks() {
  loading.value = true;
  try {
    const res = await webhooksApi.list();
    webhooks.value = res.data ?? [];
  } finally {
    loading.value = false;
  }
}

function cancelCreate() {
  showCreate.value = false;
  newURL.value = '';
  newEvents.value = [];
  createError.value = '';
}

async function createWebhook() {
  if (!newURL.value || newEvents.value.length === 0) return;
  creating.value = true;
  createError.value = '';
  try {
    const res = await webhooksApi.create({ url: newURL.value, events: newEvents.value });
    webhooks.value.unshift(res.data);
    cancelCreate();
  } catch (err: any) {
    createError.value = err?.response?.data?.description ?? 'Failed to create webhook.';
  } finally {
    creating.value = false;
  }
}

function startEdit(wh: Webhook) {
  editingId.value = wh.id;
  editURL.value = wh.url;
  editEvents.value = [...wh.events];
  editIsActive.value = wh.is_active;
}

function cancelEdit() {
  editingId.value = null;
}

async function saveEdit(id: string) {
  saving.value = true;
  try {
    const res = await webhooksApi.update(id, {
      url: editURL.value,
      events: editEvents.value,
      is_active: editIsActive.value,
    });
    const idx = webhooks.value.findIndex((w) => w.id === id);
    if (idx !== -1) webhooks.value[idx] = res.data;
    editingId.value = null;
  } finally {
    saving.value = false;
  }
}

async function deleteWebhook(id: string) {
  if (!confirm('Delete this webhook? This cannot be undone.')) return;
  await webhooksApi.delete(id);
  webhooks.value = webhooks.value.filter((w) => w.id !== id);
  delete deliveries.value[id];
  delete deliveriesLoading.value[id];
  delete deliveryTotals.value[id];
  delete deliveryPages.value[id];
  delete testResults.value[id];
  delete revealedSecret.value[id];
  if (openDeliveriesId.value === id) openDeliveriesId.value = null;
}

function eventBadgeClass(event: string) {
  if (event === 'link.created') return 'm3-badge--success';
  if (event === 'link.deleted') return 'm3-badge--error';
  if (event === 'link.clicked') return 'm3-badge--primary';
  return 'm3-badge--neutral';
}

function formatDate(iso: string) {
  if (!iso) return '—';
  return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
}

function formatDateTime(iso: string) {
  if (!iso) return '—';
  return new Date(iso).toLocaleString(undefined, {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  });
}

async function toggleDeliveries(id: string) {
  if (openDeliveriesId.value === id) {
    openDeliveriesId.value = null;
    return;
  }
  openDeliveriesId.value = id;
  if (!deliveries.value[id]) {
    await loadDeliveries(id, 1);
  }
}

async function loadDeliveries(id: string, page: number) {
  deliveriesLoading.value[id] = true;
  try {
    const res = await webhooksApi.getDeliveries(id, page);
    deliveries.value[id] = res.data.deliveries ?? [];
    deliveryTotals.value[id] = res.data.total ?? 0;
    deliveryPages.value[id] = res.data.page ?? page;
  } finally {
    deliveriesLoading.value[id] = false;
  }
}

async function sendTest(id: string) {
  testingId.value = id;
  if (openDeliveriesId.value !== id) {
    openDeliveriesId.value = id;
  }
  try {
    const res = await webhooksApi.test(id);
    testResults.value[id] = res.data;
    await loadDeliveries(id, 1);
  } catch (err: any) {
    testResults.value[id] = {
      id: '',
      webhook_id: id,
      user_id: '',
      event: 'test',
      request_body: '',
      response_code: 0,
      response_body: '',
      error_message: err?.response?.data?.description ?? 'Request failed',
      success: false,
      duration_ms: 0,
      created_at: new Date().toISOString(),
    };
  } finally {
    testingId.value = null;
  }
}

async function toggleSecret(id: string) {
  if (revealedSecret.value[id]) {
    delete revealedSecret.value[id];
    return;
  }
  try {
    const res = await webhooksApi.getSecret(id);
    revealedSecret.value[id] = res.data.secret;
  } catch {
    // ignore
  }
}

async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    // ignore
  }
}

async function resendDelivery(webhookId: string, deliveryId: string) {
  resendingId.value = deliveryId;
  try {
    const res = await webhooksApi.resendDelivery(webhookId, deliveryId);
    testResults.value[webhookId] = res.data;
    await loadDeliveries(webhookId, deliveryPages.value[webhookId] ?? 1);
  } finally {
    resendingId.value = null;
  }
}

watch(showCreate, async (val) => {
  if (val) {
    await nextTick();
    urlInputRef.value?.focus();
  }
});

onMounted(async () => {
  await fetchWebhooks();
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

<style scoped>
.page-wrapper {
  padding: 24px;
  max-width: 1100px;
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
}

.m3-card--outlined {
  border: 1px solid var(--md-sys-color-outline-variant);
  margin-bottom: 20px;
}

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

.create-form {
  padding: 20px;
}

.table-container {
  overflow-x: auto;
}

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
  vertical-align: middle;
}

.m3-table tbody tr:last-child td {
  border-bottom: none;
}

.m3-table tbody tr:hover td {
  background: var(--md-sys-color-surface-container-low);
}

.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
}

.m3-badge--primary {
  background: var(--md-sys-color-primary-container, #e8def8);
  color: var(--md-sys-color-on-primary-container, #21005d);
}

.m3-badge--neutral {
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  border: 1px solid var(--md-sys-color-outline-variant);
}

.m3-badge--success {
  background: #dcfce7;
  color: #16a34a;
}

.m3-badge--error {
  background: #fee2e2;
  color: #dc2626;
}

.deliveries-panel {
  background: var(--md-sys-color-surface-container-low);
  border-top: 1px solid var(--md-sys-color-outline-variant);
  padding: 16px 20px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 24px;
  text-align: center;
}

.info-card {
  padding: 20px;
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
