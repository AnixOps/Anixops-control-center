<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Users</h2>
        <p class="text-slate-400 text-sm mt-1">Manage platform users</p>
      </div>
      <button class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center gap-2">
        <PlusIcon class="w-5 h-5" />
        Add User
      </button>
    </div>

    <div class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden">
      <table class="w-full">
        <thead class="bg-slate-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Email</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Plan</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase">Traffic</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-700">
          <tr v-for="user in users" :key="user.id" class="hover:bg-slate-700/50">
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center">
                  <span class="text-sm font-medium text-white">{{ user.email.charAt(0).toUpperCase() }}</span>
                </div>
                <span class="text-white">{{ user.email }}</span>
              </div>
            </td>
            <td class="px-6 py-4">
              <span class="px-2 py-1 text-xs rounded bg-slate-700 text-slate-300">{{ user.plan }}</span>
            </td>
            <td class="px-6 py-4">
              <span
                class="px-2 py-1 text-xs rounded-full"
                :class="user.status === 'active' ? 'bg-green-900/30 text-green-400' : 'bg-red-900/30 text-red-400'"
              >
                {{ user.status }}
              </span>
            </td>
            <td class="px-6 py-4 text-slate-400">{{ user.traffic }}</td>
            <td class="px-6 py-4 text-right">
              <button
                @click="toggleBan(user)"
                class="px-3 py-1 text-sm rounded-lg transition-colors"
                :class="user.status === 'active' ? 'bg-red-900/30 text-red-400 hover:bg-red-900/50' : 'bg-green-900/30 text-green-400 hover:bg-green-900/50'"
              >
                {{ user.status === 'active' ? 'Ban' : 'Unban' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'

const PlusIcon = () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4v16m8-8H4' })
])

const users = ref([
  { id: 1, email: 'admin@example.com', plan: 'Enterprise', status: 'active', traffic: '12.4 GB' },
  { id: 2, email: 'user1@example.com', plan: 'Pro', status: 'active', traffic: '5.2 GB' },
  { id: 3, email: 'user2@example.com', plan: 'Basic', status: 'banned', traffic: '0 B' },
  { id: 4, email: 'user3@example.com', plan: 'Pro', status: 'active', traffic: '45.6 GB' },
])

function toggleBan(user) {
  user.status = user.status === 'active' ? 'banned' : 'active'
}
</script>