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

      <!-- ── Left: Page Settings ─────────────────────────────────────────── -->
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

      <!-- ── Right: Links Management ──────────────────────────────────────── -->
      <div class="links-column">

        <!-- Add Link Card -->
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">add_link</span>
            </div>
            <span class="an-card-title">Add Link</span>
          </div>
          <div class="an-card-body">

            <!-- ── Social Platforms ── -->
            <div class="section-label">
              <span class="material-symbols-outlined" style="font-size:15px;">share</span>
              Social Platforms
            </div>

            <div class="social-grid">
              <button
                v-for="p in SOCIAL_PLATFORMS"
                :key="p.id"
                :class="['social-btn', activePlatform?.id === p.id ? 'social-btn--active' : '']"
                :style="activePlatform?.id === p.id
                  ? { background: p.bg, color: p.fg, borderColor: p.bg }
                  : {}"
                @click="selectPlatform(p)"
                :title="p.name"
                type="button"
              >
                <svg class="social-btn__icon" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                  <path :d="p.icon" />
                </svg>
                <span class="social-btn__name">{{ p.name }}</span>
              </button>
            </div>

            <!-- Platform Input Panel -->
            <Transition name="panel-slide">
              <div v-if="activePlatform" class="platform-panel">
                <div class="platform-panel__header">
                  <span class="platform-badge" :style="{ background: activePlatform.bg, color: activePlatform.fg }">
                    <svg viewBox="0 0 24 24" fill="currentColor" style="width:13px;height:13px;flex-shrink:0;">
                      <path :d="activePlatform.icon" />
                    </svg>
                    {{ activePlatform.name }}
                  </span>
                  <button class="btn-icon btn-sm" type="button" @click="closePlatform" title="Cancel">
                    <span class="material-symbols-outlined">close</span>
                  </button>
                </div>

                <div class="platform-panel__inputs">
                  <!-- Username-based: show URL prefix + username input -->
                  <template v-if="activePlatform.inputType === 'username'">
                    <div class="prefixed-input" style="flex:1;min-width:0;">
                      <span class="prefixed-input__prefix">{{ activePlatform.urlBase }}</span>
                      <input
                        ref="platformInputRef"
                        class="form-input prefixed-input__field"
                        :value="platformUsername"
                        @input="platformUsername = ($event.target as HTMLInputElement).value"
                        :placeholder="activePlatform.placeholder"
                        @keydown.enter="addSocialLink"
                      />
                    </div>
                  </template>
                  <!-- URL-based: full URL input -->
                  <template v-else>
                    <input
                      ref="platformInputRef"
                      class="form-input"
                      type="url"
                      :value="platformUsername"
                      @input="platformUsername = ($event.target as HTMLInputElement).value"
                      :placeholder="activePlatform.placeholder"
                      @keydown.enter="addSocialLink"
                      style="flex:1;min-width:0;"
                    />
                  </template>

                  <button
                    class="btn-filled"
                    type="button"
                    :disabled="addingLink || !platformUsername.trim()"
                    @click="addSocialLink"
                    style="flex-shrink:0;"
                  >
                    <span v-if="addingLink" class="css-spinner css-spinner--sm css-spinner--white"></span>
                    <span v-else class="material-symbols-outlined">add</span>
                    Add
                  </button>
                </div>

                <div v-if="platformError" class="error-box" style="margin-top:0;">{{ platformError }}</div>
              </div>
            </Transition>

            <!-- ── Divider ── -->
            <div class="section-divider">
              <span>or add a custom link</span>
            </div>

            <!-- ── Custom Link Form ── -->
            <div class="section-label">
              <span class="material-symbols-outlined" style="font-size:15px;">link</span>
              Custom Link
            </div>

            <div class="add-link-row">
              <input
                class="form-input"
                :value="newTitle"
                @input="newTitle = ($event.target as HTMLInputElement).value"
                placeholder="Link title"
                maxlength="100"
                style="flex:1;min-width:0;"
                @keydown.enter="addLink"
              />
              <input
                class="form-input"
                type="url"
                :value="newUrl"
                @input="newUrl = ($event.target as HTMLInputElement).value"
                placeholder="https://..."
                style="flex:1;min-width:0;"
                @keydown.enter="addLink"
              />
              <button
                class="btn-filled"
                type="button"
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
              :style="{ opacity: link.is_active ? '1' : '0.55' }"
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

                <!-- Platform icon bubble or status dot -->
                <span
                  v-if="getLinkPlatform(link.url)"
                  class="link-platform-icon"
                  :style="{ background: getLinkPlatform(link.url)!.bg }"
                  :title="getLinkPlatform(link.url)!.name"
                >
                  <svg viewBox="0 0 24 24" fill="currentColor" :style="{ color: getLinkPlatform(link.url)!.fg }">
                    <path :d="getLinkPlatform(link.url)!.icon" />
                  </svg>
                </span>
                <span
                  v-else
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
                <button class="btn-icon btn-sm" title="Edit" @click="startEdit(link)">
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
                  @input="editTitle = ($event.target as HTMLInputElement).value"
                  placeholder="Title"
                  maxlength="100"
                  style="flex:1;min-width:120px;"
                />
                <input
                  class="form-input"
                  type="url"
                  :value="editUrl"
                  @input="editUrl = ($event.target as HTMLInputElement).value"
                  placeholder="URL"
                  style="flex:1;min-width:120px;"
                />
                <div class="link-edit-actions">
                  <button class="btn-icon btn-sm" :disabled="savingEdit" @click="saveEdit(link)">
                    <span v-if="savingEdit" class="css-spinner css-spinner--sm"></span>
                    <span v-else class="material-symbols-outlined">check</span>
                  </button>
                  <button class="btn-icon btn-sm" @click="cancelEdit">
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
import { ref, computed, onMounted, nextTick } from 'vue';
import bioApi from '@/api/bio';
import type { BioPage, BioLinkItem } from '@/types/bio';

