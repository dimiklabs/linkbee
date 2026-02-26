<template>
  <div class="modal fade" id="geoRulesModal" tabindex="-1" aria-labelledby="geoRulesModalLabel" aria-hidden="true" ref="modalEl">
    <div class="modal-dialog modal-lg modal-dialog-scrollable">
      <div class="modal-content">

        <!-- Header -->
        <div class="modal-header">
          <h5 class="modal-title" id="geoRulesModalLabel">
            <i class="bi bi-geo-alt me-2 text-primary"></i>Geo Routing
            <span class="text-muted fw-normal fs-6 ms-2">{{ slug }}</span>
          </h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>

        <!-- Body -->
        <div class="modal-body">

          <!-- Enable / Disable toggle -->
          <div class="d-flex align-items-center justify-content-between mb-4 p-3 bg-light rounded">
            <div>
              <div class="fw-semibold">Geo routing</div>
              <div class="text-muted small">Redirect visitors from specific countries to custom URLs.</div>
            </div>
            <div class="form-check form-switch mb-0">
              <input
                class="form-check-input"
                type="checkbox"
                role="switch"
                id="geoRoutingToggle"
                v-model="geoRoutingEnabled"
                :disabled="toggling"
                @change="onToggle"
                style="width:2.5em;height:1.4em;"
              />
            </div>
          </div>

          <!-- Status badge -->
          <div class="mb-3">
            <span v-if="geoRoutingEnabled" class="badge bg-success-subtle text-success border border-success-subtle px-3 py-2">
              <i class="bi bi-check-circle me-1"></i>Geo routing enabled
            </span>
            <span v-else class="badge bg-secondary-subtle text-secondary border border-secondary-subtle px-3 py-2">
              <i class="bi bi-dash-circle me-1"></i>Geo routing disabled
            </span>
          </div>

          <!-- Rules list -->
          <div class="mb-3">
            <div class="d-flex align-items-center justify-content-between mb-2">
              <h6 class="mb-0">Country Rules</h6>
              <span class="text-muted small">{{ rules.length }} rule{{ rules.length !== 1 ? 's' : '' }}</span>
            </div>

            <div v-if="loadingRules" class="text-center py-4 text-muted">
              <div class="spinner-border spinner-border-sm me-2"></div>Loading rules…
            </div>

            <div v-else-if="rules.length === 0" class="text-center py-4 text-muted">
              <i class="bi bi-globe2 fs-2 d-block mb-2 opacity-50"></i>
              No geo rules yet. Add one below.
            </div>

            <div v-else class="list-group mb-3">
              <div
                v-for="rule in rules"
                :key="rule.id"
                class="list-group-item list-group-item-action py-3"
              >
                <!-- View mode -->
                <div v-if="editingId !== rule.id" class="d-flex align-items-center gap-2">
                  <span class="flag-emoji me-1" :title="countryName(rule.country_code)">
                    {{ countryFlag(rule.country_code) }}
                  </span>
                  <span class="badge bg-secondary-subtle text-secondary fw-bold" style="min-width:3.5rem">
                    {{ rule.country_code }}
                  </span>
                  <span class="text-muted small me-1">Priority {{ rule.priority }}</span>
                  <span class="text-truncate flex-grow-1 text-primary small" :title="rule.destination_url">
                    {{ rule.destination_url }}
                  </span>
                  <div class="d-flex gap-1 ms-auto flex-shrink-0">
                    <button class="btn btn-sm btn-outline-secondary" @click="startEdit(rule)" title="Edit">
                      <i class="bi bi-pencil"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger" @click="deleteRule(rule.id)" title="Delete">
                      <i class="bi bi-trash"></i>
                    </button>
                  </div>
                </div>

                <!-- Edit mode -->
                <div v-else class="row g-2">
                  <div class="col-auto">
                    <select class="form-select form-select-sm" v-model="editForm.country_code" style="width:10rem">
                      <option v-for="c in countries" :key="c.code" :value="c.code">
                        {{ countryFlag(c.code) }} {{ c.name }} ({{ c.code }})
                      </option>
                    </select>
                  </div>
                  <div class="col">
                    <input class="form-control form-control-sm" v-model="editForm.destination_url" placeholder="Destination URL" />
                  </div>
                  <div class="col-auto" style="width:6rem">
                    <input class="form-control form-control-sm" type="number" v-model.number="editForm.priority" placeholder="Priority" min="0" />
                  </div>
                  <div class="col-auto d-flex gap-1">
                    <button class="btn btn-sm btn-success" @click="saveEdit(rule.id)" :disabled="saving">
                      <i class="bi bi-check-lg"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-secondary" @click="cancelEdit">
                      <i class="bi bi-x-lg"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Add new rule form -->
          <div class="border rounded p-3 bg-light">
            <h6 class="mb-3"><i class="bi bi-plus-circle me-1 text-primary"></i>Add Rule</h6>
            <div class="row g-2">
              <div class="col-12 col-sm-auto">
                <select class="form-select form-select-sm" v-model="newForm.country_code" style="min-width:12rem">
                  <option value="" disabled>Select country…</option>
                  <option v-for="c in countries" :key="c.code" :value="c.code">
                    {{ countryFlag(c.code) }} {{ c.name }} ({{ c.code }})
                  </option>
                </select>
              </div>
              <div class="col">
                <input
                  class="form-control form-control-sm"
                  v-model="newForm.destination_url"
                  placeholder="https://destination.example.com"
                  type="url"
                />
              </div>
              <div class="col-auto" style="width:7rem">
                <input
                  class="form-control form-control-sm"
                  type="number"
                  v-model.number="newForm.priority"
                  placeholder="Priority"
                  min="0"
                />
              </div>
              <div class="col-auto">
                <button
                  class="btn btn-sm btn-primary"
                  @click="addRule"
                  :disabled="!newForm.country_code || !newForm.destination_url || saving"
                >
                  <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
                  Add
                </button>
              </div>
            </div>
            <div v-if="formError" class="text-danger small mt-2">{{ formError }}</div>
          </div>

        </div><!-- /modal-body -->

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import type { Modal as BsModal } from 'bootstrap';
import geoApi from '@/api/geo';
import type { LinkGeoRule } from '@/types/links';

