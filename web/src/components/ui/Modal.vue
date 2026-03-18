<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/60 backdrop-blur-sm"
          @click="closeOnBackdrop && close()"
        ></div>

        <!-- Modal -->
        <div
          :class="sizeClasses"
          class="relative bg-dark-800 rounded-2xl border border-dark-700 shadow-2xl max-h-[90vh] overflow-hidden flex flex-col animate-in fade-in zoom-in-95 duration-200"
        >
          <!-- Header -->
          <div v-if="title || $slots.header" class="px-6 py-4 border-b border-dark-700 flex items-center justify-between shrink-0">
            <slot name="header">
              <h2 class="text-xl font-semibold text-white">{{ title }}</h2>
            </slot>
            <button
              v-if="showClose"
              @click="close()"
              class="p-2 hover:bg-dark-700 rounded-lg transition-colors"
            >
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <!-- Body -->
          <div class="flex-1 overflow-y-auto p-6">
            <slot></slot>
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="px-6 py-4 border-t border-dark-700 flex items-center justify-end gap-3 shrink-0">
            <slot name="footer"></slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl', 'full'].includes(value)
  },
  showClose: {
    type: Boolean,
    default: true
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'close'])

const sizeClasses = computed(() => {
  const sizes = {
    sm: 'w-full max-w-sm',
    md: 'w-full max-w-md',
    lg: 'w-full max-w-lg',
    xl: 'w-full max-w-xl',
    full: 'w-full max-w-4xl'
  }
  return sizes[props.size]
})

function close() {
  emit('update:modelValue', false)
  emit('close')
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.95);
}
</style>