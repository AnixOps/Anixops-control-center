<template>
  <div class="h-full overflow-y-auto p-6">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-2xl font-bold text-white">Web3 Dashboard</h1>
          <p class="text-slate-400 mt-1">IPFS & Blockchain Integration</p>
        </div>
        <div class="flex items-center gap-3">
          <div v-if="web3Store.isConnected" class="flex items-center gap-2 px-4 py-2 bg-green-900/50 border border-green-700 rounded-lg">
            <div class="w-2 h-2 bg-green-500 rounded-full"></div>
            <span class="text-green-400 text-sm">Connected</span>
          </div>
          <button
            v-else
            @click="web3Store.connectWallet"
            class="px-4 py-2 bg-orange-600 text-white rounded-lg text-sm font-medium hover:bg-orange-700 transition-colors"
          >
            Connect Wallet
          </button>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-900/50 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
              </svg>
            </div>
            <div>
              <p class="text-slate-400 text-sm">IPFS Files</p>
              <p class="text-2xl font-bold text-white">{{ stats.ipfsFiles }}</p>
            </div>
          </div>
        </div>

        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-purple-900/50 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <div>
              <p class="text-slate-400 text-sm">On-Chain Audits</p>
              <p class="text-2xl font-bold text-white">{{ stats.onChainAudits }}</p>
            </div>
          </div>
        </div>

        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-orange-900/50 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-orange-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" />
              </svg>
            </div>
            <div>
              <p class="text-slate-400 text-sm">DID</p>
              <p class="text-sm font-mono text-white truncate">{{ web3Store.did || 'Not set' }}</p>
            </div>
          </div>
        </div>

        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-green-900/50 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-slate-400 text-sm">Network</p>
              <p class="text-sm font-medium text-white">{{ network }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- IPFS Upload -->
        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <h2 class="text-lg font-semibold text-white mb-4">IPFS Storage</h2>

          <div class="border-2 border-dashed border-slate-600 rounded-lg p-8 text-center mb-4">
            <input
              type="file"
              ref="fileInput"
              @change="handleFileSelect"
              class="hidden"
            />
            <button
              @click="$refs.fileInput.click()"
              class="px-6 py-3 bg-slate-700 text-white rounded-lg hover:bg-slate-600 transition-colors"
            >
              Select File to Upload
            </button>
            <p class="text-slate-400 text-sm mt-2">Files are stored on IPFS decentralized network</p>
          </div>

          <div v-if="selectedFile" class="p-3 bg-slate-700 rounded-lg mb-4">
            <p class="text-white text-sm">{{ selectedFile.name }}</p>
            <p class="text-slate-400 text-xs">{{ formatBytes(selectedFile.size) }}</p>
          </div>

          <button
            v-if="selectedFile"
            @click="uploadFile"
            :disabled="uploading"
            class="w-full py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            {{ uploading ? 'Uploading...' : 'Upload to IPFS' }}
          </button>

          <div v-if="ipfsResult" class="mt-4 p-4 bg-green-900/30 border border-green-800 rounded-lg">
            <p class="text-green-400 text-sm font-medium">Uploaded Successfully!</p>
            <p class="text-white font-mono text-sm mt-1 break-all">CID: {{ ipfsResult }}</p>
            <a
              :href="`https://cloudflare-ipfs.com/ipfs/${ipfsResult}`"
              target="_blank"
              class="text-blue-400 text-sm hover:underline mt-2 inline-block"
            >
              View on IPFS Gateway
            </a>
          </div>
        </div>

        <!-- On-Chain Audit -->
        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <h2 class="text-lg font-semibold text-white mb-4">On-Chain Audit Log</h2>

          <div class="space-y-3 mb-4">
            <div>
              <label class="block text-sm text-slate-400 mb-1">Action</label>
              <select v-model="auditForm.action" class="w-full bg-slate-700 border border-slate-600 rounded-lg px-4 py-2 text-white">
                <option value="node.restart">Node Restart</option>
                <option value="node.create">Node Create</option>
                <option value="node.delete">Node Delete</option>
                <option value="playbook.run">Playbook Run</option>
                <option value="user.login">User Login</option>
                <option value="settings.change">Settings Change</option>
              </select>
            </div>
            <div>
              <label class="block text-sm text-slate-400 mb-1">Details</label>
              <textarea
                v-model="auditForm.details"
                rows="3"
                class="w-full bg-slate-700 border border-slate-600 rounded-lg px-4 py-2 text-white"
                placeholder="Audit details..."
              ></textarea>
            </div>
          </div>

          <button
            @click="storeAudit"
            :disabled="!web3Store.isConnected || storing"
            class="w-full py-3 bg-purple-600 text-white rounded-lg font-medium hover:bg-purple-700 disabled:opacity-50 transition-colors"
          >
            {{ storing ? 'Recording...' : 'Record on Blockchain' }}
          </button>

          <div v-if="auditResult" class="mt-4 p-4 bg-green-900/30 border border-green-800 rounded-lg">
            <p class="text-green-400 text-sm font-medium">Audit Recorded!</p>
            <p class="text-white font-mono text-xs mt-1 break-all">Tx: {{ auditResult.txHash }}</p>
            <p class="text-slate-300 text-xs mt-1">IPFS: {{ auditResult.ipfsCid }}</p>
          </div>
        </div>
      </div>

      <!-- Recent Audits -->
      <div class="mt-6 bg-slate-800 rounded-xl border border-slate-700 p-6">
        <h2 class="text-lg font-semibold text-white mb-4">Recent Audit Records</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-slate-700">
                <th class="text-left text-slate-400 text-sm font-medium py-3">Action</th>
                <th class="text-left text-slate-400 text-sm font-medium py-3">Details</th>
                <th class="text-left text-slate-400 text-sm font-medium py-3">Timestamp</th>
                <th class="text-left text-slate-400 text-sm font-medium py-3">Tx Hash</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="audit in recentAudits" :key="audit.id" class="border-b border-slate-700/50">
                <td class="py-3 text-white text-sm">{{ audit.action }}</td>
                <td class="py-3 text-slate-300 text-sm">{{ audit.details }}</td>
                <td class="py-3 text-slate-400 text-sm">{{ formatTime(audit.timestamp) }}</td>
                <td class="py-3 text-slate-300 font-mono text-xs">{{ audit.txHash?.slice(0, 20) }}...</td>
              </tr>
              <tr v-if="recentAudits.length === 0">
                <td colspan="4" class="py-8 text-center text-slate-500">No audit records yet</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useWeb3Store } from '@/stores/web3'

