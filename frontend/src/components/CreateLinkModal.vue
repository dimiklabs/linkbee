<template>
  <div ref="modalEl" class="modal fade" tabindex="-1" aria-hidden="true">
    <div class="modal-dialog modal-lg modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title fw-semibold">{{ isEditMode ? 'Edit Link' : 'Create Link' }}</h5>
          <button type="button" class="btn-close" @click="hide" aria-label="Close"></button>
        </div>

        <div class="modal-body">
          <div v-if="error" class="alert alert-danger alert-dismissible" role="alert">
            {{ error }}
            <button type="button" class="btn-close" @click="error = ''" aria-label="Close"></button>
          </div>

          <form @submit.prevent="handleSave" novalidate>
            <!-- Destination URL -->
            <div class="mb-3">
              <label for="destinationUrl" class="form-label fw-medium">
                Destination URL <span class="text-danger">*</span>
              </label>
              <input
                id="destinationUrl"
                v-model="form.destination_url"
                type="url"
                class="form-control"
                :class="{ 'is-invalid': validationErrors.destination_url }"
                placeholder="https://example.com/your-long-url"
                required
              />
              <div v-if="validationErrors.destination_url" class="invalid-feedback">
                {{ validationErrors.destination_url }}
              </div>
            </div>

            <!-- Custom Slug -->
            <div class="mb-3">
              <label for="slug" class="form-label fw-medium">Custom Slug</label>
              <input
                id="slug"
                v-model="form.slug"
                type="text"
                class="form-control"
                placeholder="my-custom-slug"
                :readonly="isEditMode"
                :class="{ 'bg-light text-muted': isEditMode }"
              />
              <div class="form-text">
                {{ isEditMode ? 'Slug cannot be changed after creation.' : 'Leave blank to auto-generate.' }}
              </div>
            </div>

            <!-- Title -->
            <div class="mb-3">
              <label for="title" class="form-label fw-medium">Title</label>
              <input
                id="title"
                v-model="form.title"
                type="text"
                class="form-control"
                placeholder="My Link Title"
              />
            </div>

            <!-- Password -->
            <div class="mb-3">
              <label for="password" class="form-label fw-medium">Password</label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                class="form-control"
                placeholder="Leave blank for no password"
                autocomplete="new-password"
              />
            </div>

            <div class="row g-3 mb-3">
              <!-- Expires At -->
              <div class="col-md-6">
                <label for="expiresAt" class="form-label fw-medium">Expires At</label>
                <input
                  id="expiresAt"
                  v-model="form.expires_at"
                  type="datetime-local"
                  class="form-control"
                />
              </div>

              <!-- Max Clicks -->
              <div class="col-md-6">
                <label for="maxClicks" class="form-label fw-medium">Max Clicks</label>
                <input
                  id="maxClicks"
                  v-model.number="form.max_clicks"
                  type="number"
                  class="form-control"
                  min="1"
                  placeholder="Unlimited"
                />
              </div>
            </div>

            <!-- Redirect Type -->
            <div class="mb-3">
              <label for="redirectType" class="form-label fw-medium">Redirect Type</label>
              <select id="redirectType" v-model.number="form.redirect_type" class="form-select">
                <option :value="302">302 — Temporary Redirect</option>
                <option :value="301">301 — Permanent Redirect</option>
              </select>
            </div>

            <!-- Tags -->
            <div class="mb-3">
              <label for="tags" class="form-label fw-medium">Tags</label>
              <input
                id="tags"
                v-model="tagsInput"
                type="text"
                class="form-control"
                placeholder="marketing, social, campaign"
              />
              <div class="form-text">Comma-separated list of tags.</div>
            </div>

            <!-- UTM Parameters (collapsible) -->
            <div class="accordion accordion-flush border rounded mb-3" id="utmAccordion">
              <div class="accordion-item">
                <h2 class="accordion-header">
                  <button
                    class="accordion-button collapsed fw-medium"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target="#utmCollapse"
                    aria-expanded="false"
                    aria-controls="utmCollapse"
                  >
                    UTM Parameters
                    <span class="badge bg-secondary ms-2 fw-normal" style="font-size: 0.7rem;">Optional</span>
                  </button>
                </h2>
                <div id="utmCollapse" class="accordion-collapse collapse" data-bs-parent="#utmAccordion">
                  <div class="accordion-body">
                    <div class="row g-3">
                      <div class="col-md-4">
                        <label for="utmSource" class="form-label fw-medium">UTM Source</label>
                        <input
                          id="utmSource"
                          v-model="form.utm_source"
                          type="text"
                          class="form-control"
                          placeholder="google"
                        />
                      </div>
                      <div class="col-md-4">
                        <label for="utmMedium" class="form-label fw-medium">UTM Medium</label>
                        <input
                          id="utmMedium"
                          v-model="form.utm_medium"
                          type="text"
                          class="form-control"
                          placeholder="cpc"
                        />
                      </div>
                      <div class="col-md-4">
                        <label for="utmCampaign" class="form-label fw-medium">UTM Campaign</label>
                        <input
                          id="utmCampaign"
                          v-model="form.utm_campaign"
                          type="text"
                          class="form-control"
                          placeholder="spring_sale"
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary" @click="hide" :disabled="saving">
            Cancel
          </button>
          <button type="button" class="btn btn-primary" @click="handleSave" :disabled="saving">
            <span v-if="saving" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
            {{ saving ? 'Saving...' : (isEditMode ? 'Save Changes' : 'Create Link') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue';
import { Modal } from 'bootstrap';
import type { LinkResponse, CreateLinkRequest, UpdateLinkRequest } from '@/types/links';
import { useLinksStore } from '@/stores/links';

interface Props {
  link?: LinkResponse;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  saved: [link: LinkResponse];
}>();

const linksStore = useLinksStore();

const modalEl = ref<HTMLElement | null>(null);
let modalInstance: Modal | null = null;

const saving = ref(false);
const error = ref('');
const validationErrors = ref<Record<string, string>>({});
const tagsInput = ref('');

const isEditMode = computed(() => !!props.link);

interface FormState {
  destination_url: string;
  slug: string;
  title: string;
  password: string;
  expires_at: string;
  max_clicks: number | null;
  redirect_type: 301 | 302;
  utm_source: string;
  utm_medium: string;
  utm_campaign: string;
}

const defaultForm = (): FormState => ({
  destination_url: '',
  slug: '',
  title: '',
  password: '',
  expires_at: '',
  max_clicks: null,
  redirect_type: 302,
  utm_source: '',
  utm_medium: '',
  utm_campaign: '',
});

const form = ref<FormState>(defaultForm());

function populateForm(link: LinkResponse) {
  form.value.destination_url = link.destination_url;
  form.value.slug = link.slug;
  form.value.title = link.title ?? '';
  form.value.password = '';
  form.value.expires_at = link.expires_at
    ? new Date(link.expires_at).toISOString().slice(0, 16)
    : '';
  form.value.max_clicks = link.max_clicks ?? null;
  form.value.redirect_type = (link.redirect_type as 301 | 302) ?? 302;
  form.value.utm_source = link.utm_source ?? '';
  form.value.utm_medium = link.utm_medium ?? '';
  form.value.utm_campaign = link.utm_campaign ?? '';
  tagsInput.value = link.tags?.join(', ') ?? '';
}

function resetForm() {
  form.value = defaultForm();
  tagsInput.value = '';
  error.value = '';
  validationErrors.value = {};
}

function validate(): boolean {
  validationErrors.value = {};
  if (!form.value.destination_url.trim()) {
    validationErrors.value.destination_url = 'Destination URL is required.';
    return false;
  }
  try {
    new URL(form.value.destination_url.trim());
  } catch {
    validationErrors.value.destination_url = 'Please enter a valid URL.';
    return false;
  }
  return true;
}

async function handleSave() {
  if (!validate()) return;

  saving.value = true;
  error.value = '';

  const tags = tagsInput.value
    .split(',')
    .map((t) => t.trim())
    .filter((t) => t.length > 0);

  try {
    let savedLink: LinkResponse;

    if (isEditMode.value && props.link) {
      const payload: UpdateLinkRequest = {
        destination_url: form.value.destination_url.trim() || undefined,
        title: form.value.title.trim() || undefined,
        password: form.value.password || undefined,
        expires_at: form.value.expires_at
          ? new Date(form.value.expires_at).toISOString()
          : null,
        max_clicks: form.value.max_clicks ?? undefined,
        redirect_type: form.value.redirect_type,
        tags: tags.length > 0 ? tags : undefined,
        utm_source: form.value.utm_source.trim() || undefined,
        utm_medium: form.value.utm_medium.trim() || undefined,
        utm_campaign: form.value.utm_campaign.trim() || undefined,
      };
      savedLink = await linksStore.updateLink(props.link.id, payload);
    } else {
      const payload: CreateLinkRequest = {
        destination_url: form.value.destination_url.trim(),
        slug: form.value.slug.trim() || undefined,
        title: form.value.title.trim() || undefined,
        password: form.value.password || undefined,
        expires_at: form.value.expires_at
          ? new Date(form.value.expires_at).toISOString()
          : undefined,
        max_clicks: form.value.max_clicks ?? undefined,
        redirect_type: form.value.redirect_type,
        tags: tags.length > 0 ? tags : undefined,
        utm_source: form.value.utm_source.trim() || undefined,
        utm_medium: form.value.utm_medium.trim() || undefined,
        utm_campaign: form.value.utm_campaign.trim() || undefined,
      };
      savedLink = await linksStore.createLink(payload);
    }

    emit('saved', savedLink);
    hide();
  } catch (err: unknown) {
    if (err instanceof Error) {
      error.value = err.message;
    } else {
      error.value = 'An unexpected error occurred. Please try again.';
    }
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  if (modalEl.value) {
    modalInstance = new Modal(modalEl.value, { backdrop: 'static' });

    modalEl.value.addEventListener('hidden.bs.modal', () => {
      resetForm();
    });
  }

  if (props.link) {
    populateForm(props.link);
  }
});

watch(
  () => props.link,
  (newLink) => {
    if (newLink) {
      populateForm(newLink);
    } else {
      resetForm();
    }
  },
  { deep: true }
);

function show() {
  if (props.link) {
    populateForm(props.link);
  } else {
    resetForm();
  }
  modalInstance?.show();
}

function hide() {
  modalInstance?.hide();
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

.accordion-button:not(.collapsed) {
  color: #635bff;
  background-color: rgba(99, 91, 255, 0.05);
  box-shadow: inset 0 -1px 0 rgba(99, 91, 255, 0.15);
}

.accordion-button:focus {
  box-shadow: 0 0 0 0.25rem rgba(99, 91, 255, 0.25);
}
</style>