const props = defineProps<{
  linkId: string;
  slug: string;
  isGeoRoutingInitial: boolean;
}>();

const emit = defineEmits<{
  (e: 'updated', isGeoRouting: boolean): void;
}>();

// ── Country data (top 50 countries + comprehensive list) ──────────────────────
const countries = [
  { code: 'US', name: 'United States' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'CA', name: 'Canada' },
  { code: 'AU', name: 'Australia' },
  { code: 'DE', name: 'Germany' },
  { code: 'FR', name: 'France' },
  { code: 'JP', name: 'Japan' },
  { code: 'IN', name: 'India' },
  { code: 'BR', name: 'Brazil' },
  { code: 'MX', name: 'Mexico' },
  { code: 'NL', name: 'Netherlands' },
  { code: 'SE', name: 'Sweden' },
  { code: 'NO', name: 'Norway' },
  { code: 'DK', name: 'Denmark' },
  { code: 'FI', name: 'Finland' },
  { code: 'CH', name: 'Switzerland' },
  { code: 'AT', name: 'Austria' },
  { code: 'BE', name: 'Belgium' },
  { code: 'PL', name: 'Poland' },
  { code: 'ES', name: 'Spain' },
  { code: 'IT', name: 'Italy' },
  { code: 'PT', name: 'Portugal' },
  { code: 'RU', name: 'Russia' },
  { code: 'CN', name: 'China' },
  { code: 'KR', name: 'South Korea' },
  { code: 'SG', name: 'Singapore' },
  { code: 'NZ', name: 'New Zealand' },
  { code: 'ZA', name: 'South Africa' },
  { code: 'NG', name: 'Nigeria' },
  { code: 'EG', name: 'Egypt' },
  { code: 'SA', name: 'Saudi Arabia' },
  { code: 'AE', name: 'United Arab Emirates' },
  { code: 'TR', name: 'Turkey' },
  { code: 'AR', name: 'Argentina' },
  { code: 'CL', name: 'Chile' },
  { code: 'CO', name: 'Colombia' },
  { code: 'PH', name: 'Philippines' },
  { code: 'ID', name: 'Indonesia' },
  { code: 'MY', name: 'Malaysia' },
  { code: 'TH', name: 'Thailand' },
  { code: 'VN', name: 'Vietnam' },
  { code: 'PK', name: 'Pakistan' },
  { code: 'BD', name: 'Bangladesh' },
  { code: 'UA', name: 'Ukraine' },
  { code: 'CZ', name: 'Czech Republic' },
  { code: 'HU', name: 'Hungary' },
  { code: 'RO', name: 'Romania' },
  { code: 'GR', name: 'Greece' },
  { code: 'HK', name: 'Hong Kong' },
  { code: 'TW', name: 'Taiwan' },
];

