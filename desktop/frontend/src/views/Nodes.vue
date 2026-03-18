<template>
  <div class="space-y-6">
    <!-- Offline Indicator -->
    <OfflineIndicator />

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Nodes</h2>
        <p class="text-slate-400 text-sm mt-1">Manage your server nodes</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="refreshNodes"
          :disabled="loading"
          class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <svg
            :class="['w-5 h-5', { 'animate-spin': loading }]"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
        <button
          @click="showAddModal = true"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <PlusIcon class="w-5 h-5" />
          Add Node
        </button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-slate-800 rounded-xl p-4 border border-slate-700">
        <p class="text-slate-400 text-sm">Total</p>
        <p class="text-2xl font-bold text-white">{{ nodeStats.total }}</p>
      </div>
      <div class="bg-slate-800 rounded-xl p-4 border border-slate-700">
        <p class="text-slate-400 text-sm">Online</p>
        <p class="text-2xl font-bold text-green-400">{{ nodeStats.online }}</p>
      </div>
      <div class="bg-slate-800 rounded-xl p-4 border border-slate-700">
        <p class="text-slate-400 text-sm">Offline</p>
        <p class="text-2xl font-bold text-red-400">{{ nodeStats.offline }}</p>
      </div>
      <div class="bg-slate-800 rounded-xl p-4 border border-slate-700">
        <p class="text-slate-400 text-sm">Total Traffic</p>
        <p class="text-2xl font-bold text-white">{{ formatBytes(totalTraffic) }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex gap-4">
      <div class="flex-1">
        <div class="relative">
          <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
          <input
            v-model="search"
            type="text"
            placeholder="Search nodes..."
            class="w-full pl-10 pr-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>
      <select
        v-model="statusFilter"
        class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">All Status</option>
        <option value="online">Online</option>
        <option value="offline">Offline</option>
      </select>
    </div>

    <!-- Nodes Table -->
    <div class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden">
      <table class="w-full">
        <thead class="bg-slate-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Name</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Host</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Users</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Traffic</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-700">
          <tr v-if="loading">
            <td colspan="6" class="px-6 py-12 text-center">
              <div class="flex flex-col items-center gap-3">
                <svg class="w-8 h-8 text-blue-400 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p class="text-slate-400">Loading nodes...</p>
              </div>
            </td>
          </tr>
          <tr v-else-if="filteredNodes.length === 0">
            <td colspan="6" class="px-6 py-12 text-center">
              <div class="flex flex-col items-center gap-3">
                <ServerIcon class="w-12 h-12 text-slate-600" />
                <p class="text-slate-400">No nodes found</p>
                <button
                  @click="showAddModal = true"
                  class="text-blue-400 hover:text-blue-300"
                >
                  Add your first node
                </button>
              </div>
            </td>
          </tr>
          <tr
            v-else
            v-for="node in filteredNodes"
            :key="node.id"
            class="hover:bg-slate-700/50"
          >
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-blue-600/20 flex items-center justify-center">
                  <ServerIcon class="w-4 h-4 text-blue-400" />
                </div>
                <div>
                  <span class="text-white font-medium">{{ node.name }}</span>
                  <p class="text-slate-500 text-xs">{{ node.type || 'v2ray' }}</p>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-slate-400">{{ node.host }}:{{ node.port || 443 }}</td>
            <td class="px-6 py-4">
              <span
                class="px-2 py-1 text-xs rounded-full"
                :class="getStatusColor(node.status)"
              >
                {{ node.status }}
              </span>
            </td>
            <td class="px-6 py-4 text-slate-400">{{ node.users || 0 }}</td>
            <td class="px-6 py-4 text-slate-400">{{ formatBytes(node.traffic || 0) }}</td>
            <td class="px-6 py-4 text-right">
              <div class="flex items-center justify-end gap-2">
                <button
                  v-if="node.status === 'offline'"
                  @click="startNode(node)"
                  :disabled="actionLoading[node.id]"
                  class="p-2 hover:bg-slate-600 rounded-lg transition-colors"
                  title="Start"
                >
                  <PlayIcon class="w-4 h-4 text-green-400" />
                </button>
                <button
                  v-else
                  @click="stopNode(node)"
                  :disabled="actionLoading[node.id]"
                  class="p-2 hover:bg-slate-600 rounded-lg transition-colors"
                  title="Stop"
                >
                  <StopIcon class="w-4 h-4 text-yellow-400" />
                </button>
                <button
                  @click="restartNode(node)"
                  :disabled="actionLoading[node.id]"
                  class="p-2 hover:bg-slate-600 rounded-lg transition-colors"
                  title="Restart"
                >
                  <RefreshIcon class="w-4 h-4 text-blue-400" />
                </button>
                <button
                  @click="editNode(node)"
                  class="p-2 hover:bg-slate-600 rounded-lg transition-colors"
                  title="Edit"
                >
                  <EditIcon class="w-4 h-4 text-slate-400" />
                </button>
                <button
                  @click="deleteNode(node)"
                  class="p-2 hover:bg-slate-600 rounded-lg transition-colors"
                  title="Delete"
                >
                  <TrashIcon class="w-4 h-4 text-red-400" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add/Edit Modal -->
    <Teleport to="body">
      <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showAddModal = false"></div>
        <div class="relative bg-slate-800 rounded-xl border border-slate-700 w-full max-w-md p-6">
          <h3 class="text-lg font-semibold text-white mb-4">
            {{ editingNode ? 'Edit Node' : 'Add New Node' }}
          </h3>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div>
              <label class="block text-sm text-slate-400 mb-1">Name *</label>
              <input
                v-model="formData.name"
                type="text"
                required
                class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-slate-400 mb-1">Host *</label>
                <input
                  v-model="formData.host"
                  type="text"
                  required
                  class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label class="block text-sm text-slate-400 mb-1">Port *</label>
                <input
                  v-model.number="formData.port"
                  type="number"
                  required
                  class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
            <div>
              <label class="block text-sm text-slate-400 mb-1">Type</label>
              <select
                v-model="formData.type"
                class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="v2ray">V2Ray</option>
                <option value="xray">XRay</option>
                <option value="trojan">Trojan</option>
                <option value="shadowsocks">Shadowsocks</option>
              </select>
            </div>
            <div class="flex items-center justify-end gap-3 pt-4">
              <button
                type="button"
                @click="showAddModal = false"
                class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
              >
                {{ submitting ? 'Saving...' : (editingNode ? 'Update' : 'Create') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation -->
    <Teleport to="body">
      <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showDeleteConfirm = false"></div>
        <div class="relative bg-slate-800 rounded-xl border border-slate-700 w-full max-w-sm p-6">
          <h3 class="text-lg font-semibold text-white mb-2">Delete Node</h3>
          <p class="text-slate-400 mb-4">Are you sure you want to delete "{{ nodeToDelete?.name }}"?</p>
          <div class="flex items-center justify-end gap-3">
            <button
              @click="showDeleteConfirm = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              @click="confirmDelete"
              :disabled="deleting"
              class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
            >
              {{ deleting ? 'Deleting...' : 'Delete' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, h, onMounted, onUnmounted } from 'vue'
import { useNodesStore } from '../stores/nodes'
import { useOfflineStore } from '../stores/offline'
import OfflineIndicator from '../components/OfflineIndicator.vue'

const nodesStore = useNodesStore()
const offlineStore = useOfflineStore()

// State
const loading = ref(false)
const search = ref('')
const statusFilter = ref('')
const showAddModal = ref(false)
const showDeleteConfirm = ref(false)
const editingNode = ref(null)
const nodeToDelete = ref(null)
const submitting = ref(false)
const deleting = ref(false)
const actionLoading = ref({})

const formData = ref({
  name: '',
  host: '',
  port: 443,
  type: 'v2ray'
})

// Computed
const nodes = computed(() => nodesStore.nodes)
const nodeStats = computed(() => nodesStore.nodeStats)
const totalTraffic = computed(() => nodes.value.reduce((sum, n) => sum + (n.traffic || 0), 0))

const filteredNodes = computed(() => {
  let result = nodes.value
  if (search.value) {
    result = result.filter(n =>
      n.name.toLowerCase().includes(search.value.toLowerCase()) ||
      n.host.includes(search.value)
    )
  }
  if (statusFilter.value) {
    result = result.filter(n => n.status === statusFilter.value)
  }
  return result
})

// Icons
const PlusIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4v16m8-8H4' })
])

const SearchIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z' })
])

const ServerIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01' })
])

const PlayIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z' }),
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
])

const StopIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M21 12a9 9 0 11-18 0 9 9 0 0118 0z' }),
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z' })
])

const RefreshIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15' })
])

const EditIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z' })
])

const TrashIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16' })
])

// Methods
function getStatusColor(status) {
  return status === 'online'
    ? 'bg-green-900/30 text-green-400'
    : 'bg-red-900/30 text-red-400'
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function refreshNodes() {
  loading.value = true
  try {
    await nodesStore.fetchNodes()
    // Cache nodes for offline use
    await offlineStore.cacheNodes(nodes.value)
  } catch (error) {
    console.error('Failed to fetch nodes:', error)
    // Try to load from cache if offline
    if (!offlineStore.isOnline.value) {
      const cachedNodes = await offlineStore.getCachedNodes()
      nodesStore.setNodes(cachedNodes)
    }
  } finally {
    loading.value = false
  }
}

async function startNode(node) {
  actionLoading.value[node.id] = true
  try {
    await nodesStore.startNode(node.id)
  } catch (error) {
    // Add to pending actions if offline
    if (!offlineStore.isOnline.value) {
      await offlineStore.addPendingAction('update', 'node', { id: node.id, status: 'online' })
    }
    console.error('Failed to start node:', error)
  } finally {
    actionLoading.value[node.id] = false
  }
}

async function stopNode(node) {
  actionLoading.value[node.id] = true
  try {
    await nodesStore.stopNode(node.id)
  } catch (error) {
    if (!offlineStore.isOnline.value) {
      await offlineStore.addPendingAction('update', 'node', { id: node.id, status: 'offline' })
    }
    console.error('Failed to stop node:', error)
  } finally {
    actionLoading.value[node.id] = false
  }
}

async function restartNode(node) {
  actionLoading.value[node.id] = true
  try {
    await nodesStore.restartNode(node.id)
  } catch (error) {
    console.error('Failed to restart node:', error)
  } finally {
    actionLoading.value[node.id] = false
  }
}

function editNode(node) {
  editingNode.value = node
  formData.value = {
    name: node.name,
    host: node.host,
    port: node.port || 443,
    type: node.type || 'v2ray'
  }
  showAddModal.value = true
}

function deleteNode(node) {
  nodeToDelete.value = node
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  if (!nodeToDelete.value) return
  deleting.value = true
  try {
    await nodesStore.deleteNode(nodeToDelete.value.id)
    showDeleteConfirm.value = false
    nodeToDelete.value = null
  } catch (error) {
    if (!offlineStore.isOnline.value) {
      await offlineStore.addPendingAction('delete', 'node', { id: nodeToDelete.value.id })
      showDeleteConfirm.value = false
      nodeToDelete.value = null
    }
    console.error('Failed to delete node:', error)
  } finally {
    deleting.value = false
  }
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (editingNode.value) {
      await nodesStore.updateNode(editingNode.value.id, formData.value)
    } else {
      await nodesStore.createNode(formData.value)
    }
    showAddModal.value = false
    editingNode.value = null
    formData.value = { name: '', host: '', port: 443, type: 'v2ray' }
  } catch (error) {
    if (!offlineStore.isOnline.value) {
      await offlineStore.addPendingAction(
        editingNode.value ? 'update' : 'create',
        'node',
        formData.value
      )
    }
    console.error('Failed to save node:', error)
  } finally {
    submitting.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await offlineStore.initialize()
  await refreshNodes()
})

onUnmounted(() => {
  offlineStore.cleanup()
})
</script>