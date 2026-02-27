<template>
  <md-dialog :open="isOpen" @closed="onDialogClosed" style="--md-dialog-container-shape:16px">
    <div slot="headline">{{ isEditMode ? 'Edit Link' : 'Create Link' }}</div>

    <div slot="content" style="min-width:540px;max-width:100%;padding:0 4px">
      <!-- Error banner -->
      <div v-if="error" class="create-error-banner">
        <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-error)">error</span>
        <span style="flex:1;font-size:0.875rem">{{ error }}</span>
        <md-icon-button @click="error = ''" style="width:32px;height:32px">
          <span class="material-symbols-outlined" style="font-size:18px">close</span>
        </md-icon-button>
      </div>

      <!-- Destination URL -->
      <div class="field-group">
        <md-outlined-text-field
          :value="form.destination_url"
          @input="form.destination_url = ($event.target as HTMLInputElement).value"
          label="Destination URL *"
          type="url"
          placeholder="https://example.com/your-long-url"
          style="width:100%"
          :error="!!validationErrors.destination_url"
          :error-text="validationErrors.destination_url"
        />
        <!-- Duplicate check feedback -->
        <div v-if="!isEditMode && checkingDuplicate" class="dup-checking">
          <md-circular-progress indeterminate style="--md-circular-progress-size:16px" />
          <span>Checking for duplicates…</span>
        </div>
        <div v-else-if="!isEditMode && duplicateLink && !ignoreDuplicate" class="dup-warning">
          <span class="material-symbols-outlined" style="font-size:18px;color:#F4A100">warning</span>
          <div style="flex:1;font-size:0.85rem">
            <strong>Duplicate URL detected.</strong>
            This destination is already shortened as
            <a :href="duplicateLink.short_url" target="_blank" rel="noopener noreferrer" style="color:var(--md-sys-color-primary);font-weight:600">
              {{ duplicateLink.short_url }}
            </a>
            <span style="color:var(--md-sys-color-on-surface-variant)"> (slug: <code>{{ duplicateLink.slug }}</code>)</span>.
          </div>
          <md-text-button @click="ignoreDuplicate = true">Create anyway</md-text-button>
        </div>
      </div>

      <!-- Custom Slug -->
      <div class="field-group">
        <md-outlined-text-field
          :value="form.slug"
          @input="form.slug = ($event.target as HTMLInputElement).value"
          label="Custom Slug"
          placeholder="my-custom-slug"
          :disabled="isEditMode"
          style="width:100%"
          :supporting-text="isEditMode ? 'Slug cannot be changed after creation.' : 'Leave blank to auto-generate.'"
        />
      </div>

      <!-- Title -->
      <div class="field-group">
        <md-outlined-text-field
          :value="form.title"
          @input="form.title = ($event.target as HTMLInputElement).value"
          label="Title"
          placeholder="My Link Title"
          style="width:100%"
        />
      </div>

      <!-- Password -->
      <div class="field-group">
        <md-outlined-text-field
          :value="form.password"
          @input="form.password = ($event.target as HTMLInputElement).value"
          label="Password"
          type="password"
          placeholder="Leave blank for no password"
          autocomplete="new-password"
          style="width:100%"
        />
      </div>

      <!-- Expires At / Max Clicks row -->
      <div class="field-row">
        <md-outlined-text-field
          :value="form.expires_at"
          @input="form.expires_at = ($event.target as HTMLInputElement).value"
          label="Expires At"
          type="datetime-local"
          style="flex:1"
        />
        <md-outlined-text-field
          :value="form.max_clicks !== null ? String(form.max_clicks) : ''"
          @input="form.max_clicks = ($event.target as HTMLInputElement).value ? Number(($event.target as HTMLInputElement).value) : null"
          label="Max Clicks"
          type="number"
          placeholder="Unlimited"
          style="flex:1"
        />
      </div>

      <!-- Redirect Type -->
      <div class="field-group">
        <md-outlined-select
          :value="String(form.redirect_type)"
          @change="form.redirect_type = Number(($event.target as HTMLSelectElement).value) as 301 | 302"
          label="Redirect Type"
          style="width:100%"
        >
          <md-select-option value="302"><div slot="headline">302 — Temporary Redirect</div></md-select-option>
          <md-select-option value="301"><div slot="headline">301 — Permanent Redirect</div></md-select-option>
        </md-outlined-select>
      </div>

      <!-- Folder -->
      <div class="field-group">
        <md-outlined-select
          :value="form.folder_id"
          @change="form.folder_id = ($event.target as HTMLSelectElement).value"
          label="Folder"
          style="width:100%"
        >
          <md-select-option value=""><div slot="headline">— No folder —</div></md-select-option>
          <md-select-option v-for="f in props.folders" :key="f.id" :value="f.id">
            <div slot="headline">{{ f.name }}</div>
          </md-select-option>
        </md-outlined-select>
      </div>

      <!-- Tags -->
      <div class="field-group">
        <md-outlined-text-field
          :value="tagsInput"
          @input="tagsInput = ($event.target as HTMLInputElement).value"
          label="Tags"
          placeholder="marketing, social, campaign"
          style="width:100%"
          supporting-text="Comma-separated list of tags."
        />
      </div>

      <!-- UTM Parameters (collapsible) -->
      <div class="utm-section">
        <button class="utm-toggle" type="button" @click="utmExpanded = !utmExpanded">
          <span>UTM Parameters</span>
          <span class="m3-badge m3-badge--neutral" style="font-size:0.7rem">Optional</span>
          <span class="material-symbols-outlined utm-chevron" :class="{ 'utm-chevron--open': utmExpanded }">expand_more</span>
        </button>
        <div v-if="utmExpanded" class="utm-fields">
          <div class="field-row">
            <md-outlined-text-field
              :value="form.utm_source"
              @input="form.utm_source = ($event.target as HTMLInputElement).value"
              label="UTM Source"
              placeholder="google"
              style="flex:1"
            />
            <md-outlined-text-field
              :value="form.utm_medium"
              @input="form.utm_medium = ($event.target as HTMLInputElement).value"
              label="UTM Medium"
              placeholder="cpc"
              style="flex:1"
            />
            <md-outlined-text-field
              :value="form.utm_campaign"
              @input="form.utm_campaign = ($event.target as HTMLInputElement).value"
              label="UTM Campaign"
              placeholder="spring_sale"
              style="flex:1"
            />
          </div>
        </div>
      </div>
    </div>

    <div slot="actions">
      <md-text-button @click="hide" :disabled="saving">Cancel</md-text-button>
      <md-filled-button @click="handleSave" :disabled="saving">
        <md-circular-progress v-if="saving" indeterminate style="--md-circular-progress-size:18px" slot="icon" />
        {{ saving ? 'Saving...' : (isEditMode ? 'Save Changes' : 'Create Link') }}
      </md-filled-button>
    </div>
  </md-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue';
