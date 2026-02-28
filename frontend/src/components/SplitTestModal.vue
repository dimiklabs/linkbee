<template>
  <BaseModal v-model="isOpen" size="lg" @closed="onDialogClosed">
    <template #headline>
      A/B Split Test
      <span style="font-size:0.8rem;font-weight:400;color:var(--md-sys-color-on-surface-variant);display:block;margin-top:2px">{{ props.slug }}</span>
    </template>

    <div style="min-width:560px;max-width:100%;padding:0 4px;">
      <!-- Error -->
      <div v-if="error" class="split-error-banner">
        <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-error)">error</span>
        <span style="flex:1;font-size:0.875rem">{{ error }}</span>
        <button class="btn-icon" @click="error = ''" style="width:32px;height:32px">
          <span class="material-symbols-outlined" style="font-size:18px">close</span>
        </button>
      </div>

      <!-- Enable toggle -->
      <div class="split-toggle-row">
        <div>
          <div class="split-toggle-title">Split Test {{ isSplitTest ? 'Enabled' : 'Disabled' }}</div>
          <div class="split-toggle-subtitle">
            {{ isSplitTest ? 'Traffic is being distributed across variants.' : 'Enable to start routing traffic to variants.' }}
          </div>
        </div>
        <md-switch
          :selected="isSplitTest"
          :disabled="toggling || loading"
          @change="toggleSplitTest($event)"
        />
      </div>

      <!-- Weight summary bar -->
      <div v-if="variants.length > 0" class="split-distribution">
        <div class="split-distribution-header">
          <span class="split-distribution-label">Traffic distribution</span>
          <span class="split-distribution-total">Total weight: {{ totalWeight }}</span>
        </div>
        <div class="split-bar">
          <div
            v-for="(v, i) in variants"
            :key="v.id"
            class="split-bar-segment"
            :style="{ width: variantPercent(v.weight) + '%', backgroundColor: variantColor(i) }"
            :title="`${v.name}: ${variantPercent(v.weight).toFixed(1)}%`"
          ></div>
        </div>
        <div class="split-legend">
          <span v-for="(v, i) in variants" :key="v.id" class="split-legend-item">
            <span class="split-legend-dot" :style="{ backgroundColor: variantColor(i) }"></span>
            {{ v.name }}: {{ variantPercent(v.weight).toFixed(1) }}%
          </span>
        </div>
      </div>

      <!-- Variants list loading -->
      <div v-if="loading" class="split-loading">
        <md-circular-progress indeterminate />
      </div>

      <div v-else>
        <div v-if="variants.length === 0" class="split-empty">
          No variants yet. Add one below to get started.
        </div>

        <div v-for="(v, i) in variants" :key="v.id" class="variant-card">
          <!-- Edit mode -->
          <div v-if="editingId === v.id" class="variant-edit-form">
            <md-outlined-text-field
              :value="editForm.name"
              @input="editForm.name = ($event.target as HTMLInputElement).value"
              label="Variant name"
              style="flex:2"
            />
            <md-outlined-text-field
              :value="editForm.destination_url"
              @input="editForm.destination_url = ($event.target as HTMLInputElement).value"
              label="Destination URL"
              style="flex:3"
            />
            <md-outlined-text-field
              :value="String(editForm.weight)"
              @input="editForm.weight = Number(($event.target as HTMLInputElement).value)"
              label="Weight"
              type="number"
              style="width:80px"
            />
            <button class="btn-icon" @click="submitEdit(v)" :disabled="saving" style="width:36px;height:36px">
              <span class="material-symbols-outlined">check</span>
            </button>
            <button class="btn-icon" @click="cancelEdit" style="width:36px;height:36px">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <!-- View mode -->
          <div v-else class="variant-view">
            <span class="variant-dot" :style="{ backgroundColor: variantColor(i) }"></span>
            <div class="variant-info">
              <div class="variant-name-row">
                <span class="variant-name">{{ v.name }}</span>
                <span class="m3-badge m3-badge--neutral" style="font-size:0.72rem">weight {{ v.weight }}</span>
                <span class="m3-badge m3-badge--secondary" style="font-size:0.72rem">{{ v.click_count.toLocaleString() }} clicks</span>
              </div>
              <div class="variant-url" :title="v.destination_url">{{ v.destination_url }}</div>
            </div>
            <div class="variant-actions">
              <button class="btn-icon" @click="startEdit(v)" title="Edit" style="width:32px;height:32px">
                <span class="material-symbols-outlined" style="font-size:15px">edit</span>
              </button>
              <button class="btn-icon btn-sm btn-danger" @click="deleteVariant(v)" title="Delete" style="width:32px;height:32px;">
                <span class="material-symbols-outlined" style="font-size:15px">delete</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Add variant form -->
        <div v-if="showAddForm" class="add-variant-card">
          <p class="add-variant-title">New Variant</p>
          <div class="variant-edit-form">
            <md-outlined-text-field
              :value="newForm.name"
              @input="newForm.name = ($event.target as HTMLInputElement).value"
              label="Variant name"
              placeholder="Variant B"
              style="flex:2"
            />
            <md-outlined-text-field
              :value="newForm.destination_url"
              @input="newForm.destination_url = ($event.target as HTMLInputElement).value"
              label="Destination URL"
              placeholder="https://destination.com"
              style="flex:3"
            />
            <md-outlined-text-field
              :value="String(newForm.weight)"
              @input="newForm.weight = Number(($event.target as HTMLInputElement).value)"
              label="Weight"
              type="number"
              style="width:80px"
            />
            <button class="btn-icon" @click="submitAdd" :disabled="saving" style="width:36px;height:36px">
              <span class="material-symbols-outlined">check</span>
            </button>
            <button class="btn-icon" @click="showAddForm = false" style="width:36px;height:36px">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>
          <div class="add-variant-hint">
            Weight is relative (e.g., 50 + 50 = 50/50 split; 70 + 30 = 70/30 split).
          </div>
        </div>

        <button class="btn-text"
          v-if="!showAddForm && editingId === ''"
          @click="openAddForm"
          style="margin-top:12px"
        >
          <span class="material-symbols-outlined">add</span>
          Add Variant
        </button>
      </div>
    </div>

    <template #actions>
      <button class="btn-text" @click="hide">Close</button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import variantsApi from '@/api/variants';
