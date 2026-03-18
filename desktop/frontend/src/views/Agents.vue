<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Agents</h2>
        <p class="text-slate-400 text-sm mt-1">Connected devices and remote agents</p>
      </div>
      <button class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors">
        Register Agent
      </button>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div
        v-for="agent in agents"
        :key="agent.id"
        class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden"
      >
        <div
          class="h-1"
          :class="agent.status === 'online' ? 'bg-green-500' : 'bg-red-500'"
        ></div>
        <div class="p-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div :class="agent.status === 'online' ? 'bg-green-600' : 'bg-red-600'" class="w-10 h-10 rounded-lg flex items-center justify-center">
                <ComputerIcon class="w-5 h-5 text-white" />
              </div>
              <div>
                <h3 class="text-white font-medium">{{ agent.name }}</h3>
                <p class="text-slate-400 text-sm">{{ agent.host }}</p>
              </div>
            </div>
            <span
              class="px-2 py-1 text-xs rounded-full"
              :class="agent.status === 'online' ? 'bg-green-900/30 text-green-400' : 'bg-red-900/30 text-red-400'"
            >
              {{ agent.status }}
            </span>
          </div>

          <div class="mt-4 grid grid-cols-3 gap-2 text-center text-sm">
            <div>
              <p class="text-slate-400 text-xs">CPU</p>
              <p class="text-white">{{ agent.cpu }}%</p>
            </div>
            <div>
              <p class="text-slate-400 text-xs">Memory</p>
              <p class="text-white">{{ agent.memory }}%</p>
            </div>
            <div>
              <p class="text-slate-400 text-xs">Disk</p>
              <p class="text-white">{{ agent.disk }}%</p>
            </div>
          </div>

          <div class="mt-4 flex items-center justify-between">
            <span class="text-slate-500 text-sm">Uptime: {{ agent.uptime }}</span>
            <button
              :disabled="agent.status === 'offline'"
              class="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors disabled:opacity-50"
            >
              Connect
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'

const ComputerIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' })
])

const agents = ref([
  { id: 1, name: 'web-01', host: '10.0.0.101', status: 'online', cpu: 45, memory: 62, disk: 35, uptime: '2d 4h' },
  { id: 2, name: 'db-01', host: '10.0.0.102', status: 'online', cpu: 78, memory: 85, disk: 60, uptime: '5d 12h' },
  { id: 3, name: 'cache-01', host: '10.0.0.103', status: 'offline', cpu: 0, memory: 0, disk: 45, uptime: '-' },
  { id: 4, name: 'dev-ws', host: '10.0.0.200', status: 'online', cpu: 25, memory: 45, disk: 70, uptime: '8h 30m' },
])
</script>