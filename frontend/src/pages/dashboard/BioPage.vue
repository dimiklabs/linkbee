<template>
  <div class="page-section" style="max-width: 900px;">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Bio Page</h1>
        <p class="page-subtitle">Share a single page with all your important links.</p>
      </div>
      <div v-if="bioPage" style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
        <span :class="['m3-badge', bioPage.is_published ? 'm3-badge--success' : 'm3-badge--neutral']" style="padding:6px 14px;font-size:0.78rem;">
          {{ bioPage.is_published ? 'Published' : 'Draft' }}
        </span>
        <a
          v-if="bioPage.is_published"
          :href="publicUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          <button class="btn-outlined">
            <span class="material-symbols-outlined">open_in_new</span>
            View page
          </button>
        </a>
        <button class="btn-filled" :disabled="saving" @click="saveSettings">
          <md-circular-progress v-if="saving" indeterminate style="margin-right:6px;" />
          Save settings
        </button>
      </div>
    </div>

    <!-- Copy Public URL -->
    <div v-if="bioPage && bioPage.is_published && bioPage.username" style="margin-bottom:16px;">
      <button class="btn-text" @click="copyBioUrl">
        <span class="material-symbols-outlined">{{ bioCopied ? 'check' : 'content_copy' }}</span>
        {{ bioCopied ? 'Copied!' : 'Copy public URL' }}
      </button>
    </div>

    <!-- Stats Row -->
    <div v-if="bioPage" class="stats-row">
      <div class="m3-card m3-card--elevated stat-card">
        <div class="stat-label">Total Links</div>
        <div class="stat-value">{{ bioStats.total }}</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <div class="stat-label">Active Links</div>
        <div class="stat-value" style="color:#16a34a;">{{ bioStats.active }}</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <div class="stat-label">Inactive Links</div>
        <div class="stat-value" style="color:var(--md-sys-color-on-surface-variant);">{{ bioStats.inactive }}</div>
      </div>
      <div class="m3-card m3-card--elevated stat-card">
        <div class="stat-label">Status</div>
        <span :class="['m3-badge', bioPage.is_published ? 'm3-badge--success' : 'm3-badge--neutral']" style="margin-top:4px;">
          {{ bioPage.is_published ? 'Published' : 'Draft' }}
        </span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;justify-content:center;padding:48px;">
      <md-circular-progress indeterminate />
    </div>

    <div v-else-if="bioPage" class="two-column-layout">

      <!-- Left: Settings Form -->
      <div class="m3-card m3-card--elevated settings-card">
        <div class="card-section-header">
          <div class="card-section-header__left">
            <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">settings</span>
            <span class="md-label-large">Page Settings</span>
          </div>
        </div>
        <div class="card-body">

          <!-- Username -->
          <div style="margin-bottom:16px;">
            <div style="display:flex;align-items:center;gap:0;">
              <span style="padding:0 12px;height:56px;display:flex;align-items:center;background:var(--md-sys-color-surface-container-low);border:1px solid var(--md-sys-color-outline);border-right:none;border-radius:4px 0 0 4px;font-size:0.875rem;color:var(--md-sys-color-on-surface-variant);">/bio/</span>
              <md-outlined-text-field
                :value="form.username"
                @input="form.username=($event.target as HTMLInputElement).value"
                label="Username *"
                placeholder="your-username"
                maxlength="30"
                style="flex:1;"
              />
            </div>
            <div style="font-size:0.75rem;color:var(--md-sys-color-on-surface-variant);margin-top:4px;">Letters, numbers, underscores only. Min 3 characters.</div>
          </div>

          <!-- Title -->
          <div style="margin-bottom:16px;">
            <md-outlined-text-field
              :value="form.title"
              @input="form.title=($event.target as HTMLInputElement).value"
              label="Page title"
              placeholder="My Links"
              maxlength="100"
              style="width:100%;"
            />
          </div>

          <!-- Description -->
          <div style="margin-bottom:16px;">
            <md-outlined-text-field
              :value="form.description"
              @input="form.description=($event.target as HTMLInputElement).value"
              label="Bio / description"
              placeholder="A short bio or description..."
              type="textarea"
              rows="3"
              style="width:100%;"
            />
          </div>

          <!-- Avatar URL -->
          <div style="margin-bottom:16px;">
            <md-outlined-text-field
              :value="form.avatar_url"
              @input="form.avatar_url=($event.target as HTMLInputElement).value"
              label="Avatar URL"
              type="url"
              placeholder="https://..."
              style="width:100%;"
            />
          </div>

          <!-- Theme -->
          <div style="margin-bottom:16px;">
            <div style="font-size:0.875rem;font-weight:500;color:var(--md-sys-color-on-surface);margin-bottom:8px;">Theme</div>
            <div style="display:flex;gap:8px;">
              <button :class="['btn-outlined', form.theme === 'light' ? 'btn-active' : '']"
                @click="form.theme = 'light'"
                style="flex:1;"
              >
                <span class="material-symbols-outlined">light_mode</span>
                Light
              </button>
              <button :class="['btn-outlined', form.theme === 'dark' ? 'btn-active' : '']"
                @click="form.theme = 'dark'"
                style="flex:1;"
              >
                <span class="material-symbols-outlined">dark_mode</span>
                Dark
              </button>
            </div>
          </div>

          <!-- Published toggle -->
          <div :class="['published-toggle', form.is_published ? 'published-toggle--active' : '']">
            <div>
              <div style="font-weight:600;font-size:0.875rem;">Published</div>
              <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.8rem;margin-top:2px;">Make your bio page publicly accessible</div>
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

          <div v-if="saveError" style="margin-top:12px;padding:10px 14px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:8px;font-size:0.875rem;">
            {{ saveError }}
          </div>
        </div>
      </div>

      <!-- Right: Links Management -->
      <div>

        <!-- Add Link Form -->
        <div class="m3-card m3-card--elevated" style="margin-bottom:16px;">
          <div class="card-section-header">
            <div class="card-section-header__left">
              <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">add_link</span>
              <span class="md-label-large">Add Link</span>
            </div>
          </div>
          <div class="card-body">
            <div class="add-link-row">
              <md-outlined-text-field
                :value="newTitle"
                @input="newTitle=($event.target as HTMLInputElement).value"
                label="Link title"
                maxlength="100"
                style="flex:1;min-width:160px;"
                @keydown.enter="addLink"
              />
              <md-outlined-text-field
                :value="newUrl"
                @input="newUrl=($event.target as HTMLInputElement).value"
                label="URL"
                type="url"
                placeholder="https://..."
                style="flex:1;min-width:160px;"
                @keydown.enter="addLink"
              />
              <button class="btn-filled"
                :disabled="addingLink || !newTitle.trim() || !newUrl.trim()"
                @click="addLink"
                style="flex-shrink:0;"
              >
                <md-circular-progress v-if="addingLink" indeterminate style="margin-right:4px;" />
                <span v-else class="material-symbols-outlined">add</span>
                Add
              </button>
            </div>
            <div v-if="addError" style="margin-top:8px;padding:8px 12px;background:var(--md-sys-color-error-container);color:var(--md-sys-color-on-error-container);border-radius:8px;font-size:0.8rem;">
              {{ addError }}
            </div>
          </div>
        </div>

        <!-- Links List -->
        <div class="m3-card m3-card--elevated">
          <div class="card-section-header">
            <div class="card-section-header__left">
              <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">list</span>
              <span class="md-label-large">Links ({{ links.length }})</span>
            </div>
            <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Drag to reorder</span>
          </div>

          <div v-if="links.length === 0" class="m3-empty-state">
            <div class="m3-empty-state__icon">
              <span class="material-symbols-outlined">link_off</span>
            </div>
            <div class="md-title-small" style="margin-bottom:4px;">No links yet</div>
            <p class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin:0;">Add your first link using the form above.</p>
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
              <div v-if="editingLinkId !== link.id" style="display:flex;align-items:center;gap:8px;width:100%;">
                <span class="material-symbols-outlined" style="cursor:grab;color:var(--md-sys-color-on-surface-variant);font-size:1.2rem;flex-shrink:0;">drag_indicator</span>
                <span
                  style="width:1.4rem;height:1.4rem;border-radius:50%;display:flex;align-items:center;justify-content:center;background:var(--md-sys-color-surface-container-low);border:1px solid var(--md-sys-color-outline-variant);font-size:0.68rem;color:var(--md-sys-color-on-surface-variant);flex-shrink:0;"
                  :title="`Position ${idx + 1}`"
                >{{ idx + 1 }}</span>
                <span
                  style="width:8px;height:8px;border-radius:50%;display:inline-block;flex-shrink:0;"
                  :style="{ backgroundColor: link.is_active ? '#16a34a' : 'var(--md-sys-color-outline-variant)' }"
                  :title="link.is_active ? 'Active' : 'Inactive'"
                ></span>
                <div style="flex:1;min-width:0;">
                  <div style="font-weight:500;font-size:0.875rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ link.title }}</div>
                  <div style="color:var(--md-sys-color-on-surface-variant);font-size:0.75rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ link.url }}</div>
                </div>
                <!-- Active toggle -->
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
              <div v-else style="display:flex;align-items:center;gap:8px;width:100%;flex-wrap:wrap;">
                <md-outlined-text-field
                  :value="editTitle"
                  @input="editTitle=($event.target as HTMLInputElement).value"
                  label="Title"
                  maxlength="100"
                  style="flex:1;min-width:160px;"
                />
                <md-outlined-text-field
                  :value="editUrl"
                  @input="editUrl=($event.target as HTMLInputElement).value"
                  label="URL"
                  type="url"
                  style="flex:1;min-width:160px;"
                />
                <div style="display:flex;gap:4px;flex-shrink:0;">
                  <button class="btn-icon" :disabled="savingEdit" @click="saveEdit(link)">
                    <md-circular-progress v-if="savingEdit" indeterminate />
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
/* page-section (global) handles padding; max-width set via style attribute on root */

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
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

/* Stats */
.stats-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.stat-card {
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
}

/* Cards */
.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;
  margin-bottom: 0;

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
  }

  &--outlined {
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* Badges */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
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

/* Two column layout */
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
  /* card styling applied from m3-card classes */
}

.card-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &__left {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.card-body {
  padding: 16px 20px;
}

/* Published toggle */
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
}

/* Add link row */
.add-link-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: flex-end;
}

/* Empty state */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  text-align: center;

  &__icon {
    width: 64px;
    height: 64px;
    border-radius: 50%;
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
}

/* Link items */
.link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  cursor: default;
  user-select: none;
  transition: background 0.12s;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: var(--md-sys-color-surface-container-low);
  }
}
</style>
