<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { useUiStore } from '@/stores/ui';

const uiStore = useUiStore();

onMounted(() => {
  uiStore.initDarkMode();
  window.addEventListener('resize', uiStore.handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', uiStore.handleResize);
});
</script>

<template>
  <router-view/>
</template>

<style lang="scss">
@import '@/assets/styles/main.scss';
</style>
