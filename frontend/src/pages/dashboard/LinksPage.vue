<template>
  <div class="links-page">
    <div class="links-layout">

      <!-- ── Folder Sidebar ──────────────────────────────────────── -->
      <aside class="folder-sidebar">
        <div class="m3-card m3-card--elevated sidebar-card">
          <div class="sidebar-header">
            <span class="sidebar-title">Folders</span>
          </div>
          <div class="sidebar-list">
            <!-- All Links -->
            <button
              class="sidebar-item"
              :class="{ 'sidebar-item--active': selectedFolderID === '' && !starredOnly && healthFilter === '' && !expiringSoonFilter }"
              @click="selectFolder('')"
            >
              <span class="material-symbols-outlined" style="font-size:16px">bar_chart</span>
              <span class="sidebar-item-label">All Links</span>
            </button>

            <!-- Starred -->
            <button
              class="sidebar-item"
              :class="{ 'sidebar-item--active': starredOnly }"
              @click="selectStarred"
            >
              <span class="material-symbols-outlined" style="font-size:16px">star</span>
              <span class="sidebar-item-label">Starred</span>
            </button>

            <!-- Unhealthy -->
            <button
              class="sidebar-item"
              :class="{ 'sidebar-item--active': healthFilter === 'unhealthy' }"
              @click="selectHealthFilter('unhealthy')"
            >
              <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-error)">error</span>
              <span class="sidebar-item-label">Unhealthy</span>
            </button>

            <!-- Expiring Soon -->
            <button
              class="sidebar-item"
              :class="{ 'sidebar-item--active': expiringSoonFilter }"
              @click="selectExpiringSoon"
            >
              <span class="material-symbols-outlined" style="font-size:16px;color:#F4A100">schedule</span>
              <span class="sidebar-item-label">Expiring Soon</span>
            </button>

            <!-- Folder rows -->
            <template v-for="folder in folders" :key="folder.id">
              <!-- Rename mode -->
              <div v-if="renamingFolderID === folder.id" class="sidebar-rename-row">
                <md-outlined-text-field
                  :value="renameValue"
                  @input="renameValue = ($event.target as HTMLInputElement).value"
                  label="Folder name"
                  style="flex:1;--md-outlined-text-field-container-shape:8px;--md-outlined-text-field-input-text-size:0.82rem"
                  @keydown.enter="submitRename(folder)"
                  @keydown.esc="cancelRename"
                  ref="renameInputRef"
                />
                <md-icon-button @click="submitRename(folder)">
                  <span class="material-symbols-outlined">check</span>
                </md-icon-button>
              </div>

              <!-- Normal mode -->
              <div v-else class="folder-row" :class="{ 'folder-row--active': selectedFolderID === folder.id }">
                <button
                  class="folder-row-btn"
                  :class="{ 'folder-row-btn--active': selectedFolderID === folder.id }"
                  @click="selectFolder(folder.id)"
                >
                  <span class="folder-dot" :style="{ backgroundColor: folder.color }"></span>
                  <span class="folder-row-name">{{ folder.name }}</span>
                  <span class="folder-row-count">{{ folder.click_count?.toLocaleString() }}</span>
                </button>
                <div class="folder-actions">
                  <md-icon-button @click.stop="startRename(folder)" title="Rename">
                    <span class="material-symbols-outlined" style="font-size:14px">edit</span>
                  </md-icon-button>
                  <md-icon-button @click.stop="deleteFolder(folder)" title="Delete folder" style="--md-icon-button-icon-color:var(--md-sys-color-error)">
                    <span class="material-symbols-outlined" style="font-size:14px">delete</span>
                  </md-icon-button>
                </div>
              </div>
            </template>
          </div>

          <!-- New folder -->
          <div class="sidebar-footer">
            <div v-if="showNewFolderInput" class="new-folder-row">
              <md-outlined-text-field
                :value="newFolderName"
                @input="newFolderName = ($event.target as HTMLInputElement).value"
                label="Folder name"
                style="flex:1;--md-outlined-text-field-container-shape:8px;--md-outlined-text-field-input-text-size:0.82rem"
                ref="newFolderInputRef"
                @keydown.enter="createFolder"
                @keydown.esc="showNewFolderInput = false; newFolderName = ''"
              />
              <md-icon-button @click="createFolder" :disabled="!newFolderName.trim()">
                <span class="material-symbols-outlined">check</span>
              </md-icon-button>
            </div>
            <md-text-button v-else @click="openNewFolderInput" style="width:100%">
              <span class="material-symbols-outlined" slot="icon">create_new_folder</span>
              New Folder
            </md-text-button>
          </div>
        </div>

        <!-- Tag filter card -->
        <div v-if="allTags.length > 0" class="m3-card m3-card--elevated sidebar-tags-card">
          <div class="sidebar-tags-header">
            <span class="sidebar-title">Tags</span>
            <md-text-button v-if="selectedTags.length > 0" @click="clearTagFilter" style="--md-text-button-label-text-size:0.75rem">Clear</md-text-button>
          </div>
          <div class="sidebar-tags-body">
            <md-chip-set>
              <md-filter-chip
                v-for="tag in allTags"
                :key="tag"
                :label="tag"
                :selected="selectedTags.includes(tag)"
                @click="toggleTag(tag)"
              />
            </md-chip-set>
          </div>
        </div>
      </aside>

      <!-- ── Main content ──────────────────────────────────────────── -->
      <div class="links-main">

        <!-- Page Header -->
        <div class="page-header">
          <div class="page-header-left">
            <h4 class="page-title">My Links</h4>
            <p class="page-subtitle">
              {{ linksStore.total }} link{{ linksStore.total !== 1 ? 's' : '' }} total
            </p>
          </div>
          <div class="page-header-actions">
            <md-outlined-text-field
              :value="searchQuery"
              @input="searchQuery = ($event.target as HTMLInputElement).value"
              label="Search links..."
              style="min-width:240px;max-width:320px"
            >
              <span class="material-symbols-outlined" slot="leading-icon">search</span>
            </md-outlined-text-field>

            <md-outlined-button @click="exportCSV" :disabled="exporting">
              <md-circular-progress v-if="exporting" indeterminate style="--md-circular-progress-size:18px" slot="icon" />
              <span v-else class="material-symbols-outlined" slot="icon">download</span>
              Export CSV
            </md-outlined-button>

            <md-outlined-button @click="openImportModal">
              <span class="material-symbols-outlined" slot="icon">upload_file</span>
              Import CSV
            </md-outlined-button>

            <md-filled-button @click="openCreateModal">
              <span class="material-symbols-outlined" slot="icon">add</span>
              Create Link
            </md-filled-button>
          </div>
        </div>

        <!-- Usage warning banner -->
        <div
          v-if="usageWarning"
          class="m3-card m3-card--outlined usage-banner"
          :style="{ borderColor: usageWarning.level === 'danger' ? 'var(--md-sys-color-error)' : '#F4A100' }"
        >
          <span class="material-symbols-outlined" :style="{ color: usageWarning.level === 'danger' ? 'var(--md-sys-color-error)' : '#F4A100' }">
            {{ usageWarning.level === 'danger' ? 'block' : 'warning' }}
          </span>
          <span class="usage-banner-msg">{{ usageWarning.msg }}</span>
          <router-link to="/dashboard/billing">
            <md-filled-tonal-button style="--md-filled-tonal-button-container-color:var(--md-sys-color-error-container)">Upgrade</md-filled-tonal-button>
          </router-link>
        </div>

        <!-- Error alert -->
        <div v-if="linksStore.error" class="m3-card m3-card--outlined error-banner">
          <span class="material-symbols-outlined" style="color:var(--md-sys-color-error)">error</span>
          <span>{{ linksStore.error }}</span>
        </div>

        <!-- Loading state -->
        <div v-if="linksStore.loading" class="loading-state">
          <md-circular-progress indeterminate style="--md-circular-progress-size:48px" />
          <p class="loading-text">Loading your links...</p>
        </div>

        <!-- Empty state -->
        <div
          v-else-if="!linksStore.loading && linksStore.links.length === 0"
          class="empty-state"
        >
          <span class="material-symbols-outlined empty-state-icon">link_off</span>
          <h5 class="empty-state-title">
            {{ searchQuery ? 'No links match your search' : selectedTags.length ? 'No links with these tags' : 'No links yet' }}
          </h5>
          <p class="empty-state-subtitle">
            {{ searchQuery ? 'Try a different search term.' : selectedTags.length ? 'Try selecting different tags or clear the tag filter.' : 'Get started by creating your first shortened link.' }}
          </p>
          <md-filled-button v-if="!searchQuery && !selectedTags.length" @click="openCreateModal">
            Create your first link
          </md-filled-button>
          <md-outlined-button v-else-if="selectedTags.length" @click="clearTagFilter">
            Clear tag filter
          </md-outlined-button>
          <md-outlined-button v-else @click="searchQuery = ''">
            Clear search
          </md-outlined-button>
        </div>

        <!-- Links table -->
        <div v-else class="m3-card m3-card--elevated links-table-card">

          <!-- Bulk action bar -->
          <div v-if="selectedIDs.size > 0" class="bulk-bar">
            <span class="bulk-bar-count">{{ selectedIDs.size }} selected</span>
            <md-divider vertical style="height:24px;margin:0 4px" />

            <md-outlined-button :disabled="bulkLoading" @click="bulkActivate">Activate</md-outlined-button>
            <md-outlined-button :disabled="bulkLoading" @click="bulkDeactivate">Deactivate</md-outlined-button>

            <!-- Move to folder -->
            <div class="bulk-popover-wrap">
              <md-outlined-button :disabled="bulkLoading" @click="showBulkFolderBar = !showBulkFolderBar; showBulkTagsBar = false">
                Move to Folder
              </md-outlined-button>
              <div v-if="showBulkFolderBar" class="bulk-popover">
                <button class="bulk-popover-item" @click="bulkMoveFolder(null)">— Remove from folder</button>
                <button
                  v-for="f in folders"
                  :key="f.id"
                  class="bulk-popover-item"
                  @click="bulkMoveFolder(f.id)"
                >{{ f.name }}</button>
              </div>
            </div>

            <!-- Tags -->
            <div class="bulk-popover-wrap">
              <md-outlined-button :disabled="bulkLoading" @click="bulkTagsMode = 'add_tags'; showBulkTagsBar = !showBulkTagsBar; showBulkFolderBar = false">+ Tags</md-outlined-button>
              <md-outlined-button :disabled="bulkLoading" @click="bulkTagsMode = 'remove_tags'; showBulkTagsBar = !showBulkTagsBar; showBulkFolderBar = false">− Tags</md-outlined-button>
              <div v-if="showBulkTagsBar" class="bulk-popover bulk-tags-popover">
                <md-outlined-text-field
                  :value="bulkTagsInput"
                  @input="bulkTagsInput = ($event.target as HTMLInputElement).value"
                  label="tag1, tag2, ..."
                  style="flex:1"
                  @keydown.enter="bulkApplyTags"
                />
                <md-filled-button @click="bulkApplyTags">
                  {{ bulkTagsMode === 'add_tags' ? 'Add' : 'Remove' }}
                </md-filled-button>
              </div>
            </div>

            <md-outlined-button :disabled="qrDownloading" @click="downloadBulkQR">
              <md-circular-progress v-if="qrDownloading" indeterminate style="--md-circular-progress-size:16px" slot="icon" />
              <span v-else class="material-symbols-outlined" slot="icon">qr_code_2</span>
              QR Codes
            </md-outlined-button>

            <div class="bulk-bar-end">
              <md-outlined-button :disabled="bulkLoading" @click="bulkDelete" style="--md-outlined-button-outline-color:var(--md-sys-color-error);--md-outlined-button-label-text-color:var(--md-sys-color-error)">
                <md-circular-progress v-if="bulkLoading" indeterminate style="--md-circular-progress-size:16px" slot="icon" />
                Delete
              </md-outlined-button>
              <md-text-button @click="clearSelection">Clear</md-text-button>
            </div>
          </div>

          <div class="table-scroll">
            <table class="m3-table">
              <thead>
                <tr>
                  <th style="width:44px">
                    <md-checkbox
                      :checked="isAllSelected"
                      :indeterminate="selectedIDs.size > 0 && !isAllSelected"
                      @change="toggleSelectAll"
                    />
                  </th>
                  <th>Title / URL</th>
                  <th>Short URL</th>
                  <th>Clicks</th>
                  <th>Status</th>
                  <th>Health</th>
                  <th>Created</th>
                  <th style="text-align:right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="link in linksStore.links"
                  :key="link.id"
                  :class="{ 'row--selected': selectedIDs.has(link.id) }"
                >
                  <!-- Select -->
                  <td>
                    <md-checkbox
                      :checked="selectedIDs.has(link.id)"
                      @change="toggleSelect(link.id)"
                    />
                  </td>

                  <!-- Title / URL -->
                  <td style="max-width:280px">
                    <div class="link-title" :title="link.title || link.destination_url">
                      {{ link.title || '—' }}
                    </div>
                    <div class="link-dest" :title="link.destination_url">
                      {{ link.destination_url }}
                    </div>
                    <div v-if="link.tags && link.tags.length > 0" class="link-tags">
                      <span v-for="tag in link.tags.slice(0, 3)" :key="tag" class="m3-badge m3-badge--neutral" style="font-size:0.68rem">{{ tag }}</span>
                      <span v-if="link.tags.length > 3" class="m3-badge m3-badge--neutral" style="font-size:0.68rem">+{{ link.tags.length - 3 }}</span>
                    </div>
                    <div class="link-meta">
                      <span v-if="link.has_password" class="link-meta-item">
                        <span class="material-symbols-outlined" style="font-size:13px">lock</span> Protected
                      </span>
                      <span
                        v-if="link.expires_at"
                        class="link-meta-item"
                        :style="{ color: isExpired(link) ? 'var(--md-sys-color-error)' : isExpiringSoon(link) ? '#F4A100' : 'inherit', fontWeight: (isExpired(link) || isExpiringSoon(link)) ? '600' : 'normal' }"
                        :title="isExpired(link) ? 'This link has expired' : isExpiringSoon(link) ? 'Expiring within 3 days' : 'Expiry date'"
                      >
                        <span class="material-symbols-outlined" style="font-size:13px">schedule</span>
                        {{ isExpired(link) ? 'Expired' : 'Expires' }} {{ formatDate(link.expires_at) }}
                      </span>
                      <span
                        v-if="link.max_clicks"
                        class="link-meta-item"
                        :style="{ color: isClickLimitReached(link) ? '#F4A100' : 'inherit', fontWeight: isClickLimitReached(link) ? '600' : 'normal' }"
                        :title="isClickLimitReached(link) ? 'Click limit reached' : 'Max clicks'"
                      >
                        <span class="material-symbols-outlined" style="font-size:13px">bolt</span>
                        {{ link.click_count }}/{{ link.max_clicks }} clicks
                      </span>
                    </div>
                  </td>

                  <!-- Short URL -->
                  <td>
                    <div class="short-url-cell">
                      <a :href="link.short_url" target="_blank" rel="noopener noreferrer" class="short-url-link">
                        {{ link.short_url }}
                      </a>
                      <span v-if="expiryBadgeClass(link)" :class="['m3-badge', expiryBadgeClass(link) === 'text-bg-danger' ? 'm3-badge--error' : 'm3-badge--warning']" style="font-size:0.65rem">
                        {{ expiryBadgeLabel(link) }}
                      </span>
                      <span v-if="clickLimitBadgeClass(link)" :class="['m3-badge', clickLimitBadgeClass(link) === 'text-bg-danger' ? 'm3-badge--error' : 'm3-badge--warning']" style="font-size:0.65rem">
                        {{ clickLimitBadgeLabel(link) }}
                      </span>
                      <md-icon-button
                        :title="copiedId === link.id ? 'Copied!' : 'Copy to clipboard'"
                        @click="copyShortUrl(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px"
                      >
                        <span class="material-symbols-outlined" :style="{ color: copiedId === link.id ? '#1AA563' : 'inherit' }">
                          {{ copiedId === link.id ? 'check_circle' : 'content_copy' }}
                        </span>
                      </md-icon-button>
                    </div>
                  </td>

                  <!-- Clicks -->
                  <td>
                    <span class="click-count">{{ link.click_count.toLocaleString() }}</span>
                  </td>

                  <!-- Status -->
                  <td>
                    <span :class="['m3-badge', statusM3BadgeClass(link)]">
                      {{ statusLabel(link) }}
                    </span>
                  </td>

                  <!-- Health -->
                  <td>
                    <span
                      :class="['m3-badge', healthM3BadgeClass(link.health_status)]"
                      :title="link.health_checked_at ? `Last checked: ${formatDate(link.health_checked_at)}` : 'Not yet checked'"
                    >
                      {{ healthLabel(link.health_status) }}
                    </span>
                  </td>

                  <!-- Created -->
                  <td class="date-cell">
                    {{ formatDate(link.created_at) }}
                  </td>

                  <!-- Actions -->
                  <td style="text-align:right">
                    <div class="actions-cell">
                      <!-- Star -->
                      <md-icon-button
                        :title="link.is_starred ? 'Unstar' : 'Star'"
                        @click="toggleStar(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px"
                      >
                        <span class="material-symbols-outlined" :style="{ color: link.is_starred ? '#F4A100' : 'inherit', fontVariationSettings: link.is_starred ? '\'FILL\' 1' : '\'FILL\' 0' }">star</span>
                      </md-icon-button>

                      <!-- Health check -->
                      <md-icon-button
                        :title="checkingHealthId === link.id ? 'Checking...' : 'Check link health'"
                        :disabled="checkingHealthId === link.id"
                        @click="runHealthCheck(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px"
                      >
                        <md-circular-progress v-if="checkingHealthId === link.id" indeterminate style="--md-circular-progress-size:16px" />
                        <span v-else class="material-symbols-outlined">health_and_safety</span>
                      </md-icon-button>

                      <!-- Analytics -->
                      <RouterLink :to="`/dashboard/links/${link.id}`">
                        <md-icon-button title="View Analytics" style="--md-icon-button-icon-size:18px;width:32px;height:32px">
                          <span class="material-symbols-outlined">bar_chart</span>
                        </md-icon-button>
                      </RouterLink>

                      <!-- Preview -->
                      <md-icon-button title="Link Preview" @click="openPreviewModal(link)" style="--md-icon-button-icon-size:18px;width:32px;height:32px">
                        <span class="material-symbols-outlined">visibility</span>
                      </md-icon-button>

                      <!-- QR Code -->
                      <md-icon-button title="View QR Code" @click="openQRModal(link)" style="--md-icon-button-icon-size:18px;width:32px;height:32px">
                        <span class="material-symbols-outlined">qr_code_2</span>
                      </md-icon-button>

                      <!-- A/B Split Test -->
                      <md-icon-button
                        title="A/B Split Test"
                        @click="openSplitModal(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px"
                        :style="link.is_split_test ? { '--md-icon-button-icon-color': 'var(--md-sys-color-primary)' } : {}"
                      >
                        <span class="material-symbols-outlined">join</span>
                      </md-icon-button>

                      <!-- Geo Routing -->
                      <md-icon-button
                        title="Geo Routing"
                        @click="openGeoModal(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px"
                        :style="link.is_geo_routing ? { '--md-icon-button-icon-color': '#1AA563' } : {}"
                      >
                        <span class="material-symbols-outlined">public</span>
                      </md-icon-button>

                      <!-- Retargeting Pixels -->
                      <md-icon-button
                        title="Retargeting Pixels"
                        @click="openPixelModal(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px"
                        :style="link.is_pixel_tracking ? { '--md-icon-button-icon-color': 'var(--md-sys-color-error)' } : {}"
                      >
                        <span class="material-symbols-outlined">track_changes</span>
                      </md-icon-button>

                      <!-- Edit -->
                      <md-icon-button title="Edit Link" @click="openEditModal(link)" style="--md-icon-button-icon-size:18px;width:32px;height:32px">
                        <span class="material-symbols-outlined">edit</span>
                      </md-icon-button>

                      <!-- Clone -->
                      <md-icon-button title="Clone Link" @click="openCloneModal(link)" style="--md-icon-button-icon-size:18px;width:32px;height:32px">
                        <span class="material-symbols-outlined">file_copy</span>
                      </md-icon-button>

                      <!-- Delete -->
                      <md-icon-button
                        title="Delete Link"
                        :disabled="deletingId === link.id"
                        @click="confirmDelete(link)"
                        style="--md-icon-button-icon-size:18px;width:32px;height:32px;--md-icon-button-icon-color:var(--md-sys-color-error)"
                      >
                        <md-circular-progress v-if="deletingId === link.id" indeterminate style="--md-circular-progress-size:16px" />
                        <span v-else class="material-symbols-outlined">delete</span>
                      </md-icon-button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagination -->
          <div v-if="linksStore.totalPages > 1" class="pagination-bar">
            <p class="pagination-info">
              Page {{ linksStore.page }} of {{ linksStore.totalPages }} &mdash; {{ linksStore.total }} total
            </p>
            <div class="pagination-controls">
              <md-outlined-button @click="goToPage(linksStore.page - 1)" :disabled="linksStore.page <= 1">
                Previous
              </md-outlined-button>

              <template v-for="(p, idx) in visiblePages" :key="idx">
                <md-filled-button v-if="p === linksStore.page" style="min-width:40px">{{ p }}</md-filled-button>
                <md-outlined-button v-else-if="p !== '...'" @click="goToPage(Number(p))" style="min-width:40px">{{ p }}</md-outlined-button>
                <span v-else class="pagination-ellipsis">...</span>
              </template>

              <md-outlined-button @click="goToPage(linksStore.page + 1)" :disabled="linksStore.page >= linksStore.totalPages">
                Next
              </md-outlined-button>
            </div>
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

        <!-- Link Preview Modal -->
        <LinkPreviewModal
          v-if="previewLink"
          ref="previewModalRef"
          modal-id="preview-modal"
          :link="previewLink"
        />

        <!-- Clone Link Modal -->
        <CloneLinkModal
          v-if="cloneLink"
          ref="cloneModalRef"
          modal-id="clone-link-modal"
          :link="cloneLink"
          @cloned="onLinkCloned"
        />

        <!-- Copy snackbar -->
        <Transition name="snack">
          <div v-if="copiedId" class="m3-snackbar">
            <span class="material-symbols-outlined" style="font-size:20px">check_circle</span>
            <span style="flex:1">Copied to clipboard!</span>
            <md-text-button @click="copiedId = null" style="--md-text-button-label-text-color:#CFBCFF">Dismiss</md-text-button>
          </div>
        </Transition>

      </div><!-- end main content -->
    </div><!-- end layout -->
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue';
import { RouterLink } from 'vue-router';
import { useLinksStore } from '@/stores/links';
import type { LinkResponse } from '@/types/links';
import type { FolderResponse } from '@/types/folders';
import linksApi from '@/api/links';
import foldersApi from '@/api/folders';
import billingApi from '@/api/billing';
import type { UsageCounts, PlanInfo } from '@/types/billing';
import CreateLinkModal from '@/components/CreateLinkModal.vue';
import QRCodeModal from '@/components/QRCodeModal.vue';
import ImportLinksModal from '@/components/ImportLinksModal.vue';
import SplitTestModal from '@/components/SplitTestModal.vue';
import GeoRulesModal from '@/components/GeoRulesModal.vue';
import PixelsModal from '@/components/PixelsModal.vue';
import LinkPreviewModal from '@/components/LinkPreviewModal.vue';
import CloneLinkModal from '@/components/CloneLinkModal.vue';
import type { ImportLinksResponse } from '@/types/links';

