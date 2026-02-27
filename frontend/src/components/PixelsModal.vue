<template>
  <BaseModal v-model="isOpen" size="md" @closed="onDialogClosed">
    <template #headline>
      Retargeting Pixels
      <span style="font-size:0.82rem;font-weight:400;color:var(--md-sys-color-on-surface-variant);display:block;margin-top:2px">
        <code>{{ link?.slug }}</code>
      </span>
    </template>

    <div style="min-width:520px;max-width:100%;padding:0 4px;">
      <!-- Enable/disable toggle -->
      <div class="pixels-toggle-row" :class="isEnabled ? 'pixels-toggle-row--active' : ''">
        <div>
          <div class="pixels-toggle-title">Pixel tracking</div>
          <div class="pixels-toggle-subtitle">
            When enabled, clicks serve an HTML page that fires all pixels then redirects
          </div>
        </div>
        <md-switch
          :selected="isEnabled"
          :disabled="toggling"
          @change="isEnabled = ($event.target as HTMLInputElement).checked; toggleTracking()"
        />
      </div>

      <!-- Warning when enabled with no pixels -->
      <div v-if="isEnabled && pixels.length === 0 && !loading" class="pixels-warning-banner">
        <span class="material-symbols-outlined" style="font-size:18px;color:#F4A100">warning</span>
        <span>Pixel tracking is enabled but no pixels are configured. Add at least one pixel below.</span>
      </div>

      <!-- Existing pixels loading -->
      <div v-if="loading" class="pixels-loading">
        <md-circular-progress indeterminate />
      </div>

      <!-- Existing pixels list -->
      <div v-else-if="pixels.length > 0" class="pixels-list-section">
        <div class="pixels-list-label">Configured pixels</div>
        <div class="pixels-list">
          <div v-for="px in pixels" :key="px.id" class="pixel-item">
            <span class="pixel-item-icon">{{ pixelIcon(px.pixel_type) }}</span>
            <div class="pixel-item-info">
              <div class="pixel-item-name">{{ pixelLabel(px.pixel_type) }}</div>
              <div class="pixel-item-id">
                {{ px.pixel_type === 'custom' ? 'Custom script' : px.pixel_id }}
              </div>
            </div>
            <button class="btn-icon"
              @click="deletePixel(px.id)"
              title="Remove pixel"
              style="width:32px;height:32px;"
            >
              <span class="material-symbols-outlined" style="font-size:18px">delete</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Add pixel form -->
      <div class="add-pixel-section">
        <div class="add-pixel-title">Add pixel</div>

        <!-- Pixel type selector -->
        <div class="pixel-type-section">
          <div class="pixel-type-label">Pixel type</div>
          <div class="pixel-type-buttons">
            <button
              v-for="pt in PIXEL_TYPES"
              :key="pt.value"
              class="pixel-type-btn"
              :class="{ 'pixel-type-btn--active': selectedType === pt.value }"
              @click="selectedType = pt.value; pixelID = ''; customScript = ''"
            >
              {{ pt.icon }} {{ pt.label }}
            </button>
          </div>
        </div>

        <!-- Pixel ID input (for named types) -->
        <div v-if="selectedTypeMeta?.requiresId" class="field-group">
          <md-outlined-text-field
            :value="pixelID"
            @input="pixelID = ($event.target as HTMLInputElement).value"
            label="Pixel / Tag ID"
            :placeholder="selectedTypeMeta?.placeholder"
            :supporting-text="selectedTypeMeta?.description"
            style="width:100%"
          />
        </div>

        <!-- Custom script textarea -->
        <div v-if="selectedType === 'custom'" class="field-group">
          <md-outlined-text-field
            type="textarea"
            :value="customScript"
            @input="customScript = ($event.target as HTMLInputElement).value"
            label="Custom script"
            :placeholder="selectedTypeMeta?.placeholder"
            rows="5"
            style="width:100%;font-family:monospace"
          />
          <div class="custom-script-warning">
            <span class="material-symbols-outlined" style="font-size:16px;color:#F4A100">warning</span>
            Custom scripts are inserted verbatim into the redirect page. Ensure the code is trusted.
          </div>
        </div>

        <div v-if="addError" class="add-error-banner">
          <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-error)">error</span>
          {{ addError }}
        </div>

        <button class="btn-filled" :disabled="adding || !canAdd" @click="addPixel">
          <md-circular-progress v-if="adding" indeterminate />
          <span v-else class="material-symbols-outlined">add</span>
          Add pixel
        </button>
      </div>

      <!-- Info box -->
      <div class="pixels-info-box">
        <strong>How it works:</strong> When pixel tracking is enabled, clicking your short link
        loads a lightweight HTML page that fires all your tracking pixels, then immediately
        redirects visitors to the destination URL via JavaScript.
      </div>
    </div>

    <template #actions>
      <button class="btn-text" @click="hide">Close</button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import pixelsApi from '@/api/pixels';
