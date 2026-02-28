<template>
  <div class="links-page">

    <!-- ── Top Toolbar ─────────────────────────────────────────────────── -->
    <div class="links-toolbar">
      <div class="toolbar-left">
        <h1 class="page-title">My Links</h1>
        <span class="total-badge">{{ linksStore.total }}</span>
      </div>
      <div class="toolbar-right">
        <div class="search-wrap">
          <span class="search-icon material-symbols-outlined">search</span>
          <input
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="Search links…"
            autocomplete="off"
          />
          <button v-if="searchQuery" class="search-clear-btn" @click="searchQuery = ''">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="toolbar-more-wrap">
          <button class="toolbar-icon-btn" title="More options" @click.stop="showToolbarMore = !showToolbarMore">
            <span class="material-symbols-outlined">more_vert</span>
          </button>
          <Transition name="menu-pop">
            <div v-if="showToolbarMore" class="toolbar-more-menu" @click.stop>
              <button class="toolbar-menu-item" @click="exportCSV(); showToolbarMore = false" :disabled="exporting">
                <md-circular-progress v-if="exporting" indeterminate style="width:18px;height:18px" />
                <span v-else class="material-symbols-outlined">download</span>
                Export CSV
              </button>
              <button class="toolbar-menu-item" @click="openImportModal(); showToolbarMore = false">
                <span class="material-symbols-outlined">upload_file</span>
                Import CSV
              </button>
              <div class="toolbar-menu-divider" />
              <button class="toolbar-menu-item" @click="showFolderManager = true; showToolbarMore = false">
                <span class="material-symbols-outlined">create_new_folder</span>
                Manage Folders
              </button>
            </div>
          </Transition>
        </div>
        <button class="btn-filled new-link-btn" @click="openCreateModal">
          <span class="material-symbols-outlined">add</span>
          New Link
        </button>
      </div>
    </div>

    <!-- ── Filter Chips ─────────────────────────────────────────────────── -->
    <div class="filter-chips-bar">
      <div class="filter-chips-scroll">
        <button class="m3-chip" :class="{ 'active': isAllActive }" @click="selectFolder('')">
          <span v-if="isAllActive" class="material-symbols-outlined chip-check-icon">check</span>
          All
        </button>
        <button class="m3-chip" :class="{ 'active': starredOnly }" @click="selectStarred">
          <span v-if="starredOnly" class="material-symbols-outlined chip-check-icon">check</span>
          <span v-else class="material-symbols-outlined chip-lead-icon">star</span>
          Starred
        </button>
        <button class="m3-chip" :class="{ 'active': healthFilter === 'unhealthy' }" @click="selectHealthFilter('unhealthy')">
          <span v-if="healthFilter === 'unhealthy'" class="material-symbols-outlined chip-check-icon">check</span>
          <span v-else class="material-symbols-outlined chip-lead-icon chip-icon--error">error</span>
          Unhealthy
        </button>
        <button class="m3-chip" :class="{ 'active': expiringSoonFilter }" @click="selectExpiringSoon">
          <span v-if="expiringSoonFilter" class="material-symbols-outlined chip-check-icon">check</span>
          <span v-else class="material-symbols-outlined chip-lead-icon chip-icon--warn">schedule</span>
          Expiring
        </button>
        <div v-if="folders.length" class="chip-divider"></div>
        <button
          v-for="folder in folders"
          :key="folder.id"
          class="m3-chip m3-chip--folder"
          :class="{ 'active': selectedFolderID === folder.id }"
          @click="selectFolder(folder.id)"
        >
          <span v-if="selectedFolderID === folder.id" class="material-symbols-outlined chip-check-icon">check</span>
          <span v-else class="folder-dot-chip" :style="{ backgroundColor: folder.color }"></span>
          {{ folder.name }}
        </button>
        <button
          v-if="allTags.length"
          class="m3-chip m3-chip--tags"
          :class="{ 'active': showTagFilter || selectedTags.length > 0 }"
          @click="showTagFilter = !showTagFilter"
        >
          <span v-if="showTagFilter || selectedTags.length" class="material-symbols-outlined chip-check-icon">check</span>
          <span v-else class="material-symbols-outlined chip-lead-icon">label</span>
          Tags
          <span v-if="selectedTags.length" class="chip-count">{{ selectedTags.length }}</span>
        </button>
      </div>
    </div>

    <!-- ── Tag Filter Chips ─────────────────────────────────────────────── -->
    <div v-if="allTags.length && (showTagFilter || selectedTags.length > 0)" class="tag-filter-bar">
      <button
        v-for="tag in allTags"
        :key="tag"
        class="tag-chip"
        :class="{ 'tag-chip--active': selectedTags.includes(tag) }"
        @click="toggleTag(tag)"
      >{{ tag }}</button>
      <button v-if="selectedTags.length" class="tag-filter-clear-btn" @click="clearTagFilter">
        <span class="material-symbols-outlined">close</span>
        Clear filters
      </button>
    </div>

    <!-- ── Usage Warning ────────────────────────────────────────────────── -->
    <div v-if="usageWarning" class="usage-banner" :class="{ 'usage-banner--danger': usageWarning.level === 'danger' }">
      <span class="material-symbols-outlined">{{ usageWarning.level === 'danger' ? 'block' : 'warning' }}</span>
      <span class="usage-msg">{{ usageWarning.msg }}</span>
      <router-link to="/dashboard/billing">
        <button class="btn-tonal btn-sm">Upgrade</button>
      </router-link>
    </div>

    <!-- ── Error ────────────────────────────────────────────────────────── -->
    <div v-if="linksStore.error" class="error-banner">
      <span class="material-symbols-outlined">error</span>
      {{ linksStore.error }}
    </div>

    <!-- ── Bulk Action Bar ──────────────────────────────────────────────── -->
    <div v-if="selectedIDs.size > 0" class="bulk-bar">
      <span class="bulk-count">{{ selectedIDs.size }} selected</span>
      <div class="bulk-actions">
        <button class="btn-outlined btn-sm" :disabled="bulkLoading" @click="bulkActivate">Activate</button>
        <button class="btn-outlined btn-sm" :disabled="bulkLoading" @click="bulkDeactivate">Deactivate</button>
        <div class="bulk-popover-wrap">
          <button class="btn-outlined btn-sm" :disabled="bulkLoading" @click="showBulkFolderBar = !showBulkFolderBar; showBulkTagsBar = false">
            <span class="material-symbols-outlined" style="font-size:16px">drive_file_move</span>Move
          </button>
          <div v-if="showBulkFolderBar" class="bulk-popover">
            <button class="bulk-popover-item" @click="bulkMoveFolder(null)">— Remove from folder</button>
            <button v-for="f in folders" :key="f.id" class="bulk-popover-item" @click="bulkMoveFolder(f.id)">{{ f.name }}</button>
          </div>
        </div>
        <div class="bulk-popover-wrap">
          <button class="btn-outlined btn-sm" :disabled="bulkLoading" @click="bulkTagsMode = 'add_tags'; showBulkTagsBar = !showBulkTagsBar; showBulkFolderBar = false">
            <span class="material-symbols-outlined" style="font-size:16px">label</span>Tags
          </button>
          <div v-if="showBulkTagsBar" class="bulk-popover bulk-tags-popover">
            <md-outlined-text-field
              :value="bulkTagsInput"
              @input="bulkTagsInput = ($event.target as HTMLInputElement).value"
              label="tag1, tag2, ..."
              style="flex:1"
              @keydown.enter="bulkApplyTags"
            />
            <button class="btn-filled btn-sm" @click="bulkApplyTags">{{ bulkTagsMode === 'add_tags' ? 'Add' : 'Remove' }}</button>
          </div>
        </div>
        <button class="btn-outlined btn-sm" :disabled="qrDownloading" @click="downloadBulkQR">
          <md-circular-progress v-if="qrDownloading" indeterminate />
          <span v-else class="material-symbols-outlined" style="font-size:16px">qr_code_2</span>QR
        </button>
        <button class="btn-outlined btn-sm btn-danger" :disabled="bulkLoading" @click="bulkDelete">
          <md-circular-progress v-if="bulkLoading" indeterminate />
          <span v-else class="material-symbols-outlined" style="font-size:16px">delete</span>Delete
        </button>
      </div>
      <button class="btn-text btn-sm" @click="clearSelection">Clear</button>
    </div>

    <!-- ── Loading Skeleton ─────────────────────────────────────────────── -->
    <div v-if="linksStore.loading" class="skeleton-list" aria-busy="true" aria-label="Loading links">
      <div v-for="i in 5" :key="i" class="skeleton-card">
        <div class="skeleton skel-icon"></div>
        <div class="skel-body">
          <div class="skeleton skel-title"></div>
          <div class="skeleton skel-url"></div>
          <div class="skeleton skel-short"></div>
        </div>
        <div class="skeleton skel-stat"></div>
        <div class="skeleton skel-actions"></div>
      </div>
    </div>

    <!-- ── Empty State ──────────────────────────────────────────────────── -->
    <div v-else-if="!linksStore.loading && linksStore.links.length === 0" class="empty-state">
      <div class="empty-icon-wrap">
        <span class="material-symbols-outlined empty-icon">
          {{ searchQuery ? 'search_off' : selectedTags.length ? 'label_off' : 'add_link' }}
        </span>
      </div>
      <h2 class="empty-title">
        {{ searchQuery ? 'No links match your search' : selectedTags.length ? 'No links with these tags' : 'No links yet' }}
      </h2>
      <p class="empty-sub">
        {{ searchQuery ? 'Try a different search term or clear the search.' : selectedTags.length ? 'Try different tags or clear the filter.' : 'Shorten your first URL and start tracking clicks in seconds.' }}
      </p>
      <div class="empty-actions">
        <button class="btn-filled" v-if="!searchQuery && !selectedTags.length" @click="openCreateModal">
          <span class="material-symbols-outlined">add</span>Create your first link
        </button>
        <button class="btn-outlined" v-else-if="selectedTags.length" @click="clearTagFilter">Clear tag filter</button>
        <button class="btn-outlined" v-else @click="searchQuery = ''">Clear search</button>
      </div>
    </div>

    <!-- ── Link Card List ───────────────────────────────────────────────── -->
    <div v-else class="link-list">
      <div class="list-header">
        <label class="select-all-label">
          <md-checkbox
            :checked="isAllSelected"
            :indeterminate="selectedIDs.size > 0 && !isAllSelected"
            @change="toggleSelectAll"
          />
          <span class="select-all-text">{{ isAllSelected ? 'Deselect all' : 'Select all' }}</span>
        </label>
        <span class="list-count">{{ linksStore.total }} link{{ linksStore.total !== 1 ? 's' : '' }}</span>
      </div>

      <div
        v-for="link in linksStore.links"
        :key="link.id"
        class="link-card"
        :class="{ 'link-card--selected': selectedIDs.has(link.id) }"
      >
        <!-- Checkbox -->
        <div class="card-check">
          <md-checkbox :checked="selectedIDs.has(link.id)" @change="toggleSelect(link.id)" />
        </div>

        <!-- Icon letter -->
        <div class="card-icon" :style="{ background: getIconColor(link) }">
          {{ getFaviconLetter(link) }}
        </div>

        <!-- Main body -->
        <div class="card-main">
          <div class="card-title-row">
            <span class="card-title" :title="link.title || link.destination_url">
              {{ link.title || getDomain(link.destination_url) }}
            </span>
            <button
              class="star-btn"
              :class="{ 'star-btn--active': link.is_starred }"
              :title="link.is_starred ? 'Unstar' : 'Star'"
              @click="toggleStar(link)"
            >
              <span class="material-symbols-outlined" :style="{ fontVariationSettings: link.is_starred ? '\'FILL\' 1' : '\'FILL\' 0' }">star</span>
            </button>
          </div>

          <div class="card-dest" :title="link.destination_url">{{ link.destination_url }}</div>

          <div class="card-short-row">
            <a :href="link.short_url" target="_blank" rel="noopener noreferrer" class="card-short-url">
              <span class="material-symbols-outlined" style="font-size:13px">link</span>
              {{ link.short_url }}
            </a>
            <button class="copy-btn" :title="copiedId === link.id ? 'Copied!' : 'Copy'" @click="copyShortUrl(link)">
              <span class="material-symbols-outlined" :style="{ color: copiedId === link.id ? '#1AA563' : undefined }">
                {{ copiedId === link.id ? 'check_circle' : 'content_copy' }}
              </span>
            </button>
            <span :class="['m3-badge', statusM3BadgeClass(link)]">{{ statusLabel(link) }}</span>
            <span
              :class="['m3-badge', healthM3BadgeClass(link.health_status)]"
              :title="link.health_checked_at ? `Checked: ${formatDate(link.health_checked_at)}` : 'Not yet checked'"
            >{{ healthLabel(link.health_status) }}</span>
            <span v-if="link.has_password" class="meta-pill" title="Password protected">
              <span class="material-symbols-outlined" style="font-size:12px">lock</span>
            </span>
            <span
              v-if="expiryBadgeClass(link)"
              :class="['meta-pill', expiryBadgeClass(link) === 'text-bg-danger' ? 'meta-pill--error' : 'meta-pill--warn']"
              :title="isExpired(link) ? 'Expired' : 'Expiring soon'"
            >
              <span class="material-symbols-outlined" style="font-size:12px">schedule</span>
              {{ expiryBadgeLabel(link) }}
            </span>
            <span
              v-if="clickLimitBadgeClass(link)"
              :class="['meta-pill', clickLimitBadgeClass(link) === 'text-bg-danger' ? 'meta-pill--error' : 'meta-pill--warn']"
              title="Click limit"
            >
              <span class="material-symbols-outlined" style="font-size:12px">bolt</span>
              {{ clickLimitBadgeLabel(link) }}
            </span>
          </div>

          <div v-if="link.tags && link.tags.length > 0" class="card-tags">
            <span v-for="tag in link.tags.slice(0, 4)" :key="tag" class="tag-pill">{{ tag }}</span>
            <span v-if="link.tags.length > 4" class="tag-pill tag-pill--more">+{{ link.tags.length - 4 }}</span>
          </div>
        </div>

        <!-- Stats -->
        <div class="card-stats">
          <div class="stat-value">{{ link.click_count.toLocaleString() }}</div>
          <div class="stat-label">clicks</div>
          <div class="stat-date">{{ formatDate(link.created_at) }}</div>
        </div>

        <!-- Actions -->
        <div class="card-actions">
          <RouterLink :to="`/dashboard/links/${link.id}`">
            <button class="action-btn" title="Analytics">
              <span class="material-symbols-outlined">bar_chart</span>
            </button>
          </RouterLink>
          <button class="action-btn" title="Edit" @click="openEditModal(link)">
            <span class="material-symbols-outlined">edit</span>
          </button>
          <button
            class="action-btn action-btn--danger"
            title="Delete"
            :disabled="deletingId === link.id"
            @click="confirmDelete(link)"
          >
            <md-circular-progress v-if="deletingId === link.id" indeterminate />
            <span v-else class="material-symbols-outlined">delete</span>
          </button>
          <div class="row-menu-wrap">
            <button class="action-btn" title="More actions" @click.stop="toggleRowMenu($event, link.id)">
              <span class="material-symbols-outlined">more_vert</span>
            </button>
            <Teleport to="body">
              <div
                v-if="openMenuId === link.id"
                class="row-menu"
                :style="{ top: menuPos.top + 'px', right: menuPos.right + 'px' }"
                @click.stop
              >
                <button class="row-menu-item" @click="toggleStar(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon" :style="{ color: link.is_starred ? '#F4A100' : undefined, fontVariationSettings: link.is_starred ? '\'FILL\' 1' : '\'FILL\' 0' }">star</span>
                  {{ link.is_starred ? 'Unstar' : 'Star' }}
                </button>
                <button class="row-menu-item" :disabled="checkingHealthId === link.id" @click="runHealthCheck(link); closeAllRowMenus()">
                  <md-circular-progress v-if="checkingHealthId === link.id" indeterminate style="margin-right:2px" />
                  <span v-else class="material-symbols-outlined row-menu-item__icon">health_and_safety</span>
                  Check Health
                </button>
                <button class="row-menu-item" @click="openPreviewModal(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon">visibility</span>Preview
                </button>
                <button class="row-menu-item" @click="openQRModal(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon">qr_code_2</span>QR Code
                </button>
                <button class="row-menu-item" @click="openCloneModal(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon">file_copy</span>Clone
                </button>
                <div class="row-menu-divider" />
                <button class="row-menu-item" :class="{ 'row-menu-item--active': link.is_split_test }" @click="openSplitModal(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon">join</span>
                  A/B Split Test
                  <span v-if="link.is_split_test" class="row-menu-item__badge">On</span>
                </button>
                <button class="row-menu-item" :class="{ 'row-menu-item--active': link.is_geo_routing }" @click="openGeoModal(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon">public</span>
                  Geo Routing
                  <span v-if="link.is_geo_routing" class="row-menu-item__badge">On</span>
                </button>
                <button class="row-menu-item" :class="{ 'row-menu-item--active': link.is_pixel_tracking }" @click="openPixelModal(link); closeAllRowMenus()">
                  <span class="material-symbols-outlined row-menu-item__icon">track_changes</span>
                  Retargeting Pixels
                  <span v-if="link.is_pixel_tracking" class="row-menu-item__badge">On</span>
                </button>
              </div>
            </Teleport>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Pagination ────────────────────────────────────────────────────── -->
    <div v-if="linksStore.totalPages > 1" class="pagination">
      <span class="pagination-info">Page {{ linksStore.page }} of {{ linksStore.totalPages }} — {{ linksStore.total }} total</span>
      <div class="pagination-controls">
        <button class="page-btn" :disabled="linksStore.page <= 1" @click="goToPage(linksStore.page - 1)">
          <span class="material-symbols-outlined">chevron_left</span>
        </button>
        <template v-for="(p, idx) in visiblePages" :key="idx">
          <button
            v-if="p !== '...'"
            class="page-btn"
            :class="{ 'page-btn--active': p === linksStore.page }"
            @click="goToPage(Number(p))"
          >{{ p }}</button>
          <span v-else class="page-ellipsis">…</span>
        </template>
        <button class="page-btn" :disabled="linksStore.page >= linksStore.totalPages" @click="goToPage(linksStore.page + 1)">
          <span class="material-symbols-outlined">chevron_right</span>
        </button>
      </div>
    </div>

    <!-- ── Folder Manager Panel ─────────────────────────────────────────── -->
    <Teleport to="body">
      <Transition name="panel">
        <div v-if="showFolderManager" class="fm-overlay" @click="showFolderManager = false">
          <div class="fm-panel" @click.stop>
            <div class="fm-header">
              <span class="fm-title">Manage Folders</span>
              <button class="action-btn" @click="showFolderManager = false">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>
            <div class="fm-body">
              <template v-for="folder in folders" :key="folder.id">
                <div v-if="renamingFolderID === folder.id" class="fm-item-rename">
                  <input
                    class="fm-input"
                    :value="renameValue"
                    @input="renameValue = ($event.target as HTMLInputElement).value"
                    @keydown.enter="submitRename(folder)"
                    @keydown.esc="cancelRename"
                    ref="renameInputRef"
                    placeholder="Folder name"
                  />
                  <button class="action-btn" @click="submitRename(folder)"><span class="material-symbols-outlined">check</span></button>
                  <button class="action-btn" @click="cancelRename"><span class="material-symbols-outlined">close</span></button>
                </div>
                <div v-else class="fm-item">
                  <button
                    class="fm-item-btn"
                    :class="{ 'fm-item-btn--active': selectedFolderID === folder.id }"
                    @click="selectFolder(folder.id); showFolderManager = false"
                  >
                    <span class="folder-dot-sm" :style="{ backgroundColor: folder.color }"></span>
                    <span class="fm-item-name">{{ folder.name }}</span>
                    <span class="fm-item-count">{{ folder.click_count?.toLocaleString() ?? 0 }}</span>
                  </button>
                  <button class="action-btn" title="Rename" @click.stop="startRename(folder)">
                    <span class="material-symbols-outlined" style="font-size:16px">edit</span>
                  </button>
                  <button class="action-btn action-btn--danger" title="Delete folder" @click.stop="deleteFolder(folder)">
                    <span class="material-symbols-outlined" style="font-size:16px">delete</span>
                  </button>
                </div>
              </template>
              <div v-if="folders.length === 0" class="fm-empty">No folders yet. Create your first one below.</div>
            </div>
            <div class="fm-footer">
              <div v-if="showNewFolderInput" class="fm-new-row">
                <input
                  class="fm-input"
                  :value="newFolderName"
                  @input="newFolderName = ($event.target as HTMLInputElement).value"
                  placeholder="Folder name"
                  @keydown.enter="createFolder"
                  @keydown.esc="showNewFolderInput = false; newFolderName = ''"
                  ref="newFolderInputRef"
                />
                <button class="action-btn" @click="createFolder" :disabled="!newFolderName.trim()">
                  <span class="material-symbols-outlined">check</span>
                </button>
              </div>
              <button v-else class="btn-outlined" style="width:100%" @click="openNewFolderInput">
                <span class="material-symbols-outlined">create_new_folder</span>
                New Folder
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- ── Modals ────────────────────────────────────────────────────────── -->
    <CreateLinkModal
      ref="createModalRef"
      :link="editingLink ?? undefined"
      :folders="folders"
      @saved="onLinkSaved"
    />
    <QRCodeModal v-if="qrLink" ref="qrModalRef" :link-id="qrLink.id" :slug="qrLink.slug" />
    <ImportLinksModal ref="importModalRef" @imported="onImportDone" />
    <SplitTestModal
      v-if="splitLink"
      ref="splitModalRef"
      :link-id="splitLink.id"
      :slug="splitLink.slug"
      :is-split-test-initial="splitLink.is_split_test"
      @updated="onSplitTestUpdated"
    />
    <GeoRulesModal
      v-if="geoLink"
      ref="geoModalRef"
      :link-id="geoLink.id"
      :slug="geoLink.slug"
      :is-geo-routing-initial="geoLink.is_geo_routing"
      @updated="onGeoRoutingUpdated"
    />
    <PixelsModal
      v-if="pixelLink"
      ref="pixelModalRef"
      modal-id="pixels-modal"
      :link="pixelLink"
      @pixel-tracking-updated="onPixelTrackingUpdated"
    />
    <LinkPreviewModal v-if="previewLink" ref="previewModalRef" modal-id="preview-modal" :link="previewLink" />
    <CloneLinkModal v-if="cloneLink" ref="cloneModalRef" modal-id="clone-link-modal" :link="cloneLink" @cloned="onLinkCloned" />

    <!-- ── Copy Snackbar ─────────────────────────────────────────────────── -->
    <Transition name="snack">
      <div v-if="copiedId" class="m3-snackbar">
        <span class="material-symbols-outlined" style="font-size:20px">check_circle</span>
        <span style="flex:1">Copied to clipboard!</span>
        <button class="btn-text" @click="copiedId = null">Dismiss</button>
      </div>
    </Transition>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch, computed } from 'vue';
