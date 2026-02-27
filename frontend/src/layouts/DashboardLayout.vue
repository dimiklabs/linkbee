<template>
  <div class="app-shell" :class="{ 'dark': uiStore.darkMode }">

    <!-- Left navigation drawer -->
    <nav
      class="nav-drawer"
      :class="{
        collapsed: uiStore.sidebarCollapsed,
        'mobile-open': uiStore.sidebarOpen,
      }"
    >
      <!-- Logo / header area -->
      <div class="nav-header">
        <div class="nav-logo-icon">S</div>
        <span class="nav-logo-text">Shortlink</span>
      </div>

      <md-divider />

      <!-- ── Main section ─────────────────────────────────── -->
      <div class="nav-section-label">Main</div>

      <router-link
        to="/dashboard/overview"
        class="nav-item"
        :class="{ active: $route.path === '/dashboard/overview' }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">home</span>
        <span class="nav-label">Overview</span>
      </router-link>

      <router-link
        to="/dashboard/links"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/links') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">link</span>
        <span class="nav-label">Links</span>
      </router-link>

      <router-link
        to="/dashboard/analytics"
        class="nav-item"
        active-class="active"
        exact-active-class="active"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">bar_chart</span>
        <span class="nav-label">Analytics</span>
      </router-link>

      <router-link
        to="/dashboard/comparison"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/comparison') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">compare_arrows</span>
        <span class="nav-label">Compare</span>
      </router-link>

      <router-link
        to="/dashboard/reports"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/reports') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">summarize</span>
        <span class="nav-label">Reports</span>
      </router-link>

      <!-- ── Tools section ────────────────────────────────── -->
      <div class="nav-section-label">Tools</div>

      <router-link
        to="/dashboard/teams"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/teams') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">group</span>
        <span class="nav-label">Teams</span>
      </router-link>

      <router-link
        to="/dashboard/api-keys"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/api-keys') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">key</span>
        <span class="nav-label">API Keys</span>
      </router-link>

      <router-link
        to="/dashboard/webhooks"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/webhooks') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">webhook</span>
        <span class="nav-label">Webhooks</span>
      </router-link>

      <router-link
        to="/dashboard/bio"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/bio') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">account_circle</span>
        <span class="nav-label">Bio</span>
      </router-link>

      <router-link
        to="/dashboard/domains"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/domains') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">domain</span>
        <span class="nav-label">Domains</span>
      </router-link>

      <!-- ── Account section ──────────────────────────────── -->
      <div class="nav-section-label">Account</div>

      <router-link
        to="/dashboard/audit-logs"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/audit-logs') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">history</span>
        <span class="nav-label">Audit Logs</span>
      </router-link>

      <router-link
        to="/dashboard/billing"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/billing') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">credit_card</span>
        <span class="nav-label">Billing</span>
      </router-link>

      <router-link
        to="/dashboard/settings"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/settings') }"
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">settings</span>
        <span class="nav-label">Settings</span>
      </router-link>

      <router-link
        to="/dashboard/security"
        class="nav-item"
        :class="{ active: $route.path.startsWith('/dashboard/security') }"
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
        @click="uiStore.closeSidebar()"
      >
        <span class="material-symbols-outlined nav-icon">admin_panel_settings</span>
        <span class="nav-label">Admin</span>
      </router-link>

      <div class="nav-spacer"></div>

      <!-- Bottom user section -->
      <div class="nav-user">
        <div class="nav-user-row">
          <div class="nav-user-avatar">{{ authStore.userInitials }}</div>
          <div class="nav-user-info">
            <div class="nav-user-name">{{ authStore.userName }}</div>
            <div class="nav-user-email">{{ authStore.profile?.email }}</div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Mobile overlay -->
    <div
      v-if="uiStore.sidebarOpen"
      class="nav-overlay"
      @click="uiStore.closeSidebar()"
    />

    <!-- Main area -->
    <div
      class="main-area"
      :class="{ collapsed: uiStore.sidebarCollapsed }"
    >
      <!-- Top app bar -->
      <header class="top-app-bar">
        <!-- Hamburger / collapse toggle -->
        <md-icon-button
          class="sidebar-toggle"
          @click="uiStore.toggleSidebar()"
          title="Toggle sidebar"
        >
          <span class="material-symbols-outlined">menu</span>
        </md-icon-button>

        <!-- Page title -->
        <span class="page-title">{{ pageTitle }}</span>

        <!-- Dark mode toggle -->
        <md-icon-button
          :title="uiStore.darkMode ? 'Switch to light mode' : 'Switch to dark mode'"
          @click="uiStore.toggleDarkMode()"
        >
          <span class="material-symbols-outlined">
            {{ uiStore.darkMode ? 'light_mode' : 'dark_mode' }}
          </span>
        </md-icon-button>

        <!-- Notifications bell -->
        <md-icon-button title="Notifications">
          <span class="material-symbols-outlined">notifications</span>
        </md-icon-button>

        <!-- User avatar with dropdown menu -->
        <div style="position: relative">
          <div
            class="user-avatar"
            id="user-menu-anchor"
            style="cursor: pointer; width: 36px; height: 36px; border-radius: 50%; background: var(--md-sys-color-primary); color: var(--md-sys-color-on-primary); display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 14px; flex-shrink: 0;"
            @click="userMenuOpen = !userMenuOpen"
          >
            {{ authStore.userInitials }}
          </div>
          <md-menu
            :open="userMenuOpen"
            @closed="userMenuOpen = false"
            anchor="user-menu-anchor"
            y-offset="8"
          >
            <md-menu-item @click="goToSettings">
              <span class="material-symbols-outlined" slot="start">settings</span>
              <div slot="headline">Settings</div>
            </md-menu-item>
            <md-divider />
            <md-menu-item @click="handleLogout">
              <span class="material-symbols-outlined" slot="start">logout</span>
              <div slot="headline">Sign out</div>
            </md-menu-item>
          </md-menu>
        </div>
      </header>

      <!-- Page content -->
      <main class="page-content">
        <router-view />
      </main>
    </div>

  </div>
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
const userMenuOpen = ref(false);

