<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" style="--md-dialog-container-shape:16px">
    <div slot="headline">Clone Link</div>

    <div slot="content" style="min-width:440px;max-width:100%;padding:0 4px">
      <!-- Source summary -->
      <div v-if="link" class="clone-source-card">
        <div class="clone-source-label">Cloning</div>
        <div class="clone-source-name">{{ link.title || link.slug }}</div>
        <div class="clone-source-url">{{ link.short_url }}</div>
      </div>

      <!-- New Title -->
      <div class="field-group">
        <md-outlined-text-field
          :value="newTitle"
          @input="newTitle = ($event.target as HTMLInputElement).value"
          label="New Title (optional)"
          placeholder="Leave blank to keep original title"
          maxlength="500"
          style="width:100%"
        />
      </div>

      <!-- Custom Slug -->
      <div class="field-group">
        <div class="slug-field">
          <span class="slug-prefix">{{ baseSlug }}/</span>
          <md-outlined-text-field
            :value="newSlug"
            @input="newSlug = ($event.target as HTMLInputElement).value; slugError = ''"
            label="Custom Slug (optional)"
            placeholder="auto-generated if blank"
            maxlength="20"
            style="flex:1"
            :error="!!slugError"
            :error-text="slugError"
            supporting-text="3–20 alphanumeric characters."
          />
        </div>
      </div>

      <!-- Error message -->
      <div v-if="errorMsg" class="clone-error-banner">
        <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-error)">error</span>
        <span>{{ errorMsg }}</span>
      </div>
    </div>

    <div slot="actions">
      <md-text-button @click="hide">Cancel</md-text-button>
      <md-filled-button :disabled="cloning" @click="doClone">
        <md-circular-progress v-if="cloning" indeterminate style="--md-circular-progress-size:18px" slot="icon" />
        Clone Link
      </md-filled-button>
    </div>
  </md-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import linksApi from '@/api/links';
import type { LinkResponse } from '@/types/links';

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

<style scoped>
.clone-source-card {
  padding: 12px 16px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 12px;
  margin-bottom: 20px;
}

.clone-source-label {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 4px;
}

.clone-source-name {
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 2px;
  color: var(--md-sys-color-on-surface);
}

.clone-source-url {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.field-group {
  margin-bottom: 16px;
}

.slug-field {
  display: flex;
  align-items: center;
  gap: 8px;
}

.slug-prefix {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  flex-shrink: 0;
  padding-top: 4px;
}

.clone-error-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: var(--md-sys-color-error-container, #FFDAD6);
  color: var(--md-sys-color-on-error-container, #410002);
  border-radius: 8px;
  font-size: 0.875rem;
}
</style>