// ── Social Platforms ─────────────────────────────────────────────────────────

interface SocialPlatform {
  id: string;
  name: string;
  bg: string;          // brand background color
  fg: string;          // icon/text color on brand bg
  domains: string[];   // URL substrings for detection
  urlBase: string;     // URL prefix for username-based platforms
  inputType: 'username' | 'url';
  placeholder: string;
  icon: string;        // SVG path (viewBox 0 0 24 24)
}

const SOCIAL_PLATFORMS: SocialPlatform[] = [
  {
    id: 'instagram',
    name: 'Instagram',
    bg: '#E1306C',
    fg: '#fff',
    domains: ['instagram.com'],
    urlBase: 'https://instagram.com/',
    inputType: 'username',
    placeholder: 'your_username',
    icon: 'M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z',
  },
  {
    id: 'x',
    name: 'Twitter / X',
    bg: '#000',
    fg: '#fff',
    domains: ['x.com', 'twitter.com'],
    urlBase: 'https://x.com/',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.741l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z',
  },
  {
    id: 'youtube',
    name: 'YouTube',
    bg: '#FF0000',
    fg: '#fff',
    domains: ['youtube.com', 'youtu.be'],
    urlBase: 'https://youtube.com/@',
    inputType: 'username',
    placeholder: 'yourchannel',
    icon: 'M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z',
  },
  {
    id: 'tiktok',
    name: 'TikTok',
    bg: '#010101',
    fg: '#fff',
    domains: ['tiktok.com'],
    urlBase: 'https://tiktok.com/@',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z',
  },
  {
    id: 'linkedin',
    name: 'LinkedIn',
    bg: '#0A66C2',
    fg: '#fff',
    domains: ['linkedin.com'],
    urlBase: 'https://linkedin.com/in/',
    inputType: 'username',
    placeholder: 'yourprofile',
    icon: 'M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z',
  },
  {
    id: 'facebook',
    name: 'Facebook',
    bg: '#1877F2',
    fg: '#fff',
    domains: ['facebook.com', 'fb.com'],
    urlBase: 'https://facebook.com/',
    inputType: 'username',
    placeholder: 'yourpage',
    icon: 'M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z',
  },
  {
    id: 'github',
    name: 'GitHub',
    bg: '#181717',
    fg: '#fff',
    domains: ['github.com'],
    urlBase: 'https://github.com/',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12',
  },
  {
    id: 'discord',
    name: 'Discord',
    bg: '#5865F2',
    fg: '#fff',
    domains: ['discord.gg', 'discord.com'],
    urlBase: 'https://discord.gg/',
    inputType: 'username',
    placeholder: 'invite-code',
    icon: 'M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057.1 18.082.118 18.1.135 18.1a19.927 19.927 0 0 0 5.993 3.031.074.074 0 0 0 .084-.026 13.918 13.918 0 0 0 1.19-1.938.07.07 0 0 0-.038-.098 13.21 13.21 0 0 1-1.885-.9.07.07 0 0 1-.007-.115c.127-.095.254-.194.375-.293a.072.072 0 0 1 .075-.01c3.927 1.793 8.18 1.793 12.062 0a.072.072 0 0 1 .076.01c.122.1.249.198.376.293a.07.07 0 0 1-.006.115 12.36 12.36 0 0 1-1.886.9.07.07 0 0 0-.038.097c.36.698.772 1.362 1.19 1.94a.074.074 0 0 0 .084.026 19.843 19.843 0 0 0 6.002-3.03.077.077 0 0 0 .032-.028c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03z',
  },
  {
    id: 'twitch',
    name: 'Twitch',
    bg: '#9146FF',
    fg: '#fff',
    domains: ['twitch.tv'],
    urlBase: 'https://twitch.tv/',
    inputType: 'username',
    placeholder: 'channel',
    icon: 'M11.571 4.714h1.715v5.143H11.57zm4.715 0H18v5.143h-1.714zM6 0L1.714 4.286v15.428h5.143V24l4.286-4.286h3.428L22.286 12V0zm14.571 11.143l-3.428 3.428h-3.429l-3 3v-3H6.857V1.714h13.714z',
  },
  {
    id: 'reddit',
    name: 'Reddit',
    bg: '#FF4500',
    fg: '#fff',
    domains: ['reddit.com'],
    urlBase: 'https://reddit.com/u/',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M12 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0zm5.01 4.744c.688 0 1.25.561 1.25 1.249a1.25 1.25 0 0 1-2.498.056l-2.597-.547-.8 3.747c1.824.07 3.48.632 4.674 1.488.308-.309.73-.491 1.207-.491.968 0 1.754.786 1.754 1.754 0 .716-.435 1.333-1.01 1.614a3.111 3.111 0 0 1 .042.52c0 2.694-3.13 4.87-7.004 4.87-3.874 0-7.004-2.176-7.004-4.87 0-.183.015-.366.043-.534A1.748 1.748 0 0 1 4.028 12c0-.968.786-1.754 1.754-1.754.463 0 .898.196 1.207.49 1.207-.883 2.878-1.43 4.744-1.487l.885-4.182a.342.342 0 0 1 .14-.197.35.35 0 0 1 .238-.042l2.906.617a1.214 1.214 0 0 1 1.108-.701zM9.25 12C8.561 12 8 12.562 8 13.25c0 .687.561 1.248 1.25 1.248.687 0 1.248-.561 1.248-1.249 0-.688-.561-1.249-1.249-1.249zm5.5 0c-.687 0-1.248.561-1.248 1.25 0 .687.561 1.248 1.249 1.248.688 0 1.249-.561 1.249-1.249 0-.687-.562-1.249-1.25-1.249zm-5.466 3.99a.327.327 0 0 0-.231.094.33.33 0 0 0 0 .463c.842.842 2.484.913 2.961.913.477 0 2.105-.056 2.961-.913a.361.361 0 0 0 .029-.463.33.33 0 0 0-.464 0c-.547.533-1.684.73-2.512.73-.828 0-1.979-.196-2.512-.73a.326.326 0 0 0-.232-.095z',
  },
  {
    id: 'spotify',
    name: 'Spotify',
    bg: '#1DB954',
    fg: '#fff',
    domains: ['spotify.com', 'open.spotify.com'],
    urlBase: 'https://open.spotify.com/',
    inputType: 'url',
    placeholder: 'https://open.spotify.com/artist/...',
    icon: 'M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z',
  },
  {
    id: 'telegram',
    name: 'Telegram',
    bg: '#26A5E4',
    fg: '#fff',
    domains: ['t.me', 'telegram.me'],
    urlBase: 'https://t.me/',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z',
  },
  {
    id: 'whatsapp',
    name: 'WhatsApp',
    bg: '#25D366',
    fg: '#fff',
    domains: ['wa.me', 'whatsapp.com'],
    urlBase: 'https://wa.me/',
    inputType: 'url',
    placeholder: 'https://wa.me/1234567890',
    icon: 'M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413z',
  },
  {
    id: 'pinterest',
    name: 'Pinterest',
    bg: '#E60023',
    fg: '#fff',
    domains: ['pinterest.com'],
    urlBase: 'https://pinterest.com/',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M12 0C5.373 0 0 5.372 0 12c0 5.084 3.163 9.426 7.627 11.174-.105-.949-.2-2.405.042-3.441.218-.937 1.407-5.965 1.407-5.965s-.359-.719-.359-1.782c0-1.668.967-2.914 2.171-2.914 1.023 0 1.518.769 1.518 1.69 0 1.029-.655 2.568-.994 3.995-.283 1.194.599 2.169 1.777 2.169 2.133 0 3.772-2.249 3.772-5.495 0-2.873-2.064-4.882-5.012-4.882-3.414 0-5.418 2.561-5.418 5.207 0 1.031.397 2.138.893 2.738.098.119.112.224.083.345l-.333 1.36c-.053.22-.174.267-.402.161-1.499-.698-2.436-2.889-2.436-4.649 0-3.785 2.75-7.262 7.929-7.262 4.163 0 7.398 2.967 7.398 6.931 0 4.136-2.607 7.464-6.227 7.464-1.216 0-2.359-.632-2.75-1.378l-.748 2.853c-.271 1.043-1.002 2.35-1.492 3.146C9.57 23.812 10.763 24 12 24c6.627 0 12-5.373 12-12S18.627 0 12 0z',
  },
  {
    id: 'threads',
    name: 'Threads',
    bg: '#000',
    fg: '#fff',
    domains: ['threads.net'],
    urlBase: 'https://threads.net/@',
    inputType: 'username',
    placeholder: 'username',
    icon: 'M12.186 24h-.007c-3.581-.024-6.334-1.205-8.184-3.509C2.35 18.44 1.5 15.586 1.5 12.068V12c.002-3.51.857-6.354 2.544-8.418C5.893 1.204 8.638.022 12.215 0h.013c2.819.019 5.125.837 6.856 2.432 1.678 1.545 2.666 3.72 2.939 6.468l.002.018h-3.467l-.002-.013c-.488-3.786-2.531-5.71-6.294-5.715h-.009c-4.455.003-6.784 2.842-6.784 7.809v.002c0 4.97 2.329 7.808 6.784 7.811h.007c1.873 0 3.183-.414 4.164-1.3.91-.823 1.52-2.044 1.785-3.527H12.71v-3.2h7.946v.001l.045.538c.055.671.061 1.327.017 1.968-.21 3.038-1.29 5.393-3.21 7.003C15.839 23.213 14.185 24 12.186 24z',
  },
  {
    id: 'bluesky',
    name: 'Bluesky',
    bg: '#0085FF',
    fg: '#fff',
    domains: ['bsky.app', 'bsky.social'],
    urlBase: 'https://bsky.app/profile/',
    inputType: 'username',
    placeholder: 'username.bsky.social',
    icon: 'M12 10.8c-1.087-2.114-4.046-6.053-6.798-7.995C2.566.944 1.561 1.266.902 1.565.139 1.908 0 3.08 0 3.768c0 .69.378 5.65.624 6.479.815 2.736 3.713 3.66 6.383 3.364.136-.02.275-.039.415-.056-.138.022-.276.04-.415.056-3.912.58-7.387 2.005-2.83 7.078 5.013 5.19 6.87-1.113 7.823-4.308.953 3.195 2.05 9.271 7.733 4.308 4.267-4.308 1.172-6.498-2.74-7.078a8.741 8.741 0 0 1-.415-.056c.14.017.279.036.415.056 2.67.297 5.568-.628 6.383-3.364.246-.828.624-5.79.624-6.478 0-.69-.139-1.861-.902-2.206-.659-.298-1.664-.62-4.3 1.24C16.046 4.748 13.087 8.687 12 10.8z',
  },
];