const linksStore = useLinksStore();

// ── Billing / Usage ───────────────────────────────────────────────────────────
const usage = ref<UsageCounts | null>(null);
const plan = ref<PlanInfo | null>(null);

const usageWarning = computed(() => {
  if (!usage.value || !plan.value) return null;
  const used = usage.value.links;
  const max = plan.value.max_links;
  if (max === -1) return null;
  const pct = used / max;
  if (pct >= 1) return { level: 'danger', msg: `You've reached your limit of ${max} links. Upgrade to add more.` };
  if (pct >= 0.8) return { level: 'warning', msg: `You've used ${used} of ${max} links (${Math.round(pct * 100)}%). Consider upgrading.` };
  return null;
});

// ── Expiry / Click-limit badge helpers ────────────────────────────────────────
function expiryBadgeClass(link: LinkResponse): string | null {
  if (!link.expires_at) return null;
  const days = Math.ceil((new Date(link.expires_at).getTime() - Date.now()) / 86400000);
  if (days <= 0) return 'text-bg-danger';
  if (days <= 3) return 'text-bg-danger';
  if (days <= 7) return 'text-bg-warning';
  return null;
}

function expiryBadgeLabel(link: LinkResponse): string | null {
  if (!link.expires_at) return null;
  const days = Math.ceil((new Date(link.expires_at).getTime() - Date.now()) / 86400000);
  if (days <= 0) return 'Expired';
  return `Exp. ${days}d`;
}

