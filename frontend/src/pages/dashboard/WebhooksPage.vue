<template>
  <div class="container-fluid py-4">
    <!-- Header -->
    <div class="d-flex align-items-center justify-content-between mb-4">
      <div>
        <h4 class="fw-bold mb-1">Webhooks</h4>
        <p class="text-muted mb-0 small">Receive real-time HTTP notifications when events happen.</p>
      </div>
      <button class="btn btn-primary" @click="showCreate = !showCreate">
        <i class="bi bi-plus-lg me-1"></i>Add Webhook
      </button>
    </div>

    <!-- Create form -->
    <div v-if="showCreate" class="card border-0 shadow-sm mb-4">
      <div class="card-body">
        <h6 class="fw-semibold mb-3">New Webhook</h6>
        <div class="mb-3">
          <label class="form-label small fw-medium">Endpoint URL</label>
          <input
            v-model="newURL"
            type="url"
            class="form-control"
            placeholder="https://your-server.com/webhook"
            ref="urlInputRef"
          />
        </div>
        <div class="mb-3">
          <label class="form-label small fw-medium">Events to subscribe</label>
          <div class="d-flex flex-column gap-2">
            <div v-for="ev in WEBHOOK_EVENTS" :key="ev.value" class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                :id="`new-ev-${ev.value}`"
                :value="ev.value"
                v-model="newEvents"
              />
              <label class="form-check-label" :for="`new-ev-${ev.value}`">
                <span class="fw-medium">{{ ev.label }}</span>
                <span class="text-muted ms-1 small">— {{ ev.description }}</span>
              </label>
            </div>
          </div>
          <div v-if="newEvents.length === 0" class="text-danger small mt-1">
            Select at least one event.
          </div>
        </div>
        <div class="d-flex gap-2">
          <button class="btn btn-primary btn-sm" :disabled="creating || !newURL || newEvents.length === 0" @click="createWebhook">
            <span v-if="creating" class="spinner-border spinner-border-sm me-1"></span>
            Create
          </button>
          <button class="btn btn-outline-secondary btn-sm" @click="cancelCreate">Cancel</button>
        </div>
        <div v-if="createError" class="alert alert-danger mt-3 py-2 small mb-0">{{ createError }}</div>
      </div>
    </div>

    <!-- Webhook list -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary"></div>
    </div>

    <div v-else-if="webhooks.length === 0 && !showCreate" class="card border-0 shadow-sm">
      <div class="card-body text-center py-5">
        <div class="mb-3" style="font-size:2.5rem">🔔</div>
        <h6 class="fw-semibold">No webhooks yet</h6>
        <p class="text-muted small mb-3">Add a webhook endpoint to start receiving event notifications.</p>
        <button class="btn btn-primary btn-sm" @click="showCreate = true">Add Webhook</button>
      </div>
    </div>

    <div v-else-if="webhooks.length > 0" class="card border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th class="ps-3">Endpoint</th>
              <th>Events</th>
              <th>Status</th>
              <th>Created</th>
              <th class="pe-3 text-end">Actions</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="wh in webhooks" :key="wh.id">
              <tr>
                <!-- Endpoint -->
                <td class="ps-3">
                  <div v-if="editingId === wh.id">
                    <input v-model="editURL" type="url" class="form-control form-control-sm" style="min-width:240px" />
                  </div>
                  <code v-else class="small">{{ wh.url }}</code>
                </td>

                <!-- Events -->
                <td>
                  <div v-if="editingId === wh.id" class="d-flex flex-column gap-1">
                    <div v-for="ev in WEBHOOK_EVENTS" :key="ev.value" class="form-check mb-0">
                      <input
                        class="form-check-input"
                        type="checkbox"
                        :id="`edit-ev-${wh.id}-${ev.value}`"
                        :value="ev.value"
                        v-model="editEvents"
                      />
                      <label class="form-check-label small" :for="`edit-ev-${wh.id}-${ev.value}`">{{ ev.label }}</label>
                    </div>
                  </div>
                  <div v-else class="d-flex flex-wrap gap-1">
                    <span
                      v-for="ev in wh.events"
                      :key="ev"
                      class="badge rounded-pill"
                      :class="eventBadgeClass(ev)"
                    >{{ ev }}</span>
                  </div>
                </td>

                <!-- Status -->
                <td>
                  <div v-if="editingId === wh.id" class="form-check form-switch mb-0">
                    <input class="form-check-input" type="checkbox" v-model="editIsActive" />
                    <label class="form-check-label small">Active</label>
                  </div>
                  <span v-else :class="wh.is_active ? 'badge bg-success' : 'badge bg-secondary'">
                    {{ wh.is_active ? 'Active' : 'Paused' }}
                  </span>
                </td>

                <!-- Created -->
                <td class="small text-muted">{{ formatDate(wh.created_at) }}</td>

                <!-- Actions -->
                <td class="pe-3 text-end">
                  <div v-if="editingId === wh.id" class="d-flex justify-content-end gap-2">
                    <button class="btn btn-primary btn-sm" :disabled="saving" @click="saveEdit(wh.id)">
                      <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
                      Save
                    </button>
                    <button class="btn btn-outline-secondary btn-sm" @click="cancelEdit">Cancel</button>
                  </div>
                  <div v-else class="d-flex justify-content-end gap-1 flex-wrap">
                    <button class="btn btn-outline-secondary btn-sm" @click="sendTest(wh.id)" :disabled="testingId === wh.id" title="Send test event">
                      <span v-if="testingId === wh.id" class="spinner-border spinner-border-sm"></span>
                      <span v-else>Test</span>
                    </button>
                    <button class="btn btn-outline-secondary btn-sm" @click="toggleDeliveries(wh.id)" title="Delivery logs">
                      Logs
                    </button>
                    <button class="btn btn-outline-secondary btn-sm" @click="startEdit(wh)" title="Edit">
                      <i class="bi bi-pencil"></i>
                    </button>
                    <button class="btn btn-outline-danger btn-sm" @click="deleteWebhook(wh.id)" title="Delete">
                      <i class="bi bi-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>

              <!-- Delivery logs expandable row -->
              <tr v-if="openDeliveriesId === wh.id">
                <td colspan="5" class="p-0">
                  <div class="bg-light border-top px-3 py-3">

                    <!-- Secret reveal -->
                    <div class="d-flex align-items-center gap-2 mb-3">
                      <span class="small fw-medium text-muted">Signing secret:</span>
                      <code v-if="revealedSecret[wh.id]" class="small">{{ revealedSecret[wh.id] }}</code>
                      <span v-else class="small text-muted font-monospace">••••••••••••••••</span>
                      <button class="btn btn-outline-secondary btn-sm py-0" style="font-size:0.7rem" @click="toggleSecret(wh.id)">
                        {{ revealedSecret[wh.id] ? 'Hide' : 'Reveal' }}
                      </button>
                      <button v-if="revealedSecret[wh.id]" class="btn btn-outline-secondary btn-sm py-0" style="font-size:0.7rem" @click="copyText(revealedSecret[wh.id])">
                        Copy
                      </button>
                    </div>

                    <!-- Test result banner -->
                    <div v-if="testResults[wh.id]" class="mb-3">
                      <div :class="['alert', 'py-2', 'small', 'mb-0', testResults[wh.id].success ? 'alert-success' : 'alert-danger']">
                        <strong>Test delivery:</strong>
                        {{ testResults[wh.id].success ? `✓ HTTP ${testResults[wh.id].response_code}` : `✗ ${testResults[wh.id].error_message || 'Failed'}` }}
                        <span class="text-muted ms-2">({{ testResults[wh.id].duration_ms }}ms)</span>
                      </div>
                    </div>

                    <!-- Deliveries table -->
                    <h6 class="small fw-semibold mb-2">Recent Deliveries</h6>
                    <div v-if="deliveriesLoading[wh.id]" class="text-center py-3">
                      <div class="spinner-border spinner-border-sm text-primary"></div>
                    </div>
                    <div v-else-if="!deliveries[wh.id] || deliveries[wh.id].length === 0" class="text-muted small py-2">
                      No deliveries recorded yet. Send a test event to get started.
                    </div>
                    <div v-else class="table-responsive">
                      <table class="table table-sm table-borderless mb-0 align-middle" style="font-size:0.8rem">
                        <thead>
                          <tr class="text-muted">
                            <th>Time</th>
                            <th>Event</th>
                            <th>Status</th>
                            <th>Duration</th>
                            <th></th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr v-for="d in deliveries[wh.id]" :key="d.id">
                            <td class="text-muted text-nowrap">{{ formatDateTime(d.created_at) }}</td>
                            <td><code>{{ d.event }}</code></td>
                            <td>
                              <span v-if="d.success" class="badge text-bg-success">{{ d.response_code }}</span>
                              <span v-else class="badge text-bg-danger">{{ d.response_code || 'Error' }}</span>
                            </td>
                            <td class="text-muted">{{ d.duration_ms }}ms</td>
                            <td class="text-end">
                              <button
                                class="btn btn-outline-secondary btn-sm py-0"
                                style="font-size:0.7rem"
                                @click="resendDelivery(wh.id, d.id)"
                                :disabled="resendingId === d.id"
                              >
                                <span v-if="resendingId === d.id" class="spinner-border spinner-border-sm"></span>
                                <span v-else>Resend</span>
                              </button>
                            </td>
                          </tr>
                        </tbody>
                      </table>

                      <!-- Pagination -->
                      <div v-if="deliveryTotals[wh.id] > 20" class="d-flex justify-content-between align-items-center mt-2">
                        <span class="text-muted small">{{ deliveryTotals[wh.id] }} total</span>
                        <div class="d-flex gap-1">
                          <button class="btn btn-sm btn-outline-secondary" :disabled="deliveryPages[wh.id] <= 1" @click="loadDeliveries(wh.id, deliveryPages[wh.id] - 1)">←</button>
                          <button class="btn btn-sm btn-outline-secondary" :disabled="deliveryPages[wh.id] * 20 >= deliveryTotals[wh.id]" @click="loadDeliveries(wh.id, deliveryPages[wh.id] + 1)">→</button>
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
    <div class="card border-0 shadow-sm mt-4">
      <div class="card-body">
        <h6 class="fw-semibold mb-2">Verifying webhook signatures</h6>
        <p class="small text-muted mb-2">
          Every delivery includes an <code>X-Webhook-Signature</code> header with an HMAC-SHA256 signature of the raw request body, prefixed with <code>sha256=</code>.
          Your signing secret was generated when you created the webhook (stored server-side only).
          Verify on your server to ensure the request is genuine:
        </p>
        <pre class="bg-light rounded p-3 small mb-0"><code>import hmac, hashlib

def verify(secret: str, body: bytes, header: str) -> bool:
    expected = "sha256=" + hmac.new(
        secret.encode(), body, hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, header)</code></pre>
      </div>
    </div>

    <!-- Example payload card -->
    <div class="card border-0 shadow-sm mt-3">
      <div class="card-body">
        <h6 class="fw-semibold mb-2">Example payload</h6>
        <pre class="bg-light rounded p-3 small mb-0"><code>{
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue';
import webhooksApi from '@/api/webhooks';
import { WEBHOOK_EVENTS } from '@/types/webhooks';
import type { Webhook, WebhookDelivery } from '@/types/webhooks';

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
  if (event === 'link.created') return 'bg-success bg-opacity-10 text-success';
  if (event === 'link.deleted') return 'bg-danger bg-opacity-10 text-danger';
  if (event === 'link.clicked') return 'bg-primary bg-opacity-10 text-primary';
  return 'bg-secondary bg-opacity-10 text-secondary';
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

onMounted(fetchWebhooks);
</script>