import { RouterLink, useRoute } from 'vue-router';
import { useLinksStore } from '@/stores/links';
import { useUiStore } from '@/stores/ui';
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

const route = useRoute();
const linksStore = useLinksStore();
const uiStore = useUiStore();

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
const showFolderManager = ref(false);
const showTagFilter = ref(false);

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
const openMenuId = ref<string | null>(null);
const menuPos = ref({ top: 0, right: 0 });
const showToolbarMore = ref(false);

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
    // silent
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

// ── Folder / filter state ─────────────────────────────────────────────────────
const folders = ref<FolderResponse[]>([]);
const selectedFolderID = ref('');
const starredOnly = ref(false);
const healthFilter = ref('');
const expiringSoonFilter = ref(false);
const checkingHealthId = ref<string | null>(null);

const isAllActive = computed(
  () => selectedFolderID.value === '' && !starredOnly.value && healthFilter.value === '' && !expiringSoonFilter.value,
);

const showNewFolderInput = ref(false);
const newFolderName = ref('');
const newFolderInputRef = ref<HTMLInputElement | null>(null);

const renamingFolderID = ref('');
const renameValue = ref('');
const renameInputRef = ref<HTMLInputElement | null>(null);

async function loadFolders() {
  try {
    const res = await foldersApi.list();
    if (res.data) folders.value = res.data;
  } catch {
    // non-critical
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

// ── Icon helpers ──────────────────────────────────────────────────────────────
const ICON_COLORS = ['#635BFF', '#0d6efd', '#198754', '#14B8A6', '#dc3545', '#fd7e14', '#6f42c1', '#0dcaf0'];

function getIconColor(link: LinkResponse): string {
  let sum = 0;
  for (let i = 0; i < link.id.length; i++) sum += link.id.charCodeAt(i);
  return ICON_COLORS[sum % ICON_COLORS.length] ?? '#635BFF';
}

function getFaviconLetter(link: LinkResponse): string {
  if (link.title) return link.title.charAt(0).toUpperCase();
  try { return new URL(link.destination_url).hostname.charAt(0).toUpperCase(); } catch { return 'L'; }
}

function getDomain(url: string): string {
  try { return new URL(url).hostname; } catch { return url; }
}

// Debounced search
let searchTimer: ReturnType<typeof setTimeout> | null = null;
watch(searchQuery, (val) => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    linksStore.fetchLinks(1, linksStore.limit, val, selectedFolderID.value, starredOnly.value || undefined, healthFilter.value || undefined, selectedTags.value.length ? selectedTags.value : undefined, expiringSoonFilter.value || undefined);
  }, 400);
});

