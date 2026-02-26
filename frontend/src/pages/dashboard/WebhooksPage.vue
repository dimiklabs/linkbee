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
            <tr v-for="wh in webhooks" :key="wh.id">
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
                <div v-else class="d-flex justify-content-end gap-2">
                  <button class="btn btn-outline-secondary btn-sm" @click="startEdit(wh)" title="Edit">
                    <i class="bi bi-pencil"></i>
                  </button>
                  <button class="btn btn-outline-danger btn-sm" @click="deleteWebhook(wh.id)" title="Delete">
                    <i class="bi bi-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
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
import { ref, onMounted, nextTick } from 'vue';
import webhooksApi from '@/api/webhooks';
import { WEBHOOK_EVENTS } from '@/types/webhooks';
import type { Webhook } from '@/types/webhooks';

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

// Focus URL input when create form opens
import { watch } from 'vue';
watch(showCreate, async (val) => {
  if (val) {
    await nextTick();
    urlInputRef.value?.focus();
  }
});

onMounted(fetchWebhooks);
</script>
