<template>
  <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
    <div class="flex items-center justify-between">
      <div>
        <p class="text-dark-400 text-sm">{{ title }}</p>
        <p v-if="loading" class="text-2xl font-bold text-white mt-1">
          <span class="animate-pulse">...</span>
        </p>
        <p v-else class="text-2xl font-bold text-white mt-1">{{ value }}</p>
      </div>
      <div
        class="w-12 h-12 rounded-lg flex items-center justify-center"
        :class="color"
      >
        <component :is="icon" class="w-6 h-6 text-white" />
      </div>
    </div>
    <div v-if="change !== 0" class="mt-2">
      <span
        :class="change > 0 ? 'text-green-400' : 'text-red-400'"
        class="text-sm"
      >
        {{ change > 0 ? '↑' : '↓' }} {{ Math.abs(change) }}%
      </span>
      <span class="text-dark-500 text-sm ml-1">vs last week</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  change: {
    type: Number,
    default: 0
  },
  icon: {
    type: Object,
    required: true
  },
  color: {
    type: String,
    default: 'bg-blue-600'
  },
  loading: {
    type: Boolean,
    default: false
  }
})
</script>