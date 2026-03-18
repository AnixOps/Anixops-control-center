<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Settings</h2>
        <p class="text-slate-400 text-sm mt-1">Configure your AnixOps Control Center</p>
      </div>
      <button class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors">
        Save Changes
      </button>
    </div>

    <!-- Server Configuration -->
    <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
      <h3 class="text-lg font-semibold text-white mb-4">Server Configuration</h3>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm text-slate-400 mb-1">Host</label>
          <input
            v-model="settings.host"
            type="text"
            class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label class="block text-sm text-slate-400 mb-1">Port</label>
          <input
            v-model.number="settings.port"
            type="number"
            class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>
    </div>

    <!-- Security -->
    <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
      <h3 class="text-lg font-semibold text-white mb-4">Security</h3>
      <div class="space-y-4">
        <label class="flex items-center justify-between p-3 bg-slate-700/50 rounded-lg cursor-pointer">
          <div>
            <p class="text-white font-medium">Two-Factor Authentication</p>
            <p class="text-slate-400 text-sm">Require 2FA for admin accounts</p>
          </div>
          <input
            v-model="settings.twoFactor"
            type="checkbox"
            class="w-5 h-5 rounded border-slate-500 bg-slate-600 text-blue-500 focus:ring-blue-500"
          />
        </label>
        <label class="flex items-center justify-between p-3 bg-slate-700/50 rounded-lg cursor-pointer">
          <div>
            <p class="text-white font-medium">Auto Backup</p>
            <p class="text-slate-400 text-sm">Create daily backups automatically</p>
          </div>
          <input
            v-model="settings.autoBackup"
            type="checkbox"
            class="w-5 h-5 rounded border-slate-500 bg-slate-600 text-blue-500 focus:ring-blue-500"
          />
        </label>
      </div>
    </div>

    <!-- Plugins -->
    <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
      <h3 class="text-lg font-semibold text-white mb-4">Plugins</h3>
      <div class="space-y-3">
        <div
          v-for="plugin in plugins"
          :key="plugin.name"
          class="flex items-center justify-between p-3 bg-slate-700/50 rounded-lg"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-2 h-2 rounded-full"
              :class="plugin.status === 'running' ? 'bg-green-400' : 'bg-red-400'"
            ></div>
            <span class="text-white">{{ plugin.name }}</span>
          </div>
          <span class="text-slate-400 text-sm">{{ plugin.version }}</span>
        </div>
      </div>
    </div>

    <!-- Danger Zone -->
    <div class="bg-red-900/20 rounded-xl border border-red-800 p-6">
      <h3 class="text-lg font-semibold text-red-400 mb-4">Danger Zone</h3>
      <div class="flex items-center justify-between">
        <div>
          <p class="text-white font-medium">Restart Server</p>
          <p class="text-slate-400 text-sm">This will restart the AnixOps server</p>
        </div>
        <button class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors">
          Restart
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const settings = ref({
  host: '0.0.0.0',
  port: 8080,
  twoFactor: false,
  autoBackup: true
})

const plugins = ref([
  { name: 'ansible', status: 'running', version: '1.0.0' },
  { name: 'v2board', status: 'running', version: '1.0.0' },
  { name: 'v2bx', status: 'running', version: '1.0.0' },
  { name: 'agent', status: 'running', version: '1.0.0' }
])
</script>