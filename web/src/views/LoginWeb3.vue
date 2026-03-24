<template>
  <div class="min-h-screen bg-slate-900 flex items-center justify-center p-6">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
          <svg class="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-white">AnixOps Control Center</h1>
        <p class="text-slate-400 mt-2">Connect with Web3</p>
      </div>

      <!-- Login Card -->
      <div class="bg-slate-800 rounded-2xl border border-slate-700 p-8">
        <!-- Tabs -->
        <div class="flex mb-6 bg-slate-700 rounded-lg p-1">
          <button
            @click="loginMethod = 'email'"
            :class="['flex-1 py-2 text-sm font-medium rounded-md transition-colors', loginMethod === 'email' ? 'bg-blue-600 text-white' : 'text-slate-400']"
          >
            Email
          </button>
          <button
            @click="loginMethod = 'web3'"
            :class="['flex-1 py-2 text-sm font-medium rounded-md transition-colors', loginMethod === 'web3' ? 'bg-blue-600 text-white' : 'text-slate-400']"
          >
            Web3 / MetaMask
          </button>
        </div>

        <!-- Email Login -->
        <form v-if="loginMethod === 'email'" @submit.prevent="handleEmailLogin">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Email</label>
              <input
                v-model="email"
                type="email"
                required
                class="w-full bg-slate-700 border border-slate-600 rounded-lg px-4 py-3 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="admin@anixops.com"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Password</label>
              <input
                v-model="password"
                type="password"
                required
                class="w-full bg-slate-700 border border-slate-600 rounded-lg px-4 py-3 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter password"
              />
            </div>

            <div v-if="error" class="p-3 bg-red-900/50 border border-red-800 rounded-lg">
              <p class="text-red-300 text-sm">{{ error }}</p>
            </div>

            <button
              type="submit"
              :disabled="loading"
              class="w-full py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              {{ loading ? 'Signing in...' : 'Sign In' }}
            </button>
          </div>
        </form>

        <!-- Web3 Login -->
        <div v-if="loginMethod === 'web3'" class="space-y-4">
          <div class="text-center py-6">
            <div class="w-20 h-20 bg-gradient-to-br from-orange-500 to-orange-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
              <svg class="w-12 h-12 text-white" viewBox="0 0 40 40" fill="currentColor">
                <path fill-rule="evenodd" d="M20 0C8.954 0 0 8.954 0 20s8.954 20 20 20 20-8.954 20-20S31.046 0 20 0zm-4.5 29.5l-6-9 1.5-4.5 4.5 1.5 4.5-4.5 4.5 4.5 4.5-1.5 1.5 4.5-6 9h-9z" clip-rule="evenodd"/>
              </svg>
            </div>
            <p class="text-slate-300 text-sm">Connect your Ethereum wallet to authenticate using SIWE (Sign-In with Ethereum)</p>
          </div>

          <div v-if="web3Store.isConnected" class="p-4 bg-slate-700 rounded-lg">
            <p class="text-sm text-slate-400">Connected Wallet</p>
            <p class="text-white font-mono text-sm truncate">{{ web3Store.walletAddress }}</p>
            <p class="text-blue-400 text-xs mt-1">{{ web3Store.did }}</p>
          </div>

          <div v-if="web3Store.error" class="p-3 bg-red-900/50 border border-red-800 rounded-lg">
            <p class="text-red-300 text-sm">{{ web3Store.error }}</p>
          </div>

          <button
            @click="handleWeb3Login"
            :disabled="web3Store.isConnecting"
            class="w-full py-3 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-lg font-medium hover:from-orange-600 hover:to-orange-700 disabled:opacity-50 transition-all flex items-center justify-center gap-2"
          >
            <svg v-if="web3Store.isConnecting" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ web3Store.isConnecting ? 'Connecting...' : web3Store.isConnected ? 'Sign to Login' : 'Connect MetaMask' }}
          </button>

          <p class="text-xs text-slate-500 text-center">
            Requires MetaMask or compatible Web3 wallet
          </p>
        </div>
      </div>

      <!-- Footer -->
      <p class="text-center text-slate-500 text-sm mt-6">
        Powered by Cloudflare Workers AI & Vectorize
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useWeb3Store } from '@/stores/web3'

const router = useRouter()
const authStore = useAuthStore()
const web3Store = useWeb3Store()

const loginMethod = ref('email')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleEmailLogin = async () => {
  loading.value = true
  error.value = ''

  const result = await authStore.login(email.value, password.value)

  if (result.success) {
    router.push('/')
  } else {
    error.value = result.error || 'Login failed'
  }

  loading.value = false
}

const handleWeb3Login = async () => {
  const result = await web3Store.web3Login()

  if (result.success) {
    // Create a session with Web3 authentication
    router.push('/')
  }
}
</script>