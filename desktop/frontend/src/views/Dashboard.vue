<template>
  <div class="space-y-6">
    <!-- Offline Indicator -->
    <OfflineIndicator />

    <!-- Stats Grid -->
    <div class="grid grid-cols-4 gap-4">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-slate-800 rounded-xl p-4 border border-slate-700"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-slate-400 text-sm">{{ stat.label }}</p>
            <p class="text-2xl font-bold text-white mt-1">
              <span v-if="loading" class="animate-pulse">...</span>
              <span v-else>{{ stat.value }}</span>
            </p>
          </div>
          <div :class="stat.iconBg" class="w-10 h-10 rounded-lg flex items-center justify-center">
            <component :is="stat.icon" class="w-5 h-5 text-white" />
          </div>
        </div>
        <div class="mt-3 flex items-center gap-2">
          <span :class="stat.changeColor" class="text-sm font-medium">
            {{ stat.change }}
          </span>
          <span class="text-slate-500 text-sm">vs last month</span>
        </div>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-2 gap-6">
      <!-- Traffic Chart -->
      <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-white font-semibold">Traffic Overview</h3>
          <select
            v-model="trafficPeriod"
            @change="fetchTraffic"
            class="px-3 py-1 bg-slate-700 border border-slate-600 rounded-lg text-white text-sm"
          >
            <option value="24h">Last 24 hours</option>
            <option value="7d">Last 7 days</option>
            <option value="30d">Last 30 days</option>
          </select>
        </div>
        <div v-if="trafficData.length > 0" class="h-64 flex items-end gap-2">
          <div
            v-for="(value, index) in trafficData"
            :key="index"
            class="flex-1 bg-blue-600 rounded-t transition-all hover:bg-blue-500"
            :style="{ height: `${value}%` }"
          ></div>
        </div>
        <div v-else class="h-64 flex items-center justify-center text-slate-500">
          <span v-if="loading">Loading...</span>
          <span v-else>No data available</span>
        </div>
        <div class="mt-4 flex justify-between text-slate-500 text-xs">
          <span>00:00</span>
          <span>06:00</span>
          <span>12:00</span>
          <span>18:00</span>
          <span>24:00</span>
        </div>
      </div>

      <!-- Node Status -->
      <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-white font-semibold">Node Status</h3>
          <button
            @click="refreshNodes"
            :disabled="nodesLoading"
            class="p-2 hover:bg-slate-700 rounded-lg transition-colors"
          >
            <svg
              :class="['w-4 h-4 text-slate-400', { 'animate-spin': nodesLoading }]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        </div>
        <div v-if="nodes.length > 0" class="space-y-4">
          <div v-for="node in nodes.slice(0, 5)" :key="node.id" class="flex items-center gap-4">
            <div
              class="w-2 h-2 rounded-full"
              :class="node.status === 'online' ? 'bg-green-500' : 'bg-red-500'"
            ></div>
            <div class="flex-1">
              <div class="flex items-center justify-between">
                <span class="text-white">{{ node.name }}</span>
                <span class="text-slate-400 text-sm">{{ node.users || 0 }} users</span>
              </div>
              <div class="mt-1 h-1.5 bg-slate-700 rounded-full overflow-hidden">
                <div
                  class="h-full bg-blue-500 rounded-full"
                  :style="{ width: `${node.load || 50}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="py-8 text-center text-slate-500">
          <span v-if="nodesLoading">Loading nodes...</span>
          <span v-else>No nodes available</span>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="grid grid-cols-3 gap-6">
      <div class="col-span-2 bg-slate-800 rounded-xl border border-slate-700 p-6">
        <h3 class="text-white font-semibold mb-4">Recent Activity</h3>
        <div v-if="activities.length > 0" class="space-y-4">
          <div
            v-for="activity in activities"
            :key="activity.id"
            class="flex items-start gap-4 p-3 bg-slate-700/50 rounded-lg"
          >
            <div :class="activity.iconBg" class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0">
              <component :is="activity.icon" class="w-4 h-4 text-white" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-white text-sm">{{ activity.message }}</p>
              <p class="text-slate-500 text-xs mt-1">{{ activity.time }}</p>
            </div>
          </div>
        </div>
        <div v-else class="py-8 text-center text-slate-500">
          No recent activity
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
        <h3 class="text-white font-semibold mb-4">Quick Actions</h3>
        <div class="space-y-2">
          <button
            v-for="action in quickActions"
            :key="action.label"
            @click="handleAction(action.action)"
            class="w-full flex items-center gap-3 px-4 py-3 bg-slate-700 hover:bg-slate-600 rounded-lg transition-colors"
          >
            <component :is="action.icon" class="w-5 h-5 text-blue-400" />
            <span class="text-white text-sm">{{ action.label }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { h, ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNodesStore } from '../stores/nodes'
import { useOfflineStore } from '../stores/offline'
import OfflineIndicator from '../components/OfflineIndicator.vue'

const router = useRouter()
const nodesStore = useNodesStore()
const offlineStore = useOfflineStore()

// State
const loading = ref(false)
const nodesLoading = ref(false)
const trafficPeriod = ref('24h')
const trafficData = ref([])
const activities = ref([])

// Computed
const nodes = computed(() => nodesStore.nodes)

const stats = computed(() => [
  {
    label: 'Total Nodes',
    value: nodesStore.nodeStats.total.toString(),
    change: '+3',
    changeColor: 'text-green-500',
    iconBg: 'bg-blue-600',
    icon: ServerIcon
  },
  {
    label: 'Active Users',
    value: '1,284',
    change: '+12%',
    changeColor: 'text-green-500',
    iconBg: 'bg-green-600',
    icon: UsersIcon
  },
  {
    label: 'Traffic Today',
    value: '2.4 TB',
    change: '+8%',
    changeColor: 'text-green-500',
    iconBg: 'bg-purple-600',
    icon: ChartIcon
  },
  {
    label: 'Uptime',
    value: '99.9%',
    change: 'Stable',
    changeColor: 'text-blue-500',
    iconBg: 'bg-orange-600',
    icon: CheckIcon
  }
])

const quickActions = computed(() => [
  { label: 'Deploy New Node', action: 'deploy', icon: PlusIcon },
  { label: 'Run Playbook', action: 'playbook', icon: PlayIcon },
  { label: 'Add User', action: 'user', icon: UsersIcon },
  { label: 'View Logs', action: 'logs', icon: ChartIcon }
])

// Icons
const ServerIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01' })
])

const UsersIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' })
])

const ChartIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' })
])

const CheckIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' })
])

const PlusIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4v16m8-8H4' })
])

const RefreshIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15' })
])

const PlayIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z' }),
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
])

// Methods
async function refreshNodes() {
  nodesLoading.value = true
  try {
    await nodesStore.fetchNodes()
    // Cache nodes for offline use
    await offlineStore.cacheNodes(nodes.value)
  } catch (error) {
    console.error('Failed to fetch nodes:', error)
    // Try to load from cache if offline
    if (!offlineStore.isOnline.value) {
      const cachedNodes = await offlineStore.getCachedNodes()
      nodesStore.setNodes(cachedNodes)
    }
  } finally {
    nodesLoading.value = false
  }
}

function fetchTraffic() {
  // Simulated traffic data
  trafficData.value = Array.from({ length: 12 }, () => Math.floor(Math.random() * 80) + 20)
}

function handleAction(action) {
  const routes = {
    deploy: '/nodes',
    playbook: '/playbooks',
    user: '/users',
    logs: '/logs'
  }
  if (routes[action]) {
    router.push(routes[action])
  }
}

// Lifecycle
onMounted(async () => {
  loading.value = true

  try {
    // Initialize offline store
    await offlineStore.initialize()

    // Load data
    await Promise.all([
      refreshNodes(),
      fetchTraffic()
    ])

    // Set up activity feed (simulated)
    activities.value = [
      { id: 1, message: 'Node tokyo-01 deployed successfully', time: '2 minutes ago', icon: CheckIcon, iconBg: 'bg-green-600' },
      { id: 2, message: 'User john@example.com registered', time: '15 minutes ago', icon: UsersIcon, iconBg: 'bg-blue-600' },
      { id: 3, message: 'Playbook deploy_node.yml completed', time: '1 hour ago', icon: PlayIcon, iconBg: 'bg-purple-600' },
      { id: 4, message: 'Node singapore-01 restarted', time: '2 hours ago', icon: RefreshIcon, iconBg: 'bg-yellow-600' }
    ]
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  offlineStore.cleanup()
})
</script>