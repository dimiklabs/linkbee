<template>
  <div class="upgrade-wall">
    <div class="upgrade-wall__icon-wrap">
      <span class="material-symbols-outlined upgrade-wall__icon">{{ icon }}</span>
    </div>
    <h2 class="md-headline-small upgrade-wall__title">{{ title }}</h2>
    <p class="md-body-medium upgrade-wall__desc">{{ description }}</p>
    <router-link to="/dashboard/billing">
      <button class="btn-filled upgrade-wall__btn">
        <span class="material-symbols-outlined">workspace_premium</span>
        Upgrade Plan
      </button>
    </router-link>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  icon?: string;
  title?: string;
  description?: string;
}>(), {
  icon: 'lock',
  title: 'Pro feature',
  description: 'Upgrade to Pro or Growth to unlock this feature.',
});
</script>

<style scoped lang="scss">
.upgrade-wall {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  text-align: center;
  gap: 16px;
  padding: 40px 24px;
}

.upgrade-wall__icon-wrap {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  background: var(--md-sys-color-primary-container);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.upgrade-wall__icon {
  font-size: 44px;
  color: var(--md-sys-color-on-primary-container);
}

.upgrade-wall__title {
  color: var(--md-sys-color-on-surface);
  margin: 0;
}

.upgrade-wall__desc {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  max-width: 420px;
  line-height: 1.6;
}

.upgrade-wall__btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding: 0 24px;
  height: 44px;
  border: none;
  border-radius: 10px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.15s;

  .material-symbols-outlined { font-size: 20px; }

  &:hover { opacity: 0.88; }
}
</style>