function toggleRowMenu(event: MouseEvent, id: string) {
  if (openMenuId.value === id) {
    openMenuId.value = null;
    return;
  }
  const btn = event.currentTarget as HTMLElement;
  const rect = btn.getBoundingClientRect();
  menuPos.value = {
    top: rect.bottom + 4,
    right: window.innerWidth - rect.right,
  };
  openMenuId.value = id;
}

function closeAllRowMenus() {
  openMenuId.value = null;
  showToolbarMore.value = false;
}

onMounted(async () => {
  await linksStore.fetchLinks(1, 20, '');
  loadFolders();
  loadTags();
  document.addEventListener('click', closeAllRowMenus);
  // Open create modal triggered by nav button or post-login redirect
  if (uiStore.pendingCreateLink || route.query.create === '1') {
    uiStore.clearPendingCreateLink();
    nextTick(() => openCreateModal());
  }
  try {
    const res = await billingApi.getUsage();
    if (res.data.data) usage.value = res.data.data;
  } catch {}
  try {
    const res = await billingApi.getSubscription();
    if (res.data.data) plan.value = res.data.data.plan;
  } catch {}
});

onUnmounted(() => {
  document.removeEventListener('click', closeAllRowMenus);
});

// When the nav button is clicked while already on this page, open modal immediately
watch(() => uiStore.pendingCreateLink, (pending) => {
  if (pending) {
    uiStore.clearPendingCreateLink();
    nextTick(() => openCreateModal());
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
    setTimeout(() => {
      if (copiedId.value === link.id) copiedId.value = null;
    }, 2000);
  } catch {
    // silent
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
    // silent
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
    if (idx !== -1) linksStore.links[idx] = { ...linksStore.links[idx]!, is_split_test: enabled };
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
    if (idx !== -1) linksStore.links[idx] = { ...linksStore.links[idx]!, is_geo_routing: enabled };
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
    if (idx !== -1) linksStore.links[idx] = { ...linksStore.links[idx]!, is_pixel_tracking: enabled };
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

<style scoped lang="scss">
/* ── Page ───────────────────────────────────────────────────────────────────── */
.links-page {
  display: flex;
  flex-direction: column;
  gap: 0;
  min-height: 100%;
}

/* ── Toolbar ────────────────────────────────────────────────────────────────── */
.links-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 20px 0;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0;
  letter-spacing: -0.01em;
}

.total-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 22px;
  padding: 0 8px;
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.toolbar-icon-btn {
  width: 38px;
  height: 38px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface-variant);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, border-color 0.15s, color 0.15s;

  .material-symbols-outlined { font-size: 20px; }

  &:hover {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
    border-color: var(--md-sys-color-outline);
  }

  &:disabled { opacity: 0.5; cursor: default; }
}

.new-link-btn {
  height: 38px;
  padding: 0 18px;
  font-size: 0.875rem;
}

/* ── Search ─────────────────────────────────────────────────────────────────── */
.search-wrap {
  position: relative;
  display: flex;
  align-items: center;
  width: 260px;
}

.search-icon {
  position: absolute;
  left: 10px;
  font-size: 18px;
  color: var(--md-sys-color-on-surface-variant);
  pointer-events: none;
  z-index: 1;
}

.search-input {
  width: 100%;
  height: 38px;
  padding: 0 32px 0 36px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface);
  font-family: 'Roboto', sans-serif;
  font-size: 14px;
  outline: none;
  transition: border-color 0.15s, background 0.15s;

  &::placeholder { color: var(--md-sys-color-on-surface-variant); }
  &:focus {
    border-color: var(--md-sys-color-primary);
    border-width: 2px;
    background: var(--md-sys-color-surface);
  }
}