import { PIXEL_TYPES } from '@/types/pixels';
import type { RetargetingPixel, PixelType } from '@/types/pixels';
import type { LinkResponse } from '@/types/links';
import BaseModal from '@/components/BaseModal.vue';

const props = defineProps<{
  modalId: string;
  link: LinkResponse | null;
}>();

const emit = defineEmits<{
  (e: 'pixel-tracking-updated', enabled: boolean): void;
}>();

const modalEl = ref<HTMLElement | null>(null);
let bsModal: any = null;

const isOpen = ref(false);
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
  // Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
});

onBeforeUnmount(() => {
  bsModal?.dispose();
});

function show() {
  isOpen.value = true;
  isEnabled.value = props.link?.is_pixel_tracking ?? false;
  loadPixels();
  bsModal?.show();
}

function hide() {
  isOpen.value = false;
  bsModal?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

defineExpose({ show, hide });
</script>

<style scoped>
/* Toggle row */
.pixels-toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  border-radius: 12px;
  background: var(--md-sys-color-surface-container-low);
  margin-bottom: 16px;
  transition: background 0.2s;
}

.pixels-toggle-row--active {
  background: rgba(99, 91, 255, 0.08);
  border: 1px solid rgba(99, 91, 255, 0.2);
}

.pixels-toggle-title {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.pixels-toggle-subtitle {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
}

/* Warning banner */
.pixels-warning-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: rgba(244, 161, 0, 0.1);
  border: 1px solid rgba(244, 161, 0, 0.4);
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 0.875rem;
}

/* Loading */
.pixels-loading {
  display: flex;
  justify-content: center;
  padding: 20px;
}

/* Pixels list */
.pixels-list-section {
  margin-bottom: 20px;
}

.pixels-list-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 8px;
}

.pixels-list {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  overflow: hidden;
}

.pixel-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.pixel-item:last-child {
  border-bottom: none;
}

.pixel-item-icon {
  font-size: 1.2rem;
  flex-shrink: 0;
}

.pixel-item-info {
  flex: 1;
  min-width: 0;
}

.pixel-item-name {
  font-weight: 500;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.pixel-item-id {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
  word-break: break-all;
}

/* Add pixel section */
.add-pixel-section {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
}

.add-pixel-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 12px;
}

.pixel-type-section {
  margin-bottom: 16px;
}

.pixel-type-label {
  font-size: 0.82rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.pixel-type-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.pixel-type-btn {
  padding: 6px 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface);
  transition: background 0.15s, border-color 0.15s;
}

.pixel-type-btn:hover {
  background: var(--md-sys-color-surface-container-low);
}

.pixel-type-btn--active {
  background: var(--md-sys-color-primary-container, rgba(99,91,255,0.12));
  border-color: var(--md-sys-color-primary);
  color: var(--md-sys-color-primary);
}

.field-group {
  margin-bottom: 16px;
}

.custom-script-warning {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 0.78rem;
  color: #F4A100;
}

.add-error-banner {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: var(--md-sys-color-error-container, #FFDAD6);
  color: var(--md-sys-color-on-error-container, #410002);
  border-radius: 8px;
  font-size: 0.82rem;
}

/* Info box */
.pixels-info-box {
  padding: 12px 16px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 10px;
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.5;
}
</style>
