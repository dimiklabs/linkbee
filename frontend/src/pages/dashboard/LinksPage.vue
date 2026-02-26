<template>
  <div class="container-fluid py-4">
    <div class="row g-3 align-items-start">

      <!-- ── Folder Sidebar ──────────────────────────────────────── -->
      <div class="col-md-3 col-lg-2 d-none d-md-block">
        <div class="card border-0 shadow-sm">
          <div class="card-header bg-white border-bottom py-3 px-3">
            <span class="fw-semibold small">Folders</span>
          </div>
          <div class="list-group list-group-flush">
            <!-- All Links -->
            <button
              class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2 px-3 border-0"
              :class="{ 'folder-active': selectedFolderID === '' && !starredOnly }"
              style="font-size: 0.85rem;"
              @click="selectFolder('')"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" class="flex-shrink-0" viewBox="0 0 16 16">
                <path d="M0 0h1v15h15v1H0zm14.817 3.113a.5.5 0 0 1 .07.704l-4.5 5.5a.5.5 0 0 1-.74.037L7.06 6.767l-3.656 5.027a.5.5 0 0 1-.808-.588l4-5.5a.5.5 0 0 1 .758-.06l2.609 2.61 4.15-5.073a.5.5 0 0 1 .704-.07"/>
              </svg>
              <span class="flex-grow-1 text-truncate">All Links</span>
            </button>

            <!-- Starred -->
            <button
              class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2 px-3 border-0"
              :class="{ 'folder-active': starredOnly }"
              style="font-size: 0.85rem;"
              @click="selectStarred"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" class="flex-shrink-0" viewBox="0 0 16 16">
                <path d="M3.612 15.443c-.386.198-.824-.149-.746-.592l.83-4.73L.173 6.765c-.329-.314-.158-.888.283-.95l4.898-.696L7.538.792c.197-.39.73-.39.927 0l2.184 4.327 4.898.696c.441.062.612.636.282.95l-3.522 3.356.83 4.73c.078.443-.36.79-.746.592L8 13.187l-4.389 2.256z"/>
              </svg>
              <span class="flex-grow-1 text-truncate">Starred</span>
            </button>

            <!-- Unhealthy -->
            <button
              class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2 px-3 border-0"
              :class="{ 'folder-active': healthFilter === 'unhealthy' }"
              style="font-size: 0.85rem;"
              @click="selectHealthFilter('unhealthy')"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" class="flex-shrink-0 text-danger" viewBox="0 0 16 16">
                <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/>
                <path d="M7.002 11a1 1 0 1 1 2 0 1 1 0 0 1-2 0zM7.1 4.995a.905.905 0 1 1 1.8 0l-.35 3.507a.552.552 0 0 1-1.1 0z"/>
              </svg>
              <span class="flex-grow-1 text-truncate">Unhealthy</span>
            </button>

            <!-- Folder rows -->
            <template v-for="folder in folders" :key="folder.id">
              <!-- Rename mode -->
              <div v-if="renamingFolderID === folder.id"
                class="d-flex align-items-center gap-1 px-2 py-1">
                <input
                  v-model="renameValue"
                  class="form-control form-control-sm flex-grow-1"
                  style="font-size: 0.8rem;"
                  @keydown.enter="submitRename(folder)"
                  @keydown.esc="cancelRename"
                  ref="renameInputRef"
                  maxlength="100"
                />
                <button class="btn btn-sm btn-primary px-1 py-0 lh-1" style="font-size: 0.8rem;" @click="submitRename(folder)">✓</button>
              </div>

              <!-- Normal mode -->
              <div v-else class="folder-row d-flex align-items-center gap-1 px-2 py-1"
                :class="{ 'folder-row-active': selectedFolderID === folder.id }">
                <button
                  class="d-flex align-items-center gap-2 border-0 bg-transparent text-start flex-grow-1 py-1 rounded"
                  style="font-size: 0.85rem; min-width: 0;"
                  :class="{ 'fw-semibold': selectedFolderID === folder.id }"
                  @click="selectFolder(folder.id)"
                >
                  <span class="flex-shrink-0 rounded-circle"
                    :style="{ width: '8px', height: '8px', backgroundColor: folder.color, display: 'inline-block' }">
                  </span>
                  <span class="text-truncate" :title="folder.name">{{ folder.name }}</span>
                </button>
                <div class="folder-actions d-flex gap-0 flex-shrink-0">
                  <button class="btn btn-sm border-0 p-1 lh-1 text-muted folder-action-btn" title="Rename" @click.stop="startRename(folder)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325"/>
                    </svg>
                  </button>
                  <button class="btn btn-sm border-0 p-1 lh-1 text-danger folder-action-btn" title="Delete folder" @click.stop="deleteFolder(folder)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                      <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                    </svg>
                  </button>
                </div>
              </div>
            </template>
          </div>

          <!-- New folder -->
          <div class="card-footer bg-white border-top p-2">
            <div v-if="showNewFolderInput" class="d-flex gap-1">
              <input
                v-model="newFolderName"
                class="form-control form-control-sm flex-grow-1"
                placeholder="Folder name"
                maxlength="100"
                ref="newFolderInputRef"
                @keydown.enter="createFolder"
                @keydown.esc="showNewFolderInput = false; newFolderName = ''"
              />
              <button class="btn btn-sm btn-primary px-2" @click="createFolder" :disabled="!newFolderName.trim()">✓</button>
            </div>
            <button
              v-else
              class="btn btn-sm btn-outline-secondary w-100 d-flex align-items-center justify-content-center gap-1"
              @click="openNewFolderInput"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
                <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
              </svg>
              <span style="font-size: 0.8rem;">New Folder</span>
            </button>
          </div>
        </div>
      </div>

      <!-- ── Main content ──────────────────────────────────────────── -->
      <div class="col-md-9 col-lg-10">

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

        <!-- Import button -->
        <button class="btn btn-outline-secondary d-flex align-items-center gap-2" @click="openImportModal">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
            <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708z"/>
          </svg>
          Import CSV
        </button>

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
              <th class="py-3 fw-semibold text-muted small text-uppercase" style="letter-spacing: 0.05em;">Health</th>
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

              <!-- Health -->
              <td class="py-3">
                <span
                  class="badge rounded-pill px-2 py-1"
                  :class="healthBadgeClass(link.health_status)"
                  style="font-size: 0.72rem;"
                  :title="link.health_checked_at ? `Last checked: ${formatDate(link.health_checked_at)}` : 'Not yet checked'"
                >
                  {{ healthLabel(link.health_status) }}
                </span>
              </td>

              <!-- Created -->
              <td class="py-3 text-muted small">
                {{ formatDate(link.created_at) }}
              </td>

              <!-- Actions -->
              <td class="pe-4 py-3 text-end">
                <div class="d-flex align-items-center justify-content-end gap-1">
                  <!-- Star -->
                  <button
                    class="btn btn-sm border-0 p-1"
                    :class="link.is_starred ? 'text-warning' : 'btn-outline-secondary'"
                    :title="link.is_starred ? 'Unstar' : 'Star'"
                    @click="toggleStar(link)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" :fill="link.is_starred ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="1" viewBox="0 0 16 16">
                      <path d="M3.612 15.443c-.386.198-.824-.149-.746-.592l.83-4.73L.173 6.765c-.329-.314-.158-.888.283-.95l4.898-.696L7.538.792c.197-.39.73-.39.927 0l2.184 4.327 4.898.696c.441.062.612.636.282.95l-3.522 3.356.83 4.73c.078.443-.36.79-.746.592L8 13.187l-4.389 2.256z"/>
                    </svg>
                  </button>

                  <!-- Health check -->
                  <button
                    class="btn btn-sm btn-outline-secondary border-0 p-1"
                    :title="checkingHealthId === link.id ? 'Checking...' : 'Check link health'"
                    :disabled="checkingHealthId === link.id"
                    @click="runHealthCheck(link)"
                  >
                    <span v-if="checkingHealthId === link.id" class="spinner-border spinner-border-sm" style="width: 14px; height: 14px;" role="status"></span>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M8 3.5a.5.5 0 0 0-1 0V9a.5.5 0 0 0 .252.434l3.5 2a.5.5 0 0 0 .496-.868L8 8.71z"/>
                      <path d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16m7-8A7 7 0 1 1 1 8a7 7 0 0 1 14 0"/>
                    </svg>
                  </button>

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

                  <!-- A/B Split Test -->
                  <button
                    class="btn btn-sm border-0 p-1"
                    :class="link.is_split_test ? 'text-primary' : 'btn-outline-secondary'"
                    title="A/B Split Test"
                    @click="openSplitModal(link)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path fill-rule="evenodd" d="M11.5 2a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3M9.05 3.5a2.5 2.5 0 0 1 4.9 0H16v1h-2.05a2.5 2.5 0 0 1-4.9 0H0v-1zM0 9.5h2.05a2.5 2.5 0 0 1 4.9 0H16v1H6.95a2.5 2.5 0 0 1-4.9 0zm4.5-1a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3"/>
                    </svg>
                  </button>

                  <!-- Geo Routing -->
                  <button
                    class="btn btn-sm border-0 p-1"
                    :class="link.is_geo_routing ? 'text-success' : 'btn-outline-secondary'"
                    title="Geo Routing"
                    @click="openGeoModal(link)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0M2.04 4.326c.325 1.329 2.532 2.54 3.717 3.19.48.263.793.434.743.484q-.121.12-.242.234c-.416.396-.787.749-.758 1.266.035.634.618.824 1.214 1.017.577.188 1.168.38 1.286.983.082.417-.075.988-.22 1.52-.215.782-.406 1.48.22 1.48 1.5-.5 3.798-3.186 4-5 .138-1.243-2-2-3.5-2.5-.478-.16-.755.081-.99.284-.172.15-.322.279-.51.216-.445-.148-2.614-1.208-2.614-2 0-.56.23-1.272.55-2zm6.959 0c.922.065 1.459.22 1.459.22 -.12.1-.32.3-.39.5-.069.2.049.3.2.4.15.1.29.2.31.4.019.2-.11.3-.23.4-.11.1-.22.2-.19.4.03.2.18.3.33.4.16.1.32.2.36.4.039.2-.08.3-.19.4-.12.1-.25.2-.25.4 0 .2.17.3.33.4.17.1.35.2.37.4.02.2-.1.4-.24.5-.14.1-.3.2-.32.4l-.015.046c-.023.068-.04.12-.06.16.9-.2 1.6-.5 1.9-1.1.3-.6.1-1.3-.2-2z"/>
                    </svg>
                  </button>

                  <!-- Retargeting Pixels -->
                  <button
                    class="btn btn-sm border-0 p-1"
                    :class="link.is_pixel_tracking ? 'text-danger' : 'btn-outline-secondary'"
                    title="Retargeting Pixels"
                    @click="openPixelModal(link)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0m.5 11.5a.5.5 0 0 1-1 0v-4a.5.5 0 0 1 1 0zm-.5-6a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5"/>
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
      :folders="folders"
      @saved="onLinkSaved"
    />

    <!-- QR Code Modal -->
    <QRCodeModal
      v-if="qrLink"
      ref="qrModalRef"
      :link-id="qrLink.id"
      :slug="qrLink.slug"
    />

    <!-- Import CSV Modal -->
    <ImportLinksModal
      ref="importModalRef"
      @imported="onImportDone"
    />

    <!-- A/B Split Test Modal -->
    <SplitTestModal
      v-if="splitLink"
      ref="splitModalRef"
      :link-id="splitLink.id"
      :slug="splitLink.slug"
      :is-split-test-initial="splitLink.is_split_test"
      @updated="onSplitTestUpdated"
    />

    <!-- Geo Routing Modal -->
    <GeoRulesModal
      v-if="geoLink"
      ref="geoModalRef"
      :link-id="geoLink.id"
      :slug="geoLink.slug"
      :is-geo-routing-initial="geoLink.is_geo_routing"
      @updated="onGeoRoutingUpdated"
    />

    <!-- Retargeting Pixels Modal -->
    <PixelsModal
      v-if="pixelLink"
      ref="pixelModalRef"
      modal-id="pixels-modal"
      :link="pixelLink"
      @pixel-tracking-updated="onPixelTrackingUpdated"
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

      </div><!-- end main content col -->
    </div><!-- end row -->
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue';
import { RouterLink } from 'vue-router';
import { Toast } from 'bootstrap';
import { useLinksStore } from '@/stores/links';
import type { LinkResponse } from '@/types/links';
import type { FolderResponse } from '@/types/folders';
import foldersApi from '@/api/folders';
import CreateLinkModal from '@/components/CreateLinkModal.vue';
import QRCodeModal from '@/components/QRCodeModal.vue';
import ImportLinksModal from '@/components/ImportLinksModal.vue';
import SplitTestModal from '@/components/SplitTestModal.vue';
import GeoRulesModal from '@/components/GeoRulesModal.vue';
import PixelsModal from '@/components/PixelsModal.vue';
import type { ImportLinksResponse } from '@/types/links';

