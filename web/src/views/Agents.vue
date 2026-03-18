<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">Agents</h1>
        <p class="text-dark-400 text-sm mt-1">Connected devices and remote agents</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="refreshAgents"
          :disabled="loading"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <ArrowPathIcon :class="['w-5 h-5', { 'animate-spin': loading }]" />
          Refresh
        </button>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <PlusIcon class="w-5 h-5" />
          Register Agent
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Total Agents</p>
            <p class="text-2xl font-bold text-white">{{ agentStats.total }}</p>
          </div>
          <ComputerDesktopIcon class="w-8 h-8 text-primary-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Online</p>
            <p class="text-2xl font-bold text-green-400">{{ agentStats.online }}</p>
          </div>
          <SignalIcon class="w-8 h-8 text-green-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Offline</p>
            <p class="text-2xl font-bold text-red-400">{{ agentStats.offline }}</p>
          </div>
          <WifiIcon class="w-8 h-8 text-red-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Active Sessions</p>
            <p class="text-2xl font-bold text-blue-400">{{ agentStats.activeSessions }}</p>
          </div>
          <UserGroupIcon class="w-8 h-8 text-blue-400" />
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-4">
      <div class="flex-1 min-w-[200px]">
        <div class="relative">
          <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-400" />
          <input
            v-model="search"
            type="text"
            placeholder="Search agents..."
            class="w-full pl-10 pr-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-400 focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>
      </div>
      <select
        v-model="statusFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
      >
        <option value="">All Status</option>
        <option value="online">Online</option>
        <option value="offline">Offline</option>
        <option value="busy">Busy</option>
      </select>
      <select
        v-model="typeFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
      >
        <option value="">All Types</option>
        <option value="server">Server</option>
        <option value="workstation">Workstation</option>
        <option value="container">Container</option>
      </select>
    </div>

    <!-- Agents Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="agent in filteredAgents"
        :key="agent.id"
        class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden hover:border-dark-600 transition-colors"
      >
        <!-- Status Bar -->
        <div
          class="h-1"
          :class="{
            'bg-green-500': agent.status === 'online',
            'bg-red-500': agent.status === 'offline',
            'bg-yellow-500': agent.status === 'busy'
          }"
        ></div>

        <!-- Header -->
        <div class="p-4 border-b border-dark-700">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-lg flex items-center justify-center"
                :class="getAgentTypeColor(agent.type)"
              >
                <component :is="getAgentTypeIcon(agent.type)" class="w-5 h-5 text-white" />
              </div>
              <div>
                <h3 class="text-white font-medium">{{ agent.name }}</h3>
                <p class="text-dark-400 text-sm">{{ agent.host }}</p>
              </div>
            </div>
            <span
              class="px-2 py-1 text-xs rounded-full"
              :class="getStatusColor(agent.status)"
            >
              {{ agent.status }}
            </span>
          </div>
        </div>

        <!-- Stats -->
        <div class="p-4 grid grid-cols-3 gap-2 text-center">
          <div>
            <p class="text-dark-400 text-xs">CPU</p>
            <p class="text-white font-medium">{{ agent.cpu || 0 }}%</p>
          </div>
          <div>
            <p class="text-dark-400 text-xs">Memory</p>
            <p class="text-white font-medium">{{ agent.memory || 0 }}%</p>
          </div>
          <div>
            <p class="text-dark-400 text-xs">Disk</p>
            <p class="text-white font-medium">{{ agent.disk || 0 }}%</p>
          </div>
        </div>

        <!-- Info -->
        <div class="px-4 py-3 bg-dark-700/30 text-sm space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-dark-400">OS</span>
            <span class="text-white">{{ agent.os || 'Unknown' }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-dark-400">Uptime</span>
            <span class="text-white">{{ agent.uptime || '-' }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-dark-400">Version</span>
            <span class="text-white">{{ agent.version || '-' }}</span>
          </div>
        </div>

        <!-- Actions -->
        <div class="p-4 border-t border-dark-700 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <button
              @click="connectAgent(agent)"
              :disabled="agent.status === 'offline'"
              class="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white text-sm rounded-lg transition-colors flex items-center gap-1 disabled:opacity-50"
            >
              <TerminalIcon class="w-4 h-4" />
              Connect
            </button>
            <button
              @click="openTerminal(agent)"
              :disabled="agent.status === 'offline'"
              class="px-3 py-1.5 bg-dark-700 hover:bg-dark-600 text-white text-sm rounded-lg transition-colors disabled:opacity-50"
            >
              Terminal
            </button>
          </div>
          <div class="flex items-center gap-1">
            <button
              @click="viewStats(agent)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="Stats"
            >
              <ChartBarIcon class="w-4 h-4 text-dark-400" />
            </button>
            <button
              @click="restartAgent(agent)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="Restart"
            >
              <ArrowPathIcon class="w-4 h-4 text-dark-400" />
            </button>
            <button
              @click="confirmDelete(agent)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="Remove"
            >
              <TrashIcon class="w-4 h-4 text-red-400" />
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="filteredAgents.length === 0" class="col-span-3 py-12 text-center">
        <ComputerDesktopIcon class="w-12 h-12 text-dark-500 mx-auto mb-4" />
        <p class="text-dark-400 mb-2">No agents found</p>
        <button
          @click="openCreateModal"
          class="text-primary-400 hover:text-primary-300"
        >
          Register your first agent
        </button>
      </div>
    </div>

    <!-- Terminal Modal -->
    <Teleport to="body">
      <div v-if="showTerminal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeTerminal"></div>
        <div class="relative bg-dark-900 rounded-2xl border border-dark-700 w-full max-w-4xl h-[80vh] flex flex-col">
          <div class="px-4 py-3 border-b border-dark-700 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <TerminalIcon class="w-5 h-5 text-primary-400" />
              <span class="text-white font-medium">{{ terminalAgent?.name }}</span>
              <span class="text-dark-400 text-sm">{{ terminalAgent?.host }}</span>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="clearTerminal"
                class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white text-sm rounded-lg"
              >
                Clear
              </button>
              <button @click="closeTerminal" class="p-2 hover:bg-dark-700 rounded-lg">
                <XMarkIcon class="w-5 h-5 text-dark-400" />
              </button>
            </div>
          </div>

          <!-- Terminal Output -->
          <div
            ref="terminalOutput"
            class="flex-1 p-4 font-mono text-sm overflow-y-auto bg-black"
          >
            <div
              v-for="(line, index) in terminalLines"
              :key="index"
              class="py-0.5"
              :class="{
                'text-green-400': line.type === 'output',
                'text-yellow-400': line.type === 'command',
                'text-red-400': line.type === 'error'
              }"
            >
              <span v-if="line.type === 'command'" class="text-primary-400">$ </span>
              {{ line.content }}
            </div>
          </div>

          <!-- Terminal Input -->
          <div class="p-4 border-t border-dark-700 flex items-center gap-2">
            <span class="text-primary-400 font-mono">$</span>
            <input
              v-model="terminalInput"
              @keydown.enter="executeCommand"
              @keydown.up="historyUp"
              @keydown.down="historyDown"
              type="text"
              placeholder="Enter command..."
              class="flex-1 bg-transparent text-white font-mono focus:outline-none"
              autofocus
            />
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Agent Stats Modal -->
    <Teleport to="body">
      <div v-if="showStatsModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showStatsModal = false"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">Agent Statistics</h2>
            <button @click="showStatsModal = false" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <div v-if="statsAgent" class="p-6 space-y-6">
            <!-- CPU Chart -->
            <div>
              <h3 class="text-white font-medium mb-3">CPU Usage</h3>
              <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
                <div
                  class="h-full bg-blue-500"
                  :style="{ width: `${statsAgent.cpu || 0}%` }"
                ></div>
              </div>
              <p class="text-dark-400 text-sm mt-1">{{ statsAgent.cpu || 0 }}%</p>
            </div>

            <!-- Memory Chart -->
            <div>
              <h3 class="text-white font-medium mb-3">Memory Usage</h3>
              <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
                <div
                  class="h-full bg-green-500"
                  :style="{ width: `${statsAgent.memory || 0}%` }"
                ></div>
              </div>
              <p class="text-dark-400 text-sm mt-1">{{ statsAgent.memory || 0 }}% ({{ formatBytes((statsAgent.memory || 0) * 100 * 1024 * 1024) }} / {{ formatBytes(100 * 1024 * 1024) }})</p>
            </div>

            <!-- Disk Chart -->
            <div>
              <h3 class="text-white font-medium mb-3">Disk Usage</h3>
              <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
                <div
                  class="h-full bg-purple-500"
                  :style="{ width: `${statsAgent.disk || 0}%` }"
                ></div>
              </div>
              <p class="text-dark-400 text-sm mt-1">{{ statsAgent.disk || 0 }}%</p>
            </div>

            <!-- Network -->
            <div>
              <h3 class="text-white font-medium mb-3">Network I/O</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="p-3 bg-dark-700 rounded-lg">
                  <p class="text-dark-400 text-sm">Download</p>
                  <p class="text-white font-medium">{{ formatBytes(statsAgent.networkIn || 0) }}/s</p>
                </div>
                <div class="p-3 bg-dark-700 rounded-lg">
                  <p class="text-dark-400 text-sm">Upload</p>
                  <p class="text-white font-medium">{{ formatBytes(statsAgent.networkOut || 0) }}/s</p>
                </div>
              </div>
            </div>

            <!-- Processes -->
            <div>
              <h3 class="text-white font-medium mb-3">Running Processes</h3>
              <div class="space-y-2 max-h-40 overflow-y-auto">
                <div
                  v-for="proc in statsAgent.processes || []"
                  :key="proc.pid"
                  class="flex items-center justify-between p-2 bg-dark-700 rounded"
                >
                  <div class="flex items-center gap-2">
                    <span class="text-dark-400 text-sm">{{ proc.pid }}</span>
                    <span class="text-white text-sm">{{ proc.name }}</span>
                  </div>
                  <span class="text-dark-400 text-sm">{{ proc.cpu }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Create Agent Modal -->
    <Teleport to="body">
      <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showCreateModal = false"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-md p-6">
          <h2 class="text-xl font-semibold text-white mb-4">Register New Agent</h2>
          <form @submit.prevent="handleCreateAgent" class="space-y-4">
            <div>
              <label class="block text-sm text-dark-400 mb-1">Agent Name</label>
              <input
                v-model="newAgent.name"
                type="text"
                required
                placeholder="e.g., web-server-01"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
            <div>
              <label class="block text-sm text-dark-400 mb-1">Host / IP</label>
              <input
                v-model="newAgent.host"
                type="text"
                required
                placeholder="e.g., 192.168.1.100"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
            <div>
              <label class="block text-sm text-dark-400 mb-1">Type</label>
              <select
                v-model="newAgent.type"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                <option value="server">Server</option>
                <option value="workstation">Workstation</option>
                <option value="container">Container</option>
              </select>
            </div>
            <div>
              <label class="block text-sm text-dark-400 mb-1">Agent Key (optional)</label>
              <input
                v-model="newAgent.key"
                type="text"
                placeholder="Auto-generated if empty"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>

            <div class="flex items-center justify-end gap-3 pt-4">
              <button
                type="button"
                @click="showCreateModal = false"
                class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors"
              >
                Register
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Toast Notifications -->
    <Teleport to="body">
      <div class="fixed bottom-4 right-4 z-50 space-y-2">
        <TransitionGroup name="toast">
          <div
            v-for="toast in toasts"
            :key="toast.id"
            class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 min-w-[300px]"
            :class="{
              'bg-green-900/90 border border-green-700': toast.type === 'success',
              'bg-red-900/90 border border-red-700': toast.type === 'error',
              'bg-yellow-900/90 border border-yellow-700': toast.type === 'warning'
            }"
          >
            <CheckCircleIcon v-if="toast.type === 'success'" class="w-5 h-5 text-green-400" />
            <XCircleIcon v-else-if="toast.type === 'error'" class="w-5 h-5 text-red-400" />
            <ExclamationTriangleIcon v-else class="w-5 h-5 text-yellow-400" />
            <span class="text-white flex-1">{{ toast.message }}</span>
            <button @click="removeToast(toast.id)" class="text-dark-400 hover:text-white">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </TransitionGroup>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { agentsApi } from '@/api'
import {
  ComputerDesktopIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon,
  TrashIcon,
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
  TerminalIcon,
  ChartBarIcon,
  SignalIcon,
  WifiIcon,
  UserGroupIcon,
  ExclamationTriangleIcon,
  ServerIcon,
  CubeIcon
} from '@heroicons/vue/24/outline'

// State
const search = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const showTerminal = ref(false)
const showStatsModal = ref(false)
const showCreateModal = ref(false)
const terminalAgent = ref(null)
const statsAgent = ref(null)
const terminalInput = ref('')
const terminalLines = ref([])
const terminalOutput = ref(null)
const commandHistory = ref([])
const historyIndex = ref(-1)
const loading = ref(false)
const toasts = ref([])
const agents = ref([])

const newAgent = ref({
  name: '',
  host: '',
  type: 'server',
  key: ''
})

// Computed
const filteredAgents = computed(() => {
  let result = agents.value

  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(a =>
      a.name.toLowerCase().includes(searchLower) ||
      a.host.includes(search.value)
    )
  }

  if (statusFilter.value) {
    result = result.filter(a => a.status === statusFilter.value)
  }

  if (typeFilter.value) {
    result = result.filter(a => a.type === typeFilter.value)
  }

  return result
})