.search-clear-btn {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--md-sys-color-on-surface-variant);
  .material-symbols-outlined { font-size: 16px; }
  &:hover { background: var(--md-sys-color-surface-container); }
}

/* ── Filter Chips Bar ───────────────────────────────────────────────────────── */
.filter-chips-bar {
  margin-bottom: 16px;
}

.filter-chips-scroll {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow-x: auto;
  scrollbar-width: none;
  padding-bottom: 4px;
  &::-webkit-scrollbar { display: none; }
}

/* M3 Filter Chip */
.m3-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 32px;
  padding: 0 16px 0 8px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: transparent;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
  transition: background 0.15s, border-color 0.15s, color 0.15s, box-shadow 0.15s;
  font-family: 'Roboto', sans-serif;

  &:hover:not(.active) {
    background: var(--md-sys-color-surface-container-low);
    border-color: var(--md-sys-color-outline);
  }

  &.active {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    border-color: transparent;
    font-weight: 600;
    padding-left: 4px;
  }
}

.chip-check-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.chip-lead-icon {
  font-size: 18px;
  flex-shrink: 0;
  color: var(--md-sys-color-on-surface-variant);

  &.chip-icon--error { color: var(--md-sys-color-error); }
  &.chip-icon--warn  { color: #F4A100; }
}

.folder-dot-chip {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-left: 4px;
  margin-right: 2px;
}

.chip-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border-radius: 999px;
  font-size: 0.65rem;
  font-weight: 700;
  margin-left: 2px;
}