const linksStore = useLinksStore();

const searchQuery = ref('');
const editingLink = ref<LinkResponse | null>(null);
const qrLink = ref<LinkResponse | null>(null);
const copiedId = ref<string | null>(null);
const deletingId = ref<string | null>(null);

const createModalRef = ref<InstanceType<typeof CreateLinkModal> | null>(null);
const qrModalRef = ref<InstanceType<typeof QRCodeModal> | null>(null);
const importModalRef = ref<InstanceType<typeof ImportLinksModal> | null>(null);
const splitModalRef = ref<InstanceType<typeof SplitTestModal> | null>(null);
const splitLink = ref<LinkResponse | null>(null);
const geoModalRef = ref<InstanceType<typeof GeoRulesModal> | null>(null);
const geoLink = ref<LinkResponse | null>(null);
const pixelModalRef = ref<InstanceType<typeof PixelsModal> | null>(null);
const pixelLink = ref<LinkResponse | null>(null);
const toastEl = ref<HTMLElement | null>(null);
let toastInstance: Toast | null = null;

// ── Folder / filter state ────────────────────────────────────────────────────
const folders = ref<FolderResponse[]>([]);
const selectedFolderID = ref('');
const starredOnly = ref(false);
const healthFilter = ref('');
const checkingHealthId = ref<string | null>(null);

