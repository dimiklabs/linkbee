import { defineStore } from 'pinia';
import { computed, ref, watch } from 'vue';

export const useUIStore = defineStore('ui', () => {
  const sidebarCollapsed = ref(false);
  const sidebarOpen = ref(false);
  const userMenuOpen = ref(false);
  const darkMode = ref(false);

  // Initialize from localStorage or system preference
  const storedDark = localStorage.getItem('darkMode');
  if (storedDark !== null) {
    darkMode.value = storedDark === 'true';
  } else {
    darkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches;
  }

  const storedCollapsed = localStorage.getItem('sidebarCollapsed');
  if (storedCollapsed !== null) {
    sidebarCollapsed.value = storedCollapsed === 'true';
  }

  // Persist darkMode to localStorage
  watch(darkMode, (newValue) => {
    localStorage.setItem('darkMode', String(newValue));
  });

  watch(sidebarCollapsed, (newValue) => {
    localStorage.setItem('sidebarCollapsed', String(newValue));
  });

  const isMobile = computed(() => window.innerWidth < 992);

  const toggleSidebar = () => {
    if (isMobile.value) {
      sidebarOpen.value = !sidebarOpen.value;
    } else {
      sidebarCollapsed.value = !sidebarCollapsed.value;
    }
  };

  const closeSidebar = () => { sidebarOpen.value = false; };
  const toggleUserMenu = () => { userMenuOpen.value = !userMenuOpen.value; };
  const closeUserMenu = () => { userMenuOpen.value = false; };
  const handleResize = () => { if (window.innerWidth >= 992) sidebarOpen.value = false; };
  const toggleDarkMode = () => { darkMode.value = !darkMode.value; };
  const setDarkMode = (isDark: boolean) => { darkMode.value = isDark; };

  // Signals LinksPage to open the create-link modal after navigation
  const pendingCreateLink = ref(false);
  const triggerCreateLink = () => { pendingCreateLink.value = true; };
  const clearPendingCreateLink = () => { pendingCreateLink.value = false; };

  return {
    sidebarCollapsed, sidebarOpen, userMenuOpen, darkMode, isMobile,
    toggleSidebar, closeSidebar, toggleUserMenu, closeUserMenu,
    handleResize, toggleDarkMode, setDarkMode,
    pendingCreateLink, triggerCreateLink, clearPendingCreateLink,
  };
});

// Alias for backward compatibility
export const useUiStore = useUIStore;