import type { LinkVariant } from '@/types/links';
import BaseModal from '@/components/BaseModal.vue';

interface Props {
  linkId: string;
  slug: string;
  isSplitTestInitial: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  updated: [isSplitTest: boolean];
}>();

const modalEl = ref<HTMLElement | null>(null);
let modalInstance: any = null;

const isOpen = ref(false);
const variants = ref<LinkVariant[]>([]);
const isSplitTest = ref(props.isSplitTestInitial);
const loading = ref(false);
const saving = ref(false);
const toggling = ref(false);
const error = ref('');
const editingId = ref('');
const showAddForm = ref(false);

interface VariantForm {
  name: string;
  destination_url: string;
  weight: number;
}

const editForm = ref<VariantForm>({ name: '', destination_url: '', weight: 50 });
const newForm = ref<VariantForm>({ name: '', destination_url: '', weight: 50 });

const VARIANT_COLORS = ['#635bff', '#0d6efd', '#198754', '#dc3545', '#fd7e14', '#6f42c1'];

const totalWeight = computed(() => variants.value.reduce((s, v) => s + v.weight, 0));

function variantPercent(w: number): number {
  return totalWeight.value > 0 ? (w / totalWeight.value) * 100 : 0;
}

function variantColor(i: number): string {
  return VARIANT_COLORS[i % VARIANT_COLORS.length] ?? '#635BFF';
}

onMounted(() => {
  // Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
});

async function loadVariants() {
  loading.value = true;
  error.value = '';
  try {
    const res = await variantsApi.list(props.linkId);
    if (res.data) variants.value = res.data;
  } catch {
    error.value = 'Failed to load variants.';
  } finally {
    loading.value = false;
  }
}

function show() {
  isSplitTest.value = props.isSplitTestInitial;
  variants.value = [];
  error.value = '';
  showAddForm.value = false;
  editingId.value = '';
  isOpen.value = true;
  loadVariants();
  modalInstance?.show();
}