import type { LinkResponse, CreateLinkRequest, UpdateLinkRequest } from '@/types/links';
import type { FolderResponse } from '@/types/folders';
import { useLinksStore } from '@/stores/links';
import linksApi from '@/api/links';

interface Props {
  link?: LinkResponse;
  folders?: FolderResponse[];
}

const props = defineProps<Props>();
const emit = defineEmits<{
  saved: [link: LinkResponse];
}>();

const linksStore = useLinksStore();

const modalEl = ref<HTMLElement | null>(null);
let modalInstance: any = null;

const isOpen = ref(false);
const saving = ref(false);
const error = ref('');
const validationErrors = ref<Record<string, string>>({});
const tagsInput = ref('');
const utmExpanded = ref(false);

// Duplicate detection
const duplicateLink = ref<LinkResponse | null>(null);
const checkingDuplicate = ref(false);
const ignoreDuplicate = ref(false);
let duplicateTimer: ReturnType<typeof setTimeout> | null = null;

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
  folder_id: string;
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
  folder_id: '',
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
  form.value.folder_id = link.folder_id ?? '';
  tagsInput.value = link.tags?.join(', ') ?? '';
}

function resetForm() {
  form.value = defaultForm();
  tagsInput.value = '';
  error.value = '';
  validationErrors.value = {};
  duplicateLink.value = null;
  ignoreDuplicate.value = false;
  utmExpanded.value = false;
  if (duplicateTimer) { clearTimeout(duplicateTimer); duplicateTimer = null; }
}

async function runDuplicateCheck(url: string) {
  checkingDuplicate.value = true;
  try {
    duplicateLink.value = await linksApi.checkDuplicate(url);
  } finally {
    checkingDuplicate.value = false;
  }
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

// Watch destination URL to trigger duplicate check in create mode
watch(
  () => form.value.destination_url,
  (url) => {
    if (isEditMode.value) return;
    duplicateLink.value = null;
    ignoreDuplicate.value = false;
    if (duplicateTimer) clearTimeout(duplicateTimer);
    if (!url.trim()) return;
    try { new URL(url.trim()); } catch { return; }
    duplicateTimer = setTimeout(() => runDuplicateCheck(url.trim()), 600);
  }
);

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
        folder_id: form.value.folder_id || null,
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
        folder_id: form.value.folder_id || undefined,
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
  // Bootstrap modal lifecycle removed — component will be rewritten for Vuetify
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
  isOpen.value = true;
  modalInstance?.show();
}

function hide() {
  isOpen.value = false;
  modalInstance?.hide();
}

function onDialogClosed() {
  isOpen.value = false;
}

defineExpose({ show, hide });
</script>

<style scoped>
.field-group {
  margin-bottom: 16px;
}

.field-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.create-error-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  margin-bottom: 16px;
  background: var(--md-sys-color-error-container, #FFDAD6);
  color: var(--md-sys-color-on-error-container, #410002);
  border-radius: 8px;
}

.dup-checking {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
}

.dup-warning {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 8px;
  padding: 10px 12px;
  background: rgba(244, 161, 0, 0.1);
  border: 1px solid rgba(244, 161, 0, 0.4);
  border-radius: 8px;
}

/* UTM section */
.utm-section {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 8px;
}

.utm-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 12px 16px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  text-align: left;
}

.utm-toggle:hover {
  background: var(--md-sys-color-surface-container-low);
}

.utm-chevron {
  margin-left: auto;
  transition: transform 0.2s;
}

.utm-chevron--open {
  transform: rotate(180deg);
}

.utm-fields {
  padding: 12px 16px 16px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
}
</style>
