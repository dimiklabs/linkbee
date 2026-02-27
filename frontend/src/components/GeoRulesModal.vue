<template>
  <BaseModal v-model="isOpen" size="lg" @closed="onDialogClosed">
    <template #headline>
      <span class="material-symbols-outlined" style="font-size:20px;vertical-align:middle;margin-right:6px;color:var(--md-sys-color-primary)">public</span>
      Geo Routing
      <span style="font-size:0.85rem;font-weight:400;color:var(--md-sys-color-on-surface-variant);margin-left:8px">{{ slug }}</span>
    </template>

    <div style="min-width:560px;max-width:100%;padding:0 4px;">

      <!-- Enable / Disable toggle -->
      <div class="geo-toggle-row">
        <div class="geo-toggle-info">
          <div class="geo-toggle-title">Geo routing</div>
          <div class="geo-toggle-subtitle">Redirect visitors from specific countries to custom URLs.</div>
        </div>
        <md-switch
          :selected="geoRoutingEnabled"
          @change="geoRoutingEnabled = ($event.target as HTMLInputElement).checked; onToggle()"
          :disabled="toggling"
        />
      </div>

      <!-- Status badge -->
      <div class="geo-status-row">
        <span v-if="geoRoutingEnabled" class="m3-badge m3-badge--success">
          <span class="material-symbols-outlined" style="font-size:14px;vertical-align:middle;margin-right:2px">check_circle</span>
          Geo routing enabled
        </span>
        <span v-else class="m3-badge m3-badge--neutral">
          <span class="material-symbols-outlined" style="font-size:14px;vertical-align:middle;margin-right:2px">remove_circle</span>
          Geo routing disabled
        </span>
      </div>

      <!-- Rules list -->
      <div class="geo-rules-section">
        <div class="geo-rules-header">
          <h6 class="geo-rules-title">Country Rules</h6>
          <span class="geo-rules-count">{{ rules.length }} rule{{ rules.length !== 1 ? 's' : '' }}</span>
        </div>

        <div v-if="loadingRules" class="geo-loading">
          <md-circular-progress indeterminate style="--md-circular-progress-size:24px" />
          <span>Loading rules…</span>
        </div>

        <div v-else-if="rules.length === 0" class="geo-empty">
          <span class="material-symbols-outlined" style="font-size:36px;opacity:0.4">language</span>
          <span>No geo rules yet. Add one below.</span>
        </div>

        <div v-else class="geo-rules-list">
          <div v-for="rule in rules" :key="rule.id" class="geo-rule-item">
            <!-- View mode -->
            <div v-if="editingId !== rule.id" class="geo-rule-view">
              <span class="flag-emoji" :title="countryName(rule.country_code)">
                {{ countryFlag(rule.country_code) }}
              </span>
              <span class="m3-badge m3-badge--secondary" style="min-width:3.5rem;text-align:center">
                {{ rule.country_code }}
              </span>
              <span class="geo-rule-priority">Priority {{ rule.priority }}</span>
              <span class="geo-rule-url" :title="rule.destination_url">
                {{ rule.destination_url }}
              </span>
              <div class="geo-rule-actions">
                <md-icon-button @click="startEdit(rule)" title="Edit" style="width:32px;height:32px">
                  <span class="material-symbols-outlined" style="font-size:16px">edit</span>
                </md-icon-button>
                <md-icon-button @click="deleteRule(rule.id)" title="Delete" style="width:32px;height:32px;--md-icon-button-icon-color:var(--md-sys-color-error)">
                  <span class="material-symbols-outlined" style="font-size:16px">delete</span>
                </md-icon-button>
              </div>
            </div>

            <!-- Edit mode -->
            <div v-else class="geo-rule-edit">
              <AppSelect
                v-model="editForm.country_code"
                label="Country"
                style="min-width:160px"
              >
                <option v-for="c in countries" :key="c.code" :value="c.code">{{ countryFlag(c.code) }} {{ c.name }} ({{ c.code }})</option>
              </AppSelect>
              <md-outlined-text-field
                :value="editForm.destination_url"
                @input="editForm.destination_url = ($event.target as HTMLInputElement).value"
                label="Destination URL"
                style="flex:1"
              />
              <md-outlined-text-field
                :value="String(editForm.priority)"
                @input="editForm.priority = Number(($event.target as HTMLInputElement).value)"
                label="Priority"
                type="number"
                style="width:90px"
              />
              <div class="geo-rule-edit-actions">
                <md-icon-button @click="saveEdit(rule.id)" :disabled="saving" style="width:32px;height:32px;--md-icon-button-icon-color:#1AA563">
                  <span class="material-symbols-outlined" style="font-size:18px">check</span>
                </md-icon-button>
                <md-icon-button @click="cancelEdit" style="width:32px;height:32px">
                  <span class="material-symbols-outlined" style="font-size:18px">close</span>
                </md-icon-button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Add new rule form -->
      <div class="geo-add-section">
        <div class="geo-add-title">
          <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-primary)">add_circle</span>
          Add Rule
        </div>
        <div class="geo-add-form">
          <AppSelect
            v-model="newForm.country_code"
            label="Select country…"
            style="min-width:180px"
          >
            <option value="" disabled>Select country…</option>
            <option v-for="c in countries" :key="c.code" :value="c.code">{{ countryFlag(c.code) }} {{ c.name }} ({{ c.code }})</option>
          </AppSelect>
          <md-outlined-text-field
            :value="newForm.destination_url"
            @input="newForm.destination_url = ($event.target as HTMLInputElement).value"
            label="Destination URL"
            type="url"
            placeholder="https://destination.example.com"
            style="flex:1"
          />
          <md-outlined-text-field
            :value="String(newForm.priority)"
            @input="newForm.priority = Number(($event.target as HTMLInputElement).value)"
            label="Priority"
            type="number"
            style="width:100px"
          />
          <md-filled-button
            @click="addRule"
            :disabled="!newForm.country_code || !newForm.destination_url || saving"
          >
            <md-circular-progress v-if="saving" indeterminate style="--md-circular-progress-size:16px" slot="icon" />
            Add
          </md-filled-button>
        </div>
        <div v-if="formError" class="geo-form-error">
          <span class="material-symbols-outlined" style="font-size:16px">error</span>
          {{ formError }}
        </div>
      </div>

    </div>

    <template #actions>
      <md-text-button @click="hide">Close</md-text-button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import geoApi from '@/api/geo';
