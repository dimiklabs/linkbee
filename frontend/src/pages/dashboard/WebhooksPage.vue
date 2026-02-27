<template>
  <div class="page-section" style="max-width: 1100px;">

    <!-- Page Header -->
    <div class="dash-page-header">
      <div class="dash-page-header__left">
        <h1 class="dash-page-header__title">Webhooks</h1>
        <p class="dash-page-header__subtitle">Receive real-time HTTP notifications when events happen.</p>
      </div>
      <div class="dash-page-header__actions">
        <button class="btn-filled" @click="showCreate = !showCreate">
          <span class="material-symbols-outlined" style="font-size:18px;margin-right:6px;">add</span>
          Add Webhook
        </button>
      </div>
    </div>

    <!-- Usage Warning Banner -->
    <div v-if="usageWarning" :class="['warning-banner', usageWarning.level === 'danger' ? 'warning-banner--error' : 'warning-banner--warning']">
      <span class="material-symbols-outlined">{{ usageWarning.level === 'danger' ? 'block' : 'warning' }}</span>
      <span style="flex:1;font-size:0.875rem;">{{ usageWarning.msg }}</span>
      <RouterLink to="/dashboard/billing">
        <button class="btn-tonal">Upgrade</button>
      </RouterLink>
    </div>

    <!-- Create form -->
    <div v-if="showCreate" class="m3-card m3-card--elevated create-form">
      <div class="m3-card-header">
        <div style="display:flex;align-items:center;gap:8px;">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">add_circle</span>
          <span style="font-size:16px;font-weight:600;">New Webhook Endpoint</span>
        </div>
      </div>
      <md-divider />
      <div style="padding:20px;">
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
          <button class="btn-filled" :disabled="creating || !newURL || newEvents.length === 0" @click="createWebhook">
            <md-circular-progress v-if="creating" indeterminate style="margin-right:6px;" />
            Create
          </button>
          <button class="btn-outlined" @click="cancelCreate">Cancel</button>
        </div>
        <div v-if="createError" style="margin-top:12px;padding:10px 14px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:8px;font-size:0.875rem;">
          {{ createError }}
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate />
    </div>

    <!-- Empty state -->
    <div v-else-if="webhooks.length === 0 && !showCreate" class="m3-card m3-card--elevated m3-empty-state">
      <div class="m3-empty-state__icon">
        <span class="material-symbols-outlined">notifications_active</span>
      </div>
      <div class="md-title-medium" style="margin-bottom:8px;">No webhooks yet</div>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0 0 20px;max-width:380px;">Add a webhook endpoint to start receiving real-time event notifications for link activity.</p>
      <button class="btn-filled" @click="showCreate = true">
        <span class="material-symbols-outlined" style="font-size:18px;margin-right:6px;">add</span>
        Add Webhook
      </button>
    </div>

    <!-- Webhooks table -->
    <div v-else-if="webhooks.length > 0" class="m3-card m3-card--elevated">
      <div class="m3-card-header">
        <div style="display:flex;align-items:center;gap:8px;">
          <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">notifications</span>
          <span style="font-size:16px;font-weight:600;">Webhook Endpoints</span>
        </div>
        <span class="m3-badge m3-badge--neutral">{{ webhooks.length }} endpoint{{ webhooks.length !== 1 ? 's' : '' }}</span>
      </div>
      <md-divider />
      <div class="m3-table-wrapper">
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
                  <div v-else class="event-chips">
                    <span
                      v-for="ev in wh.events"
                      :key="ev"
                      :class="['event-chip', eventChipClass(ev)]"
                    >{{ ev }}</span>
                  </div>
                </td>

                <!-- Status -->
                <td>
                  <div v-if="editingId === wh.id" style="display:flex;align-items:center;gap:8px;">
                    <input type="checkbox" v-model="editIsActive" style="accent-color:var(--md-sys-color-primary);" />
                    <span style="font-size:0.875rem;">Active</span>
                  </div>
                  <span v-else :class="['m3-badge', wh.is_active ? 'm3-badge--success' : 'm3-badge--warning']">
                    <span class="material-symbols-outlined" style="font-size:12px;vertical-align:middle;margin-right:3px;">{{ wh.is_active ? 'circle' : 'pause_circle' }}</span>
                    {{ wh.is_active ? 'Active' : 'Paused' }}
                  </span>
                </td>

                <!-- Created -->
                <td style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;white-space:nowrap;">{{ formatDate(wh.created_at) }}</td>

                <!-- Actions -->
                <td>
                  <div v-if="editingId === wh.id" style="display:flex;justify-content:flex-end;gap:8px;">
                    <button class="btn-filled" :disabled="saving" @click="saveEdit(wh.id)">
                      <md-circular-progress v-if="saving" indeterminate style="margin-right:6px;" />
                      Save
                    </button>
                    <button class="btn-outlined" @click="cancelEdit">Cancel</button>
                  </div>
                  <div v-else style="display:flex;justify-content:flex-end;gap:4px;flex-wrap:wrap;">
                    <button class="btn-icon" @click="sendTest(wh.id)" :disabled="testingId === wh.id" title="Send test event">
                      <md-circular-progress v-if="testingId === wh.id" indeterminate />
                      <span v-else class="material-symbols-outlined">send</span>
                    </button>
                    <button class="btn-icon" @click="toggleDeliveries(wh.id)" title="Delivery logs">
                      <span class="material-symbols-outlined">history</span>
                    </button>
                    <button class="btn-icon" @click="startEdit(wh)" title="Edit">
                      <span class="material-symbols-outlined">edit</span>
                    </button>
                    <button class="btn-icon btn-sm btn-danger" @click="promptDelete(wh.id)" title="Delete">
                      <span class="material-symbols-outlined">delete</span>
                    </button>
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
                      <div class="copy-field" style="flex:1;min-width:200px;max-width:480px;">
                        <span class="copy-field__value">
                          <span v-if="revealedSecret[wh.id]">{{ revealedSecret[wh.id] }}</span>
                          <span v-else style="font-family:monospace;">••••••••••••••••</span>
                        </span>
                        <button class="copy-field__btn" @click="toggleSecret(wh.id)">
                          <span class="material-symbols-outlined">{{ revealedSecret[wh.id] ? 'visibility_off' : 'visibility' }}</span>
                          {{ revealedSecret[wh.id] ? 'Hide' : 'Reveal' }}
                        </button>
                        <button v-if="revealedSecret[wh.id]" class="copy-field__btn" @click="copyText(revealedSecret[wh.id])">
                          <span class="material-symbols-outlined">content_copy</span>
                          Copy
                        </button>
                      </div>
                    </div>

                    <!-- Test result banner -->
                    <div v-if="testResults[wh.id]" style="margin-bottom:16px;">
                      <div :style="`padding:10px 14px;border-radius:8px;font-size:0.875rem;background:${testResults[wh.id].success ? '#dcfce7' : '#fee2e2'};color:${testResults[wh.id].success ? '#16a34a' : '#dc2626'};`">
                        <strong>Test delivery:</strong>
                        {{ testResults[wh.id].success ? `HTTP ${testResults[wh.id].response_code}` : `${testResults[wh.id].error_message || 'Failed'}` }}
                        <span style="opacity:0.7;margin-left:8px;">({{ testResults[wh.id].duration_ms }}ms)</span>
                      </div>
                    </div>

                    <!-- Deliveries table -->
                    <div style="font-size:0.875rem;font-weight:600;margin-bottom:8px;">Recent Deliveries</div>
                    <div v-if="deliveriesLoading[wh.id]" style="display:flex;justify-content:center;padding:16px;">
                      <md-circular-progress indeterminate />
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
                              <button class="btn-text"
                                @click="resendDelivery(wh.id, d.id)"
                                :disabled="resendingId === d.id"
                                style="font-size:0.75rem;"
                              >
                                <md-circular-progress v-if="resendingId === d.id" indeterminate />
                                <span v-else>Resend</span>
                              </button>
                            </td>
                          </tr>
                        </tbody>
                      </table>

                      <!-- Pagination -->
                      <div v-if="deliveryTotals[wh.id] > 20" style="display:flex;justify-content:space-between;align-items:center;margin-top:8px;">
                        <span style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;">{{ deliveryTotals[wh.id] }} total</span>
                        <div style="display:flex;gap:4px;">
                          <button class="btn-icon" :disabled="deliveryPages[wh.id] <= 1" @click="loadDeliveries(wh.id, deliveryPages[wh.id] - 1)">
                            <span class="material-symbols-outlined">chevron_left</span>
                          </button>
                          <button class="btn-icon" :disabled="deliveryPages[wh.id] * 20 >= deliveryTotals[wh.id]" @click="loadDeliveries(wh.id, deliveryPages[wh.id] + 1)">
                            <span class="material-symbols-outlined">chevron_right</span>
                          </button>
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

    <!-- Confirm Delete Webhook Dialog -->
    <BaseModal v-model="showDeleteConfirm" size="sm" :persistent="false">
      <template #headline>
        <span class="material-symbols-outlined" style="color:var(--md-sys-color-error)">delete</span>
        Delete Webhook
      </template>
      <p style="color:var(--md-sys-color-on-surface-variant);">Delete this webhook endpoint? This cannot be undone and any subscribed integrations will stop receiving events.</p>
      <template #actions>
        <button class="btn-text" @click="showDeleteConfirm = false">Cancel</button>
        <button class="btn-filled btn-danger"
          @click="confirmDelete"
        >
          Delete
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue';
import { RouterLink } from 'vue-router';
import BaseModal from '@/components/BaseModal.vue';
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

