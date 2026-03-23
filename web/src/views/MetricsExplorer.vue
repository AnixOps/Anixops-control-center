<template>
  <div class="metrics-page">
    <div class="page-header">
      <h1>Metrics Explorer</h1>
      <div class="header-actions">
        <select v-model="selectedService" class="filter-select">
          <option value="">All Services</option>
          <option v-for="s in services" :key="s" :value="s">{{ s }}</option>
        </select>
        <select v-model="timeRange" class="filter-select">
          <option value="5m">5 minutes</option>
          <option value="15m">15 minutes</option>
          <option value="1h">1 hour</option>
          <option value="6h">6 hours</option>
          <option value="24h">24 hours</option>
        </select>
        <button class="btn btn-secondary" @click="refresh">Refresh</button>
      </div>
    </div>

    <!-- Key Metrics -->
    <div class="metrics-grid">
      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-name">Request Rate</span>
          <span class="metric-trend up">+12%</span>
        </div>
        <span class="metric-value">{{ formatNumber(metrics.requestRate) }}</span>
        <span class="metric-unit">req/s</span>
      </div>
      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-name">Error Rate</span>
          <span class="metric-trend down">-5%</span>
        </div>
        <span class="metric-value">{{ metrics.errorRate.toFixed(2) }}%</span>
        <span class="metric-unit">errors</span>
      </div>
      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-name">P50 Latency</span>
          <span class="metric-trend neutral">0%</span>
        </div>
        <span class="metric-value">{{ metrics.p50Latency }}</span>
        <span class="metric-unit">ms</span>
      </div>
      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-name">P99 Latency</span>
          <span class="metric-trend up">+8%</span>
        </div>
        <span class="metric-value">{{ metrics.p99Latency }}</span>
        <span class="metric-unit">ms</span>
      </div>
    </div>

    <!-- Service Metrics -->
    <div class="services-section">
      <h2>Service Metrics</h2>
      <div class="services-table">
        <div class="table-header">
          <span class="col-name">Service</span>
          <span class="col-requests">Requests/s</span>
          <span class="col-errors">Errors</span>
          <span class="col-latency">Latency</span>
          <span class="col-cpu">CPU</span>
          <span class="col-memory">Memory</span>
        </div>
        <div v-for="service in serviceMetrics" :key="service.name" class="table-row">
          <span class="col-name">{{ service.name }}</span>
          <span class="col-requests">{{ formatNumber(service.requests) }}</span>
          <span :class="['col-errors', { warning: service.errors > 1 }]">{{ service.errors.toFixed(2) }}%</span>
          <span class="col-latency">{{ service.latency }}ms</span>
          <span class="col-cpu">{{ service.cpu }}%</span>
          <span class="col-memory">{{ service.memory }}%</span>
        </div>
      </div>
    </div>

    <!-- Time Series -->
    <div class="timeseries-section">
      <h2>Request Rate Over Time</h2>
      <div class="chart-placeholder">
        <div class="chart-area">
          <div v-for="(point, i) in timeSeriesData" :key="i" class="chart-bar" :style="{ height: point + '%' }"></div>
        </div>
        <div class="chart-labels">
          <span v-for="(label, i) in timeLabels" :key="i">{{ label }}</span>
        </div>
      </div>
    </div>

    <!-- Top Endpoints -->
    <div class="endpoints-section">
      <h2>Top Endpoints</h2>
      <div class="endpoints-list">
        <div v-for="ep in topEndpoints" :key="ep.path" class="endpoint-item">
          <div class="endpoint-info">
            <span class="endpoint-method" :class="ep.method.toLowerCase()">{{ ep.method }}</span>
            <span class="endpoint-path">{{ ep.path }}</span>
          </div>
          <div class="endpoint-stats">
            <span class="endpoint-requests">{{ formatNumber(ep.requests) }} req</span>
            <span class="endpoint-latency">{{ ep.avgLatency }}ms</span>
          </div>
          <div class="endpoint-bar">
            <div class="bar-fill" :style="{ width: ep.percentage + '%' }"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const selectedService = ref('')
const timeRange = ref('1h')

