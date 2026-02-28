<template>
  <BaseModal
    v-model="isOpen"
    :size="modalSize"
    @closed="onDialogClosed"
  >
    <template #headline>
      <span class="material-symbols-outlined dialog-headline-icon">{{ headlineIcon }}</span>
      <span>{{ headlineText }}</span>
    </template>

    <div class="dialog-content">

      <!-- Error banner -->
      <div v-if="error" class="alert-banner alert-banner--error">
        <span class="material-symbols-outlined alert-icon">error</span>
        <span class="alert-text">{{ error }}</span>
        <button class="btn-icon alert-close" @click="error = ''">
          <span class="material-symbols-outlined">close</span>
        </button>
      </div>

      <!-- ═══════════════════════════════════════════════════════ -->
      <!-- STEP 1 — Enter URL (create mode only)                  -->
      <!-- ═══════════════════════════════════════════════════════ -->
      <template v-if="!isEditMode && step === 'create'">
        <div class="step1-wrap">

          <div class="field-group">
            <md-outlined-text-field
              :value="form.destination_url"
              @input="form.destination_url = ($event.target as HTMLInputElement).value"
              label="Paste your long URL"
              type="url"
              placeholder="https://example.com/your-very-long-url"
              class="field-full step1-url-field"
              :error="!!validationErrors.destination_url"
              :error-text="validationErrors.destination_url"
            >
              <span class="material-symbols-outlined" slot="leading-icon">link</span>
            </md-outlined-text-field>

            <!-- Duplicate check -->
            <div v-if="checkingDuplicate" class="feedback-row feedback-row--muted">
              <md-circular-progress indeterminate style="flex-shrink:0" />
              <span>Checking for duplicate URLs…</span>
            </div>
            <div v-else-if="duplicateLink && !ignoreDuplicate" class="dup-warning">
              <span class="material-symbols-outlined dup-warning-icon">warning</span>
              <div class="dup-warning-body">
                <strong>Duplicate URL detected.</strong>
                Already shortened as
                <a :href="duplicateLink.short_url" target="_blank" rel="noopener noreferrer" class="dup-link">
                  {{ duplicateLink.short_url }}
                </a>.
              </div>
              <button class="btn-text dup-ignore-btn" @click="ignoreDuplicate = true">Create anyway</button>
            </div>
          </div>

          <!-- Short link preview -->
          <div v-if="form.destination_url && !validationErrors.destination_url" class="short-link-preview">
            <span class="material-symbols-outlined short-link-preview-icon">bolt</span>
            <span class="short-link-preview-label">Your short link</span>
            <div class="short-link-preview-pill">
              <span class="short-link-preview-domain">sl.ink/</span>
              <span class="short-link-preview-slug">{{ form.slug || '·····' }}</span>
            </div>
          </div>

          <!-- Custom slug (secondary / optional) -->
          <div class="step1-slug-wrap">
            <div class="step1-slug-label">
              <span class="material-symbols-outlined" style="font-size:14px">tag</span>
              Custom slug
              <span class="optional-badge">optional</span>
            </div>
            <md-outlined-text-field
              :value="form.slug"
              @input="form.slug = ($event.target as HTMLInputElement).value"
              label="Custom slug"
              placeholder="my-custom-slug"
              class="field-full"
              supporting-text="Leave blank to auto-generate a short code."
            >
              <span class="material-symbols-outlined" slot="leading-icon">tag</span>
            </md-outlined-text-field>
          </div>

        </div>
      </template>

      <!-- ═══════════════════════════════════════════════════════ -->
      <!-- STEP 2 — Customize (after creation OR in edit mode)    -->
      <!-- ═══════════════════════════════════════════════════════ -->
      <template v-else>

        <!-- Created link success card (only right after creation) -->
        <div v-if="!isEditMode && createdLink" class="created-card">
          <div class="created-card-icon">
            <span class="material-symbols-outlined">check_circle</span>
          </div>
          <div class="created-card-body">
            <span class="created-card-label">Short link ready</span>
            <a :href="createdLink.short_url" target="_blank" rel="noopener noreferrer" class="created-card-url">
              {{ createdLink.short_url }}
            </a>
          </div>
          <button class="btn-icon created-card-copy" @click="copyShortLink" :title="copied ? 'Copied!' : 'Copy'">
            <span class="material-symbols-outlined">{{ copied ? 'check' : 'content_copy' }}</span>
          </button>
        </div>

        <!-- Edit mode: destination URL editable -->
        <div v-if="isEditMode" class="form-section">
          <div class="section-label">
            <span class="material-symbols-outlined section-label-icon">travel_explore</span>
            Destination URL
          </div>
          <md-outlined-text-field
            :value="form.destination_url"
            @input="form.destination_url = ($event.target as HTMLInputElement).value"
            label="Destination URL"
            type="url"
            class="field-full"
            :error="!!validationErrors.destination_url"
            :error-text="validationErrors.destination_url"
          >
            <span class="material-symbols-outlined" slot="leading-icon">link</span>
          </md-outlined-text-field>
        </div>

        <!-- Section: Details (Title, Redirect, Tags) -->
        <div class="form-section">
          <div class="section-label">
            <span class="material-symbols-outlined section-label-icon">badge</span>
            Details
          </div>

          <div class="field-row" style="margin-bottom:12px">
            <md-outlined-text-field
              :value="form.title"
              @input="form.title = ($event.target as HTMLInputElement).value"
              label="Title"
              placeholder="My Link Title"
              class="field-flex"
            >
              <span class="material-symbols-outlined" slot="leading-icon">title</span>
            </md-outlined-text-field>

            <AppSelect
              :model-value="String(form.redirect_type)"
              @update:model-value="form.redirect_type = Number($event) as 301 | 302"
              label="Redirect Type"
              class="field-flex"
            >
              <option value="302">302 — Temporary</option>
              <option value="301">301 — Permanent</option>
            </AppSelect>
          </div>

          <md-outlined-text-field
            :value="tagsInput"
            @input="tagsInput = ($event.target as HTMLInputElement).value"
            label="Tags"
            placeholder="marketing, social, campaign"
            class="field-full"
            supporting-text="Comma-separated tags to organise your links."
          >
            <span class="material-symbols-outlined" slot="leading-icon">sell</span>
          </md-outlined-text-field>
          <div v-if="parsedTags.length > 0" class="tag-chips">
            <span v-for="tag in parsedTags" :key="tag" class="tag-chip">
              <span class="material-symbols-outlined tag-chip-icon">label</span>
              {{ tag }}
            </span>
          </div>
        </div>

        <!-- UTM Parameters (collapsible — open by default after creation) -->
        <div class="advanced-section">
          <button class="advanced-toggle" type="button" @click="utmExpanded = !utmExpanded">
            <span class="material-symbols-outlined advanced-toggle-icon">analytics</span>
            <span class="advanced-toggle-label">UTM Parameters</span>
            <span class="adv-badge adv-badge--optional">Optional</span>
            <span class="material-symbols-outlined advanced-chevron" :class="{ 'advanced-chevron--open': utmExpanded }">
              expand_more
            </span>
          </button>
          <div v-if="utmExpanded" class="advanced-fields">
            <div class="field-row">
              <md-outlined-text-field
                :value="form.utm_source"
                @input="form.utm_source = ($event.target as HTMLInputElement).value"
                label="UTM Source"
                placeholder="google"
                class="field-flex"
                supporting-text="e.g. google, newsletter"
              />
              <md-outlined-text-field
                :value="form.utm_medium"
                @input="form.utm_medium = ($event.target as HTMLInputElement).value"
                label="UTM Medium"
                placeholder="cpc"
                class="field-flex"
                supporting-text="e.g. cpc, email, social"
              />
              <md-outlined-text-field
                :value="form.utm_campaign"
                @input="form.utm_campaign = ($event.target as HTMLInputElement).value"
                label="UTM Campaign"
                placeholder="spring_sale"
                class="field-flex"
                supporting-text="e.g. spring_sale"
              />
            </div>
          </div>
        </div>

        <!-- Advanced Options (collapsible) -->
        <div class="advanced-section" style="margin-top:8px">
          <button class="advanced-toggle" type="button" @click="advancedExpanded = !advancedExpanded">
            <span class="material-symbols-outlined advanced-toggle-icon">tune</span>
            <span class="advanced-toggle-label">Advanced Options</span>
            <span class="adv-badge">Password · Expiry · Limits</span>
            <span class="material-symbols-outlined advanced-chevron" :class="{ 'advanced-chevron--open': advancedExpanded }">
              expand_more
            </span>
          </button>
          <div v-if="advancedExpanded" class="advanced-fields">
            <AppSelect
              v-model="form.folder_id"
              label="Folder"
              class="field-full"
            >
              <option value="">— No folder —</option>
              <option v-for="f in props.folders" :key="f.id" :value="f.id">{{ f.name }}</option>
            </AppSelect>

            <md-outlined-text-field
              :value="form.password"
              @input="form.password = ($event.target as HTMLInputElement).value"
              label="Link Password"
              type="password"
              placeholder="Leave blank for no password"
              autocomplete="new-password"
              class="field-full"
              supporting-text="Visitors will need this password to access the destination."
            >
              <span class="material-symbols-outlined" slot="leading-icon">lock</span>
            </md-outlined-text-field>

            <div class="field-row">
              <md-outlined-text-field
                :value="form.expires_at"
                @input="form.expires_at = ($event.target as HTMLInputElement).value"
                label="Expires At"
                type="datetime-local"
                class="field-flex"
                supporting-text="Leave blank to never expire."
              >
                <span class="material-symbols-outlined" slot="leading-icon">event</span>
              </md-outlined-text-field>

              <md-outlined-text-field
                :value="form.max_clicks !== null ? String(form.max_clicks) : ''"
                @input="form.max_clicks = ($event.target as HTMLInputElement).value ? Number(($event.target as HTMLInputElement).value) : null"
                label="Max Clicks"
                type="number"
                placeholder="Unlimited"
                class="field-flex"
                supporting-text="Disable after this many clicks."
              >
                <span class="material-symbols-outlined" slot="leading-icon">ads_click</span>
              </md-outlined-text-field>
            </div>
          </div>
        </div>

      </template>
    </div>

    <template #actions>
      <!-- Step 1 actions -->
      <template v-if="!isEditMode && step === 'create'">
        <button class="btn-text" @click="hide" :disabled="saving">Cancel</button>
        <button
          class="btn-filled save-btn"
          @click="handleCreate"
          :disabled="saving || (!ignoreDuplicate && !!duplicateLink)"
        >
          <md-circular-progress v-if="saving" indeterminate style="margin-right:6px;" />
          <span v-else class="material-symbols-outlined" style="margin-right:6px;font-size:18px;">bolt</span>
          {{ saving ? 'Creating…' : 'Shorten Link' }}
        </button>
      </template>

      <!-- Step 2 / Edit mode actions -->
      <template v-else>
        <button class="btn-text" @click="handleDone" :disabled="saving">
          {{ isEditMode ? 'Cancel' : 'Done' }}
        </button>
        <button class="btn-filled save-btn" @click="handleUpdate" :disabled="saving">
          <md-circular-progress v-if="saving" indeterminate style="margin-right:6px;" />
          <span v-else class="material-symbols-outlined" style="margin-right:6px;font-size:18px;">save</span>
          {{ saving ? 'Saving…' : 'Save Changes' }}
        </button>
      </template>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { LinkResponse, CreateLinkRequest, UpdateLinkRequest } from '@/types/links';
