<template>
  <div class="app-shell" :class="{ dark: uiStore.darkMode }">

    <!-- Mobile overlay backdrop -->
    <transition name="overlay-fade">
      <div
        v-if="uiStore.sidebarOpen"
        class="nav-overlay"
        @click="uiStore.closeSidebar()"
      />
    </transition>

    <!-- ── Navigation drawer ──────────────────────────────────────────── -->
    <nav
      class="nav-drawer"
      :class="{
        collapsed: uiStore.sidebarCollapsed,
        'mobile-open': uiStore.sidebarOpen,
      }"
    >
      <!-- Logo / header -->
      <div class="nav-header">
        <router-link to="/dashboard" class="nav-logo-link" :title="uiStore.sidebarCollapsed ? 'Linkbee — Home' : undefined">
          <img src="/logo.png" alt="Linkbee" class="nav-logo-img" />
        </router-link>
      </div>

      <md-divider />

      <!-- Scrollable nav body -->
      <div class="nav-body">

        <!-- ── Main section ──────────────────────────────────────────── -->
        <div class="nav-section-label">Main</div>

        <router-link
          to="/dashboard/overview"
          class="nav-item"
          :class="{ active: $route.path === '/dashboard/overview' }"
          :title="uiStore.sidebarCollapsed ? 'Overview' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">home</span>
          <span class="nav-label">Overview</span>
        </router-link>

        <router-link
          to="/dashboard/links"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/links') }"
          :title="uiStore.sidebarCollapsed ? 'Links' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">link</span>
          <span class="nav-label">Links</span>
        </router-link>

        <router-link
          to="/dashboard/analytics"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/analytics') }"
          :title="uiStore.sidebarCollapsed ? 'Analytics' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">bar_chart</span>
          <span class="nav-label">Analytics</span>
        </router-link>

        <router-link
          to="/dashboard/comparison"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/comparison') }"
          :title="uiStore.sidebarCollapsed ? 'Compare' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">compare_arrows</span>
          <span class="nav-label">Compare</span>
        </router-link>

        <router-link
          to="/dashboard/reports"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/reports') }"
          :title="uiStore.sidebarCollapsed ? 'Reports' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">summarize</span>
          <span class="nav-label">Reports</span>
        </router-link>

        <!-- ── Tools section ─────────────────────────────────────────── -->
        <div class="nav-section-label">Tools</div>

        <router-link
          to="/dashboard/teams"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/teams') }"
          :title="uiStore.sidebarCollapsed ? 'Teams' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">group</span>
          <span class="nav-label">Teams</span>
        </router-link>

        <router-link
          to="/dashboard/api-keys"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/api-keys') }"
          :title="uiStore.sidebarCollapsed ? 'API Keys' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">key</span>
          <span class="nav-label">API Keys</span>
        </router-link>

        <router-link
          to="/dashboard/webhooks"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/webhooks') }"
          :title="uiStore.sidebarCollapsed ? 'Webhooks' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">webhook</span>
          <span class="nav-label">Webhooks</span>
        </router-link>

        <router-link
          to="/dashboard/bio"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/bio') }"
          :title="uiStore.sidebarCollapsed ? 'Bio' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">account_circle</span>
          <span class="nav-label">Bio</span>
        </router-link>

        <router-link
          to="/dashboard/domains"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/domains') }"
          :title="uiStore.sidebarCollapsed ? 'Domains' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">domain</span>
          <span class="nav-label">Domains</span>
        </router-link>

        <!-- ── Account section ───────────────────────────────────────── -->
        <div class="nav-section-label">Account</div>

        <router-link
          to="/dashboard/audit-logs"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/audit-logs') }"
          :title="uiStore.sidebarCollapsed ? 'Audit Logs' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">history</span>
          <span class="nav-label">Audit Logs</span>
        </router-link>

        <router-link
          to="/dashboard/billing"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/billing') }"
          :title="uiStore.sidebarCollapsed ? 'Billing' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">credit_card</span>
          <span class="nav-label">Billing</span>
        </router-link>

        <router-link
          to="/dashboard/settings"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/settings') }"
          :title="uiStore.sidebarCollapsed ? 'Settings' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">settings</span>
          <span class="nav-label">Settings</span>
        </router-link>

        <router-link
          to="/dashboard/security"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/security') }"
          :title="uiStore.sidebarCollapsed ? 'Security' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">shield</span>
          <span class="nav-label">Security</span>
        </router-link>

        <!-- Admin (only visible to admin role) -->
        <router-link
          v-if="authStore.isAdmin"
          to="/admin"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/admin') }"
          :title="uiStore.sidebarCollapsed ? 'Admin' : undefined"
          @click="uiStore.closeSidebar()"
        >
          <span class="material-symbols-outlined nav-icon">admin_panel_settings</span>
          <span class="nav-label">Admin</span>
        </router-link>

      </div><!-- /nav-body -->

      <!-- Bottom user card -->
      <div class="nav-user">
        <div class="nav-user-inner">
          <div
            class="nav-user-avatar"
            :title="uiStore.sidebarCollapsed ? authStore.userName : undefined"
          >{{ authStore.userInitials }}</div>
          <div class="nav-user-info">
            <div class="nav-user-name">{{ authStore.userName }}</div>
            <div class="nav-user-email">{{ authStore.profile?.email }}</div>
          </div>
          <button
            class="nav-logout-btn"
            title="Sign out"
            @click="handleLogout"
          >
            <span class="material-symbols-outlined">logout</span>
          </button>
        </div>
      </div>
    </nav><!-- /nav-drawer -->

    <!-- ── Main area ──────────────────────────────────────────────────── -->
    <div
      class="main-area"
      :class="{ collapsed: uiStore.sidebarCollapsed }"
    >

      <!-- Top app bar -->
      <header class="top-app-bar">
        <!-- Hamburger / collapse toggle -->
        <button
          class="icon-btn sidebar-toggle"
          title="Toggle sidebar"
          @click="uiStore.toggleSidebar()"
        >
          <span class="material-symbols-outlined">menu</span>
        </button>

        <!-- Page title -->
        <span class="page-title">{{ pageTitle }}</span>

        <!-- Dark mode toggle -->
        <button
          class="icon-btn"
          :title="uiStore.darkMode ? 'Switch to light mode' : 'Switch to dark mode'"
          @click="uiStore.toggleDarkMode()"
        >
          <span class="material-symbols-outlined">
            {{ uiStore.darkMode ? 'light_mode' : 'dark_mode' }}
          </span>
        </button>

        <!-- New Link shortcut -->
        <button class="topbar-new-link-btn" @click="handleNewLink">
          <span class="material-symbols-outlined">add</span>
          New Link
        </button>

        <!-- User avatar dropdown -->
        <div class="user-menu-wrapper" ref="userMenuWrapperRef">
          <button
            class="user-avatar-btn"
            :title="authStore.userName"
            :aria-expanded="userDropdownOpen"
            @click="userDropdownOpen = !userDropdownOpen"
          >
            {{ authStore.userInitials }}
          </button>

          <transition name="dropdown">
            <div
              v-if="userDropdownOpen"
              class="user-dropdown"
              role="menu"
            >
              <!-- User info header -->
              <div class="user-dropdown-header">
                <div class="user-dropdown-avatar">{{ authStore.userInitials }}</div>
                <div class="user-dropdown-info">
                  <div class="user-dropdown-name">{{ authStore.userName }}</div>
                  <div class="user-dropdown-email">{{ authStore.profile?.email }}</div>
                </div>
              </div>

              <div class="user-dropdown-divider" />

              <button
                class="user-dropdown-item"
                role="menuitem"
                @click="goToProfile"
              >
                <span class="material-symbols-outlined">person</span>
                Profile
              </button>

              <button
                class="user-dropdown-item"
                role="menuitem"
                @click="goToSettings"
              >
                <span class="material-symbols-outlined">settings</span>
                Settings
              </button>

              <div class="user-dropdown-divider" />

              <button
                class="user-dropdown-item user-dropdown-item--danger"
                role="menuitem"
                @click="handleLogout"
              >
                <span class="material-symbols-outlined">logout</span>
                Sign out
              </button>
            </div>
          </transition>
        </div>
      </header>

      <!-- Page content -->
      <main class="page-content">
        <router-view />
      </main>
    </div><!-- /main-area -->

  </div><!-- /app-shell -->
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useUiStore } from '@/stores/ui';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const uiStore = useUiStore();

