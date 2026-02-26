<template>
  <div class="modal fade" :id="modalId" tabindex="-1" ref="modalEl">
    <div class="modal-dialog modal-lg">
      <div class="modal-content border-0 shadow">
        <div class="modal-header border-0 pb-0">
          <div>
            <h5 class="modal-title fw-bold">Retargeting Pixels</h5>
            <p class="text-muted small mb-0">
              <code class="small">{{ link?.slug }}</code>
            </p>
          </div>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>

        <div class="modal-body">
          <!-- Enable/disable toggle -->
          <div class="d-flex align-items-center justify-content-between p-3 rounded mb-4"
               :class="isEnabled ? 'bg-primary bg-opacity-10 border border-primary border-opacity-25' : 'bg-light'">
            <div>
              <div class="fw-semibold small">Pixel tracking</div>
              <div class="text-muted" style="font-size:0.8rem">
                When enabled, clicks serve an HTML page that fires all pixels then redirects
              </div>
            </div>
            <div class="form-check form-switch mb-0">
              <input
                class="form-check-input"
                type="checkbox"
                role="switch"
                v-model="isEnabled"
                :disabled="toggling"
                @change="toggleTracking"
                style="width:2.5rem;height:1.25rem;cursor:pointer"
              />
            </div>
          </div>

          <!-- Warning when enabled with no pixels -->
          <div v-if="isEnabled && pixels.length === 0 && !loading" class="alert alert-warning py-2 small mb-3">
            Pixel tracking is enabled but no pixels are configured. Add at least one pixel below.
          </div>

          <!-- Existing pixels -->
          <div v-if="loading" class="text-center py-3">
            <div class="spinner-border spinner-border-sm text-primary"></div>
          </div>

          <div v-else-if="pixels.length > 0" class="mb-3">
            <h6 class="fw-semibold small text-muted text-uppercase mb-2">Configured pixels</h6>
            <div class="list-group list-group-flush border rounded">
              <div
                v-for="px in pixels"
                :key="px.id"
                class="list-group-item d-flex align-items-center gap-3 py-2"
              >
                <span style="font-size:1.25rem">{{ pixelIcon(px.pixel_type) }}</span>
                <div class="flex-fill min-w-0">
                  <div class="fw-medium small">{{ pixelLabel(px.pixel_type) }}</div>
                  <div class="text-muted" style="font-size:0.78rem;word-break:break-all">
                    {{ px.pixel_type === 'custom' ? 'Custom script' : px.pixel_id }}
                  </div>
                </div>
                <button
                  class="btn btn-outline-danger btn-sm flex-shrink-0"
                  @click="deletePixel(px.id)"
                  title="Remove pixel"
                >
                  <i class="bi bi-trash"></i>
                </button>
              </div>
            </div>
          </div>

          <!-- Add pixel form -->
          <div class="border rounded p-3">
            <h6 class="fw-semibold small mb-3">Add pixel</h6>

            <!-- Pixel type selector -->
            <div class="mb-3">
              <label class="form-label small fw-medium">Pixel type</label>
              <div class="d-flex flex-wrap gap-2">
                <button
                  v-for="pt in PIXEL_TYPES"
                  :key="pt.value"
                  class="btn btn-sm"
                  :class="selectedType === pt.value ? 'btn-primary' : 'btn-outline-secondary'"
                  @click="selectedType = pt.value; pixelID = ''; customScript = ''"
                >
                  {{ pt.icon }} {{ pt.label }}
                </button>
              </div>
            </div>

            <!-- Pixel ID input (for named types) -->
            <div v-if="selectedTypeMeta?.requiresId" class="mb-3">
              <label class="form-label small fw-medium">Pixel / Tag ID</label>
              <input
                v-model="pixelID"
                type="text"
                class="form-control"
                :placeholder="selectedTypeMeta?.placeholder"
              />
              <div class="form-text">{{ selectedTypeMeta?.description }}</div>
            </div>

            <!-- Custom script textarea -->
            <div v-if="selectedType === 'custom'" class="mb-3">
              <label class="form-label small fw-medium">Custom script</label>
              <textarea
                v-model="customScript"
                class="form-control font-monospace small"
                rows="5"
                :placeholder="selectedTypeMeta?.placeholder"
              ></textarea>
              <div class="form-text text-warning">
                ⚠️ Custom scripts are inserted verbatim into the redirect page. Ensure the code is trusted.
              </div>
            </div>

            <div v-if="addError" class="alert alert-danger py-2 small mb-2">{{ addError }}</div>

            <button
              class="btn btn-primary btn-sm"
              :disabled="adding || !canAdd"
              @click="addPixel"
            >
              <span v-if="adding" class="spinner-border spinner-border-sm me-1"></span>
              Add pixel
            </button>
          </div>

          <!-- Info box -->
          <div class="mt-3 p-3 bg-light rounded small text-muted">
            <strong>How it works:</strong> When pixel tracking is enabled, clicking your short link
            loads a lightweight HTML page that fires all your tracking pixels, then immediately
            redirects visitors to the destination URL via JavaScript.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import type { Modal } from 'bootstrap';