// ── State ────────────────────────────────────────────────────────────────────

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

// Social platform
const activePlatform = ref<SocialPlatform | null>(null);
const platformUsername = ref('');
const platformError = ref('');
const platformInputRef = ref<HTMLInputElement | null>(null);

// Add custom link
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

// ── Computed ─────────────────────────────────────────────────────────────────

const publicUrl = computed(() =>
  bioPage.value?.username
    ? `${window.location.origin}/bio/${bioPage.value.username}`
    : ''
);

const bioStats = computed(() => ({
  total: bioPage.value?.links.length ?? 0,
  active: bioPage.value?.links.filter(l => l.is_active).length ?? 0,
  inactive: bioPage.value?.links.filter(l => !l.is_active).length ?? 0,
}));

const bioCopied = ref(false);

// ── Platform helpers ──────────────────────────────────────────────────────────

function getLinkPlatform(url: string): SocialPlatform | null {
  const lower = url.toLowerCase();
  return SOCIAL_PLATFORMS.find(p => p.domains.some(d => lower.includes(d))) ?? null;
}

async function selectPlatform(p: SocialPlatform) {
  if (activePlatform.value?.id === p.id) {
    activePlatform.value = null;
    platformUsername.value = '';
    platformError.value = '';
    return;
  }
  activePlatform.value = p;
  platformUsername.value = '';
  platformError.value = '';
  addError.value = '';
  await nextTick();
  platformInputRef.value?.focus();
}

