<template>
  <div class="flex items-start gap-3 p-3 bg-dark-700/50 rounded-lg">
    <div
      class="w-8 h-8 rounded-full flex items-center justify-center shrink-0"
      :class="typeColor"
    >
      <component :is="typeIcon" class="w-4 h-4" />
    </div>
    <div class="flex-1 min-w-0">
      <p class="text-white text-sm">{{ activity.message }}</p>
      <p class="text-dark-500 text-xs mt-1">{{ activity.time || formatTime(activity.timestamp) }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  CheckCircleIcon,
  UserIcon,
  ClockIcon,
  BoltIcon,
  XCircleIcon,
  ServerIcon,
  DocumentTextIcon
} from '@heroicons/vue/24/outline'

const props = defineProps({
  activity: {
    type: Object,
    required: true
  }
})

const typeColor = computed(() => {
  const colors = {
    deploy: 'bg-green-600',
    user: 'bg-blue-600',
    backup: 'bg-purple-600',
    alert: 'bg-yellow-600',
    error: 'bg-red-600',
    node: 'bg-cyan-600',
    playbook: 'bg-indigo-600'
  }
  return colors[props.activity.type] || 'bg-gray-600'
})

const typeIcon = computed(() => {
  const icons = {
    deploy: CheckCircleIcon,
    user: UserIcon,
    backup: ClockIcon,
    alert: BoltIcon,
    error: XCircleIcon,
    node: ServerIcon,
    playbook: DocumentTextIcon
  }
  return icons[props.activity.type] || ClockIcon
})

function formatTime(timestamp) {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) return 'Just now'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} minutes ago`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} hours ago`
  return date.toLocaleDateString()
}
</script>