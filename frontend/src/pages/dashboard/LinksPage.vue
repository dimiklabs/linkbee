<template>
  <div class="container-fluid py-4">
    <!-- Page Header -->
    <div class="d-flex align-items-center justify-content-between flex-wrap gap-3 mb-4">
      <div>
        <h4 class="mb-1 fw-bold">My Links</h4>
        <p class="text-muted small mb-0">
          {{ linksStore.total }} link{{ linksStore.total !== 1 ? 's' : '' }} total
        </p>
      </div>
      <div class="d-flex align-items-center gap-2 flex-wrap">
        <!-- Search -->
        <div class="input-group" style="min-width: 240px; max-width: 320px;">
          <span class="input-group-text bg-white border-end-0">
            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" fill="currentColor" class="text-muted" viewBox="0 0 16 16">
              <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"/>
            </svg>
          </span>
          <input
            v-model="searchQuery"
            type="search"
            class="form-control border-start-0 ps-0"
            placeholder="Search links..."
            aria-label="Search links"
          />
        </div>

        <!-- Create button -->
        <button class="btn btn-primary d-flex align-items-center gap-2" @click="openCreateModal">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
          </svg>
          Create Link
        </button>
      </div>
    </div>

    <!-- Error alert -->
    <div v-if="linksStore.error" class="alert alert-danger d-flex align-items-center gap-2" role="alert">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5m.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
      </svg>
      {{ linksStore.error }}
    </div>

    <!-- Loading state -->
    <div v-if="linksStore.loading" class="text-center py-5">
      <div class="spinner-border text-primary" style="width: 2.5rem; height: 2.5rem;" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="text-muted mt-3 mb-0">Loading your links...</p>
    </div>

    <!-- Empty state -->
    <div
      v-else-if="!linksStore.loading && linksStore.links.length === 0"
      class="text-center py-5"
    >
      <div class="mb-3">
        <svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" fill="currentColor" class="text-muted opacity-50" viewBox="0 0 16 16">
          <path d="M4.715 6.542 3.343 7.914a3 3 0 1 0 4.243 4.243l1.828-1.829A3 3 0 0 0 8.586 5.5L8 6.086a1 1 0 0 0-.154.199 2 2 0 0 1 .861 3.337L6.88 11.45a2 2 0 1 1-2.83-2.83l.793-.792a4 4 0 0 1-.128-1.287z"/>
          <path d="M6.586 4.672A3 3 0 0 0 7.414 9.5l.775-.776a2 2 0 0 1-.896-3.346L9.12 3.55a2 2 0 1 1 2.83 2.83l-.793.792c.112.42.155.855.128 1.287l1.372-1.372a3 3 0 1 0-4.243-4.243z"/>
        </svg>
      </div>
      <h5 class="fw-semibold text-muted">
        {{ searchQuery ? 'No links match your search' : 'No links yet' }}
      </h5>
      <p class="text-muted small mb-4">
        {{ searchQuery ? 'Try a different search term.' : 'Get started by creating your first shortened link.' }}
      </p>
      <button v-if="!searchQuery" class="btn btn-primary" @click="openCreateModal">
        Create your first link
      </button>
      <button v-else class="btn btn-outline-secondary btn-sm" @click="searchQuery = ''">
        Clear search
      </button>
    </div>

    <!-- Links table -->
    <div v-else class="card shadow-sm border-0">
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th class="ps-4 py-3 fw-semibold text-muted small text-uppercase" style="letter-spacing: 0.05em;">Title / URL</th>
              <th class="py-3 fw-semibold text-muted small text-uppercase" style="letter-spacing: 0.05em;">Short URL</th>
              <th class="py-3 fw-semibold text-muted small text-uppercase" style="letter-spacing: 0.05em;">Clicks</th>
              <th class="py-3 fw-semibold text-muted small text-uppercase" style="letter-spacing: 0.05em;">Status</th>
              <th class="py-3 fw-semibold text-muted small text-uppercase" style="letter-spacing: 0.05em;">Created</th>
              <th class="pe-4 py-3 fw-semibold text-muted small text-uppercase text-end" style="letter-spacing: 0.05em;">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="link in linksStore.links" :key="link.id">
              <!-- Title / URL -->
              <td class="ps-4 py-3" style="max-width: 280px;">
                <div class="fw-medium text-truncate" :title="link.title || link.destination_url">
                  {{ link.title || '—' }}
                </div>
                <div class="text-muted small text-truncate" :title="link.destination_url">
                  {{ link.destination_url }}
                </div>
                <div v-if="link.tags && link.tags.length > 0" class="mt-1 d-flex flex-wrap gap-1">
                  <span
                    v-for="tag in link.tags.slice(0, 3)"
                    :key="tag"
                    class="badge rounded-pill bg-light text-secondary border"
                    style="font-size: 0.68rem; font-weight: 500;"
                  >
                    {{ tag }}
                  </span>
                  <span v-if="link.tags.length > 3" class="badge rounded-pill bg-light text-secondary border" style="font-size: 0.68rem;">
                    +{{ link.tags.length - 3 }}
                  </span>
                </div>
              </td>

              <!-- Short URL -->
              <td class="py-3">
                <div class="d-flex align-items-center gap-2">
                  <a :href="link.short_url" target="_blank" rel="noopener noreferrer" class="link-primary text-decoration-none fw-medium small">
                    {{ link.short_url }}
                  </a>
                  <button
                    class="btn btn-sm btn-outline-secondary border-0 p-1 copy-btn"
                    :title="copiedId === link.id ? 'Copied!' : 'Copy to clipboard'"
                    @click="copyShortUrl(link)"
                    style="line-height: 1;"
                  >
                    <template v-if="copiedId === link.id">
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="#198754" viewBox="0 0 16 16">
                        <path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0"/>
                      </svg>
                    </template>
                    <template v-else>
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                        <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1z"/>
                        <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0z"/>
                      </svg>
                    </template>
                  </button>
                </div>
              </td>

              <!-- Clicks -->
              <td class="py-3">
                <span class="fw-semibold">{{ link.click_count.toLocaleString() }}</span>
              </td>

              <!-- Status -->
              <td class="py-3">
                <span
                  class="badge rounded-pill px-3 py-2"
                  :class="link.is_active ? 'text-bg-success' : 'text-bg-danger'"
                  style="font-size: 0.75rem;"
                >
                  {{ link.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>

              <!-- Created -->
              <td class="py-3 text-muted small">
                {{ formatDate(link.created_at) }}
              </td>

              <!-- Actions -->
              <td class="pe-4 py-3 text-end">
                <div class="d-flex align-items-center justify-content-end gap-1">
                  <!-- Analytics -->
                  <RouterLink
                    :to="`/dashboard/links/${link.id}`"
                    class="btn btn-sm btn-outline-secondary border-0 p-1"
                    title="View Analytics"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M4 11H2v3h2zm5-4H7v7h2zm5-5v12h-2V2zm-2-1a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1zM6 7a1 1 0 0 0-1 1v7a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1V8a1 1 0 0 0-1-1zm-5 4a1 1 0 0 0-1 1v3a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1z"/>
                    </svg>
                  </RouterLink>

                  <!-- QR Code -->
                  <button
                    class="btn btn-sm btn-outline-secondary border-0 p-1"
                    title="View QR Code"
                    @click="openQRModal(link)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M2 2h2v2H2z"/>
                      <path d="M6 0v6H0V0zM5 1H1v4h4zM4 12H2v2h2z"/>
                      <path d="M6 10v6H0v-6zm-5 1v4h4v-4zm11-9h2v2h-2z"/>
                      <path d="M10 0v6h6V0zm5 1v4h-4V1zM8 1V0h1v2H8v2H7V1zm0 5H7V4h1zM6 8V7h1V6h1v2h1V7h5v1h-4v1H7V8zm0 0v1H2V8H1v1H0V7h3v1zm10 1h-1V7h1zm-1 0h-1v2h2v-1h-1zm-4 0h2v1h-1v1h-1zm2 3v-1h-1v1h-1v1H9v1h3v-2zm0 0h3v1h-2v1h-1zm-4-1v1h1v-2H7v1z"/>
                    </svg>
                  </button>

                  <!-- Edit -->
                  <button
                    class="btn btn-sm btn-outline-secondary border-0 p-1"
                    title="Edit Link"
                    @click="openEditModal(link)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325"/>
                    </svg>
                  </button>

                  <!-- Delete -->
                  <button
                    class="btn btn-sm btn-outline-danger border-0 p-1"
                    title="Delete Link"
                    :disabled="deletingId === link.id"
                    @click="confirmDelete(link)"
                  >
                    <span v-if="deletingId === link.id" class="spinner-border spinner-border-sm" style="width: 14px; height: 14px;" role="status"></span>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                      <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div
        v-if="linksStore.totalPages > 1"
        class="card-footer bg-white border-top d-flex align-items-center justify-content-between flex-wrap gap-2 py-3 px-4"
      >
        <p class="text-muted small mb-0">
          Page {{ linksStore.page }} of {{ linksStore.totalPages }}
          &mdash; {{ linksStore.total }} total
        </p>
        <nav aria-label="Links pagination">
          <ul class="pagination pagination-sm mb-0">
            <li class="page-item" :class="{ disabled: linksStore.page <= 1 }">
              <button class="page-link" @click="goToPage(linksStore.page - 1)" :disabled="linksStore.page <= 1">
                Previous
              </button>
            </li>

            <li
              v-for="(p, idx) in visiblePages"
              :key="idx"
              class="page-item"
              :class="{ active: p === linksStore.page, disabled: p === '...' }"
            >
              <button
                v-if="p !== '...'"
                class="page-link"
                @click="goToPage(Number(p))"
              >
                {{ p }}
              </button>
              <span v-else class="page-link border-0 bg-transparent text-muted">...</span>
            </li>

            <li class="page-item" :class="{ disabled: linksStore.page >= linksStore.totalPages }">
              <button
                class="page-link"
                @click="goToPage(linksStore.page + 1)"
                :disabled="linksStore.page >= linksStore.totalPages"
              >
                Next
              </button>
            </li>
          </ul>
        </nav>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <CreateLinkModal
      ref="createModalRef"
      :link="editingLink ?? undefined"
      @saved="onLinkSaved"
    />

    <!-- QR Code Modal -->
    <QRCodeModal
      v-if="qrLink"
      ref="qrModalRef"
      :link-id="qrLink.id"
      :slug="qrLink.slug"
    />

    <!-- Copy toast -->
    <div
      class="position-fixed bottom-0 end-0 p-3"
      style="z-index: 1100;"
      aria-live="polite"
      aria-atomic="true"
    >
      <div
        ref="toastEl"
        class="toast align-items-center text-bg-dark border-0"
        role="alert"
        aria-live="assertive"
        aria-atomic="true"
      >
        <div class="d-flex">
          <div class="toast-body d-flex align-items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="#4ade80" viewBox="0 0 16 16">
              <path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0"/>
            </svg>
            Copied to clipboard!
          </div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { RouterLink } from 'vue-router';
import { Toast } from 'bootstrap';
import { useLinksStore } from '@/stores/links';
import type { LinkResponse } from '@/types/links';
import CreateLinkModal from '@/components/CreateLinkModal.vue';
import QRCodeModal from '@/components/QRCodeModal.vue';

const linksStore = useLinksStore();

const searchQuery = ref('');
const editingLink = ref<LinkResponse | null>(null);
const qrLink = ref<LinkResponse | null>(null);
const copiedId = ref<string | null>(null);
const deletingId = ref<string | null>(null);

const createModalRef = ref<InstanceType<typeof CreateLinkModal> | null>(null);
const qrModalRef = ref<InstanceType<typeof QRCodeModal> | null>(null);
const toastEl = ref<HTMLElement | null>(null);
let toastInstance: Toast | null = null;

// Debounced search
let searchTimer: ReturnType<typeof setTimeout> | null = null;
watch(searchQuery, (val) => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    linksStore.fetchLinks(1, linksStore.limit, val);
  }, 400);
});

