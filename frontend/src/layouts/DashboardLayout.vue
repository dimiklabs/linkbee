<template>
  <div class="dashboard-layout" :class="{ 'sidebar-collapsed': uiStore.sidebarCollapsed, 'dark-mode': uiStore.darkMode }">

    <!-- Sidebar overlay (mobile) -->
    <div
      v-if="uiStore.sidebarOpen"
      class="sidebar-overlay"
      @click="uiStore.closeSidebar()"
    ></div>

    <!-- ─── Sidebar ──────────────────────────────────────────────────────── -->
    <aside class="sidebar" :class="{ open: uiStore.sidebarOpen }">

      <!-- Sidebar header / logo -->
      <div class="sidebar-header">
        <div class="sidebar-logo">
          <div class="logo-icon">S</div>
          <span v-if="!uiStore.sidebarCollapsed" class="logo-text">Shortlink</span>
        </div>
        <button
          class="btn-icon d-none d-lg-flex"
          title="Toggle sidebar"
          @click="uiStore.toggleSidebar()"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M4.5 11.5A.5.5 0 0 1 5 11h10a.5.5 0 0 1 0 1H5a.5.5 0 0 1-.5-.5zm-2-4A.5.5 0 0 1 3 7h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm-2-4A.5.5 0 0 1 1 3h10a.5.5 0 0 1 0 1H1a.5.5 0 0 1-.5-.5z"/>
          </svg>
        </button>
      </div>

      <!-- Nav items -->
      <nav class="sidebar-nav">
        <router-link
          to="/dashboard/links"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/dashboard/links') }"
          @click="uiStore.closeSidebar()"
        >
          <span class="nav-icon">🔗</span>
          <span v-if="!uiStore.sidebarCollapsed" class="nav-label">Links</span>
        </router-link>
      </nav>

      <!-- Sidebar footer: user info + logout -->
      <div class="sidebar-footer">
        <!-- Dark mode toggle -->
        <button
          class="nav-item w-100 border-0 text-start mb-1"
          @click="uiStore.toggleDarkMode()"
          title="Toggle dark mode"
        >
          <span class="nav-icon">{{ uiStore.darkMode ? '☀️' : '🌙' }}</span>
          <span v-if="!uiStore.sidebarCollapsed" class="nav-label">
            {{ uiStore.darkMode ? 'Light mode' : 'Dark mode' }}
          </span>
        </button>

        <!-- User row -->
        <div class="user-row">
          <div class="user-avatar">{{ authStore.userInitials }}</div>
          <div v-if="!uiStore.sidebarCollapsed" class="user-info flex-fill min-w-0">
            <div class="user-name">{{ authStore.userName }}</div>
            <div class="user-email">{{ authStore.profile?.email }}</div>
          </div>
          <button
            v-if="!uiStore.sidebarCollapsed"
            class="btn-icon flex-shrink-0"
            title="Sign out"
            @click="handleLogout"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0v2z"/>
              <path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3z"/>
            </svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- ─── Main content ─────────────────────────────────────────────────── -->
    <div class="main-content">

      <!-- Top bar -->
      <header class="topbar">
        <!-- Hamburger (mobile) -->
        <button
          class="btn-icon sidebar-toggle"
          title="Open menu"
          @click="uiStore.toggleSidebar()"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M2.5 12a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5z"/>
          </svg>
        </button>

        <!-- Page title -->
        <span class="topbar-title">{{ pageTitle }}</span>

        <div class="topbar-spacer"></div>

        <!-- Dark mode toggle (topbar) -->
        <button
          class="btn-icon d-none d-md-flex"
          :title="uiStore.darkMode ? 'Switch to light mode' : 'Switch to dark mode'"
          @click="uiStore.toggleDarkMode()"
        >
          <span style="font-size: 1rem;">{{ uiStore.darkMode ? '☀️' : '🌙' }}</span>
        </button>

        <!-- User menu -->
        <div class="user-menu" @click.stop="userDropdownOpen = !userDropdownOpen">
          <button class="user-avatar-btn">
            <div class="user-avatar-sm">{{ authStore.userInitials }}</div>
            <span class="d-none d-md-block user-display-name">{{ authStore.userName }}</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
              <path fill-rule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
            </svg>
          </button>

          <div v-if="userDropdownOpen" class="user-dropdown card shadow">
            <div class="dropdown-user-info px-3 py-2 border-bottom">
              <div class="fw-medium small">{{ authStore.userName }}</div>
              <div class="text-muted" style="font-size: 0.75rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
                {{ authStore.profile?.email }}
              </div>
            </div>
            <button
              class="dropdown-action d-flex align-items-center gap-2 px-3 py-2 text-danger w-100 border-0"
              @click="handleLogout"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                <path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0v2z"/>
                <path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3z"/>
              </svg>
              Sign out
            </button>
          </div>
        </div>
      </header>

      <!-- Router outlet -->
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

const pageTitle = computed(() => {
  const name = route.name as string | undefined;
  if (!name) return 'Dashboard';
  const titles: Record<string, string> = {
    links: 'My Links',
    'link-analytics': 'Link Analytics',
  };
  return titles[name] || 'Dashboard';
});