// ── State ──────────────────────────────────────────────────────────────────────
const modalEl = ref<HTMLElement | null>(null);
let bsModal: BsModal | null = null;

const rules = ref<LinkGeoRule[]>([]);
const loadingRules = ref(false);
const saving = ref(false);
const toggling = ref(false);
const geoRoutingEnabled = ref(props.isGeoRoutingInitial);
const formError = ref('');

const editingId = ref<string | null>(null);
const editForm = ref({ country_code: '', destination_url: '', priority: 0 });

const newForm = ref({ country_code: '', destination_url: '', priority: 0 });

// ── Helpers ────────────────────────────────────────────────────────────────────
function countryFlag(code: string): string {
  // Regional indicator symbol letters for emoji flags
  return code
    .toUpperCase()
    .split('')
    .map(c => String.fromCodePoint(0x1F1E0 + c.charCodeAt(0) - 65))
    .join('');
}

function countryName(code: string): string {
  return countries.find(c => c.code === code)?.name ?? code;
}

// ── Data loading ───────────────────────────────────────────────────────────────
async function loadRules() {
  loadingRules.value = true;
  try {
    const res = await geoApi.list(props.linkId);
    rules.value = res.data ?? [];
  } catch {
    rules.value = [];
  } finally {
    loadingRules.value = false;
  }
}

// ── Toggle ─────────────────────────────────────────────────────────────────────
async function onToggle() {
  toggling.value = true;
  try {
    await geoApi.toggleGeoRouting(props.linkId, geoRoutingEnabled.value);
    emit('updated', geoRoutingEnabled.value);
  } catch {
    // revert on failure
    geoRoutingEnabled.value = !geoRoutingEnabled.value;
  } finally {
    toggling.value = false;
  }
}

// ── Add rule ───────────────────────────────────────────────────────────────────
async function addRule() {
  formError.value = '';
  if (!newForm.value.country_code || !newForm.value.destination_url) return;
  saving.value = true;
  try {
    const res = await geoApi.create(props.linkId, {
      country_code: newForm.value.country_code,
      destination_url: newForm.value.destination_url,
      priority: newForm.value.priority,
    });
    rules.value.push(res.data);
    newForm.value = { country_code: '', destination_url: '', priority: 0 };
  } catch (err: any) {
    formError.value = err?.response?.data?.description ?? 'Failed to add rule';
  } finally {
    saving.value = false;
  }
}

// ── Edit rule ──────────────────────────────────────────────────────────────────
function startEdit(rule: LinkGeoRule) {
  editingId.value = rule.id;
  editForm.value = {
    country_code: rule.country_code,
    destination_url: rule.destination_url,
    priority: rule.priority,
  };
}

function cancelEdit() {
  editingId.value = null;
}

async function saveEdit(ruleId: string) {
  saving.value = true;
  try {
    const res = await geoApi.update(props.linkId, ruleId, {
      country_code: editForm.value.country_code,
      destination_url: editForm.value.destination_url,
      priority: editForm.value.priority,
    });
    const idx = rules.value.findIndex(r => r.id === ruleId);
    if (idx !== -1) rules.value[idx] = res.data;
    editingId.value = null;
  } catch {
    // ignore — keep editing
  } finally {
    saving.value = false;
  }
}

// ── Delete rule ────────────────────────────────────────────────────────────────
async function deleteRule(ruleId: string) {
  if (!confirm('Delete this geo rule?')) return;
  try {
    await geoApi.delete(props.linkId, ruleId);
    rules.value = rules.value.filter(r => r.id !== ruleId);
  } catch {
    // ignore
  }
}

// ── Bootstrap modal lifecycle ──────────────────────────────────────────────────
onMounted(async () => {
  const { Modal } = await import('bootstrap');
  if (modalEl.value) {
    bsModal = new Modal(modalEl.value);
    modalEl.value.addEventListener('hidden.bs.modal', () => {
      editingId.value = null;
      formError.value = '';
    });
    modalEl.value.addEventListener('shown.bs.modal', loadRules);
  }
});

onUnmounted(() => {
  bsModal?.dispose();
});

watch(() => props.isGeoRoutingInitial, (v) => {
  geoRoutingEnabled.value = v;
});

defineExpose({
  show() { bsModal?.show(); },
  hide() { bsModal?.hide(); },
});
</script>

<style scoped>
.flag-emoji {
  font-size: 1.3rem;
  line-height: 1;
}
</style>