const web3Store = useWeb3Store()

const stats = reactive({
  ipfsFiles: 0,
  onChainAudits: 0
})

const network = ref('Ethereum Mainnet')
const selectedFile = ref(null)
const uploading = ref(false)
const ipfsResult = ref(null)
const storing = ref(false)
const auditResult = ref(null)
const recentAudits = ref([])

const auditForm = reactive({
  action: 'node.restart',
  details: ''
})

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
  }
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const uploadFile = async () => {
  if (!selectedFile.value) return

  uploading.value = true
  ipfsResult.value = null

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    const cid = await web3Store.uploadToIPFS(formData)
    if (cid) {
      ipfsResult.value = cid
      stats.ipfsFiles++
      selectedFile.value = null
    }
  } catch (err) {
    console.error('Upload failed:', err)
  } finally {
    uploading.value = false
  }
}

const storeAudit = async () => {
  if (!web3Store.isConnected) return

  storing.value = true
  auditResult.value = null

  try {
    const result = await web3Store.storeAudit({
      action: auditForm.action,
      details: auditForm.details
    })

    if (result) {
      auditResult.value = result
      stats.onChainAudits++
      recentAudits.value.unshift({
        id: Date.now(),
        action: auditForm.action,
        details: auditForm.details,
        timestamp: new Date().toISOString(),
        txHash: result.txHash
      })
      auditForm.details = ''
    }
  } catch (err) {
    console.error('Audit failed:', err)
  } finally {
    storing.value = false
  }
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

onMounted(() => {
  // Load initial data
})
</script>