function clickLimitBadgeClass(link: LinkResponse): string | null {
  if (!link.max_clicks) return null;
  const pct = link.click_count / link.max_clicks;
  if (pct >= 1) return 'text-bg-danger';
  if (pct >= 0.9) return 'text-bg-warning';
  return null;
}

function clickLimitBadgeLabel(link: LinkResponse): string | null {
  if (!link.max_clicks) return null;
  const pct = link.click_count / link.max_clicks;
  if (pct >= 1) return 'At limit';
  if (pct >= 0.9) return `${Math.round(pct * 100)}% used`;
  return null;
}

const exporting = ref(false);
const searchQuery = ref('');
const allTags = ref<string[]>([]);
const selectedTags = ref<string[]>([]);
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
const previewModalRef = ref<InstanceType<typeof LinkPreviewModal> | null>(null);
const previewLink = ref<LinkResponse | null>(null);
const cloneModalRef = ref<InstanceType<typeof CloneLinkModal> | null>(null);
const cloneLink = ref<LinkResponse | null>(null);
const toastEl = ref<HTMLElement | null>(null);
let toastInstance: any = null;

// ── Bulk selection ────────────────────────────────────────────────────────────
const selectedIDs = ref<Set<string>>(new Set());
const bulkLoading = ref(false);
const qrDownloading = ref(false);
const bulkTagsInput = ref('');
const bulkTagsMode = ref<'add_tags' | 'remove_tags'>('add_tags');
const showBulkTagsBar = ref(false);
const showBulkFolderBar = ref(false);

