<template>
  <div class="monitoring-page">
    <div class="page-header">
      <h1>Monitoring</h1>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="refresh">
          Refresh
        </button>
      </div>
    </div>

    <!-- Health Status Cards -->
    <div class="health-grid">
      <div
        v-for="check in healthChecks"
        :key="check.name"
        :class="['health-card', check.status]"
      >
        <div class="health-icon">
          {{ check.status === 'healthy' ? '✓' : check.status === 'degraded' ? '⚠' : '✕' }}
        </div>
        <div class="health-info">
          <span class="health-name">{{ check.name }}</span>
          <span class="health-status">{{ check.status }}</span>
          <span class="health-latency">{{ check.latency }}ms</span>
        </div>
      </div>
    </div>

    <!-- Metrics Section -->
    <div class="metrics-section">
      <h2>Key Metrics</h2>
      <div class="metrics-grid">
        <div class="metric-card">
          <span class="metric-value">{{ formatNumber(metrics.requestRate) }}</span>
          <span class="metric-label">Requests/sec</span>
        </div>
        <div class="metric-card">
          <span class="metric-value">{{ metrics.errorRate.toFixed(2) }}%</span>
          <span class="metric-label">Error Rate</span>
        </div>
        <div class="metric-card">
          <span class="metric-value">{{ metrics.avgLatency }}ms</span>
          <span class="metric-label">Avg Latency (P50)</span>
        </div>
        <div class="metric-card">
          <span class="metric-value">{{ metrics.p99Latency }}ms</span>
          <span class="metric-label">P99 Latency</span>
        </div>
      </div>
    </div>

    <!-- Active Alerts -->
    <div class="alerts-section">
      <h2>Active Alerts</h2>
      <div v-if="alerts.length === 0" class="no-alerts">
        No active alerts
      </div>
      <div v-else class="alerts-list">
        <div
          v-for="alert in alerts"
          :key="alert.id"
          :class="['alert-item', alert.severity]"
        >
          <div class="alert-header">
            <span class="alert-name">{{ alert.name }}</span>
            <span :class="['alert-severity', alert.severity]">{{ alert.severity }}</span>
          </div>
          <div class="alert-details">
            <span>{{ alert.metric }}: {{ alert.value }}</span>
            <span>Threshold: {{ alert.threshold }}</span>
            <span>Since: {{ formatTime(alert.startedAt) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Services Grid -->
    <div class="services-section">
      <h2>Services</h2>
      <div class="services-grid">
        <div
          v-for="service in services"
          :key="service.name"
          :class="['service-card', service.health]"
        >
          <div class="service-header">
            <span class="service-name">{{ service.name }}</span>
            <span :class="['service-health', service.health]">{{ service.health }}</span>
          </div>
          <div class="service-metrics">
            <div class="service-metric">
              <span class="metric-value">{{ formatNumber(service.requestRate) }}/s</span>
              <span class="metric-label">Requests</span>
            </div>
            <div class="service-metric">
              <span class="metric-value">{{ service.errorRate.toFixed(2) }}%</span>
              <span class="metric-label">Errors</span>
            </div>
            <div class="service-metric">
              <span class="metric-value">{{ service.latency }}ms</span>
              <span class="metric-label">Latency</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const healthChecks = ref([
  { name: 'API', status: 'healthy', latency: 12 },
  { name: 'Database', status: 'healthy', latency: 5 },
  { name: 'Cache', status: 'healthy', latency: 2 },
  { name: 'Queue', status: 'degraded', latency: 50 },
])

const metrics = ref({
  requestRate: 1250,
  errorRate: 0.15,
  avgLatency: 45,
  p99Latency: 120,
})

const alerts = ref([
  // { id: '1', name: 'High Memory Usage', severity: 'warning', metric: 'memory_percent', value: 85, threshold: 80, startedAt: '2026-03-23T10:00:00Z' },
])

const services = ref([
  { name: 'api-gateway', health: 'healthy', requestRate: 500, errorRate: 0.1, latency: 25 },
  { name: 'auth-service', health: 'healthy', requestRate: 200, errorRate: 0.0, latency: 15 },
  { name: 'task-runner', health: 'healthy', requestRate: 100, errorRate: 0.5, latency: 50 },
  { name: 'log-processor', health: 'healthy', requestRate: 450, errorRate: 0.0, latency: 10 },
])

function formatNumber(num) {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

function formatTime(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)
  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  return `${Math.floor(diff / 3600)}h ago`
}

function refresh() {
  // TODO: Implement refresh
}

onMounted(() => {
  // TODO: Fetch data from API
})
</script>

<style scoped>
.monitoring-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.health-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.health-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-left: 4px solid;
}

.health-card.healthy {
  border-left-color: #4caf50;
}

.health-card.degraded {
  border-left-color: #ff9800;
}

.health-card.unhealthy {
  border-left-color: #f44336;
}

.health-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
}

.health-card.healthy .health-icon {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.health-card.degraded .health-icon {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.health-card.unhealthy .health-icon {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.health-info {
  display: flex;
  flex-direction: column;
}

.health-name {
  font-weight: 600;
}

.health-status {
  color: var(--text-secondary);
  font-size: 12px;
  text-transform: capitalize;
}

.health-latency {
  color: var(--text-secondary);
  font-size: 12px;
}

.metrics-section,
.alerts-section,
.services-section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 24px;
}

.metrics-section h2,
.alerts-section h2,
.services-section h2 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.metric-card {
  text-align: center;
  padding: 16px;
  background: var(--bg);
  border-radius: 8px;
}

.metric-value {
  display: block;
  font-size: 24px;
  font-weight: 600;
}

.metric-label {
  display: block;
  color: var(--text-secondary);
  font-size: 12px;
  margin-top: 4px;
}

.no-alerts {
  color: var(--text-secondary);
  text-align: center;
  padding: 24px;
}

.alerts-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  padding: 12px;
  background: var(--bg);
  border-radius: 8px;
  border-left: 4px solid;
}

.alert-item.warning {
  border-left-color: #ff9800;
}

.alert-item.critical {
  border-left-color: #f44336;
}

.alert-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.alert-name {
  font-weight: 600;
}

.alert-severity {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  text-transform: uppercase;
}

.alert-severity.warning {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.alert-severity.critical {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.alert-details {
  display: flex;
  gap: 16px;
  color: var(--text-secondary);
  font-size: 12px;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.service-card {
  background: var(--bg);
  border-radius: 8px;
  padding: 16px;
  border-left: 4px solid;
}

.service-card.healthy {
  border-left-color: #4caf50;
}

.service-card.degraded {
  border-left-color: #ff9800;
}

.service-card.unhealthy {
  border-left-color: #f44336;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.service-name {
  font-weight: 600;
}

.service-health {
  font-size: 11px;
  text-transform: uppercase;
  padding: 2px 8px;
  border-radius: 4px;
}

.service-health.healthy {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.service-health.degraded {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.service-health.unhealthy {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.service-metrics {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.service-metric {
  text-align: center;
}

.service-metric .metric-value {
  font-size: 14px;
}

.service-metric .metric-label {
  font-size: 10px;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
}

.btn-secondary {
  background: transparent;
  color: var(--text);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--hover-bg);
}
</style>