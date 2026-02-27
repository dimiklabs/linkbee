<template>
  <div class="container-fluid py-4" style="max-width: 900px;">

    <!-- Header -->
    <div class="d-flex align-items-center justify-content-between flex-wrap gap-3 mb-4">
      <div>
        <h4 class="mb-1 fw-bold">Link-in-Bio</h4>
        <p class="text-muted small mb-0">Share a single page with all your important links.</p>
      </div>
      <div v-if="bioPage" class="d-flex align-items-center gap-2 flex-wrap">
        <!-- Published badge -->
        <span
          class="badge rounded-pill px-3 py-2"
          :class="bioPage.is_published ? 'text-bg-success' : 'bg-light text-secondary border'"
          style="font-size: 0.78rem;"
        >
          {{ bioPage.is_published ? 'Published' : 'Draft' }}
        </span>
        <!-- Public URL -->
        <a
          v-if="bioPage.is_published"
          :href="publicUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="btn btn-sm btn-outline-secondary d-flex align-items-center gap-1"
        >
          <i class="bi bi-box-arrow-up-right"></i>
          View page
        </a>
        <!-- Save button -->
        <button
          class="btn btn-primary btn-sm d-flex align-items-center gap-1"
          :disabled="saving"
          @click="saveSettings"
        >
          <span v-if="saving" class="spinner-border spinner-border-sm"></span>
          Save settings
        </button>
      </div>
    </div>

    <!-- Copy Public URL button in header area -->
    <div v-if="bioPage && bioPage.is_published && bioPage.username" class="mb-3">
      <button
        class="btn btn-sm btn-outline-primary d-flex align-items-center gap-1"
        @click="copyBioUrl"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
          <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1z"/>
          <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0z"/>
        </svg>
        {{ bioCopied ? 'Copied!' : 'Copy public URL' }}
      </button>
    </div>

    <!-- Page Overview Stats Bar -->
    <div v-if="bioPage" class="d-flex gap-3 flex-wrap mb-4">
      <!-- Total Links -->
      <div class="stat-chip">
        <span class="stat-chip__label">Total Links</span>
        <span class="stat-chip__value">{{ bioStats.total }}</span>
      </div>
      <!-- Active Links -->
      <div class="stat-chip">
        <span class="stat-chip__label">Active Links</span>
        <span class="stat-chip__value text-success">{{ bioStats.active }}</span>
      </div>
      <!-- Inactive Links -->
      <div class="stat-chip">
        <span class="stat-chip__label">Inactive Links</span>
        <span class="stat-chip__value text-secondary">{{ bioStats.inactive }}</span>
      </div>
      <!-- Status -->
      <div class="stat-chip">
        <span class="stat-chip__label">Status</span>
        <span
          class="badge rounded-pill px-2 py-1"
          :class="bioPage.is_published ? 'text-bg-success' : 'bg-secondary bg-opacity-25 text-secondary'"
          style="font-size: 0.72rem;"
        >
          {{ bioPage.is_published ? 'Published' : 'Draft' }}
        </span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status"></div>
    </div>

    <div v-else-if="bioPage" class="row g-4">

      <!-- Left: Settings -->
      <div class="col-lg-5">
        <div class="card border-0 shadow-sm mb-4">
          <div class="card-header bg-white border-bottom py-3 px-4">
            <span class="fw-semibold small">Page Settings</span>
          </div>
          <div class="card-body px-4 py-3">

            <!-- Username -->
            <div class="mb-3">
              <label class="form-label small fw-medium">Username <span class="text-danger">*</span></label>
              <div class="input-group">
                <span class="input-group-text bg-light text-muted small">/bio/</span>
                <input
                  v-model="form.username"
                  type="text"
                  class="form-control"
                  placeholder="your-username"
                  maxlength="30"
                />
              </div>
              <div class="form-text">Letters, numbers, underscores only. Min 3 characters.</div>
            </div>

            <!-- Title -->
            <div class="mb-3">
              <label class="form-label small fw-medium">Page title</label>
              <input
                v-model="form.title"
                type="text"
                class="form-control"
                placeholder="My Links"
                maxlength="100"
              />
            </div>

            <!-- Description -->
            <div class="mb-3">
              <label class="form-label small fw-medium">Bio / description</label>
              <textarea
                v-model="form.description"
                class="form-control"
                rows="3"
                placeholder="A short bio or description..."
              ></textarea>
            </div>

            <!-- Avatar URL -->
            <div class="mb-3">
              <label class="form-label small fw-medium">Avatar URL</label>
              <input
                v-model="form.avatar_url"
                type="url"
                class="form-control"
                placeholder="https://..."
              />
            </div>

            <!-- Theme -->
            <div class="mb-3">
              <label class="form-label small fw-medium">Theme</label>
              <div class="d-flex gap-2">
                <button
                  class="btn btn-sm flex-fill"
                  :class="form.theme === 'light' ? 'btn-primary' : 'btn-outline-secondary'"
                  @click="form.theme = 'light'"
                >
                  ☀️ Light
                </button>
                <button
                  class="btn btn-sm flex-fill"
                  :class="form.theme === 'dark' ? 'btn-primary' : 'btn-outline-secondary'"
                  @click="form.theme = 'dark'"
                >
                  🌙 Dark
                </button>
              </div>
            </div>

            <!-- Published toggle -->
            <div class="d-flex align-items-center justify-content-between p-3 rounded"
                 :class="form.is_published ? 'bg-success bg-opacity-10 border border-success border-opacity-25' : 'bg-light'">
              <div>
                <div class="fw-semibold small">Published</div>
                <div class="text-muted" style="font-size: 0.8rem;">Make your bio page publicly accessible</div>
              </div>
              <div class="form-check form-switch mb-0">
                <input
                  v-model="form.is_published"
                  class="form-check-input"
                  type="checkbox"
                  role="switch"
                  style="width: 2.5rem; height: 1.25rem; cursor: pointer;"
                />
              </div>
            </div>

            <div v-if="saveError" class="alert alert-danger py-2 small mt-3 mb-0">{{ saveError }}</div>
          </div>
        </div>
      </div>

      <!-- Right: Links + Preview -->
      <div class="col-lg-7">

        <!-- Add link form -->
        <div class="card border-0 shadow-sm mb-4">
          <div class="card-header bg-white border-bottom py-3 px-4">
            <span class="fw-semibold small">Add Link</span>
          </div>
          <div class="card-body px-4 py-3">
            <div class="row g-2">
              <div class="col-12 col-sm-5">
                <input
                  v-model="newTitle"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="Link title"
                  maxlength="100"
                  @keydown.enter="addLink"
                />
              </div>
              <div class="col-12 col-sm-5">
                <input
                  v-model="newUrl"
                  type="url"
                  class="form-control form-control-sm"
                  placeholder="https://..."
                  @keydown.enter="addLink"
                />
              </div>
              <div class="col-12 col-sm-2">
                <button
                  class="btn btn-primary btn-sm w-100"
                  :disabled="addingLink || !newTitle.trim() || !newUrl.trim()"
                  @click="addLink"
                >
                  <span v-if="addingLink" class="spinner-border spinner-border-sm"></span>
                  <span v-else>Add</span>
                </button>
              </div>
            </div>
            <div v-if="addError" class="alert alert-danger py-2 small mt-2 mb-0">{{ addError }}</div>
          </div>
        </div>

        <!-- Links list -->
        <div class="card border-0 shadow-sm">
          <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
            <span class="fw-semibold small">Links ({{ links.length }})</span>
            <span class="text-muted small">Drag to reorder</span>
          </div>

          <div v-if="links.length === 0" class="card-body text-center text-muted py-4">
            <p class="mb-0 small">No links yet. Add your first link above.</p>
          </div>

          <div v-else class="list-group list-group-flush">
            <div
              v-for="(link, idx) in links"
              :key="link.id"
              class="list-group-item py-2 px-4"
              :style="{ opacity: link.is_active ? '1' : '0.6' }"
              draggable="true"
              @dragstart="onDragStart(idx)"
              @dragover.prevent="onDragOver(idx)"
              @drop="onDrop"
              @dragend="onDragEnd"
            >
              <!-- View mode -->
              <div v-if="editingLinkId !== link.id" class="d-flex align-items-center gap-2">
                <span class="text-muted flex-shrink-0" style="cursor: grab; font-size: 1rem;">⠿</span>
                <!-- Rank indicator -->
                <span
                  class="badge rounded-circle d-flex align-items-center justify-content-center flex-shrink-0 text-secondary bg-light border"
                  style="width: 1.4rem; height: 1.4rem; font-size: 0.68rem; padding: 0;"
                  :title="`Position ${idx + 1}`"
                >{{ idx + 1 }}</span>
                <!-- Active/Inactive dot -->
                <span
                  class="flex-shrink-0"
                  style="width: 8px; height: 8px; border-radius: 50%; display: inline-block;"
                  :style="{ backgroundColor: link.is_active ? '#198754' : '#adb5bd' }"
                  :title="link.is_active ? 'Active' : 'Inactive'"
                ></span>
                <div class="flex-fill min-w-0">
                  <div class="fw-medium small text-truncate">{{ link.title }}</div>
                  <div class="text-muted text-truncate" style="font-size: 0.78rem;">{{ link.url }}</div>
                </div>
                <!-- Active toggle -->
                <div class="form-check form-switch mb-0 flex-shrink-0">
                  <input
                    class="form-check-input"
                    type="checkbox"
                    role="switch"
                    :checked="link.is_active"
                    style="cursor: pointer;"
                    @change="toggleActive(link)"
                  />
                </div>
                <!-- Edit -->
                <button class="btn btn-sm border-0 p-1 text-muted" title="Edit" @click="startEdit(link)">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                    <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325"/>
                  </svg>
                </button>
                <!-- Delete -->
                <button class="btn btn-sm border-0 p-1 text-danger" title="Delete" @click="deleteLink(link)">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                    <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                    <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                  </svg>
                </button>
              </div>

              <!-- Edit mode -->
              <div v-else class="d-flex align-items-center gap-2">
                <div class="flex-fill row g-1">
                  <div class="col-12 col-sm-5">
                    <input v-model="editTitle" type="text" class="form-control form-control-sm" placeholder="Title" maxlength="100" />
                  </div>
                  <div class="col-12 col-sm-5">
                    <input v-model="editUrl" type="url" class="form-control form-control-sm" placeholder="URL" />
                  </div>
                  <div class="col-12 col-sm-2 d-flex gap-1">
                    <button class="btn btn-sm btn-primary px-2" :disabled="savingEdit" @click="saveEdit(link)">
                      <span v-if="savingEdit" class="spinner-border spinner-border-sm"></span>
                      <span v-else>✓</span>
                    </button>
                    <button class="btn btn-sm btn-outline-secondary px-2" @click="cancelEdit">✕</button>
                  </div>
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

<style scoped>
.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}
.btn-primary:hover {
  background-color: #5249e0;
  border-color: #5249e0;
}
.list-group-item[draggable="true"] {
  cursor: default;
  user-select: none;
}

/* Page Overview stat chips */
.stat-chip {
  display: inline-flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  background: #fff;
  border: 1px solid #e9ecef;
  border-radius: 10px;
  padding: 8px 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  min-width: 90px;
}

.stat-chip__label {
  font-size: 0.7rem;
  color: #6c757d;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  white-space: nowrap;
}

.stat-chip__value {
  font-size: 1rem;
  font-weight: 700;
  color: #212529;
  line-height: 1.2;
}
</style>
