<template>
  <div
    class="flex items-start gap-3 p-3 rounded-lg"
    :class="alertClasses"
  >
    <div class="shrink-0 mt-0.5">
      <ExclamationTriangleIcon v-if="alert.level === 'warning'" class="w-5 h-5 text-yellow-400" />
      <XCircleIcon v-else class="w-5 h-5 text-red-400" />
    </div>
    <div class="flex-1 min-w-0">
      <p class="text-sm font-medium" :class="alert.level === 'warning' ? 'text-yellow-400' : 'text-red-400'">
        {{ alert.title }}
      </p>
      <p class="text-sm text-dark-300 mt-1">{{ alert.message }}</p>
      <p class="text-xs text-dark-500 mt-1">{{ formatTime(alert.timestamp) }}</p>
    </div>
    <button
      @click="$emit('dismiss')"
      class="text-dark-400 hover:text-white transition-colors"
    >
      <XMarkIcon class="w-4 h-4" />
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ExclamationTriangleIcon, XCircleIcon, XMarkIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  alert: {
    type: Object,
    required: true
  }
})

defineEmits(['dismiss'])

const alertClasses = computed(() => {
  return props.alert.level === 'warning'
    ? 'bg-yellow-500/10 border border-yellow-500/20'
    : 'bg-red-500/10 border border-red-500/20'
})

function formatTime(timestamp) {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString()
}
</script>