<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">System Logs</h2>
        <p class="text-slate-400 text-sm mt-1">Real-time log monitoring</p>
      </div>
      <button
        @click="isStreaming = !isStreaming"
        :class="isStreaming ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'"
        class="px-4 py-2 text-white rounded-lg transition-colors"
      >
        {{ isStreaming ? 'Stop Stream' : 'Start Stream' }}
      </button>
    </div>

    <div class="flex gap-4">
      <select
        v-model="levelFilter"
        class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">All Levels</option>
        <option value="INFO">INFO</option>
        <option value="WARN">WARN</option>
        <option value="ERROR">ERROR</option>
      </select>
      <button
        @click="logs = []"
        class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
      >
        Clear
      </button>
    </div>

    <div class="bg-slate-900 rounded-xl border border-slate-700 font-mono text-sm max-h-[500px] overflow-y-auto">
      <div
        v-for="(log, index) in filteredLogs"
        :key="index"
        class="px-4 py-2 border-b border-slate-700/50 hover:bg-slate-700/30"
      >
        <span class="text-slate-500">{{ log.time }}</span>
        <span
          class="ml-2 px-1 rounded"
          :class="getLevelColor(log.level)"
        >
          {{ log.level }}
        </span>
        <span class="ml-2 text-slate-400">[{{ log.source }}]</span>
        <span class="ml-2 text-white">{{ log.message }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const levelFilter = ref('')
const isStreaming = ref(false)
const logs = ref([])
let interval = null

const filteredLogs = computed(() => {
  if (!levelFilter.value) return logs.value
  return logs.value.filter(log => log.level === levelFilter.value)
})

function getLevelColor(level) {
  const colors = {
    'INFO': 'text-blue-400',
    'WARN': 'text-yellow-400',
    'ERROR': 'text-red-400'
  }
  return colors[level] || 'text-gray-400'
}

function addRandomLog() {
  const levels = ['INFO', 'INFO', 'INFO', 'WARN', 'ERROR']
  const sources = ['api', 'node', 'auth', 'plugin', 'system']
  const messages = [
    'Request processed successfully',
    'Node connection established',
    'User authenticated',
    'High memory usage detected',
    'Connection timeout',
  ]

  logs.value.unshift({
    time: new Date().toLocaleTimeString(),
    level: levels[Math.floor(Math.random() * levels.length)],
    source: sources[Math.floor(Math.random() * sources.length)],
    message: messages[Math.floor(Math.random() * messages.length)]
  })

  if (logs.value.length > 100) {
    logs.value.pop()
  }
}

onMounted(() => {
  for (let i = 0; i < 50; i++) {
    addRandomLog()
  }
})
</script>