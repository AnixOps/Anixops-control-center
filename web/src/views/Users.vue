<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">{{ t('users.title') }}</h1>
        <p class="text-dark-400 text-sm mt-1">{{ t('users.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="handleExport"
          :disabled="exporting"
          class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <ArrowDownTrayIcon class="w-5 h-5" />
          {{ t('users.exportUsers') }}
        </button>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
        >
          <PlusIcon class="w-5 h-5" />
          {{ t('users.addUser') }}
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('users.total') }}</p>
            <p class="text-2xl font-bold text-white">{{ userStats.total }}</p>
          </div>
          <UsersIcon class="w-8 h-8 text-primary-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('users.active') }}</p>
            <p class="text-2xl font-bold text-green-400">{{ userStats.active }}</p>
          </div>
          <CheckCircleIcon class="w-8 h-8 text-green-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('users.banned') }}</p>
            <p class="text-2xl font-bold text-red-400">{{ userStats.banned }}</p>
          </div>
          <NoSymbolIcon class="w-8 h-8 text-red-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('users.admins') }}</p>
            <p class="text-2xl font-bold text-purple-400">{{ userStats.admins }}</p>
          </div>
          <ShieldCheckIcon class="w-8 h-8 text-purple-400" />
        </div>
      </div>
      <div class="bg-dark-800 rounded-xl p-4 border border-dark-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-dark-400 text-sm">{{ t('common.online') }}</p>
            <p class="text-2xl font-bold text-blue-400">{{ onlineCount }}</p>
          </div>
          <SignalIcon class="w-8 h-8 text-blue-400" />
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
        <option value="active">{{ t('common.active') }}</option>
        <option value="banned">{{ t('users.banned') }}</option>
        <option value="suspended">Suspended</option>
      </select>
      <select
        v-model="planFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        @change="handleFilterChange"
      >
        <option value="">{{ t('common.all') }} {{ t('users.plan') }}</option>
        <option value="free">Free</option>
        <option value="basic">Basic</option>
        <option value="pro">Pro</option>
        <option value="enterprise">Enterprise</option>
      </select>
      <select
        v-model="roleFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        @change="handleFilterChange"
      >
        <option value="">{{ t('common.all') }} {{ t('users.role') }}</option>
        <option value="user">{{ t('common.active') }}</option>
        <option value="admin">Admin</option>
        <option value="superadmin">Super Admin</option>
      </select>
      <button
        @click="refreshUsers"
        :disabled="loading"
        class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
      >
        <ArrowPathIcon :class="['w-5 h-5', { 'animate-spin': loading }]" />
        {{ t('common.refresh') }}
      </button>
    </div>

    <!-- Users Table -->
    <div class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-dark-700">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('navigation.users') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('users.plan') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('users.role') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('common.status') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('nodes.traffic') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('users.expiresAt') }}</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-dark-300 uppercase tracking-wider">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-dark-700">
            <tr v-if="loading">
              <td colspan="7" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center gap-3">
                  <ArrowPathIcon class="w-8 h-8 text-primary-400 animate-spin" />
                  <p class="text-dark-400">{{ t('common.loading') }}</p>
                </div>
              </td>
            </tr>
            <tr v-else-if="filteredUsers.length === 0">
              <td colspan="7" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center gap-3">
                  <UsersIcon class="w-12 h-12 text-dark-500" />
                  <p class="text-dark-400">{{ t('common.noResults') }}</p>
                </div>
              </td>
            </tr>
            <tr
              v-else
              v-for="user in filteredUsers"
              :key="user.id"
              class="hover:bg-dark-700/50 transition-colors"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 rounded-full flex items-center justify-center text-white font-medium"
                    :class="getAvatarColor(user.email)"
                  >
                    {{ getInitials(user.email) }}
                  </div>
                  <div>
                    <p class="text-white font-medium">{{ user.email }}</p>
                    <p class="text-dark-400 text-sm">{{ user.name || 'No name' }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  class="px-2 py-1 text-xs rounded"
                  :class="getPlanColor(user.plan)"
                >
                  {{ user.plan || 'Free' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <ShieldCheckIcon v-if="user.role === 'admin' || user.role === 'superadmin'" class="w-4 h-4 text-purple-400" />
                  <span class="text-dark-300 capitalize">{{ user.role || 'user' }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <div
                    class="w-2 h-2 rounded-full"
                    :class="getStatusDotColor(user.status)"
                  ></div>
                  <span
                    class="px-2 py-1 text-xs rounded-full"
                    :class="getStatusColor(user.status)"
                  >
                    {{ user.status }}
                  </span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div>
                  <p class="text-white text-sm">{{ formatBytes(user.usedTraffic || 0) }}</p>
                  <p class="text-dark-500 text-xs">/ {{ formatBytes(user.totalTraffic || 0) }}</p>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <p
                  class="text-sm"
                  :class="isExpired(user.expiresAt) ? 'text-red-400' : 'text-dark-300'"
                >
                  {{ formatDate(user.expiresAt) }}
                </p>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click="openDetailModal(user)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="View Details"
                  >
                    <EyeIcon class="w-4 h-4 text-dark-400" />
                  </button>
                  <button
                    @click="openEditModal(user)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <PencilIcon class="w-4 h-4 text-dark-400" />
                  </button>
                  <button
                    v-if="user.status === 'active'"
                    @click="handleBan(user)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Ban"
                  >
                    <NoSymbolIcon class="w-4 h-4 text-red-400" />
                  </button>
                  <button
                    v-else
                    @click="handleUnban(user)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Unban"
                  >
                    <CheckIcon class="w-4 h-4 text-green-400" />
                  </button>
                  <button
                    @click="handleResetPassword(user)"
                    class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                    title="Reset Password"
                  >
                    <KeyIcon class="w-4 h-4 text-yellow-400" />
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
          Showing {{ (currentPage - 1) * pageSize + 1 }} to {{ Math.min(currentPage * pageSize, total) }} of {{ total }} users
        </p>
        <div class="flex items-center gap-2">
          <button
            @click="changePage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            Previous
          </button>
          <div class="flex items-center gap-1">
            <button
              v-for="page in displayedPages"
              :key="page"
              @click="changePage(page)"
              class="w-8 h-8 rounded-lg transition-colors"
              :class="page === currentPage ? 'bg-primary-600 text-white' : 'bg-dark-700 text-dark-300 hover:bg-dark-600'"
            >
              {{ page }}
            </button>
          </div>
          <button
            @click="changePage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            Next
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
              {{ editingUser ? 'Edit User' : 'Add New User' }}
            </h2>
            <button @click="closeModal" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
            <!-- Account Info -->
            <div class="space-y-4">
              <h3 class="text-lg font-medium text-white">Account Information</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="col-span-2">
                  <label class="block text-sm text-dark-400 mb-1">Email *</label>
                  <input
                    v-model="formData.email"
                    type="email"
                    required
                    placeholder="user@example.com"
                    :disabled="!!editingUser"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50"
                  />
                </div>
                <div v-if="!editingUser">
                  <label class="block text-sm text-dark-400 mb-1">Password *</label>
                  <input
                    v-model="formData.password"
                    type="password"
                    :required="!editingUser"
                    placeholder="Min 8 characters"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
                <div v-if="!editingUser">
                  <label class="block text-sm text-dark-400 mb-1">Confirm Password *</label>
                  <input
                    v-model="formData.confirmPassword"
                    type="password"
                    :required="!editingUser"
                    placeholder="Confirm password"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>

            <!-- Subscription -->
            <div class="space-y-4">
              <h3 class="text-lg font-medium text-white">Subscription</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Plan</label>
                  <select
                    v-model="formData.plan"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="free">Free</option>
                    <option value="basic">Basic</option>
                    <option value="pro">Pro</option>
                    <option value="enterprise">Enterprise</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Role</label>
                  <select
                    v-model="formData.role"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="user">User</option>
                    <option value="admin">Admin</option>
                    <option value="superadmin">Super Admin</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Traffic Limit (GB)</label>
                  <input
                    v-model.number="formData.totalTraffic"
                    type="number"
                    min="0"
                    placeholder="0 = unlimited"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label class="block text-sm text-dark-400 mb-1">Expires At</label>
                  <input
                    v-model="formData.expiresAt"
                    type="date"
                    class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>

            <!-- Node Access -->
            <div class="space-y-4">
              <h3 class="text-lg font-medium text-white">Node Access</h3>
              <div class="grid grid-cols-3 gap-2">
                <label
                  v-for="node in availableNodes"
                  :key="node.id"
                  class="flex items-center gap-2 p-3 bg-dark-700 rounded-lg cursor-pointer hover:bg-dark-600"
                  :class="{ 'ring-2 ring-primary-500': formData.nodeAccess.includes(node.id) }"
                >
                  <input
                    type="checkbox"
                    :value="node.id"
                    v-model="formData.nodeAccess"
                    class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
                  />
                  <span class="text-white text-sm">{{ node.name }}</span>
                </label>
              </div>
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
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
              >
                <ArrowPathIcon v-if="submitting" class="w-4 h-4 animate-spin" />
                {{ editingUser ? 'Update' : 'Create' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- User Detail Modal -->
    <Teleport to="body">
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showDetailModal = false"></div>
        <div class="relative bg-dark-800 rounded-2xl border border-dark-700 w-full max-w-3xl max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-dark-800 px-6 py-4 border-b border-dark-700 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">User Details</h2>
            <button @click="showDetailModal = false" class="p-2 hover:bg-dark-700 rounded-lg">
              <XMarkIcon class="w-5 h-5 text-dark-400" />
            </button>
          </div>

          <div v-if="selectedUser" class="p-6 space-y-6">
            <!-- User Header -->
            <div class="flex items-center gap-4">
              <div
                class="w-16 h-16 rounded-full flex items-center justify-center text-white text-xl font-medium"
                :class="getAvatarColor(selectedUser.email)"
              >
                {{ getInitials(selectedUser.email) }}
              </div>
              <div>
                <h3 class="text-xl font-semibold text-white">{{ selectedUser.email }}</h3>
                <p class="text-dark-400">{{ selectedUser.name || 'No display name' }}</p>
              </div>
              <span
                class="ml-auto px-3 py-1 rounded-full text-sm"
                :class="getStatusColor(selectedUser.status)"
              >
                {{ selectedUser.status }}
              </span>
            </div>

            <!-- Stats Grid -->
            <div class="grid grid-cols-4 gap-4">
              <div class="bg-dark-700 rounded-lg p-4">
                <p class="text-dark-400 text-sm">Plan</p>
                <p class="text-white font-medium capitalize">{{ selectedUser.plan || 'Free' }}</p>
              </div>
              <div class="bg-dark-700 rounded-lg p-4">
                <p class="text-dark-400 text-sm">Role</p>
                <p class="text-white font-medium capitalize">{{ selectedUser.role || 'User' }}</p>
              </div>
              <div class="bg-dark-700 rounded-lg p-4">
                <p class="text-dark-400 text-sm">Traffic Used</p>
                <p class="text-white font-medium">{{ formatBytes(selectedUser.usedTraffic || 0) }}</p>
              </div>
              <div class="bg-dark-700 rounded-lg p-4">
                <p class="text-dark-400 text-sm">Expires</p>
                <p class="text-white font-medium">{{ formatDate(selectedUser.expiresAt) }}</p>
              </div>
            </div>

            <!-- Activity Log -->
            <div>
              <h4 class="text-lg font-medium text-white mb-3">Recent Activity</h4>
              <div class="space-y-2">
                <div
                  v-for="activity in userActivities"
                  :key="activity.id"
                  class="flex items-center gap-3 p-3 bg-dark-700 rounded-lg"
                >
                  <component :is="getActivityIcon(activity.type)" class="w-5 h-5 text-dark-400" />
                  <div class="flex-1">
                    <p class="text-white text-sm">{{ activity.description }}</p>
                    <p class="text-dark-500 text-xs">{{ activity.time }}</p>
                  </div>
                </div>
              </div>
            </div>
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
import { useI18n } from 'vue-i18n'
import { useUsersStore } from '@/stores/users'
import { useNodesStore } from '@/stores/nodes'
import {
  UsersIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon,
  PencilIcon,
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
  EyeIcon,
  NoSymbolIcon,
  CheckIcon,
  KeyIcon,
  ShieldCheckIcon,
  ArrowDownTrayIcon,
  SignalIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const usersStore = useUsersStore()
const nodesStore = useNodesStore()

// State
const search = ref('')
const statusFilter = ref('')
const planFilter = ref('')
const roleFilter = ref('')
const showModal = ref(false)
const showDetailModal = ref(false)
const editingUser = ref(null)
const selectedUser = ref(null)
const submitting = ref(false)
const exporting = ref(false)
const formError = ref('')
const toasts = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const searchTimeout = ref(null)
const onlineCount = ref(0)

// Form data
const defaultFormData = {
  email: '',
  password: '',
  confirmPassword: '',
  plan: 'free',
  role: 'user',
  totalTraffic: 0,
  expiresAt: '',
  nodeAccess: []
}
const formData = ref({ ...defaultFormData })

// Mock data for demo
const userActivities = ref([
  { id: 1, type: 'login', description: 'Logged in from 192.168.1.1', time: '2 minutes ago' },
  { id: 2, type: 'traffic', description: 'Used 1.2GB traffic', time: '1 hour ago' },
  { id: 3, type: 'node', description: 'Connected to tokyo-01', time: '3 hours ago' }
])

// Computed
const loading = computed(() => usersStore.loading)
const users = computed(() => usersStore.users)
const total = computed(() => usersStore.total)
const userStats = computed(() => usersStore.userStats)
const availableNodes = computed(() => nodesStore.nodes)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const filteredUsers = computed(() => {
  let result = users.value

  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(user =>
      user.email.toLowerCase().includes(searchLower) ||
      (user.name && user.name.toLowerCase().includes(searchLower))
    )
  }

  if (statusFilter.value) {
    result = result.filter(user => user.status === statusFilter.value)
  }

  if (planFilter.value) {
    result = result.filter(user => user.plan === planFilter.value)
  }

  if (roleFilter.value) {
    result = result.filter(user => user.role === roleFilter.value)
  }

  return result
})

const displayedPages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Methods
function getInitials(email) {
  return email.slice(0, 2).toUpperCase()
}

function getAvatarColor(email) {
  const colors = [
    'bg-red-500', 'bg-orange-500', 'bg-amber-500', 'bg-yellow-500',
    'bg-lime-500', 'bg-green-500', 'bg-emerald-500', 'bg-teal-500',
    'bg-cyan-500', 'bg-sky-500', 'bg-blue-500', 'bg-indigo-500',
    'bg-violet-500', 'bg-purple-500', 'bg-fuchsia-500', 'bg-pink-500'
  ]
  const hash = email.split('').reduce((a, b) => {
    a = ((a << 5) - a) + b.charCodeAt(0)
    return a & a
  }, 0)
  return colors[Math.abs(hash) % colors.length]
}

function getStatusColor(status) {
  const colors = {
    active: 'bg-green-900/30 text-green-400',
    banned: 'bg-red-900/30 text-red-400',
    suspended: 'bg-yellow-900/30 text-yellow-400',
    pending: 'bg-blue-900/30 text-blue-400'
  }
  return colors[status] || 'bg-gray-900/30 text-gray-400'
}

function getStatusDotColor(status) {
  const colors = {
    active: 'bg-green-400',
    banned: 'bg-red-400',
    suspended: 'bg-yellow-400',
    pending: 'bg-blue-400'
  }
  return colors[status] || 'bg-gray-400'
}

function getPlanColor(plan) {
  const colors = {
    free: 'bg-gray-700 text-gray-300',
    basic: 'bg-blue-900/30 text-blue-400',
    pro: 'bg-purple-900/30 text-purple-400',
    enterprise: 'bg-amber-900/30 text-amber-400'
  }
  return colors[plan] || 'bg-gray-700 text-gray-300'
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(dateStr) {
  if (!dateStr) return 'Never'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return 'Invalid'
  return date.toLocaleDateString()
}

function isExpired(dateStr) {
  if (!dateStr) return false
  return new Date(dateStr) < new Date()
}

function getActivityIcon(type) {
  const icons = {
    login: UsersIcon,
    traffic: SignalIcon,
    node: UsersIcon
  }
  return icons[type] || UsersIcon
}

function debouncedSearch() {
  clearTimeout(searchTimeout.value)
  searchTimeout.value = setTimeout(() => {
    handleFilterChange()
  }, 300)
}

function handleFilterChange() {
  currentPage.value = 1
  usersStore.setFilters({
    search: search.value,
    status: statusFilter.value,
    plan: planFilter.value,
    role: roleFilter.value,
    page: currentPage.value,
    limit: pageSize.value
  })
}

async function refreshUsers() {
  try {
    await usersStore.fetchUsers()
    showToast('Users refreshed', 'success')
  } catch (error) {
    showToast(error.message || 'Failed to refresh users', 'error')
  }
}

function openCreateModal() {
  editingUser.value = null
  formData.value = { ...defaultFormData }
  formError.value = ''
  showModal.value = true
}

function openEditModal(user) {
  editingUser.value = user
  formData.value = {
    email: user.email,
    password: '',
    confirmPassword: '',
    plan: user.plan || 'free',
    role: user.role || 'user',
    totalTraffic: user.totalTraffic || 0,
    expiresAt: user.expiresAt ? user.expiresAt.split('T')[0] : '',
    nodeAccess: user.nodeAccess || []
  }
  formError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingUser.value = null
  formData.value = { ...defaultFormData }
  formError.value = ''
}

async function handleSubmit() {
  if (!editingUser.value && formData.value.password !== formData.value.confirmPassword) {
    formError.value = 'Passwords do not match'
    return
  }

  submitting.value = true
  formError.value = ''

  try {
    const data = {
      email: formData.value.email,
      plan: formData.value.plan,
      role: formData.value.role,
      totalTraffic: formData.value.totalTraffic * 1024 * 1024 * 1024, // GB to bytes
      expiresAt: formData.value.expiresAt,
      nodeAccess: formData.value.nodeAccess
    }

    if (!editingUser.value) {
      data.password = formData.value.password
    }

    if (editingUser.value) {
      await usersStore.updateUser(editingUser.value.id, data)
      showToast('User updated successfully', 'success')
    } else {
      await usersStore.createUser(data)
      showToast('User created successfully', 'success')
    }
    closeModal()
  } catch (error) {
    formError.value = error.message || 'Operation failed'
  } finally {
    submitting.value = false
  }
}

function openDetailModal(user) {
  selectedUser.value = user
  showDetailModal.value = true
}

async function handleBan(user) {
  try {
    await usersStore.banUser(user.id)
    showToast(`${user.email} has been banned`, 'success')
  } catch (error) {
    showToast(error.message || 'Failed to ban user', 'error')
  }
}

async function handleUnban(user) {
  try {
    await usersStore.unbanUser(user.id)
    showToast(`${user.email} has been unbanned`, 'success')
  } catch (error) {
    showToast(error.message || 'Failed to unban user', 'error')
  }
}

async function handleResetPassword(user) {
  try {
    await usersStore.resetPassword(user.id)
    showToast(`Password reset email sent to ${user.email}`, 'success')
  } catch (error) {
    showToast(error.message || 'Failed to reset password', 'error')
  }
}

async function handleExport() {
  exporting.value = true
  try {
    await usersStore.exportUsers()
    showToast('Users exported successfully', 'success')
  } catch (error) {
    showToast(error.message || 'Failed to export users', 'error')
  } finally {
    exporting.value = false
  }
}

function changePage(page) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  usersStore.setFilters({ page })
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
  await usersStore.fetchUsers()
  await nodesStore.fetchNodes()
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