const isAllSelected = computed(
  () => linksStore.links.length > 0 && selectedIDs.value.size === linksStore.links.length,
);

function toggleSelect(id: string) {
  const s = new Set(selectedIDs.value);
  if (s.has(id)) s.delete(id);
  else s.add(id);
  selectedIDs.value = s;
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIDs.value = new Set();
  } else {
    selectedIDs.value = new Set(linksStore.links.map((l) => l.id));
  }
}

function clearSelection() {
  selectedIDs.value = new Set();
  bulkTagsInput.value = '';
  showBulkTagsBar.value = false;
  showBulkFolderBar.value = false;
}

async function executeBulkAction(action: string, extra: Record<string, unknown> = {}) {
  if (selectedIDs.value.size === 0) return;
  bulkLoading.value = true;
  try {
    await linksApi.bulkAction({
      ids: [...selectedIDs.value],
      action: action as 'delete' | 'activate' | 'deactivate' | 'move_folder' | 'add_tags' | 'remove_tags',
      ...extra,
    });
    clearSelection();
    await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
    await loadTags();
  } catch {
    // silent — store errors are shown by the table error banner
  } finally {
    bulkLoading.value = false;
  }
}

async function bulkActivate() { await executeBulkAction('activate'); }
async function bulkDeactivate() { await executeBulkAction('deactivate'); }

