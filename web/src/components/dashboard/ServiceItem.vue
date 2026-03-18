<template>
  <div class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg">
    <div class="flex items-center gap-3">
      <div
        class="w-2 h-2 rounded-full"
        :class="statusColor"
      />
      <span class="text-white">{{ name }}</span>
    </div>
    <div class="text-right">
      <span class="text-dark-400 text-sm">{{ count }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  status: {
    type: String,
    default: 'offline'
  },
  count: {
    type: String,
    default: ''
  }
})

const statusColor = computed(() => {
  switch (props.status) {
    case 'online':
      return 'bg-green-400'
    case 'degraded':
      return 'bg-yellow-400'
    case 'offline':
    default:
      return 'bg-red-400'
  }
})
</script>