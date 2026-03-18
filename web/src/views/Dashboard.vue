<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">{{ t('dashboard.title') }}</h1>
        <p class="text-dark-400 text-sm mt-1">{{ t('dashboard.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 text-sm">
          <span
            class="w-2 h-2 rounded-full"
            :class="connected ? 'bg-green-400 animate-pulse' : 'bg-red-400'"
          />
          <span class="text-dark-400">{{ connected ? t('dashboard.connected') : t('dashboard.disconnected') }}</span>
        </div>
        <button
          @click="refresh"
          :disabled="loading"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <ArrowPathIcon :class="{ 'animate-spin': loading }" class="w-4 h-4" />
          {{ t('common.refresh') }}
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard
        v-for="stat in statCards"
        :key="stat.title"
        :title="stat.title"
        :value="stat.value"
        :change="stat.change"
        :icon="stat.icon"
        :color="stat.color"
        :loading="loading"
      />
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Traffic Chart -->
      <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-white">{{ t('dashboard.trafficOverview') }}</h2>
          <select
            v-model="trafficPeriod"
            @change="fetchTraffic(trafficPeriod)"
            class="bg-dark-700 text-white text-sm rounded-lg px-3 py-1.5 border border-dark-600"
          >
            <option value="24h">24 Hours</option>
            <option value="7d">7 Days</option>
            <option value="30d">30 Days</option>
          </select>
        </div>
        <LineChart
          v-if="trafficData.length > 0"
          :data="trafficData"
          :height="250"
        />
        <div v-else class="h-[250px] flex items-center justify-center text-dark-400">
          <span v-if="loading">{{ t('common.loading') }}</span>
          <span v-else>{{ t('common.noData') }}</span>
        </div>
      </div>

      <!-- System Metrics -->
      <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
        <h2 class="text-lg font-semibold text-white mb-4">{{ t('dashboard.systemMetrics') }}</h2>
        <div class="space-y-4">
          <MetricBar
            :label="t('dashboard.cpuUsage')"
            :value="systemMetrics.cpu"
            :max="100"
            color="bg-blue-500"
          />
          <MetricBar
            :label="t('dashboard.memory')"
            :value="systemMetrics.memory"
            :max="100"
            color="bg-green-500"
          />
          <MetricBar
            :label="t('dashboard.disk')"
            :value="systemMetrics.disk"
            :max="100"
            color="bg-purple-500"
          />
          <MetricBar
            :label="t('dashboard.network')"
            :value="systemMetrics.network"
            :max="100"
            color="bg-orange-500"
          />
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Infrastructure Status -->
      <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
        <h2 class="text-lg font-semibold text-white mb-4">{{ t('dashboard.infrastructure') }}</h2>
        <div class="space-y-3">
          <ServiceItem
            v-for="service in services"
            :key="service.name"
            :name="service.name"
            :status="service.status"
            :count="service.count"
          />
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="lg:col-span-2 bg-dark-800 rounded-xl p-6 border border-dark-700">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-white">{{ t('dashboard.recentActivity') }}</h2>
          <router-link to="/logs" class="text-primary-400 text-sm hover:underline">
            {{ t('common.details') }} →
          </router-link>
        </div>
        <div class="space-y-3 max-h-[400px] overflow-y-auto">
          <ActivityItem
            v-for="activity in activities"
            :key="activity.id"
            :activity="activity"
          />
        </div>
      </div>
    </div>

    <!-- Quick Actions & Alerts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Quick Actions -->
      <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
        <h2 class="text-lg font-semibold text-white mb-4">{{ t('dashboard.quickActions') }}</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <QuickAction
            v-for="action in quickActions"
            :key="action.name"
            :name="action.name"
            :icon="action.icon"
            :color="action.color"
            @click="handleAction(action.action)"
          />
        </div>
      </div>

      <!-- Active Alerts -->
      <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-white">{{ t('dashboard.activeAlerts') }}</h2>
          <span
            v-if="alerts.length > 0"
            class="px-2 py-1 bg-red-500/20 text-red-400 text-xs rounded-full"
          >
            {{ alerts.length }}
          </span>
        </div>
        <div v-if="alerts.length > 0" class="space-y-3">
          <AlertItem
            v-for="alert in alerts.slice(0, 5)"
            :key="alert.id"
            :alert="alert"
            @dismiss="dismissAlert(alert.id)"
          />
        </div>
        <div v-else class="py-8 text-center text-dark-400">
          <CheckCircleIcon class="w-12 h-12 mx-auto mb-2 text-green-400" />
          <p>{{ t('dashboard.noAlerts') }}</p>
        </div>
      </div>
    </div>

    <!-- Plugins Status -->
    <div class="bg-dark-800 rounded-xl p-6 border border-dark-700">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-white">{{ t('plugins.title') }}</h2>
        <router-link to="/plugins" class="text-primary-400 text-sm hover:underline">
          {{ t('common.details') }} →
        </router-link>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <PluginCard
          v-for="plugin in plugins"
          :key="plugin.name"
          :plugin="plugin"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  ArrowPathIcon,
  CheckCircleIcon,
  ServerIcon,
  UsersIcon,
  ComputerDesktopIcon,
  ChartBarIcon,
  DocumentTextIcon
} from '@heroicons/vue/24/outline'

import { useDashboardStore } from '@/stores/dashboard'
import { usePluginsStore } from '@/stores/plugins'
import { useWebSocket, WSEvents } from '@/services/websocket'

import StatCard from '@/components/dashboard/StatCard.vue'
import MetricBar from '@/components/dashboard/MetricBar.vue'
import ServiceItem from '@/components/dashboard/ServiceItem.vue'
import ActivityItem from '@/components/dashboard/ActivityItem.vue'
import QuickAction from '@/components/dashboard/QuickAction.vue'
import AlertItem from '@/components/dashboard/AlertItem.vue'
import PluginCard from '@/components/dashboard/PluginCard.vue'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()
const router = useRouter()
const dashboardStore = useDashboardStore()
const pluginsStore = usePluginsStore()
const { connected, subscribe, connect } = useWebSocket()

// State
const loading = ref(false)
const trafficPeriod = ref('24h')
const systemMetrics = ref({
  cpu: 45,
  memory: 67,
  disk: 78,
  network: 23
})

// Computed
const stats = computed(() => dashboardStore.stats)
const activities = computed(() => dashboardStore.activities)
const alerts = computed(() => dashboardStore.alerts)
const trafficData = computed(() => dashboardStore.traffic)
const plugins = computed(() => pluginsStore.plugins.slice(0, 4))

const statCards = computed(() => [
  {
    title: t('dashboard.totalNodes'),
    value: stats.value.nodes.total.toString(),
    change: 0,
    icon: ServerIcon,
    color: 'bg-blue-600'
  },
  {
    title: t('dashboard.activeUsers'),
    value: stats.value.users.active.toString(),
    change: 12,
    icon: UsersIcon,
    color: 'bg-green-600'
  },
  {
    title: t('dashboard.onlineAgents'),
    value: stats.value.agents.online.toString(),
    change: 0,
    icon: ComputerDesktopIcon,
    color: 'bg-purple-600'
  },
  {
    title: t('dashboard.trafficToday'),
    value: formatBytes(stats.value.traffic.today),
    change: -5,
    icon: ChartBarIcon,
    color: 'bg-orange-600'
  }
])

const services = computed(() => [
  { name: 'Panel', status: 'online', count: t('common.running') },
  { name: t('navigation.nodes'), status: stats.value.nodes.offline > 0 ? 'degraded' : 'online', count: `${stats.value.nodes.online} / ${stats.value.nodes.total}` },
  { name: t('navigation.agents'), status: stats.value.agents.online > 0 ? 'online' : 'offline', count: `${stats.value.agents.online} / ${stats.value.agents.total}` },
  { name: t('navigation.plugins'), status: stats.value.plugins.active > 0 ? 'online' : 'offline', count: `${stats.value.plugins.active} / ${stats.value.plugins.total}` }
])

const quickActions = computed(() => [
  { name: t('dashboard.deployNode'), action: 'deploy', icon: ServerIcon, color: 'bg-blue-500' },
  { name: t('dashboard.manageUsers'), action: 'users', icon: UsersIcon, color: 'bg-green-500' },
  { name: t('dashboard.runPlaybook'), action: 'playbook', icon: DocumentTextIcon, color: 'bg-purple-500' },
  { name: t('dashboard.viewLogs'), action: 'logs', icon: ChartBarIcon, color: 'bg-orange-500' }
])

// Methods
async function refresh() {
  loading.value = true
  try {
    await Promise.all([
      dashboardStore.fetchAll(),
      pluginsStore.fetchPlugins()
    ])
  } finally {
    loading.value = false
  }
}

function fetchTraffic(period) {
  dashboardStore.fetchTraffic(period)
}

function handleAction(action) {
  const routes = {
    deploy: '/nodes',
    users: '/users',
    playbook: '/playbooks',
    logs: '/logs'
  }
  router.push(routes[action])
}

function dismissAlert(id) {
  dashboardStore.removeAlert(id)
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// WebSocket subscriptions
let unsubscribers = []

onMounted(async () => {
  await refresh()

  // Connect WebSocket
  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/v1/ws`
  connect(wsUrl)

  // Subscribe to events
  unsubscribers.push(
    subscribe(WSEvents.SYSTEM_METRICS, (data) => {
      systemMetrics.value = data
    })
  )

  unsubscribers.push(
    subscribe(WSEvents.ACTIVITY, (activity) => {
      dashboardStore.addActivity(activity)
    })
  )

  unsubscribers.push(
    subscribe(WSEvents.ALERT_CREATED, (alert) => {
      dashboardStore.addAlert(alert)
    })
  )

  unsubscribers.push(
    subscribe(WSEvents.NODE_STATUS, (data) => {
      // Update node stats in real-time
    })
  )
})

onUnmounted(() => {
  unsubscribers.forEach(unsub => unsub())
})
</script>