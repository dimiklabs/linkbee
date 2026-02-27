<template>
  <div :id="modalId" class="modal fade" tabindex="-1" aria-hidden="true" ref="modalEl">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow">

        <div class="modal-header border-0 pb-0">
          <h5 class="modal-title fw-semibold">Clone Link</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>

        <div class="modal-body pt-3">
          <!-- Source summary -->
          <div v-if="link" class="rounded-3 p-3 mb-4" style="background: var(--bs-secondary-bg, #f8f9fa);">
            <div class="small text-muted mb-1">Cloning</div>
            <div class="fw-semibold text-truncate mb-1">{{ link.title || link.slug }}</div>
            <div class="small text-muted text-truncate">{{ link.short_url }}</div>
          </div>

          <!-- Optional overrides -->
          <div class="mb-3">
            <label class="form-label fw-medium small">New Title <span class="text-muted fw-normal">(optional)</span></label>
            <input
              v-model="newTitle"
              type="text"
              class="form-control"
              placeholder="Leave blank to keep original title"
              maxlength="500"
            />
          </div>

          <div class="mb-3">
            <label class="form-label fw-medium small">Custom Slug <span class="text-muted fw-normal">(optional)</span></label>
            <div class="input-group">
              <span class="input-group-text text-muted small" style="font-size: 0.8rem;">
                {{ baseSlug }}/
              </span>
              <input
                v-model="newSlug"
                type="text"
                class="form-control"
                :class="{ 'is-invalid': slugError }"
                placeholder="auto-generated if blank"
                maxlength="20"
                @input="slugError = ''"
              />
              <div v-if="slugError" class="invalid-feedback">{{ slugError }}</div>
            </div>
            <div class="form-text">3–20 alphanumeric characters.</div>
          </div>

          <div v-if="errorMsg" class="alert alert-danger py-2 small">{{ errorMsg }}</div>
        </div>

        <div class="modal-footer border-0 pt-0">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
          <button
            type="button"
            class="btn btn-primary"
            :disabled="cloning"
            @click="doClone"
          >
            <span v-if="cloning" class="spinner-border spinner-border-sm me-1" role="status"></span>
            Clone Link
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { Modal } from 'bootstrap';
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
let bsModal: Modal | null = null;

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
  if (modalEl.value) {
    bsModal = new Modal(modalEl.value);
    modalEl.value.addEventListener('hidden.bs.modal', resetForm);
  }
});

onBeforeUnmount(() => {
  bsModal?.dispose();
});

function show() {
  resetForm();
  bsModal?.show();
}

function hide() {
  bsModal?.hide();
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