.chip-divider {
  width: 1px;
  height: 20px;
  background: var(--md-sys-color-outline-variant);
  flex-shrink: 0;
  margin: 0 4px;
}

/* ── Toolbar More Dropdown ──────────────────────────────────────────────────── */
.toolbar-more-wrap {
  position: relative;
}

.toolbar-more-menu {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  z-index: 300;
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  padding: 4px;
  min-width: 192px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.14);
}

.toolbar-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 9px 14px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.875rem;
  font-family: 'Roboto', sans-serif;
  color: var(--md-sys-color-on-surface);
  border-radius: 8px;
  text-align: left;
  transition: background 0.12s;

  .material-symbols-outlined { font-size: 18px; color: var(--md-sys-color-on-surface-variant); flex-shrink: 0; }

  &:hover { background: var(--md-sys-color-surface-container); }
  &:disabled { opacity: 0.5; cursor: default; }
}

.toolbar-menu-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  margin: 3px 4px;
}

/* Menu pop animation */
.menu-pop-enter-active { transition: opacity 0.12s ease, transform 0.12s ease; }
.menu-pop-leave-active { transition: opacity 0.1s ease, transform 0.1s ease; }
.menu-pop-enter-from, .menu-pop-leave-to { opacity: 0; transform: translateY(-6px) scale(0.97); }