onMounted(() => {
  linksStore.fetchLinks(1, 20, '');
  if (toastEl.value) {
    toastInstance = new Toast(toastEl.value, { delay: 2000 });
  }
});

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

async function copyShortUrl(link: LinkResponse) {
  try {
    await navigator.clipboard.writeText(link.short_url);
    copiedId.value = link.id;
    toastInstance?.show();
    setTimeout(() => {
      if (copiedId.value === link.id) copiedId.value = null;
    }, 2000);
  } catch {
    // fallback silent fail
  }
}

function openCreateModal() {
  editingLink.value = null;
  createModalRef.value?.show();
}

function openEditModal(link: LinkResponse) {
  editingLink.value = link;
  createModalRef.value?.show();
}

function openQRModal(link: LinkResponse) {
  qrLink.value = link;
  // Allow DOM to update with v-if before showing
  setTimeout(() => {
    qrModalRef.value?.show();
  }, 50);
}

async function onLinkSaved(_link: LinkResponse) {
  await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value);
  editingLink.value = null;
}

async function confirmDelete(link: LinkResponse) {
  const confirmed = window.confirm(
    `Are you sure you want to delete "${link.title || link.slug}"? This action cannot be undone.`
  );
  if (!confirmed) return;

  deletingId.value = link.id;
  try {
    await linksStore.deleteLink(link.id);
    const newTotal = linksStore.total - 1;
    const maxPage = Math.ceil(newTotal / linksStore.limit) || 1;
    const targetPage = Math.min(linksStore.page, maxPage);
    await linksStore.fetchLinks(targetPage, linksStore.limit, searchQuery.value);
  } finally {
    deletingId.value = null;
  }
}

function goToPage(page: number) {
  if (page < 1 || page > linksStore.totalPages) return;
  linksStore.fetchLinks(page, linksStore.limit, searchQuery.value);
}

const visiblePages = computed<(number | string)[]>(() => {
  const total = linksStore.totalPages;
  const current = linksStore.page;
  const pages: (number | string)[] = [];

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i);
    return pages;
  }

  pages.push(1);
  if (current > 3) pages.push('...');
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let i = start; i <= end; i++) pages.push(i);
  if (current < total - 2) pages.push('...');
  pages.push(total);

  return pages;
});
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

.link-primary {
  color: #635bff;
}

.link-primary:hover {
  color: #5249e0;
}

.page-item.active .page-link {
  background-color: #635bff;
  border-color: #635bff;
}

.page-link:focus {
  box-shadow: 0 0 0 0.25rem rgba(99, 91, 255, 0.25);
}

.copy-btn:hover {
  background-color: rgba(0, 0, 0, 0.05) !important;
}
</style>
