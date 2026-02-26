<template>
  <div ref="modalEl" class="modal fade" tabindex="-1" aria-hidden="true">
    <div class="modal-dialog modal-lg modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <div>
            <h5 class="modal-title fw-semibold mb-0">A/B Split Test</h5>
            <p class="text-muted mb-0" style="font-size: 0.8rem;">{{ props.slug }}</p>
          </div>
          <button type="button" class="btn-close" @click="hide" aria-label="Close"></button>
        </div>

        <div class="modal-body">
          <!-- Error -->
          <div v-if="error" class="alert alert-danger alert-dismissible" role="alert">
            {{ error }}
            <button type="button" class="btn-close" @click="error = ''" aria-label="Close"></button>
          </div>

          <!-- Enable toggle -->
          <div class="d-flex align-items-center justify-content-between mb-4 p-3 border rounded bg-light">
            <div>
              <div class="fw-semibold" style="font-size: 0.9rem;">Split Test {{ isSplitTest ? 'Enabled' : 'Disabled' }}</div>
              <div class="text-muted" style="font-size: 0.8rem;">
                {{ isSplitTest ? 'Traffic is being distributed across variants.' : 'Enable to start routing traffic to variants.' }}
              </div>
            </div>
            <div class="form-check form-switch mb-0">
              <input
                class="form-check-input"
                type="checkbox"
                role="switch"
                :checked="isSplitTest"
                :disabled="toggling || loading"
                @change="toggleSplitTest"
                style="width: 2.5rem; height: 1.25rem; cursor: pointer;"
              />
            </div>
          </div>

          <!-- Weight summary bar -->
          <div v-if="variants.length > 0" class="mb-4">
            <div class="d-flex justify-content-between mb-1">
              <span class="small fw-medium">Traffic distribution</span>
              <span class="small text-muted">Total weight: {{ totalWeight }}</span>
            </div>
            <div class="progress" style="height: 18px; border-radius: 6px;">
              <div
                v-for="(v, i) in variants"
                :key="v.id"
                class="progress-bar"
                :style="{ width: variantPercent(v.weight) + '%', backgroundColor: variantColor(i) }"
                :title="`${v.name}: ${variantPercent(v.weight).toFixed(1)}%`"
              ></div>
            </div>
            <div class="d-flex flex-wrap gap-2 mt-2">
              <span v-for="(v, i) in variants" :key="v.id" class="d-flex align-items-center gap-1" style="font-size: 0.78rem;">
                <span class="rounded-circle d-inline-block" :style="{ width: '8px', height: '8px', backgroundColor: variantColor(i) }"></span>
                {{ v.name }}: {{ variantPercent(v.weight).toFixed(1) }}%
              </span>
            </div>
          </div>

          <!-- Variants list -->
          <div v-if="loading" class="text-center py-4">
            <div class="spinner-border spinner-border-sm text-primary" role="status"></div>
          </div>

          <div v-else>
            <div v-if="variants.length === 0" class="text-center text-muted py-3" style="font-size: 0.88rem;">
              No variants yet. Add one below to get started.
            </div>

            <div v-for="(v, i) in variants" :key="v.id" class="variant-card border rounded p-3 mb-2">
              <div v-if="editingId === v.id">
                <!-- Edit mode -->
                <div class="row g-2">
                  <div class="col-12 col-md-4">
                    <input v-model="editForm.name" class="form-control form-control-sm" placeholder="Variant name" />
                  </div>
                  <div class="col-12 col-md-5">
                    <input v-model="editForm.destination_url" class="form-control form-control-sm" placeholder="https://..." />
                  </div>
                  <div class="col-6 col-md-2">
                    <input v-model.number="editForm.weight" class="form-control form-control-sm" type="number" min="1" max="1000" placeholder="Weight" />
                  </div>
                  <div class="col-6 col-md-1 d-flex gap-1 justify-content-end">
                    <button class="btn btn-sm btn-primary px-2" @click="submitEdit(v)" :disabled="saving">✓</button>
                    <button class="btn btn-sm btn-outline-secondary px-2" @click="cancelEdit">✕</button>
                  </div>
                </div>
              </div>

              <div v-else class="d-flex align-items-start gap-3">
                <span class="variant-dot flex-shrink-0 rounded-circle mt-1" :style="{ backgroundColor: variantColor(i) }"></span>
                <div class="flex-grow-1 min-width-0">
                  <div class="d-flex align-items-center gap-2 flex-wrap">
                    <span class="fw-semibold" style="font-size: 0.88rem;">{{ v.name }}</span>
                    <span class="badge bg-light text-secondary border" style="font-size: 0.72rem;">weight {{ v.weight }}</span>
                    <span class="badge text-bg-secondary" style="font-size: 0.72rem;">{{ v.click_count.toLocaleString() }} clicks</span>
                  </div>
                  <div class="text-muted text-truncate mt-1" style="font-size: 0.8rem;" :title="v.destination_url">
                    {{ v.destination_url }}
                  </div>
                </div>
                <div class="d-flex gap-1 flex-shrink-0">
                  <button class="btn btn-sm border-0 p-1 text-muted" title="Edit" @click="startEdit(v)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325"/>
                    </svg>
                  </button>
                  <button class="btn btn-sm border-0 p-1 text-danger" title="Delete" @click="deleteVariant(v)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                      <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>

            <!-- Add variant form -->
            <div class="border rounded p-3 mt-3 bg-light" v-if="showAddForm">
              <p class="fw-medium small mb-2">New Variant</p>
              <div class="row g-2">
                <div class="col-12 col-md-4">
                  <input v-model="newForm.name" class="form-control form-control-sm" placeholder="Variant name (e.g. Variant B)" />
                </div>
                <div class="col-12 col-md-5">
                  <input v-model="newForm.destination_url" class="form-control form-control-sm" placeholder="https://destination.com" />
                </div>
                <div class="col-6 col-md-2">
                  <input v-model.number="newForm.weight" class="form-control form-control-sm" type="number" min="1" max="1000" placeholder="Weight" />
                </div>
                <div class="col-6 col-md-1 d-flex gap-1 justify-content-end">
                  <button class="btn btn-sm btn-primary px-2" @click="submitAdd" :disabled="saving">✓</button>
                  <button class="btn btn-sm btn-outline-secondary px-2" @click="showAddForm = false">✕</button>
                </div>
              </div>
              <div class="form-text">Weight is relative (e.g., 50 + 50 = 50/50 split; 70 + 30 = 70/30 split).</div>
            </div>

            <button
              v-if="!showAddForm && editingId === ''"
              class="btn btn-sm btn-outline-primary mt-3 d-flex align-items-center gap-1"
              @click="openAddForm"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
              </svg>
              Add Variant
            </button>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary" @click="hide">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import variantsApi from '@/api/variants';
import type { LinkVariant } from '@/types/links';

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
let modalInstance: Modal | null = null;

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
  return VARIANT_COLORS[i % VARIANT_COLORS.length];
}

onMounted(() => {
  if (modalEl.value) {
    modalInstance = new Modal(modalEl.value, { backdrop: 'static' });
  }
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
  loadVariants();
  modalInstance?.show();
}

function hide() {
  modalInstance?.hide();
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
  color: white;
}

.variant-card {
  transition: box-shadow 0.15s;
}

.variant-card:hover {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.variant-dot {
  width: 10px;
  height: 10px;
  display: inline-block;
  flex-shrink: 0;
}

.form-check-input:checked {
  background-color: #635bff;
  border-color: #635bff;
}
</style>