const services = ref(['api-gateway', 'auth-service', 'task-runner', 'log-processor'])

const metrics = ref({
  requestRate: 1250,
  errorRate: 0.15,
  p50Latency: 45,
  p99Latency: 120
})

const serviceMetrics = ref([
  { name: 'api-gateway', requests: 500, errors: 0.1, latency: 25, cpu: 35, memory: 45 },
  { name: 'auth-service', requests: 200, errors: 0.0, latency: 15, cpu: 20, memory: 30 },
  { name: 'task-runner', requests: 100, errors: 0.5, latency: 50, cpu: 45, memory: 55 },
  { name: 'log-processor', requests: 450, errors: 0.0, latency: 10, cpu: 25, memory: 40 }
])

const timeSeriesData = ref([65, 72, 68, 80, 75, 82, 78, 90, 85, 88, 92, 95, 88, 85, 80])
const timeLabels = ref(['10:00', '10:05', '10:10', '10:15', '10:20', '10:25', '10:30', '10:35', '10:40', '10:45', '10:50', '10:55', '11:00', '11:05', '11:10'])

const topEndpoints = ref([
  { method: 'GET', path: '/api/users', requests: 12500, avgLatency: 25, percentage: 100 },
  { method: 'POST', path: '/api/tasks', requests: 8500, avgLatency: 45, percentage: 68 },
  { method: 'GET', path: '/api/nodes', requests: 6200, avgLatency: 15, percentage: 50 },
  { method: 'GET', path: '/api/playbooks', requests: 4800, avgLatency: 20, percentage: 38 },
  { method: 'POST', path: '/api/auth/login', requests: 3500, avgLatency: 35, percentage: 28 }
])

function formatNumber(num) {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

function refresh() {
  // TODO: Implement refresh
}

onMounted(() => {
  // TODO: Fetch data from API
})
</script>

<style scoped>
.metrics-page {
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
  align-items: center;
}

.filter-select {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.metric-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.metric-name {
  color: var(--text-secondary);
  font-size: 12px;
}

.metric-trend {
  font-size: 11px;
  font-weight: 600;
}

.metric-trend.up {
  color: #f44336;
}

.metric-trend.down {
  color: #4caf50;
}

.metric-trend.neutral {
  color: var(--text-secondary);
}

.metric-value {
  display: block;
  font-size: 28px;
  font-weight: 600;
}

.metric-unit {
  color: var(--text-secondary);
  font-size: 12px;
}

.services-section,
.timeseries-section,
.endpoints-section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 24px;
}

.services-section h2,
.timeseries-section h2,
.endpoints-section h2 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.services-table {
  display: flex;
  flex-direction: column;
}

.table-header,
.table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr 1fr;
  gap: 12px;
  padding: 8px 0;
}

.table-header {
  color: var(--text-secondary);
  font-size: 12px;
  border-bottom: 1px solid var(--border-color);
}

.table-row {
  font-size: 13px;
  border-bottom: 1px solid var(--border-color);
}

.table-row:last-child {
  border-bottom: none;
}

.col-errors.warning {
  color: #ff9800;
}

.chart-placeholder {
  padding: 16px 0;
}

.chart-area {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  height: 150px;
}

.chart-bar {
  flex: 1;
  background: var(--primary);
  border-radius: 4px 4px 0 0;
  min-height: 4px;
}

.chart-labels {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  font-size: 10px;
  color: var(--text-secondary);
}

.chart-labels span {
  flex: 1;
  text-align: center;
}

.endpoints-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.endpoint-item {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 12px;
  align-items: center;
}

.endpoint-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.endpoint-method {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
}

.endpoint-method.get {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.endpoint-method.post {
  background: rgba(33, 150, 243, 0.2);
  color: #2196f3;
}

.endpoint-method.put {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.endpoint-method.delete {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.endpoint-path {
  font-family: monospace;
  font-size: 12px;
}

.endpoint-stats {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--text-secondary);
}

.endpoint-bar {
  height: 4px;
  background: var(--border-color);
  border-radius: 2px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: var(--primary);
  border-radius: 2px;
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
</style>