/* ── Tag Filter Bar ─────────────────────────────────────────────────────────── */
.tag-filter-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  padding: 8px 0 12px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  margin-bottom: 16px;
}

.tag-filter-clear-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 28px;
  padding: 0 10px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.78rem;
  font-weight: 500;
  font-family: 'Roboto', sans-serif;
  color: var(--md-sys-color-on-surface-variant);
  border-radius: 6px;
  transition: background 0.12s, color 0.12s;
  .material-symbols-outlined { font-size: 15px; }
  &:hover { background: var(--md-sys-color-surface-container); color: var(--md-sys-color-on-surface); }
}

.tag-chip {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 999px;
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.78rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, color 0.15s;

  &:hover {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
  }

  &--active {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
    border-color: var(--md-sys-color-primary);
    font-weight: 600;
  }
}

/* ── Banners ────────────────────────────────────────────────────────────────── */
.usage-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: rgba(244, 161, 0, 0.1);
  border: 1px solid rgba(244, 161, 0, 0.4);
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);

  .material-symbols-outlined { color: #F4A100; font-size: 20px; flex-shrink: 0; }

  &--danger {
    background: var(--md-sys-color-error-container);
    border-color: var(--md-sys-color-error);
    .material-symbols-outlined { color: var(--md-sys-color-error); }
  }
}