// Confirm delete dialog
const showDeleteConfirm = ref(false);
const webhookIdToDelete = ref<string | null>(null);

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

function promptDelete(id: string) {
  webhookIdToDelete.value = id;
  showDeleteConfirm.value = true;
}

async function confirmDelete() {
  if (!webhookIdToDelete.value) return;
  const id = webhookIdToDelete.value;
  await webhooksApi.delete(id);
  webhooks.value = webhooks.value.filter((w) => w.id !== id);
  delete deliveries.value[id];
  delete deliveriesLoading.value[id];
  delete deliveryTotals.value[id];
  delete deliveryPages.value[id];
  delete testResults.value[id];
  delete revealedSecret.value[id];
  if (openDeliveriesId.value === id) openDeliveriesId.value = null;
  showDeleteConfirm.value = false;
  webhookIdToDelete.value = null;
}

function eventChipClass(event: string): string {
  if (event === 'link.created') return 'event-chip--success';
  if (event === 'link.deleted') return 'event-chip--error';
  if (event === 'link.clicked') return 'event-chip--primary';
  return 'event-chip--neutral';
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

<style scoped lang="scss">
/* page-section (global) handles padding; max-width set via style attribute on root */

/* ── Page Header ──────────────────────────────────────────────────────────── */
.dash-page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;

  &__left {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  &__title {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  &__subtitle {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.875rem;
    margin: 0;
  }

  &__actions {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-shrink: 0;
  }
}

/* ── Cards ────────────────────────────────────────────────────────────────── */
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

/* ── M3 Card header ───────────────────────────────────────────────────────── */
.m3-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  gap: 1rem;
  flex-wrap: wrap;
}

/* ── M3 Table wrapper ─────────────────────────────────────────────────────── */
.m3-table-wrapper {
  overflow-x: auto;
}

/* ── M3 Empty state ───────────────────────────────────────────────────────── */
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

/* ── Warning banner ───────────────────────────────────────────────────────── */
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

/* ── Table ────────────────────────────────────────────────────────────────── */
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

/* ── M3 Badges ────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
}

.m3-badge--neutral {
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  border: 1px solid var(--md-sys-color-outline-variant);
}

.m3-badge--success {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
}

.m3-badge--warning {
  background: rgba(245, 158, 11, 0.12);
  color: #b45309;
}

.m3-badge--error {
  background: rgba(220, 38, 38, 0.12);
  color: #dc2626;
}

/* ── Event chips ──────────────────────────────────────────────────────────── */
.event-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.event-chip {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  font-family: 'SFMono-Regular', Consolas, monospace;

  &--success {
    background: rgba(22, 163, 74, 0.10);
    color: #15803d;
    border: 1px solid rgba(22, 163, 74, 0.20);
  }

  &--error {
    background: rgba(220, 38, 38, 0.10);
    color: #dc2626;
    border: 1px solid rgba(220, 38, 38, 0.20);
  }

  &--primary {
    background: rgba(99, 91, 255, 0.10);
    color: var(--md-sys-color-primary);
    border: 1px solid rgba(99, 91, 255, 0.20);
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Copy field ───────────────────────────────────────────────────────────── */
.copy-field {
  display: flex;
  align-items: center;
  gap: 6px;
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

    .material-symbols-outlined {
      font-size: 16px;
    }
  }
}

/* ── Deliveries panel ─────────────────────────────────────────────────────── */
.deliveries-panel {
  background: var(--md-sys-color-surface-container-low);
  border-top: 1px solid var(--md-sys-color-outline-variant);
  padding: 16px 20px;
}

/* ── Info card ────────────────────────────────────────────────────────────── */
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
