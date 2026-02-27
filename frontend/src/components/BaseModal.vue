<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="modelValue"
        class="modal-backdrop"
        @click.self="handleBackdropClick"
        @keydown.esc="handleEsc"
      >
        <Transition name="modal-slide">
          <div
            v-if="modelValue"
            class="modal-container"
            :class="[`modal-container--${size}`, containerClass]"
            role="dialog"
            aria-modal="true"
          >
            <!-- Header -->
            <div class="modal-header">
              <div class="modal-header-content">
                <slot name="headline" />
              </div>
              <button
                class="modal-close-btn"
                type="button"
                aria-label="Close"
                @click="$emit('update:modelValue', false)"
              >
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>

            <!-- Divider -->
            <div class="modal-divider" />

            <!-- Scrollable content -->
            <div class="modal-body">
              <slot />
            </div>

            <!-- Actions footer -->
            <template v-if="$slots.actions">
              <div class="modal-divider" />
              <div class="modal-footer">
                <slot name="actions" />
              </div>
            </template>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { watch, onMounted, onUnmounted } from 'vue';

const props = withDefaults(defineProps<{
  modelValue: boolean;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  persistent?: boolean;
  containerClass?: string;
}>(), {
  size: 'md',
  persistent: false,
  containerClass: '',
});

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  'closed': [];
}>();

function handleBackdropClick() {
  if (!props.persistent) {
    emit('update:modelValue', false);
    emit('closed');
  }
}

function handleEsc() {
  if (!props.persistent) {
    emit('update:modelValue', false);
    emit('closed');
  }
}

// Lock body scroll when modal is open
watch(() => props.modelValue, (open) => {
  if (open) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
  }
}, { immediate: true });

onUnmounted(() => {
  document.body.style.overflow = '';
});
</script>

<style scoped lang="scss">
/* ── Backdrop ───────────────────────────────────────────────────────────── */
.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.52);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  backdrop-filter: blur(2px);
  -webkit-backdrop-filter: blur(2px);
}

/* ── Container ──────────────────────────────────────────────────────────── */
.modal-container {
  background: var(--md-sys-color-surface);
  border-radius: 28px;
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 48px);
  width: 100%;
  box-shadow:
    0 8px 12px 6px rgba(0,0,0,0.15),
    0 4px 4px rgba(0,0,0,0.3);
  position: relative;
  overflow: hidden;

  /* Size variants */
  &--sm  { max-width: 400px; }
  &--md  { max-width: 560px; }
  &--lg  { max-width: 720px; }
  &--xl  { max-width: 900px; }
}

/* ── Header ─────────────────────────────────────────────────────────────── */
.modal-header {
  display: flex;
  align-items: center;
  padding: 20px 24px 18px;
  flex-shrink: 0;
  gap: 8px;
}

.modal-header-content {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.modal-close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
  transition: background 0.15s;

  &:hover { background: color-mix(in srgb, var(--md-sys-color-on-surface) 10%, transparent); }
  &:focus-visible {
    outline: 2px solid var(--md-sys-color-primary);
    outline-offset: 2px;
  }

  .material-symbols-outlined { font-size: 20px; }
}

/* ── Divider ────────────────────────────────────────────────────────────── */
.modal-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  flex-shrink: 0;
}

/* ── Scrollable body ────────────────────────────────────────────────────── */
.modal-body {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 20px 24px;
  min-height: 0;

  /* thin scrollbar */
  &::-webkit-scrollbar { width: 4px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb {
    background: var(--md-sys-color-outline-variant);
    border-radius: 4px;
  }
}

/* ── Footer actions ─────────────────────────────────────────────────────── */
.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 14px 20px;
  flex-shrink: 0;
}

/* ── Animations ─────────────────────────────────────────────────────────── */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-slide-enter-active {
  transition: opacity 0.22s ease, transform 0.22s cubic-bezier(0.34, 1.2, 0.64, 1);
}
.modal-slide-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}
.modal-slide-enter-from {
  opacity: 0;
  transform: translateY(16px) scale(0.97);
}
.modal-slide-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.98);
}

/* ── Mobile full-screen ─────────────────────────────────────────────────── */
@media (max-width: 600px) {
  .modal-backdrop { padding: 0; align-items: flex-end; }
  .modal-container {
    max-width: 100%;
    border-radius: 20px 20px 0 0;
    max-height: 92vh;
  }
}
</style>