import type { FolderResponse } from '@/types/folders';
import { useLinksStore } from '@/stores/links';
import linksApi from '@/api/links';
import BaseModal from '@/components/BaseModal.vue';
import AppSelect from '@/components/AppSelect.vue';

interface Props {
  link?: LinkResponse;
  folders?: FolderResponse[];
}

const props = defineProps<Props>();
const emit = defineEmits<{
  saved: [link: LinkResponse];
}>();

const linksStore = useLinksStore();

const isOpen = ref(false);
const saving = ref(false);
const error = ref('');
const validationErrors = ref<Record<string, string>>({});
const tagsInput = ref('');
const utmExpanded = ref(false);
const advancedExpanded = ref(false);
const step = ref<'create' | 'customize'>('create');
const createdLink = ref<LinkResponse | null>(null);
const copied = ref(false);

// Duplicate detection
const duplicateLink = ref<LinkResponse | null>(null);
const checkingDuplicate = ref(false);
const ignoreDuplicate = ref(false);
let duplicateTimer: ReturnType<typeof setTimeout> | null = null;

const isEditMode = computed(() => !!props.link);

const modalSize = computed(() =>
  !isEditMode.value && step.value === 'create' ? 'md' : 'lg'
);

const headlineIcon = computed(() => {
  if (isEditMode.value) return 'edit_square';
  return step.value === 'create' ? 'add_link' : 'check_circle';
});

