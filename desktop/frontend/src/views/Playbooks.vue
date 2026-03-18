<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Playbooks</h2>
        <p class="text-slate-400 text-sm mt-1">Automated deployment and configuration</p>
      </div>
      <button class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center gap-2">
        <PlusIcon class="w-5 h-5" />
        New Playbook
      </button>
    </div>

    <div class="space-y-4">
      <div
        v-for="playbook in playbooks"
        :key="playbook.name"
        class="bg-slate-800 rounded-xl border border-slate-700 p-6"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div :class="getStatusBgColor(playbook.status)" class="w-10 h-10 rounded-lg flex items-center justify-center">
              <component :is="getStatusIcon(playbook.status)" class="w-5 h-5 text-white" />
            </div>
            <div>
              <h3 class="text-white font-semibold">{{ playbook.name }}</h3>
              <p class="text-slate-400 text-sm">Last run: {{ playbook.lastRun }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span
              class="px-2 py-1 text-xs rounded-full"
              :class="getStatusColor(playbook.status)"
            >
              {{ playbook.status }}
            </span>
            <button class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">
              Run
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'

const PlusIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4v16m8-8H4' })
])

const CheckIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' })
])

const XIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z' })
])

const ClockIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' })
])

const playbooks = ref([
  { name: 'deploy_node.yml', lastRun: '2024-03-15 12:34', status: 'Success' },
  { name: 'update_certificates.yml', lastRun: '2024-03-14 08:00', status: 'Success' },
  { name: 'backup_database.yml', lastRun: '2024-03-15 00:00', status: 'Running' },
  { name: 'cleanup_logs.yml', lastRun: '2024-03-13 06:00', status: 'Failed' },
  { name: 'deploy_agent.yml', lastRun: '2024-03-12 18:30', status: 'Success' },
])

function getStatusColor(status) {
  const colors = {
    'Success': 'bg-green-900/30 text-green-400',
    'Failed': 'bg-red-900/30 text-red-400',
    'Running': 'bg-yellow-900/30 text-yellow-400'
  }
  return colors[status] || 'bg-gray-900/30 text-gray-400'
}

function getStatusBgColor(status) {
  const colors = {
    'Success': 'bg-green-600',
    'Failed': 'bg-red-600',
    'Running': 'bg-yellow-600'
  }
  return colors[status] || 'bg-gray-600'
}

function getStatusIcon(status) {
  const icons = {
    'Success': CheckIcon,
    'Failed': XIcon,
    'Running': ClockIcon
  }
  return icons[status] || ClockIcon
}
</script>