import pixelsApi from '@/api/pixels';
import { PIXEL_TYPES } from '@/types/pixels';
import type { RetargetingPixel, PixelType } from '@/types/pixels';
import type { LinkResponse } from '@/types/links';

const props = defineProps<{
  modalId: string;
  link: LinkResponse | null;
}>();

const emit = defineEmits<{
  (e: 'pixel-tracking-updated', enabled: boolean): void;
}>();

const modalEl = ref<HTMLElement | null>(null);
let bsModal: Modal | null = null;

const pixels = ref<RetargetingPixel[]>([]);
const loading = ref(false);
const isEnabled = ref(false);
const toggling = ref(false);

// Add form
const selectedType = ref<PixelType>('facebook');
const pixelID = ref('');
const customScript = ref('');
const adding = ref(false);
const addError = ref('');

const selectedTypeMeta = computed(() => PIXEL_TYPES.find((pt) => pt.value === selectedType.value));

const canAdd = computed(() => {
  if (selectedType.value === 'custom') return customScript.value.trim().length > 0;
  return pixelID.value.trim().length > 0;
});

function pixelIcon(type: string) {
  return PIXEL_TYPES.find((pt) => pt.value === type)?.icon ?? '⚙️';
}
function pixelLabel(type: string) {
  return PIXEL_TYPES.find((pt) => pt.value === type)?.label ?? type;
}

async function loadPixels() {
  if (!props.link) return;
  loading.value = true;
  try {
    const res = await pixelsApi.list(props.link.id);
    pixels.value = res.data ?? [];
  } finally {
    loading.value = false;
  }
}

async function toggleTracking() {
  if (!props.link) return;
  toggling.value = true;
  try {
    await pixelsApi.toggle(props.link.id, isEnabled.value);
    emit('pixel-tracking-updated', isEnabled.value);
  } catch {
    isEnabled.value = !isEnabled.value; // revert on error
  } finally {
    toggling.value = false;
  }
}

async function addPixel() {
  if (!props.link || !canAdd.value) return;
  adding.value = true;
  addError.value = '';
  try {
    const res = await pixelsApi.create(props.link.id, {
      pixel_type: selectedType.value,
      pixel_id: selectedType.value !== 'custom' ? pixelID.value.trim() : undefined,
      custom_script: selectedType.value === 'custom' ? customScript.value.trim() : undefined,
    });
    pixels.value.push(res.data);
    pixelID.value = '';
    customScript.value = '';
  } catch (err: any) {
    addError.value = err?.response?.data?.description ?? 'Failed to add pixel.';
  } finally {
    adding.value = false;
  }
}

async function deletePixel(pixelId: string) {
  if (!props.link) return;
  if (!confirm('Remove this pixel?')) return;
  await pixelsApi.delete(props.link.id, pixelId);
  pixels.value = pixels.value.filter((p) => p.id !== pixelId);
}

onMounted(() => {
  if (modalEl.value) {
    import('bootstrap').then(({ Modal }) => {
      bsModal = new Modal(modalEl.value!);
      modalEl.value!.addEventListener('shown.bs.modal', () => {
        isEnabled.value = props.link?.is_pixel_tracking ?? false;
        loadPixels();
      });
    });
  }
});

onBeforeUnmount(() => {
  bsModal?.dispose();
});

defineExpose({ show: () => bsModal?.show(), hide: () => bsModal?.hide() });
</script>