async function handleLogout() {
  userDropdownOpen.value = false;
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
$sidebar-width: 240px;
$sidebar-collapsed-width: 64px;
$topbar-height: 60px;
$primary: #635bff;

.dashboard-layout {
  display: flex;
  min-height: 100vh;
  position: relative;
}

// ─── Sidebar ─────────────────────────────────────────────────────────────────
.sidebar {
  width: $sidebar-width;
  min-height: 100vh;
  background: #fff;
  border-right: 1px solid #e3e8ee;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width 0.2s ease, transform 0.2s ease;
  position: sticky;
  top: 0;
  align-self: flex-start;
  height: 100vh;
  overflow-y: auto;
  z-index: 100;

  @media (max-width: 991px) {
    position: fixed;
    left: 0;
    top: 0;
    transform: translateX(-100%);
    z-index: 200;
    box-shadow: none;

    &.open {
      transform: translateX(0);
      box-shadow: 4px 0 24px rgba(0, 0, 0, 0.15);
    }
  }
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 199;

  @media (min-width: 992px) {
    display: none;
  }
}

.dashboard-layout.sidebar-collapsed .sidebar {
  width: $sidebar-collapsed-width;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  height: $topbar-height;
  border-bottom: 1px solid #e3e8ee;
  flex-shrink: 0;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  overflow: hidden;

  .logo-icon {
    width: 34px;
    height: 34px;
    background: $primary;
    color: #fff;
    font-weight: 700;
    font-size: 1rem;
    border-radius: 9px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .logo-text {
    font-size: 1rem;
    font-weight: 700;
    color: #1a1f36;
    white-space: nowrap;
  }
}

.sidebar-nav {
  flex: 1;
  padding: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.625rem 0.75rem;
  border-radius: 8px;
  color: #697386;
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: background 0.15s, color 0.15s;
  cursor: pointer;
  background: transparent;

  &:hover {
    background: #f7f9fc;
    color: #1a1f36;
  }

  &.active {
    background: #eef0ff;
    color: $primary;
  }

  .nav-icon {
    font-size: 1rem;
    flex-shrink: 0;
    line-height: 1;
  }

  .nav-label {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.sidebar-footer {
  padding: 0.75rem;
  border-top: 1px solid #e3e8ee;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.user-row {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.5rem 0.5rem;
  border-radius: 8px;
  transition: background 0.15s;

  &:hover {
    background: #f7f9fc;
  }
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: $primary;
  color: #fff;
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-info {
  overflow: hidden;
}

.user-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #1a1f36;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: 0.7rem;
  color: #697386;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.min-w-0 {
  min-width: 0;
}

// ─── Main content ─────────────────────────────────────────────────────────────
.main-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

// ─── Topbar ──────────────────────────────────────────────────────────────────
.topbar {
  height: $topbar-height;
  background: #fff;
  border-bottom: 1px solid #e3e8ee;
  display: flex;
  align-items: center;
  padding: 0 1.25rem;
  gap: 0.75rem;
  position: sticky;
  top: 0;
  z-index: 99;
  flex-shrink: 0;
}

.topbar-title {
  font-weight: 600;
  font-size: 0.9375rem;
  color: #1a1f36;
  white-space: nowrap;
}

.topbar-spacer {
  flex: 1;
}

.btn-icon {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #697386;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  flex-shrink: 0;

  &:hover {
    background: #f7f9fc;
    color: #1a1f36;
  }
}

// ─── User menu ────────────────────────────────────────────────────────────────
.user-menu {
  position: relative;
}

.user-avatar-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 0.375rem 0.5rem;
  border-radius: 8px;
  transition: background 0.15s;

  &:hover {
    background: #f7f9fc;
  }
}

.user-avatar-sm {
  width: 32px;
  height: 32px;
  background: $primary;
  color: #fff;
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-display-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: #1a1f36;
  max-width: 130px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-dropdown {
  position: absolute;
  right: 0;
  top: calc(100% + 0.5rem);
  width: 220px;
  z-index: 300;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #e3e8ee;
}

.dropdown-user-info {
  background: #f7f9fc;
}

.dropdown-action {
  font-size: 0.875rem;
  background: transparent;
  text-align: left;
  transition: background 0.15s;
  cursor: pointer;

  &:hover {
    background: #fff0f0;
  }
}

// ─── Page content ──────────────────────────────────────────────────────────────
.page-content {
  flex: 1;
  padding: 1.5rem;
  background: #f7f9fc;
  overflow-x: hidden;

  @media (max-width: 575px) {
    padding: 1rem;
  }
}

// ─── Dark mode ────────────────────────────────────────────────────────────────
.dark-mode {
  .sidebar {
    background: #161b22;
    border-right-color: #30363d;
  }

  .sidebar-header {
    border-bottom-color: #30363d;
  }

  .sidebar-logo .logo-text {
    color: #e6edf3;
  }

  .nav-item {
    color: #8b949e;

    &:hover {
      background: #21262d;
      color: #e6edf3;
    }

    &.active {
      background: rgba(99, 91, 255, 0.2);
      color: $primary;
    }
  }

  .sidebar-footer {
    border-top-color: #30363d;
  }

  .user-row {
    &:hover {
      background: #21262d;
    }
  }

  .user-name {
    color: #e6edf3;
  }

  .topbar {
    background: #161b22;
    border-bottom-color: #30363d;
  }

  .topbar-title {
    color: #e6edf3;
  }

  .btn-icon {
    color: #8b949e;

    &:hover {
      background: #21262d;
      color: #e6edf3;
    }
  }

  .user-avatar-btn:hover {
    background: #21262d;
  }

  .user-display-name {
    color: #e6edf3;
  }

  .user-dropdown {
    background: #161b22;
    border-color: #30363d;
  }

  .dropdown-user-info {
    background: #0d1117;

    .fw-medium {
      color: #e6edf3;
    }
  }

  .dropdown-action:hover {
    background: rgba(255, 59, 48, 0.12);
  }

  .page-content {
    background: #0d1117;
  }
}
</style>
