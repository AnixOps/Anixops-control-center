<template>
  <div
    :class="toastClasses"
    class="fixed z-50 rounded-lg shadow-lg flex items-center gap-3 min-w-[300px] max-w-md transition-all"
  >
    <div class="shrink-0">
      <CheckCircleIcon v-if="type === 'success'" class="w-5 h-5 text-green-400" />
      <XCircleIcon v-else-if="type === 'error'" class="w-5 h-5 text-red-400" />
      <ExclamationTriangleIcon v-else-if="type === 'warning'" class="w-5 h-5 text-yellow-400" />
      <InformationCircleIcon v-else class="w-5 h-5 text-blue-400" />
    </div>
    <div class="flex-1 py-1">
      <p v-if="title" class="text-white font-medium">{{ title }}</p>
      <p class="text-dark-200 text-sm" :class="{ 'mt-1': title }">{{ message }}</p>
    </div>
    <button
      @click="$emit('close')"
      class="shrink-0 p-1 hover:bg-white/10 rounded transition-colors"
    >
      <XMarkIcon class="w-4 h-4 text-dark-300" />
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

const props = defineProps({
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    required: true
  }
})

defineEmits(['close'])

const toastClasses = computed(() => {
  const variants = {
    success: 'bg-green-900/90 border border-green-700',
    error: 'bg-red-900/90 border border-red-700',
    warning: 'bg-yellow-900/90 border border-yellow-700',
    info: 'bg-blue-900/90 border border-blue-700'
  }
  return variants[props.type]
})
</script>