// New folder creation
const showNewFolderInput = ref(false);
const newFolderName = ref('');
const newFolderInputRef = ref<HTMLInputElement | null>(null);

// Inline rename
const renamingFolderID = ref('');
const renameValue = ref('');
const renameInputRef = ref<HTMLInputElement | null>(null);

async function loadFolders() {
  try {
    const res = await foldersApi.list();
    if (res.data) folders.value = res.data;
  } catch {
    // non-critical — sidebar just stays empty
  }
}

function selectFolder(id: string) {
  selectedFolderID.value = id;
  starredOnly.value = false;
  healthFilter.value = '';
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, id);
}

function selectStarred() {
  starredOnly.value = true;
  selectedFolderID.value = '';
  healthFilter.value = '';
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', true);
}

function selectHealthFilter(status: string) {
  healthFilter.value = status;
  selectedFolderID.value = '';
  starredOnly.value = false;
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', undefined, status);
}

function openNewFolderInput() {
  showNewFolderInput.value = true;
  nextTick(() => newFolderInputRef.value?.focus());
}

async function createFolder() {
  const name = newFolderName.value.trim();
  if (!name) return;
  try {
    const res = await foldersApi.create({ name, color: '#635bff' });
    if (res.data) {
      folders.value.push(res.data);
      folders.value.sort((a, b) => a.name.localeCompare(b.name));
    }
  } catch {
    // silent
  } finally {
    newFolderName.value = '';
    showNewFolderInput.value = false;
  }
}

