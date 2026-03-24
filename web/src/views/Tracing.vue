<template>
  <div class="tracing-page">
    <div class="page-header">
      <h1>Distributed Tracing</h1>
      <div class="header-actions">
        <select v-model="serviceFilter" class="filter-select">
          <option value="">All Services</option>
          <option v-for="s in services" :key="s" :value="s">{{ s }}</option>
        </select>
        <select v-model="statusFilter" class="filter-select">
          <option value="">All Status</option>
          <option value="ok">OK</option>
          <option value="error">Error</option>
        </select>
        <button class="btn btn-secondary" @click="refresh">
          Refresh
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-value">{{ stats.totalTraces }}</span>
        <span class="stat-label">Total Traces</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.totalSpans }}</span>
        <span class="stat-label">Total Spans</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.averageDuration.toFixed(0) }}ms</span>
        <span class="stat-label">Avg Duration</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ (stats.errorRate * 100).toFixed(1) }}%</span>
        <span class="stat-label">Error Rate</span>
      </div>
    </div>

    <!-- Traces List -->
    <div class="traces-section">
      <h2>Recent Traces</h2>
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else-if="filteredTraces.length === 0" class="empty">
        No traces found
      </div>
      <div v-else class="traces-list">
        <div
          v-for="trace in filteredTraces"
          :key="trace.traceId"
          class="trace-item"
          @click="selectTrace(trace)"
        >
          <div class="trace-header">
            <span class="trace-name">{{ trace.rootSpan?.name || 'Unknown' }}</span>
            <span :class="['trace-status', trace.status]">{{ trace.status }}</span>
          </div>
          <div class="trace-details">
            <span class="trace-id">{{ trace.traceId.substring(0, 16) }}...</span>
            <span class="trace-duration">{{ trace.duration }}ms</span>
            <span class="trace-spans">{{ trace.spanCount }} spans</span>
            <span class="trace-services">{{ trace.serviceCount }} services</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Trace Detail Modal -->
    <div v-if="selectedTrace" class="modal-overlay" @click.self="selectedTrace = null">
      <div class="modal">
        <div class="modal-header">
          <h2>Trace Details</h2>
          <button class="close-btn" @click="selectedTrace = null">✕</button>
        </div>
        <div class="modal-body">
          <div class="trace-info">
            <div class="info-row">
              <span class="label">Trace ID:</span>
              <span class="value code">{{ selectedTrace.traceId }}</span>
            </div>
            <div class="info-row">
              <span class="label">Duration:</span>
              <span class="value">{{ selectedTrace.duration }}ms</span>
            </div>
            <div class="info-row">
              <span class="label">Status:</span>
              <span :class="['value', selectedTrace.status]">{{ selectedTrace.status }}</span>
            </div>
          </div>

          <h3>Spans</h3>
          <div class="spans-list">
            <div
              v-for="span in selectedTrace.spans"
              :key="span.spanId"
              :class="['span-item', span.status.code]"
            >
              <div class="span-header">
                <span class="span-name">{{ span.name }}</span>
                <span class="span-duration">{{ span.duration }}ms</span>
              </div>
              <div class="span-details">
                <span class="span-service">{{ span.resource['service.name'] || 'unknown' }}</span>
                <span class="span-kind">{{ span.kind }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const loading = ref(false)
const serviceFilter = ref('')
const statusFilter = ref('')
const selectedTrace = ref(null)

const stats = ref({
  totalTraces: 1234,
  totalSpans: 5678,
  averageDuration: 145,
  errorRate: 0.02,
})

const services = ref(['api-gateway', 'auth-service', 'task-runner', 'log-processor'])

const traces = ref([
  {
    traceId: '0af7651916cd43dd8448eb211c80319c',
    status: 'ok',
    duration: 150,
    spanCount: 5,
    serviceCount: 3,
    rootSpan: { name: 'HTTP GET /api/users' },
    spans: [
      { spanId: '1', name: 'HTTP GET /api/users', duration: 150, kind: 'server', status: { code: 'ok' }, resource: { 'service.name': 'api-gateway' } },
      { spanId: '2', name: 'auth:validate_token', duration: 20, kind: 'client', status: { code: 'ok' }, resource: { 'service.name': 'api-gateway' } },
      { spanId: '3', name: 'db:query', duration: 50, kind: 'client', status: { code: 'ok' }, resource: { 'service.name': 'api-gateway' } },
    ],
  },
  {
    traceId: '1bf7651916cd43dd8448eb211c80319d',
    status: 'error',
    duration: 300,
    spanCount: 4,
    serviceCount: 2,
    rootSpan: { name: 'HTTP POST /api/tasks' },
    spans: [
      { spanId: '1', name: 'HTTP POST /api/tasks', duration: 300, kind: 'server', status: { code: 'error', message: 'Timeout' }, resource: { 'service.name': 'api-gateway' } },
      { spanId: '2', name: 'task:execute', duration: 280, kind: 'client', status: { code: 'error', message: 'Timeout' }, resource: { 'service.name': 'task-runner' } },
    ],
  },
])

const filteredTraces = computed(() => {
  let result = traces.value

  if (statusFilter.value) {
    result = result.filter(t => t.status === statusFilter.value)
  }

  if (serviceFilter.value) {
    result = result.filter(t =>
      t.spans.some(s => s.resource['service.name'] === serviceFilter.value)
    )
  }

  return result
})

function selectTrace(trace) {
  selectedTrace.value = trace
}

function refresh() {
  // TODO: Implement refresh
}

onMounted(() => {
  // TODO: Fetch data from API
})
</script>

<style scoped>
.tracing-page {
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

.traces-section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
}

.traces-section h2 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.loading,
.empty {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.traces-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.trace-item {
  padding: 12px 16px;
  background: var(--bg);
  border-radius: 8px;
  cursor: pointer;
  border-left: 4px solid;
}

.trace-item:hover {
  background: var(--hover-bg);
}

.trace-item .ok {
  border-left-color: #4caf50;
}

.trace-item .error {
  border-left-color: #f44336;
}

.trace-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.trace-name {
  font-weight: 600;
}

.trace-status {
  font-size: 11px;
  text-transform: uppercase;
  padding: 2px 8px;
  border-radius: 4px;
}

.trace-status.ok {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.trace-status.error {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.trace-details {
  display: flex;
  gap: 16px;
  color: var(--text-secondary);
  font-size: 12px;
}

.trace-id {
  font-family: monospace;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--card-bg);
  border-radius: 12px;
  width: 600px;
  max-height: 80vh;
  overflow: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text-secondary);
}

.modal-body {
  padding: 24px;
}

.trace-info {
  margin-bottom: 24px;
}

.info-row {
  display: flex;
  gap: 8px;
  padding: 8px 0;
}

.info-row .label {
  color: var(--text-secondary);
  min-width: 100px;
}

.info-row .value.code {
  font-family: monospace;
  font-size: 12px;
}

.info-row .value.error {
  color: #f44336;
}

.info-row .value.ok {
  color: #4caf50;
}

h3 {
  margin: 0 0 16px 0;
  font-size: 14px;
}

.spans-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.span-item {
  padding: 12px;
  background: var(--bg);
  border-radius: 8px;
  border-left: 3px solid;
}

.span-item.ok {
  border-left-color: #4caf50;
}

.span-item.error {
  border-left-color: #f44336;
}

.span-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.span-name {
  font-weight: 500;
}

.span-duration {
  color: var(--text-secondary);
}

.span-details {
  display: flex;
  gap: 12px;
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