<template>
  <div class="tasks-page">
    <div class="page-header">
      <h1>Tasks</h1>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateDialog = true">
          <span class="icon">➕</span> New Task
        </button>
        <button class="btn btn-secondary" @click="fetchTasks">
          <span class="icon">🔄</span> Refresh
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon pending">⏳</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.pending }}</span>
          <span class="stat-label">Pending</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon running">🏃</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.running }}</span>
          <span class="stat-label">Running</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon completed">✅</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.completed }}</span>
          <span class="stat-label">Completed</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon failed">❌</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.failed }}</span>
          <span class="stat-label">Failed</span>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="filters-bar">
      <div class="filter-group">
        <label>Status:</label>
        <select v-model="statusFilter" @change="applyFilters">
          <option value="">All</option>
          <option value="pending">Pending</option>
          <option value="running">Running</option>
          <option value="completed">Completed</option>
          <option value="failed">Failed</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>
      <div class="filter-group">
        <label>Playbook:</label>
        <select v-model="playbookFilter" @change="applyFilters">
          <option value="">All Playbooks</option>
          <option v-for="p in playbooks" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>
      </div>
      <div class="search-group">
        <input type="text" v-model="searchQuery" placeholder="Search tasks..." @input="applyFilters" />
      </div>
    </div>

    <!-- Tasks Table -->
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Playbook</th>
            <th>Status</th>
            <th>Target Nodes</th>
            <th>Created</th>
            <th>Duration</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="7" class="loading-row">Loading...</td>
          </tr>
          <tr v-else-if="filteredTasks.length === 0">
            <td colspan="7" class="empty-row">No tasks found</td>
          </tr>
          <tr v-else v-for="task in filteredTasks" :key="task.id" @click="viewTask(task)">
            <td><code>{{ task.id }}</code></td>
            <td>{{ task.playbook_name }}</td>
            <td>
              <span :class="['status-badge', task.status]">
                {{ task.status }}
              </span>
            </td>
            <td>{{ task.target_nodes?.length || 0 }} nodes</td>
            <td>{{ formatDate(task.created_at) }}</td>
            <td>{{ formatDuration(task) }}</td>
            <td class="actions">
              <button v-if="task.status === 'running'" class="btn-icon danger" @click.stop="cancelTask(task)" title="Cancel">
                ⏹️
              </button>
              <button v-if="task.status === 'failed'" class="btn-icon warning" @click.stop="retryTask(task)" title="Retry">
                🔄
              </button>
              <button class="btn-icon" @click.stop="viewTask(task)" title="View Details">
                👁️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Task Dialog -->
    <div v-if="showCreateDialog" class="modal-overlay" @click.self="showCreateDialog = false">
      <div class="modal">
        <div class="modal-header">
          <h2>Create Task</h2>
          <button class="close-btn" @click="showCreateDialog = false">✕</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="createTask">
            <div class="form-group">
              <label>Playbook</label>
              <select v-model="newTask.playbook_id" required>
                <option value="">Select playbook</option>
                <option v-for="p in playbooks" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Target Nodes</label>
              <div class="node-selector">
                <label v-for="node in nodes" :key="node.id" class="checkbox-label">
                  <input type="checkbox" :value="node.id" v-model="newTask.target_nodes" />
                  {{ node.name }}
                </label>
              </div>
            </div>
            <div class="form-group">
              <label>Variables (JSON)</label>
              <textarea v-model="newTask.variables" placeholder='{"key": "value"}' rows="4"></textarea>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showCreateDialog = false">Cancel</button>
          <button class="btn btn-primary" @click="createTask">Create Task</button>
        </div>
      </div>
    </div>

    <!-- Task Detail Modal -->
    <div v-if="selectedTask" class="modal-overlay" @click.self="selectedTask = null">
      <div class="modal modal-lg">
        <div class="modal-header">
          <h2>Task #{{ selectedTask.id }}</h2>
          <button class="close-btn" @click="selectedTask = null">✕</button>
        </div>
        <div class="modal-body">
          <div class="detail-grid">
            <div class="detail-item">
              <label>Playbook:</label>
              <span>{{ selectedTask.playbook_name }}</span>
            </div>
            <div class="detail-item">
              <label>Status:</label>
              <span :class="['status-badge', selectedTask.status]">{{ selectedTask.status }}</span>
            </div>
            <div class="detail-item">
              <label>Created:</label>
              <span>{{ formatDate(selectedTask.created_at) }}</span>
            </div>
            <div class="detail-item">
              <label>Duration:</label>
              <span>{{ formatDuration(selectedTask) }}</span>
            </div>
          </div>
          <div class="logs-section">
            <h3>Execution Logs</h3>
            <div class="logs-container">
              <pre>{{ selectedTask.logs || 'No logs available' }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { get, post } = useApi()

const tasks = ref([])
const playbooks = ref([])
const nodes = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const selectedTask = ref(null)
const statusFilter = ref('')
const playbookFilter = ref('')
const searchQuery = ref('')

const newTask = ref({
  playbook_id: '',
  target_nodes: [],
  variables: ''
})

const stats = computed(() => ({
  pending: tasks.value.filter(t => t.status === 'pending').length,
  running: tasks.value.filter(t => t.status === 'running').length,
  completed: tasks.value.filter(t => t.status === 'completed').length,
  failed: tasks.value.filter(t => t.status === 'failed').length
}))

const filteredTasks = computed(() => {
  let result = tasks.value
  if (statusFilter.value) {
    result = result.filter(t => t.status === statusFilter.value)
  }
  if (playbookFilter.value) {
    result = result.filter(t => t.playbook_id === playbookFilter.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(t =>
      t.playbook_name?.toLowerCase().includes(q) ||
      t.id.toString().includes(q)
    )
  }
  return result
})

async function fetchTasks() {
  loading.value = true
  try {
    const res = await get('/tasks')
    tasks.value = res.data?.items || []
  } catch (e) {
    console.error('Failed to fetch tasks:', e)
    tasks.value = []
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

async function createTask() {
  try {
    const payload = {
      playbook_id: newTask.value.playbook_id,
      target_nodes: newTask.value.target_nodes,
    }
    if (newTask.value.variables) {
      try {
        payload.variables = JSON.parse(newTask.value.variables)
      } catch (e) {
        alert('Invalid JSON in variables')
        return
      }
    }
    await post('/tasks', payload)
    showCreateDialog.value = false
    newTask.value = { playbook_id: '', target_nodes: [], variables: '' }
    fetchTasks()
  } catch (e) {
    alert('Failed to create task: ' + e.message)
  }
}

async function cancelTask(task) {
  if (!confirm('Cancel this task?')) return
  try {
    await post(`/tasks/${task.id}/cancel`)
    fetchTasks()
  } catch (e) {
    alert('Failed to cancel task')
  }
}

async function retryTask(task) {
  try {
    await post(`/tasks/${task.id}/retry`)
    fetchTasks()
  } catch (e) {
    alert('Failed to retry task')
  }
}

async function viewTask(task) {
  try {
    const res = await get(`/tasks/${task.id}`)
    selectedTask.value = res.data
  } catch (e) {
    selectedTask.value = task
  }
}

function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

function formatDuration(task) {
  if (!task.started_at) return '-'
  const start = new Date(task.started_at)
  const end = task.completed_at ? new Date(task.completed_at) : new Date()
  const diff = Math.floor((end - start) / 1000)
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ${diff % 60}s`
  return `${Math.floor(diff / 3600)}h ${Math.floor((diff % 3600) / 60)}m`
}

onMounted(() => {
  fetchTasks()
  fetchPlaybooks()
  fetchNodes()
})
</script>

<style scoped>
.tasks-page {
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
  grid-template-columns: repeat(4, 1fr);
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

.stat-icon.pending { background: rgba(255, 193, 7, 0.2); }
.stat-icon.running { background: rgba(33, 150, 243, 0.2); }
.stat-icon.completed { background: rgba(76, 175, 80, 0.2); }
.stat-icon.failed { background: rgba(244, 67, 54, 0.2); }

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

.filters-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: var(--card-bg);
  border-radius: 8px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-group label {
  color: var(--text-secondary);
  font-size: 14px;
}

.filter-group select {
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
}

.search-group input {
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
  width: 250px;
}

.table-container {
  background: var(--card-bg);
  border-radius: 8px;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.data-table th {
  background: var(--header-bg);
  font-weight: 600;
  color: var(--text-secondary);
}

.data-table tbody tr:hover {
  background: var(--hover-bg);
  cursor: pointer;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
}

.status-badge.pending { background: rgba(255, 193, 7, 0.2); color: #ffc107; }
.status-badge.running { background: rgba(33, 150, 243, 0.2); color: #2196f3; }
.status-badge.completed { background: rgba(76, 175, 80, 0.2); color: #4caf50; }
.status-badge.failed { background: rgba(244, 67, 54, 0.2); color: #f44336; }
.status-badge.cancelled { background: rgba(158, 158, 158, 0.2); color: #9e9e9e; }

.actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  font-size: 16px;
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

.loading-row,
.empty-row {
  text-align: center;
  color: var(--text-secondary);
  padding: 40px;
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

.modal-lg {
  width: 700px;
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

.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  background: var(--input-bg);
  color: var(--text);
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

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item label {
  color: var(--text-secondary);
  font-size: 12px;
}

.logs-section h3 {
  margin: 0 0 12px;
  font-size: 14px;
  color: var(--text-secondary);
}

.logs-container {
  background: #1e1e1e;
  border-radius: 8px;
  padding: 16px;
  max-height: 300px;
  overflow: auto;
}

.logs-container pre {
  margin: 0;
  color: #d4d4d4;
  font-family: monospace;
  font-size: 12px;
  white-space: pre-wrap;
}
</style>