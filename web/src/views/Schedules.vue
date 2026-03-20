<template>
  <div class="schedules-page">
    <div class="page-header">
      <h1>Schedules</h1>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateDialog = true">
          <span class="icon">➕</span> New Schedule
        </button>
        <button class="btn btn-secondary" @click="fetchSchedules">
          <span class="icon">🔄</span> Refresh
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon total">📅</div>
        <div class="stat-info">
          <span class="stat-value">{{ schedules.length }}</span>
          <span class="stat-label">Total Schedules</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon enabled">✅</div>
        <div class="stat-info">
          <span class="stat-value">{{ enabledCount }}</span>
          <span class="stat-label">Enabled</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon disabled">⏸️</div>
        <div class="stat-info">
          <span class="stat-value">{{ disabledCount }}</span>
          <span class="stat-label">Disabled</span>
        </div>
      </div>
    </div>

    <!-- Schedules Grid -->
    <div class="schedules-grid">
      <div v-if="loading" class="loading-state">Loading...</div>
      <div v-else-if="schedules.length === 0" class="empty-state">
        No schedules found. Create one to get started.
      </div>
      <div v-else class="schedule-cards">
        <div v-for="schedule in schedules" :key="schedule.id" class="schedule-card">
          <div class="card-header">
            <div class="schedule-name">
              <h3>{{ schedule.name }}</h3>
              <span :class="['status-badge', schedule.enabled ? 'enabled' : 'disabled']">
                {{ schedule.enabled ? 'Enabled' : 'Disabled' }}
              </span>
            </div>
            <div class="card-actions">
              <button class="btn-icon" @click="toggleSchedule(schedule)" :title="schedule.enabled ? 'Disable' : 'Enable'">
                {{ schedule.enabled ? '⏸️' : '▶️' }}
              </button>
              <button class="btn-icon" @click="runSchedule(schedule)" title="Run Now">
                🚀
              </button>
              <button class="btn-icon" @click="editSchedule(schedule)" title="Edit">
                ✏️
              </button>
              <button class="btn-icon danger" @click="deleteSchedule(schedule)" title="Delete">
                🗑️
              </button>
            </div>
          </div>
          <div class="card-body">
            <div class="info-row">
              <span class="label">Playbook:</span>
              <span class="value">{{ schedule.playbook_name || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="label">Schedule:</span>
              <span class="value">
                <code>{{ schedule.cron }}</code>
                <span class="cron-desc">{{ getCronDescription(schedule.cron) }}</span>
              </span>
            </div>
            <div class="info-row">
              <span class="label">Timezone:</span>
              <span class="value">{{ schedule.timezone }}</span>
            </div>
            <div class="info-row">
              <span class="label">Target Nodes:</span>
              <span class="value">{{ schedule.target_nodes?.length || 0 }} nodes</span>
            </div>
            <div v-if="schedule.next_run" class="info-row">
              <span class="label">Next Run:</span>
              <span class="value">{{ formatDate(schedule.next_run) }}</span>
            </div>
            <div v-if="schedule.last_run" class="info-row">
              <span class="label">Last Run:</span>
              <span class="value">{{ formatDate(schedule.last_run) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Schedule Dialog -->
    <div v-if="showCreateDialog || editingSchedule" class="modal-overlay" @click.self="closeDialog">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingSchedule ? 'Edit Schedule' : 'Create Schedule' }}</h2>
          <button class="close-btn" @click="closeDialog">✕</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveSchedule">
            <div class="form-group">
              <label>Name</label>
              <input v-model="formData.name" type="text" required placeholder="Schedule name" />
            </div>
            <div class="form-group">
              <label>Playbook</label>
              <select v-model="formData.playbook_id" required>
                <option value="">Select playbook</option>
                <option v-for="p in playbooks" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>Cron Expression</label>
                <input v-model="formData.cron" type="text" required placeholder="0 * * * *" />
              </div>
              <div class="form-group">
                <label>Timezone</label>
                <select v-model="formData.timezone">
                  <option>UTC</option>
                  <option>America/New_York</option>
                  <option>America/Los_Angeles</option>
                  <option>Europe/London</option>
                  <option>Asia/Shanghai</option>
                  <option>Asia/Tokyo</option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label>Target Nodes</label>
              <div class="node-selector">
                <label v-for="node in nodes" :key="node.id" class="checkbox-label">
                  <input type="checkbox" :value="node.id" v-model="formData.target_nodes" />
                  {{ node.name }}
                </label>
              </div>
            </div>
            <div class="form-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="formData.enabled" />
                Enabled
              </label>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeDialog">Cancel</button>
          <button class="btn btn-primary" @click="saveSchedule">
            {{ editingSchedule ? 'Save' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { get, post, put, del } = useApi()

const schedules = ref([])
const playbooks = ref([])
const nodes = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const editingSchedule = ref(null)

const formData = ref({
  name: '',
  playbook_id: '',
  cron: '0 * * * *',
  timezone: 'UTC',
  target_nodes: [],
  enabled: true
})

const enabledCount = computed(() => schedules.value.filter(s => s.enabled).length)
const disabledCount = computed(() => schedules.value.filter(s => !s.enabled).length)

async function fetchSchedules() {
  loading.value = true
  try {
    const res = await get('/schedules')
    schedules.value = res.data?.items || []
  } catch (e) {
    console.error('Failed to fetch schedules:', e)
    schedules.value = []
  } finally {
    loading.value = false
  }
}

async function fetchPlaybooks() {
  try {
    const res = await get('/playbooks')
    playbooks.value = res.data?.items || []
  } catch (e) {
    playbooks.value = []
  }
}

async function fetchNodes() {
  try {
    const res = await get('/nodes')
    nodes.value = res.data?.items || []
  } catch (e) {
    nodes.value = []
  }
}

async function saveSchedule() {
  try {
    const payload = {
      name: formData.value.name,
      playbook_id: formData.value.playbook_id,
      cron: formData.value.cron,
      timezone: formData.value.timezone,
      target_nodes: formData.value.target_nodes,
      enabled: formData.value.enabled
    }

    if (editingSchedule.value) {
      await put(`/schedules/${editingSchedule.value.id}`, payload)
    } else {
      await post('/schedules', payload)
    }

    closeDialog()
    fetchSchedules()
  } catch (e) {
    alert('Failed to save schedule: ' + e.message)
  }
}

async function toggleSchedule(schedule) {
  try {
    await post(`/schedules/${schedule.id}/toggle`)
    fetchSchedules()
  } catch (e) {
    alert('Failed to toggle schedule')
  }
}

async function runSchedule(schedule) {
  if (!confirm(`Run "${schedule.name}" now?`)) return
  try {
    const res = await post(`/schedules/${schedule.id}/run`)
    alert(`Task created: ${res.data?.task_id || 'OK'}`)
    fetchSchedules()
  } catch (e) {
    alert('Failed to run schedule')
  }
}

async function deleteSchedule(schedule) {
  if (!confirm(`Delete "${schedule.name}"?`)) return
  try {
    await del(`/schedules/${schedule.id}`)
    fetchSchedules()
  } catch (e) {
    alert('Failed to delete schedule')
  }
}

function editSchedule(schedule) {
  editingSchedule.value = schedule
  formData.value = {
    name: schedule.name,
    playbook_id: schedule.playbook_id,
    cron: schedule.cron,
    timezone: schedule.timezone || 'UTC',
    target_nodes: schedule.target_nodes || [],
    enabled: schedule.enabled
  }
}

function closeDialog() {
  showCreateDialog.value = false
  editingSchedule.value = null
  formData.value = {
    name: '',
    playbook_id: '',
    cron: '0 * * * *',
    timezone: 'UTC',
    target_nodes: [],
    enabled: true
  }
}

function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

function getCronDescription(cron) {
  if (!cron) return ''
  const parts = cron.split(' ')
  if (parts.length !== 5) return cron

  if (parts[0].startsWith('*/')) {
    return `Every ${parts[0].substring(2)} minutes`
  }
  if (parts[0] === '0' && parts[1] === '*') {
    return 'Hourly'
  }
  if (parts[1] !== '*' && !parts[1].contains('/')) {
    const hour = parseInt(parts[1])
    const minute = parseInt(parts[0])
    return `Daily at ${hour}:${minute.toString().padStart(2, '0')}`
  }
  return cron
}

onMounted(() => {
  fetchSchedules()
  fetchPlaybooks()
  fetchNodes()
})
</script>

<style scoped>
.schedules-page {
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.stat-icon.total { background: rgba(33, 150, 243, 0.2); }
.stat-icon.enabled { background: rgba(76, 175, 80, 0.2); }
.stat-icon.disabled { background: rgba(158, 158, 158, 0.2); }

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
}

.stat-label {
  color: var(--text-secondary);
  font-size: 14px;
}

.schedules-grid {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 16px;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.schedule-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
}

.schedule-card {
  background: var(--bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--header-bg);
  border-bottom: 1px solid var(--border-color);
}

.schedule-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.schedule-name h3 {
  margin: 0;
  font-size: 16px;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
}

.status-badge.enabled {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
}

.status-badge.disabled {
  background: rgba(158, 158, 158, 0.2);
  color: #9e9e9e;
}

.card-actions {
  display: flex;
  gap: 4px;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  font-size: 14px;
  border-radius: 4px;
}

.btn-icon:hover {
  background: var(--hover-bg);
}

.btn-icon.danger:hover {
  background: rgba(244, 67, 54, 0.2);
}

.card-body {
  padding: 12px 16px;
}

.info-row {
  display: flex;
  gap: 8px;
  padding: 4px 0;
  font-size: 13px;
}

.info-row .label {
  color: var(--text-secondary);
  min-width: 100px;
}

.info-row code {
  background: var(--code-bg);
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
}

.cron-desc {
  color: var(--text-secondary);
  font-size: 12px;
  margin-left: 8px;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
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
  width: 500px;
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

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.node-selector {
  max-height: 150px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  cursor: pointer;
}
</style>