.usage-msg { flex: 1; }

.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: var(--md-sys-color-error-container);
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-error-container);
  .material-symbols-outlined { color: var(--md-sys-color-error); font-size: 20px; }
}

/* ── Bulk Bar ───────────────────────────────────────────────────────────────── */
.bulk-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: var(--md-sys-color-secondary-container);
  border-radius: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.bulk-count {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-on-secondary-container);
  white-space: nowrap;
}

.bulk-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  flex-wrap: wrap;
}

.bulk-popover-wrap {
  position: relative;
}

.bulk-popover {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  z-index: 200;
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  padding: 4px;
  min-width: 160px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}

.bulk-tags-popover {
  display: flex;
  align-items: center;
  gap: 6px;
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
  color: var(--md-sys-color-on-surface);
  text-align: left;
  border-radius: 6px;
  transition: background 0.15s;

  &:hover { background: var(--md-sys-color-surface-container); }
}

/* ── Skeleton ───────────────────────────────────────────────────────────────── */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-card {
  display: grid;
  grid-template-columns: 44px 1fr 72px 140px;
  gap: 14px;
  align-items: center;
  padding: 16px;
  background: var(--md-sys-color-surface-container-low);
  border-radius: 12px;
}

.skeleton {
  background: linear-gradient(90deg, var(--md-sys-color-surface-container) 25%, var(--md-sys-color-surface-container-high) 50%, var(--md-sys-color-surface-container) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
  border-radius: 6px;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.skel-icon  { width: 44px; height: 44px; border-radius: 10px; }
.skel-body  { display: flex; flex-direction: column; gap: 7px; }
.skel-title { height: 15px; width: 55%; border-radius: 4px; }
.skel-url   { height: 12px; width: 80%; border-radius: 4px; }
.skel-short { height: 12px; width: 65%; border-radius: 4px; }
.skel-stat  { height: 40px; border-radius: 6px; }
.skel-actions { height: 34px; border-radius: 8px; }

/* ── Empty State ────────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80px 20px;
  text-align: center;
}

.empty-icon-wrap {
  width: 96px;
  height: 96px;
  border-radius: 28px;
  background: var(--md-sys-color-surface-container);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}

.empty-icon {
  font-size: 48px;
  color: var(--md-sys-color-primary);
  opacity: 0.55;
}

.empty-title {
  font-size: 1.15rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0 0 8px;
}

.empty-sub {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 24px;
  max-width: 360px;
}

.empty-actions {
  display: flex;
  gap: 10px;
}

/* ── Link List ──────────────────────────────────────────────────────────────── */
.link-list {
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
  background: var(--md-sys-color-surface);
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px 6px 10px;
  background: var(--md-sys-color-surface-container-low);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  min-height: 40px;
}

.select-all-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.select-all-text {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 500;
}

.list-count {
  font-size: 0.78rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Link Card ──────────────────────────────────────────────────────────────── */
.link-card {
  display: grid;
  grid-template-columns: 40px 48px 1fr auto auto;
  gap: 12px;
  align-items: start;
  padding: 14px 14px 14px 10px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface);
  transition: background 0.15s;

  &:last-child { border-bottom: none; }

  &:hover {
    background: var(--md-sys-color-surface-container-low);
    cursor: default;
  }

  &--selected {
    background: color-mix(in srgb, var(--md-sys-color-primary) 8%, var(--md-sys-color-surface));
  }
}

.card-check {
  display: flex;
  align-items: flex-start;
  padding-top: 8px;
}

.card-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.15rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  letter-spacing: 0;
  margin-top: 2px;
}

.card-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

.star-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--md-sys-color-on-surface-variant);
  padding: 0;
  flex-shrink: 0;
  transition: color 0.15s;

  .material-symbols-outlined { font-size: 17px; }
  &:hover { color: #F4A100; }
  &--active { color: #F4A100; }
}

.card-dest {
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 1px;
}

.card-short-row {
  display: flex;
  align-items: center;
  gap: 5px;
  flex-wrap: wrap;
  margin-top: 3px;
}

.card-short-url {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--md-sys-color-primary);
  text-decoration: none;

  &:hover { text-decoration: underline; }
}

.copy-btn {
  width: 22px;
  height: 22px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  color: var(--md-sys-color-on-surface-variant);
  padding: 0;
  transition: background 0.15s;

  .material-symbols-outlined { font-size: 15px; }
  &:hover { background: var(--md-sys-color-surface-container); }
}

.meta-pill {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 1px 6px;
  border-radius: 999px;
  font-size: 0.65rem;
  font-weight: 500;
  background: var(--md-sys-color-surface-container);
  color: var(--md-sys-color-on-surface-variant);

  &--error {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }
  &--warn {
    background: rgba(244, 161, 0, 0.15);
    color: #925F00;
  }
}

.card-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.tag-pill {
  display: inline-flex;
  align-items: center;
  padding: 1px 8px;
  border-radius: 999px;
  font-size: 0.65rem;
  font-weight: 500;
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);

  &--more {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface-variant);
  }
}