const userDropdownOpen = ref(false);
const userMenuWrapperRef = ref<HTMLElement | null>(null);

const pageTitle = computed(() => {
  const name = route.name as string | undefined;
  if (!name) return 'Dashboard';
  const titles: Record<string, string> = {
    overview: 'Overview',
    links: 'My Links',
    'link-analytics': 'Link Analytics',
    'global-analytics': 'Analytics',
    'link-comparison': 'Compare Links',
    reports: 'Scheduled Reports',
    'api-keys': 'API Keys',
    webhooks: 'Webhooks',
    domains: 'Custom Domains',
    'audit-logs': 'Audit Logs',
    billing: 'Billing & Plan',
    security: 'Security',
    settings: 'Settings',
    bio: 'Link-in-Bio',
    teams: 'Teams',
    admin: 'Admin',
  };
  return titles[name] ?? 'Dashboard';
});

function handleNewLink() {
  uiStore.triggerCreateLink();
  if (route.name !== 'links') {
    router.push('/dashboard/links');
  }
}

function goToProfile() {
  userDropdownOpen.value = false;
  router.push('/dashboard/settings');
}

function goToSettings() {
  userDropdownOpen.value = false;
  router.push('/dashboard/settings');
}

async function handleLogout() {
  userDropdownOpen.value = false;
  await authStore.logout();
  router.push('/login');
}

