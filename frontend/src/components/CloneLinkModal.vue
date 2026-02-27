<template>
  <BaseModal v-model="isOpen" size="md" @closed="onDialogClosed">
    <template #headline>
      <span class="material-symbols-outlined dialog-headline-icon">content_copy</span>
      Clone Link
    </template>

    <div class="dialog-content">

      <!-- Source link card -->
      <div v-if="link" class="source-card">
        <div class="source-card-header">
          <span class="material-symbols-outlined source-card-icon">link</span>
          <span class="source-card-eyebrow">Cloning this link</span>
        </div>

        <div class="source-card-body">
          <div class="source-title">{{ link.title || link.slug }}</div>

          <div class="source-url-row">
            <span class="source-url-badge">Short</span>
            <span class="source-url-value source-url-value--short">{{ link.short_url }}</span>
          </div>

          <div class="source-url-row">
            <span class="source-url-badge source-url-badge--dest">Dest</span>
            <span class="source-url-value" :title="link.destination_url">{{ link.destination_url }}</span>
          </div>
        </div>

        <div class="source-card-footer">
          <span class="material-symbols-outlined" style="font-size:14px;color:var(--md-sys-color-on-surface-variant)">info</span>
          <span>A new link will be created with the same destination URL and settings.</span>
        </div>
      </div>

      <!-- Divider with label -->
      <div class="section-divider">
        <span class="section-divider-line"></span>
        <span class="section-divider-label">Customise the clone</span>
        <span class="section-divider-line"></span>
      </div>

      <!-- New Title -->
      <div class="field-group">
        <md-outlined-text-field
          :value="newTitle"
          @input="newTitle = ($event.target as HTMLInputElement).value"
          label="New Title"
          placeholder="Leave blank to keep original title"
          maxlength="500"
          class="field-full"
          supporting-text="Optionally give this clone a different title."
        >
          <span class="material-symbols-outlined" slot="leading-icon">title</span>
        </md-outlined-text-field>
      </div>

      <!-- Custom Slug -->
      <div class="field-group">
        <div class="slug-row">
          <div class="slug-prefix-chip">
            <span class="material-symbols-outlined" style="font-size:14px">link</span>
            {{ baseSlug }}/
          </div>
          <md-outlined-text-field
            :value="newSlug"
            @input="newSlug = ($event.target as HTMLInputElement).value; slugError = ''"
            label="Custom Slug"
            placeholder="auto-generated if blank"
            maxlength="20"
            class="slug-field-input"
            :error="!!slugError"
            :error-text="slugError"
            supporting-text="3–20 alphanumeric characters, or leave blank."
          >
            <span class="material-symbols-outlined" slot="leading-icon">tag</span>
          </md-outlined-text-field>
        </div>
      </div>

      <!-- Error message -->
      <div v-if="errorMsg" class="alert-banner">
        <span class="material-symbols-outlined alert-icon">error</span>
        <span class="alert-text">{{ errorMsg }}</span>
      </div>

    </div>

    <template #actions>
      <md-text-button @click="hide" :disabled="cloning">Cancel</md-text-button>
      <md-filled-button :disabled="cloning" @click="doClone" class="clone-btn">
        <md-circular-progress v-if="cloning" indeterminate style="--md-circular-progress-size:18px;margin-right:6px" />
        <span class="material-symbols-outlined" v-else style="font-size:18px;margin-right:6px">content_copy</span>
        {{ cloning ? 'Cloning…' : 'Clone Link' }}
      </md-filled-button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import linksApi from '@/api/links';
import type { LinkResponse } from '@/types/links';
import BaseModal from '@/components/BaseModal.vue';

const props = defineProps<{
  modalId?: string;
  link: LinkResponse | null;
}>();

const emit = defineEmits<{
  (e: 'cloned', link: LinkResponse): void;
}>();

const modalId = props.modalId ?? 'clone-link-modal';
const modalEl = ref<HTMLElement | null>(null);
let bsModal: any = null;