function closePlatform() {
  activePlatform.value = null;
  platformUsername.value = '';
  platformError.value = '';
}

async function addSocialLink() {
  if (!activePlatform.value || !platformUsername.value.trim()) return;
  addingLink.value = true;
  platformError.value = '';
  const p = activePlatform.value;
  const raw = platformUsername.value.trim().replace(/^@/, '');
  const url = p.inputType === 'username' ? p.urlBase + raw : raw;
  try {
    const res = await bioApi.createLink({ title: p.name, url });
    if (res.data) {
      links.value.push(res.data);
      platformUsername.value = '';
      activePlatform.value = null;
    }
  } catch (err: any) {
    platformError.value = err?.response?.data?.description ?? 'Failed to add link.';
  } finally {
    addingLink.value = false;
  }
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────

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

// ── Actions ───────────────────────────────────────────────────────────────────

async function copyBioUrl() {
  if (!publicUrl.value) return;
  await navigator.clipboard.writeText(publicUrl.value);
  bioCopied.value = true;
  setTimeout(() => (bioCopied.value = false), 2000);
}

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
    // silent
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

// ── Drag & Drop ───────────────────────────────────────────────────────────────

function onDragStart(idx: number) { dragFrom = idx; }

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
  bioApi.reorderLinks({ ids: links.value.map(l => l.id) }).catch(() => {});
}