const headlineText = computed(() => {
  if (isEditMode.value) return 'Edit Link';
  return step.value === 'create' ? 'Shorten a Link' : 'Link Shortened!';
});

const parsedTags = computed(() =>
  tagsInput.value
    .split(',')
    .map((t) => t.trim())
    .filter((t) => t.length > 0)
);

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
  utmExpanded.value = !!(link.utm_source || link.utm_medium || link.utm_campaign);
}

function resetForm() {
  form.value = defaultForm();
  tagsInput.value = '';
  error.value = '';
  validationErrors.value = {};
  duplicateLink.value = null;
  ignoreDuplicate.value = false;
  utmExpanded.value = false;
  advancedExpanded.value = false;
  createdLink.value = null;
  step.value = 'create';
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

async function handleCreate() {
  if (!validate()) return;
  saving.value = true;
  error.value = '';
  try {
    const payload: CreateLinkRequest = {
      destination_url: form.value.destination_url.trim(),
      slug: form.value.slug.trim() || undefined,
      redirect_type: 302,
    };
    createdLink.value = await linksStore.createLink(payload);
    step.value = 'customize';
    utmExpanded.value = true; // open UTM by default so user can add tracking
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Failed to create link.';
  } finally {
    saving.value = false;
  }
}

async function handleUpdate() {
  if (isEditMode.value && !validate()) return;
  saving.value = true;
  error.value = '';
  const tags = parsedTags.value;
  const linkId = isEditMode.value ? props.link!.id : createdLink.value!.id;
  try {
    const payload: UpdateLinkRequest = {
      destination_url: isEditMode.value ? form.value.destination_url.trim() || undefined : undefined,
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
    const savedLink = await linksStore.updateLink(linkId, payload);
    emit('saved', savedLink);
    hide();
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Failed to save changes.';
  } finally {
    saving.value = false;
  }
}

function handleDone() {
  // Close step 2 without saving extra fields; emit the already-created link
  if (createdLink.value) {
    emit('saved', createdLink.value);
  }
  hide();
}

async function copyShortLink() {
  const url = createdLink.value?.short_url;
  if (!url) return;
  try {
    await navigator.clipboard.writeText(url);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 2000);
  } catch {
    // Clipboard API unavailable
  }
}

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
    step.value = 'customize';
  } else {
    resetForm();
    step.value = 'create';
  }
  isOpen.value = true;
}