import type { LinkGeoRule } from '@/types/links';
import BaseModal from '@/components/BaseModal.vue';
import AppSelect from '@/components/AppSelect.vue';

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
let bsModal: any = null;

const isOpen = ref(false);
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

// Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
onMounted(async () => {
  // noop
});

onUnmounted(() => {
  bsModal?.dispose();
});

watch(() => props.isGeoRoutingInitial, (v) => {
  geoRoutingEnabled.value = v;
});

function show() {
  isOpen.value = true;
  loadRules();
  bsModal?.show();
}

function hide() {
  isOpen.value = false;
  bsModal?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

defineExpose({
  show,
  hide,
});
</script>

<style scoped>
/* Toggle row */
.geo-toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 12px;
  margin-bottom: 16px;
}

.geo-toggle-title {
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface);
}

.geo-toggle-subtitle {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
}

/* Status badge row */
.geo-status-row {
  margin-bottom: 20px;
}

/* Rules section */
.geo-rules-section {
  margin-bottom: 20px;
}

.geo-rules-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.geo-rules-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0;
}

.geo-rules-count {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.geo-loading,
.geo-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  text-align: center;
}

.geo-rules-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 12px;
}

.geo-rule-item {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  overflow: hidden;
}

.geo-rule-view {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
}

.flag-emoji {
  font-size: 1.3rem;
  line-height: 1;
  flex-shrink: 0;
}

.geo-rule-priority {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  flex-shrink: 0;
}

.geo-rule-url {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.82rem;
  color: var(--md-sys-color-primary);
  min-width: 0;
}

.geo-rule-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.geo-rule-edit {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  flex-wrap: wrap;
}

.geo-rule-edit-actions {
  display: flex;
  gap: 4px;
}

/* Add section */
.geo-add-section {
  background: var(--md-sys-color-surface-container-low);
  border-radius: 12px;
  padding: 16px;
}

.geo-add-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 12px;
}

.geo-add-form {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  flex-wrap: wrap;
}

.geo-form-error {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 0.82rem;
  color: var(--md-sys-color-error);
}
</style>
