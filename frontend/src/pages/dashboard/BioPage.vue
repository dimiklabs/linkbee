<template>
  <div class="bio-page">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Bio Page</h1>
        <p class="page-subtitle">Share a single page with all your important links.</p>
      </div>
      <div v-if="bioPage" class="page-header__actions">
        <span :class="['m3-badge', bioPage.is_published ? 'm3-badge--success' : 'm3-badge--neutral']" style="padding:6px 14px;">
          {{ bioPage.is_published ? 'Published' : 'Draft' }}
        </span>
        <a v-if="bioPage.is_published" :href="publicUrl" target="_blank" rel="noopener noreferrer">
          <button class="btn-outlined">
            <span class="material-symbols-outlined">open_in_new</span>
            View page
          </button>
        </a>
        <button class="btn-filled" :disabled="saving" @click="saveSettings">
          <span v-if="saving" class="css-spinner css-spinner--sm css-spinner--white"></span>
          Save settings
        </button>
      </div>
    </div>

    <!-- Copy Public URL -->
    <div v-if="bioPage && bioPage.is_published && bioPage.username">
      <button class="btn-text" @click="copyBioUrl">
        <span class="material-symbols-outlined">{{ bioCopied ? 'check' : 'content_copy' }}</span>
        {{ bioCopied ? 'Copied!' : 'Copy public URL' }}
      </button>
    </div>

    <!-- Stats Row -->
    <div v-if="bioPage" class="stats-row">
      <div class="stat-card">
        <div class="stat-label">Total Links</div>
        <div class="stat-value">{{ bioStats.total }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Active Links</div>
        <div class="stat-value stat-value--success">{{ bioStats.active }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Inactive Links</div>
        <div class="stat-value stat-value--muted">{{ bioStats.inactive }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Status</div>
        <span :class="['m3-badge', bioPage.is_published ? 'm3-badge--success' : 'm3-badge--neutral']" style="margin-top:4px;">
          {{ bioPage.is_published ? 'Published' : 'Draft' }}
        </span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <span class="css-spinner"></span>
    </div>

    <div v-else-if="bioPage" class="two-column-layout">

      <!-- Left: Settings Form -->
      <div class="an-card settings-card">
        <div class="an-card-header">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">settings</span>
          </div>
          <span class="an-card-title">Page Settings</span>
        </div>
        <div class="an-card-body">

          <!-- Username -->
          <div class="form-field">
            <label class="form-field__label">Username *</label>
            <div class="bio-username-row">
              <span class="bio-prefix">/bio/</span>
              <input
                class="form-input bio-username-input"
                :value="form.username"
                @input="form.username=($event.target as HTMLInputElement).value"
                placeholder="your-username"
                maxlength="30"
              />
            </div>
            <span class="form-hint">Letters, numbers, underscores only. Min 3 characters.</span>
          </div>

          <!-- Title -->
          <div class="form-field">
            <label class="form-field__label">Page title</label>
            <input
              class="form-input"
              :value="form.title"
              @input="form.title=($event.target as HTMLInputElement).value"
              placeholder="My Links"
              maxlength="100"
            />
          </div>

          <!-- Description -->
          <div class="form-field">
            <label class="form-field__label">Bio / description</label>
            <textarea
              class="form-textarea"
              :value="form.description"
              @input="form.description=($event.target as HTMLTextAreaElement).value"
              placeholder="A short bio or description..."
              rows="3"
            ></textarea>
          </div>

          <!-- Avatar URL -->
          <div class="form-field">
            <label class="form-field__label">Avatar URL</label>
            <input
              class="form-input"
              type="url"
              :value="form.avatar_url"
              @input="form.avatar_url=($event.target as HTMLInputElement).value"
              placeholder="https://..."
            />
          </div>

          <!-- Theme -->
          <div class="form-field">
            <label class="form-field__label">Theme</label>
            <div class="theme-toggle">
              <button :class="['btn-outlined', form.theme === 'light' ? 'btn-active' : '']"
                @click="form.theme = 'light'" style="flex:1;">
                <span class="material-symbols-outlined">light_mode</span>
                Light
              </button>
              <button :class="['btn-outlined', form.theme === 'dark' ? 'btn-active' : '']"
                @click="form.theme = 'dark'" style="flex:1;">
                <span class="material-symbols-outlined">dark_mode</span>
                Dark
              </button>
            </div>
          </div>

          <!-- Published toggle -->
          <div :class="['published-toggle', form.is_published ? 'published-toggle--active' : '']">
            <div>
              <div class="published-toggle__label">Published</div>
              <div class="published-toggle__hint">Make your bio page publicly accessible</div>
            </div>
            <label style="display:flex;align-items:center;cursor:pointer;">
              <input
                v-model="form.is_published"
                type="checkbox"
                role="switch"
                style="width:2.5rem;height:1.25rem;cursor:pointer;accent-color:var(--md-sys-color-primary);"
              />
            </label>
          </div>

          <div v-if="saveError" class="error-box">{{ saveError }}</div>
        </div>
      </div>

      <!-- Right: Links Management -->
      <div class="links-column">

        <!-- Add Link Form -->
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">add_link</span>
            </div>
            <span class="an-card-title">Add Link</span>
          </div>
          <div class="an-card-body">
            <div class="add-link-row">
              <input
                class="form-input"
                :value="newTitle"
                @input="newTitle=($event.target as HTMLInputElement).value"
                placeholder="Link title"
                maxlength="100"
                style="flex:1;min-width:160px;"
                @keydown.enter="addLink"
              />
              <input
                class="form-input"
                type="url"
                :value="newUrl"
                @input="newUrl=($event.target as HTMLInputElement).value"
                placeholder="https://..."
                style="flex:1;min-width:160px;"
                @keydown.enter="addLink"
              />
              <button class="btn-filled"
                :disabled="addingLink || !newTitle.trim() || !newUrl.trim()"
                @click="addLink"
                style="flex-shrink:0;"
              >
                <span v-if="addingLink" class="css-spinner css-spinner--sm css-spinner--white"></span>
                <span v-else class="material-symbols-outlined">add</span>
                Add
              </button>
            </div>
            <div v-if="addError" class="error-box">{{ addError }}</div>
          </div>
        </div>

        <!-- Links List -->
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">list</span>
            </div>
            <span class="an-card-title">Links ({{ links.length }})</span>
            <span class="an-card-hint">Drag to reorder</span>
          </div>

          <div v-if="links.length === 0" class="empty-state">
            <div class="empty-icon">
              <span class="material-symbols-outlined">link_off</span>
            </div>
            <div class="empty-title">No links yet</div>
            <p class="empty-desc">Add your first link using the form above.</p>
          </div>

          <div v-else>
            <div
              v-for="(link, idx) in links"
              :key="link.id"
              class="link-item"
              :style="{ opacity: link.is_active ? '1' : '0.6' }"
              draggable="true"
              @dragstart="onDragStart(idx)"
              @dragover.prevent="onDragOver(idx)"
              @drop="onDrop"
              @dragend="onDragEnd"
            >
              <!-- View mode -->
              <div v-if="editingLinkId !== link.id" class="link-item__view">
                <span class="material-symbols-outlined link-drag">drag_indicator</span>
                <span class="link-pos">{{ idx + 1 }}</span>
                <span
                  class="link-status-dot"
                  :style="{ backgroundColor: link.is_active ? '#16a34a' : 'var(--md-sys-color-outline-variant)' }"
                  :title="link.is_active ? 'Active' : 'Inactive'"
                ></span>
                <div class="link-info">
                  <div class="link-info__title">{{ link.title }}</div>
                  <div class="link-info__url">{{ link.url }}</div>
                </div>
                <input
                  type="checkbox"
                  role="switch"
                  :checked="link.is_active"
                  style="cursor:pointer;accent-color:var(--md-sys-color-primary);flex-shrink:0;"
                  @change="toggleActive(link)"
                />
                <button class="btn-icon" title="Edit" @click="startEdit(link)">
                  <span class="material-symbols-outlined">edit</span>
                </button>
                <button class="btn-icon btn-sm btn-danger" title="Delete" @click="deleteLink(link)">
                  <span class="material-symbols-outlined">delete</span>
                </button>
              </div>

              <!-- Edit mode -->
              <div v-else class="link-item__edit">
                <input
                  class="form-input"
                  :value="editTitle"
                  @input="editTitle=($event.target as HTMLInputElement).value"
                  placeholder="Title"
                  maxlength="100"
                  style="flex:1;min-width:160px;"
                />
                <input
                  class="form-input"
                  type="url"
                  :value="editUrl"
                  @input="editUrl=($event.target as HTMLInputElement).value"
                  placeholder="URL"
                  style="flex:1;min-width:160px;"
                />
                <div class="link-edit-actions">
                  <button class="btn-icon" :disabled="savingEdit" @click="saveEdit(link)">
                    <span v-if="savingEdit" class="css-spinner css-spinner--sm"></span>
                    <span v-else class="material-symbols-outlined">check</span>
                  </button>
                  <button class="btn-icon" @click="cancelEdit">
                    <span class="material-symbols-outlined">close</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import bioApi from '@/api/bio';
import type { BioPage, BioLinkItem } from '@/types/bio';

const loading = ref(true);
const saving = ref(false);
const saveError = ref('');
const bioPage = ref<BioPage | null>(null);
const links = ref<BioLinkItem[]>([]);

const form = ref({
  username: '',
  title: '',
  description: '',
  avatar_url: '',
  theme: 'light' as 'light' | 'dark',
  is_published: false,
});

// Add link
const newTitle = ref('');
const newUrl = ref('');
const addingLink = ref(false);
const addError = ref('');

// Edit link
const editingLinkId = ref<string | null>(null);
const editTitle = ref('');
const editUrl = ref('');
const savingEdit = ref(false);

// Drag reorder
let dragFrom = -1;

const publicUrl = computed(() =>
  bioPage.value ? `${window.location.origin}/bio/${bioPage.value.username}` : ''
);

const bioPagePublicUrl = computed(() => {
  if (!bioPage.value?.username) return '';
  return `${window.location.origin}/bio/${bioPage.value.username}`;
});

const bioStats = computed(() => ({
  total: bioPage.value?.links.length ?? 0,
  active: bioPage.value?.links.filter(l => l.is_active).length ?? 0,
  inactive: bioPage.value?.links.filter(l => !l.is_active).length ?? 0,
}));

const bioCopied = ref(false);
async function copyBioUrl() {
  if (!bioPagePublicUrl.value) return;
  await navigator.clipboard.writeText(bioPagePublicUrl.value);
  bioCopied.value = true;
  setTimeout(() => (bioCopied.value = false), 2000);
}

onMounted(async () => {
  try {
    const res = await bioApi.get();
    if (res.data) {
      bioPage.value = res.data;
      links.value = [...(res.data.links ?? [])];
      form.value = {
        username: res.data.username,
        title: res.data.title,
        description: res.data.description,
        avatar_url: res.data.avatar_url,
        theme: res.data.theme,
        is_published: res.data.is_published,
      };
    }
  } finally {
    loading.value = false;
  }
});

async function saveSettings() {
  saving.value = true;
  saveError.value = '';
  try {
    const res = await bioApi.update(form.value);
    if (res.data) {
      bioPage.value = res.data;
      form.value.username = res.data.username;
    }
  } catch (err: any) {
    saveError.value = err?.response?.data?.description ?? 'Failed to save settings.';
  } finally {
    saving.value = false;
  }
}

async function addLink() {
  if (!newTitle.value.trim() || !newUrl.value.trim()) return;
  addingLink.value = true;
  addError.value = '';
  try {
    const res = await bioApi.createLink({ title: newTitle.value.trim(), url: newUrl.value.trim() });
    if (res.data) {
      links.value.push(res.data);
      newTitle.value = '';
      newUrl.value = '';
    }
  } catch (err: any) {
    addError.value = err?.response?.data?.description ?? 'Failed to add link.';
  } finally {
    addingLink.value = false;
  }
}

async function toggleActive(link: BioLinkItem) {
  try {
    const res = await bioApi.updateLink(link.id, {
      title: link.title,
      url: link.url,
      is_active: !link.is_active,
    });
    if (res.data) {
      const idx = links.value.findIndex(l => l.id === link.id);
      if (idx !== -1) links.value[idx] = res.data;
    }
  } catch {
    // silent — checkbox stays in place
  }
}

async function deleteLink(link: BioLinkItem) {
  if (!confirm(`Delete "${link.title}"?`)) return;
  await bioApi.deleteLink(link.id);
  links.value = links.value.filter(l => l.id !== link.id);
}

function startEdit(link: BioLinkItem) {
  editingLinkId.value = link.id;
  editTitle.value = link.title;
  editUrl.value = link.url;
}

function cancelEdit() {
  editingLinkId.value = null;
}

async function saveEdit(link: BioLinkItem) {
  savingEdit.value = true;
  try {
    const res = await bioApi.updateLink(link.id, {
      title: editTitle.value.trim() || link.title,
      url: editUrl.value.trim() || link.url,
      is_active: link.is_active,
    });
    if (res.data) {
      const idx = links.value.findIndex(l => l.id === link.id);
      if (idx !== -1) links.value[idx] = res.data;
    }
    editingLinkId.value = null;
  } catch {
    // silent
  } finally {
    savingEdit.value = false;
  }
}

// ── Drag & drop reorder ───────────────────────────────────────────────────────
function onDragStart(idx: number) {
  dragFrom = idx;
}

function onDragOver(idx: number) {
  if (dragFrom === idx) return;
  const items = [...links.value];
  const [moved] = items.splice(dragFrom, 1);
  if (!moved) return;
  items.splice(idx, 0, moved);
  links.value = items;
  dragFrom = idx;
}

function onDrop() {
  // Persist order
  bioApi.reorderLinks({ ids: links.value.map(l => l.id) }).catch(() => {});
}

function onDragEnd() {
  dragFrom = -1;
}
</script>

<style scoped lang="scss">
/* ── Root ─────────────────────────────────────────────────────────────────── */
.bio-page {
  max-width: 900px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page header ──────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;

  &__actions {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 4px;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  margin: 0;
}

/* ── Stats ────────────────────────────────────────────────────────────────── */
.stats-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.stat-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  padding: 14px 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 110px;
}

.stat-label {
  font-size: 0.72rem;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.stat-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  line-height: 1.2;

  &--success { color: #16a34a; }
  &--muted { color: var(--md-sys-color-on-surface-variant); }
}

/* ── Loading ──────────────────────────────────────────────────────────────── */
.loading-state {
  display: flex;
  justify-content: center;
  padding: 48px;
}

/* ── AN Card ──────────────────────────────────────────────────────────────── */
.an-card {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);
}

.an-card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 18px; }

  &--primary {
    background: var(--md-sys-color-primary-container);
    .material-symbols-outlined { color: var(--md-sys-color-on-primary-container); }
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  flex: 1;
}

.an-card-hint {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.an-card-body {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ── CSS Spinner ──────────────────────────────────────────────────────────── */
.css-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2.5px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm {
    width: 16px;
    height: 16px;
    border-width: 2px;
  }

  &--white {
    border-color: rgba(255,255,255,0.35);
    border-top-color: #fff;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Form fields ──────────────────────────────────────────────────────────── */
.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;

  &__label {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
  }
}

.form-input {
  height: 40px;
  padding: 0 12px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9375rem;
  font-family: inherit;
  transition: border-color 0.15s, box-shadow 0.15s;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }
}

.form-textarea {
  padding: 8px 12px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9375rem;
  font-family: inherit;
  resize: vertical;
  transition: border-color 0.15s, box-shadow 0.15s;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }
}

.form-hint {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Username row ─────────────────────────────────────────────────────────── */
.bio-username-row {
  display: flex;
  align-items: stretch;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  overflow: hidden;
  transition: border-color 0.15s, box-shadow 0.15s;

  &:focus-within {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.12);
  }
}

.bio-prefix {
  padding: 0 12px;
  display: flex;
  align-items: center;
  background: var(--md-sys-color-surface-container-low);
  border-right: 1.5px solid var(--md-sys-color-outline-variant);
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
  white-space: nowrap;
}

.bio-username-input {
  flex: 1;
  border: none;
  border-radius: 0;
  height: 40px;

  &:focus {
    box-shadow: none;
    outline: none;
  }
}

/* ── Theme toggle ─────────────────────────────────────────────────────────── */
.theme-toggle {
  display: flex;
  gap: 8px;
}

/* ── Published toggle ─────────────────────────────────────────────────────── */
.published-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-radius: 8px;
  background: var(--md-sys-color-surface-container-low);
  gap: 12px;

  &--active {
    background: rgba(22, 163, 74, 0.06);
    border: 1px solid rgba(22, 163, 74, 0.25);
  }

  &__label {
    font-weight: 600;
    font-size: 0.875rem;
    color: var(--md-sys-color-on-surface);
  }

  &__hint {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.8rem;
    margin-top: 2px;
  }
}

/* ── Error box ────────────────────────────────────────────────────────────── */
.error-box {
  padding: 10px 14px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  font-size: 0.875rem;
}

/* ── Two column layout ────────────────────────────────────────────────────── */
.two-column-layout {
  display: grid;
  grid-template-columns: 1fr 1.4fr;
  gap: 20px;
  align-items: start;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.settings-card {
  /* styling from .an-card */
}

.links-column {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ── Add link row ─────────────────────────────────────────────────────────── */
.add-link-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: flex-end;
}

/* ── Empty state ──────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  text-align: center;
}

.empty-icon {
  width: 64px;
  height: 64px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;

  .material-symbols-outlined {
    font-size: 1.75rem;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.6;
  }
}

.empty-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 4px;
}

.empty-desc {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

/* ── Badges ───────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--success {
    background: rgba(22, 163, 74, 0.12);
    color: #16a34a;
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Link items ───────────────────────────────────────────────────────────── */
.link-item {
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  cursor: default;
  user-select: none;
  transition: background 0.12s;

  &:last-child { border-bottom: none; }
  &:hover { background: var(--md-sys-color-surface-container-low); }
}

.link-item__view {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  width: 100%;
}

.link-item__edit {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  width: 100%;
  flex-wrap: wrap;
  box-sizing: border-box;
}

.link-drag {
  cursor: grab;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 1.2rem;
  flex-shrink: 0;
}

.link-pos {
  width: 1.4rem;
  height: 1.4rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  font-size: 0.68rem;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
}

.link-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.link-info {
  flex: 1;
  min-width: 0;

  &__title {
    font-weight: 500;
    font-size: 0.875rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__url {
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.75rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.link-edit-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}
</style>
