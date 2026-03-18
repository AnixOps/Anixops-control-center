<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">{{ t('nodes.title') }}</h1>
        <p class="text-dark-400 text-sm mt-1">{{ t('nodes.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="showBatchActions = !showBatchActions"
          v-if="selectedNodes.length > 0"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
        >
          {{ selectedNodes.length }} {{ t('common.selected') || 'selected' }}
        </button>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <PlusIcon class="w-5 h-5" />
          {{ t('nodes.addNode') }}
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('nodes.total') }}</p>
            <p class="text-2xl font-bold text-white">{{ nodeStats.total }}</p>
          </div>
          <ServerIcon class="w-8 h-8 text-primary-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('nodes.online') }}</p>
            <p class="text-2xl font-bold text-green-400">{{ nodeStats.online }}</p>
          </div>
          <CheckCircleIcon class="w-8 h-8 text-green-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('nodes.offline') }}</p>
            <p class="text-2xl font-bold text-red-400">{{ nodeStats.offline }}</p>
          </div>
          <XCircleIcon class="w-8 h-8 text-red-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('nodes.traffic') }}</p>
            <p class="text-2xl font-bold text-white">{{ formatBytes(totalTraffic) }}</p>
          </div>
          <ChartBarIcon class="w-8 h-8 text-blue-400" />
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-4">
      <div class="flex-1 min-w-[200px]">
        <div class="relative">
          <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-400" />
          <input
            v-model="search"
            type="text"
            :placeholder="t('common.search') + '...'"
            class="w-full pl-10 pr-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-400 focus:outline-none focus:ring-2 focus:ring-primary-500"
            @input="debouncedSearch"
          />
        </div>
      </div>
      <select
        v-model="statusFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        @change="handleFilterChange"
      >
        <option value="">{{ t('common.all') }} {{ t('common.status') }}</option>
        <option value="online">{{ t('common.online') }}</option>
        <option value="offline">{{ t('common.offline') }}</option>
        <option value="starting">Starting</option>
        <option value="stopping">Stopping</option>
      </select>
      <select
        v-model="typeFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        @change="handleFilterChange"
      >
        <option value="">{{ t('common.all') }} {{ t('nodes.type') }}</option>
        <option value="v2ray">V2Ray</option>
        <option value="xray">XRay</option>
        <option value="trojan">Trojan</option>
        <option value="shadowsocks">Shadowsocks</option>
      </select>
      <button
        @click="refreshNodes"
        :disabled="loading"
        class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
      >
        <ArrowPathIcon :class="['w-5 h-5', { 'animate-spin': loading }]" />
        {{ t('common.refresh') }}
      </button>
    </div>

    <!-- Nodes Table -->
    <div class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-dark-700">
            <tr>
              <th class="px-4 py-3 text-left">
                <input
                  type="checkbox"
                  :checked="allSelected"
                  @change="toggleSelectAll"
                  class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
                />
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('nodes.nodeName') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('nodes.type') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('common.status') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('nodes.users') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('nodes.traffic') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('nodes.uptime') }}</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-dark-700">
            <tr v-if="loading">
              <td colspan="8" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center gap-3">
                  <ArrowPathIcon class="w-8 h-8 text-primary-400 animate-spin" />
                  <p class="text-dark-400">{{ t('common.loading') }}</p>
                </div>
              </td>
            </tr>
            <tr v-else-if="filteredNodes.length === 0">
              <td colspan="8" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center gap-3">
                  <ServerIcon class="w-12 h-12 text-dark-500" />
                  <p class="text-dark-400">{{ t('common.noResults') }}</p>
                  <button
                    @click="openCreateModal"
                    class="text-primary-400 hover:text-primary-300"
                  >
                    {{ t('nodes.addNode') }}
                  </button>
                </div>
              </td>
            </tr>
            <tr
              v-else
              v-for="node in filteredNodes"
              :key="node.id"
              class="hover:bg-dark-700/50 transition-colors"
              :class="{ 'bg-primary-900/10': selectedNodes.includes(node.id) }"
            >
              <td class="px-4 py-4">
                <input
                  type="checkbox"
                  :checked="selectedNodes.includes(node.id)"
                  @change="toggleSelect(node.id)"
                  class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
                />
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-primary-600/20 flex items-center justify-center">
                    <ServerIcon class="w-5 h-5 text-primary-400" />
                  </div>
                  <div>
                    <p class="text-white font-medium">{{ node.name }}</p>
                    <p class="text-dark-400 text-sm">{{ node.host }}:{{ node.port || 443 }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs rounded bg-dark-600 text-dark-300">
                  {{ node.type || 'v2ray' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <div
                    class="w-2 h-2 rounded-full"
                    :class="getStatusDotColor(node.status)"
                  ></div>
                  <span
                    class="px-2 py-1 text-xs rounded-full"
                    :class="getStatusColor(node.status)"
                  >
                    {{ node.status }}
                  </span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-dark-300">{{ node.users || 0 }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-dark-300">{{ formatBytes(node.traffic || 0) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-dark-300">{{ node.uptime || '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    v-if="node.status === 'offline'"
                    @click="handleStartNode(node)"
                    :disabled="actionLoading[node.id]"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Start"
                  >
                    <PlayIcon class="w-4 h-4 text-green-400" />
                  </button>
                  <button
                    v-if="node.status === 'online'"
                    @click="handleStopNode(node)"
                    :disabled="actionLoading[node.id]"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Stop"
                  >
                    <StopIcon class="w-4 h-4 text-yellow-400" />
                  </button>
                  <button
                    @click="handleRestartNode(node)"
                    :disabled="actionLoading[node.id]"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Restart"
                  >
                    <ArrowPathIcon class="w-4 h-4 text-blue-400" />
                  </button>
                  <button
                    @click="openEditModal(node)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <PencilIcon class="w-4 h-4 text-dark-400" />
                  </button>
                  <button
                    @click="confirmDelete(node)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
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

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="px-6 py-4 border-t border-dark-700 flex items-center justify-between">
        <p class="text-dark-400 text-sm">
          {{ t('common.showing') || 'Showing' }} {{ (currentPage - 1) * pageSize + 1 }} {{ t('common.to') || 'to' }} {{ Math.min(currentPage * pageSize, total) }} {{ t('common.of') || 'of' }} {{ total }} {{ t('nodes.title').toLowerCase() }}
        </p>
        <div class="flex items-center gap-2">
          <button
            @click="changePage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            {{ t('common.previous') }}
          </button>
          <span class="text-dark-400">{{ currentPage }} / {{ totalPages }}</span>
          <button
            @click="changePage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            {{ t('common.next') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeModal"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">
              {{ editingNode ? t('nodes.editNode') : t('nodes.addNode') }}
            </h2>
            <button @click="closeModal" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
            <!-- Basic Info -->
            <div class="space-y-4">
              <h3 class="text-lg font-medium text-white">{{ t('settings.server') }}</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-dark-400 mb-1">{{ t('nodes.nodeName') }} *</label>
                  <input
                    v-model="formData.name"
                    type="text"
                    required
                    placeholder="e.g., tokyo-01"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">{{ t('nodes.type') }} *</label>
                  <select
                    v-model="formData.type"
                    required
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="v2ray">V2Ray</option>
                    <option value="xray">XRay</option>
                    <option value="trojan">Trojan</option>
                    <option value="shadowsocks">Shadowsocks</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">{{ t('nodes.host') }} *</label>
                  <input
                    v-model="formData.host"
                    type="text"
                    required
                    placeholder="e.g., 192.168.1.1 or domain.com"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">{{ t('nodes.port') }} *</label>
                  <input
                    v-model.number="formData.port"
                    type="number"
                    required
                    placeholder="443"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>

            <!-- Network Settings -->
            <div class="space-y-4">
              <h3 class="text-lg font-medium text-white">Network Settings</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Network</label>
                  <select
                    v-model="formData.network"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="tcp">TCP</option>
                    <option value="ws">WebSocket</option>
                    <option value="grpc">gRPC</option>
                    <option value="http2">HTTP/2</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Security</label>
                  <select
                    v-model="formData.security"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="tls">TLS</option>
                    <option value="reality">Reality</option>
                    <option value="none">None</option>
                  </select>
                </div>
                <div v-if="formData.security === 'tls'" class="col-span-2">
                  <label class="block text-sm text-dark-400 mb-1">Server Name (SNI)</label>
                  <input
                    v-model="formData.serverName"
                    type="text"
                    placeholder="domain.com"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>

            <!-- Rate Limit -->
            <div class="space-y-4">
              <h3 class="text-lg font-medium text-white">Rate Limit</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Upload Limit (MB/s)</label>
                  <input
                    v-model.number="formData.uploadLimit"
                    type="number"
                    min="0"
                    placeholder="0 = unlimited"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Download Limit (MB/s)</label>
                  <input
                    v-model.number="formData.downloadLimit"
                    type="number"
                    min="0"
                    placeholder="0 = unlimited"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>

            <!-- Tags -->
            <div>
              <label class="block text-sm text-dark-400 mb-1">Tags</label>
              <input
                v-model="formData.tags"
                type="text"
                placeholder="Comma-separated tags: premium, asia, beta"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>

            <!-- Error Message -->
            <div v-if="formError" class="p-4 bg-red-900/20 border border-red-800 rounded-lg">
              <p class="text-red-400 text-sm">{{ formError }}</p>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3 pt-4 border-t border-dark-700">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
              >
                <ArrowPathIcon v-if="submitting" class="w-4 h-4 animate-spin" />
                {{ editingNode ? t('common.update') : t('common.create') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showDeleteConfirm = false"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-md p-6">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-12 h-12 rounded-full bg-red-900/30 flex items-center justify-center">
              <ExclamationTriangleIcon class="w-6 h-6 text-red-400" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-white">{{ t('nodes.deleteNode') }}</h3>
              <p class="text-dark-400 text-sm">{{ t('common.confirm') }}</p>
            </div>
          </div>
          <p class="text-dark-300 mb-6">
            {{ t('nodes.confirmDelete', { name: nodeToDelete?.name }) }}
          </p>
          <div class="flex items-center justify-end gap-3">
            <button
              @click="showDeleteConfirm = false"
              class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
            >
              {{ t('common.cancel') }}
            </button>
            <button
              @click="handleDelete"
              :disabled="deleting"
              class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors flex items-center gap-2"
            >
              <ArrowPathIcon v-if="deleting" class="w-4 h-4 animate-spin" />
              {{ t('common.delete') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Toast Notifications -->
    <Teleport to="body">
      <div class="fixed bottom-4 right-4 z-50 space-y-2">
        <TransitionGroup name="toast">
          <div
            v-for="toast in toasts"
            :key="toast.id"
            class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 min-w-[300px]"
            :class="{
              'bg-green-900/90 border border-green-700': toast.type === 'success',
              'bg-red-900/90 border border-red-700': toast.type === 'error',
              'bg-yellow-900/90 border border-yellow-700': toast.type === 'warning'
            }"
          >
            <CheckCircleIcon v-if="toast.type === 'success'" class="w-5 h-5 text-green-400" />
            <XCircleIcon v-else-if="toast.type === 'error'" class="w-5 h-5 text-red-400" />
            <ExclamationTriangleIcon v-else class="w-5 h-5 text-yellow-400" />
            <span class="text-white flex-1">{{ toast.message }}</span>
            <button @click="removeToast(toast.id)" class="text-dark-400 hover:text-white">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </TransitionGroup>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNodesStore } from '@/stores/nodes'
import { useWebSocket } from '@/services/websocket'
import {
  ServerIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
  PlayIcon,
  StopIcon,
  ExclamationTriangleIcon,
  ChartBarIcon
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const nodesStore = useNodesStore()
const { subscribe, unsubscribe } = useWebSocket()

// State
const search = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const showModal = ref(false)
const showDeleteConfirm = ref(false)
const showBatchActions = ref(false)
const editingNode = ref(null)
const nodeToDelete = ref(null)
const selectedNodes = ref([])
const actionLoading = ref({})
const submitting = ref(false)
const deleting = ref(false)
const formError = ref('')
const toasts = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const searchTimeout = ref(null)

// Form data
const defaultFormData = {
  name: '',
  type: 'v2ray',
  host: '',
  port: 443,
  network: 'tcp',
  security: 'tls',
  serverName: '',
  uploadLimit: 0,
  downloadLimit: 0,
  tags: ''
}
const formData = ref({ ...defaultFormData })

// Computed
const loading = computed(() => nodesStore.loading)
const nodes = computed(() => nodesStore.nodes)
const total = computed(() => nodesStore.total)
const nodeStats = computed(() => nodesStore.nodeStats)
const totalTraffic = computed(() => {
  return nodes.value.reduce((sum, node) => sum + (node.traffic || 0), 0)
})
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const filteredNodes = computed(() => {
  let result = nodes.value

  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(node =>
      node.name.toLowerCase().includes(searchLower) ||
      node.host.includes(search.value)
    )
  }

  if (statusFilter.value) {
    result = result.filter(node => node.status === statusFilter.value)
  }

  if (typeFilter.value) {
    result = result.filter(node => node.type === typeFilter.value)
  }

  return result
})

const allSelected = computed(() => {
  return filteredNodes.value.length > 0 &&
    filteredNodes.value.every(node => selectedNodes.value.includes(node.id))
})

// Methods
function getStatusColor(status) {
  const colors = {
    online: 'bg-green-900/30 text-green-400',
    offline: 'bg-red-900/30 text-red-400',
    starting: 'bg-yellow-900/30 text-yellow-400',
    stopping: 'bg-yellow-900/30 text-yellow-400',
    error: 'bg-red-900/30 text-red-400'
  }
  return colors[status] || 'bg-gray-900/30 text-gray-400'
}

function getStatusDotColor(status) {
  const colors = {
    online: 'bg-green-400',
    offline: 'bg-red-400',
    starting: 'bg-yellow-400',
    stopping: 'bg-yellow-400',
    error: 'bg-red-400'
  }
  return colors[status] || 'bg-gray-400'
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function debouncedSearch() {
  clearTimeout(searchTimeout.value)
  searchTimeout.value = setTimeout(() => {
    handleFilterChange()
  }, 300)
}

function handleFilterChange() {
  currentPage.value = 1
  nodesStore.setFilters({
    search: search.value,
    status: statusFilter.value,
    type: typeFilter.value,
    page: currentPage.value,
    limit: pageSize.value
  })
}

async function refreshNodes() {
  try {
    await nodesStore.fetchNodes()
    showToast(t('common.success'), 'success')
  } catch (error) {
    showToast(error.message || t('common.error'), 'error')
  }
}

function openCreateModal() {
  editingNode.value = null
  formData.value = { ...defaultFormData }
  formError.value = ''
  showModal.value = true
}

function openEditModal(node) {
  editingNode.value = node
  formData.value = {
    name: node.name,
    type: node.type || 'v2ray',
    host: node.host,
    port: node.port || 443,
    network: node.network || 'tcp',
    security: node.security || 'tls',
    serverName: node.serverName || '',
    uploadLimit: node.uploadLimit || 0,
    downloadLimit: node.downloadLimit || 0,
    tags: (node.tags || []).join(', ')
  }
  formError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingNode.value = null
  formData.value = { ...defaultFormData }
  formError.value = ''
}

async function handleSubmit() {
  submitting.value = true
  formError.value = ''

  try {
    const data = {
      ...formData.value,
      tags: formData.value.tags
        .split(',')
        .map(t => t.trim())
        .filter(Boolean)
    }

    if (editingNode.value) {
      await nodesStore.updateNode(editingNode.value.id, data)
      showToast(t('common.updateSuccess'), 'success')
    } else {
      await nodesStore.createNode(data)
      showToast(t('common.createSuccess'), 'success')
    }
    closeModal()
  } catch (error) {
    formError.value = error.message || 'Operation failed'
  } finally {
    submitting.value = false
  }
}

function confirmDelete(node) {
  nodeToDelete.value = node
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!nodeToDelete.value) return

  deleting.value = true
  try {
    await nodesStore.deleteNode(nodeToDelete.value.id)
    showToast(t('common.deleteSuccess'), 'success')
    showDeleteConfirm.value = false
    nodeToDelete.value = null
  } catch (error) {
    showToast(error.message || t('common.deleteError'), 'error')
  } finally {
    deleting.value = false
  }
}

async function handleStartNode(node) {
  actionLoading.value[node.id] = true
  try {
    await nodesStore.startNode(node.id)
    showToast(t('nodes.starting'), 'success')
  } catch (error) {
    showToast(error.message || t('common.error'), 'error')
  } finally {
    actionLoading.value[node.id] = false
  }
}

async function handleStopNode(node) {
  actionLoading.value[node.id] = true
  try {
    await nodesStore.stopNode(node.id)
    showToast(t('nodes.stopping'), 'success')
  } catch (error) {
    showToast(error.message || t('common.error'), 'error')
  } finally {
    actionLoading.value[node.id] = false
  }
}

async function handleRestartNode(node) {
  actionLoading.value[node.id] = true
  try {
    await nodesStore.restartNode(node.id)
    showToast(t('nodes.restarting'), 'success')
  } catch (error) {
    showToast(error.message || t('common.error'), 'error')
  } finally {
    actionLoading.value[node.id] = false
  }
}

function toggleSelect(id) {
  const index = selectedNodes.value.indexOf(id)
  if (index === -1) {
    selectedNodes.value.push(id)
  } else {
    selectedNodes.value.splice(index, 1)
  }
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedNodes.value = []
  } else {
    selectedNodes.value = filteredNodes.value.map(n => n.id)
  }
}

function changePage(page) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  nodesStore.setFilters({ page })
}

function showToast(message, type = 'success') {
  const id = Date.now()
  toasts.value.push({ id, message, type })
  setTimeout(() => removeToast(id), 5000)
}

function removeToast(id) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

// Lifecycle
onMounted(async () => {
  await refreshNodes()

  // Subscribe to real-time node updates
  subscribe('node:status', (data) => {
    nodesStore.updateNodeStatus(data.id, data.status)
  })

  subscribe('node:stats', (data) => {
    const index = nodes.value.findIndex(n => n.id === data.id)
    if (index !== -1) {
      nodes.value[index] = { ...nodes.value[index], ...data }
    }
  })
})

onUnmounted(() => {
  unsubscribe('node:status')
  unsubscribe('node:stats')
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>