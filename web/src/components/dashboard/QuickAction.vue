<template>
  <button
    class="flex flex-col items-center gap-2 p-4 bg-dark-700 hover:bg-dark-600 rounded-lg transition-colors group"
    @click="$emit('click')"
  >
    <div
      class="w-10 h-10 rounded-lg flex items-center justify-center transition-transform group-hover:scale-110"
      :class="color"
    >
      <component :is="icon" class="w-5 h-5 text-white" />
    </div>
    <span class="text-sm text-white">{{ name }}</span>
  </button>
</template>

<script setup>
defineProps({
  name: {
    type: String,
    required: true
  },
  icon: {
    type: Object,
    required: true
  },
  color: {
    type: String,
    default: 'bg-primary-500'
  }
})

defineEmits(['click'])
</script>