const pageTitle = computed(() => {
  const name = route.name as string | undefined;
  if (!name) return 'Dashboard';
  const titles: Record<string, string> = {
    overview: 'Overview',
    links: 'My Links',
    'link-analytics': 'Link Analytics',
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
  return titles[name] || 'Dashboard';
});

function goToSettings() {
  userMenuOpen.value = false;
  router.push('/dashboard/settings');
}

async function handleLogout() {
  userDropdownOpen.value = false;
  userMenuOpen.value = false;
  await authStore.logout();
  router.push('/auth/login');
}

function handleClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement;
  if (!target.closest('.user-menu')) {
    userDropdownOpen.value = false;
  }
  if (!target.closest('.sidebar') && !target.closest('.sidebar-toggle')) {
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
/* ── App shell ──────────────────────────────────────────────────────────── */
.app-shell {
  display: flex;
  min-height: 100vh;
  position: relative;
  background: var(--md-sys-color-background);
}

/* ── Navigation drawer ──────────────────────────────────────────────────── */
.nav-drawer {
  width: 256px;
  min-height: 100vh;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 200;
  transition: width 0.25s cubic-bezier(0.2, 0, 0, 1), transform 0.25s cubic-bezier(0.2, 0, 0, 1);
  overflow: hidden;
  overflow-y: auto;
}

.nav-drawer.collapsed {
  width: 80px;
}

/* Logo / header area */
.nav-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  height: 64px;
  flex-shrink: 0;
}

.nav-logo-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 18px;
  flex-shrink: 0;
}

.nav-logo-text {
  font-weight: 600;
  font-size: 18px;
  white-space: nowrap;
  color: var(--md-sys-color-on-surface);
}

/* Section label */
.nav-section-label {
  padding: 8px 24px 4px;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;
  overflow: hidden;
}

/* Nav items — M3 Navigation Drawer Item */
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px;
  height: 56px;
  border-radius: 100px;
  margin: 1px 12px;
  text-decoration: none;
  color: var(--md-sys-color-on-surface-variant);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
  white-space: nowrap;
  overflow: hidden;
  position: relative;
}

.nav-item:hover {
  background: color-mix(in srgb, var(--md-sys-color-on-surface) 8%, transparent);
}

.nav-item.router-link-active,
.nav-item.active {
  background: var(--md-sys-color-secondary-container);
  color: var(--md-sys-color-on-secondary-container);
}

.nav-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.nav-label {
  white-space: nowrap;
}

/* Collapsed state */
.nav-drawer.collapsed .nav-section-label,
.nav-drawer.collapsed .nav-label,
.nav-drawer.collapsed .nav-logo-text {
  display: none;
}

.nav-drawer.collapsed .nav-item {
  justify-content: center;
  padding: 0;
  margin: 1px 8px;
}

.nav-spacer {
  flex: 1;
}

/* Bottom user section */
.nav-user {
  padding: 12px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  flex-shrink: 0;
}

.nav-user-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 8px;
  border-radius: 12px;
  overflow: hidden;
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
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.nav-user-info {
  overflow: hidden;
  min-width: 0;
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
}

.nav-drawer.collapsed .nav-user-info {
  display: none;
}

.nav-drawer.collapsed .nav-user-row {
  justify-content: center;
  padding: 4px;
}

/* ── Mobile overlay ─────────────────────────────────────────────────────── */
.nav-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 199;
}

/* ── Main area ──────────────────────────────────────────────────────────── */
.main-area {
  margin-left: 256px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.25s cubic-bezier(0.2, 0, 0, 1);
  background: var(--md-sys-color-background);
  flex: 1;
  min-width: 0;
}

.main-area.collapsed {
  margin-left: 80px;
}

/* ── Top app bar ────────────────────────────────────────────────────────── */
.top-app-bar {
  height: 64px;
  background: var(--md-sys-color-surface);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  align-items: center;
  padding: 0 16px 0 8px;
  gap: 8px;
  position: sticky;
  top: 0;
  z-index: 100;
  flex-shrink: 0;
}

.top-app-bar .page-title {
  font-size: 20px;
  font-weight: 500;
  flex: 1;
  color: var(--md-sys-color-on-surface);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-avatar {
  flex-shrink: 0;
  user-select: none;
}

/* ── Page content ───────────────────────────────────────────────────────── */
.page-content {
  flex: 1;
  padding: 24px;
  overflow-x: hidden;

  @media (max-width: 575px) {
    padding: 16px;
  }
}

/* ── Mobile responsive ──────────────────────────────────────────────────── */
@media (max-width: 768px) {
  .nav-drawer {
    transform: translateX(-100%);
    width: 256px;
  }

  .nav-drawer.mobile-open {
    transform: translateX(0);
  }

  .main-area {
    margin-left: 0 !important;
  }
}
</style>
