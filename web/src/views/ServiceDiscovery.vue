<template>
  <div class="discovery-page">
    <div class="page-header">
      <h1>Service Discovery</h1>
      <div class="header-actions">
        <select v-model="namespaceFilter" class="filter-select">
          <option value="">All Namespaces</option>
          <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
        </select>
        <select v-model="statusFilter" class="filter-select">
          <option value="">All Status</option>
          <option value="healthy">Healthy</option>
          <option value="unhealthy">Unhealthy</option>
          <option value="starting">Starting</option>
          <option value="draining">Draining</option>
        </select>
        <button class="btn btn-secondary" @click="refresh">
          Refresh
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-value">{{ stats.totalServices }}</span>
        <span class="stat-label">Services</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.totalInstances }}</span>
        <span class="stat-label">Instances</span>
      </div>
      <div class="stat-card healthy">
        <span class="stat-value">{{ stats.healthyInstances }}</span>
        <span class="stat-label">Healthy</span>
      </div>
      <div class="stat-card unhealthy">
        <span class="stat-value">{{ stats.unhealthyInstances }}</span>
        <span class="stat-label">Unhealthy</span>
      </div>
    </div>

    <!-- Services List -->
    <div class="services-section">
      <h2>Registered Services</h2>
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else-if="filteredServices.length === 0" class="empty">
        No services found
      </div>
      <div v-else class="services-list">
        <div
          v-for="service in filteredServices"
          :key="service.id"
          :class="['service-item', service.health]"
        >
          <div class="service-header">
            <span class="service-name">{{ service.name }}</span>
            <span :class="['service-status', service.health]">{{ service.health }}</span>
          </div>
          <div class="service-details">
            <span class="service-namespace">{{ service.namespace }}</span>
            <span class="service-endpoint">{{ service.protocol }}://{{ service.host }}:{{ service.port }}</span>
            <span class="service-weight">Weight: {{ service.weight }}</span>
          </div>
          <div class="service-meta">
            <span v-for="(value, key) in service.metadata" :key="key" class="meta-tag">
              {{ key }}: {{ value }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Load Balancer Config -->
    <div class="lb-section">
      <h2>Load Balancer</h2>
      <div class="lb-config">
        <div class="lb-option">
          <label>Algorithm:</label>
          <select v-model="lbAlgorithm" class="lb-select">
            <option value="round-robin">Round Robin</option>
            <option value="weighted">Weighted</option>
            <option value="least-connections">Least Connections</option>
            <option value="ip-hash">IP Hash</option>
            <option value="random">Random</option>
          </select>
        </div>
        <div class="lb-stats">
          <div class="lb-stat">
            <span class="lb-stat-value">{{ lbStats.requests }}</span>
            <span class="lb-stat-label">Requests</span>
          </div>
          <div class="lb-stat">
            <span class="lb-stat-value">{{ lbStats.successRate.toFixed(1) }}%</span>
            <span class="lb-stat-label">Success Rate</span>
          </div>
          <div class="lb-stat">
            <span class="lb-stat-value">{{ lbStats.avgLatency.toFixed(0) }}ms</span>
            <span class="lb-stat-label">Avg Latency</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const loading = ref(false)
const namespaceFilter = ref('')
const statusFilter = ref('')
const lbAlgorithm = ref('round-robin')

const stats = ref({
  totalServices: 4,
  totalInstances: 8,
  healthyInstances: 7,
  unhealthyInstances: 1,
})

const lbStats = ref({
  requests: 15420,
  successRate: 99.8,
  avgLatency: 24,
})

const namespaces = ref(['default', 'production', 'staging'])

const services = ref([
  {
    id: '1',
    name: 'api-gateway',
    namespace: 'production',
    host: '10.0.0.1',
    port: 8080,
    protocol: 'https',
    health: 'healthy',
    weight: 100,
    metadata: { version: 'v2.1.0', region: 'us-east' },
  },
  {
    id: '2',
    name: 'auth-service',
    namespace: 'production',
    host: '10.0.0.2',
    port: 8081,
    protocol: 'https',
    health: 'healthy',
    weight: 100,
    metadata: { version: 'v1.5.0', region: 'us-east' },
  },
  {
    id: '3',
    name: 'task-runner',
    namespace: 'production',
    host: '10.0.0.3',
    port: 8082,
    protocol: 'grpc',
    health: 'healthy',
    weight: 50,
    metadata: { version: 'v1.2.0', region: 'us-east' },
  },
  {
    id: '4',
    name: 'log-processor',
    namespace: 'default',
    host: '10.0.0.4',
    port: 8083,
    protocol: 'http',
    health: 'unhealthy',
    weight: 100,
    metadata: { version: 'v1.0.0', region: 'us-east' },
  },
])

const filteredServices = computed(() => {
  let result = services.value

  if (namespaceFilter.value) {
    result = result.filter(s => s.namespace === namespaceFilter.value)
  }

  if (statusFilter.value) {
    result = result.filter(s => s.health === statusFilter.value)
  }

  return result
})

function refresh() {
  // TODO: Implement refresh
}

onMounted(() => {
  // TODO: Fetch data from API
})
</script>

<style scoped>
.discovery-page {
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

.filter-select,
.lb-select {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  border-left: 4px solid var(--border-color);
}

.stat-card.healthy {
  border-left-color: #4caf50;
}

.stat-card.unhealthy {
  border-left-color: #f44336;
}

.stat-value {
  display: block;
  font-size: 28px;
  font-weight: 600;
}

.stat-label {
  color: var(--text-secondary);
  font-size: 12px;
}

.services-section,
.lb-section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 24px;
}

.services-section h2,
.lb-section h2 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.loading,
.empty {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.services-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.service-item {
  padding: 12px 16px;
  background: var(--bg);
  border-radius: 8px;
  border-left: 4px solid;
}

.service-item.healthy {
  border-left-color: #4caf50;
}

.service-item.unhealthy {
  border-left-color: #f44336;
}

.service-item.starting {
  border-left-color: #ff9800;
}

.service-item.draining {
  border-left-color: #9e9e9e;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.service-name {
  font-weight: 600;
}

.service-status {
  font-size: 11px;
  text-transform: uppercase;
  padding: 2px 8px;
  border-radius: 4px;
}

.service-status.healthy {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.service-status.unhealthy {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.service-details {
  display: flex;
  gap: 16px;
  color: var(--text-secondary);
  font-size: 12px;
  margin-bottom: 8px;
}

.service-endpoint {
  font-family: monospace;
}

.service-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.meta-tag {
  font-size: 10px;
  padding: 2px 6px;
  background: var(--hover-bg);
  border-radius: 4px;
  color: var(--text-secondary);
}

.lb-config {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.lb-option {
  display: flex;
  align-items: center;
  gap: 12px;
}

.lb-option label {
  font-weight: 500;
}

.lb-stats {
  display: flex;
  gap: 24px;
}

.lb-stat {
  text-align: center;
}

.lb-stat-value {
  display: block;
  font-size: 20px;
  font-weight: 600;
}

.lb-stat-label {
  color: var(--text-secondary);
  font-size: 11px;
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