<template>
  <router-view />
</template>

<script setup lang="ts">
import { watch, onMounted } from 'vue';
import { useUIStore } from '@/stores/ui';

const uiStore = useUIStore();

function applyTheme(dark: boolean) {
  document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
  document.documentElement.style.colorScheme = dark ? 'dark' : 'light';
}

onMounted(() => applyTheme(uiStore.darkMode));
watch(() => uiStore.darkMode, applyTheme);
</script>
