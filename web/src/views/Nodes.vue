<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-white">Nodes</h1>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors"
      >
        Add Node
      </button>
    </div>

    <!-- Filters -->
    <div class="flex gap-4">
      <input
        v-model="search"
        type="text"
        placeholder="Search nodes..."
        class="flex-1 px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-400 focus:outline-none focus:ring-2 focus:ring-primary-500"
      />
      <select
        v-model="statusFilter"
        class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
      >
        <option value="">All Status</option>
        <option value="online">Online</option>
        <option value="offline">Offline</option>
      </select>
    </div>

    <!-- Nodes Table -->
    <div class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden">
      <table class="w-full">
        <thead class="bg-dark-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">Name</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">Host</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">Users</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-dark-300 uppercase tracking-wider">Traffic</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-dark-300 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-dark-700">
          <tr
            v-for="node in filteredNodes"
            :key="node.id"
            class="hover:bg-dark-700/50"
          >
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-primary-600/20 flex items-center justify-center">
                  <ServerIcon class="w-4 h-4 text-primary-400" />
                </div>
                <span class="text-white font-medium">{{ node.name }}</span>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-dark-300">{{ node.host }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 py-1 text-xs rounded-full"
                :class="node.status === 'online' ? 'bg-green-900/30 text-green-400' : 'bg-red-900/30 text-red-400'"
              >
                {{ node.status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-dark-300">{{ node.users }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-dark-300">{{ node.traffic }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right">
              <div class="flex items-center justify-end gap-2">
                <button
                  @click="viewNode(node)"
                  class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                  title="View"
                >
                  <EyeIcon class="w-4 h-4 text-dark-400" />
                </button>
                <button
                  @click="editNode(node)"
                  class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                  title="Edit"
                >
                  <PencilIcon class="w-4 h-4 text-dark-400" />
                </button>
                <button
                  @click="deleteNode(node)"
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
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  ServerIcon,
  EyeIcon,
  PencilIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'

const search = ref('')
const statusFilter = ref('')
const showCreateModal = ref(false)

const nodes = ref([
  { id: 1, name: 'tokyo-01', host: '192.168.1.101', status: 'online', users: 156, traffic: '1.2TB' },
  { id: 2, name: 'singapore-01', host: '192.168.1.102', status: 'online', users: 89, traffic: '890GB' },
  { id: 3, name: 'la-01', host: '192.168.1.103', status: 'offline', users: 0, traffic: '0B' },
  { id: 4, name: 'frankfurt-01', host: '192.168.1.104', status: 'online', users: 45, traffic: '234GB' },
  { id: 5, name: 'london-01', host: '192.168.1.105', status: 'online', users: 67, traffic: '456GB' }
])

const filteredNodes = computed(() => {
  return nodes.value.filter(node => {
    const matchesSearch = node.name.toLowerCase().includes(search.value.toLowerCase()) ||
                          node.host.includes(search.value)
    const matchesStatus = !statusFilter.value || node.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

function viewNode(node) {
  console.log('View node:', node)
}

function editNode(node) {
  console.log('Edit node:', node)
}

function deleteNode(node) {
  console.log('Delete node:', node)
}
</script>