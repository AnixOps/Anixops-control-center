<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">System Logs</h1>
        <p class="text-dark-400 text-sm mt-1">Real-time log monitoring and analysis</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="toggleStream"
          :class="[
            'px-4 py-2 rounded-lg transition-colors flex items-center gap-2',
            isStreaming ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'
          ]"
        >
          <span :class="isStreaming ? 'bg-white' : 'bg-white'" class="w-2 h-2 rounded-full"></span>
          {{ isStreaming ? 'Stop Stream' : 'Start Stream' }}
        </button>
        <button
          @click="exportLogs"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <ArrowDownTrayIcon class="w-5 h-5" />
          Export
        </button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Total Logs</p>
            <p class="text-2xl font-bold text-white">{{ logStats.total.toLocaleString() }}</p>
          </div>
          <DocumentTextIcon class="w-8 h-8 text-primary-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Errors</p>
            <p class="text-2xl font-bold text-red-400">{{ logStats.errors }}</p>
          </div>
          <ExclamationCircleIcon class="w-8 h-8 text-red-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Warnings</p>
            <p class="text-2xl font-bold text-yellow-400">{{ logStats.warnings }}</p>
          </div>
          <ExclamationTriangleIcon class="w-8 h-8 text-yellow-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Info</p>
            <p class="text-2xl font-bold text-blue-400">{{ logStats.info }}</p>
          </div>
          <InformationCircleIcon class="w-8 h-8 text-blue-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Rate</p>
            <p class="text-2xl font-bold text-green-400">{{ logStats.rate }}/s</p>
          </div>
          <BoltIcon class="w-8 h-8 text-green-400" />
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-dark-800 rounded-xl border border-dark-700 p-4">
      <div class="flex flex-wrap gap-4">
        <div class="flex-1 min-w-[200px]">
          <div class="relative">
            <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-400" />
            <input
              v-model="search"
              type="text"
              placeholder="Search logs..."
              class="w-full pl-10 pr-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-400 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
        </div>
        <select
          v-model="levelFilter"
          multiple
          class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          <option value="error">ERROR</option>
          <option value="warn">WARN</option>
          <option value="info">INFO</option>
          <option value="debug">DEBUG</option>
        </select>
        <select
          v-model="sourceFilter"
          class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          <option value="">All Sources</option>
          <option value="api">API</option>
          <option value="node">Node</option>
          <option value="auth">Auth</option>
          <option value="plugin">Plugin</option>
          <option value="system">System</option>
        </select>
        <input
          v-model="dateFilter"
          type="datetime-local"
          class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        />
        <button
          @click="clearFilters"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Log Viewer -->
    <div class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden">
      <!-- Toolbar -->
      <div class="px-4 py-3 border-b border-dark-700 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              v-model="autoScroll"
              type="checkbox"
              class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
            />
            <span class="text-dark-400 text-sm">Auto-scroll</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              v-model="showTimestamp"
              type="checkbox"
              class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
            />
            <span class="text-dark-400 text-sm">Show timestamps</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              v-model="showSource"
              type="checkbox"
              class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
            />
            <span class="text-dark-400 text-sm">Show source</span>
          </label>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="clearLogs"
            class="px-3 py-1.5 bg-dark-700 hover:bg-dark-600 text-white text-sm rounded-lg transition-colors flex items-center gap-1"
          >
            <TrashIcon class="w-4 h-4" />
            Clear
          </button>
          <button
            @click="loadMore"
            :disabled="loading"
            class="px-3 py-1.5 bg-dark-700 hover:bg-dark-600 text-white text-sm rounded-lg transition-colors flex items-center gap-1"
          >
            <ArrowPathIcon :class="['w-4 h-4', { 'animate-spin': loading }]" />
            Load More
          </button>
        </div>
      </div>

      <!-- Log Content -->
      <div
        ref="logContainer"
        class="font-mono text-sm max-h-[600px] overflow-y-auto bg-dark-900"
      >
        <div
          v-for="(log, index) in filteredLogs"
          :key="index"
          class="px-4 py-2 border-b border-dark-700/50 hover:bg-dark-700/30 cursor-pointer"
          @click="showLogDetail(log)"
        >
          <div class="flex items-start gap-3">
            <span v-if="showTimestamp" class="text-dark-500 shrink-0">{{ log.time }}</span>
            <span
              class="px-1.5 py-0.5 rounded text-xs font-bold shrink-0"
              :class="getLevelColor(log.level)"
            >
              {{ log.level }}
            </span>
            <span v-if="showSource" class="text-dark-400 text-xs shrink-0">[{{ log.source }}]</span>
            <span class="text-white flex-1 break-all">{{ log.message }}</span>
          </div>
          <div v-if="log.metadata" class="mt-1 ml-20 text-xs text-dark-500">
            {{ JSON.stringify(log.metadata) }}
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="filteredLogs.length === 0" class="py-12 text-center">
          <DocumentTextIcon class="w-12 h-12 text-dark-500 mx-auto mb-4" />
          <p class="text-dark-400">No logs found</p>
        </div>

        <!-- Loading Indicator -->
        <div v-if="loading" class="py-4 text-center">
          <ArrowPathIcon class="w-6 h-6 text-primary-400 animate-spin mx-auto" />
        </div>
      </div>
    </div>

    <!-- Log Detail Modal -->
    <Teleport to="body">
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showDetailModal = false"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">Log Details</h2>
            <button @click="showDetailModal = false" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <div v-if="selectedLog" class="p-6 space-y-4">
            <div class="grid grid-cols-3 gap-4">
              <div class="bg-dark-700 rounded-lg p-3">
                <p class="text-dark-400 text-xs">Level</p>
                <span
                  class="px-2 py-1 text-xs rounded font-bold"
                  :class="getLevelColor(selectedLog.level)"
                >
                  {{ selectedLog.level }}
                </span>
              </div>
              <div class="bg-dark-700 rounded-lg p-3">
                <p class="text-dark-400 text-xs">Source</p>
                <p class="text-white">{{ selectedLog.source }}</p>
              </div>
              <div class="bg-dark-700 rounded-lg p-3">
                <p class="text-dark-400 text-xs">Time</p>
                <p class="text-white">{{ selectedLog.time }}</p>
              </div>
            </div>

            <div>
              <p class="text-dark-400 text-sm mb-1">Message</p>
              <div class="bg-dark-700 rounded-lg p-3">
                <p class="text-white font-mono text-sm whitespace-pre-wrap">{{ selectedLog.message }}</p>
              </div>
            </div>

            <div v-if="selectedLog.stack">
              <p class="text-dark-400 text-sm mb-1">Stack Trace</p>
              <div class="bg-dark-700 rounded-lg p-3 max-h-60 overflow-y-auto">
                <pre class="text-red-400 font-mono text-xs whitespace-pre-wrap">{{ selectedLog.stack }}</pre>
              </div>
            </div>

            <div v-if="selectedLog.metadata">
              <p class="text-dark-400 text-sm mb-1">Metadata</p>
              <div class="bg-dark-700 rounded-lg p-3">
                <pre class="text-white font-mono text-xs whitespace-pre-wrap">{{ JSON.stringify(selectedLog.metadata, null, 2) }}</pre>
              </div>
            </div>

            <div class="flex items-center justify-end gap-3 pt-4 border-t border-dark-700">
              <button
                @click="copyLog"
                class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
              >
                <ClipboardDocumentIcon class="w-4 h-4" />
                Copy
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { logsApi } from '@/api'
import {
  DocumentTextIcon,
  MagnifyingGlassIcon,
  ArrowDownTrayIcon,
  TrashIcon,
  XMarkIcon,
  ExclamationCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  BoltIcon,
  ArrowPathIcon,
  ClipboardDocumentIcon
} from '@heroicons/vue/24/outline'

