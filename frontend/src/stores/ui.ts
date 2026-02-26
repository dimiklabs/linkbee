import { defineStore } from 'pinia';
import { computed, ref, watch } from 'vue';

export const useUiStore = defineStore('ui', () => {
  const sidebarCollapsed = ref(false);
  const sidebarOpen = ref(false);
  const userMenuOpen = ref(false);
  const darkMode = ref(false);

  const initDarkMode = () => {
    // Restore dark mode preference
    const storedDark = localStorage.getItem('darkMode');
    if (storedDark !== null) {
      darkMode.value = storedDark === 'true';
    } else {
      darkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    applyDarkMode(darkMode.value);

    // Restore sidebar collapsed state (desktop only)
    const storedCollapsed = localStorage.getItem('sidebarCollapsed');
    if (storedCollapsed !== null) {
      sidebarCollapsed.value = storedCollapsed === 'true';
    }
  };

  const applyDarkMode = (isDark: boolean) => {
    if (isDark) {
      document.documentElement.classList.add('dark-mode');
      document.body.classList.add('dark-mode');
    } else {
      document.documentElement.classList.remove('dark-mode');
      document.body.classList.remove('dark-mode');
    }
  };

  watch(darkMode, (newValue) => {
    localStorage.setItem('darkMode', String(newValue));
    applyDarkMode(newValue);
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

  return {
    sidebarCollapsed, sidebarOpen, userMenuOpen, darkMode, isMobile,
    toggleSidebar, closeSidebar, toggleUserMenu, closeUserMenu,
    handleResize, initDarkMode, toggleDarkMode, setDarkMode,
  };
});