function hide() {
  isOpen.value = false;
  modalInstance?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

async function toggleSplitTest(e: Event) {
  const enabled = (e.target as HTMLInputElement).checked;
  toggling.value = true;
  error.value = '';
  try {
    await variantsApi.toggleSplitTest(props.linkId, enabled);
    isSplitTest.value = enabled;
    emit('updated', enabled);
  } catch {
    error.value = 'Failed to toggle split test.';
  } finally {
    toggling.value = false;
  }
}

function openAddForm() {
  newForm.value = { name: `Variant ${String.fromCharCode(65 + variants.value.length)}`, destination_url: '', weight: 50 };
  showAddForm.value = true;
}

async function submitAdd() {
  if (!newForm.value.name.trim() || !newForm.value.destination_url.trim()) {
    error.value = 'Name and destination URL are required.';
    return;
  }
  saving.value = true;
  error.value = '';
  try {
    const res = await variantsApi.create(props.linkId, {
      name: newForm.value.name.trim(),
      destination_url: newForm.value.destination_url.trim(),
      weight: newForm.value.weight || 50,
    });
    if (res.data) variants.value.push(res.data);
    showAddForm.value = false;
  } catch {
    error.value = 'Failed to create variant.';
  } finally {
    saving.value = false;
  }
}

function startEdit(v: LinkVariant) {
  editingId.value = v.id;
  editForm.value = { name: v.name, destination_url: v.destination_url, weight: v.weight };
  showAddForm.value = false;
}

function cancelEdit() {
  editingId.value = '';
}

async function submitEdit(v: LinkVariant) {
  saving.value = true;
  error.value = '';
  try {
    const res = await variantsApi.update(props.linkId, v.id, {
      name: editForm.value.name.trim() || undefined,
      destination_url: editForm.value.destination_url.trim() || undefined,
      weight: editForm.value.weight || undefined,
    });
    if (res.data) {
      const idx = variants.value.findIndex(x => x.id === v.id);
      if (idx !== -1) variants.value[idx] = res.data!;
    }
    editingId.value = '';
  } catch {
    error.value = 'Failed to update variant.';
  } finally {
    saving.value = false;
  }
}

async function deleteVariant(v: LinkVariant) {
  if (!window.confirm(`Delete variant "${v.name}"?`)) return;
  error.value = '';
  try {
    await variantsApi.delete(props.linkId, v.id);
    variants.value = variants.value.filter(x => x.id !== v.id);
  } catch {
    error.value = 'Failed to delete variant.';
  }
}

defineExpose({ show, hide });
</script>

<style scoped>
.split-error-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  margin-bottom: 16px;
  background: var(--md-sys-color-error-container, #FFDAD6);
  color: var(--md-sys-color-on-error-container, #410002);
  border-radius: 8px;
}

/* Toggle row */
.split-toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  background: var(--md-sys-color-surface-container-low);
  margin-bottom: 20px;
}

.split-toggle-title {
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface);
}

.split-toggle-subtitle {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
}

/* Distribution bar */
.split-distribution {
  margin-bottom: 20px;
}

.split-distribution-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.split-distribution-label {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
}

.split-distribution-total {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
}

.split-bar {
  height: 18px;
  border-radius: 6px;
  overflow: hidden;
  display: flex;
}

.split-bar-segment {
  height: 100%;
  transition: width 0.3s;
}

.split-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 8px;
}

.split-legend-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
}

.split-legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

/* Loading / empty */
.split-loading {
  display: flex;
  justify-content: center;
  padding: 24px;
}

.split-empty {
  text-align: center;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  padding: 16px;
}

/* Variant cards */
.variant-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  margin-bottom: 8px;
  overflow: hidden;
  transition: box-shadow 0.15s;
}

.variant-card:hover {
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}

.variant-view {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
}

.variant-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
  margin-top: 4px;
}

.variant-info {
  flex: 1;
  min-width: 0;
}

.variant-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 4px;
}

.variant-name {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.variant-url {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.variant-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
}

/* Edit form */
.variant-edit-form {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  flex-wrap: wrap;
}

/* Add variant */
.add-variant-card {
  border: 1px dashed var(--md-sys-color-outline-variant);
  border-radius: 10px;
  padding: 12px 14px;
  margin-top: 8px;
  background: var(--md-sys-color-surface-container-low);
}

.add-variant-title {
  font-size: 0.85rem;
  font-weight: 500;
  margin: 0 0 10px;
  color: var(--md-sys-color-on-surface);
}

.add-variant-hint {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 8px;
}
</style>
