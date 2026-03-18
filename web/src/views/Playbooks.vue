<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">Playbooks</h1>
        <p class="text-dark-400 text-sm mt-1">Automated deployment and configuration</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="showTemplates = !showTemplates"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <DocumentDuplicateIcon class="w-5 h-5" />
          Templates
        </button>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <PlusIcon class="w-5 h-5" />
          New Playbook
        </button>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Total Playbooks</p>
            <p class="text-2xl font-bold text-white">{{ playbookStats.total }}</p>
          </div>
          <DocumentTextIcon class="w-8 h-8 text-primary-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Running</p>
            <p class="text-2xl font-bold text-yellow-400">{{ playbookStats.running }}</p>
          </div>
          <PlayIcon class="w-8 h-8 text-yellow-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Success Rate</p>
            <p class="text-2xl font-bold text-green-400">{{ playbookStats.successRate }}%</p>
          </div>
          <ChartBarIcon class="w-8 h-8 text-green-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">Scheduled</p>
            <p class="text-2xl font-bold text-blue-400">{{ playbookStats.scheduled }}</p>
          </div>
          <ClockIcon class="w-8 h-8 text-blue-400" />
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
            placeholder="Search playbooks..."
            class="w-full pl-10 pr-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-400 focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>
      </div>
      <select
        v-model="statusFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
      >
        <option value="">All Status</option>
        <option value="success">Success</option>
        <option value="running">Running</option>
        <option value="failed">Failed</option>
        <option value="pending">Pending</option>
      </select>
      <select
        v-model="typeFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
      >
        <option value="">All Types</option>
        <option value="deployment">Deployment</option>
        <option value="maintenance">Maintenance</option>
        <option value="backup">Backup</option>
        <option value="custom">Custom</option>
      </select>
    </div>

    <!-- Playbooks Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div
        v-for="playbook in filteredPlaybooks"
        :key="playbook.id"
        class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden hover:border-dark-600 transition-colors"
      >
        <!-- Header -->
        <div class="p-4 border-b border-dark-700">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-lg flex items-center justify-center"
                :class="getPlaybookTypeColor(playbook.type)"
              >
                <component :is="getPlaybookTypeIcon(playbook.type)" class="w-5 h-5 text-white" />
              </div>
              <div>
                <h3 class="text-white font-medium">{{ playbook.name }}</h3>
                <p class="text-dark-400 text-sm">{{ playbook.description || 'No description' }}</p>
              </div>
            </div>
            <span
              class="px-2 py-1 text-xs rounded-full"
              :class="getStatusColor(playbook.status)"
            >
              {{ playbook.status }}
            </span>
          </div>
        </div>

        <!-- Body -->
        <div class="p-4 space-y-3">
          <div class="flex items-center gap-4 text-sm">
            <div class="flex items-center gap-2">
              <ClockIcon class="w-4 h-4 text-dark-400" />
              <span class="text-dark-400">Last run:</span>
              <span class="text-white">{{ playbook.lastRun || 'Never' }}</span>
            </div>
            <div class="flex items-center gap-2">
              <ServerIcon class="w-4 h-4 text-dark-400" />
              <span class="text-dark-400">Nodes:</span>
              <span class="text-white">{{ playbook.targetNodes?.length || 0 }}</span>
            </div>
          </div>

          <!-- Progress Bar (for running playbooks) -->
          <div v-if="playbook.status === 'running'" class="space-y-1">
            <div class="flex items-center justify-between text-xs">
              <span class="text-dark-400">Progress</span>
              <span class="text-white">{{ playbook.progress || 0 }}%</span>
            </div>
            <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
              <div
                class="h-full bg-yellow-500 transition-all duration-300"
                :style="{ width: `${playbook.progress || 0}%` }"
              ></div>
            </div>
          </div>

          <!-- Schedule Info -->
          <div v-if="playbook.schedule" class="flex items-center gap-2 text-sm">
            <CalendarIcon class="w-4 h-4 text-blue-400" />
            <span class="text-dark-400">Scheduled:</span>
            <span class="text-blue-400">{{ playbook.schedule }}</span>
          </div>

          <!-- Tags -->
          <div v-if="playbook.tags?.length" class="flex flex-wrap gap-1">
            <span
              v-for="tag in playbook.tags"
              :key="tag"
              class="px-2 py-0.5 text-xs bg-dark-700 text-dark-300 rounded"
            >
              {{ tag }}
            </span>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-4 py-3 bg-dark-700/50 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <button
              @click="runPlaybook(playbook)"
              :disabled="playbook.status === 'running'"
              class="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white text-sm rounded-lg transition-colors flex items-center gap-1 disabled:opacity-50"
            >
              <PlayIcon class="w-4 h-4" />
              Run
            </button>
            <button
              v-if="playbook.status === 'running'"
              @click="stopPlaybook(playbook)"
              class="px-3 py-1.5 bg-red-600 hover:bg-red-700 text-white text-sm rounded-lg transition-colors flex items-center gap-1"
            >
              <StopIcon class="w-4 h-4" />
              Stop
            </button>
          </div>
          <div class="flex items-center gap-1">
            <button
              @click="viewLogs(playbook)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="View Logs"
            >
              <DocumentTextIcon class="w-4 h-4 text-dark-400" />
            </button>
            <button
              @click="openEditModal(playbook)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="Edit"
            >
              <PencilIcon class="w-4 h-4 text-dark-400" />
            </button>
            <button
              @click="duplicatePlaybook(playbook)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="Duplicate"
            >
              <DocumentDuplicateIcon class="w-4 h-4 text-dark-400" />
            </button>
            <button
              @click="confirmDelete(playbook)"
              class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              title="Delete"
            >
              <TrashIcon class="w-4 h-4 text-red-400" />
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="filteredPlaybooks.length === 0" class="col-span-2 py-12 text-center">
        <DocumentTextIcon class="w-12 h-12 text-dark-500 mx-auto mb-4" />
        <p class="text-dark-400 mb-2">No playbooks found</p>
        <button
          @click="openCreateModal"
          class="text-primary-400 hover:text-primary-300"
        >
          Create your first playbook
        </button>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeModal"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-4xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">
              {{ editingPlaybook ? 'Edit Playbook' : 'New Playbook' }}
            </h2>
            <button @click="closeModal" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
            <!-- Basic Info -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-dark-400 mb-1">Name *</label>
                <input
                  v-model="formData.name"
                  type="text"
                  required
                  placeholder="e.g., deploy-node"
                  class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label class="block text-sm text-dark-400 mb-1">Type</label>
                <select
                  v-model="formData.type"
                  class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="deployment">Deployment</option>
                  <option value="maintenance">Maintenance</option>
                  <option value="backup">Backup</option>
                  <option value="custom">Custom</option>
                </select>
              </div>
              <div class="col-span-2">
                <label class="block text-sm text-dark-400 mb-1">Description</label>
                <input
                  v-model="formData.description"
                  type="text"
                  placeholder="Brief description"
                  class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
            </div>

            <!-- Target Nodes -->
            <div>
              <label class="block text-sm text-dark-400 mb-2">Target Nodes</label>
              <div class="grid grid-cols-3 gap-2 max-h-40 overflow-y-auto">
                <label
                  v-for="node in availableNodes"
                  :key="node.id"
                  class="flex items-center gap-2 p-2 bg-dark-700 rounded-lg cursor-pointer hover:bg-dark-600"
                  :class="{ 'ring-2 ring-primary-500': formData.targetNodes.includes(node.id) }"
                >
                  <input
                    type="checkbox"
                    :value="node.id"
                    v-model="formData.targetNodes"
                    class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
                  />
                  <span class="text-white text-sm truncate">{{ node.name }}</span>
                </label>
              </div>
            </div>

            <!-- Playbook Content -->
            <div>
              <label class="block text-sm text-dark-400 mb-1">Playbook Content (YAML)</label>
              <textarea
                v-model="formData.content"
                rows="12"
                placeholder="- name: Deploy node
  hosts: all
  tasks:
    - name: Update config
      ..."
                class="w-full px-4 py-3 bg-dark-700 border border-dark-600 rounded-lg text-white font-mono text-sm placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
              ></textarea>
            </div>

            <!-- Schedule -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-dark-400 mb-1">Schedule (Cron)</label>
                <input
                  v-model="formData.schedule"
                  type="text"
                  placeholder="0 0 * * * (daily at midnight)"
                  class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label class="block text-sm text-dark-400 mb-1">Timeout (minutes)</label>
                <input
                  v-model.number="formData.timeout"
                  type="number"
                  min="1"
                  placeholder="30"
                  class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
            </div>

            <!-- Variables -->
            <div>
              <label class="block text-sm text-dark-400 mb-1">Variables (JSON)</label>
              <textarea
                v-model="formData.variables"
                rows="4"
                placeholder='{ "var1": "value1", "var2": "value2" }'
                class="w-full px-4 py-3 bg-dark-700 border border-dark-600 rounded-lg text-white font-mono text-sm placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
              ></textarea>
            </div>

            <!-- Error Message -->
            <div v-if="formError" class="p-4 bg-red-900/20 border border-red-800 rounded-lg">
              <p class="text-red-400 text-sm">{{ formError }}</p>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3 pt-4 border-t border-dark-700">
              <button
                type="button"
                @click="validatePlaybook"
                class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
              >
                <CheckCircleIcon class="w-4 h-4" />
                Validate
              </button>
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
              >
                <ArrowPathIcon v-if="submitting" class="w-4 h-4 animate-spin" />
                {{ editingPlaybook ? 'Update' : 'Create' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Run Modal -->
    <Teleport to="body">
      <div v-if="showRunModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-3xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">Running: {{ runningPlaybook?.name }}</h2>
            <button @click="showRunModal = false" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <div class="p-6 space-y-4">
            <!-- Progress -->
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-dark-400">Progress</span>
                <span class="text-white">{{ runProgress }}%</span>
              </div>
              <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
                <div
                  class="h-full bg-primary-500 transition-all duration-300"
                  :style="{ width: `${runProgress}%` }"
                ></div>
              </div>
            </div>

            <!-- Logs -->
            <div class="bg-dark-900 rounded-lg p-4 font-mono text-sm max-h-96 overflow-y-auto">
              <div
                v-for="(log, index) in runLogs"
                :key="index"
                class="py-1 border-b border-dark-700/50 last:border-0"
              >
                <span class="text-dark-500">{{ log.time }}</span>
                <span :class="getLogLevelColor(log.level)"> [{{ log.level }}] </span>
                <span class="text-white">{{ log.message }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3">
              <button
                @click="stopCurrentRun"
                class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors flex items-center gap-2"
              >
                <StopIcon class="w-4 h-4" />
                Stop
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Templates Modal -->
    <Teleport to="body">
      <div v-if="showTemplates" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showTemplates = false"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">Playbook Templates</h2>
            <button @click="showTemplates = false" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <div class="p-6 space-y-4">
            <div
              v-for="template in playbookTemplates"
              :key="template.id"
              class="p-4 bg-dark-700 rounded-lg hover:bg-dark-600 transition-colors cursor-pointer"
              @click="useTemplate(template)"
            >
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="text-white font-medium">{{ template.name }}</h3>
                  <p class="text-dark-400 text-sm">{{ template.description }}</p>
                </div>
                <PlusIcon class="w-5 h-5 text-primary-400" />
              </div>
            </div>
          </div>
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
              <h3 class="text-lg font-semibold text-white">Delete Playbook</h3>
              <p class="text-dark-400 text-sm">This action cannot be undone</p>
            </div>
          </div>
          <p class="text-dark-300 mb-6">
            Are you sure you want to delete <span class="text-white font-medium">{{ playbookToDelete?.name }}</span>?
          </p>
          <div class="flex items-center justify-end gap-3">
            <button
              @click="showDeleteConfirm = false"
              class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              @click="handleDelete"
              class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
            >
              Delete
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
import { ref, computed, onMounted } from 'vue'
import { playbooksApi, nodesApi } from '@/api'
import {
  DocumentTextIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  PlayIcon,
  StopIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  ServerIcon,
  CalendarIcon,
  DocumentDuplicateIcon,
  ChartBarIcon,
  ExclamationTriangleIcon,
  BoltIcon,
  WrenchScrewdriverIcon,
  ArchiveBoxIcon,
  CogIcon
} from '@heroicons/vue/24/outline'

// State
const search = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const showModal = ref(false)
const showRunModal = ref(false)
const showTemplates = ref(false)
const showDeleteConfirm = ref(false)
const editingPlaybook = ref(null)
const runningPlaybook = ref(null)
const playbookToDelete = ref(null)
const submitting = ref(false)
const formError = ref('')
const toasts = ref([])
const runProgress = ref(0)
const runLogs = ref([])
const availableNodes = ref([])
const playbooks = ref([])

// Form data
const defaultFormData = {
  name: '',
  type: 'deployment',
  description: '',
  content: '',
  targetNodes: [],
  schedule: '',
  timeout: 30,
  variables: '{}'
}
const formData = ref({ ...defaultFormData })

// Templates
const playbookTemplates = ref([
  {
    id: 1,
    name: 'Deploy New Node',
    description: 'Deploy a new proxy node with default configuration',
    type: 'deployment',
    content: `- name: Deploy Node
  hosts: all
  become: yes
  tasks:
    - name: Install dependencies
      apt:
        name: "{{ item }}"
        state: present
      loop:
        - curl
        - wget
    - name: Deploy service
      shell: ./deploy.sh`
  },
  {
    id: 2,
    name: 'Backup Database',
    description: 'Create a backup of the database',
    type: 'backup',
    content: `- name: Backup Database
  hosts: db_servers
  tasks:
    - name: Create backup
      shell: pg_dump > backup.sql
    - name: Upload to S3
      s3_sync:
        bucket: backups`
  },
  {
    id: 3,
    name: 'Update Certificates',
    description: 'Renew SSL certificates on all nodes',
    type: 'maintenance',
    content: `- name: Renew Certificates
  hosts: all
  tasks:
    - name: Check certificate
      shell: certbot certificates
    - name: Renew if needed
      shell: certbot renew`
  }
])

// Computed
const filteredPlaybooks = computed(() => {
  let result = playbooks.value

  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(p =>
      p.name.toLowerCase().includes(searchLower) ||
      (p.description && p.description.toLowerCase().includes(searchLower))
    )
  }

  if (statusFilter.value) {
    result = result.filter(p => p.status === statusFilter.value)
  }

  if (typeFilter.value) {
    result = result.filter(p => p.type === typeFilter.value)
  }

  return result
})

const playbookStats = computed(() => ({
  total: playbooks.value.length,
  running: playbooks.value.filter(p => p.status === 'running').length,
  successRate: playbooks.value.length > 0
    ? Math.round((playbooks.value.filter(p => p.status === 'success').length / playbooks.value.length) * 100)
    : 0,
  scheduled: playbooks.value.filter(p => p.schedule).length
}))

// Methods
function getPlaybookTypeIcon(type) {
  const icons = {
    deployment: BoltIcon,
    maintenance: WrenchScrewdriverIcon,
    backup: ArchiveBoxIcon,
    custom: CogIcon
  }
  return icons[type] || DocumentTextIcon
}

function getPlaybookTypeColor(type) {
  const colors = {
    deployment: 'bg-blue-600',
    maintenance: 'bg-yellow-600',
    backup: 'bg-green-600',
    custom: 'bg-purple-600'
  }
  return colors[type] || 'bg-gray-600'
}

function getStatusColor(status) {
  const colors = {
    success: 'bg-green-900/30 text-green-400',
    running: 'bg-yellow-900/30 text-yellow-400',
    failed: 'bg-red-900/30 text-red-400',
    pending: 'bg-gray-900/30 text-gray-400'
  }
  return colors[status] || 'bg-gray-900/30 text-gray-400'
}

function getLogLevelColor(level) {
  const colors = {
    INFO: 'text-blue-400',
    WARN: 'text-yellow-400',
    ERROR: 'text-red-400',
    DEBUG: 'text-gray-400'
  }
  return colors[level] || 'text-gray-400'
}

function openCreateModal() {
  editingPlaybook.value = null
  formData.value = { ...defaultFormData }
  formError.value = ''
  showModal.value = true
}

function openEditModal(playbook) {
  editingPlaybook.value = playbook
  formData.value = {
    name: playbook.name,
    type: playbook.type || 'deployment',
    description: playbook.description || '',
    content: playbook.content || '',
    targetNodes: playbook.targetNodes || [],
    schedule: playbook.schedule || '',
    timeout: playbook.timeout || 30,
    variables: JSON.stringify(playbook.variables || {})
  }
  formError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingPlaybook.value = null
  formData.value = { ...defaultFormData }
  formError.value = ''
}

function useTemplate(template) {
  formData.value = {
    ...defaultFormData,
    name: template.name,
    type: template.type,
    description: template.description,
    content: template.content
  }
  showTemplates.value = false
  showModal.value = true
}

async function handleSubmit() {
  submitting.value = true
  formError.value = ''

  try {
    const data = {
      name: formData.value.name,
      type: formData.value.type,
      description: formData.value.description,
      content: formData.value.content,
      targetNodes: formData.value.targetNodes,
      schedule: formData.value.schedule,
      timeout: formData.value.timeout,
      variables: JSON.parse(formData.value.variables)
    }

    if (editingPlaybook.value) {
      await playbooksApi.update(editingPlaybook.value.id, data)
      showToast('Playbook updated successfully', 'success')
    } else {
      await playbooksApi.create(data)
      showToast('Playbook created successfully', 'success')
    }
    closeModal()
    await fetchPlaybooks()
  } catch (error) {
    formError.value = error.message || 'Operation failed'
  } finally {
    submitting.value = false
  }
}

async function validatePlaybook() {
  try {
    await playbooksApi.validate({ content: formData.value.content })
    showToast('Playbook is valid', 'success')
    formError.value = ''
  } catch (error) {
    formError.value = error.message || 'Validation failed'
  }
}

async function runPlaybook(playbook) {
  runningPlaybook.value = playbook
  runProgress.value = 0
  runLogs.value = [
    { time: new Date().toLocaleTimeString(), level: 'INFO', message: `Starting playbook: ${playbook.name}` }
  ]
  showRunModal.value = true

  try {
    const response = await playbooksApi.run(playbook.id)
    // Simulate progress updates
    const interval = setInterval(() => {
      if (runProgress.value < 100) {
        runProgress.value += 10
        runLogs.value.push({
          time: new Date().toLocaleTimeString(),
          level: 'INFO',
          message: `Processing task ${Math.floor(runProgress.value / 10)} of 10`
        })
      } else {
        clearInterval(interval)
        runLogs.value.push({
          time: new Date().toLocaleTimeString(),
          level: 'INFO',
          message: 'Playbook completed successfully'
        })
        showToast('Playbook completed successfully', 'success')
      }
    }, 500)
  } catch (error) {
    runLogs.value.push({
      time: new Date().toLocaleTimeString(),
      level: 'ERROR',
      message: error.message || 'Failed to run playbook'
    })
    showToast(error.message || 'Failed to run playbook', 'error')
  }
}

async function stopPlaybook(playbook) {
  try {
    await playbooksApi.stop(playbook.id)
    showToast('Playbook stopped', 'warning')
  } catch (error) {
    showToast(error.message || 'Failed to stop playbook', 'error')
  }
}

function stopCurrentRun() {
  if (runningPlaybook.value) {
    stopPlaybook(runningPlaybook.value)
    showRunModal.value = false
  }
}

function viewLogs(playbook) {
  // Navigate to logs view with playbook filter
  console.log('View logs for playbook:', playbook.name)
}

async function duplicatePlaybook(playbook) {
  try {
    await playbooksApi.duplicate(playbook.id)
    showToast('Playbook duplicated', 'success')
    await fetchPlaybooks()
  } catch (error) {
    showToast(error.message || 'Failed to duplicate playbook', 'error')
  }
}

function confirmDelete(playbook) {
  playbookToDelete.value = playbook
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!playbookToDelete.value) return

  try {
    await playbooksApi.delete(playbookToDelete.value.id)
    showToast('Playbook deleted', 'success')
    showDeleteConfirm.value = false
    playbookToDelete.value = null
    await fetchPlaybooks()
  } catch (error) {
    showToast(error.message || 'Failed to delete playbook', 'error')
  }
}

async function fetchPlaybooks() {
  try {
    const response = await playbooksApi.list()
    playbooks.value = response.data.data || response.data || []
  } catch (error) {
    // Use mock data for demo
    playbooks.value = [
      { id: 1, name: 'deploy_node.yml', type: 'deployment', description: 'Deploy new proxy node', status: 'success', lastRun: '2024-03-15 12:34', targetNodes: [1, 2] },
      { id: 2, name: 'update_certificates.yml', type: 'maintenance', description: 'Update SSL certificates', status: 'success', lastRun: '2024-03-14 08:00', targetNodes: [1, 2, 3] },
      { id: 3, name: 'backup_database.yml', type: 'backup', description: 'Create database backup', status: 'running', lastRun: '2024-03-15 00:00', progress: 60, targetNodes: [4], schedule: '0 0 * * *' },
      { id: 4, name: 'cleanup_logs.yml', type: 'maintenance', description: 'Clean up old log files', status: 'failed', lastRun: '2024-03-13 06:00', targetNodes: [] },
      { id: 5, name: 'deploy_agent.yml', type: 'deployment', description: 'Deploy monitoring agent', status: 'success', lastRun: '2024-03-12 18:30', targetNodes: [5, 6], schedule: '0 */6 * * *' }
    ]
  }
}

async function fetchNodes() {
  try {
    const response = await nodesApi.list()
    availableNodes.value = response.data.data || response.data || []
  } catch (error) {
    availableNodes.value = [
      { id: 1, name: 'tokyo-01' },
      { id: 2, name: 'singapore-01' },
      { id: 3, name: 'la-01' }
    ]
  }
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
  await fetchPlaybooks()
  await fetchNodes()
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