function handleClickOutside(e: MouseEvent) {
  // Close user dropdown when clicking outside
  if (
    userDropdownOpen.value &&
    userMenuWrapperRef.value &&
    !userMenuWrapperRef.value.contains(e.target as Node)
  ) {
    userDropdownOpen.value = false;
  }

  // Close mobile sidebar when clicking outside drawer & toggle
  const target = e.target as HTMLElement;
  if (
    uiStore.sidebarOpen &&
    !target.closest('.nav-drawer') &&
    !target.closest('.sidebar-toggle')
  ) {
    uiStore.closeSidebar();
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped lang="scss">
// ── Shared easing ─────────────────────────────────────────────────────────
$ease-standard: cubic-bezier(0.2, 0, 0, 1);
$drawer-width: 256px;
$rail-width: 80px;
$appbar-height: 64px;

// ── App shell ──────────────────────────────────────────────────────────────
.app-shell {
  display: flex;
  min-height: 100vh;
  background: var(--md-sys-color-background);
  position: relative;
}

// ── Navigation drawer ──────────────────────────────────────────────────────
.nav-drawer {
  width: $drawer-width;
  min-height: 100vh;
  background: var(--md-sys-color-surface-container-low);
  border-right: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 200;
  transition:
    width 0.25s $ease-standard,
    transform 0.25s $ease-standard;
  overflow: hidden;
}

.nav-drawer.collapsed {
  width: $rail-width;
}

// Logo / header
.nav-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
  height: $appbar-height;
  flex-shrink: 0;
}

.nav-logo-link {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  text-decoration: none;
  border-radius: 8px;
  transition: opacity 0.15s;

  &:hover { opacity: 0.8; }
  &:focus-visible { outline: 2px solid var(--md-sys-color-primary); outline-offset: 3px; }
}

.nav-logo-img {
  height: 36px;
  width: auto;
  max-width: 140px;
  object-fit: contain;
  display: block;
  transition: max-width 0.25s $ease-standard, opacity 0.15s $ease-standard;
}

// Collapsed: shrink logo to icon-only width
.nav-drawer.collapsed .nav-logo-img {
  max-width: 36px;
  object-fit: cover;
  object-position: left center;
}

// ── Scrollable nav body ────────────────────────────────────────────────────
.nav-body {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px 0;

  // Thin scrollbar
  scrollbar-width: thin;
  scrollbar-color: var(--md-sys-color-outline-variant) transparent;

  &::-webkit-scrollbar {
    width: 4px;
  }
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  &::-webkit-scrollbar-thumb {
    background: var(--md-sys-color-outline-variant);
    border-radius: 4px;
  }
}

// Section labels
.nav-section-label {
  padding: 12px 28px 4px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.8px;
  text-transform: uppercase;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  overflow: hidden;
  opacity: 1;
  transition: opacity 0.15s $ease-standard;
}

.nav-drawer.collapsed .nav-section-label {
  opacity: 0;
  height: 0;
  padding: 0;
  overflow: hidden;
}

// Nav items — M3 Navigation Drawer Item
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px;
  height: 52px;
  border-radius: 100px;
  margin: 2px 12px;
  text-decoration: none;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.1px;
  cursor: pointer;
  transition:
    background 0.15s $ease-standard,
    color 0.15s $ease-standard;
  white-space: nowrap;
  overflow: hidden;
  position: relative;

  // State layer
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    background: var(--md-sys-color-on-surface);
    opacity: 0;
    transition: opacity 0.15s $ease-standard;
    pointer-events: none;
  }

  &:hover::before {
    opacity: 0.08;
  }

  &:active::before {
    opacity: 0.12;
  }
}