function hide() {
  isOpen.value = false;
}

function onDialogClosed() {
  isOpen.value = false;
}

defineExpose({ show, hide });
</script>

<style scoped lang="scss">
/* ── Headline icon ────────────────────────────────────── */
.dialog-headline-icon {
  font-size: 22px;
  color: var(--md-sys-color-primary);
}

/* ── Content wrapper ──────────────────────────────────── */
.dialog-content {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* ── Alerts ───────────────────────────────────────────── */
.alert-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 12px;
  margin-bottom: 20px;
  font-size: 0.875rem;
  line-height: 1.4;

  &--error {
    background: var(--md-sys-color-error-container, #ffdad6);
    color: var(--md-sys-color-on-error-container, #410002);
    border: 1px solid color-mix(in srgb, var(--md-sys-color-error, #ba1a1a) 30%, transparent);
  }
}

.alert-icon { font-size: 18px; flex-shrink: 0; }
.alert-text { flex: 1; }
.alert-close { width: 32px; height: 32px; flex-shrink: 0; }

/* ── STEP 1 layout ────────────────────────────────────── */
.step1-wrap {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.step1-url-field {
  --md-outlined-text-field-container-shape: 14px;
}

.step1-slug-wrap {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
}

.step1-slug-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.optional-badge {
  font-size: 0.65rem;
  font-weight: 500;
  padding: 1px 7px;
  border-radius: 10px;
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-on-surface-variant);
  text-transform: none;
  letter-spacing: 0;
}

/* ── Created link card ────────────────────────────────── */
.created-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: color-mix(in srgb, var(--md-sys-color-primary, #635bff) 8%, transparent);
  border: 1px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 25%, transparent);
  border-radius: 14px;
  margin-bottom: 20px;
}

.created-card-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--md-sys-color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined {
    font-size: 20px;
    color: var(--md-sys-color-on-primary);
  }
}

.created-card-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.created-card-label {
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--md-sys-color-primary);
}

.created-card-url {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  text-decoration: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

  &:hover { text-decoration: underline; }
}

.created-card-copy {
  flex-shrink: 0;
  color: var(--md-sys-color-primary);
}

/* ── Form sections ────────────────────────────────────── */
.form-section {
  padding: 16px 0;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &:first-of-type { padding-top: 4px; }
  &:last-of-type  { border-bottom: none; }
}

.section-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--md-sys-color-primary);
  margin-bottom: 12px;
}

