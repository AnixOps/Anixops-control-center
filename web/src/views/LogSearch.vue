<template>
  <div class="logs-search-page">
    <div class="page-header">
      <h1>Log Search</h1>
      <div class="header-actions">
        <select v-model="selectedService" class="filter-select">
          <option value="">All Services</option>
          <option v-for="s in services" :key="s" :value="s">{{ s }}</option>
        </select>
        <select v-model="selectedLevel" class="filter-select">
          <option value="">All Levels</option>
          <option value="ERROR">ERROR</option>
          <option value="WARN">WARN</option>
          <option value="INFO">INFO</option>
          <option value="DEBUG">DEBUG</option>
        </select>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search logs..."
          class="search-input"
          @keyup.enter="search"
        />
        <button class="btn btn-primary" @click="search">Search</button>
      </div>
    </div>

    <!-- Search Stats -->
    <div class="stats-bar">
      <span class="stat-item">
        <span class="stat-value">{{ totalHits }}</span> results
      </span>
      <span class="stat-item">
        <span class="stat-value">{{ searchTime }}ms</span> search time
      </span>
    </div>

    <!-- Time Range Filter -->
    <div class="time-filter">
      <button
        v-for="range in timeRanges"
        :key="range.value"
        :class="['time-btn', { active: selectedTimeRange === range.value }]"
        @click="selectedTimeRange = range.value"
      >
        {{ range.label }}
      </button>
    </div>

    <!-- Results -->
    <div class="results-section">
      <div v-if="loading" class="loading">Searching...</div>
      <div v-else-if="logs.length === 0" class="empty">No logs found</div>
      <div v-else class="logs-list">
        <div
          v-for="log in logs"
          :key="log._id"
          :class="['log-item', log._source.level.toLowerCase()]"
        >
          <div class="log-header">
            <span class="log-timestamp">{{ formatTime(log._source['@timestamp']) }}</span>
            <span :class="['log-level', log._source.level]">{{ log._source.level }}</span>
            <span class="log-service">{{ log._source.service }}</span>
          </div>
          <div class="log-message">{{ log._source.message }}</div>
          <div v-if="log._source.trace_id" class="log-trace">
            <span class="trace-label">Trace:</span>
            <span class="trace-value">{{ log._source.trace_id }}</span>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="logs.length > 0" class="pagination">
        <button class="btn btn-secondary" :disabled="page === 1" @click="prevPage">
          Previous
        </button>
        <span class="page-info">Page {{ page }} of {{ totalPages }}</span>
        <button class="btn btn-secondary" :disabled="page >= totalPages" @click="nextPage">
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const loading = ref(false)
const searchQuery = ref('')
const selectedService = ref('')
const selectedLevel = ref('')
const selectedTimeRange = ref('1h')
const page = ref(1)
const totalHits = ref(0)
const searchTime = ref(0)
const totalPages = ref(1)

const services = ref(['api-gateway', 'auth-service', 'task-runner', 'log-processor'])

const timeRanges = [
  { label: '15m', value: '15m' },
  { label: '1h', value: '1h' },
  { label: '6h', value: '6h' },
  { label: '24h', value: '24h' },
  { label: '7d', value: '7d' },
  { label: '30d', value: '30d' },
]

const logs = ref([
  {
    _id: 'abc123',
    _score: 1.5,
    _source: {
      '@timestamp': '2026-03-23T10:00:00.000Z',
      message: 'Request processed successfully',
      level: 'INFO',
      service: 'api-gateway',
      trace_id: '0af7651916cd43dd8448eb211c80319c',
      host: 'node-1'
    }
  },
  {
    _id: 'def456',
    _score: 1.2,
    _source: {
      '@timestamp': '2026-03-23T10:01:00.000Z',
      message: 'Database connection timeout after 30s',
      level: 'ERROR',
      service: 'auth-service',
      trace_id: '1bf7651916cd43dd8448eb211c80319d',
      host: 'node-2'
    }
  },
  {
    _id: 'ghi789',
    _score: 1.0,
    _source: {
      '@timestamp': '2026-03-23T10:02:00.000Z',
      message: 'High memory usage detected: 85%',
      level: 'WARN',
      service: 'task-runner',
      host: 'node-3'
    }
  }
])

function formatTime(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleString()
}

function search() {
  loading.value = true
  totalHits.value = 1250
  searchTime.value = 5
  totalPages.value = 125
  loading.value = false
}

function prevPage() {
  if (page.value > 1) {
    page.value--
    search()
  }
}

function nextPage() {
  if (page.value < totalPages.value) {
    page.value++
    search()
  }
}

onMounted(() => {
  totalHits.value = 1250
  totalPages.value = 125
})
</script>

<style scoped>
.logs-search-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
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
.search-input {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
}

.search-input {
  width: 250px;
}

.stats-bar {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
  padding: 8px 16px;
  background: var(--card-bg);
  border-radius: 8px;
}

.stat-item {
  color: var(--text-secondary);
  font-size: 13px;
}

.stat-value {
  color: var(--text);
  font-weight: 600;
}

.time-filter {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.time-btn {
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 12px;
}

.time-btn.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.results-section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
}

.loading,
.empty {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.logs-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-item {
  padding: 12px;
  background: var(--bg);
  border-radius: 8px;
  border-left: 4px solid;
}

.log-item.info {
  border-left-color: #2196f3;
}

.log-item.warn {
  border-left-color: #ff9800;
}

.log-item.error {
  border-left-color: #f44336;
}

.log-item.debug {
  border-left-color: #9e9e9e;
}

.log-header {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 8px;
}

.log-timestamp {
  font-family: monospace;
  font-size: 12px;
  color: var(--text-secondary);
}

.log-level {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
  font-weight: 600;
}

.log-level.INFO {
  background: rgba(33, 150, 243, 0.2);
  color: #2196f3;
}

.log-level.WARN {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.log-level.ERROR {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.log-level.DEBUG {
  background: rgba(158, 158, 158, 0.2);
  color: #9e9e9e;
}

.log-service {
  font-size: 11px;
  padding: 2px 8px;
  background: var(--hover-bg);
  border-radius: 4px;
}

.log-message {
  font-size: 13px;
  line-height: 1.5;
}

.log-trace {
  margin-top: 8px;
  font-size: 11px;
  color: var(--text-secondary);
}

.trace-label {
  margin-right: 4px;
}

.trace-value {
  font-family: monospace;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.page-info {
  color: var(--text-secondary);
  font-size: 13px;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
}

.btn-primary {
  background: var(--primary);
  color: white;
  border: none;
}

.btn-secondary {
  background: transparent;
  color: var(--text);
  border: 1px solid var(--border-color);
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>