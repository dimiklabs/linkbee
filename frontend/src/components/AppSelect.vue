<template>
  <div class="app-select-wrap" :class="{ 'app-select-wrap--error': !!errorText, 'app-select-wrap--disabled': disabled }">
    <div class="app-select-field">
      <label v-if="label" class="app-select-label" :class="{ 'app-select-label--float': modelValue !== '' && modelValue !== null && modelValue !== undefined }">
        {{ label }}
      </label>
      <select
        class="app-select-native"
        :value="modelValue"
        :disabled="disabled"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
      >
        <slot />
      </select>
      <span class="app-select-arrow material-symbols-outlined">expand_more</span>
    </div>
    <div v-if="errorText" class="app-select-error">{{ errorText }}</div>
    <div v-else-if="supportingText" class="app-select-support">{{ supportingText }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue?: string | number;
  label?: string;
  disabled?: boolean;
  errorText?: string;
  supportingText?: string;
}>();

defineEmits<{ 'update:modelValue': [value: string] }>();
</script>

<style scoped lang="scss">
.app-select-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 100%;
}

.app-select-field {
  position: relative;
  display: flex;
  align-items: center;
  background: transparent;
  border: 1px solid var(--md-sys-color-outline);
  border-radius: 4px;
  height: 56px;
  transition: border-color 0.15s;

  &:focus-within {
    border-color: var(--md-sys-color-primary);
    border-width: 2px;
  }

  .app-select-wrap--error & {
    border-color: var(--md-sys-color-error);
    &:focus-within { border-color: var(--md-sys-color-error); }
  }

  .app-select-wrap--disabled & {
    border-color: color-mix(in srgb, var(--md-sys-color-outline) 38%, transparent);
    background: color-mix(in srgb, var(--md-sys-color-on-surface) 4%, transparent);
  }
}

.app-select-label {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  color: var(--md-sys-color-on-surface-variant);
  pointer-events: none;
  transition: all 0.15s;
  background: var(--md-sys-color-surface);
  padding: 0 4px;

  .app-select-field:focus-within & {
    top: 0;
    font-size: 12px;
    color: var(--md-sys-color-primary);
  }

  &--float {
    top: 0 !important;
    font-size: 12px !important;
    color: var(--md-sys-color-on-surface-variant) !important;
  }

  .app-select-wrap--error .app-select-field:focus-within & {
    color: var(--md-sys-color-error);
  }
}

.app-select-native {
  width: 100%;
  height: 100%;
  padding: 0 40px 0 14px;
  border: none;
  background: transparent;
  color: var(--md-sys-color-on-surface);
  font-family: 'Roboto', sans-serif;
  font-size: 16px;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  outline: none;

  &:disabled {
    color: color-mix(in srgb, var(--md-sys-color-on-surface) 38%, transparent);
    cursor: not-allowed;
  }

  option {
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface);
  }
}

.app-select-arrow {
  position: absolute;
  right: 10px;
  font-size: 20px;
  color: var(--md-sys-color-on-surface-variant);
  pointer-events: none;
  transition: transform 0.2s;

  .app-select-field:focus-within & {
    transform: rotate(180deg);
    color: var(--md-sys-color-primary);
  }
}

.app-select-error {
  font-size: 12px;
  color: var(--md-sys-color-error);
  padding: 0 14px;
}

.app-select-support {
  font-size: 12px;
  color: var(--md-sys-color-on-surface-variant);
  padding: 0 14px;
}
</style>
