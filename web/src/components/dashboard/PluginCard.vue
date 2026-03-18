<template>
  <div class="bg-dark-700/50 rounded-lg p-4 border border-dark-600">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-white font-medium">{{ plugin.name }}</h3>
      <span
        class="px-2 py-0.5 text-xs rounded-full"
        :class="statusClasses"
      >
        {{ plugin.status }}
      </span>
    </div>
    <p class="text-dark-400 text-sm mb-3 line-clamp-2">{{ plugin.description }}</p>
    <div class="flex items-center justify-between text-xs text-dark-500">
      <span>v{{ plugin.version }}</span>
      <span>{{ plugin.capabilities?.length || 0 }} capabilities</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  plugin: {
    type: Object,
    required: true
  }
})

const statusClasses = computed(() => {
  switch (props.plugin.status) {
    case 'running':
      return 'bg-green-500/20 text-green-400'
    case 'stopped':
      return 'bg-red-500/20 text-red-400'
    case 'error':
      return 'bg-yellow-500/20 text-yellow-400'
    default:
      return 'bg-gray-500/20 text-gray-400'
  }
})
</script>