const isOpen = ref(false);
const newTitle = ref('');
const newSlug = ref('');
const slugError = ref('');
const errorMsg = ref('');
const cloning = ref(false);

const baseSlug = computed(() => {
  if (!props.link) return '';
  try {
    return new URL(props.link.short_url).host;
  } catch {
    return '';
  }
});

onMounted(() => {
  // Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
});

onBeforeUnmount(() => {
  bsModal?.dispose();
});

function show() {
  resetForm();
  isOpen.value = true;
  bsModal?.show();
}

function hide() {
  isOpen.value = false;
  bsModal?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

function resetForm() {
  newTitle.value = '';
  newSlug.value = '';
  slugError.value = '';
  errorMsg.value = '';
  cloning.value = false;
}

async function doClone() {
  if (!props.link) return;

  const slug = newSlug.value.trim();
  if (slug && !/^[a-zA-Z0-9]{3,20}$/.test(slug)) {
    slugError.value = 'Slug must be 3–20 alphanumeric characters.';
    return;
  }

  cloning.value = true;
  errorMsg.value = '';
  try {
    const res = await linksApi.clone(props.link.id, {
      new_title: newTitle.value.trim() || undefined,
      new_slug: slug || undefined,
    });
    if (res.data) {
      emit('cloned', res.data);
      hide();
    }
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: { description?: string } } } };
    errorMsg.value = axiosErr.response?.data?.error?.description ?? 'Failed to clone link. Please try again.';
  } finally {
    cloning.value = false;
  }
}

defineExpose({ show, hide });
</script>

<style scoped lang="scss">
/* ── Headline ─────────────────────────────────────────── */
.dialog-headline-icon {
  font-size: 22px;
  color: var(--md-sys-color-primary);
}

/* ── Content ──────────────────────────────────────────── */
.dialog-content {
  min-width: 440px;
  max-width: 100%;
  padding: 0 2px;
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* ── Source card ──────────────────────────────────────── */
.source-card {
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 20px;
}

.source-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px 8px;
  background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 8%, transparent);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.source-card-icon {
  font-size: 16px;
  color: var(--md-sys-color-primary);
}

.source-card-eyebrow {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--md-sys-color-primary);
}

.source-card-body {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.source-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 4px;
}

.source-url-row {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
}

.source-url-badge {
  flex-shrink: 0;
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 2px 6px;
  border-radius: 6px;
  background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 14%, transparent);
  color: var(--md-sys-color-primary);

  &--dest {
    background: var(--md-sys-color-surface-container-highest);
    color: var(--md-sys-color-on-surface-variant);
  }
}

.source-url-value {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;

  &--short {
    color: var(--md-sys-color-primary);
    font-weight: 500;
  }
}

.source-card-footer {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px 10px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.4;
}

/* ── Section divider ──────────────────────────────────── */
.section-divider {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 18px;
}

.section-divider-line {
  flex: 1;
  height: 1px;
  background: var(--md-sys-color-outline-variant);
}

.section-divider-label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.07em;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
}

/* ── Fields ───────────────────────────────────────────── */
.field-group {
  margin-bottom: 16px;
}

.field-full {
  width: 100%;
}

.slug-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.slug-prefix-chip {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 10px 12px;
  margin-top: 8px;
  background: var(--md-sys-color-surface-container);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  font-size: 0.78rem;
  font-family: monospace;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  flex-shrink: 0;
}

.slug-field-input {
  flex: 1;
}

/* ── Alert banner ─────────────────────────────────────── */
.alert-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 12px;
  background: var(--md-sys-color-error-container, #ffdad6);
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 30%, transparent);
  color: var(--md-sys-color-on-error-container, #410002);
  font-size: 0.875rem;
}

.alert-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.alert-text {
  flex: 1;
}

/* ── Actions ──────────────────────────────────────────── */
.clone-btn {
  min-width: 130px;
}
</style>
