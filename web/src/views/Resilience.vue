<template>
  <div class="resilience-page">
    <div class="page-header">
      <h1>Resilience</h1>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="refresh">Refresh</button>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-value">{{ stats.circuitBreakers.total }}</span>
        <span class="stat-label">Circuit Breakers</span>
      </div>
      <div class="stat-card warning">
        <span class="stat-value">{{ stats.circuitBreakers.open }}</span>
        <span class="stat-label">Open</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.rateLimiters }}</span>
        <span class="stat-label">Rate Limiters</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.retryConfigs }}</span>
        <span class="stat-label">Retry Configs</span>
      </div>
    </div>

    <!-- Circuit Breakers -->
    <div class="section">
      <h2>Circuit Breakers</h2>
      <div class="breakers-grid">
        <div v-for="breaker in circuitBreakers" :key="breaker.name" :class="['breaker-card', breaker.state]">
          <div class="breaker-header">
            <span class="breaker-name">{{ breaker.name }}</span>
            <span :class="['breaker-state', breaker.state]">{{ breaker.state }}</span>
          </div>
          <div class="breaker-stats">
            <div class="breaker-stat">
              <span class="stat-label">Failures</span>
              <span class="stat-value">{{ breaker.failureCount }}</span>
            </div>
            <div class="breaker-stat">
              <span class="stat-label">Successes</span>
              <span class="stat-value">{{ breaker.successCount }}</span>
            </div>
            <div class="breaker-stat">
              <span class="stat-label">Threshold</span>
              <span class="stat-value">{{ breaker.config.failureThreshold }}</span>
            </div>
          </div>
          <div class="breaker-timeline">
            <div class="timeline-bar">
              <div class="timeline-fill" :style="{ width: (breaker.failureCount / breaker.config.failureThreshold * 100) + '%' }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Rate Limiters -->
    <div class="section">
      <h2>Rate Limiters</h2>
      <div class="limiters-list">
        <div v-for="limiter in rateLimiters" :key="limiter.name" class="limiter-item">
          <div class="limiter-info">
            <span class="limiter-name">{{ limiter.name }}</span>
            <span class="limiter-tokens">{{ limiter.tokens }} / {{ limiter.config.maxTokens }} tokens</span>
          </div>
          <div class="limiter-bar">
            <div class="bar-fill" :style="{ width: (limiter.tokens / limiter.config.maxTokens * 100) + '%' }"></div>
          </div>
          <div class="limiter-rate">
            +{{ limiter.config.refillRate }} tokens/sec
          </div>
        </div>
      </div>
    </div>

    <!-- Retry Configurations -->
    <div class="section">
      <h2>Retry Configurations</h2>
      <div class="retries-table">
        <div class="table-header">
          <span class="col-name">Name</span>
          <span class="col-retries">Max Retries</span>
          <span class="col-backoff">Backoff</span>
          <span class="col-initial">Initial</span>
          <span class="col-max">Max Delay</span>
        </div>
        <div v-for="retry in retryConfigs" :key="retry.name" class="table-row">
          <span class="col-name">{{ retry.name }}</span>
          <span class="col-retries">{{ retry.maxRetries }}</span>
          <span class="col-backoff">{{ retry.backoffMultiplier }}x</span>
          <span class="col-initial">{{ retry.initialDelay }}ms</span>
          <span class="col-max">{{ formatDelay(retry.maxDelay) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const stats = ref({
  circuitBreakers: { total: 3, open: 1, closed: 2 },
  rateLimiters: 2,
  retryConfigs: 3
})

const circuitBreakers = ref([
  {
    name: 'api-gateway',
    state: 'closed',
    failureCount: 0,
    successCount: 150,
    config: { failureThreshold: 5, successThreshold: 3, timeout: 60000 }
  },
  {
    name: 'auth-service',
    state: 'open',
    failureCount: 7,
    successCount: 0,
    config: { failureThreshold: 5, successThreshold: 3, timeout: 60000 }
  },
  {
    name: 'database',
    state: 'half-open',
    failureCount: 0,
    successCount: 2,
    config: { failureThreshold: 5, successThreshold: 3, timeout: 60000 }
  }
])

const rateLimiters = ref([
  {
    name: 'api-default',
    tokens: 85,
    config: { maxTokens: 100, refillRate: 10, refillInterval: 1000 }
  },
  {
    name: 'api-burst',
    tokens: 450,
    config: { maxTokens: 500, refillRate: 50, refillInterval: 1000 }
  }
])

const retryConfigs = ref([
  { name: 'default', maxRetries: 3, backoffMultiplier: 2, initialDelay: 100, maxDelay: 30000 },
  { name: 'database', maxRetries: 5, backoffMultiplier: 2, initialDelay: 50, maxDelay: 5000 },
  { name: 'external-api', maxRetries: 2, backoffMultiplier: 3, initialDelay: 500, maxDelay: 10000 }
])

function formatDelay(ms) {
  if (ms >= 1000) return (ms / 1000) + 's'
  return ms + 'ms'
}

function refresh() {
  // TODO: Implement refresh
}

onMounted(() => {
  // TODO: Fetch data from API
})
</script>

<style scoped>
.resilience-page {
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

.stat-card.warning {
  border-left: 4px solid #f44336;
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

.section {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 24px;
}

.section h2 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.breakers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.breaker-card {
  background: var(--bg);
  border-radius: 8px;
  padding: 16px;
  border-left: 4px solid;
}

.breaker-card.closed {
  border-left-color: #4caf50;
}

.breaker-card.open {
  border-left-color: #f44336;
}

.breaker-card.half-open {
  border-left-color: #ff9800;
}

.breaker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.breaker-name {
  font-weight: 600;
}

.breaker-state {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.breaker-state.closed {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.breaker-state.open {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
}

.breaker-state.half-open {
  background: rgba(255, 152, 0, 0.2);
  color: #ff9800;
}

.breaker-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.breaker-stat {
  text-align: center;
}

.breaker-stat .stat-label {
  display: block;
  font-size: 10px;
}

.breaker-stat .stat-value {
  font-size: 16px;
}

.breaker-timeline {
  margin-top: 8px;
}

.timeline-bar {
  height: 4px;
  background: var(--border-color);
  border-radius: 2px;
  overflow: hidden;
}

.timeline-fill {
  height: 100%;
  background: #f44336;
  border-radius: 2px;
}

.limiters-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.limiter-item {
  display: grid;
  grid-template-columns: 1fr 2fr auto;
  gap: 16px;
  align-items: center;
  padding: 12px;
  background: var(--bg);
  border-radius: 8px;
}

.limiter-name {
  font-weight: 500;
}

.limiter-tokens {
  font-size: 12px;
  color: var(--text-secondary);
}

.limiter-bar {
  height: 8px;
  background: var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.limiter-bar .bar-fill {
  height: 100%;
  background: #4caf50;
  border-radius: 4px;
}

.limiter-rate {
  font-size: 11px;
  color: var(--text-secondary);
}

.retries-table {
  display: flex;
  flex-direction: column;
}

.table-header,
.table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr;
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