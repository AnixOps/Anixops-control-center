<template>
  <div class="min-h-screen bg-slate-900 flex items-center justify-center p-6">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 rounded-2xl bg-blue-600 flex items-center justify-center mx-auto">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-white mt-4">AnixOps Control Center</h1>
        <p class="text-slate-400 mt-2">Sign in to your account</p>
      </div>

      <!-- Login Form -->
      <div class="bg-slate-800 rounded-2xl border border-slate-700 p-8">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-sm text-slate-400 mb-2">Email</label>
            <input
              v-model="email"
              type="email"
              required
              placeholder="admin@example.com"
              class="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label class="block text-sm text-slate-400 mb-2">Password</label>
            <input
              v-model="password"
              type="password"
              required
              placeholder="••••••••"
              class="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="flex items-center justify-between">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="rememberMe"
                type="checkbox"
                class="rounded border-slate-500 bg-slate-600 text-blue-500 focus:ring-blue-500"
              />
              <span class="text-slate-400 text-sm">Remember me</span>
            </label>
            <button type="button" class="text-blue-400 text-sm hover:text-blue-300">
              Forgot password?
            </button>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <svg v-if="loading" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
          </button>
        </form>

        <div class="mt-6 pt-6 border-t border-slate-700">
          <button
            type="button"
            class="w-full py-3 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132c.645-2.266.99-4.659.99-7.132A8 8 0 008 4.07M3 15.364c.64-1.319 1-2.8 1-4.364 0-1.457.39-2.823 1.07-4" />
            </svg>
            Sign in with Biometrics
          </button>
        </div>
      </div>

      <!-- Server Settings -->
      <div class="mt-4 text-center">
        <button
          @click="showServerSettings = true"
          class="text-slate-400 text-sm hover:text-white transition-colors"
        >
          Server Settings
        </button>
      </div>

      <!-- Server Settings Modal -->
      <div v-if="showServerSettings" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
        <div class="bg-slate-800 rounded-2xl border border-slate-700 w-full max-w-md p-6">
          <h3 class="text-lg font-semibold text-white mb-4">Server Settings</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm text-slate-400 mb-1">API URL</label>
              <input
                v-model="apiUrl"
                type="text"
                placeholder="http://localhost:8080/api/v1"
                class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <button
              @click="showServerSettings = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              @click="showServerSettings = false"
              class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const loading = ref(false)
const showServerSettings = ref(false)
const apiUrl = ref('http://localhost:8080/api/v1')

async function handleLogin() {
  loading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    router.push('/')
  } finally {
    loading.value = false
  }
}
</script>