/* ── Card Stats ─────────────────────────────────────────────────────────────── */
.card-stats {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
  min-width: 68px;
  padding-top: 4px;
}

.stat-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  line-height: 1;
}

.stat-label {
  font-size: 0.65rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.stat-date {
  font-size: 0.68rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 6px;
  white-space: nowrap;
}

/* ── Card Actions ───────────────────────────────────────────────────────────── */
.card-actions {
  display: flex;
  align-items: flex-start;
  gap: 2px;
  padding-top: 4px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.15s;

  a { text-decoration: none; }
}

.link-card:hover .card-actions,
.link-card:focus-within .card-actions,
.link-card--selected .card-actions {
  opacity: 1;
}

.action-btn {
  width: 34px;
  height: 34px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: var(--md-sys-color-on-surface-variant);
  padding: 0;
  transition: background 0.15s, color 0.15s;

  .material-symbols-outlined { font-size: 18px; }

  &:hover {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
  }

  &--danger:hover {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-error);
  }

  &:disabled { opacity: 0.4; cursor: default; }
}

/* ── Row Menu ───────────────────────────────────────────────────────────────── */
.row-menu-wrap { position: relative; }

.row-menu {
  position: fixed;
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  z-index: 1000;
  min-width: 188px;
  padding: 4px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.15);
}

.row-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 9px 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
  border-radius: 8px;
  text-align: left;
  transition: background 0.12s;
  position: relative;

  &:hover { background: var(--md-sys-color-surface-container); }
  &:disabled { opacity: 0.5; cursor: default; }

  &--active {
    color: var(--md-sys-color-primary);
    font-weight: 500;
  }
}

.row-menu-item__icon {
  font-size: 18px;
  flex-shrink: 0;
  color: var(--md-sys-color-on-surface-variant);
}

.row-menu-item__badge {
  margin-left: auto;
  font-size: 0.65rem;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 999px;
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
}

.row-menu-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  margin: 3px 0;
}

/* ── Pagination ─────────────────────────────────────────────────────────────── */
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  flex-wrap: wrap;
  gap: 12px;
}

.pagination-info {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 2px;
}

.page-btn {
  min-width: 34px;
  height: 34px;
  padding: 0 6px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 8px;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;

  .material-symbols-outlined { font-size: 20px; }

  &:hover { background: var(--md-sys-color-surface-container); }
  &:disabled { opacity: 0.4; cursor: default; }

  &--active {
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    font-weight: 700;

    &:hover { background: var(--md-sys-color-primary); }
  }
}

.page-ellipsis {
  padding: 0 6px;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.875rem;
}

/* ── Folder Manager ─────────────────────────────────────────────────────────── */
.fm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.32);
  z-index: 800;
  display: flex;
  align-items: stretch;
  justify-content: flex-end;
}

.fm-panel {
  background: var(--md-sys-color-surface);
  width: 320px;
  max-width: 90vw;
  height: 100%;
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 24px rgba(0,0,0,0.16);
}

.fm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);
}

.fm-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.fm-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.fm-item {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px 8px 2px 4px;
}

.fm-item-btn {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 8px;
  color: var(--md-sys-color-on-surface-variant);
  text-align: left;
  min-width: 0;
  transition: background 0.15s;

  &:hover { background: var(--md-sys-color-surface-container-low); }
  &--active {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
    font-weight: 600;
  }
}

.fm-item-name {
  flex: 1;
  font-size: 0.875rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fm-item-count {
  font-size: 0.7rem;
  font-weight: 600;
  background: var(--md-sys-color-surface-container);
  color: var(--md-sys-color-on-surface-variant);
  padding: 1px 7px;
  border-radius: 999px;
}

.fm-item-rename {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
}

.fm-input {
  flex: 1;
  height: 38px;
  padding: 0 12px;
  border: 1px solid var(--md-sys-color-outline);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.875rem;
  font-family: 'Roboto', sans-serif;
  outline: none;
  transition: border-color 0.15s;

  &:focus { border-color: var(--md-sys-color-primary); border-width: 2px; }
  &::placeholder { color: var(--md-sys-color-on-surface-variant); }
}

.fm-empty {
  padding: 28px 20px;
  text-align: center;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.5;
}

.fm-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
}

.fm-new-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Panel slide-in transition */
.panel-enter-active,
.panel-leave-active {
  transition: opacity 0.22s ease;
  .fm-panel { transition: transform 0.22s ease; }
}

.panel-enter-from,
.panel-leave-to {
  opacity: 0;
  .fm-panel { transform: translateX(100%); }
}

/* ── Snackbar ───────────────────────────────────────────────────────────────── */
.m3-snackbar {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--md-sys-color-inverse-surface);
  color: var(--md-sys-color-inverse-on-surface);
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 0.875rem;
  z-index: 2000;
  box-shadow: 0 4px 16px rgba(0,0,0,0.2);
  white-space: nowrap;
  min-width: 240px;

  .material-symbols-outlined { color: #1AA563; }
  .btn-text { color: var(--md-sys-color-inverse-primary); }
}

.snack-enter-active,
.snack-leave-active { transition: all 0.25s ease; }
.snack-enter-from,
.snack-leave-to { opacity: 0; transform: translateX(-50%) translateY(12px); }

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 768px) {
  .links-toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  .toolbar-right {
    flex-wrap: wrap;
  }
  .search-wrap {
    width: 100%;
  }
  .link-card {
    grid-template-columns: 36px 40px 1fr auto;
    .card-stats { display: none; }
  }
}
</style>
