<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Plugins</h2>
        <p class="text-slate-400 text-sm mt-1">Manage installed plugins</p>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div
        v-for="plugin in plugins"
        :key="plugin.name"
        class="bg-slate-800 rounded-xl border border-slate-700 p-6"
      >
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-4">
            <div :class="plugin.enabled ? 'bg-blue-600' : 'bg-slate-700'" class="w-12 h-12 rounded-lg flex items-center justify-center">
              <ExtensionIcon class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="text-white font-semibold">{{ plugin.name }}</h3>
              <p class="text-slate-400 text-sm mt-1">{{ plugin.description }}</p>
            </div>
          </div>
          <span
            class="px-2 py-1 text-xs rounded-full"
            :class="plugin.enabled ? 'bg-green-900/30 text-green-400' : 'bg-slate-700 text-slate-400'"
          >
            {{ plugin.enabled ? 'Enabled' : 'Disabled' }}
          </span>
        </div>

        <div class="mt-4 flex items-center justify-between">
          <span class="text-slate-500 text-sm">v{{ plugin.version }}</span>
          <div class="flex items-center gap-2">
            <button
              @click="plugin.enabled = !plugin.enabled"
              class="px-3 py-1.5 text-sm rounded-lg transition-colors"
              :class="plugin.enabled ? 'bg-red-900/30 text-red-400 hover:bg-red-900/50' : 'bg-green-900/30 text-green-400 hover:bg-green-900/50'"
            >
              {{ plugin.enabled ? 'Disable' : 'Enable' }}
            </button>
            <button class="px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-white text-sm rounded-lg transition-colors">
              Settings
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'

const ExtensionIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z' })
])

const plugins = ref([
  { name: 'Ansible', description: 'Automated deployment and configuration', enabled: true, version: '1.0.0' },
  { name: 'V2Board', description: 'V2Ray panel integration', enabled: true, version: '1.0.0' },
  { name: 'V2BX', description: 'Extended V2Board integration', enabled: true, version: '1.0.0' },
  { name: 'Agent', description: 'Remote agent management', enabled: true, version: '1.0.0' },
])
</script>