<template>
  <div class="elk-page">
    <div class="page-header">
      <h1>ELK Stack</h1>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="refresh">
          Refresh
        </button>
      </div>
    </div>

    <!-- Cluster Health -->
    <div class="health-section">
      <h2>Cluster Health</h2>
      <div class="health-grid">
        <div :class="['health-card', clusterHealth.status]">
          <span class="health-status">{{ clusterHealth.status.toUpperCase() }}</span>
          <span class="health-label">Status</span>
        </div>
        <div class="health-card">
          <span class="health-value">{{ clusterHealth.number_of_nodes }}</span>
          <span class="health-label">Nodes</span>
        </div>
        <div class="health-card">
          <span class="health-value">{{ clusterHealth.active_shards }}</span>
          <span class="health-label">Active Shards</span>
        </div>
        <div class="health-card">
          <span class="health-value">{{ clusterHealth.unassigned_shards }}</span>
          <span class="health-label">Unassigned</span>
        </div>
      </div>
    </div>

    <!-- Indices -->
    <div class="indices-section">
      <h2>Indices</h2>
      <div class="indices-grid">
        <div v-for="index in indices" :key="index.name" class="index-card">
          <div class="index-header">
            <span class="index-name">{{ index.name }}</span>
            <span :class="['index-health', index.health]">{{ index.health }}</span>
          </div>
          <div class="index-stats">
            <div class="stat">
              <span class="stat-value">{{ formatNumber(index.docs) }}</span>
              <span class="stat-label">Docs</span>
            </div>
            <div class="stat">
              <span class="stat-value">{{ formatBytes(index.size) }}</span>
              <span class="stat-label">Size</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Index Templates -->
    <div class="templates-section">
      <h2>Index Templates</h2>
      <div class="templates-list">
        <div v-for="template in templates" :key="template.name" class="template-item">
          <span class="template-name">{{ template.name }}</span>
          <span class="template-patterns">{{ template.index_patterns.join(', ') }}</span>
        </div>
      </div>
    </div>

    <!-- ILM Policies -->
    <div class="ilm-section">
      <h2>Index Lifecycle Policies</h2>
      <div class="ilm-list">
        <div v-for="policy in ilmPolicies" :key="policy.name" class="ilm-item">
          <span class="ilm-name">{{ policy.name }}</span>
          <div class="ilm-phases">
            <span v-for="phase in policy.phases" :key="phase" :class="['phase-badge', phase]">
              {{ phase }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Pipelines -->
    <div class="pipelines-section">
      <h2>Logstash Pipelines</h2>
      <div class="pipelines-list">
        <div v-for="pipeline in pipelines" :key="pipeline.name" class="pipeline-item">
          <span class="pipeline-name">{{ pipeline.name }}</span>
          <span class="pipeline-description">{{ pipeline.description }}</span>
        </div>
      </div>
    </div>

    <!-- Dashboards -->
    <div class="dashboards-section">
      <h2>Kibana Dashboards</h2>
      <div class="dashboards-grid">
        <div v-for="dashboard in dashboards" :key="dashboard.id" class="dashboard-card">
          <span class="dashboard-title">{{ dashboard.title }}</span>
          <span class="dashboard-description">{{ dashboard.description }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const clusterHealth = ref({
  cluster_name: 'anixops-logs',
  status: 'green',
  number_of_nodes: 3,
  active_shards: 90,
  unassigned_shards: 0
})

const indices = ref([
  { name: 'logs-app-2026.03.23', health: 'green', docs: 500000, size: 1073741824 },
  { name: 'logs-app-2026.03.22', health: 'green', docs: 750000, size: 1610612736 },
  { name: 'metrics-app-2026.03.23', health: 'green', docs: 125000, size: 214748364 },
])

const templates = ref([
  { name: 'logs-app', index_patterns: ['logs-app-*'] },
  { name: 'metrics-app', index_patterns: ['metrics-app-*'] },
])

const ilmPolicies = ref([
  { name: 'logs-policy', phases: ['hot', 'warm', 'cold', 'delete'] },
  { name: 'metrics-policy', phases: ['hot', 'delete'] },
])

const pipelines = ref([
  { name: 'logs-pipeline', description: 'Process application logs' },
])

const dashboards = ref([
  { id: 'logs-overview', title: 'Application Logs Overview', description: 'Overview of application logs' },
  { id: 'metrics-overview', title: 'Metrics Overview', description: 'System and application metrics' },
])

function formatNumber(num) {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

function formatBytes(bytes) {
  if (bytes >= 1073741824) return (bytes / 1073741824).toFixed(1) + 'GB'
  if (bytes >= 1048576) return (bytes / 1048576).toFixed(1) + 'MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(1) + 'KB'
  return bytes + 'B'
}

function refresh() {
  // TODO: Implement refresh
}

onMounted(() => {
  // TODO: Fetch data from API
})
</script>

<style scoped>
.elk-page {
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

.health-section,
.indices-section,
.templates-section,
.ilm-section,
.pipelines-section,
.dashboards-section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 24px;
}

.health-section h2,
.indices-section h2,
.templates-section h2,
.ilm-section h2,
.pipelines-section h2,
.dashboards-section h2 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.health-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.health-card {
  background: var(--bg);
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  border-left: 4px solid var(--border-color);
}

.health-card.green {
  border-left-color: #4caf50;
}

.health-card.yellow {
  border-left-color: #ff9800;
}

.health-card.red {
  border-left-color: #f44336;
}

.health-status {
  display: block;
  font-size: 24px;
  font-weight: 600;
}

.health-card.green .health-status {
  color: #4caf50;
}

.health-card.yellow .health-status {
  color: #ff9800;
}

.health-card.red .health-status {
  color: #f44336;
}

.health-value {
  display: block;
  font-size: 24px;
  font-weight: 600;
}

.health-label {
  color: var(--text-secondary);
  font-size: 12px;
}

.indices-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.index-card {
  background: var(--bg);
  border-radius: 8px;
  padding: 12px;
}

.index-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.index-name {
  font-family: monospace;
  font-size: 12px;
}

.index-health {
  font-size: 10px;
  text-transform: uppercase;
  padding: 2px 6px;
  border-radius: 4px;
}

.index-health.green {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.index-health.yellow {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.index-health.red {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.index-stats {
  display: flex;
  gap: 16px;
}

.stat {
  text-align: center;
}

.stat-value {
  display: block;
  font-weight: 600;
}

.stat-label {
  font-size: 10px;
  color: var(--text-secondary);
}

.templates-list,
.ilm-list,
.pipelines-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.template-item,
.ilm-item,
.pipeline-item {
  background: var(--bg);
  border-radius: 8px;
  padding: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.template-name,
.ilm-name,
.pipeline-name {
  font-weight: 500;
}

.template-patterns,
.pipeline-description {
  font-size: 12px;
  color: var(--text-secondary);
}

.ilm-phases {
  display: flex;
  gap: 8px;
}

.phase-badge {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 4px;
  text-transform: capitalize;
}

.phase-badge.hot {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.phase-badge.warm {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.phase-badge.cold {
  background: rgba(33, 150, 243, 0.2);
  color: #2196f3;
}

.phase-badge.delete {
  background: rgba(158, 158, 158, 0.2);
  color: #9e9e9e;
}

.dashboards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.dashboard-card {
  background: var(--bg);
  border-radius: 8px;
  padding: 12px;
  cursor: pointer;
}

.dashboard-card:hover {
  background: var(--hover-bg);
}

.dashboard-title {
  display: block;
  font-weight: 500;
  margin-bottom: 4px;
}

.dashboard-description {
  font-size: 12px;
  color: var(--text-secondary);
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