// State
const search = ref('')
const levelFilter = ref(['error', 'warn', 'info'])
const sourceFilter = ref('')
const dateFilter = ref('')
const autoScroll = ref(true)
const showTimestamp = ref(true)
const showSource = ref(true)
const isStreaming = ref(false)
const loading = ref(false)
const showDetailModal = ref(false)
const selectedLog = ref(null)
const logs = ref([])
const logContainer = ref(null)
let ws = null

// Computed
const logStats = computed(() => ({
  total: logs.value.length,
  errors: logs.value.filter(l => l.level === 'ERROR').length,
  warnings: logs.value.filter(l => l.level === 'WARN').length,
  info: logs.value.filter(l => l.level === 'INFO').length,
  rate: Math.floor(Math.random() * 10) + 1
}))

const filteredLogs = computed(() => {
  let result = logs.value

  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(log =>
      log.message.toLowerCase().includes(searchLower) ||
      log.source.toLowerCase().includes(searchLower)
    )
  }

  if (levelFilter.value.length > 0) {
    result = result.filter(log => levelFilter.value.includes(log.level.toLowerCase()))
  }

  if (sourceFilter.value) {
    result = result.filter(log => log.source === sourceFilter.value)
  }

  return result
})

// Methods
function getLevelColor(level) {
  const colors = {
    ERROR: 'bg-red-900/50 text-red-400',
    WARN: 'bg-yellow-900/50 text-yellow-400',
    INFO: 'bg-blue-900/50 text-blue-400',
    DEBUG: 'bg-gray-900/50 text-gray-400'
  }
  return colors[level] || 'bg-gray-900/50 text-gray-400'
}