.section-label-icon { font-size: 14px; }

/* ── Fields ───────────────────────────────────────────── */
.field-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.field-full  { width: 100%; }
.field-flex  { flex: 1; min-width: 0; }

/* ── Feedback rows ────────────────────────────────────── */
.feedback-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.82rem;

  &--muted { color: var(--md-sys-color-on-surface-variant); }
}

/* ── Duplicate warning ────────────────────────────────── */
.dup-warning {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  background: color-mix(in srgb, #f4a100 10%, transparent);
  border: 1px solid color-mix(in srgb, #f4a100 40%, transparent);
  border-radius: 10px;
  font-size: 0.85rem;
}

.dup-warning-icon {
  font-size: 18px;
  color: #f4a100;
  flex-shrink: 0;
  margin-top: 1px;
}

.dup-warning-body { flex: 1; line-height: 1.5; }

.dup-link {
  color: var(--md-sys-color-primary);
  font-weight: 600;
  text-decoration: none;
  &:hover { text-decoration: underline; }
}

.dup-ignore-btn { flex-shrink: 0; align-self: center; }

/* ── Short link preview ───────────────────────────────── */
.short-link-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.short-link-preview-icon {
  font-size: 16px;
  color: var(--md-sys-color-primary);
}

.short-link-preview-label { flex-shrink: 0; }

.short-link-preview-pill {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  background: var(--md-sys-color-primary-container, rgba(99, 91, 255, 0.1));
  border: 1px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 30%, transparent);
  border-radius: 20px;
  font-family: monospace;
  font-size: 0.8rem;
  font-weight: 500;
}

.short-link-preview-domain { color: var(--md-sys-color-on-surface-variant); }
.short-link-preview-slug   { color: var(--md-sys-color-primary); font-weight: 700; }

/* ── Tag chips ────────────────────────────────────────── */
.tag-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.tag-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: var(--md-sys-color-secondary-container, rgba(99, 91, 255, 0.08));
  color: var(--md-sys-color-on-secondary-container, var(--md-sys-color-primary));
  border-radius: 20px;
  font-size: 0.78rem;
  font-weight: 500;
  border: 1px solid color-mix(in srgb, var(--md-sys-color-primary, #635bff) 20%, transparent);
}

.tag-chip-icon { font-size: 12px; opacity: 0.7; }

/* ── Advanced / collapsible sections ─────────────────── */
.advanced-section {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
  margin-bottom: 4px;
}

.advanced-toggle {
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
  transition: background 0.15s;

  &:hover { background: var(--md-sys-color-surface-container-low); }
}

.advanced-toggle-icon {
  font-size: 18px;
  color: var(--md-sys-color-on-surface-variant);
}

.advanced-toggle-label { flex: 1; }

.adv-badge {
  font-size: 0.68rem;
  font-weight: 500;
  padding: 2px 8px;
  border-radius: 10px;
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-on-surface-variant);

  &--optional {
    background: color-mix(in srgb, var(--md-sys-color-tertiary, #7e5260) 12%, transparent);
    color: var(--md-sys-color-tertiary, #7e5260);
  }
}

.advanced-chevron {
  font-size: 20px;
  color: var(--md-sys-color-on-surface-variant);
  transition: transform 0.2s ease;
  margin-left: auto;

  &--open { transform: rotate(180deg); }
}

.advanced-fields {
  padding: 4px 16px 16px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: var(--md-sys-color-surface-container-lowest, transparent);
}

/* ── Save button ──────────────────────────────────────── */
.save-btn { min-width: 140px; }
</style>