const agentStats = computed(() => ({
  total: agents.value.length,
  online: agents.value.filter(a => a.status === 'online').length,
  offline: agents.value.filter(a => a.status === 'offline').length,
  activeSessions: agents.value.filter(a => a.status === 'online').length
}))

// Methods
function getAgentTypeIcon(type) {
  const icons = {
    server: ServerIcon,
    workstation: ComputerDesktopIcon,
    container: CubeIcon
  }
  return icons[type] || ComputerDesktopIcon
}

function getAgentTypeColor(type) {
  const colors = {
    server: 'bg-blue-600',
    workstation: 'bg-purple-600',
    container: 'bg-green-600'
  }
  return colors[type] || 'bg-gray-600'
}

function getStatusColor(status) {
  const colors = {
    online: 'bg-green-900/30 text-green-400',
    offline: 'bg-red-900/30 text-red-400',
    busy: 'bg-yellow-900/30 text-yellow-400'
  }
  return colors[status] || 'bg-gray-900/30 text-gray-400'
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function refreshAgents() {
  loading.value = true
  try {
    const response = await agentsApi.list()
    agents.value = response.data.data || response.data || []
    showToast('Agents refreshed', 'success')
  } catch (error) {
    // Use mock data
    agents.value = [
      { id: 1, name: 'web-01', host: '10.0.0.101', status: 'online', type: 'server', os: 'Ubuntu 22.04', uptime: '2d 4h', version: '1.0.0', cpu: 45, memory: 62, disk: 35 },
      { id: 2, name: 'db-01', host: '10.0.0.102', status: 'online', type: 'server', os: 'Debian 12', uptime: '5d 12h', version: '1.0.0', cpu: 78, memory: 85, disk: 60 },
      { id: 3, name: 'cache-01', host: '10.0.0.103', status: 'offline', type: 'server', os: 'CentOS 8', uptime: '-', version: '1.0.0', cpu: 0, memory: 0, disk: 45 },
      { id: 4, name: 'dev-ws', host: '10.0.0.200', status: 'online', type: 'workstation', os: 'Windows 11', uptime: '8h 30m', version: '1.0.1', cpu: 25, memory: 45, disk: 70 }
    ]
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  newAgent.value = { name: '', host: '', type: 'server', key: '' }
  showCreateModal.value = true
}

async function handleCreateAgent() {
  try {
    await agentsApi.create(newAgent.value)
    showToast('Agent registered successfully', 'success')
    showCreateModal.value = false
    await refreshAgents()
  } catch (error) {
    showToast(error.message || 'Failed to register agent', 'error')
  }
}

function connectAgent(agent) {
  openTerminal(agent)
}

function openTerminal(agent) {
  if (agent.status === 'offline') return
  terminalAgent.value = agent
  terminalLines.value = [
    { type: 'output', content: `Connected to ${agent.name} (${agent.host})` },
    { type: 'output', content: `OS: ${agent.os || 'Unknown'}` },
    { type: 'output', content: `Type 'help' for available commands` }
  ]
  showTerminal.value = true
  nextTick(() => {
    terminalOutput.value?.scrollTo(0, terminalOutput.value.scrollHeight)
  })
}

function closeTerminal() {
  showTerminal.value = false
  terminalAgent.value = null
  terminalLines.value = []
  commandHistory.value = []
  historyIndex.value = -1
}

function clearTerminal() {
  terminalLines.value = []
}

function executeCommand() {
  const cmd = terminalInput.value.trim()
  if (!cmd) return

  terminalLines.value.push({ type: 'command', content: cmd })
  commandHistory.value.unshift(cmd)
  historyIndex.value = -1

  // Simulate command execution
  setTimeout(() => {
    const output = simulateCommand(cmd)
    terminalLines.value.push(...output)
    nextTick(() => {
      terminalOutput.value?.scrollTo(0, terminalOutput.value.scrollHeight)
    })
  }, 100)

  terminalInput.value = ''
}

function simulateCommand(cmd) {
  const parts = cmd.split(' ')
  const command = parts[0].toLowerCase()

  switch (command) {
    case 'help':
      return [
        { type: 'output', content: 'Available commands:' },
        { type: 'output', content: '  help     - Show this help' },
        { type: 'output', content: '  ls       - List files' },
        { type: 'output', content: '  ps       - Show processes' },
        { type: 'output', content: '  top      - System stats' },
        { type: 'output', content: '  netstat  - Network connections' },
        { type: 'output', content: '  exit     - Close terminal' }
      ]
    case 'ls':
      return [
        { type: 'output', content: 'total 24' },
        { type: 'output', content: 'drwxr-xr-x  5 root root 4096 Mar 15 10:30 .' },
        { type: 'output', content: 'drwxr-xr-x 20 root root 4096 Mar 10 08:00 ..' },
        { type: 'output', content: '-rw-r--r--  1 root root  123 Mar 15 10:30 config.yml' },
        { type: 'output', content: 'drwxr-xr-x  2 root root 4096 Mar 14 16:20 logs' }
      ]
    case 'ps':
      return [
        { type: 'output', content: 'PID   USER     CPU  MEM  COMMAND' },
        { type: 'output', content: '1     root     0.0  0.1  /sbin/init' },
        { type: 'output', content: '842   root     1.2  2.5  /usr/bin/anixops-agent' },
        { type: 'output', content: '1203  root     0.5  1.2  /usr/bin/node' }
      ]
    case 'top':
      return [
        { type: 'output', content: `CPU: ${terminalAgent.value?.cpu || 0}% | Memory: ${terminalAgent.value?.memory || 0}%` }
      ]
    case 'exit':
      closeTerminal()
      return []
    default:
      return [{ type: 'error', content: `Command not found: ${command}` }]
  }
}

function historyUp() {
  if (commandHistory.value.length === 0) return
  if (historyIndex.value < commandHistory.value.length - 1) {
    historyIndex.value++
    terminalInput.value = commandHistory.value[historyIndex.value]
  }
}

function historyDown() {
  if (historyIndex.value > 0) {
    historyIndex.value--
    terminalInput.value = commandHistory.value[historyIndex.value]
  } else if (historyIndex.value === 0) {
    historyIndex.value = -1
    terminalInput.value = ''
  }
}

function viewStats(agent) {
  statsAgent.value = {
    ...agent,
    networkIn: Math.random() * 1000000,
    networkOut: Math.random() * 500000,
    processes: [
      { pid: 1, name: 'init', cpu: 0.1 },
      { pid: 842, name: 'anixops-agent', cpu: 1.5 },
      { pid: 1203, name: 'node', cpu: 2.3 }
    ]
  }
  showStatsModal.value = true
}

async function restartAgent(agent) {
  try {
    await agentsApi.restart(agent.id)
    showToast(`Restarting ${agent.name}...`, 'success')
  } catch (error) {
    showToast(error.message || 'Failed to restart agent', 'error')
  }
}

function confirmDelete(agent) {
  if (confirm(`Are you sure you want to remove ${agent.name}?`)) {
    deleteAgent(agent)
  }
}

async function deleteAgent(agent) {
  try {
    await agentsApi.delete(agent.id)
    showToast('Agent removed', 'success')
    await refreshAgents()
  } catch (error) {
    showToast(error.message || 'Failed to remove agent', 'error')
  }
}

function showToast(message, type = 'success') {
  const id = Date.now()
  toasts.value.push({ id, message, type })
  setTimeout(() => removeToast(id), 5000)
}

function removeToast(id) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

// Lifecycle
onMounted(async () => {
  await refreshAgents()
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>