function onDragEnd() { dragFrom = -1; }
</script>

<style scoped lang="scss">
/* ── Root ─────────────────────────────────────────────────────────────────── */
.bio-page {
  max-width: 1000px;
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
  gap: 14px;
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

  &--sm { width: 16px; height: 16px; border-width: 2px; }
  &--white { border-color: rgba(255,255,255,.35); border-top-color: #fff; }
}

@keyframes spin { to { transform: rotate(360deg); } }

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
  transition: border-color .15s, box-shadow .15s;
  box-sizing: border-box;
  width: 100%;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, .12);
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
  transition: border-color .15s, box-shadow .15s;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, .12);
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
  transition: border-color .15s, box-shadow .15s;

  &:focus-within {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, .12);
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
  width: auto;

  &:focus { box-shadow: none; outline: none; }
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
    background: rgba(22, 163, 74, .06);
    border: 1px solid rgba(22, 163, 74, .25);
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
  grid-template-columns: 1fr 1.5fr;
  gap: 20px;
  align-items: start;

  @media (max-width: 800px) {
    grid-template-columns: 1fr;
  }
}

.links-column {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ── Section label ────────────────────────────────────────────────────────── */
.section-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.07em;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Social Platform Grid ─────────────────────────────────────────────────── */
.social-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;

  @media (max-width: 500px) {
    grid-template-columns: repeat(3, 1fr);
  }
}

