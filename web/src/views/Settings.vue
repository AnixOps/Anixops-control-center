<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">Settings</h1>
        <p class="text-dark-400 text-sm mt-1">Configure your AnixOps Control Center</p>
      </div>
      <button
        @click="saveAllSettings"
        :disabled="saving"
        class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
      >
        <CheckIcon v-if="!saving" class="w-5 h-5" />
        <ArrowPathIcon v-else class="w-5 h-5 animate-spin" />
        {{ saving ? 'Saving...' : 'Save Changes' }}
      </button>
    </div>

    <!-- Settings Navigation -->
    <div class="flex gap-2 border-b border-dark-700 pb-2">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="[
          'px-4 py-2 rounded-lg transition-colors flex items-center gap-2',
          activeTab === tab.id ? 'bg-primary-600 text-white' : 'text-dark-400 hover:text-white hover:bg-dark-700'
        ]"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.name }}
      </button>
    </div>

    <!-- Server Configuration -->
    <div v-show="activeTab === 'server'" class="space-y-6">
      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <ServerIcon class="w-5 h-5 text-primary-400" />
          Server Configuration
        </h2>
        <div class="grid grid-cols-2 gap-6">
          <div>
            <label class="block text-sm text-dark-400 mb-1">Host</label>
            <input
              v-model="settings.server.host"
              type="text"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">Port</label>
            <input
              v-model.number="settings.server.port"
              type="number"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">SSL Certificate</label>
            <select
              v-model="settings.server.sslMode"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="none">None (HTTP)</option>
              <option value="auto">Auto (Let's Encrypt)</option>
              <option value="custom">Custom Certificate</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">Max Connections</label>
            <input
              v-model.number="settings.server.maxConnections"
              type="number"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div class="col-span-2">
            <label class="block text-sm text-dark-400 mb-1">Trusted Proxies</label>
            <input
              v-model="settings.server.trustedProxies"
              type="text"
              placeholder="Comma-separated IPs or CIDR ranges"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
        </div>
      </div>

      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <DatabaseIcon class="w-5 h-5 text-primary-400" />
          Database Configuration
        </h2>
        <div class="grid grid-cols-2 gap-6">
          <div>
            <label class="block text-sm text-dark-400 mb-1">Database Type</label>
            <select
              v-model="settings.database.type"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="sqlite">SQLite</option>
              <option value="mysql">MySQL</option>
              <option value="postgres">PostgreSQL</option>
            </select>
          </div>
          <div v-if="settings.database.type !== 'sqlite'">
            <label class="block text-sm text-dark-400 mb-1">Host</label>
            <input
              v-model="settings.database.host"
              type="text"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div v-if="settings.database.type !== 'sqlite'">
            <label class="block text-sm text-dark-400 mb-1">Port</label>
            <input
              v-model.number="settings.database.port"
              type="number"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div v-if="settings.database.type !== 'sqlite'">
            <label class="block text-sm text-dark-400 mb-1">Database Name</label>
            <input
              v-model="settings.database.name"
              type="text"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div v-if="settings.database.type !== 'sqlite'">
            <label class="block text-sm text-dark-400 mb-1">Max Connections</label>
            <input
              v-model.number="settings.database.maxConnections"
              type="number"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Security Settings -->
    <div v-show="activeTab === 'security'" class="space-y-6">
      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <ShieldCheckIcon class="w-5 h-5 text-primary-400" />
          Authentication
        </h2>
        <div class="space-y-4">
          <div class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg">
            <div>
              <p class="text-white font-medium">Two-Factor Authentication</p>
              <p class="text-dark-400 text-sm">Require 2FA for all admin accounts</p>
            </div>
            <button
              @click="settings.security.twoFactorEnabled = !settings.security.twoFactorEnabled"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                settings.security.twoFactorEnabled ? 'bg-primary-600' : 'bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  settings.security.twoFactorEnabled ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
          <div class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg">
            <div>
              <p class="text-white font-medium">Session Timeout</p>
              <p class="text-dark-400 text-sm">Automatically logout after inactivity</p>
            </div>
            <select
              v-model="settings.security.sessionTimeout"
              class="px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="15">15 minutes</option>
              <option value="30">30 minutes</option>
              <option value="60">1 hour</option>
              <option value="120">2 hours</option>
              <option value="0">Never</option>
            </select>
          </div>
          <div class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg">
            <div>
              <p class="text-white font-medium">IP Whitelist</p>
              <p class="text-dark-400 text-sm">Only allow login from specific IPs</p>
            </div>
            <button
              @click="settings.security.ipWhitelistEnabled = !settings.security.ipWhitelistEnabled"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                settings.security.ipWhitelistEnabled ? 'bg-primary-600' : 'bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  settings.security.ipWhitelistEnabled ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
          <div v-if="settings.security.ipWhitelistEnabled">
            <label class="block text-sm text-dark-400 mb-1">Allowed IPs</label>
            <textarea
              v-model="settings.security.allowedIPs"
              rows="3"
              placeholder="One IP per line (e.g., 192.168.1.0/24)"
              class="w-full px-4 py-3 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
            ></textarea>
          </div>
        </div>
      </div>

      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <KeyIcon class="w-5 h-5 text-primary-400" />
          Password Policy
        </h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-dark-400 mb-1">Minimum Length</label>
            <input
              v-model.number="settings.security.minPasswordLength"
              type="number"
              min="6"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">Password Expiry (days)</label>
            <input
              v-model.number="settings.security.passwordExpiry"
              type="number"
              min="0"
              placeholder="0 = never expires"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div class="col-span-2 space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="settings.security.requireUppercase"
                type="checkbox"
                class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-dark-300">Require uppercase letter</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="settings.security.requireLowercase"
                type="checkbox"
                class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-dark-300">Require lowercase letter</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="settings.security.requireNumber"
                type="checkbox"
                class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-dark-300">Require number</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="settings.security.requireSpecialChar"
                type="checkbox"
                class="rounded border-dark-500 bg-dark-600 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-dark-300">Require special character</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Notification Settings -->
    <div v-show="activeTab === 'notifications'" class="space-y-6">
      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <BellIcon class="w-5 h-5 text-primary-400" />
          Email Notifications
        </h2>
        <div class="grid grid-cols-2 gap-6">
          <div>
            <label class="block text-sm text-dark-400 mb-1">SMTP Host</label>
            <input
              v-model="settings.notifications.smtpHost"
              type="text"
              placeholder="smtp.example.com"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">SMTP Port</label>
            <input
              v-model.number="settings.notifications.smtpPort"
              type="number"
              placeholder="587"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">Username</label>
            <input
              v-model="settings.notifications.smtpUser"
              type="text"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">Password</label>
            <input
              v-model="settings.notifications.smtpPassword"
              type="password"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm text-dark-400 mb-1">From Email</label>
            <input
              v-model="settings.notifications.fromEmail"
              type="email"
              placeholder="noreply@example.com"
              class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div class="flex items-end">
            <button
              @click="testEmail"
              class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors"
            >
              Send Test Email
            </button>
          </div>
        </div>
      </div>

      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4">Alert Rules</h2>
        <div class="space-y-3">
          <div
            v-for="alert in alertRules"
            :key="alert.id"
            class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg"
          >
            <div class="flex items-center gap-3">
              <component :is="alert.icon" class="w-5 h-5" :class="alert.color" />
              <div>
                <p class="text-white">{{ alert.name }}</p>
                <p class="text-dark-400 text-sm">{{ alert.description }}</p>
              </div>
            </div>
            <button
              @click="alert.enabled = !alert.enabled"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                alert.enabled ? 'bg-primary-600' : 'bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  alert.enabled ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Plugins Settings -->
    <div v-show="activeTab === 'plugins'" class="space-y-6">
      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <PuzzlePieceIcon class="w-5 h-5 text-primary-400" />
          Installed Plugins
        </h2>
        <div class="space-y-3">
          <div
            v-for="plugin in plugins"
            :key="plugin.name"
            class="flex items-center justify-between p-4 bg-dark-700/50 rounded-lg"
          >
            <div class="flex items-center gap-4">
              <div
                class="w-10 h-10 rounded-lg flex items-center justify-center"
                :class="plugin.enabled ? 'bg-primary-600' : 'bg-dark-600'"
              >
                <component :is="plugin.icon" class="w-5 h-5 text-white" />
              </div>
              <div>
                <p class="text-white font-medium">{{ plugin.name }}</p>
                <p class="text-dark-400 text-sm">{{ plugin.description }}</p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <span
                class="px-2 py-1 text-xs rounded"
                :class="plugin.enabled ? 'bg-green-900/30 text-green-400' : 'bg-gray-900/30 text-gray-400'"
              >
                {{ plugin.enabled ? 'Enabled' : 'Disabled' }}
              </span>
              <span class="text-dark-400 text-sm">v{{ plugin.version }}</span>
              <div class="flex items-center gap-1">
                <button
                  @click="togglePlugin(plugin)"
                  class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                  :title="plugin.enabled ? 'Disable' : 'Enable'"
                >
                  <component
                    :is="plugin.enabled ? PauseIcon : PlayIcon"
                    class="w-4 h-4"
                    :class="plugin.enabled ? 'text-yellow-400' : 'text-green-400'"
                  />
                </button>
                <button
                  @click="restartPlugin(plugin)"
                  class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                  title="Restart"
                >
                  <ArrowPathIcon class="w-4 h-4 text-dark-400" />
                </button>
                <button
                  @click="openPluginSettings(plugin)"
                  class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
                  title="Settings"
                >
                  <CogIcon class="w-4 h-4 text-dark-400" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Backup & Restore -->
    <div v-show="activeTab === 'backup'" class="space-y-6">
      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <ArchiveBoxIcon class="w-5 h-5 text-primary-400" />
          Backup Settings
        </h2>
        <div class="space-y-4">
          <div class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg">
            <div>
              <p class="text-white font-medium">Automatic Backups</p>
              <p class="text-dark-400 text-sm">Create backups automatically</p>
            </div>
            <button
              @click="settings.backup.autoBackup = !settings.backup.autoBackup"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                settings.backup.autoBackup ? 'bg-primary-600' : 'bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  settings.backup.autoBackup ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
          <div v-if="settings.backup.autoBackup" class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm text-dark-400 mb-1">Frequency</label>
              <select
                v-model="settings.backup.frequency"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                <option value="hourly">Hourly</option>
                <option value="daily">Daily</option>
                <option value="weekly">Weekly</option>
                <option value="monthly">Monthly</option>
              </select>
            </div>
            <div>
              <label class="block text-sm text-dark-400 mb-1">Retention (days)</label>
              <input
                v-model.number="settings.backup.retention"
                type="number"
                min="1"
                class="w-full px-4 py-2 bg-dark-700 border border-dark-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
          </div>
          <div class="flex items-center gap-4 pt-4">
            <button
              @click="createBackup"
              class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors flex items-center gap-2"
            >
              <ArchiveBoxIcon class="w-5 h-5" />
              Create Backup
            </button>
            <button
              @click="restoreBackup"
              class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors flex items-center gap-2"
            >
              <ArrowUturnLeftIcon class="w-5 h-5" />
              Restore from Backup
            </button>
          </div>
        </div>
      </div>

      <div class="bg-dark-800 rounded-xl border border-dark-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4">Backup History</h2>
        <div class="space-y-2">
          <div
            v-for="backup in backupHistory"
            :key="backup.id"
            class="flex items-center justify-between p-3 bg-dark-700/50 rounded-lg"
          >
            <div class="flex items-center gap-3">
              <ArchiveBoxIcon class="w-5 h-5 text-dark-400" />
              <div>
                <p class="text-white">{{ backup.name }}</p>
                <p class="text-dark-400 text-sm">{{ backup.date }} | {{ backup.size }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="downloadBackup(backup)"
                class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              >
                <ArrowDownTrayIcon class="w-4 h-4 text-dark-400" />
              </button>
              <button
                @click="deleteBackup(backup)"
                class="p-2 hover:bg-dark-600 rounded-lg transition-colors"
              >
                <TrashIcon class="w-4 h-4 text-red-400" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

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
import { ref, reactive, onMounted } from 'vue'
import { settingsApi, pluginsApi } from '@/api'
import {
  ServerIcon,
  DatabaseIcon,
  ShieldCheckIcon,
  KeyIcon,
  BellIcon,
  PuzzlePieceIcon,
  ArchiveBoxIcon,
  CheckIcon,
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  ArrowPathIcon,
  PlayIcon,
  PauseIcon,
  CogIcon,
  ArrowDownTrayIcon,
  TrashIcon,
  ArrowUturnLeftIcon,
  ExclamationCircleIcon,
  CircleStackIcon,
  CommandLineIcon,
  ComputerDesktopIcon,
  UserGroupIcon
} from '@heroicons/vue/24/outline'

// State
const activeTab = ref('server')
const saving = ref(false)
const toasts = ref([])

const tabs = [
  { id: 'server', name: 'Server', icon: ServerIcon },
  { id: 'security', name: 'Security', icon: ShieldCheckIcon },
  { id: 'notifications', name: 'Notifications', icon: BellIcon },
  { id: 'plugins', name: 'Plugins', icon: PuzzlePieceIcon },
  { id: 'backup', name: 'Backup', icon: ArchiveBoxIcon }
]

const settings = reactive({
  server: {
    host: '0.0.0.0',
    port: 8080,
    sslMode: 'auto',
    maxConnections: 1000,
    trustedProxies: ''
  },
  database: {
    type: 'sqlite',
    host: 'localhost',
    port: 3306,
    name: 'anixops',
    maxConnections: 10
  },
  security: {
    twoFactorEnabled: false,
    sessionTimeout: 30,
    ipWhitelistEnabled: false,
    allowedIPs: '',
    minPasswordLength: 8,
    passwordExpiry: 0,
    requireUppercase: true,
    requireLowercase: true,
    requireNumber: true,
    requireSpecialChar: false
  },
  notifications: {
    smtpHost: '',
    smtpPort: 587,
    smtpUser: '',
    smtpPassword: '',
    fromEmail: ''
  },
  backup: {
    autoBackup: true,
    frequency: 'daily',
    retention: 30
  }
})

const alertRules = ref([
  { id: 1, name: 'Node Offline', description: 'Alert when a node goes offline', enabled: true, icon: ExclamationCircleIcon, color: 'text-red-400' },
  { id: 2, name: 'High CPU Usage', description: 'Alert when CPU usage exceeds 90%', enabled: true, icon: ExclamationTriangleIcon, color: 'text-yellow-400' },
  { id: 3, name: 'High Memory Usage', description: 'Alert when memory usage exceeds 85%', enabled: true, icon: ExclamationTriangleIcon, color: 'text-yellow-400' },
  { id: 4, name: 'SSL Certificate Expiry', description: 'Alert when SSL cert is about to expire', enabled: true, icon: ShieldCheckIcon, color: 'text-blue-400' },
  { id: 5, name: 'Failed Login Attempts', description: 'Alert on multiple failed logins', enabled: false, icon: UserGroupIcon, color: 'text-purple-400' }
])

const plugins = ref([
  { name: 'Ansible', description: 'Automated deployment and configuration', enabled: true, version: '1.0.0', icon: CommandLineIcon },
  { name: 'V2Board', description: 'V2Ray panel integration', enabled: true, version: '1.0.0', icon: CircleStackIcon },
  { name: 'V2BX', description: 'V2Board extended integration', enabled: true, version: '1.0.0', icon: CircleStackIcon },
  { name: 'Agent', description: 'Remote agent management', enabled: true, version: '1.0.0', icon: ComputerDesktopIcon }
])

const backupHistory = ref([
  { id: 1, name: 'backup-2024-03-15.sql', date: '2024-03-15 00:00', size: '2.4 MB' },
  { id: 2, name: 'backup-2024-03-14.sql', date: '2024-03-14 00:00', size: '2.3 MB' },
  { id: 3, name: 'backup-2024-03-13.sql', date: '2024-03-13 00:00', size: '2.2 MB' }
])

// Methods
async function saveAllSettings() {
  saving.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    showToast('Settings saved successfully', 'success')
  } catch (error) {
    showToast(error.message || 'Failed to save settings', 'error')
  } finally {
    saving.value = false
  }
}

async function testEmail() {
  try {
    showToast('Test email sent', 'success')
  } catch (error) {
    showToast('Failed to send test email', 'error')
  }
}

async function togglePlugin(plugin) {
  try {
    plugin.enabled = !plugin.enabled
    showToast(`${plugin.name} ${plugin.enabled ? 'enabled' : 'disabled'}`, 'success')
  } catch (error) {
    showToast(`Failed to toggle ${plugin.name}`, 'error')
  }
}

async function restartPlugin(plugin) {
  try {
    showToast(`Restarting ${plugin.name}...`, 'success')
  } catch (error) {
    showToast(`Failed to restart ${plugin.name}`, 'error')
  }
}

function openPluginSettings(plugin) {
  console.log('Open settings for:', plugin.name)
}

async function createBackup() {
  try {
    showToast('Creating backup...', 'success')
    backupHistory.value.unshift({
      id: Date.now(),
      name: `backup-${new Date().toISOString().split('T')[0]}.sql`,
      date: new Date().toLocaleString(),
      size: '0 MB'
    })
  } catch (error) {
    showToast('Failed to create backup', 'error')
  }
}

function restoreBackup() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.sql,.zip'
  input.onchange = (e) => {
    const file = e.target.files[0]
    if (file) {
      showToast(`Restoring from ${file.name}...`, 'success')
    }
  }
  input.click()
}

function downloadBackup(backup) {
  showToast(`Downloading ${backup.name}...`, 'success')
}

function deleteBackup(backup) {
  if (confirm(`Delete ${backup.name}?`)) {
    backupHistory.value = backupHistory.value.filter(b => b.id !== backup.id)
    showToast('Backup deleted', 'success')
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
  try {
    const response = await settingsApi.get()
    Object.assign(settings, response.data)
  } catch (error) {
    // Use default settings
  }
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