function startRename(folder: FolderResponse) {
  renamingFolderID.value = folder.id;
  renameValue.value = folder.name;
  nextTick(() => renameInputRef.value?.focus());
}

function cancelRename() {
  renamingFolderID.value = '';
  renameValue.value = '';
}

async function submitRename(folder: FolderResponse) {
  const name = renameValue.value.trim();
  if (!name) { cancelRename(); return; }
  try {
    const res = await foldersApi.update(folder.id, { name });
    if (res.data) {
      const idx = folders.value.findIndex(f => f.id === folder.id);
      if (idx !== -1) folders.value[idx] = res.data;
      folders.value.sort((a, b) => a.name.localeCompare(b.name));
    }
  } catch {
    // silent
  } finally {
    cancelRename();
  }
}

async function deleteFolder(folder: FolderResponse) {
  if (!window.confirm(`Delete folder "${folder.name}"? Links inside will be moved to All Links.`)) return;
  try {
    await foldersApi.delete(folder.id);
    folders.value = folders.value.filter(f => f.id !== folder.id);
    if (selectedFolderID.value === folder.id) {
      selectedFolderID.value = '';
      linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', starredOnly.value || undefined, healthFilter.value || undefined);
    }
  } catch {
    // silent
  }
}

// Debounced search
let searchTimer: ReturnType<typeof setTimeout> | null = null;
watch(searchQuery, (val) => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    linksStore.fetchLinks(1, linksStore.limit, val, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined);
  }, 400);
});