async function bulkDelete() {
  const n = selectedIDs.value.size;
  if (!window.confirm(`Delete ${n} selected link${n !== 1 ? 's' : ''}? This cannot be undone.`)) return;
  await executeBulkAction('delete');
}

async function bulkMoveFolder(folderID: string | null) {
  await executeBulkAction('move_folder', { folder_id: folderID });
  showBulkFolderBar.value = false;
}

async function bulkApplyTags() {
  const tags = bulkTagsInput.value.split(',').map((t) => t.trim()).filter(Boolean);
  if (!tags.length) return;
  await executeBulkAction(bulkTagsMode.value, { tags });
  bulkTagsInput.value = '';
  showBulkTagsBar.value = false;
}

async function downloadBulkQR() {
  if (selectedIDs.value.size === 0) return;
  qrDownloading.value = true;
  try {
    for (const id of selectedIDs.value) {
      const link = linksStore.links.find(l => l.id === id);
      const url = linksApi.getQRUrl(id);
      const res = await fetch(url);
      if (!res.ok) continue;
      const blob = await res.blob();
      const objUrl = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = objUrl;
      a.download = `qr-${link?.slug ?? id}.png`;
      a.click();
      URL.revokeObjectURL(objUrl);
      await new Promise(r => setTimeout(r, 150));
    }
  } finally {
    qrDownloading.value = false;
  }
}