.nav-item.router-link-active,
.nav-item.active {
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);

  &::before {
    background: var(--md-sys-color-on-secondary-container);
  }

  .nav-icon {
    font-variation-settings: 'FILL' 1;
  }
}

.nav-icon {
  font-size: 22px;
  flex-shrink: 0;
  font-variation-settings: 'FILL' 0;
  transition: font-variation-settings 0.15s $ease-standard;
}

.nav-label {
  white-space: nowrap;
  overflow: hidden;
  transition: opacity 0.15s $ease-standard;
}

// Collapsed rail: center icons, hide labels
.nav-drawer.collapsed .nav-item {
  justify-content: center;
  padding: 0;
  margin: 2px 8px;
}

.nav-drawer.collapsed .nav-label {
  opacity: 0;
  width: 0;
  overflow: hidden;
}

// ── Bottom user card ───────────────────────────────────────────────────────
.nav-user {
  flex-shrink: 0;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  padding: 12px;
}

.nav-user-inner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 8px;
  border-radius: 12px;
  overflow: hidden;
  transition: background 0.15s $ease-standard;

  &:hover {
    background: color-mix(in srgb, var(--md-sys-color-on-surface) 6%, transparent);
  }
}

.nav-user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  flex-shrink: 0;
  user-select: none;
}

.nav-user-info {
  flex: 1;
  overflow: hidden;
  min-width: 0;
  transition: opacity 0.15s $ease-standard;
}

.nav-user-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-user-email {
  font-size: 11px;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 1px;
}

.nav-logout-btn {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: var(--md-sys-color-on-surface-variant);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.15s $ease-standard, color 0.15s $ease-standard;

  .material-symbols-outlined {
    font-size: 18px;
  }

  &:hover {
    background: color-mix(in srgb, var(--md-sys-color-error) 12%, transparent);
    color: var(--md-sys-color-error);
  }
}

// Collapsed state for user section
.nav-drawer.collapsed .nav-user-info,
.nav-drawer.collapsed .nav-logout-btn {
  opacity: 0;
  width: 0;
  overflow: hidden;
  padding: 0;
  margin: 0;
}

.nav-drawer.collapsed .nav-user-inner {
  justify-content: center;
  padding: 4px;
}

// ── Mobile overlay ─────────────────────────────────────────────────────────
.nav-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 199;
  backdrop-filter: blur(2px);
}

.overlay-fade-enter-active,
.overlay-fade-leave-active {
  transition: opacity 0.25s $ease-standard;
}
.overlay-fade-enter-from,
.overlay-fade-leave-to {
  opacity: 0;
}

// ── Main area ──────────────────────────────────────────────────────────────
.main-area {
  margin-left: $drawer-width;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--md-sys-color-background);
  flex: 1;
  min-width: 0;
  transition: margin-left 0.25s $ease-standard;
}

.main-area.collapsed {
  margin-left: $rail-width;
}

// ── Top app bar ────────────────────────────────────────────────────────────
.top-app-bar {
  height: $appbar-height;
  background: var(--md-sys-color-surface);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  align-items: center;
  padding: 0 16px 0 8px;
  gap: 4px;
  position: sticky;
  top: 0;
  z-index: 100;
  flex-shrink: 0;
}

.page-title {
  font-size: 20px;
  font-weight: 500;
  flex: 1;
  color: var(--md-sys-color-on-surface);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-left: 4px;
}