.social-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  padding: 10px 6px 8px;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 11px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface-variant);
  cursor: pointer;
  transition: border-color .15s, background .15s, color .15s, transform .1s;
  min-width: 0;

  &:hover:not(.social-btn--active) {
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface);
    transform: translateY(-1px);
  }

  &--active {
    box-shadow: 0 2px 8px rgba(0,0,0,.18);
    transform: translateY(-1px);
  }

  &__icon {
    width: 20px;
    height: 20px;
    flex-shrink: 0;
  }

  &__name {
    font-size: 0.62rem;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
    line-height: 1.2;
  }
}

/* ── Platform Input Panel ─────────────────────────────────────────────────── */
.platform-panel {
  background: var(--md-sys-color-surface-container-low);
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 12px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;

  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  &__inputs {
    display: flex;
    gap: 8px;
    align-items: center;
  }
}

.platform-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px 4px 8px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.01em;
}

/* ── Prefixed input (URL base + username) ─────────────────────────────────── */
.prefixed-input {
  display: flex;
  align-items: stretch;
  border: 1.5px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  overflow: hidden;
  transition: border-color .15s, box-shadow .15s;
  background: var(--md-sys-color-surface);

  &:focus-within {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px rgba(99, 91, 255, .12);
  }

  &__prefix {
    padding: 0 8px;
    display: flex;
    align-items: center;
    background: var(--md-sys-color-surface-container);
    border-right: 1.5px solid var(--md-sys-color-outline-variant);
    font-size: 0.75rem;
    color: var(--md-sys-color-on-surface-variant);
    white-space: nowrap;
    flex-shrink: 0;
    max-width: 140px;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &__field {
    flex: 1;
    border: none;
    border-radius: 0;
    background: transparent;
    width: auto;

    &:focus { box-shadow: none; outline: none; }
  }
}

/* ── Section divider ──────────────────────────────────────────────────────── */
.section-divider {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.78rem;

  &::before, &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--md-sys-color-outline-variant);
  }
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
    opacity: .6;
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

  &--success { background: rgba(22, 163, 74, .12); color: #16a34a; }
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
  transition: background .12s;

  &:last-child { border-bottom: none; }
  &:hover { background: var(--md-sys-color-surface-container-low); }
}

.link-item__view {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  width: 100%;
  box-sizing: border-box;
}

.link-item__edit {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
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

/* Platform icon bubble in link list */
.link-platform-icon {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  svg { width: 13px; height: 13px; }
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
    font-size: 0.72rem;
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

/* ── Transition: platform panel slide ─────────────────────────────────────── */
.panel-slide-enter-active,
.panel-slide-leave-active {
  transition: opacity .2s ease, transform .2s ease;
}

.panel-slide-enter-from,
.panel-slide-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