// ── Folder / filter state ────────────────────────────────────────────────────
const folders = ref<FolderResponse[]>([]);
const selectedFolderID = ref('');
const starredOnly = ref(false);
const healthFilter = ref('');
const expiringSoonFilter = ref(false);
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
  expiringSoonFilter.value = false;
  selectedTags.value = [];
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, id);
}

function selectStarred() {
  starredOnly.value = true;
  selectedFolderID.value = '';
  healthFilter.value = '';
  expiringSoonFilter.value = false;
  selectedTags.value = [];
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', true);
}

function selectHealthFilter(status: string) {
  healthFilter.value = status;
  selectedFolderID.value = '';
  starredOnly.value = false;
  expiringSoonFilter.value = false;
  selectedTags.value = [];
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', undefined, status);
}

function selectExpiringSoon() {
  expiringSoonFilter.value = true;
  selectedFolderID.value = '';
  starredOnly.value = false;
  healthFilter.value = '';
  selectedTags.value = [];
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', undefined, undefined, undefined, true);
}

async function loadTags() {
  try {
    const res = await linksApi.getAllTags();
    if (res.data) allTags.value = res.data;
  } catch {
    // non-critical
  }
}

function toggleTag(tag: string) {
  const idx = selectedTags.value.indexOf(tag);
  if (idx === -1) {
    selectedTags.value = [...selectedTags.value, tag];
  } else {
    selectedTags.value = selectedTags.value.filter((t) => t !== tag);
  }
  selectedFolderID.value = '';
  starredOnly.value = false;
  healthFilter.value = '';
  expiringSoonFilter.value = false;
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', undefined, undefined, selectedTags.value.length ? selectedTags.value : undefined);
}