// Generic icon button
.icon-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: var(--md-sys-color-on-surface-variant);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s $ease-standard, color 0.15s $ease-standard;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    background: var(--md-sys-color-on-surface);
    opacity: 0;
    transition: opacity 0.15s $ease-standard;
    pointer-events: none;
  }

  &:hover::before {
    opacity: 0.08;
  }

  &:active::before {
    opacity: 0.12;
  }

  .material-symbols-outlined {
    font-size: 22px;
  }
}

// ── Topbar New Link button ─────────────────────────────────────────────────
.topbar-new-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 36px;
  padding: 0 16px;
  border: none;
  border-radius: 8px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
  transition: opacity 0.15s $ease-standard;

  .material-symbols-outlined { font-size: 18px; }

  &:hover { opacity: 0.88; }
}

// ── User avatar button ────────────────────────────────────────────────────
.user-menu-wrapper {
  position: relative;
  flex-shrink: 0;
}

.user-avatar-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 2px solid var(--md-sys-color-primary-container);
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  cursor: pointer;
  user-select: none;
  transition:
    border-color 0.15s $ease-standard,
    box-shadow 0.15s $ease-standard;

  &:hover {
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 20%, transparent);
  }
}

// ── User dropdown ──────────────────────────────────────────────────────────
.user-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  width: 260px;
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 16px;
  box-shadow:
    0 4px 8px 3px rgba(0, 0, 0, 0.08),
    0 1px 3px rgba(0, 0, 0, 0.12);
  z-index: 300;
  overflow: hidden;
  padding: 8px 0;
}

.user-dropdown-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px 14px;
}

.user-dropdown-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  flex-shrink: 0;
}

.user-dropdown-info {
  overflow: hidden;
  min-width: 0;
}

.user-dropdown-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-dropdown-email {
  font-size: 12px;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}

.user-dropdown-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  margin: 4px 0;
}

.user-dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 16px;
  border: none;
  background: transparent;
  color: var(--md-sys-color-on-surface);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  text-align: left;
  transition: background 0.12s $ease-standard;

  .material-symbols-outlined {
    font-size: 20px;
    color: var(--md-sys-color-on-surface-variant);
    flex-shrink: 0;
  }

  &:hover {
    background: color-mix(in srgb, var(--md-sys-color-on-surface) 8%, transparent);
  }
}

.user-dropdown-item--danger {
  color: var(--md-sys-color-error);

  .material-symbols-outlined {
    color: var(--md-sys-color-error);
  }

  &:hover {
    background: color-mix(in srgb, var(--md-sys-color-error) 10%, transparent);
  }
}

// Dropdown transition
.dropdown-enter-active,
.dropdown-leave-active {
  transition:
    opacity 0.15s $ease-standard,
    transform 0.15s $ease-standard;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.97);
}

// ── Page content ───────────────────────────────────────────────────────────
.page-content {
  flex: 1;
  padding: 24px;
  overflow-x: hidden;

  @media (max-width: 575px) {
    padding: 16px;
  }
}

// ── Mobile responsive ──────────────────────────────────────────────────────
@media (max-width: 768px) {
  .nav-drawer {
    transform: translateX(-100%);
    width: $drawer-width;
    // Override collapsed on mobile — always full width when open
    &.collapsed {
      width: $drawer-width;

      .nav-section-label {
        opacity: 1;
        height: auto;
        padding: 12px 28px 4px;
      }

      .nav-label {
        opacity: 1;
        width: auto;
        overflow: visible;
      }

      .nav-item {
        justify-content: flex-start;
        padding: 0 16px;
        margin: 2px 12px;
      }

      .nav-logo-text {
        opacity: 1;
        width: auto;
        pointer-events: auto;
      }

      .nav-user-info,
      .nav-logout-btn {
        opacity: 1;
        width: auto;
        overflow: visible;
        padding: initial;
        margin: initial;
      }

      .nav-user-inner {
        justify-content: flex-start;
        padding: 6px 8px;
      }
    }
  }

  .nav-drawer.mobile-open {
    transform: translateX(0);
  }

  .main-area,
  .main-area.collapsed {
    margin-left: 0;
  }
}
</style>