function toggleStream() {
  if (isStreaming.value) {
    stopStream()
  } else {
    startStream()
  }
}

function startStream() {
  isStreaming.value = true
  // Simulate WebSocket connection
  const interval = setInterval(() => {
    if (!isStreaming.value) {
      clearInterval(interval)
      return
    }
    addRandomLog()
  }, 1000)
}

function stopStream() {
  isStreaming.value = false
}

function addRandomLog() {
  const levels = ['INFO', 'INFO', 'INFO', 'WARN', 'ERROR']
  const sources = ['api', 'node', 'auth', 'plugin', 'system']
  const messages = [
    'Request processed successfully',
    'Node connection established',
    'User authenticated',
    'Plugin loaded',
    'Cache cleared',
    'Database query executed',
    'Configuration updated',
    'High memory usage detected',
    'Connection timeout',
    'Authentication failed'
  ]

  const log = {
    time: new Date().toLocaleTimeString(),
    level: levels[Math.floor(Math.random() * levels.length)],
    source: sources[Math.floor(Math.random() * sources.length)],
    message: messages[Math.floor(Math.random() * messages.length)]
  }

  logs.value.unshift(log)

  if (autoScroll.value) {
    nextTick(() => {
      if (logContainer.value) {
        logContainer.value.scrollTop = 0
      }
    })
  }
}

async function loadMore() {
  loading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 500))
    // Add more mock logs
    for (let i = 0; i < 50; i++) {
      logs.value.push({
        time: new Date(Date.now() - i * 60000).toLocaleTimeString(),
        level: ['INFO', 'WARN', 'ERROR'][Math.floor(Math.random() * 3)],
        source: ['api', 'node', 'auth'][Math.floor(Math.random() * 3)],
        message: `Log entry ${logs.value.length + i}`
      })
    }
  } finally {
    loading.value = false
  }
}

function clearLogs() {
  logs.value = []
}

function clearFilters() {
  search.value = ''
  levelFilter.value = ['error', 'warn', 'info']
  sourceFilter.value = ''
  dateFilter.value = ''
}

function showLogDetail(log) {
  selectedLog.value = log
  showDetailModal.value = true
}

function copyLog() {
  if (selectedLog.value) {
    navigator.clipboard.writeText(JSON.stringify(selectedLog.value, null, 2))
  }
}

async function exportLogs() {
  try {
    const data = JSON.stringify(filteredLogs.value, null, 2)
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `logs-${new Date().toISOString().split('T')[0]}.json`
    a.click()
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Export failed:', error)
  }
}

// Lifecycle
onMounted(async () => {
  // Load initial logs
  for (let i = 0; i < 100; i++) {
    logs.value.push({
      time: new Date(Date.now() - i * 60000).toLocaleTimeString(),
      level: ['INFO', 'INFO', 'INFO', 'WARN', 'ERROR'][Math.floor(Math.random() * 5)],
      source: ['api', 'node', 'auth', 'plugin', 'system'][Math.floor(Math.random() * 5)],
      message: [
        'Request processed successfully',
        'Node connection established',
        'User authenticated',
        'High memory usage detected',
        'Connection timeout'
      ][Math.floor(Math.random() * 5)]
    })
  }
})

onUnmounted(() => {
  stopStream()
})
</script>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgb(17 24 39);
}

::-webkit-scrollbar-thumb {
  background: rgb(55 65 81);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgb(75 85 99);
}
</style>