function clearTagFilter() {
  selectedTags.value = [];
  linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, undefined, expiringSoonFilter.value || undefined);
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
      linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, '', starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
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
    linksStore.fetchLinks(1, linksStore.limit, val, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
  }, 400);
});

onMounted(async () => {
  linksStore.fetchLinks(1, 20, '');
  loadFolders();
  loadTags();
  // Bootstrap Toast removed — will use Vuetify snackbar in page rewrite
  try {
    const res = await billingApi.getUsage();
    usage.value = res.data.data;
  } catch {}
  try {
    const res = await billingApi.getSubscription();
    plan.value = res.data.data.plan;
  } catch {}
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

async function exportCSV() {
  exporting.value = true;
  try {
    await linksApi.exportCSV();
  } catch {
    // silently ignore — browser will show no download
  } finally {
    exporting.value = false;
  }
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

function openPreviewModal(link: LinkResponse) {
  previewLink.value = link;
  setTimeout(() => previewModalRef.value?.show(), 50);
}

function openCloneModal(link: LinkResponse) {
  cloneLink.value = link;
  setTimeout(() => cloneModalRef.value?.show(), 50);
}

async function onLinkCloned(newLink: LinkResponse) {
  linksStore.links.unshift(newLink);
  await loadTags();
}

async function onImportDone(_result: ImportLinksResponse) {
  await linksStore.fetchLinks(1, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
  await loadTags();
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
  await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
  editingLink.value = null;
  await loadTags();
}

async function toggleStar(link: LinkResponse) {
  await linksStore.toggleStar(link.id);
  if (starredOnly.value && !linksStore.links.find(l => l.id === link.id)?.is_starred) {
    await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, '', true, healthFilter.value || undefined, undefined, expiringSoonFilter.value || undefined);
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
        await linksStore.fetchLinks(linksStore.page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, 'unhealthy', undefined, expiringSoonFilter.value || undefined);
      }
    }
  } finally {
    checkingHealthId.value = null;
  }
}

function isExpired(link: LinkResponse): boolean {
  return !!link.expires_at && new Date(link.expires_at) < new Date();
}

function isExpiringSoon(link: LinkResponse): boolean {
  if (!link.expires_at) return false;
  const exp = new Date(link.expires_at);
  const now = new Date();
  const threeDays = 3 * 24 * 60 * 60 * 1000;
  return exp > now && exp.getTime() - now.getTime() <= threeDays;
}

function isClickLimitReached(link: LinkResponse): boolean {
  return !!link.max_clicks && link.click_count >= link.max_clicks;
}

function statusBadgeClass(link: LinkResponse): string {
  if (isExpired(link)) return 'text-bg-secondary';
  if (isClickLimitReached(link)) return 'text-bg-warning';
  return link.is_active ? 'text-bg-success' : 'text-bg-danger';
}

function statusM3BadgeClass(link: LinkResponse): string {
  if (isExpired(link)) return 'm3-badge--neutral';
  if (isClickLimitReached(link)) return 'm3-badge--warning';
  return link.is_active ? 'm3-badge--success' : 'm3-badge--error';
}

function statusLabel(link: LinkResponse): string {
  if (isExpired(link)) return 'Expired';
  if (isClickLimitReached(link)) return 'Limit Reached';
  return link.is_active ? 'Active' : 'Inactive';
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

function healthM3BadgeClass(status: string): string {
  switch (status) {
    case 'healthy':   return 'm3-badge--success';
    case 'unhealthy': return 'm3-badge--error';
    case 'timeout':   return 'm3-badge--warning';
    case 'error':     return 'm3-badge--neutral';
    default:          return 'm3-badge--neutral';
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
    await linksStore.fetchLinks(targetPage, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
    await loadTags();
  } finally {
    deletingId.value = null;
  }
}

function goToPage(page: number) {
  if (page < 1 || page > linksStore.totalPages) return;
  linksStore.fetchLinks(page, linksStore.limit, searchQuery.value, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
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
/* ── Layout ──────────────────────────────────────────────────────────────────── */
.links-page {
  padding: 24px;
  min-height: 100%;
}

.links-layout {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 20px;
  align-items: start;
}

@media (max-width: 768px) {
  .links-layout {
    grid-template-columns: 1fr;
  }
  .folder-sidebar {
    display: none;
  }
}

/* ── Sidebar ─────────────────────────────────────────────────────────────────── */
.folder-sidebar {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sidebar-card {
  overflow: hidden;
}

.sidebar-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.sidebar-title {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.sidebar-list {
  padding: 4px 0;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 16px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
  text-align: left;
  transition: background 0.15s;
}

.sidebar-item:hover {
  background: var(--md-sys-color-surface-container-low);
}

.sidebar-item--active {
  background: rgba(99, 91, 255, 0.08);
  color: var(--md-sys-color-primary);
  font-weight: 600;
}

.sidebar-item-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.folder-row {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px 8px;
  border-radius: 6px;
  transition: background 0.15s;
}

.folder-row:hover {
  background: rgba(0, 0, 0, 0.04);
}

.folder-row--active {
  background: rgba(99, 91, 255, 0.08);
}

.folder-row-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  padding: 6px 4px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
  min-width: 0;
  text-align: left;
}

.folder-row-btn--active {
  color: var(--md-sys-color-primary);
  font-weight: 600;
}

.folder-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  display: inline-block;
}

.folder-row-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.folder-row-count {
  font-size: 0.7rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-left: auto;
}

.folder-actions {
  display: flex;
  gap: 0;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.15s;
}

.folder-row:hover .folder-actions {
  opacity: 1;
}

.sidebar-rename-row {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
}

.sidebar-footer {
  padding: 8px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
}

.new-folder-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.sidebar-tags-card {
  overflow: hidden;
}

.sidebar-tags-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.sidebar-tags-body {
  padding: 8px 12px 12px;
}

/* ── Page Header ─────────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
}

.page-title {
  font-size: 1.25rem;
  font-weight: 700;
  margin: 0 0 2px;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.page-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

/* ── Banners ─────────────────────────────────────────────────────────────────── */
.usage-banner,
.error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 16px;
  border-radius: 12px;
}

.usage-banner-msg {
  flex: 1;
  font-size: 0.875rem;
}

/* ── Loading / Empty ─────────────────────────────────────────────────────────── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  gap: 16px;
}

.loading-text {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.9rem;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 64px 24px;
  text-align: center;
}

.empty-state-icon {
  font-size: 56px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.4;
  margin-bottom: 16px;
}

.empty-state-title {
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 8px;
}

.empty-state-subtitle {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 24px;
  max-width: 360px;
}

/* ── Table Card ──────────────────────────────────────────────────────────────── */
.links-table-card {
  overflow: hidden;
}

/* ── Bulk bar ────────────────────────────────────────────────────────────────── */
.bulk-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  padding: 10px 16px;
  background: rgba(99, 91, 255, 0.05);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.bulk-bar-count {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-primary);
}

.bulk-bar-end {
  margin-left: auto;
  display: flex;
  gap: 8px;
}

.bulk-popover-wrap {
  position: relative;
  display: flex;
  align-items: center;
  gap: 4px;
}

.bulk-popover {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  z-index: 200;
  min-width: 160px;
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
  padding: 4px;
}

.bulk-tags-popover {
  display: flex;
  gap: 8px;
  min-width: 280px;
  padding: 8px;
}

.bulk-popover-item {
  display: block;
  width: 100%;
  padding: 8px 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.875rem;
  text-align: left;
  border-radius: 6px;
  color: var(--md-sys-color-on-surface);
}

.bulk-popover-item:hover {
  background: var(--md-sys-color-surface-container-low);
}

/* ── Table scroll ────────────────────────────────────────────────────────────── */
.table-scroll {
  overflow-x: auto;
}

/* ── Row states ──────────────────────────────────────────────────────────────── */
.row--selected {
  background: rgba(99, 91, 255, 0.04);
}

/* ── Link cell content ───────────────────────────────────────────────────────── */
.link-title {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--md-sys-color-on-surface);
}

.link-dest {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: 1px;
}

.link-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 4px;
}

.link-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
  font-size: 0.72rem;
  color: var(--md-sys-color-on-surface-variant);
}

.link-meta-item {
  display: flex;
  align-items: center;
  gap: 2px;
}

/* ── Short URL cell ──────────────────────────────────────────────────────────── */
.short-url-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
}

.short-url-link {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  font-weight: 500;
  font-size: 0.875rem;
}

.short-url-link:hover {
  text-decoration: underline;
}

/* ── Other cells ─────────────────────────────────────────────────────────────── */
.click-count {
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.date-cell {
  font-size: 0.82rem;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
}

.actions-cell {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0;
}

/* ── Pagination ──────────────────────────────────────────────────────────────── */
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
}

.pagination-info {
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 6px;
}

.pagination-ellipsis {
  padding: 0 6px;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Snackbar ────────────────────────────────────────────────────────────────── */
.m3-snackbar {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1100;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #1C1B1F;
  color: #E6E1E5;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.3);
  min-width: 280px;
  max-width: 400px;
}

.snack-enter-active,
.snack-leave-active { transition: all 0.25s; }
.snack-enter-from,
.snack-leave-to { transform: translateY(80px); opacity: 0; }
</style>
