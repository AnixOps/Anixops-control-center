<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between text-sm">
      <span class="text-dark-300">{{ label }}</span>
      <span class="text-white font-medium">{{ value }}%</span>
    </div>
    <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
      <div
        class="h-full rounded-full transition-all duration-500"
        :class="color"
        :style="{ width: `${percentage}%` }"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: Number,
    default: 0
  },
  max: {
    type: Number,
    default: 100
  },
  color: {
    type: String,
    default: 'bg-blue-500'
  }
})

const percentage = computed(() => {
  return Math.min(100, Math.max(0, (props.value / props.max) * 100))
})
</script>