onMounted(() => {
  linksStore.fetchLinks(1, 20, '');
  loadFolders();
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

function openImportModal() {
  importModalRef.value?.show();
}

function openSplitModal(link: LinkResponse) {
  splitLink.value = link;
  setTimeout(() => splitModalRef.value?.show(), 50);
}

function onSplitTestUpdated(enabled: boolean) {
  if (splitLink.value) {
    const idx = linksStore.links.findIndex(l => l.id === splitLink.value!.id);
    if (idx !== -1) linksStore.links[idx] = { ...linksStore.links[idx], is_split_test: enabled };
    splitLink.value = { ...splitLink.value, is_split_test: enabled };
  }
}

function openGeoModal(link: LinkResponse) {
  geoLink.value = link;
  setTimeout(() => geoModalRef.value?.show(), 50);
}

function onGeoRoutingUpdated(enabled: boolean) {
  if (geoLink.value) {
    const idx = linksStore.links.findIndex(l => l.id === geoLink.value!.id);
    if (idx !== -1) linksStore.links[idx] = { ...linksStore.links[idx], is_geo_routing: enabled };
    geoLink.value = { ...geoLink.value, is_geo_routing: enabled };
  }
}

function openPixelModal(link: LinkResponse) {
  pixelLink.value = link;
  setTimeout(() => pixelModalRef.value?.show(), 50);
}

function onPixelTrackingUpdated(enabled: boolean) {
  if (pixelLink.value) {
    const idx = linksStore.links.findIndex(l => l.id === pixelLink.value!.id);
    if (idx !== -1) linksStore.links[idx] = { ...linksStore.links[idx], is_pixel_tracking: enabled };
    pixelLink.value = { ...pixelLink.value, is_pixel_tracking: enabled };
  }
}

async function onImportDone(_result: ImportLinksResponse) {
  await linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined);
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
  await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined);
  editingLink.value = null;
}

async function toggleStar(link: LinkResponse) {
  await linksStore.toggleStar(link.id);
  if (starredOnly.value && !linksStore.links.find(l => l.id === link.id)?.is_starred) {
    await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, '', true, healthFilter.value || undefined);
  }
}

async function runHealthCheck(link: LinkResponse) {
  checkingHealthId.value = link.id;
  try {
    await linksStore.checkHealth(link.id);
    // If filtering by unhealthy and link just became healthy, refresh list
    if (healthFilter.value === 'unhealthy') {
      const updated = linksStore.links.find(l => l.id === link.id);
      if (updated && updated.health_status !== 'unhealthy') {
        await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, 'unhealthy');
      }
    }
  } finally {
    checkingHealthId.value = null;
  }
}

function healthBadgeClass(status: string): string {
  switch (status) {
    case 'healthy':   return 'text-bg-success';
    case 'unhealthy': return 'text-bg-danger';
    case 'timeout':   return 'text-bg-warning';
    case 'error':     return 'text-bg-secondary';
    default:          return 'bg-light text-secondary border';
  }
}

function healthLabel(status: string): string {
  switch (status) {
    case 'healthy':   return 'Healthy';
    case 'unhealthy': return 'Unhealthy';
    case 'timeout':   return 'Timeout';
    case 'error':     return 'Error';
    default:          return 'Unknown';
  }
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
    await linksStore.fetchLinks(targetPage, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined);
  } finally {
    deletingId.value = null;
  }
}

function goToPage(page: number) {
  if (page < 1 || page > linksStore.totalPages) return;
  linksStore.fetchLinks(page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined);
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

/* Folder sidebar */
.folder-active {
  background-color: rgba(99, 91, 255, 0.08) !important;
  color: #635bff !important;
  font-weight: 600;
}

.folder-row {
  border-radius: 6px;
  transition: background 0.15s;
}

.folder-row:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.folder-row-active {
  background-color: rgba(99, 91, 255, 0.08);
}

.folder-action-btn {
  opacity: 0;
  transition: opacity 0.15s;
}

.folder-row:hover .folder-action-btn {
  opacity: 1;
}
</style>
